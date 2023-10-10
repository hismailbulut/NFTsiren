package log

import (
	"bytes"
	stdlog "log"
	"testing"
)

func BenchmarkLog(b *testing.B) {
	var output bytes.Buffer
	SetOutput(&output)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info().Println("Hello world!")
	}
}

func BenchmarkLogField(b *testing.B) {
	var output bytes.Buffer
	SetOutput(&output)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info().Field("field", 10).Println("Hello world!")
	}
}

func BenchmarkStdLog(b *testing.B) {
	var output bytes.Buffer
	stdlog.SetOutput(&output)
	stdlog.SetFlags(stdlog.Ldate | stdlog.Ltime | stdlog.Lshortfile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stdlog.Println("Hello world!")
	}
}

func BenchmarkDisabledLog(b *testing.B) {
	SetMinLevel(LevelWarn)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info().Println("Hello world!")
	}
}
