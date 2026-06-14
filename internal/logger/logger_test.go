package logger

import "testing"

func Test_InfoLog(t *testing.T) {
	Init()
	Info("Info Logs Tested")
}

func Test_DebugLog(t *testing.T) {
	Init()
	Debug("Debug Logs Tested")
}

func Test_ErrorLog(t *testing.T) {
	Init()
	Error("Error Logs Tested")
}
