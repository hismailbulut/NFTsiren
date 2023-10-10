package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Level int32

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelNone
)

func (level Level) String() string {
	switch level {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	}
	return "?"
}

type Logger struct {
	mutex    sync.Mutex
	minLevel Level
	output   io.Writer
}

func (l *Logger) Write(p []byte) (n int, err error) {
	if l.output != nil {
		l.mutex.Lock()
		n, err = l.output.Write(p)
		l.mutex.Unlock()
	}
	return
}

// This is a low level function, do not use this
func (logger *Logger) put(message *Message) {
	logger.mutex.Lock()
	minLevel := logger.minLevel
	logger.mutex.Unlock()
	if message.level < minLevel {
		return
	}
	// Get file name
	_, file, line, ok := runtime.Caller(message.depth)
	if !ok {
		file = "???"
		line = 0
	}
	// Reduce filename length if required
	filename := filepath.Base(file)
	const filenamemaxlen = 19
	linestr := strconv.Itoa(line)
	maxlen := filenamemaxlen - len(linestr)
	if len(filename) > maxlen {
		prefix := ".."
		filename = prefix + filename[len(prefix)+(len(filename)-maxlen):]
		// maxlen -= len(middle)
		// filename = filename[:maxlen/2] + middle + filename[(maxlen/2)+1+len(middle):]
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%-5s %s %-*s -> %s\n",
		message.level,
		time.Now().Format(time.DateTime),
		filenamemaxlen+1, // +1 of filenamemaxlen because of the ":" between file and line
		filename+":"+linestr,
		strings.TrimSpace(message.message),
	)
	if message.fields != nil && len(message.fields) > 0 {
		for k, v := range message.fields {
			fmt.Fprintf(&buf, "\t%s: %v\n", k, v)
		}
	}
	logger.Write(buf.Bytes())
}

var Default Logger = Logger{
	output:   os.Stdout,
	minLevel: LevelDebug,
}

func SetMinLevel(minLevel Level) {
	Default.mutex.Lock()
	defer Default.mutex.Unlock()
	Default.minLevel = minLevel
}

func MinLevel() Level {
	Default.mutex.Lock()
	defer Default.mutex.Unlock()
	return Default.minLevel
}

func SetOutput(output io.Writer) {
	Default.mutex.Lock()
	defer Default.mutex.Unlock()
	Default.output = output
}

type Message struct {
	logger  *Logger
	level   Level
	fields  map[string]any
	message string
	depth   int
}

func (lm Message) Field(key string, value any) Message {
	if lm.fields == nil {
		lm.fields = make(map[string]any)
	}
	lm.fields[key] = value
	lm.depth++
	return lm
}

func (lm Message) Printf(format string, a ...any) {
	lm.message = fmt.Sprintf(format, a...)
	lm.depth++
	lm.logger.put(&lm)
}

func (lm Message) Println(a ...any) {
	lm.message = fmt.Sprintln(a...)
	lm.depth++
	lm.logger.put(&lm)
}

func defaultMessageWithLevel(level Level) Message {
	return Message{
		logger: &Default,
		level:  level,
		depth:  1,
	}
}

func Debug() Message {
	return defaultMessageWithLevel(LevelDebug)
}

func Info() Message {
	return defaultMessageWithLevel(LevelInfo)
}

func Warn() Message {
	return defaultMessageWithLevel(LevelWarn)
}

func Error() Message {
	return defaultMessageWithLevel(LevelError)
}
