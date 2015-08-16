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

// Input is a go representation of an html input element.
type Input struct {
	// El is the original html input element for the Input.
	El *dom.HTMLInputElement
	// Name is equal to the value of the input's name attribute.
	Name string
	// RawValue is equal to the input's value attribute.
	RawValue string
	// Type is equal to the input's type attribute. Note that this is sometimes
	// different than the type reported by type property of the HTMLInputElement
	// in the DOM API.
	Type InputType
}

// NewInput creates a new Input object from the given html input element.
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

// Int converts the value of the input to an int. It returns an error if the
// value could not be converted.
func (input Input) Int() (int, error) {
	return strconv.Atoi(input.RawValue)
}

// Uint converts the value of the input to a uint. It returns an error if the
// value could not be converted.
func (input Input) Uint() (uint, error) {
	u, err := strconv.ParseUint(input.RawValue, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(u), nil
}

// Float converts the value of the input to a float64. It returns an error if the
// value could not be converted.
func (input Input) Float() (float64, error) {
	return strconv.ParseFloat(input.RawValue, 64)
}

// Bool converts the value of the input to a bool. For inputs with the type
// checkbox or radio, Bool will return true iff the input has the checked
// attribute. For all other input types it will attempt to parse the input value
// as a bool. It returns an error if the value could not be converted.
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

// Time converts the value of the input to a time.Time. Time supports the time,
// date, and datetime input types and assumes the input value adheres to the
// rfc3339 standard (the default for form inputs). If the type of the input is
// anything else, it will attempt to parse it as an rfc3339 datetime. It returns
// an error if the value could not be converted.
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
