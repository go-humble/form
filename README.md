Humble/Form
===========

[![GoDoc](https://godoc.org/github.com/go-humble/form?status.svg)](https://godoc.org/github.com/go-humble/form)

Version X.X.X (develop)

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

You will also need to install gopherjs if you don't already have it. The latest version is
recommended. Install gopherjs with:

```
go get -u github.com/gopherjs/gopherjs
```


Development Status
------------------

Form is still under development and is not yet ready for use. The first usable
release will be version 0.1.0. If you are curious, feel free to poke around,
but don't expect anything to work yet.


Quickstart Guide
----------------

Coming soon!


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
[CONTRIBUTING.md](https://github.com/go-humble/view/blob/master/CONTRIBUTING.md)


License
-------

View is licensed under the MIT License. See the
[LICENSE](https://github.com/go-humble/view/blob/master/LICENSE) file for more
information.
