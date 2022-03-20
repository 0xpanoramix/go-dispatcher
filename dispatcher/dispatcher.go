package dispatcher

import (
	"errors"
	"reflect"
)

var (
	ErrNonExistentService = errors.New("service is not registered in the dispatcher")
	ErrNonExistentMethod  = errors.New("method is not registered in the dispatcher")
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
