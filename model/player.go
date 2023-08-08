package model

import "github.com/pkg/errors"

type Player struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Score   int    `json:"score"`
}

func Validate(p Player) error {
	validationchecks := []struct {
		shouldBeTrue bool
		ErrifNotTrue error
	}{
		// field mandatory checks
		{
			shouldBeTrue: len(p.Name) > 0,
			ErrifNotTrue: errors.Wrap(ErrInvalidPlayer, "name cannot be empty"),
		},
		{
			shouldBeTrue: len(p.Country) > 0,
			ErrifNotTrue: errors.Wrap(ErrInvalidPlayer, "country cannot be empty"),
		},

		// constraints checks
		{
			shouldBeTrue: len(p.Name) <= 15,
			ErrifNotTrue: errors.Wrap(ErrInvalidPlayer, "name should be 15 characters at max"),
		},
		{
			shouldBeTrue: len(p.Country) <= 2,
			ErrifNotTrue: errors.Wrap(ErrInvalidPlayer, "Country code should be of 2 characters"),
		},
	}

	for _, check := range validationchecks {
		if !check.shouldBeTrue {
			return check.ErrifNotTrue
		}
	}

	return nil
}
