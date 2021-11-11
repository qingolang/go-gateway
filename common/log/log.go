package log

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	LEVEL_FLAGS = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

const tunnel_size_default = 1024

// Writer
type Writer interface {
	Init() error
	Write(*Record) error
}

// Rotater
type Rotater interface {
	Rotate() error
	SetPathPattern(string) error
}

// Flusher
type Flusher interface {
	Flush() error
}

// Record
type Record struct {
	time  string
	code  string
	info  string
	level int
}

// String
func (r *Record) String() string {
	return fmt.Sprintf("[%s][%s][%s] %s\n", LEVEL_FLAGS[r.level], r.time, r.code, r.info)
}

// Logger
type Logger struct {
	writers     []Writer
	tunnel      chan *Record
	level       int
	lastTime    int64
	lastTimeStr string
	c           chan bool
	layout      string
	recordPool  *sync.Pool
}

// NewLogger
func NewLogger() *Logger {
	if log_def != nil && !takeup {
		takeup = true //默认启动标志
		return log_def
	}
	l := new(Logger)
	l.writers = []Writer{}
	l.tunnel = make(chan *Record, tunnel_size_default)
	l.c = make(chan bool, 2)
	l.level = DEBUG
	l.layout = "2006/01/02 15:04:05"
	l.recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
	go boostrapLogWriter(l)

	return l
}

// Register
func (l *Logger) Register(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
}

// SetLevel
func (l *Logger) SetLevel(lvl int) {
	l.level = lvl
}

// SetLayout
func (l *Logger) SetLayout(layout string) {
	l.layout = layout
}

// Trace
func (l *Logger) Trace(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(TRACE, fmt, args...)
}

// Debug
func (l *Logger) Debug(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}

// Warn
func (l *Logger) Warn(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(WARNING, fmt, args...)
}

// Info
func (l *Logger) Info(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(INFO, fmt, args...)
}

// Error
func (l *Logger) Error(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, fmt, args...)
}

// Fatal
func (l *Logger) Fatal(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(FATAL, fmt, args...)
}

// Close 关闭
func (l *Logger) Close() {
	close(l.tunnel)
	<-l.c
	for _, w := range l.writers {
		if f, ok := w.(Flusher); ok {
			if err := f.Flush(); err != nil {
				log.Printf("[ERROR] Log close :%v\n", err)
			}
		}
	}
}

// deliverRecordToWriter
func (l *Logger) deliverRecordToWriter(level int, format string, args ...interface{}) {
	var inf, code string

	if level < l.level {
		return
	}

	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
	}

	// source code, file and line num
	_, file, line, ok := runtime.Caller(2)
	if ok {
		code = path.Base(file) + ":" + strconv.Itoa(line)
	}

	// format time
	now := time.Now()
	if now.Unix() != l.lastTime {
		l.lastTime = now.Unix()
		l.lastTimeStr = now.Format(l.layout)
	}
	r := l.recordPool.Get().(*Record)
	r.info = inf
	r.code = code
	r.time = l.lastTimeStr
	r.level = level

	l.tunnel <- r
}

// boostrapLogWriter
func boostrapLogWriter(logger *Logger) {
	if logger == nil {
		panic("logger is nil")
	}

	var (
		r  *Record
		ok bool
	)

	if r, ok = <-logger.tunnel; !ok {
		logger.c <- true
		return
	}

	for _, w := range logger.writers {
		if err := w.Write(r); err != nil {
			log.Printf("[ERROR] logger.writers : %v\n", err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Second * 10)

	for {
		select {
		case r, ok = <-logger.tunnel:
			if !ok {
				logger.c <- true
				return
			}
			for _, w := range logger.writers {
				if err := w.Write(r); err != nil {
					log.Printf("[ERROR] logger.writers : %v\n", err)
				}
			}

			logger.recordPool.Put(r)

		case <-flushTimer.C:
			for _, w := range logger.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						log.Printf("[ERROR] logger.writers : %v\n", err)
					}
				}
			}
			flushTimer.Reset(time.Millisecond * 1000)

		case <-rotateTimer.C:
			for _, w := range logger.writers {
				if r, ok := w.(Rotater); ok {
					if err := r.Rotate(); err != nil {
						log.Printf("[ERROR] logger.Rotater : %v\n", err)
					}
				}
			}
			rotateTimer.Reset(time.Second * 10)
		}
	}
}

// default logger
var (
	log_def *Logger
	takeup  = false
)

// SetLevel
func SetLevel(lvl int) {
	defaultLoggerInit()
	log_def.level = lvl
}

// SetLayout
func SetLayout(layout string) {
	defaultLoggerInit()
	log_def.layout = layout
}

// Trace
func Trace(fmt string, args ...interface{}) {
	defaultLoggerInit()
	log_def.deliverRecordToWriter(TRACE, fmt, args...)
}

// Debug
func Debug(fmt string, args ...interface{}) {
	defaultLoggerInit()
	log_def.deliverRecordToWriter(DEBUG, fmt, args...)
}

// Warn
func Warn(fmt string, args ...interface{}) {
	defaultLoggerInit()
	log_def.deliverRecordToWriter(WARNING, fmt, args...)
}

// Info
func Info(fmt string, args ...interface{}) {
	defaultLoggerInit()
	log_def.deliverRecordToWriter(INFO, fmt, args...)
}

// Error
func Error(fmt string, args ...interface{}) {
	defaultLoggerInit()
	log_def.deliverRecordToWriter(ERROR, fmt, args...)
}

// Fatal
func Fatal(fmt string, args ...interface{}) {
	defaultLoggerInit()
	log_def.deliverRecordToWriter(FATAL, fmt, args...)
}

// Register
func Register(w Writer) {
	defaultLoggerInit()
	log_def.Register(w)
}

// Close
func Close() {
	defaultLoggerInit()
	log_def.Close()
	log_def = nil
	takeup = false
}

// defaultLoggerInit
func defaultLoggerInit() {
	if !takeup {
		log_def = NewLogger()
	}
}
