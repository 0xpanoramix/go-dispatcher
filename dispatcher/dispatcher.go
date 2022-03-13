package dispatcher

import (
	"errors"
	"reflect"
)

var (
	ErrInvalidServiceType    = errors.New("service must be a pointer to struct")
	ErrNonExistentService    = errors.New("service is not registered in the dispatcher")
	ErrNonExistentMethod     = errors.New("method is not registered in the dispatcher")
	ErrInvalidArgumentsCount = errors.New("invalid number of arguments provided")
	ErrInvalidArgumentType   = errors.New("invalid argument type")
)

// FuncMetadata represents one method.
// It contains metadata about the function : the method itself,
// its argument's count and types.
type FuncMetadata struct {
	function   reflect.Value
	argsCount  int
	argsTypes  []reflect.Type
	isVariadic bool
}

func (f *FuncMetadata) GetFunction() reflect.Value {
	return f.function
}

func (f *FuncMetadata) GetArgsCount() int {
	return f.argsCount
}

func (f *FuncMetadata) GetArgsTypes() []reflect.Type {
	return f.argsTypes
}

func (f *FuncMetadata) IsVariadic() bool {
	return f.isVariadic
}

// ServiceData represents a service along with its methods.
type ServiceData struct {
	service reflect.Value
	methods map[string]*FuncMetadata
}

// Dispatcher contains the services mapping.
// Each service contains its own methods.
type Dispatcher struct {
	services map[string]*ServiceData
}

// New creates a new Dispatcher and allocates memory for its service map.
func New() *Dispatcher {
	return &Dispatcher{services: make(map[string]*ServiceData)}
}

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

// getService return a registered service with his methods
func (d *Dispatcher) getService(service string) (*ServiceData, error) {
	s, ok := d.services[service]
	if !ok {
		return nil, ErrNonExistentService
	}

	return s, nil
}

// GetMethod return a method from a service
func (d *Dispatcher) GetMethod(service, method string) (*FuncMetadata, error) {
	s, err := d.getService(service)
	if err != nil {
		return nil, err
	}

	// Retrieve method.
	m, ok := s.methods[method]
	if !ok {
		return nil, ErrNonExistentMethod
	}

	return m, nil
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
