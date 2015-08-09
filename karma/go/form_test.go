// Copyright 2015 Alex Browne and Soroush Pour.
// Allrights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package main

import (
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
			<input type="text" name="text" value="789" >
			<input type="number" name="number" value="123456789" >
			</form>`)
		formEl := container.QuerySelector("form")
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		expectedValues := map[string]int{
			"default": 23,
			"tel":     8675309,
			"text":    789,
			"number":  123456789,
		}
		// Check that the parsed value for each input is correct.
		for name, expectedValue := range expectedValues {
			got, err := form.GetInt(name)
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
			"Expected form to have 2 errors because 'non-existing' and 'empty' are required.")
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
		// Check that an invalid field does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it was not less than 10."
		form.Validate("invalid").Lessf(10, customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not less than 10.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with Requiredf")
		// Check that a input which is not an integer adds the correct error
		// message.
		form.Validate("non-integer").Less(10)
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 error because 'non-integer' is not an integer.")
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
		// Check that an invalid field does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it was not less than or equal to 10."
		form.Validate("invalid").LessOrEqualf(10, customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not less than or equal to 10.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with Requiredf")
		// Check that a input which is not an integer adds the correct error
		// message.
		form.Validate("non-integer").LessOrEqual(10)
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 error because 'non-integer' is not an integer.")
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
		// Check that an invalid field does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it was not greater than 10."
		form.Validate("invalid").Greaterf(10, customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not greater than 10.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with Requiredf")
		// Check that a input which is not an integer adds the correct error
		// message.
		form.Validate("non-integer").Greater(10)
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 error because 'non-integer' is not an integer.")
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
		// Check that an invalid field does add a validation error with a custom
		// message.
		customMessage := "invalid input was invalid becase it was not greater than or equal to 10."
		form.Validate("invalid").GreaterOrEqualf(10, customMessage)
		assert.Equal(len(form.Errors), 1,
			"Expected form to have 1 error because 'invalid' is not greater than or equal to 10.")
		assert.Equal(form.Errors[0].Error(), customMessage,
			"Custom message was not set with Requiredf")
		// Check that a input which is not an integer adds the correct error
		// message.
		form.Validate("non-integer").GreaterOrEqual(10)
		assert.Equal(len(form.Errors), 2,
			"Expected form to have 2 error because 'non-integer' is not an integer.")
		assert.Equal(form.Errors[1].Error(), "non-integer must be an integer.",
			"Error was not added when input was a non-integer.")
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
	}
}
