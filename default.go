package trace

import (
	"bytes"
	"io"
	"os"
)

func init() {
	module.RegisterDriver("default", &defaultDriver{})
}

type (
	defaultDriver struct{}

	defaultConnection struct {
		instance       *Instance
		stdout, stderr io.Writer
	}
)

func (d *defaultDriver) Connect(inst *Instance) (Connection, error) {
	return &defaultConnection{
		instance: inst,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
	}, nil
}

func (c *defaultConnection) Open() error  { return nil }
func (c *defaultConnection) Close() error { return nil }

func (c *defaultConnection) Write(spans ...Span) error {
	if len(spans) == 0 {
		return nil
	}
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	for _, span := range spans {
		line := c.instance.Format(span) + "\n"
		if span.Code != 0 || (span.Status != "" && span.Status != StatusOK) {
			_, _ = errBuf.WriteString(line)
		} else {
			_, _ = outBuf.WriteString(line)
		}
	}
	if errBuf.Len() > 0 {
		_, _ = c.stderr.Write(errBuf.Bytes())
	}
	if outBuf.Len() > 0 {
		_, _ = c.stdout.Write(outBuf.Bytes())
	}
	return nil
}
