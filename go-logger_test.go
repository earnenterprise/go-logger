package logger

import (
	"testing"
)

func TestErrors(t *testing.T) {
	ApplicationTitle = "GO-Logger Testing"
	OutputDebug = true
	OutputInfo = true
	OutputConsole = true
	OutputFilename = "go-logger_test.txt"
	OutputJSON = "go-logger_test.json"
	OutputAsReadable = true
	ErrorWebhook = "https://outlook.office.com/webhook/aed5d6c5-1b3b-4e67-8ba8-ecf303be311b@fbc33648-70de-450d-8c7b-5d0ad852cec0/IncomingWebhook/c88a43ae0af3481d8dec54773a772651/79dd5408-c1a2-43ec-960f-37dae79678cd"
	// IgnoredFiles = append(IgnoredFiles, "go-logger_test.go")

	ClearLog()

	Printf("Testing printf message, should be info")
	Debug("Testing debug message")
	Info("Testing info message")
	Warning("Testing warning message")
	Error("Testing error message")
	Fatalf("Testing fatal message")

	Info(struct{ Object string }{Object: "123"})
}
