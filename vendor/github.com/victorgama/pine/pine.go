// Pine is a completely useless (but cute) logging interface
package pine

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
)

type msgType int

type pineMsg struct {
	t      msgType
	at     time.Time
	module string
	extra  *string
	msg    string
	params []interface{}
}

// PineWriter is a writter instance already associated to a module.
type PineWriter struct {
	parent *Pine
	name   string
}

// PineExtraWriter is a writer instance associated to a module and with
// a static Extra field.
type PineExtraWriter struct {
	parent *PineWriter
	extra  string
}

// Pine is a completely useless (but cute) logging interface
type Pine struct {
	timeProvider   func() time.Time
	outputProvider func(msg string)
	formatProvider func(at time.Time, t msgType, module string, extra *string, msg string, params ...interface{}) string
}

func (p *Pine) write(t msgType, module string, extra *string, msg string, params ...interface{}) {
	at := p.timeProvider()
	p.outputProvider(p.formatProvider(at, t, module, extra, msg, params...))
}

// NewWriter creates a new writer instance with a given module name
func (p *Pine) NewWriter(module string) *PineWriter {
	return &PineWriter{
		parent: p,
		name:   module,
	}
}

// WithExtra returns a new PineExtraWriter with an associated module and
// static extra value
func (p *PineWriter) WithExtra(extra string) *PineExtraWriter {
	return &PineExtraWriter{
		parent: p,
		extra:  extra,
	}
}

var stdout = colorable.NewColorableStdout()

var pine = &Pine{
	timeProvider:   func() time.Time { return time.Now() },
	outputProvider: func(msg string) { fmt.Fprint(stdout, msg) },
}

func ttyFormatProvider(at time.Time, t msgType, module string, extra *string, msg string, params ...interface{}) string {
	prefix := fmt.Sprintf("%s %s  %s ", aurora.Gray(at.Format("15:04:05")), typeEmoji[t], aurora.Magenta(module))
	ex := ""
	if extra != nil {
		ex = fmt.Sprintf("%s ", aurora.Cyan(*extra))
	}
	suffix := fmt.Sprintf(msg, params...)
	return fmt.Sprintf("%s%s%s\n", prefix, ex, suffix)
}

func basicFormatProvider(at time.Time, t msgType, module string, extra *string, msg string, params ...interface{}) string {
	prefix := fmt.Sprintf("%s level=%s module=%s ", at, strings.ToLower(typeMap[t]), module)
	ex := ""
	if extra != nil {
		ex = fmt.Sprintf("extra=%s ", *extra)
	}
	suffix := fmt.Sprintf(msg, params...)
	return fmt.Sprintf("%s%s%s\n", prefix, ex, suffix)
}

func init() {
	stat, _ := os.Stdout.Stat()
	if stat.Mode()&os.ModeNamedPipe == os.ModeNamedPipe {
		pine.formatProvider = basicFormatProvider
	} else {
		pine.formatProvider = ttyFormatProvider
	}
}

// NewWriter creates a new Writer instance using the provided module
// name
func NewWriter(module string) *PineWriter {
	return pine.NewWriter(module)
}

//go:generate go run generators/type-generator.go
