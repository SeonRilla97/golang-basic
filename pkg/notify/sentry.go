package notify

import (
	"fmt"

	"github.com/getsentry/sentry-go"
)

func SendToSentry(err any, stack string) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetExtra("stack", stack)
		scope.SetLevel(sentry.LevelFatal)

		switch e := err.(type) {
		case error:
			sentry.CaptureException(e)
		default:
			sentry.CaptureMessage(fmt.Sprintf("Panic: %v", e))
		}
	})
}
