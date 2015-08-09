// Copyright 2015 Alex Browne and Soroush Pour.
// Allrights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

package form

type InputType string

const (
	InputDefault       InputType = ""
	InputButton        InputType = "button"
	InputCheckbox      InputType = "checkbox"
	InputColor         InputType = "color"
	InputDate          InputType = "date"
	InputDateTime      InputType = "datetime"
	InputDateTimeLocal InputType = "datetime-local"
	InputEmail         InputType = "email"
	InputFile          InputType = "file"
	InputHidden        InputType = "hidden"
	InputImage         InputType = "image"
	InputMonth         InputType = "month"
	InputNumber        InputType = "number"
	InputPassword      InputType = "password"
	InputRadio         InputType = "radio"
	InputRange         InputType = "range"
	InputReset         InputType = "reset"
	InputSearch        InputType = "search"
	InputTel           InputType = "tel"
	InputText          InputType = "text"
	InputTime          InputType = "time"
	InputURL           InputType = "url"
	InputWeek          InputType = "week"
)
