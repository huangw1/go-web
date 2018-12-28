package log

import (
	"log"
	"fmt"
	"os"
	"github.com/fatih/color"
)

const calldepth = 3

type defaultLogger struct {
	*log.Logger
}

func (l *defaultLogger) Debug(v ...interface{})  {
	l.Output(calldepth, combine("DEBUG", fmt.Sprint(v)))
}

func (l *defaultLogger) Debugf(format string, v ...interface{})  {
	l.Output(calldepth, combine("DEBUG", fmt.Sprintf(format, v)))
}

func (l *defaultLogger) Info(v ...interface{})  {
	l.Output(calldepth, combine(color.GreenString("INFO"), fmt.Sprint(v)))
}

func (l *defaultLogger) Infof(format string, v ...interface{})  {
	l.Output(calldepth, combine(color.GreenString("INFO"), fmt.Sprintf(format, v)))
}

func (l *defaultLogger) Warn(v ...interface{})  {
	l.Output(calldepth, combine(color.YellowString("WARN"), fmt.Sprint(v)))
}

func (l *defaultLogger) Warnf(format string, v ...interface{})  {
	l.Output(calldepth, combine(color.YellowString("WARN"), fmt.Sprintf(format, v)))
}

func (l *defaultLogger) Error(v ...interface{})  {
	l.Output(calldepth, combine(color.RedString("ERROR"), fmt.Sprint(v)))
}

func (l *defaultLogger) Errorf(format string, v ...interface{})  {
	l.Output(calldepth, combine(color.RedString("ERROR"), fmt.Sprintf(format, v)))
}

func (l *defaultLogger) Fatal(v ...interface{})  {
	l.Output(calldepth, combine(color.MagentaString("FATAL"), fmt.Sprint(v)))
	os.Exit(1)
}

func (l *defaultLogger) Fatalf(format string, v ...interface{})  {
	l.Output(calldepth, combine(color.MagentaString("FATAL"), fmt.Sprintf(format, v)))
	os.Exit(1)
}

func (l *defaultLogger) Panic(v ...interface{})  {
	l.Panic(v)
}

func (l *defaultLogger) Panicf(format string, v ...interface{})  {
	l.Panicf(format, v)
}


func combine(level, msg string) string {
	return fmt.Sprintf("%s: %s", level, msg)
}
