// Copyright 2015 Alex Browne and Soroush Pour.
// Allrights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package main

import (
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
			<input name="default" value="foo" />
			<input type="email" name="email" value="foo@example.com" />
			<input type="hidden" name="secret" value="this is a secret" />
			<input type="password" name="password" value="password123" />
			<input type="search" name="search" value="foo bar" />
			<input type="tel" name="phone" value="867-5309" />
			<input type="text" name="text" value="This is some text." />
			<input type="url" name="url" value="http://example.com" />
			</form>`)
		formEl := container.QuerySelector("form").(*dom.HTMLFormElement)
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
			<input name="default" value="23" />
			<input type="email" name="email" value="42" />
			<input type="hidden" name="secret" value="97" />
			<input type="password" name="password" value="123" />
			<input type="search" name="search" value="456" />
			<input type="tel" name="phone" value="8675309" />
			<input type="text" name="text" value="789" />
			<input type="url" name="url" value="123456789" />
			</form>`)
		formEl := container.QuerySelector("form").(*dom.HTMLFormElement)
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		expectedValues := map[string]int{
			"default":  23,
			"email":    42,
			"secret":   97,
			"password": 123,
			"search":   456,
			"phone":    8675309,
			"text":     789,
			"url":      123456789,
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
			<input name="default" value="23.0" />
			<input type="email" name="email" value="42.1" />
			<input type="hidden" name="secret" value="97.2" />
			<input type="password" name="password" value="123.3" />
			<input type="search" name="search" value="456.4" />
			<input type="tel" name="phone" value="8675309.5" />
			<input type="text" name="text" value="789.6" />
			<input type="url" name="url" value="123456789.7" />
			</form>`)
		formEl := container.QuerySelector("form").(*dom.HTMLFormElement)
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		expectedValues := map[string]float64{
			"default":  23.0,
			"email":    42.1,
			"secret":   97.2,
			"password": 123.3,
			"search":   456.4,
			"phone":    8675309.5,
			"text":     789.6,
			"url":      123456789.7,
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
		// should be convertible to floats via GetFloat.
		container.SetInnerHTML(`<form>
			<input name="default" value="true" />
			<input type="email" name="email" value="false" />
			<input type="hidden" name="secret" value="true" />
			<input type="password" name="password" value="false" />
			<input type="search" name="search" value="true" />
			<input type="tel" name="phone" value="false" />
			<input type="text" name="text" value="true" />
			<input type="url" name="url" value="false" />
			</form>`)
		formEl := container.QuerySelector("form").(*dom.HTMLFormElement)
		form, err := form.Parse(formEl)
		assertNoError(assert, err, "")
		expectedValues := map[string]bool{
			"default":  true,
			"email":    false,
			"secret":   true,
			"password": false,
			"search":   true,
			"phone":    false,
			"text":     true,
			"url":      false,
		}
		// Check that the parsed value for each input is correct.
		for name, expectedValue := range expectedValues {
			got, err := form.GetBool(name)
			assertNoError(assert, err, "")
			assert.Equal(got, expectedValue, "Incorrect value for field: "+name)
		}
	})
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
