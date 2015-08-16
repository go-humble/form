// Copyright 2015 Alex Browne and Soroush Pour.
// Allrights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package form

import "fmt"

// InputValidation is an object which has methods for validating an input. Such
// methods always return an InputValidation and are chainable. Whenever an input
// is invalid, a validation error is added to the corresponding Form and can
// be accessed via the form.HasErrors method and the form.Errors property.
type InputValidation struct {
	Errors    []error
	Input     *Input
	Form      *Form
	InputName string
}

// ValidationError is returned whenever an error arises from a validation
// method.
type ValidationError struct {
	Input     *Input
	InputName string
	msg       string
}

// Error satisfies the Error method of the builtin error interface.
func (valErr ValidationError) Error() string {
	return valErr.msg
}

// Validate returns an InputValidation object which can be used to validate a
// single input.
func (form *Form) Validate(inputName string) *InputValidation {
	return &InputValidation{
		Input:     form.Inputs[inputName],
		Form:      form,
		InputName: inputName,
	}
}

// AddError adds a validation error to the form with the given format and args.
// The arguments format and args work exactly like they do in fmt.Sprintf.
func (val *InputValidation) AddError(format string, args ...interface{}) {
	err := fmt.Errorf(format, args...)
	val.Errors = append(val.Errors, err)
	valErr := &ValidationError{
		Input:     val.Input,
		InputName: val.InputName,
		msg:       err.Error(),
	}
	val.Form.Errors = append(val.Form.Errors, valErr)
}

// Required adds a validation error to the form if the input is not included in
// the form or if it is an empty string.
func (val *InputValidation) Required() *InputValidation {
	return val.Requiredf("%s is required.", val.InputName)
}

// Requiredf is like Required but allows you to specify a custom error message.
// The arguments format and args work exactly like they do in fmt.Sprintf.
func (val *InputValidation) Requiredf(format string, args ...interface{}) *InputValidation {
	if val.Input == nil || val.Input.RawValue == "" {
		val.AddError(format, args...)
	}
	return val
}

func lessFunc(limit int) func(int) bool {
	return func(value int) bool {
		return value < limit
	}
}

func lessOrEqualFunc(limit int) func(int) bool {
	return func(value int) bool {
		return value <= limit
	}
}

func greaterFunc(limit int) func(int) bool {
	return func(value int) bool {
		return value > limit
	}
}

func greaterOrEqualFunc(limit int) func(int) bool {
	return func(value int) bool {
		return value >= limit
	}
}

// Less adds a validation error to the form if the input is not less than limit.
// Less only works for int values.
func (val *InputValidation) Less(limit int) *InputValidation {
	return val.Lessf(limit, "%s must be less than %d.", val.InputName, limit)
}

// Lessf is like Less but allows you to specify a custom error message.
// The arguments format and args work exactly like they do in fmt.Sprintf.
func (val *InputValidation) Lessf(limit int, format string, args ...interface{}) *InputValidation {
	return val.validateInt(lessFunc(limit), format, args...)
}

// LessOrEqual adds a validation error to the form if the input is not less than
// or equal to limit. LessOrEqual only works for int values.
func (val *InputValidation) LessOrEqual(limit int) *InputValidation {
	return val.LessOrEqualf(limit, "%s must be less than or equal to %d.", val.InputName, limit)
}

// LessOrEqualf is like LessOrEqual but allows you to specify a custom error
// message. The arguments format and args work exactly like they do in
// fmt.Sprintf.
func (val *InputValidation) LessOrEqualf(limit int, format string, args ...interface{}) *InputValidation {
	return val.validateInt(lessOrEqualFunc(limit), format, args...)
}

// Greater adds a validation error to the form if the input is not greater than
// limit. Greater only works for int values.
func (val *InputValidation) Greater(limit int) *InputValidation {
	return val.Greaterf(limit, "%s must be greater than %d.", val.InputName, limit)
}

// Greaterf is like Greater but allows you to specify a custom error message.
// The arguments format and args work exactly like they do in fmt.Sprintf.
func (val *InputValidation) Greaterf(limit int, format string, args ...interface{}) *InputValidation {
	return val.validateInt(greaterFunc(limit), format, args...)
}

// GreaterOrEqual adds a validation error to the form if the input is not
// greater than or equal to limit. GreaterOrEqual only works for int values.
func (val *InputValidation) GreaterOrEqual(limit int) *InputValidation {
	return val.GreaterOrEqualf(limit, "%s must be greater than or equal to %d.", val.InputName, limit)
}

// GreaterOrEqualf is like GreaterOrEqual but allows you to specify a custom error message.
// The arguments format and args work exactly like they do in fmt.Sprintf.
func (val *InputValidation) GreaterOrEqualf(limit int, format string, args ...interface{}) *InputValidation {
	return val.validateInt(greaterOrEqualFunc(limit), format, args...)
}

// IsInt adds a validation error to the form if the input is not convertible
// to an int.
func (val *InputValidation) IsInt() *InputValidation {
	return val.IsIntf("%s must be an integer.", val.InputName)
}

// IsIntf is like IsInt but allows you to specify a custom error message.
// The arguments format and args work exactly like they do in fmt.Sprintf.
func (val *InputValidation) IsIntf(format string, args ...interface{}) *InputValidation {
	// If the input does not exist or is empty, skip this validation.
	if val.Input == nil || val.Input.RawValue == "" {
		return val
	}
	// Attempt to convert the input value to a int and if the conversion fails,
	// add a validation error.
	if _, err := val.Input.Int(); err != nil {
		val.AddError(format, args...)
	}
	return val
}

func (val *InputValidation) validateInt(validateFunc func(value int) bool, format string, args ...interface{}) *InputValidation {
	// If the input does not exist or is empty, skip this validation.
	if val.Input == nil || val.Input.RawValue == "" {
		return val
	}
	// Attempt to convert the input value to an integer.
	intVal, err := val.Input.Int()
	if err != nil {
		val.AddError("%s must be an integer.", val.InputName)
		return val
	}
	// Call validateFunc and if it returns false, add the appropriate error.
	if !validateFunc(intVal) {
		val.AddError(format, args...)
	}
	return val
}

func lessFloatFunc(limit float64) func(float64) bool {
	return func(value float64) bool {
		return value < limit
	}
}

func lessOrEqualFloatFunc(limit float64) func(float64) bool {
	return func(value float64) bool {
		return value <= limit
	}
}

func greaterFloatFunc(limit float64) func(float64) bool {
	return func(value float64) bool {
		return value > limit
	}
}

func greaterOrEqualFloatFunc(limit float64) func(float64) bool {
	return func(value float64) bool {
		return value >= limit
	}
}

// LessFloat adds a validation error to the form if the input is not less than
// limit. LessFloat only works for float values.
func (val *InputValidation) LessFloat(limit float64) *InputValidation {
	return val.LessFloatf(limit, "%s must be less than %f.", val.InputName, limit)
}

// LessFloatf is like LessFloat but allows you to specify a custom error
// message. The arguments format and args work exactly like they do in
// fmt.Sprintf.
func (val *InputValidation) LessFloatf(limit float64, format string, args ...interface{}) *InputValidation {
	return val.validateFloat(lessFloatFunc(limit), format, args...)
}

// LessOrEqualFloat adds a validation error to the form if the input is not less
// than or equal to limit. LessOrEqualFloat only works for float values.
func (val *InputValidation) LessOrEqualFloat(limit float64) *InputValidation {
	return val.LessOrEqualFloatf(limit, "%s must be less than or equal to %f.", val.InputName, limit)
}

// LessOrEqualFloatf is like LessOrEqualFloat but allows you to specify a custom
// error message. The arguments format and args work exactly like they do in
// fmt.Sprintf.
func (val *InputValidation) LessOrEqualFloatf(limit float64, format string, args ...interface{}) *InputValidation {
	return val.validateFloat(lessOrEqualFloatFunc(limit), format, args...)
}

// GreaterFloat adds a validation error to the form if the input is not greater
// than limit. GreaterFloat only works for float values.
func (val *InputValidation) GreaterFloat(limit float64) *InputValidation {
	return val.GreaterFloatf(limit, "%s must be greater than %f.", val.InputName, limit)
}

// GreaterFloatf is like GreaterFloat but allows you to specify a custom error
// message. The arguments format and args work exactly like they do in
// fmt.Sprintf.
func (val *InputValidation) GreaterFloatf(limit float64, format string, args ...interface{}) *InputValidation {
	return val.validateFloat(greaterFloatFunc(limit), format, args...)
}

// GreaterOrEqualFloat adds a validation error to the form if the input is not
// greater than or equal to limit. GreaterOrEqualFloat only works for float
// values.
func (val *InputValidation) GreaterOrEqualFloat(limit float64) *InputValidation {
	return val.GreaterOrEqualFloatf(limit, "%s must be greater than or equal to %f.", val.InputName, limit)
}

// GreaterOrEqualFloatf is like GreaterOrEqualFloat but allows you to specify a
// custom error message. The arguments format and args work exactly like they do
// in fmt.Sprintf.
func (val *InputValidation) GreaterOrEqualFloatf(limit float64, format string, args ...interface{}) *InputValidation {
	return val.validateFloat(greaterOrEqualFloatFunc(limit), format, args...)
}

// IsFloat adds a validation error to the form if the input is not convertible
// to a float64.
func (val *InputValidation) IsFloat() *InputValidation {
	return val.IsFloatf("%s must be a number.", val.InputName)
}

// IsFloatf is like IsFloat but allows you to specify a custom error message.
// The arguments format and args work exactly like they do in fmt.Sprintf.
func (val *InputValidation) IsFloatf(format string, args ...interface{}) *InputValidation {
	// If the input does not exist or is empty, skip this validation.
	if val.Input == nil || val.Input.RawValue == "" {
		return val
	}
	// Attempt to convert the input value to a float and if the conversion fails,
	// add a validation error.
	if _, err := val.Input.Float(); err != nil {
		val.AddError(format, args...)
	}
	return val
}

func (val *InputValidation) validateFloat(validateFunc func(value float64) bool, format string, args ...interface{}) *InputValidation {
	// If the input does not exist or is empty, skip this validation.
	if val.Input == nil || val.Input.RawValue == "" {
		return val
	}
	// Attempt to convert the input value to a float.
	floatVal, err := val.Input.Float()
	if err != nil {
		val.AddError("%s must be a number.", val.InputName)
		return val
	}
	// Call validateFunc and if it returns false, add the appropriate error.
	if !validateFunc(floatVal) {
		val.AddError(format, args...)
	}
	return val
}

// IsBool adds a validation error to the form if the input is not convertible
// to a bool.
func (val *InputValidation) IsBool() *InputValidation {
	return val.IsBoolf("%s must be either true or false.", val.InputName)
}

// IsBoolf is like IsBool but allows you to specify a custom error message.
// The arguments format and args work exactly like they do in fmt.Sprintf.
func (val *InputValidation) IsBoolf(format string, args ...interface{}) *InputValidation {
	// If the input does not exist or is empty, skip this validation.
	if val.Input == nil || val.Input.RawValue == "" {
		return val
	}
	// Attempt to convert the input to a boolean and if the conversion fails,
	// add a validation error.
	if _, err := val.Input.Bool(); err != nil {
		val.AddError(format, args...)
	}
	return val
}
