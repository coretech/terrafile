package stdemuxerhook

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// StdDemuxerHook demuxes logs to io.Writers based on
// severity. By default it uses the following outputs:
// error and higher -> os.Stderr
// warning and lower -> os.Stdout
type StdDemuxerHook struct {
	stdErrLogger *logrus.Logger
	stdOutLogger *logrus.Logger
	level        logrus.Level
}

// New returns a new StdDemuxerHook by silencing the parent
// logger and configuring separate loggers for stderr and
// stdout with the parents loggers properties.
func New(parent *logrus.Logger) *StdDemuxerHook {
	errLogger := logrus.New()
	errLogger.Out = os.Stderr
	errLogger.Level = logrus.DebugLevel
	outLogger := logrus.New()
	outLogger.Out = os.Stdout
	outLogger.Level = logrus.DebugLevel

	// Inherit formatter and level from parent logger
	errLogger.Formatter = parent.Formatter
	outLogger.Formatter = parent.Formatter
	level := parent.Level

	// Make sure parent Logger does not log anything by itself
	parent.Out = ioutil.Discard
	parent.Formatter = &NopFormatter{}

	return &StdDemuxerHook{
		stdErrLogger: errLogger,
		stdOutLogger: outLogger,
		level:        level,
	}
}

// Fire is triggered on new log entries
func (hook *StdDemuxerHook) Fire(entry *logrus.Entry) error {
	if hook.level < entry.Level {
		return nil
	}

	switch entry.Level {
	// stderr
	case logrus.ErrorLevel:
		hook.stdErrLogger.WithFields(entry.Data).Error(entry.Message)
	case logrus.PanicLevel:
		hook.stdErrLogger.WithFields(entry.Data).Panic(entry.Message)
	case logrus.FatalLevel:
		hook.stdErrLogger.WithFields(entry.Data).Fatal(entry.Message)
	// stdout
	case logrus.DebugLevel:
		hook.stdOutLogger.WithFields(entry.Data).Debug(entry.Message)
	case logrus.InfoLevel:
		hook.stdOutLogger.WithFields(entry.Data).Info(entry.Message)
	case logrus.WarnLevel:
		hook.stdOutLogger.WithFields(entry.Data).Warn(entry.Message)
	}

	return nil
}

// Levels returns all levels this hook should be registered to
func (*StdDemuxerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// SetOutput allows to set the info and error level outputs to arbitrary io.Writers
func (hook *StdDemuxerHook) SetOutput(infoLevel io.Writer, errorLevel io.Writer) {
	hook.stdOutLogger.Out = infoLevel
	hook.stdErrLogger.Out = errorLevel
}

// NopFormatter always yields zero bytes and consumes 0 allocs/op.
type NopFormatter struct{}

func (*NopFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte{}, nil
}
