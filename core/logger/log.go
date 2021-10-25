// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file is based on the standard library log file.
package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	traceLog = iota
	debugLog
	infoLog
	warnLog
	errorLog
	fatalLog
	panicLog
)

var messages = [7][]byte{
	[]byte(" [TRACE] "),
	[]byte(" [DEBUG] "),
	[]byte(" [INFO] "),
	[]byte(" [WARN] "),
	[]byte(" [ERROR] "),
	[]byte(" [FATAL] "),
	[]byte(" [PANIC] "),
}

var stderr io.Writer = os.Stderr
var defaultLogger logger

type logger struct {
	mu     sync.Mutex
	out    io.Writer
	buffer []byte
}

func init() {
	defaultLogger = logger{out: stderr, buffer: make([]byte, 1024)}
}

func MustOrPanic(err error) {
	if err != nil {

		panic(err)
	}
}

func LogError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (l logger) output(level int, msg string) error {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buffer = l.buffer[:0]
	appendTime(now, &l.buffer)
	l.buffer = append(l.buffer, messages[level]...)
	return nil
}

func Debug(v ...interface{}) {

}

func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func appendTime(t time.Time, buffer *[]byte) {
	year, month, day := t.Date()
	itoa(buffer, year, 4)
	*buffer = append(*buffer, '-')
	itoa(buffer, int(month), 2)
	*buffer = append(*buffer, '-')
	itoa(buffer, day, 2)
	*buffer = append(*buffer, 'T')
	hour, min, sec := t.Clock()
	itoa(buffer, hour, 2)
	*buffer = append(*buffer, ':')
	itoa(buffer, min, 2)
	*buffer = append(*buffer, ':')
	itoa(buffer, sec, 2)
}
