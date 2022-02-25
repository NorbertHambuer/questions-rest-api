package entities

import (
	"testing"
)

func TestValidateOption(t *testing.T) {
	testCases := []struct {
		name    string
		input   Option
		isError bool
	}{
		{
			name: "correct option",
			input: Option{
				Id:          0,
				QuestionId:  0,
				Body:        "is this user correct?",
				Correct:     true,
				OptionOrder: 0,
			},
			isError: false,
		},
		{
			name:    "empty option",
			input:   Option{},
			isError: true,
		},
		{
			name: "negative questionId",
			input: Option{
				Id:          0,
				QuestionId:  -1,
				Body:        "correct body",
				Correct:     false,
				OptionOrder: 0,
			},
			isError: true,
		},
		{
			name: "invalid body length",
			input: Option{
				Id:          0,
				QuestionId:  0,
				Body:        "",
				Correct:     false,
				OptionOrder: 0,
			},
			isError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()

			if (err != nil) != tc.isError {
				t.Errorf("expected error (%v), got error (%v)", tc.isError, err.Error())
			}
		})
	}
}
