package main

import (
	"go.uber.org/dig"
)

// func main() { // coverage-ignore
// 	port := flag.Int("port", 8080, "Port to listen on")
// 	jsonLogs := flag.Bool("json-logs", false, "Indicates if logs should be in JSON format or text (default)")
// 	logLevel := flag.String("log-level", slog.LevelDebug.String(), "Log level can be DEBUG, INFO, WARN and ERROR")
// 	noop := flag.Bool("noop", false, "Do not start. Just setup deps and exit. Useful for testing if setup is all working.")
// 	flag.Parse()

// 	cfg := lo.Must(config.Load())

// 	var logLevelVal slog.Level
// 	lo.Must0(logLevelVal.UnmarshalText([]byte(*logLevel)))
// 	rootLogger := diag.SetupRootLogger(
// 		diag.NewRootLoggerOpts().
// 			WithJSONLogs(*jsonLogs).
// 			WithLogLevel(logLevelVal),
// 	)
// 	run(runOpts{
// 		rootLogger:     rootLogger,
// 		noopHTTPListen: *noop,
// 		cfg:            cfg,
// 	})
// }

func main() {
	container := dig.New()
	rootCmd := newRootCmd(container)
	rootCmd.AddCommand(
		newHTTPServerCmd(container),
	)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
