package form

import (
	"fmt"
	"strconv"

	"honnef.co/go/js/dom"
)

type StringInput struct {
	value string
	defaultInput
}

func NewStringInput(el *dom.HTMLInputElement) *StringInput {
	return &StringInput{
		value:        el.Value,
		defaultInput: newDefaultInput(el),
	}
}

func (input StringInput) String() (string, error) {
	return input.value, nil
}

func (input StringInput) Int() (int, error) {
	i, err := strconv.Atoi(input.value)
	if err != nil {
		return 0, fmt.Errorf("form: Error parsing %s: %s", input.Name(), err.Error())
	}
	return i, nil
}

func (input StringInput) Float() (float64, error) {
	f, err := strconv.ParseFloat(input.value, 64)
	if err != nil {
		return 0.0, fmt.Errorf("form: Error parsing %s: %s", input.Name(), err.Error())
	}
	return f, nil
}

func (input StringInput) Bool() (bool, error) {
	b, err := strconv.ParseBool(input.value)
	if err != nil {
		return false, fmt.Errorf("form: Error parsing %s: %s", input.Name(), err.Error())
	}
	return b, nil
}
