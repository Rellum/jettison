package log_test

import (
	"os"
	"testing"

	"github.com/luno/jettison"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"
	"github.com/stretchr/testify/assert"
)

type testLogger struct {
	logs []log.Log
}

func (tl *testLogger) Log(l log.Log) string {
	tl.logs = append(tl.logs, l)

	return ""
}

func TestAddLoggers(t *testing.T) {
	defer log.SetDefaultLoggerForTesting(t, os.Stdout)
	tl := new(testLogger)
	log.SetLogger(tl)

	log.Info(nil, "message", jettison.WithKeyValueString("some", "param"))
	log.Error(nil, errors.New("errMsg"))

	assert.Equal(t, "message,info,some,param,", toStr(tl.logs[0]))
	assert.Equal(t, "errMsg,error,", toStr(tl.logs[1]))
}

func toStr(l log.Log) string {
	str := l.Message + ","
	str += string(l.Level) + ","
	if len(l.Parameters) == 0 {
		return str
	}
	for _, kv := range l.Parameters {
		str += kv.Key + "," + kv.Value + ","
	}
	return str
}
