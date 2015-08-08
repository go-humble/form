package form

import (
	"strconv"

	"honnef.co/go/js/dom"
)

type Input struct {
	El       *dom.HTMLInputElement
	Name     string
	RawValue string
	Type     InputType
}

func NewInput(el *dom.HTMLInputElement) *Input {
	return &Input{
		El:       el,
		Name:     el.Name,
		RawValue: el.Value,
		Type:     InputType(el.Type),
	}
}

func (input Input) String() string {
	return input.RawValue
}

func (input Input) Int() (int, error) {
	return strconv.Atoi(input.RawValue)
}

func (input Input) Float() (float64, error) {
	return strconv.ParseFloat(input.RawValue, 64)
}

func (input Input) Bool() (bool, error) {
	switch input.Type {
	case InputCheckbox, InputRadio:
		return input.El.Checked, nil
	default:
		return strconv.ParseBool(input.RawValue)
	}
}
