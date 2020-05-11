package main

import (
	"time"

	"github.com/BTBurke/cannon"
	"go.uber.org/zap"
)

func main() {
	// Using zap's preset constructors is the simplest way to get a feel for the
	// package, but they don't allow much customization.
	logger := zap.NewExample() // or NewProduction, or NewDevelopment
	defer logger.Sync()

	const url = "http://example.com"

	// In most circumstances, use the SugaredLogger. It's 4-10x faster than most
	// other structured logging packages and has a familiar, loosely-typed API.
	sugar := logger.Sugar()
	sugar.Infow("Failed to fetch URL.",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	// In the unusual situations where every microsecond matters, use the
	// Logger. It's even faster than the SugaredLogger, but only supports
	// structured logging.
	logger.Info("Failed to fetch URL.",
		// Structured context as strongly typed fields.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	log, _ := cannon.NewDevelopment()
	start := time.Now()

	// you can log as normal using any zap methods, such as adding common logging fields
	logger1 := log.With(
		zap.String("request_id", "001"),
	)

	// you can add additional fields at each logging call
	logger1.Info("auth success", zap.String("auth_role", "user_rw"))

	// the logger can be passed along in a context to handlers and other services
	// ctx := cannon.CtxLogger(context.Background(), logger1)
	// requestHandler(ctx, req, resp)

	// when finished with this request, call cannon.Emit (with optional additional fields) to log
	// a single wide log line with every field added throughout the entire request

	// cannonical log lines make it easy to gather all of the relevant context for each request in one place
	// and allow you to aggregate statistics across requests for a better view of how your application is performing
	cannon.Emit(logger1, zap.Duration("request_duration", time.Now().Sub(start)))
}
