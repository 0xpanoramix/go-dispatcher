package dispatcher

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDispatcher_Validate(t *testing.T) {
	testCases := []struct {
		name           string
		success        bool
		serviceName    string
		methodName     string
		param          interface{}
		expectedResult []interface{}
		expectedError  error
	}{
		{
			name:           "no arguments",
			success:        true,
			serviceName:    "mock",
			methodName:     "Exported",
			param:          nil,
			expectedResult: []interface{}{},
			expectedError:  nil,
		},
		{
			name:           "parse one arg : string",
			success:        true,
			param:          "foo",
			serviceName:    "mock",
			methodName:     "MethodWithOneArgumentString",
			expectedResult: []interface{}{"foo"},
			expectedError:  nil,
		},
		{
			name:           "parse one arg : int",
			success:        true,
			param:          4,
			serviceName:    "mock",
			methodName:     "MethodWithOneArgumentInteger",
			expectedResult: []interface{}{4},
			expectedError:  nil,
		},
		{
			name:           "parse one arg : boolean",
			success:        true,
			param:          false,
			serviceName:    "mock",
			methodName:     "MethodWithOneArgumentBoolean",
			expectedResult: []interface{}{false},
			expectedError:  nil,
		},
		{
			name:           "parse one arg : float",
			success:        true,
			param:          float64(2),
			serviceName:    "mock",
			methodName:     "MethodWithOneArgumentFloat",
			expectedResult: []interface{}{float64(2)},
			expectedError:  nil,
		},
		{
			name:           "parse one arg : object",
			success:        true,
			param:          struct{ Foo string }{Foo: "foo"},
			serviceName:    "mock",
			methodName:     "MethodWithOneArgumentObject",
			expectedResult: []interface{}{struct{ Foo string }{Foo: "foo"}},
			expectedError:  nil,
		},
		{
			name:           "parse one arg : array",
			success:        true,
			param:          []string{"foo", "bar", "baz"},
			serviceName:    "mock",
			methodName:     "MethodWithArrayArguments",
			expectedResult: []interface{}{[]string{"foo", "bar", "baz"}},
			expectedError:  nil,
		},
		{
			name:           "parse on arg : type do not match",
			success:        false,
			param:          []interface{}{"test"},
			serviceName:    "mock",
			methodName:     "MethodWithOneArgumentNumber",
			expectedResult: nil,
			expectedError:  errors.New(""),
		},
		{
			name:           "parse multi arg : mix primitive type",
			success:        true,
			serviceName:    "mock",
			methodName:     "MethodWithArguments",
			param:          []interface{}{"foo", 5},
			expectedResult: []interface{}{"foo", 5},
			expectedError:  nil,
		},
		{
			name:           "parse multi arg : type do not match",
			success:        false,
			serviceName:    "mock",
			methodName:     "MethodWithOneArguments",
			param:          []interface{}{true, 4},
			expectedResult: nil,
			expectedError:  errors.New(""),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			err := d.Register(tt.serviceName, &mockService{})
			assert.NoError(t, err)

			res, err := d.Validate(tt.serviceName, tt.methodName, tt.param)
			assert.Equal(t, tt.expectedResult, res)
			assert.IsType(t, tt.expectedError, err)
		})
	}
}
