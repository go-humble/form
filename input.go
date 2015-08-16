// Copyright 2015 Alex Browne and Soroush Pour.
// Allrights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package form

import (
	"strconv"
	"time"

	"honnef.co/go/js/dom"
)

type Input struct {
	El       *dom.HTMLInputElement
	Name     string
	RawValue string
	Type     InputType
}

func NewInput(el *dom.HTMLInputElement) *Input {
	// Attempt to determine the type by first getting the type attribute
	// directly. This is more reliable as some browsers will always return
	// "text" as the type if the type attribute is not recognized/supported.
	// If the type attribute is missing, fallback to using what the browser
	// thinks the type is (which is probably "text").
	inputType := InputType(el.GetAttribute("type"))
	if inputType == "" {
		inputType = InputType(el.Type)
	}
	return &Input{
		El:       el,
		Name:     el.Name,
		RawValue: el.Value,
		Type:     inputType,
	}
}

func (input Input) Int() (int, error) {
	return strconv.Atoi(input.RawValue)
}

func (input Input) Uint() (uint, error) {
	u, err := strconv.ParseUint(input.RawValue, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(u), nil
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

const (
	rfc3339DateLayout          = "2006-01-02"
	rfc3339DatetimeLocalLayout = "2006-01-02T15:04:05.999999999"
)

func (input Input) Time() (time.Time, error) {
	switch input.Type {
	case InputDate:
		return time.Parse(rfc3339DateLayout, input.RawValue)
	case InputDateTimeLocal:
		return time.Parse(rfc3339DatetimeLocalLayout, input.RawValue)
	default:
		return time.Parse(time.RFC3339, input.RawValue)
	}
}
