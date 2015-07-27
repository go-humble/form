// Copyright 2015 Alex Browne and Soroush Pour.
// Allrights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package main

import "github.com/rusco/qunit"

func main() {
	qunit.Test("Test", func(assert qunit.QUnitAssert) {
		assert.Ok(true, "")
	})
}
