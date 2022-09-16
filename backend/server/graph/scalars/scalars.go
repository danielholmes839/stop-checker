package scalars

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"stop-checker.com/db/model"
)

func MarshalTime(t model.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		str := t.String()
		w.Write([]byte(fmt.Sprintf("\"%s\"", str)))
	})
}

func UnmarshalTime(v interface{}) (model.Time, error) {
	return model.NewTime(0, 0), errors.New("Unmarshalling 'Time' not implemented")
}

func MarshalDate(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(t.Format("\"2006-01-02\"")))
	})
}

func UnmarshalDate(v interface{}) (time.Time, error) {
	// TODO: fix timezone
	date, ok := v.(string)
	if !ok {
		return time.Time{}, errors.New("'Date' scalar must be a string")
	}

	t, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func MarshalDateTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(t.UTC().Format("\"2006-01-02T15:04:00Z\"")))
	})
}

func UnmarshalDateTime(v interface{}) (time.Time, error) {
	date, ok := v.(string)
	if !ok {
		return time.Time{}, errors.New("'DateTime' scalar must be a string")
	}

	t, err := time.Parse("2006-01-02T15:04:00Z", date)

	if err != nil {
		return time.Time{}, err
	}

	return t.In(time.Local), nil
}
