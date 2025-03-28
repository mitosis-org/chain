package cmd

import (
	"fmt"
	"time"

	"github.com/omni-network/omni/lib/errors"

	"cosmossdk.io/log"
	"github.com/getsentry/sentry-go"
)

var _ log.Logger = (*SentryLogger)(nil)

type SentryLogger struct {
	logger log.Logger
}

func NewSentryLogger(logger log.Logger) log.Logger {
	return SentryLogger{logger}
}

func InitSentry(config *SentryConfig, logger log.Logger) (log.Logger, error) {
	if config.DSN == "" {
		return logger, nil
	}

	environment := config.Environment
	if environment == "" {
		environment = "localnet"
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         config.DSN,
		Environment: environment,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize sentry")
	}

	return NewSentryLogger(logger), nil
}

func (l SentryLogger) HandleSentryMessage(level sentry.Level, msg string, keyVals ...any) {
	hub := sentry.CurrentHub().Clone()
	evt := sentry.NewEvent()

	tags := map[string]string{}
	for i := 1; i < len(keyVals); i += 2 {
		key := fmt.Sprintf("%v", keyVals[i-1])
		val := fmt.Sprintf("%v", keyVals[i])
		tags[key] = val
	}

	evt.Message = msg
	evt.Timestamp = time.Now()
	evt.Tags = tags
	evt.Level = level

	hub.CaptureEvent(evt)
}

func (l SentryLogger) Info(msg string, keyVals ...any) {
	l.logger.Info(msg, keyVals...)
	// l.HandleSentryMessage(sentry.LevelInfo, msg, keyVals)
}

func (l SentryLogger) Warn(msg string, keyVals ...any) {
	l.logger.Warn(msg, keyVals...)
	l.HandleSentryMessage(sentry.LevelWarning, msg, keyVals)
}

func (l SentryLogger) Error(msg string, keyVals ...any) {
	l.logger.Error(msg, keyVals...)
	l.HandleSentryMessage(sentry.LevelError, msg, keyVals)
}

func (l SentryLogger) Debug(msg string, keyVals ...any) {
	l.logger.Debug(msg, keyVals...)
	// l.HandleSentryMessage(sentry.LevelDebug, msg, keyVals)
}

func (l SentryLogger) With(keyVals ...any) log.Logger {
	logger := l.logger.With(keyVals...)
	return SentryLogger{logger}
}

func (l SentryLogger) Impl() any {
	return l.logger.Impl()
}
