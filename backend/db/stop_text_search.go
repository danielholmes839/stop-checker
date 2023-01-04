package db

import (
	"regexp"
	"sort"
	"strings"

	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type StopTextResult struct {
	model.Stop
	MatchingTokens int
	MatchingCode   bool
}

type StopTextIndex struct {
	stopsByToken      map[string][]model.Stop
	stopsByCode       repository.InvertedIndex[model.Stop]
	stopRoutes        repository.StopRoutes
	removePunctuation *regexp.Regexp
}

func NewStopTextIndex(
	stopsByCode repository.InvertedIndex[model.Stop],
	stopRoutes repository.StopRoutes,
	stops []model.Stop,
) *StopTextIndex {
	re, _ := regexp.Compile(`[^\w]`)

	index := &StopTextIndex{
		stopsByCode:       stopsByCode,
		stopsByToken:      map[string][]model.Stop{},
		stopRoutes:        stopRoutes,
		removePunctuation: re,
	}

	for _, stop := range stops {
		for _, token := range index.tokenize(stop.Name) {
			if _, ok := index.stopsByToken[token]; !ok {
				index.stopsByToken[token] = []model.Stop{}
			}

			index.stopsByToken[token] = append(index.stopsByToken[token], stop)
		}
	}

	return index
}

func (s *StopTextIndex) Query(text string) []model.Stop {
	tokens := s.tokenize(text)
	resultsMap := map[string]*StopTextResult{}

	for _, token := range tokens {
		if stops, err := s.stopsByCode.Get(token); err == nil {
			for _, stop := range stops {
				if result, tracked := resultsMap[stop.ID()]; !tracked {
					resultsMap[stop.ID()] = &StopTextResult{
						Stop:           stop,
						MatchingTokens: 1,
						MatchingCode:   true,
					}
				} else {
					result.MatchingTokens++
					result.MatchingCode = true
				}
			}
		}

		stops := s.stopsByToken[token]

		for _, stop := range stops {
			if result, tracked := resultsMap[stop.ID()]; !tracked {
				resultsMap[stop.ID()] = &StopTextResult{
					Stop:           stop,
					MatchingTokens: 1,
					MatchingCode:   false,
				}
			} else {
				result.MatchingTokens++
			}
		}
	}

	results := make([]*StopTextResult, len(resultsMap))
	i := 0
	for _, result := range resultsMap {
		results[i] = result
		i++
	}

	sort.Slice(results, func(i, j int) bool {
		ri := results[i]
		rj := results[j]

		if ri.MatchingCode && !rj.MatchingCode {
			return true
		}

		if !ri.MatchingCode && rj.MatchingCode {
			return false
		}

		if ri.MatchingTokens == rj.MatchingTokens {
			riRoutes := s.stopRoutes.Get(ri.ID())
			rjRoutes := s.stopRoutes.Get(rj.ID())
			return len(riRoutes) > len(rjRoutes)
		}

		return ri.MatchingTokens > rj.MatchingTokens
	})

	stops := make([]model.Stop, len(results))
	for i, result := range results {
		stops[i] = result.Stop
	}
	return stops

}

func (s *StopTextIndex) tokenize(text string) []string {
	text = strings.ToLower(text)

	tokens := []string{}

	for _, token := range strings.Split(text, " ") {
		token = s.removePunctuation.ReplaceAllString(token, "")
		if len(token) == 0 {
			continue
		}
		tokens = append(tokens, token)
	}

	return tokens
}
