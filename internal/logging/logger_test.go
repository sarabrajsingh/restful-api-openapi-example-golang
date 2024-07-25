package logging

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerLogger_Printf(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	serverLogger := &ServerLogger{Logger: logger}

	serverLogger.Printf("This is a %s message", "test")

	output := buf.String()
	assert.Contains(t, output, "INFO: ")
	assert.Contains(t, output, "This is a test message")
}

func TestServerLogger_Println(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	serverLogger := &ServerLogger{Logger: logger}

	serverLogger.Println("This is a test message")

	output := buf.String()
	assert.Contains(t, output, "INFO: ")
	assert.Contains(t, output, "This is a test message")
}

func TestServerLogger_Fatalf(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		var buf bytes.Buffer
		logger := log.New(&buf, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		serverLogger := &ServerLogger{Logger: logger}

		serverLogger.Fatalf("This is a %s message", "fatal")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestServerLogger_Fatalf")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestServerLogger_Fatal(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		var buf bytes.Buffer
		logger := log.New(&buf, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		serverLogger := &ServerLogger{Logger: logger}

		serverLogger.Fatal("This is a fatal message")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestServerLogger_Fatal")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestNewRealLogger(t *testing.T) {
	logger := NewRealLogger()
	assert.NotNil(t, logger)
}
