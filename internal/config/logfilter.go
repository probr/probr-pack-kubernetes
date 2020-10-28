package config

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/logutils"
)

// init configures the log filter (provided by hashicorp/logutils) with a suitable level (using environment variable GODOG_LOGLEVEL).
func init() {
	//look for a log level env setting.  Start with GODOG
	level, isPresent := os.LookupEnv("GODOG_LOGLEVEL")
	if !isPresent {
		//also look for standard LOGLEVEL
		level, isPresent = os.LookupEnv("LOGLEVEL")
		if !isPresent {
			//default to error
			level = "ERROR"
		}
	}
	setLogFilter(level, os.Stderr)
}

func setLogFilter(minLevel string, writer io.Writer) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "NOTICE", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(minLevel),
		Writer:   writer,
	}
	log.SetOutput(filter)
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
