package config

import (
	"log"
	"os"
	"strings"

	"github.com/hashicorp/logutils"
)

// logLevel stores the current log level, it cannot be changed.
var logLevel string

// init configures the log filter (provided by hashicorp/logutils) with a suitable level (using environment variable GODOG_LOGLEVEL).
func init() {
	//look for a log level env setting.  Start with GODOG
	level, b := os.LookupEnv("GODOG_LOGLEVEL")	
	if !b {
		//also look for standard LOGLEVEL
		level, b = os.LookupEnv("LOGLEVEL")
		if !b {
			//default to error
			level = "ERROR"
		}
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "NOTICE", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(level),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}

// CurrentLogLevel returns the current log level. It cannot be changed.
func CurrentLogLevel() string {
	return logLevel
}

// LineBreakReplacer replaces carriage return (\r), linefeed (\n), formfeed (\f) and other similar characters with a space.
func LineBreakReplacer(s string) string {

	const space = " "
	return strings.NewReplacer(
		"\r\n", space,
		"\r", space,
		"\n", space,
		"\v", space, // vertical tab
		"\f", space,
		"\u0085", space, // Unicode 'NEXT LINE (NEL)'
		"\u2028", space, // Unicode 'LINE SEPARATOR'
		"\u2029", space, // Unicode 'PARAGRAPH SEPARATOR'
	).Replace(s)
}
