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

type Input interface {
	Name() string
	Value() string
	Type() InputType
}

type Stringer interface {
	String() (string, error)
}

type Inter interface {
	Int() (int, error)
}

type Floater interface {
	Float() (float64, error)
}

type Booler interface {
	Bool() (bool, error)
}

type defaultInput struct {
	name     string
	rawValue string
	typ      InputType
}

func (i defaultInput) Name() string {
	return i.name
}

func (i defaultInput) Value() string {
	return i.rawValue
}

func (i defaultInput) Type() InputType {
	return i.typ
}

func (i defaultInput) String() (string, error) {
	return i.rawValue, nil
}
