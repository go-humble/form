// Copyright 2015 Alex Browne.
// Allrights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package form

import (
	"fmt"
	"time"

	"honnef.co/go/js/dom"
)

// InputNotFoundError is returned whenever form.GetX is called with an
// inputName that is not found in the form inputs.
type InputNotFoundError struct {
	Name string
}

// newInputNotFoundError creates and returns an InputNotFoundError with the
// given name.
func newInputNotFoundError(name string) InputNotFoundError {
	return InputNotFoundError{
		Name: name,
	}
}

// Error satisifies the Error method of the error interface.
func (e InputNotFoundError) Error() string {
	return fmt.Sprintf("form: could not find input with name = %s", e.Name)
}

// Form is a go representation of an html input form.
type Form struct {
	Inputs map[string]*Input
	Errors []error
}

// Parse creates a and returns a Form object from the given form element.
func Parse(formElement dom.Element) (*Form, error) {
	form := &Form{
		Inputs: map[string]*Input{},
	}
	htmlFormElement, ok := formElement.(*dom.HTMLFormElement)
	if !ok {
		return nil, fmt.Errorf("form: Argument to Parse must be a *dom.HTMLFormElement. (Got %T)", formElement)
	}
	for _, el := range htmlFormElement.Elements() {
		// Cast the element to an input element.
		inputEl, ok := el.(*dom.HTMLInputElement)
		if !ok {
			// Skip elements which are not input elements.
			continue
		}
		form.Inputs[inputEl.Name] = NewInput(inputEl)
	}
	return form, nil
}

// GetString returns the value of the input identified by inputName. It returns
// an InputNotFoundError if there is no input with the given inputName.
func (form *Form) GetString(inputName string) (string, error) {
	input, found := form.Inputs[inputName]
	if !found {
		return "", newInputNotFoundError(inputName)
	}
	return input.RawValue, nil
}

// GetInt returns the value of the input identified by inputName converted
// to an int. It returns an error if the input is not found or if the input
// value could not be converted to an int.
func (form *Form) GetInt(inputName string) (int, error) {
	input, found := form.Inputs[inputName]
	if !found {
		return 0, newInputNotFoundError(inputName)
	}
	return input.Int()
}

// GetUint returns the value of the input identified by inputName converted
// to a uint. It returns an error if the input is not found or if the input
// value could not be converted to a uint.
func (form *Form) GetUint(inputName string) (uint, error) {
	input, found := form.Inputs[inputName]
	if !found {
		return 0, newInputNotFoundError(inputName)
	}
	return input.Uint()
}

// GetFloat returns the value of the input identified by inputName converted
// to a float. It returns an error if the input is not found or if the input
// value could not be converted to a float.
func (form *Form) GetFloat(inputName string) (float64, error) {
	input, found := form.Inputs[inputName]
	if !found {
		return 0, newInputNotFoundError(inputName)
	}
	return input.Float()
}

// GetBool returns the value of the input identified by inputName converted
// to a bool. It returns an error if the input is not found or if the input
// value could not be converted to a bool.
func (form *Form) GetBool(inputName string) (bool, error) {
	input, found := form.Inputs[inputName]
	if !found {
		return false, newInputNotFoundError(inputName)
	}
	return input.Bool()
}

// GetTime returns the value of the input identified by inputName converted to a
// time.Time. GetTime supports the time, date, and datetime input types and
// assumes the input value adheres to the rfc3339 standard (the default for form
// inputs). If the type of the input is anything else, it will attempt to parse
// it as an rfc3339 datetime. It returns an error if the input is not found or
// if the input value could not be converted to a time.Time.
func (form *Form) GetTime(inputName string) (time.Time, error) {
	input, found := form.Inputs[inputName]
	if !found {
		return time.Time{}, newInputNotFoundError(inputName)
	}
	return input.Time()
}

// HasErrors returns true if the form has at least one validation error.
func (form *Form) HasErrors() bool {
	return len(form.Errors) > 0
}
