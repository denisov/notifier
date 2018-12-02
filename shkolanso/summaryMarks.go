package shkolanso

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

type mark struct {
	// можно ли сразу парсить дату ?
	Date shkolaTime `json:"date"`
	Mark string     `json:"mark"`
}

type disciplineMark struct {
	Discipline string `json:"discipline"`
	Marks      []mark `json:"marks"`
}

type summaryMarks struct {
	DisciplineMarks []disciplineMark `json:"discipline_marks"`
}

type shkolaTime struct {
	time.Time
}

func (st *shkolaTime) UnmarshalJSON(p []byte) error {
	loc, err := time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		return errors.Wrap(err, "Не могу найти location")
	}
	t, err := time.ParseInLocation(
		"2006-01-02",
		strings.Replace(string(p), "\"", "", -1),
		loc,
	)

	if err != nil {
		return err
	}

	st.Time = t
	return nil
}
