package log

import (
	"sync"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/os/glog"
)

var _logger *glog.Logger
var _loggerOnce sync.Once

const logName = "cnbs"

// init init
func init() {
	if err := NewLogger().SetConfigWithMap(g.Map{
		// "path":     "/var/log",
		"level":    "all",
		"stdout":   true,
		"StStatus": 0,
		"project":  logName,
	}); err != nil {
		panic(err)
	}
}

// NewLogger logger
func NewLogger() *glog.Logger {
	if _logger != nil {
		return _logger
	}
	_loggerOnce.Do(func() {
		_logger = glog.New()

	})
	return _logger
}

// Warn warn
func Warn(args ...interface{}) {
	NewLogger().Warning(args...)
}

// Info info
func Info(args ...interface{}) {
	NewLogger().Info(args...)
}

// Info info
func InfoF(f string, args ...interface{}) {
	NewLogger().Infof(f, args...)
}

// Debug dev
func Debug(args ...interface{}) {
	NewLogger().Debug(args...)
}

// DebugF dev
func DebugF(f string, args ...interface{}) {
	NewLogger().Debugf(f, args...)
}

// // Trace dev
// func Trace(args ...interface{}) {
// 	NewLogger().(args...)
// }

// Fatal dev
func Fatal(args ...interface{}) {
	NewLogger().Fatal(args...)
}

// Error dev
func Error(args ...interface{}) {
	NewLogger().Error(args...)
}

// ErrorF ErrorF
func ErrorF(f string, args ...interface{}) {
	NewLogger().Errorf(f, args...)
}
