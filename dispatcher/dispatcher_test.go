package dispatcher

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type mockService struct {
}

//nolint:unused
func (ms *mockService) unexported() {
}

func (ms *mockService) Exported() {
}

func (ms *mockService) MethodWithArguments(str string, integer int) {
}

func (ms *mockService) MethodWithPtrArguments(ptr *mockService) {
}

func (ms *mockService) MethodWithReturnValue(str string, integer int) (string, int) {
	return str, integer
}

func (ms *mockService) MethodWithPtrArgumentsAndReturnValue(ptr *mockService) *mockService {
	return ptr
}

type mockServiceVariadic struct {
}

func (ms *mockServiceVariadic) MethodWithVariadicArguments(args ...interface{}) {
}

func (ms *mockServiceVariadic) MethodWithConstantAndVariadicArguments(integer int, args ...string) {
}

type mockServiceWithFields struct {
	name string
	age  int
	ptr  *mockService
}

func (ms *mockServiceWithFields) MethodWhichSetsFields(name string, age int, ptr *mockService) {
	ms.age = age
	ms.name = name
	ms.ptr = ptr
}

func (ms *mockServiceWithFields) GetFields() (string, int, *mockService) {
	return ms.name, ms.age, ms.ptr
}

//nolint:thelper
// newDevDispatcher creates a dev dispatcher with the mock service registered by default.
func newDevDispatcher(t *testing.T) *Dispatcher {
	d := New()
	err := d.Register("mock", &mockService{})

	if err != nil {
		t.Fatal(err)
	}

	return d
}

// TestNew creates a new dispatcher and verifies that its service map has been properly built.
func TestNew(t *testing.T) {
	d := New()
	assert.NotNilf(t, d.services, "dispatcher's service map should not be nil")
}

func TestDispatcher_Register(t *testing.T) {
	testCases := []struct {
		name          string
		serviceName   string
		service       interface{}
		success       bool
		expectedError error
	}{
		{
			name:          "Valid service type",
			serviceName:   "mock",
			service:       &mockService{},
			success:       true,
			expectedError: nil,
		},
		{
			name:          "Invalid service type #1",
			serviceName:   "mock",
			service:       mockService{},
			success:       false,
			expectedError: ErrInvalidServiceType,
		},
		{
			name:          "Invalid service type #2",
			serviceName:   "int",
			service:       int(0),
			success:       false,
			expectedError: ErrInvalidServiceType,
		},
		{
			name:          "Invalid service type #3",
			serviceName:   "string",
			service:       "",
			success:       false,
			expectedError: ErrInvalidServiceType,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			err := d.Register(tt.serviceName, tt.service)

			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestDispatcher_Run(t *testing.T) {
	mockPtr := &mockService{}

	testCases := []struct {
		name           string
		serviceName    string
		methodName     string
		arguments      []interface{}
		success        bool
		expectedError  error
		expectedOutput []reflect.Value
		expectPtr      bool
	}{
		{
			name:          "Valid method with no arguments",
			serviceName:   "mock",
			methodName:    "Exported",
			success:       true,
			expectedError: nil,
		},
		{
			name:          "Valid method with constant argument count",
			serviceName:   "mock",
			methodName:    "MethodWithArguments",
			arguments:     []interface{}{"Hello World!", 42},
			success:       true,
			expectedError: nil,
		},
		{
			name:          "Valid method with ptr argument",
			serviceName:   "mock",
			methodName:    "MethodWithPtrArguments",
			arguments:     []interface{}{mockPtr},
			success:       true,
			expectedError: nil,
		},
		{
			name:          "Valid method with ptr argument and return value",
			serviceName:   "mock",
			methodName:    "MethodWithPtrArgumentsAndReturnValue",
			arguments:     []interface{}{mockPtr},
			success:       true,
			expectedError: nil,
			expectedOutput: []reflect.Value{
				reflect.ValueOf(mockPtr),
			},
			expectPtr: true,
		},
		{
			name:        "Valid method with return value",
			serviceName: "mock",
			methodName:  "MethodWithReturnValue",
			arguments: []interface{}{
				"Hello",
				42,
			},
			success:       true,
			expectedError: nil,
			expectedOutput: []reflect.Value{
				reflect.ValueOf("Hello"),
				reflect.ValueOf(42),
			},
		},
		{
			name:          "Non existent service",
			serviceName:   "",
			methodName:    "Exported",
			success:       false,
			expectedError: ErrNonExistentService,
		},
		{
			name:          "Non existent method",
			serviceName:   "mock",
			methodName:    "unexported",
			success:       false,
			expectedError: ErrNonExistentMethod,
		},
		{
			name:          "Too many arguments",
			serviceName:   "mock",
			methodName:    "Exported",
			arguments:     []interface{}{42},
			success:       false,
			expectedError: ErrInvalidArgumentsCount,
		},
		{
			name:          "Not enough arguments",
			serviceName:   "mock",
			methodName:    "MethodWithArguments",
			arguments:     []interface{}{42},
			success:       false,
			expectedError: ErrInvalidArgumentsCount,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			d := newDevDispatcher(t)
			output, err := d.Run(tt.serviceName, tt.methodName, tt.arguments...)

			assert.Equal(t, tt.expectedError, err)
			if tt.expectedOutput != nil {
				for i, item := range output {
					assert.Equal(t, tt.expectedOutput[i].Interface(), item.Interface())
					if tt.expectPtr {
						assert.Equal(t, tt.expectedOutput[i].Pointer(), item.Pointer())
					}
				}
			} else {
				assert.Equal(t, tt.expectedOutput, output)
			}
		})
	}
}

func TestDispatcher_RunVariadic(t *testing.T) {
	testCases := []struct {
		name           string
		serviceName    string
		methodName     string
		arguments      []interface{}
		success        bool
		expectedError  error
		expectedOutput []reflect.Value
		expectPtr      bool
	}{
		{
			name:          "Valid method call with variadic arguments only #1",
			serviceName:   "mock",
			methodName:    "MethodWithVariadicArguments",
			arguments:     []interface{}{"Hello", "World", "!"},
			success:       true,
			expectedError: nil,
		},
		{
			name:          "Valid method call with variadic arguments only #2",
			serviceName:   "mock",
			methodName:    "MethodWithVariadicArguments",
			arguments:     []interface{}{},
			success:       true,
			expectedError: nil,
		},
		{
			name:          "Valid method call with predefined and variadic arguments",
			serviceName:   "mock",
			methodName:    "MethodWithConstantAndVariadicArguments",
			arguments:     []interface{}{42, "Hello", "World", "!"},
			success:       true,
			expectedError: nil,
		},
		{
			name:          "Invalid variadic arguments #1",
			serviceName:   "mock",
			methodName:    "MethodWithConstantAndVariadicArguments",
			arguments:     []interface{}{42, 43},
			success:       false,
			expectedError: ErrInvalidArgumentType,
		},
		{
			name:          "Invalid variadic arguments #2",
			serviceName:   "mock",
			methodName:    "MethodWithConstantAndVariadicArguments",
			arguments:     []interface{}{42, "Hello", 43},
			success:       false,
			expectedError: ErrInvalidArgumentType,
		},
		{
			name:          "Invalid variadic arguments #3",
			serviceName:   "mock",
			methodName:    "MethodWithConstantAndVariadicArguments",
			arguments:     []interface{}{42, "Hello", "World", 43},
			success:       false,
			expectedError: ErrInvalidArgumentType,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			err := d.Register(tt.serviceName, &mockServiceVariadic{})
			assert.NoError(t, err)

			_, err = d.Run(tt.serviceName, tt.methodName, tt.arguments...)

			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestDispatcher_RunSetFields(t *testing.T) {
	testCases := []struct {
		name           string
		serviceName    string
		methodName     string
		arguments      []interface{}
		success        bool
		expectedError  error
		expectedOutput []reflect.Value
		expectedPtr    bool
	}{
		{
			name:          "Can set service fields",
			serviceName:   "mock",
			methodName:    "MethodWhichSetsFields",
			arguments:     []interface{}{"Hello", 42, &mockService{}},
			success:       true,
			expectedError: nil,
			expectedOutput: []reflect.Value{
				reflect.ValueOf("Hello"),
				reflect.ValueOf(42),
				reflect.ValueOf(&mockService{}),
			},
			expectedPtr: false,
		},
		{
			name:          "Can set service fields and keep pointer",
			serviceName:   "mock",
			methodName:    "MethodWhichSetsFields",
			arguments:     []interface{}{"Hello", 42, &mockService{}},
			success:       true,
			expectedError: nil,
			expectedOutput: []reflect.Value{
				reflect.ValueOf("Hello"),
				reflect.ValueOf(42),
				reflect.ValueOf(&mockService{}),
			},
			expectedPtr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			err := d.Register(tt.serviceName, &mockServiceWithFields{})
			assert.NoError(t, err)

			_, err = d.Run(tt.serviceName, tt.methodName, tt.arguments...)

			assert.Equal(t, tt.expectedError, err)

			outputs, err := d.Run(tt.serviceName, "GetFields")
			assert.NoError(t, err)
			for i, output := range outputs {
				assert.Equal(t, tt.expectedOutput[i].Interface(), output.Interface())
			}

			if tt.expectedPtr {
				assert.Equal(t, tt.expectedOutput[2].Pointer(), outputs[2].Pointer())
			}
		})
	}
}

func TestEmptyNamespace(t *testing.T) {
	d := New()
	err := d.Register("", &mockService{})
	assert.NoError(t, err)

	out, err := d.Run("", "MethodWithReturnValue", "test", 42)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(out))
	assert.Equal(t, "test", out[0].String())
	assert.Equal(t, int64(42), out[1].Int())
}
