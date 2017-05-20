Humble/Form
===========

[![Version](https://img.shields.io/badge/version-0.0.2-5272B4.svg)](https://github.com/go-humble/form/releases)
[![GoDoc](https://godoc.org/github.com/go-humble/form?status.svg)](https://godoc.org/github.com/go-humble/form)

A go package which provides functions for validating and serializing html forms.
It is intended to be compiled to javascript via
[gopherjs](https://github.com/gopherjs/gopherjs) and run in the browser. Form
works great as a stand-alone package or in combination with other
[Humble](https://github.com/go-humble) packages.

Form is written in pure go. It feels like go, follows go idioms when possible,
and compiles with the go tools. But it is meant to be compiled to javascript and
run in the browser.


Browser Support
---------------

Form works with IE9+ (with a
[polyfill for typed arrays](https://github.com/inexorabletash/polyfill/blob/master/typedarray.js))
and all other modern browsers. Form compiles to javascript via [gopherjs](https://github.com/gopherjs/gopherjs)
and this is a gopherjs limitation.

Form is regularly tested with the latest versions of Firefox, Chrome, and Safari
on Mac OS. Each major or minor release is tested with IE9+ and the latest
versions of Firefox and Chrome on Windows.


Installation
------------

Install form like you would any other go package:

```bash
go get github.com/go-humble/form
```

You will also need to install gopherjs if you don't already have it. The latest
version is recommended. Install gopherjs with:

```
go get -u github.com/gopherjs/gopherjs
```

Finally, you will need to install the
[gopherjs dom bindings](http://godoc.org/honnef.co/go/js/dom):

```
go get -u honnef.co/go/js/dom
```

Quickstart Guide
----------------

### Parsing a Form

The first thing you'll want to do is create a `Form` object. To do this, you can
get an html form element using a query selector and then pass it as an argument
to the [`Parse`](http://godoc.org/github.com/go-humble/form#Parse) function.

```go
// Use a query selector to get the form element.
document := dom.GetWindow().Document()
formEl := document.QuerySelector("#form")
// Parse the form element and get a form.Form object in return.
f, err := form.Parse(formEl)
if err != nil {
	// Handle err.
}
```

### Validations

You can validate the inputs in the form by using the `Validate` method.
`Validate` expects an input name as an argument and returns an `InputValidation`
object, which has chainable methods for validating a single input.

Here's an example:

```go
// Validate the form inputs.
f.Validate("name").Required()
f.Validate("age").Required().IsInt().Greater(0).LessOrEqual(99)
// Check if there were any validation errors.
if f.HasErrors() {
	for _, err := range fmr.Errors {
		// Do something with each error.
	}
}
```

See the
[documentation on the `InputValidation` type](http://godoc.org/github.com/go-humble/form#InputValidation)
for more validation methods.

### Getting Input Values

You can use helper methods to get the value for an input and convert it to
a go type. For example, here's how you can get the value for an input field
named "age" converted to an int:

```go
age, err := f.GetInt("age")
if err != nil {
	// Handle err.
}
```

See the
[documentation on the `Form` type](http://godoc.org/github.com/go-humble/form#Form)
for more helper methods.

### Binding

You can bind a form to any struct by using the
[`Bind`](http://godoc.org/github.com/go-humble/form#Form.Bind) method. `Bind`
compares the names of the fields of the struct to the names of the form inputs.
When it finds a match, it attempts to set the value of the struct field to the
input value.

Suppose you had a form that looked like this:

```html
<form>
	<input name="name" >
	<input name="age" type="number" >
</form>
```

And a `Person` struct with the following definition:

```go
type Person struct {
	Name string
	Age  int
}
```

You could then bind the form to a `Person` using the `Bind` method:

```go
person := &Person{}
if err := f.Bind(person); err != nil {
	// Handle err.
}
```

`Bind` creates a one-way, one-time binding. Changes to the form input values
will not automatically update person, nor will changes to person automatically
change the form input values.

`Bind` supports most primative types and pointers to primative types. If your
struct contains a type that is not supported, you can implement the
[`Binder`](http://godoc.org/github.com/go-humble/form#Binder) or
[`InputBinder`](http://godoc.org/github.com/go-humble/form#InputBinder)
interfaces to define custom behavior.


Testing
-------

Form uses the [karma test runner](http://karma-runner.github.io/0.12/index.html)
to test the code running in actual browsers.

The tests require the following additional dependencies:

- [node.js](http://nodejs.org/) (If you didn't already install it above)
- [karma](http://karma-runner.github.io/0.12/index.html)
- [karma-qunit](https://github.com/karma-runner/karma-qunit)

Don't forget to also install the karma command line tools with `npm install -g karma-cli`.

You will also need to install a launcher for each browser you want to test with,
as well as the browsers themselves. Typically you install a karma launcher with
`npm install -g karma-chrome-launcher`. You can edit the config file at
`karma/test-mac.conf.js` or create a new one (e.g. `karma/test-windows.conf.js`)
if you want to change the browsers that are tested on.

Once you have installed all the dependencies, start karma with
`karma start karma/test-mac.conf.js` (or your customized config file, if
applicable). Once karma is running, you can keep it running in between tests.

Next you need to compile the test.go file to javascript so it can run in the
browsers:

```
gopherjs build karma/go/form_test.go -o karma/js/form_test.js
```

Finally run the tests with `karma run karma/test-mac.conf.js` (changing the name
of the config file if needed).

If you are on a unix-like operating system, you can recompile and run the tests
in one go by running the provided bash script: `./karma/test.sh`.


Contributing
------------

See
[CONTRIBUTING.md](https://github.com/go-humble/form/blob/master/CONTRIBUTING.md)


License
-------

Form is licensed under the MIT License. See the
[LICENSE](https://github.com/go-humble/form/blob/master/LICENSE) file for more
information.
