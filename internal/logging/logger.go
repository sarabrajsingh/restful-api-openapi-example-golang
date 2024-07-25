package logging

import (
	"log"
	"os"
)

type Logger interface {
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
}

type ServerLogger struct {
	*log.Logger
}

func (s *ServerLogger) Printf(format string, args ...interface{}) {
	s.Logger.Printf(format, args...)
}

func (s *ServerLogger) Println(args ...interface{}) {
	s.Logger.Println(args...)
}

func (s *ServerLogger) Fatalf(format string, args ...interface{}) {
	s.Logger.Fatalf(format, args...)
}

func (s *ServerLogger) Fatal(args ...interface{}) {
	s.Logger.Fatal(args...)
}

func NewRealLogger() Logger {
	return &ServerLogger{
		Logger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
