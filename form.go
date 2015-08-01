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
	Inputs map[string]Input
}

func Parse(formElement *dom.HTMLFormElement) (*Form, error) {
	form := &Form{
		Inputs: map[string]Input{},
	}
	for _, el := range formElement.Elements() {
		// Cast the element to an input element.
		inputEl, ok := el.(*dom.HTMLInputElement)
		if !ok {
			// Skip elements which are not input elements.
			continue
		}
		// Depending on the type of the form element, we want to handle it
		// differently.
		inputType := InputType(inputEl.Type)
		switch inputType {
		case InputDefault, InputEmail, InputHidden, InputPassword, InputSearch,
			InputTel, InputText, InputURL:
			form.Inputs[inputEl.Name] = NewStringInput(inputEl)
		default:
			// If the type is unknown or not supported, just store it as a string
			// input. If you try to convert it to something that we can't convert
			// it to, form will return a readable error.
			form.Inputs[inputEl.Name] = NewStringInput(inputEl)
		}
	}
	return form, nil
}

func (form *Form) GetString(name string) (string, error) {
	input, found := form.Inputs[name]
	if !found {
		return "", NewInputNotFoundError(name)
	}
	stringer, ok := input.(Stringer)
	if !ok {
		return "", fmt.Errorf("form: Cannot get string for input type %s. (%T does not implement Stringer)",
			input.Type(), input)
	}
	return stringer.String()
}

func (form *Form) GetInt(name string) (int, error) {
	input, found := form.Inputs[name]
	if !found {
		return 0, NewInputNotFoundError(name)
	}
	inter, ok := input.(Inter)
	if !ok {
		return 0, fmt.Errorf("form: Cannot get int for input type %s. (%T does not implement Inter)",
			input.Type(), input)
	}
	return inter.Int()
}

func (form *Form) GetFloat(name string) (float64, error) {
	input, found := form.Inputs[name]
	if !found {
		return 0, NewInputNotFoundError(name)
	}
	floater, ok := input.(Floater)
	if !ok {
		return 0, fmt.Errorf("form: Cannot get float for input type %s. (%T does not implement Floater)",
			input.Type(), input)
	}
	return floater.Float()
}

func (form *Form) GetBool(name string) (bool, error) {
	input, found := form.Inputs[name]
	if !found {
		return false, NewInputNotFoundError(name)
	}
	booler, ok := input.(Booler)
	if !ok {
		return false, fmt.Errorf("form: Cannot get bool for input type %s. (%T does not implement Booler)",
			input.Type(), input)
	}
	return booler.Bool()
}
