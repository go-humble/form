// Copyright 2015 Alex Browne.
// Allrights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package form

import (
	"fmt"

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
}

func Parse(formElement *dom.HTMLFormElement) (*Form, error) {
	form := &Form{
		Inputs: map[string]*Input{},
	}
	for _, el := range formElement.Elements() {
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

func (form *Form) GetString(name string) (string, error) {
	input, found := form.Inputs[name]
	if !found {
		return "", NewInputNotFoundError(name)
	}
	return input.String(), nil
}

func (form *Form) GetInt(name string) (int, error) {
	input, found := form.Inputs[name]
	if !found {
		return 0, NewInputNotFoundError(name)
	}
	return input.Int()
}

func (form *Form) GetFloat(name string) (float64, error) {
	input, found := form.Inputs[name]
	if !found {
		return 0, NewInputNotFoundError(name)
	}
	return input.Float()
}

func (form *Form) GetBool(name string) (bool, error) {
	input, found := form.Inputs[name]
	if !found {
		return false, NewInputNotFoundError(name)
	}
	return input.Bool()
}

// func (form *Form) GetTime(name string) (time.Time, error) {
// 	input, found := form.Inputs[name]
// 	if !found {
// 		return time.Time{}, NewInputNotFoundError(name)
// 	}
// 	return input.Time()
// }
