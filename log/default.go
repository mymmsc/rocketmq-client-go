package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// These flags define which text to prefix to each log entry generated by the DefaultLogger.
const (
	// Bits or'ed together to control what's printed. There is no control over the
	// order they appear (the order listed here) or the format they present (as
	// described in the comments).  A colon appears after these items:
	//      2009/0123 01:23:23.123123 /a/b/c/d.go:23: message
	Ldate         = 1 << iota                     // the date: 2009/0123
	Ltime                                         // the time: 01:23:23
	Lmicroseconds                                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                                    // final file name element and line number: d.go:23. overrides Llongfile
	Lmodule                                       // module name
	Llevel                                        // level: 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
	LstdFlags     = Ldate | Ltime | Lmicroseconds // initial values for the standard logger
	Ldefault      = Lmodule | Llevel | Lshortfile | LstdFlags
) // [prefix][time][level][module][shortfile|longfile]

// log level
const (
	Ldebug = iota
	Linfo
	Lwarn
	Lerror
	Lpanic
	Lfatal
)

var levels = []string{
	"[DEBUG]",
	"[INFO]",
	"[WARN]",
	"[ERROR]",
	"[PANIC]",
	"[FATAL]",
}

// A DefaultLogger represents an active logging object that generates lines of
// output to an io.Writer.  Each logging operation makes a single call to
// the Writer's Write method.  A DefaultLogger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type DefaultLogger struct {
	mu         sync.Mutex // ensures atomic writes; protects the following fields
	prefix     string     // prefix to write at beginning of each line
	flag       int        // properties
	Level      int
	out        io.Writer    // destination for output
	buf        bytes.Buffer // for accumulating text to write
	levelStats [6]int64
}

// New creates a new DefaultLogger.   The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.
func New(out io.Writer, prefix string, flag int) *DefaultLogger {
	return &DefaultLogger{out: out, prefix: prefix, Level: 1, flag: flag}
}

var Std = New(os.Stderr, "", Ldefault)

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
// Knows the buffer has capacity.
func itoa(buf *bytes.Buffer, i int, wid int) {
	var u uint = uint(i)
	if u == 0 && wid <= 1 {
		buf.WriteByte('0')
		return
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}

	// avoid slicing b to avoid an allocation.
	for bp < len(b) {
		buf.WriteByte(b[bp])
		bp++
	}
}

func shortFile(file string, flag int) string {
	sep := "/"
	if (flag & Lmodule) != 0 {
		sep = "/src/"
	}
	pos := strings.LastIndex(file, sep)
	if pos != -1 {
		return file[pos+len(sep):]
	}
	return file
}

func (l *DefaultLogger) formatHeader(buf *bytes.Buffer, t time.Time, file string, line int, lvl int, reqId string) {
	if l.prefix != "" {
		buf.WriteString(l.prefix)
	}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			buf.WriteByte('/')
			itoa(buf, int(month), 2)
			buf.WriteByte('/')
			itoa(buf, day, 2)
			buf.WriteByte(' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			buf.WriteByte(':')
			itoa(buf, min, 2)
			buf.WriteByte(':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				buf.WriteByte('.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			buf.WriteByte(' ')
		}
	}
	if reqId != "" {
		buf.WriteByte('[')
		buf.WriteString(reqId)
		buf.WriteByte(']')
	}
	if l.flag&Llevel != 0 {
		buf.WriteString(levels[lvl])
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			file = shortFile(file, l.flag)
		}
		buf.WriteByte(' ')
		buf.WriteString(file)
		buf.WriteByte(':')
		itoa(buf, line, -1)
		buf.WriteString(": ")
	}
}

// Output writes the output for a logging event.  The string s contains
// the text to print after the prefix specified by the flags of the
// DefaultLogger.  A newline is appended if the last character of s is not
// already a newline.  Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *DefaultLogger) Output(reqId string, lvl int, calldepth int, s string) error {
	if lvl < l.Level {
		return nil
	}
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile|Lmodule) != 0 {
		// release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.levelStats[lvl]++
	l.buf.Reset()
	l.formatHeader(&l.buf, now, file, line, lvl, reqId)
	l.buf.WriteString(s)
	if len(s) > 0 && s[len(s)-1] != '\n' {
		l.buf.WriteByte('\n')
	}
	_, err := l.out.Write(l.buf.Bytes())
	return err
}

// -----------------------------------------

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *DefaultLogger) Printf(format string, v ...interface{}) {
	l.Output("", Linfo, 2, fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *DefaultLogger) Print(v ...interface{}) { l.Output("", Linfo, 2, fmt.Sprint(v...)) }

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *DefaultLogger) Println(v ...interface{}) { l.Output("", Linfo, 2, fmt.Sprintln(v...)) }

// -----------------------------------------

func (l *DefaultLogger) Debugf(format string, v ...interface{}) {
	if Ldebug < l.Level {
		return
	}
	l.Output("", Ldebug, 2, fmt.Sprintf(format, v...))
}

func (l *DefaultLogger) Debug(v ...interface{}) {
	if Ldebug < l.Level {
		return
	}
	l.Output("", Ldebug, 2, fmt.Sprintln(v...))
}

// -----------------------------------------

func (l *DefaultLogger) Infof(format string, v ...interface{}) {
	if Linfo < l.Level {
		return
	}
	l.Output("", Linfo, 2, fmt.Sprintf(format, v...))
}

func (l *DefaultLogger) Info(v ...interface{}) {
	if Linfo < l.Level {
		return
	}
	l.Output("", Linfo, 2, fmt.Sprintln(v...))
}

// -----------------------------------------

func (l *DefaultLogger) Warnf(format string, v ...interface{}) {
	l.Output("", Lwarn, 2, fmt.Sprintf(format, v...))
}

func (l *DefaultLogger) Warn(v ...interface{}) { l.Output("", Lwarn, 2, fmt.Sprintln(v...)) }

// -----------------------------------------

func (l *DefaultLogger) Errorf(format string, v ...interface{}) {
	l.Output("", Lerror, 2, fmt.Sprintf(format, v...))
}

func (l *DefaultLogger) Error(v ...interface{}) { l.Output("", Lerror, 2, fmt.Sprintln(v...)) }

// -----------------------------------------

func (l *DefaultLogger) Fatal(v ...interface{}) {
	l.Output("", Lfatal, 2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *DefaultLogger) Fatalf(format string, v ...interface{}) {
	l.Output("", Lfatal, 2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *DefaultLogger) Fatalln(v ...interface{}) {
	l.Output("", Lfatal, 2, fmt.Sprintln(v...))
	os.Exit(1)
}

// -----------------------------------------

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *DefaultLogger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output("", Lpanic, 2, s)
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *DefaultLogger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output("", Lpanic, 2, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *DefaultLogger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output("", Lpanic, 2, s)
	panic(s)
}

// -----------------------------------------

func (l *DefaultLogger) Stack(v ...interface{}) {
	s := fmt.Sprint(v...)
	s += "\n"
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, true)
	s += string(buf[:n])
	s += "\n"
	l.Output("", Lerror, 2, s)
}

func (l *DefaultLogger) SingleStack(v ...interface{}) {
	s := fmt.Sprint(v...)
	s += "\n"
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, false)
	s += string(buf[:n])
	s += "\n"
	l.Output("", Lerror, 2, s)
}

// -----------------------------------------

func (l *DefaultLogger) Stat() (stats []int64) {
	l.mu.Lock()
	v := l.levelStats
	l.mu.Unlock()
	return v[:]
}

// Flags returns the output flags for the logger.
func (l *DefaultLogger) Flags() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flag
}

// SetFlags sets the output flags for the logger.
func (l *DefaultLogger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
}

// Prefix returns the output prefix for the logger.
func (l *DefaultLogger) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.prefix
}

// SetPrefix sets the output prefix for the logger.
func (l *DefaultLogger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}
