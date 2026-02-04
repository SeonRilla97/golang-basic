package sentry

import (
	"gorm-test/pkg/apperror"

	"github.com/getsentry/sentry-go"
)

func Init(dsn string) error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      "production",
		TracesSampleRate: 0.1,
	})
}

func CaptureError(err error) {
	if appErr, ok := err.(*apperror.apperror); ok {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("error_code", appErr.Code)
			scope.SetExtra("context", appErr.Context)
			sentry.CaptureException(appErr)
		})
	} else {
		sentry.CaptureException(err)
	}
}
