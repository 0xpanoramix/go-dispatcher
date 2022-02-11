package dispatcher

import (
	"fmt"
	"reflect"
)

var (
	ErrInvalidServiceType    = fmt.Errorf("service must be a pointer to struct")
	ErrNonExistentService    = fmt.Errorf("service is not registered in the dispatcher")
	ErrNonExistentMethod     = fmt.Errorf("method is not registered in the dispatcher")
	ErrInvalidArgumentsCount = fmt.Errorf("invalid number of arguments provided")
	ErrInvalidArgumentType   = fmt.Errorf("invalid argument type")
	ErrMethodWithVariadic    = fmt.Errorf("variadic method are not supported")
)

// funcMetadata represents one method.
// It contains metadata about the function : the method itself,
// its argument's count and types.
type funcMetadata struct {
	function  reflect.Value
	argsCount int
	argsTypes []reflect.Type
}

// serviceData represents a service along with its methods.
type serviceData struct {
	service reflect.Value
	methods map[string]*funcMetadata
}

// Dispatcher contains the services mapping.
// Each service contains its own methods.
type Dispatcher struct {
	services map[string]*serviceData
}

// New creates a new Dispatcher and allocates memory for its service map.
func New() *Dispatcher {
	return &Dispatcher{services: make(map[string]*serviceData)}
}

// Register registers a new service into the dispatcher.
//
// It saves the service's methods in a mapping,
// along with each method's metadata used when calling one of them.
//
// Please refer to funcMetadata for more information
// about the function metadata content.
func (d *Dispatcher) Register(serviceName string, service interface{}) error {
	// The service must be a pointer to struct.
	st := reflect.TypeOf(service)
	if st.Kind() == reflect.Struct || st.Kind() != reflect.Ptr {
		return ErrInvalidServiceType
	}

	// Save the service data locally.
	sd := &serviceData{
		service: reflect.ValueOf(service),
		methods: make(map[string]*funcMetadata),
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
		sd.methods[methodName] = &funcMetadata{
			function:  st.Method(i).Func,
			argsCount: st.Method(i).Func.Type().NumIn(),
			argsTypes: []reflect.Type{},
		}

		// For each method, save its argument's types.
		for j := 0; j < st.Method(i).Func.Type().NumIn(); j++ {
			sd.methods[methodName].argsTypes = append(sd.methods[methodName].argsTypes,
				st.Method(i).Func.Type().In(j))
		}

		// If method has variadic arguments, specify it in the metadata.
		if st.Method(i).Func.Type().IsVariadic() {
			return ErrMethodWithVariadic
		}
	}

	// Save the service and its methods in the dispatcher.
	d.services[serviceName] = sd

	return nil
}

// Run will call the service's method previously registered using the given arguments.
func (d *Dispatcher) Run(service, method string, args ...interface{}) ([]reflect.Value, error) {
	// Checks that the service has been registered.
	if d.services[service] == nil {
		return nil, ErrNonExistentService
	}

	// Checks that the method exists.
	if d.services[service].methods[method] == nil {
		return nil, ErrNonExistentMethod
	}

	// Checks for the given argument list length.
	// A +1 is added to the list length because argsCount counts the service
	// as the first argument to provide to the called method
	if len(args)+1 != d.services[service].methods[method].argsCount {
		return nil, ErrInvalidArgumentsCount
	}

	// Checks for the arguments type.
	for i, arg := range args {
		if reflect.TypeOf(arg) != d.services[service].methods[method].argsTypes[i+1] {
			return nil, ErrInvalidArgumentType
		}
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
