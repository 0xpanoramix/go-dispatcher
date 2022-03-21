package dispatcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockService struct {
}

//nolint:unused
func (ms *mockService) unexported() {}

func (ms *mockService) Exported() {}

func (ms *mockService) MethodWithOneArgumentString(_ string) {}

func (ms *mockService) MethodWithOneArgumentInteger(_ int) {}

func (ms *mockService) MethodWithOneArgumentBoolean(_ bool) {}

func (ms *mockService) MethodWithOneArgumentFloat(_ float64) {}

func (ms *mockService) MethodWithOneArgumentObject(_ struct{ Foo string }) {}

func (ms *mockService) MethodWithArguments(_ string, _ int) {}

func (ms *mockService) MethodWithPtrArguments(_ *mockService) {}

func (ms *mockService) MethodWithArrayArguments(_ []string) {}

func (ms *mockService) MethodWithReturnValue(str string, integer int) (string, int) {
	return str, integer
}

func (ms *mockService) MethodWithPtrArgumentsAndReturnValue(ptr *mockService) *mockService {
	return ptr
}

func (ms *mockService) MethodWithArrayArgumentsAndReturnValue(integers []int) []int {
	return integers
}

type mockServiceVariadic struct{}

func (ms *mockServiceVariadic) MethodWithVariadicArguments(_ ...interface{}) {}

func (ms *mockServiceVariadic) MethodWithConstantAndVariadicArguments(_ int, _ ...string) {}

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
