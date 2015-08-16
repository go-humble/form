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

type InputNotFoundError struct {
	Name string
}

func NewInputNotFoundError(name string) InputNotFoundError {
	return InputNotFoundError{
		Name: name,
	}
}

func (e InputNotFoundError) Error() string {
	return fmt.Sprintf("form: could not find input with name = %s", e.Name)
}

type Form struct {
	Inputs map[string]*Input
	Errors []error
}

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

func (form *Form) GetString(inputName string) (string, error) {
	input, found := form.Inputs[inputName]
	if !found {
		return "", NewInputNotFoundError(inputName)
	}
	return input.RawValue, nil
}

func (form *Form) GetInt(inputName string) (int, error) {
	input, found := form.Inputs[inputName]
	if !found {
		return 0, NewInputNotFoundError(inputName)
	}
	return input.Int()
}

func (form *Form) GetUint(inputName string) (uint, error) {
	input, found := form.Inputs[inputName]
	if !found {
		return 0, NewInputNotFoundError(inputName)
	}
	return input.Uint()
}

func (form *Form) GetFloat(inputName string) (float64, error) {
	input, found := form.Inputs[inputName]
	if !found {
		return 0, NewInputNotFoundError(inputName)
	}
	return input.Float()
}

func (form *Form) GetBool(inputName string) (bool, error) {
	input, found := form.Inputs[inputName]
	if !found {
		return false, NewInputNotFoundError(inputName)
	}
	return input.Bool()
}

func (form *Form) GetTime(inputName string) (time.Time, error) {
	input, found := form.Inputs[inputName]
	if !found {
		return time.Time{}, NewInputNotFoundError(inputName)
	}
	return input.Time()
}

func (form *Form) HasErrors() bool {
	return len(form.Errors) > 0
}
