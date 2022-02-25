package entities

import "github.com/go-playground/validator"

// Option defines the structure for the option object
// swagger: model
type Option struct {
	// the id for this option
	//
	// required: true
	Id int64 `json:"-"`
	// question foreign key
	//
	// required: true
	QuestionId int64 `json:"-" validate:"gte=0"`
	// body content for this option, representing one possible answer for a question
	//
	// required: true
	// min: 1
	Body string `json:"body" validate:"required,min=1"`
	// boolean that represents if this options is the correct one
	//
	// required: true
	Correct bool `json:"correct"`
	// integer that represents this option position inside the rest of the questions options
	//
	// required: true
	// min: 0
	OptionOrder int `json:"-" validate:"gte=0"`
}

// Validate checks and validates each field of the option object based on its definition
func (o *Option) Validate() error {
	validate := validator.New()

	return validate.Struct(o)
}
