package form

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Binder has a single method, BindForm, which binds the values in the form to
// the method receiver. It can be used to override the behavior of Bind. Note
// that because BindForm is expected to set the value of the receiver, the
// receiver should always be a pointer.
type Binder interface {
	BindForm(*Form) error
}

// InputBinder has a single method, BindInput, which binds an input value to
// the receiver. It can be used to override the behavior of Bind. Typically,
// InputBinder should be implemented by some type which is a field of the struct
// passed to Bind. Note that because BindInput is expected to set the value of
// the receiver, the receiver should always be a pointer.
type InputBinder interface {
	BindInput(*Input) error
}

// inputBinderType is the reflect.Type of the InputBinder interface.
var inputBinderType = reflect.TypeOf([]InputBinder{}).Elem()

// Bind attempts to bind the form input values to v, which must be a pointer to
// a struct. Bind performs a one-way, one-time binding. Changes to the form
// input values will not automatically update v, nor will changes to v
// automatically change the form input values.
//
// The following struct field types are supported: string, []byte, int, int8,
// int16, int32, int64, float32, float64, bool, time.Time, and pointers to any
// of the preceding types. Bind attempts to match struct field names with input
// names in a case-insensitive manner. So an input with a name attribute "foo"
// will be assigned to the field v.Foo. Bind will simply ignore fields of v
// which do not match any input names and input names which do not match any
// fields of v. Because Bind uses reflection, only exported fields of v (those
// which start with a capital letter) will be affected.
//
// If v implements Binder, form.Bind will just call v.BindForm. Similarly if any
// of the fields of v implement InputBinder, Bind will call BindInput on the
// specific field. For all other typees, Bind will attempt to do the binding
// using a series of conversion rules. If a field type is not supported and that
// field does not implement InputBinder, Bind will return an error.
//
// Bind will return an error if the type of v is not a pointer to a struct or if
// there is an error arising from binding any of the individual inputs to
// fields.
func (form *Form) Bind(v interface{}) error {
	// If v implements Binder, call v.BindForm.
	if binder, ok := v.(Binder); ok {
		return binder.BindForm(form)
	}
	// Make sure the type of v is a pointer to a struct.
	ptrType := reflect.TypeOf(v)
	if ptrType.Kind() != reflect.Ptr {
		return fmt.Errorf("form: Bind expects a pointer to a struct, but got: %T", v)
	}
	typ := ptrType.Elem()
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("form: Bind expects a pointer to a struct, but got: %T", v)
	}
	ptrVal := reflect.ValueOf(v)
	if ptrVal.IsNil() {
		return fmt.Errorf("form: Argument to Bind was nil")
	}
	val := ptrVal.Elem()
	// Iterate through the fields of v and the inputs of form and find
	// case-insenstive name matches.
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		for _, input := range form.Inputs {
			if strings.EqualFold(field.Name, input.Name) {
				// If the names match, attempt to bind the input to the field.
				if err := bindInput(field.Type, val.FieldByName(field.Name), input); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

// getUnderlyingFieldType returns the underlying type of the given fieldType.
// That is, if fieldType is a poiner, getUnderlyingFieldType will dereference
// it. If the dereferenced type is also a pointer, it will dereference that,
// and continue dereferencing until it reaches a type which is not a pointer.
// E.g., if fieldType is *int, getUnderlyingFieldType will return int.
func getUnderlyingFieldType(fieldType reflect.Type) reflect.Type {
	for fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	return fieldType
}

// bindInput attempts to bind the given input to the field represented by
// fieldType and fieldVal.
func bindInput(fieldType reflect.Type, fieldVal reflect.Value, input *Input) error {
	// Check if the field type implements InputBinder. If it does, this is an
	// indication that the caller wants to use custom behavior to bind the input
	// to the field value, so we'll skip the normal behavior.
	switch {
	case fieldType.Implements(inputBinderType):
		// In this case the fieldType implements InputBinder directly.
		// If the fieldVal is a nil pointer, initialize it.
		if fieldVal.Kind() == reflect.Ptr && fieldVal.IsNil() {
			fieldVal.Set(reflect.New(fieldType.Elem()))
		}
		// We know the first and only return value has type error, so check if it
		// is non-nil and return it if so.
		returnVals := fieldVal.MethodByName("BindInput").Call([]reflect.Value{reflect.ValueOf(input)})
		if !returnVals[0].IsNil() {
			return returnVals[0].Interface().(error)
		}
		return nil
	case reflect.PtrTo(fieldType).Implements(inputBinderType):
		// In this case, a pointer to the fieldType implements InputBinder. We take
		// the address of fieldVal and call BindInput on it.
		returnVals := fieldVal.Addr().MethodByName("BindInput").Call([]reflect.Value{reflect.ValueOf(input)})
		// We know the first and only return value has type error, so check if it
		// is non-nil and return it if so.
		if !returnVals[0].IsNil() {
			return returnVals[0].Interface().(error)
		}
		return nil
	}
	// Check the underlying type of the field. If it is one of the supported
	// types, attempt to bind the input to it. Otherwise we will return an error.
	underlyingType := getUnderlyingFieldType(fieldType)
	switch underlyingType.Kind() {
	case reflect.String:
		setUnderlyingFieldValue(fieldVal, reflect.ValueOf(input.RawValue))
		return nil
	case reflect.Slice:
		if underlyingType.Elem().Kind() == reflect.Uint8 {
			// The underlying type of the struct field is a slice of bytes, which
			// we can easily convert to.
			setUnderlyingFieldValue(fieldVal, reflect.ValueOf([]byte(input.RawValue)))
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if err := bindInt(fieldVal, underlyingType, input); err != nil {
			return err
		}
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if err := bindUint(fieldVal, underlyingType, input); err != nil {
			return err
		}
		return nil
	case reflect.Float32, reflect.Float64:
		if err := bindFloat(fieldVal, underlyingType, input); err != nil {
			return err
		}
		return nil
	case reflect.Bool:
		valBool, err := input.Bool()
		if err != nil {
			return err
		}
		setUnderlyingFieldValue(fieldVal, reflect.ValueOf(valBool))
		return nil
	}
	if underlyingType == reflect.TypeOf(time.Time{}) {
		valTime, err := input.Time()
		if err != nil {
			return err
		}
		setUnderlyingFieldValue(fieldVal, reflect.ValueOf(valTime))
		return nil
	}
	return fmt.Errorf(
		"form: Don't know how to bind input of type %s and value %v to struct field of type %s",
		string(input.Type),
		input.RawValue,
		fieldType.String())
}

// bindInt assumes that the field is an int type (int, int8, int16, int32,
// or int64) and attempts to bind the input to the field.
func bindInt(fieldVal reflect.Value, underlyingType reflect.Type, input *Input) error {
	valInt, err := input.Int()
	if err != nil {
		return err
	}
	sizedInt := reflect.ValueOf(valInt).Convert(underlyingType)
	setUnderlyingFieldValue(fieldVal, sizedInt)
	return nil
}

// bindUint assumes that the field is an uint type (uint, uint8, uint16, uint32,
// or uint64) and attempts to bind the input to the field.
func bindUint(fieldVal reflect.Value, underlyingType reflect.Type, input *Input) error {
	valUint, err := input.Uint()
	if err != nil {
		return err
	}
	sizedUint := reflect.ValueOf(valUint).Convert(underlyingType)
	setUnderlyingFieldValue(fieldVal, sizedUint)
	return nil
}

// bindFloat assumes that the field is an float type (float32, float64) and
// attempts to bind the input to the field.
func bindFloat(fieldVal reflect.Value, underlyingType reflect.Type, input *Input) error {
	valFloat, err := input.Float()
	if err != nil {
		return err
	}
	sizedFloat := reflect.ValueOf(valFloat).Convert(underlyingType)
	setUnderlyingFieldValue(fieldVal, sizedFloat)
	return nil
}

// setUnderlyingFieldValue sets the underlying value of fieldVal to the given
// val. E.g. if fieldVal has type *int it will set the underlying int value.
func setUnderlyingFieldValue(fieldVal reflect.Value, val reflect.Value) {
	for fieldVal.Kind() == reflect.Ptr {
		if fieldVal.IsNil() {
			fieldVal.Set(reflect.New(fieldVal.Type()))
		}
		fieldVal = fieldVal.Elem()
	}
	fieldVal.Set(val)
}
