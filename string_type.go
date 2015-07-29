package form

import (
	"fmt"
	"strconv"
)

type StringInput struct {
	defaultInput
}

func NewStringInput(name string, rawValue string, typ InputType) *StringInput {
	return &StringInput{
		defaultInput{
			name:     name,
			rawValue: rawValue,
			typ:      typ,
		},
	}
}

func (input StringInput) Int() (int, error) {
	i, err := strconv.Atoi(input.rawValue)
	if err != nil {
		return 0, fmt.Errorf("form: Error parsing %s: %s", input.Name(), err.Error())
	}
	return i, nil
}

func (input StringInput) Float() (float64, error) {
	f, err := strconv.ParseFloat(input.rawValue, 64)
	if err != nil {
		return 0.0, fmt.Errorf("form: Error parsing %s: %s", input.Name(), err.Error())
	}
	return f, nil
}

func (input StringInput) Bool() (bool, error) {
	b, err := strconv.ParseBool(input.rawValue)
	if err != nil {
		return false, fmt.Errorf("form: Error parsing %s: %s", input.Name(), err.Error())
	}
	return b, nil
}
