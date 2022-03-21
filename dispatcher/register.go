package dispatcher

import (
	"errors"
	"reflect"
)

var (
	ErrInvalidServiceType = errors.New("service must be a pointer to struct")
)

// Register registers a new service into the dispatcher.
//
// It saves the service's methods in a mapping,
// along with each method's metadata used when calling one of them.
//
// Please refer to FuncMetadata for more information
// about the function metadata content.
func (d *Dispatcher) Register(serviceName string, service interface{}) error {
	// The service must be a pointer to struct.
	st := reflect.TypeOf(service)
	if st.Kind() == reflect.Struct || st.Kind() != reflect.Ptr {
		return ErrInvalidServiceType
	}

	// Save the service data locally.
	sd := &ServiceData{
		service: reflect.ValueOf(service),
		methods: make(map[string]*FuncMetadata),
	}
	// Loop on the service's methods.
	for i := 0; i < st.NumMethod(); i++ {
		// Skip unexported methods
		if st.Method(i).PkgPath != "" {
			continue
		}

		// Get the method name.
		methodName := st.Method(i).Name

		// Save each method and the method's argument count.
		sd.methods[methodName] = &FuncMetadata{
			function:   st.Method(i).Func,
			argsCount:  st.Method(i).Func.Type().NumIn(),
			argsTypes:  []reflect.Type{},
			isVariadic: false,
		}

		// For each method, save its argument's types.
		for j := 0; j < st.Method(i).Func.Type().NumIn(); j++ {
			sd.methods[methodName].argsTypes = append(sd.methods[methodName].argsTypes,
				st.Method(i).Func.Type().In(j))
		}

		// If method has variadic arguments, specify it in the metadata.
		if st.Method(i).Func.Type().IsVariadic() {
			sd.methods[methodName].isVariadic = true
		}
	}

	// Save the service and its methods in the dispatcher.
	d.services[serviceName] = sd

	return nil
}
