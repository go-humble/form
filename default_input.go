package form

import "honnef.co/go/js/dom"

func newDefaultInput(el *dom.HTMLInputElement) defaultInput {
	return defaultInput{
		name:     el.Name,
		rawValue: el.Value,
		typ:      InputType(el.Type),
	}
}

type defaultInput struct {
	name     string
	rawValue string
	typ      InputType
}

func (i defaultInput) Name() string {
	return i.name
}

func (i defaultInput) RawValue() string {
	return i.rawValue
}

func (i defaultInput) Type() InputType {
	return i.typ
}
