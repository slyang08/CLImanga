package log

import (
	"log"
	"os"
	"runtime"
	"time"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func Init() {
	os.MkdirAll("logs", 0755)
	logFile, err := os.OpenFile("logs/latest.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	Info = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	Info.Printf("===== CLImanga started at %s =====\n", time.Now().Format(time.RFC3339))
}

// Wraps the parameterless function passed in, and automatically records logs before and after execution.
func WrapFunction(fn func()) func() {
	return func() {
		pc, _, _, _ := runtime.Caller(1)
		functionName := "unknown"
		if f := runtime.FuncForPC(pc); f != nil {
			functionName = f.Name()
		}
		Info.Printf("Starting function: %s", functionName)
		fn()
		Info.Printf("Finished function: %s", functionName)
	}
}

func LogFunctionName() {
	pc, _, _, ok := runtime.Caller(1)
	functionName := "unknown"
	if ok {
		funcObj := runtime.FuncForPC(pc)
		if funcObj != nil {
			functionName = funcObj.Name()
		}
	}
	Info.Printf("Starting function %s", functionName)
}
