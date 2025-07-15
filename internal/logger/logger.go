package logger

import (
	"os"
	"bytes"
	"strings"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(args ...any)
	Warn(args ...any)
	Debug(args ...any)
	Error(args ...any)
	Fatal(args ...any)
}

type customTextFormatter struct {}


var logger *logrus.Logger

func GetLogger() *logrus.Logger {
	logger = logrus.New()

	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)
	logger.SetFormatter(&customTextFormatter{})

	defer logger.Info("Logger has been init")
	return logger
}

func (f *customTextFormatter) Format(entry *logrus.Entry) ([]byte, error){
    var b bytes.Buffer
    if len(entry.Data) != 0 {
        b.WriteString(fmt.Sprintf("[%s] [%s %s] %s", 
            strings.ToUpper(entry.Level.String()),
            entry.Data["method"],
            entry.Data["path"],
            entry.Message,
        ))

        for key, value := range entry.Data {
            b.WriteString(fmt.Sprintf(" %s=%v", key, value))
        }
    } else {
        b.WriteString(fmt.Sprintf("[%s] [%s] [%s:%d] %s", 
            strings.ToUpper(entry.Level.String()),
            entry.Caller.File,
            entry.Caller.Function,
            entry.Caller.Line,
            entry.Message,
        ))
    }

    b.WriteString("\n")
    return b.Bytes(), nil
}

func LoggerMiddleware() {

}
