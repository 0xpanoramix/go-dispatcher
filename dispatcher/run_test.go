package dispatcher

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			name:          "Valid method with array arguments",
			serviceName:   "mock",
			methodName:    "MethodWithArrayArguments",
			arguments:     []interface{}{[]string{"foo", "bar", "baz"}},
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
			name:           "Valid method with array argument and return value",
			serviceName:    "mock",
			methodName:     "MethodWithArrayArgumentsAndReturnValue",
			arguments:      []interface{}{[]int{1, 2, 3}},
			success:        true,
			expectedError:  nil,
			expectedOutput: []reflect.Value{reflect.ValueOf([]int{1, 2, 3})},
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

func TestDispatcher_RunEmptyNamespace(t *testing.T) {
	d := New()
	err := d.Register("", &mockService{})
	assert.NoError(t, err)

	out, err := d.Run("", "MethodWithReturnValue", "test", 42)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(out))
	assert.Equal(t, "test", out[0].String())
	assert.Equal(t, int64(42), out[1].Int())
}
