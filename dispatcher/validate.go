package dispatcher

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	ErrInvalidArgExpectedSlice = errors.New("invalid arguments, expected slice")
)

// Validate verify and convert any param into types expected by the method
//   If no args         -> return empty
//   If 1 arg           -> directly parse and convert the param and return it
//   IF 2 or more arg   -> verify that param is an array and loop through it to
//  convert to an array of interface with correct type
func (d *Dispatcher) Validate(service, method string, param interface{}) ([]interface{}, error) {
	funcMethod, err := d.GetMethod(service, method)
	if err != nil {
		return nil, err
	}

	args := funcMethod.argsTypes[1:]

	switch len(args) {
	case 0:
		return []interface{}{}, nil
	case 1:
		p, err := verifyParam(args[0], param)
		if err != nil {
			return nil, err
		}

		return []interface{}{p}, err
	default:
		if reflect.TypeOf(param).Kind() != reflect.Slice {
			return nil, ErrInvalidArgExpectedSlice
		}

		params, err := convertInterfaceToArray(param)
		if err != nil {
			return nil, err
		}

		res := make([]interface{}, len(args))
		for i, e := range params {
			p, err := verifyParam(args[i], e)
			if err != nil {
				return nil, err
			}
			res[i] = p
		}

		return res, nil
	}
}

// verifyParam convert the param into the type of the arg
// Since a simple reflect is not enough to verify if the param is type of arg
// this function use json.Unmarshal to correctly convert the param
func verifyParam(arg reflect.Type, param interface{}) (interface{}, error) {
	expectedType := reflect.StructOf([]reflect.StructField{{
		Name: "Placeholder",
		Type: arg,
	}})
	expected := reflect.New(expectedType).Interface()

	placeholder := struct {
		Placeholder interface{}
	}{
		Placeholder: param,
	}

	data, err := json.Marshal(placeholder)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &expected)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(expected).Elem().FieldByName("Placeholder")
	return value.Interface(), nil
}

// convertInterfaceToArray is a utility function used to transform
// an interface into an array of interface
// The result can then be used to populate arguments to the dispatcher
func convertInterfaceToArray(value interface{}) ([]interface{}, error) {
	var out []interface{}

	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() != reflect.Slice {
		return nil, ErrInvalidArgExpectedSlice
	}

	for i := 0; i < reflectValue.Len(); i++ {
		out = append(out, reflectValue.Index(i).Interface())
	}

	return out, nil
}
