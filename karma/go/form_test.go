// Copyright 2015 Alex Browne and Soroush Pour.
// Allrights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package main

import (
	"strings"
	"time"

	"github.com/go-humble/form"
	"github.com/rusco/qunit"
	"honnef.co/go/js/dom"
)

var (
	document  = dom.GetWindow().Document()
	body      = document.QuerySelector("body")
	container dom.Element
)

func init() {
	container = document.CreateElement("div")
	container.SetID("container")
	body.AppendChild(container)
}

func reset() {
	container.SetInnerHTML("")
}

// CustomBinder implements form.Binder.
type CustomBinder struct {
	Int    int
	String string
}

// BindForm implements the BindForm method of form.Binder. Instead of binding
// the fields directly, We'll add 1 to the Int field and prepend a "_" to the
// string field.
func (b *CustomBinder) BindForm(form *form.Form) error {
	if _, found := form.Inputs["int"]; found {
		valInt, err := form.GetInt("int")
		if err != nil {
			return err
		}
		b.Int = valInt + 1
	}
	if _, found := form.Inputs["string"]; found {
		valString, err := form.GetString("string")
		if err != nil {
			return err
		}
		b.String = "_" + valString
	}
	return nil
}

// Name is a custom type which implements form.InputBinder.
type Name struct {
	First string
	Last  string
}

// BindInput implements the BindInput method of form.InputBinder. We split the
// input into a first name and last name and assign each field of Name manually.
func (name *Name) BindInput(input *form.Input) error {
	names := strings.Split(input.RawValue, " ")
	name.First = names[0]
	name.Last = names[1]
	return nil
}

// CustomInputBinder is a struct which has a single field of type Name. When we
// call form.Bind on it, we expect Name.BindInput to be invoked.
type CustomInputBinder struct {
	Name Name
}

func main() {
	qunit.Test("GetString", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values. All the input types here
		// should be convertible to strings via GetString.
		container.SetInnerHTML(`<form>
			<input name="default" value="foo" >
			<input type="email" name="email" value="foo@example.com" >
			<input type="hidden" name="secret" value="this is a secret" >
			<input type="password" name="password" value="password123" >
			<input type="search" name="search" value="foo bar" >
			<input type="tel" name="phone" value="867-5309" >
			<input type="text" name="text" value="This is some text." >
			<input type="url" name="url" value="http://example.com" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		expectedValues := map[string]string{
			"default":  "foo",
			"email":    "foo@example.com",
			"secret":   "this is a secret",
			"password": "password123",
			"search":   "foo bar",
			"phone":    "867-5309",
			"text":     "This is some text.",
			"url":      "http://example.com",
		}
		// Check that the parsed value for each input is correct.
		for name, expectedValue := range expectedValues {
			got, err := form.GetString(name)
			assertNoError(assert, err, "")
			assert.Equal(got, expectedValue, "Incorrect value for field: "+name)
		}
	})

	qunit.Test("GetInt", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values. All the input types here
		// should be convertible to ints via GetInt.
		container.SetInnerHTML(`<form>
			<input name="default" value="23" >
			<input type="tel" name="tel" value="8675309" >
			<input type="text" name="text" value="-789" >
			<input type="number" name="number" value="123456789" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		expectedValues := map[string]int{
			"default": 23,
			"tel":     8675309,
			"text":    -789,
			"number":  123456789,
		}
		// Check that the parsed value for each input is correct.
		for name, expectedValue := range expectedValues {
			got, err := form.GetInt(name)
			assertNoError(assert, err, "")
			assert.Equal(got, expectedValue, "Incorrect value for field: "+name)
		}
	})

	qunit.Test("GetUint", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values. All the input types here
		// should be convertible to uints via GetUint.
		container.SetInnerHTML(`<form>
			<input name="default" value="23" >
			<input type="tel" name="tel" value="8675309" >
			<input type="text" name="text" value="789" >
			<input type="number" name="number" value="123456789" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		expectedValues := map[string]uint{
			"default": 23,
			"tel":     8675309,
			"text":    789,
			"number":  123456789,
		}
		// Check that the parsed value for each input is correct.
		for name, expectedValue := range expectedValues {
			got, err := form.GetUint(name)
			assertNoError(assert, err, "")
			assert.Equal(got, expectedValue, "Incorrect value for field: "+name)
		}
	})

	qunit.Test("GetFloat", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values. All the input types here
		// should be convertible to floats via GetFloat.
		container.SetInnerHTML(`<form>
			<input name="default" value="23.0" >
			<input type="text" name="text" value="789.6" >
			<input type="number" name="number" value="123456789.7" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		expectedValues := map[string]float64{
			"default": 23.0,
			"text":    789.6,
			"number":  123456789.7,
		}
		// Check that the parsed value for each input is correct.
		for name, expectedValue := range expectedValues {
			got, err := form.GetFloat(name)
			assertNoError(assert, err, "")
			assert.Equal(got, expectedValue, "Incorrect value for field: "+name)
		}
	})

	qunit.Test("GetBool", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values. All the input types here
		// should be convertible to booleans via GetBool.
		container.SetInnerHTML(`<form>
			<input name="default" value="true" >
			<input type="text" name="text" value="true" >
			<input type="checkbox" name="checkbox" checked >
			<input type="radio" name="radio" checked >
			<input type="checkbox" name="checkbox-false" >
			<input type="radio" name="radio-false" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		expectedValues := map[string]bool{
			"default":        true,
			"text":           true,
			"checkbox":       true,
			"radio":          true,
			"checkbox-false": false,
			"radio-false":    false,
		}
		// Check that the parsed value for each input is correct.
		for name, expectedValue := range expectedValues {
			got, err := form.GetBool(name)
			assertNoError(assert, err, "")
			assert.Equal(got, expectedValue, "Incorrect value for field: "+name)
		}
	})

	qunit.Test("GetTime", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values. All the input types here
		// should be convertible to time.Time via GetTime.
		container.SetInnerHTML(`<form>
			<input name="date" type="date" value="1992-09-29" >
			<input name="datetime" type="datetime" value="1985-12-03T23:59:34-08:00" >
			<input name="datetime-local" type="datetime-local" value="1985-04-12T23:20:50.52" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		rfc3339Date := "2006-01-02"
		rfc3339DatetimeLocal := "2006-01-02T15:04:05.999999999"
		expectedValues := map[string]time.Time{
			"date":           mustParseTime(rfc3339Date, "1992-09-29"),
			"datetime":       mustParseTime(time.RFC3339, "1985-12-03T23:59:34-08:00"),
			"datetime-local": mustParseTime(rfc3339DatetimeLocal, "1985-04-12T23:20:50.52"),
		}
		// Check that the parsed value for each input is correct.
		for name, expectedValue := range expectedValues {
			got, err := form.GetTime(name)
			assertNoError(assert, err, "Error for field: "+name)
			assert.DeepEqual(got, expectedValue, "Incorrect value for field: "+name)
		}
	})

	qunit.Test("ValidateRequired", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="non-empty" value="foo" >
			<input name="empty" value="" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Check that a non-empty input does not add a validation error.
		form.Validate("non-empty").Required()
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that an empty input does add a validation error with a custom
		// message.
		customMessage := "empty cannot be blank."
		form.Validate("empty").Requiredf(customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'empty' is required.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with Requiredf")
		// Check that a non-existing input does add a validation error.
		form.Validate("non-existing").Required()
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 errorss because 'non-existing' and 'empty' are required.")
	})

	qunit.Test("ValidateLess", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="valid" value="5" >
			<input name="invalid" value="10" >
			<input name="non-integer" value="foo" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Check that a valid input does not add any validation errors.
		form.Validate("valid").Less(10)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that a non-existing input does not add any validation errors.
		form.Validate("non-existing").Less(10)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that an invalid input does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it was not less than 10."
		form.Validate("invalid").Lessf(10, customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not less than 10.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with Lessf")
		// Check that a input which is not an integer adds the correct error
		// message.
		form.Validate("non-integer").Less(10)
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 errors because 'non-integer' is not an integer.")
		assert.Equal(form.Errors[1].Error(), "non-integer must be an integer.",
			"Error was not added when input was a non-integer.")
	})

	qunit.Test("ValidateLessOrEqual", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="valid" value="5" >
			<input name="invalid" value="11" >
			<input name="non-integer" value="foo" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Check that a valid input does not add any validation errors.
		form.Validate("valid").LessOrEqual(10)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that a non-existing input does not add any validation errors.
		form.Validate("non-existing").LessOrEqual(10)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that an invalid input does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it was not less than or equal to 10."
		form.Validate("invalid").LessOrEqualf(10, customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not less than or equal to 10.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with LessOrEqualf")
		// Check that a input which is not an integer adds the correct error
		// message.
		form.Validate("non-integer").LessOrEqual(10)
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 errors because 'non-integer' is not an integer.")
		assert.Equal(form.Errors[1].Error(), "non-integer must be an integer.",
			"Error was not added when input was a non-integer.")
	})

	qunit.Test("ValidateGreater", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="valid" value="15" >
			<input name="invalid" value="10" >
			<input name="non-integer" value="foo" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Check that a valid input does not add any validation errors.
		form.Validate("valid").Greater(10)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that a non-existing input does not add any validation errors.
		form.Validate("non-existing").Greater(10)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that an invalid input does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it was not greater than 10."
		form.Validate("invalid").Greaterf(10, customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not greater than 10.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with Greaterf")
		// Check that a input which is not an integer adds the correct error
		// message.
		form.Validate("non-integer").Greater(10)
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 errors because 'non-integer' is not an integer.")
		assert.Equal(form.Errors[1].Error(), "non-integer must be an integer.",
			"Error was not added when input was a non-integer.")
	})

	qunit.Test("ValidateGreaterOrEqual", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="valid" value="15" >
			<input name="invalid" value="9" >
			<input name="non-integer" value="foo" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Check that a valid input does not add any validation errors.
		form.Validate("valid").GreaterOrEqual(10)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that a non-existing input does not add any validation errors.
		form.Validate("non-existing").GreaterOrEqual(10)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that an invalid input does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it was not greater than or equal to 10."
		form.Validate("invalid").GreaterOrEqualf(10, customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not greater than or equal to 10.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with GreaterOrEqualf")
		// Check that a input which is not an integer adds the correct error
		// message.
		form.Validate("non-integer").GreaterOrEqual(10)
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 errors because 'non-integer' is not an integer.")
		assert.Equal(form.Errors[1].Error(), "non-integer must be an integer.",
			"Error was not added when input was a non-integer.")
	})

	qunit.Test("ValidateIsInt", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="valid" value="5" >
			<input name="invalid" value="foo" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Check that a valid input does not add any validation errors.
		form.Validate("valid").IsInt()
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that a non-existing input does not add any validation errors.
		form.Validate("non-existing").IsInt()
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that an invalid input does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it not an integer."
		form.Validate("invalid").IsIntf(customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not an integer.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with IsIntf")
	})

	qunit.Test("ValidateLessFloat", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="valid" value="5.4" >
			<input name="invalid" value="10.0" >
			<input name="non-float" value="foo" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Check that a valid input does not add any validation errors.
		form.Validate("valid").LessFloat(10.0)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that a non-existing input does not add any validation errors.
		form.Validate("non-existing").LessFloat(10.0)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that an invalid input does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it was not less than 10.0."
		form.Validate("invalid").LessFloatf(10.0, customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not less than 10.0.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with LessFloatf")
		// Check that a input which is not a float adds the correct error
		// message.
		form.Validate("non-float").LessFloat(10.0)
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 errors because 'non-float' is not a number.")
		assert.Equal(form.Errors[1].Error(), "non-float must be a number.",
			"Error was not added when input was a non-float.")
	})

	qunit.Test("ValidateLessOrEqualFloat", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="valid" value="5.4" >
			<input name="invalid" value="10.1" >
			<input name="non-float" value="foo" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Check that a valid input does not add any validation errors.
		form.Validate("valid").LessOrEqualFloat(10.0)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that a non-existing input does not add any validation errors.
		form.Validate("non-existing").LessOrEqualFloat(10.0)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that an invalid input does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it was not less than or equal to 10.0."
		form.Validate("invalid").LessOrEqualFloatf(10.0, customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not less than or equal to 10.0.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with LessOrEqualFloatf")
		// Check that a input which is not a float adds the correct error
		// message.
		form.Validate("non-float").LessOrEqualFloat(10.0)
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 errors because 'non-float' is not a number.")
		assert.Equal(form.Errors[1].Error(), "non-float must be a number.",
			"Error was not added when input was a non-float.")
	})

	qunit.Test("ValidateGreaterFloat", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="valid" value="15.7" >
			<input name="invalid" value="10.0" >
			<input name="non-float" value="foo" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Check that a valid input does not add any validation errors.
		form.Validate("valid").GreaterFloat(10.0)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that a non-existing input does not add any validation errors.
		form.Validate("non-existing").GreaterFloat(10.0)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that an invalid input does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it was not greater than 10.0."
		form.Validate("invalid").GreaterFloatf(10.0, customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not greater than 10.0.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with GreaterFloatf")
		// Check that a input which is not a float adds the correct error
		// message.
		form.Validate("non-float").GreaterFloat(10.0)
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 errors because 'non-float' is not a number.")
		assert.Equal(form.Errors[1].Error(), "non-float must be a number.",
			"Error was not added when input was a non-float.")
	})

	qunit.Test("ValidateGreaterOrEqualFloat", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="valid" value="15.7" >
			<input name="invalid" value="9.9" >
			<input name="non-float" value="foo" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Check that a valid input does not add any validation errors.
		form.Validate("valid").GreaterOrEqualFloat(10.0)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that a non-existing input does not add any validation errors.
		form.Validate("non-existing").GreaterOrEqualFloat(10.0)
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that an invalid input does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it was not greater than or equal to 10.0."
		form.Validate("invalid").GreaterOrEqualFloatf(10.0, customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not greater than or equal to 10.0.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with GreaterOrEqualFloatf")
		// Check that a input which is not a float adds the correct error
		// message.
		form.Validate("non-float").GreaterOrEqualFloat(10.0)
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 errors because 'non-float' is not a number.")
		assert.Equal(form.Errors[1].Error(), "non-float must be a number.",
			"Error was not added when input was a non-float.")
	})

	qunit.Test("ValidateIsFloat", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="valid" value="5.3" >
			<input name="invalid" value="foo" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Check that a valid input does not add any validation errors.
		form.Validate("valid").IsFloat()
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that a non-existing input does not add any validation errors.
		form.Validate("non-existing").IsFloat()
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that an invalid input does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it not a float."
		form.Validate("invalid").IsFloatf(customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not a float.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with IsFloatf")
	})

	qunit.Test("ValidateIsBool", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="valid" value="true" >
			<input type="checkbox" name="checkbox" checked >
			<input name="invalid" value="foo" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Check that a valid input does not add any validation errors.
		form.Validate("valid").IsBool()
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that the checkbox input does not add any validation errors.
		form.Validate("checkbox").IsBool()
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that a non-existing input does not add any validation errors.
		form.Validate("non-existing").IsBool()
		assert.Equal(form.HasErrors(), false, "Expected form to have no errors")
		// Check that an invalid input does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it not a boolean."
		form.Validate("invalid").IsBoolf(customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not a boolean.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with IsBoolf")
	})

	qunit.Test("Bind", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="string" value="foo" >
			<input name="bytes" value="bar" >
			<input type="number" name="int" value="4" >
			<input type="number" name="int8" value="8" >
			<input type="number" name="int16" value="15" >
			<input type="number" name="int32" value="16" >
			<input type="number" name="int64" value="23" >
			<input type="number" name="uint" value="42" >
			<input type="number" name="uint8" value="1" >
			<input type="number" name="uint16" value="2" >
			<input type="number" name="uint32" value="3" >
			<input type="number" name="uint64" value="4" >
			<input type="number" name="float32" value="39.7" >
			<input type="number" name="float64" value="12.6" >
			<input type="checkbox" name="bool" checked >
			<input type="datetime" name="time" value="1985-12-03T23:59:34-08:00" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Bind the form to some target and check the results.
		target := struct {
			String  string
			Bytes   []byte
			Int     int
			Int8    int8
			Int16   int16
			Int32   int32
			Int64   int64
			Uint    uint
			Uint8   uint8
			Uint16  uint16
			Uint32  uint32
			Uint64  uint64
			Float32 float32
			Float64 float64
			Bool    bool
			Time    time.Time
		}{}
		err = form.Bind(&target)
		assertNoError(assert, err, "")
		assert.Equal(target.String, "foo", "target.String was not correct.")
		assert.DeepEqual(target.Bytes, []byte("bar"), "target.Bytes was not correct.")
		assert.Equal(target.Int, 4, "target.Int was not correct.")
		assert.Equal(target.Int8, 8, "target.Int8 was not correct.")
		assert.Equal(target.Int16, 15, "target.Int16 was not correct.")
		assert.Equal(target.Int32, 16, "target.Int32 was not correct.")
		assert.Equal(target.Int64, 23, "target.Int64 was not correct.")
		assert.Equal(target.Uint, 42, "target.Uint was not correct.")
		assert.Equal(target.Uint8, 1, "target.Uint8 was not correct.")
		assert.Equal(target.Uint16, 2, "target.Uint16 was not correct.")
		assert.Equal(target.Uint32, 3, "target.Uint32 was not correct.")
		assert.Equal(target.Uint64, 4, "target.Uint64 was not correct.")
		assert.Equal(target.Bool, true, "target.Bool was not correct.")
		assert.DeepEqual(target.Time,
			mustParseTime(time.RFC3339, "1985-12-03T23:59:34-08:00"),
			"target.Time was not correct.")
	})

	qunit.Test("BindWithPointers", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="string" value="foo" >
			<input name="bytes" value="bar" >
			<input type="number" name="int" value="4" >
			<input type="number" name="int8" value="8" >
			<input type="number" name="int16" value="15" >
			<input type="number" name="int32" value="16" >
			<input type="number" name="int64" value="23" >
			<input type="number" name="uint" value="42" >
			<input type="number" name="uint8" value="1" >
			<input type="number" name="uint16" value="2" >
			<input type="number" name="uint32" value="3" >
			<input type="number" name="uint64" value="4" >
			<input type="number" name="float32" value="39.7" >
			<input type="number" name="float64" value="12.6" >
			<input type="checkbox" name="bool" checked >
			<input type="datetime" name="time" value="1985-12-03T23:59:34-08:00" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Bind the form to some target and check the results.
		target := struct {
			String  *string
			Bytes   *[]byte
			Int     *int
			Int8    *int8
			Int16   *int16
			Int32   *int32
			Int64   *int64
			Uint    *uint
			Uint8   *uint8
			Uint16  *uint16
			Uint32  *uint32
			Uint64  *uint64
			Float32 *float32
			Float64 *float64
			Bool    *bool
			Time    *time.Time
		}{}
		err = form.Bind(&target)
		assertNoError(assert, err, "")
		assert.Equal(*target.String, "foo", "target.String was not correct.")
		assert.DeepEqual(*target.Bytes, []byte("bar"), "target.Bytes was not correct.")
		assert.Equal(*target.Int, 4, "target.Int was not correct.")
		assert.Equal(*target.Int8, 8, "target.Int8 was not correct.")
		assert.Equal(*target.Int16, 15, "target.Int16 was not correct.")
		assert.Equal(*target.Int32, 16, "target.Int32 was not correct.")
		assert.Equal(*target.Int64, 23, "target.Int64 was not correct.")
		assert.Equal(*target.Uint, 42, "target.Uint was not correct.")
		assert.Equal(*target.Uint8, 1, "target.Uint8 was not correct.")
		assert.Equal(*target.Uint16, 2, "target.Uint16 was not correct.")
		assert.Equal(*target.Uint32, 3, "target.Uint32 was not correct.")
		assert.Equal(*target.Uint64, 4, "target.Uint64 was not correct.")
		assert.Equal(*target.Bool, true, "target.Bool was not correct.")
		assert.DeepEqual(*target.Time,
			mustParseTime(time.RFC3339, "1985-12-03T23:59:34-08:00"),
			"target.Time was not correct.")
	})

	qunit.Test("Binder", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="string" value="foo" >
			<input name="int" type="number" value="42" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Bind the form to our CustomBinder type and check the results.
		binder := &CustomBinder{}
		err = form.Bind(binder)
		assertNoError(assert, err, "")
		assert.Equal(binder.String, "_foo", "binder.String was not correct.")
		assert.Equal(binder.Int, 43, "binder.Int was not correct.")
	})

	qunit.Test("InputBinder", func(assert qunit.QUnitAssert) {
		defer reset()
		// Create a form with some inputs and values.
		container.SetInnerHTML(`<form>
			<input name="name" value="Foo Bar" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		// Bind the form to our CustomInputBinder type and check the results.
		binder := &CustomInputBinder{}
		err = form.Bind(binder)
		assertNoError(assert, err, "")
		assert.Equal(binder.Name.First, "Foo", "binder.Name.First was not correct.")
		assert.Equal(binder.Name.Last, "Bar", "binder.Name.Last was not correct.")
	})
}

func mustParseTime(layout string, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func assertNoError(assert qunit.QUnitAssert, err error, msg string) {
	if err != nil {
		if msg == "" {
			assert.Equal(err, nil, err.Error())
		} else {
			assert.Equal(err, nil, msg+"\n"+err.Error())
		}
	} else {
		assert.Equal(err, nil, "")
	}
}
