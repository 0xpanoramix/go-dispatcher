package dispatcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			service:       0,
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
