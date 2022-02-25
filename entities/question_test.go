package entities

import "testing"

func TestValidateQuestion(t *testing.T) {
	testCases := []struct {
		name    string
		input   Question
		isError bool
	}{
		{
			name: "correct question",
			input: Question{
				Id:   0,
				Body: "Where does the sun set?",
				Options: []Option{{
					Id:          0,
					QuestionId:  0,
					Body:        "East",
					Correct:     false,
					OptionOrder: 0,
				}, {
					Id:          0,
					QuestionId:  0,
					Body:        "West",
					Correct:     true,
					OptionOrder: 0,
				}},
			},
			isError: false,
		},
		{
			name:    "empty question",
			input:   Question{},
			isError: true,
		},
		{
			name: "empty options",
			input: Question{
				Id:      0,
				Body:    "Where does the sun rise?",
				Options: nil,
			},
			isError: true,
		},
		{
			name: "invalid body length",
			input: Question{
				Id:   0,
				Body: "1",
				Options: []Option{{
					Id:          0,
					QuestionId:  0,
					Body:        "false",
					Correct:     false,
					OptionOrder: 0,
				}, {
					Id:          0,
					QuestionId:  0,
					Body:        "true",
					Correct:     true,
					OptionOrder: 0,
				}},
			},
			isError: true,
		},
		{
			name: "no correct options",
			input: Question{
				Id:   0,
				Body: "Where does the sun set?",
				Options: []Option{{
					Id:          0,
					QuestionId:  0,
					Body:        "East",
					Correct:     false,
					OptionOrder: 0,
				}, {
					Id:          0,
					QuestionId:  0,
					Body:        "West",
					Correct:     false,
					OptionOrder: 0,
				}},
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
