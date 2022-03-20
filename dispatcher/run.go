package dispatcher

import (
	"errors"
	"reflect"
)

var (
	ErrInvalidArgumentsCount = errors.New("invalid number of arguments provided")
	ErrInvalidArgumentType   = errors.New("invalid argument type")
)

// Run will call the service's method previously registered using the given arguments.
func (d *Dispatcher) Run(service, method string, args ...interface{}) ([]reflect.Value, error) {
	// Checks that the service has been registered and method exists.
	if _, err := d.GetMethod(service, method); err != nil {
		return nil, err
	}

	if !d.verifyArgumentCount(service, method, args...) {
		return nil, ErrInvalidArgumentsCount
	}

	if !d.verifyArgumentTypes(service, method, args...) {
		return nil, ErrInvalidArgumentType
	}

	// Prepare the arguments.
	inArgs := make([]reflect.Value, len(args)+1)
	inArgs[0] = d.services[service].service

	for i, arg := range args {
		inArgs[i+1] = reflect.ValueOf(arg)
	}

	// Run the method.
	output := d.services[service].methods[method].function.Call(inArgs)

	return output, nil
}

func (d *Dispatcher) verifyArgumentCount(service, method string, args ...interface{}) bool {
	if d.services[service].methods[method].isVariadic {
		// If no arguments are passed in the variadic slice, there's nothing to verify
		if len(args) == 0 {
			return true
		}

		if len(args)+1 < d.services[service].methods[method].argsCount {
			return false
		}
	} else {
		if len(args)+1 != d.services[service].methods[method].argsCount {
			return false
		}
	}

	return true
}

func (d *Dispatcher) verifyArgumentTypes(service, method string, args ...interface{}) bool {
	if d.services[service].methods[method].isVariadic {
		// We verify the constant arguments.
		max := d.services[service].methods[method].argsCount - 1
		for i := 0; i < max-1; i++ {
			if reflect.TypeOf(args[i]) != d.services[service].methods[method].argsTypes[i+1] {
				return false
			}
		}

		// If the variadic arguments are interfaces, we allow any type.
		if d.services[service].methods[method].argsTypes[max].Elem().Kind() == reflect.Interface {
			return true
		}

		// Otherwise, we must verify that each variadic element has the proper type.
		for i := max - 1; i < len(args); i++ {
			if reflect.TypeOf(args[i]) != d.services[service].methods[method].argsTypes[max].Elem() {
				return false
			}
		}
	} else {
		for i, arg := range args {
			if reflect.TypeOf(arg) != d.services[service].methods[method].argsTypes[i+1] {
				return false
			}
		}
	}

	return true
}
