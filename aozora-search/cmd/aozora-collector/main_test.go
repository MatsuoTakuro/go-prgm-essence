package main

import (
	"log"
	"os"
	"testing"
)

// WARN: you will get a timeout error when you press `run test` on vscode.
func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"for debugging"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func TestMain(m *testing.M) {
	teardown, err := setup()
	if err != nil {
		log.Fatalf("setup failed: %v", err)
	}

	// Run the tests and capture the exit code
	code := m.Run()

	// WARN: if defined as defer func, this teardown func will not be executed
	// cuz the subsequent os.Exit will immediately terminate the program without calling defer func.
	teardown()

	// Exit with the captured code
	os.Exit(code)
}

var testLogfilePath = "./log_test.json"

func setup() (func(), error) {
	var closeLogger func() error
	var err error
	logfilePath = testLogfilePath
	logger, closeLogger, err = setupLogger(logfilePath)
	if err != nil {
		return nil, err
	}

	// Return the teardown function
	return func() {
		if err := closeLogger(); err != nil {
			log.Println("Error during logger's teardown for testing:", err)
			return
		}
		log.Println("Closed logger for testing successfully")
	}, nil
}
