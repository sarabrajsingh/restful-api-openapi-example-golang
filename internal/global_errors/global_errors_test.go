package global_errors_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sarabrajsingh/restful-openapi/internal/global_errors"
	"github.com/sarabrajsingh/restful-openapi/mocks"
	"github.com/stretchr/testify/assert"
)

func TestErrorStore_AddError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)
	es := global_errors.NewErrorStore()

	// Expect log calls with formatted strings
	mockLogger.EXPECT().Printf("appending [%s] to errorBuffer", "Error 1").Times(1)
	mockLogger.EXPECT().Printf("appending [%s] to errorBuffer", "Error 2").Times(1)

	es.AddError(mockLogger, "Error 1")
	es.AddError(mockLogger, "Error 2")

	// Check the buffer content
	errors := es.GetErrors(mockLogger)
	assert.Len(t, errors, 2)
	assert.Equal(t, "Error 1", errors[0])
	assert.Equal(t, "Error 2", errors[1])
}

func TestErrorStore_DeleteErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)
	es := global_errors.NewErrorStore()

	mockLogger.EXPECT().Printf("appending [%s] to errorBuffer", "Error 1").Times(1)
	// Add some errors
	es.AddError(mockLogger, "Error 1")

	// Expect log call for deletion
	mockLogger.EXPECT().Printf("Successfully cleared the errors buffer").Times(1)

	// Clear the errors
	es.DeleteErrors(mockLogger)

	// Check the buffer content
	errors := es.GetErrors(mockLogger)
	assert.Len(t, errors, 0)
}

func TestErrorStore_BufferOverflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)
	es := global_errors.NewErrorStore()

	// Expect multiple adds to fill the buffer and one additional to trigger overflow
	for i := 0; i < global_errors.MaxErrorBufferSize; i++ {
		mockLogger.EXPECT().Printf("appending [%s] to errorBuffer", "Error").Times(1)
		es.AddError(mockLogger, "Error")
	}

	// The last overflow error
	mockLogger.EXPECT().Printf("appending [%s] to errorBuffer", "Overflow Error").Times(1)
	mockLogger.EXPECT().Printf("error buffer overflow; reset array").Times(1)
	es.AddError(mockLogger, "Overflow Error")

	// Check the buffer content
	errors := es.GetErrors(mockLogger)
	assert.Len(t, errors, global_errors.MaxErrorBufferSize)
	assert.Equal(t, "Overflow Error", errors[global_errors.MaxErrorBufferSize-1])
}
