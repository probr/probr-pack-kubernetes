package config

import (
	"io"
	"log"
	"strings"

	"github.com/hashicorp/logutils"
)

// Override the minimum log level.
func SetLogFilter(minLevel string, writer io.Writer) {
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
