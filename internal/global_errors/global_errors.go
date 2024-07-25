package global_errors

import (
	"sync"

	"github.com/sarabrajsingh/restful-openapi/internal/logging"
)

const MaxErrorBufferSize = 512

type ErrorStore interface {
	DeleteErrors(logging.Logger)
	GetErrors(logging.Logger) []string
	AddError(logging.Logger, string)
}

type errorStoreImpl struct {
	errorBuffer []string
	mutex       *sync.Mutex
}

func NewErrorStore() ErrorStore {
	return &errorStoreImpl{
		errorBuffer: make([]string, 0),
		mutex:       &sync.Mutex{},
	}
}

func (es *errorStoreImpl) AddError(log logging.Logger, errorMessage string) {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	if len(es.errorBuffer) >= MaxErrorBufferSize {
		log.Printf("error buffer overflow; reset array")
		es.errorBuffer = es.errorBuffer[1:]
	}

	log.Printf("appending [%s] to errorBuffer", errorMessage)
	es.errorBuffer = append(es.errorBuffer, errorMessage)
}

func (es *errorStoreImpl) GetErrors(log logging.Logger) []string {
	es.mutex.Lock()
	defer es.mutex.Unlock()
	return append([]string{}, es.errorBuffer...)
}

func (es *errorStoreImpl) DeleteErrors(log logging.Logger) {
	es.mutex.Lock()
	defer es.mutex.Unlock()
	es.errorBuffer = make([]string, 0)
	log.Printf("Successfully cleared the errors buffer")
}
