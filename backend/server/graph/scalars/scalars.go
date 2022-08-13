package scalars

import (
	"errors"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(t.Format("\"3:04pm\"")))
	})
}

func UnmarshalTime(v interface{}) (time.Time, error) {
	return time.Now(), errors.New("Unmarshalling 'Time' not implemented")
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

	location, _ := time.LoadLocation("America/Montreal")

	t, err := time.ParseInLocation("2006-01-02", date, location)
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

	location, _ := time.LoadLocation("America/Montreal")

	t, err := time.Parse("2006-01-02T15:04:00Z", date)
	if err != nil {
		return time.Time{}, err
	}

	return t.In(location), nil
}
