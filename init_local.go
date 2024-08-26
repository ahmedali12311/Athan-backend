//go:build local
// +build local

package main

import (
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/nwidger/jsoncolor"
	"github.com/rs/zerolog"
)

func init() {
	log_writer = zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	zerolog.InterfaceMarshalFunc = func(v any) ([]byte, error) {
		f := jsoncolor.NewFormatter()

		// set custom colors
		f.SpaceColor = color.New(color.FgRed, color.Bold)
		f.CommaColor = color.New(color.FgWhite, color.Bold)
		f.ColonColor = color.New(color.FgYellow, color.Bold)
		f.ObjectColor = color.New(color.FgGreen, color.Bold)
		f.ArrayColor = color.New(color.FgHiRed)
		f.FieldColor = color.New(color.FgCyan)
		f.StringColor = color.New(color.FgHiYellow)
		f.TrueColor = color.New(color.FgCyan, color.Bold)
		f.FalseColor = color.New(color.FgHiRed)
		f.NumberColor = color.New(color.FgHiMagenta)
		f.NullColor = color.New(color.FgWhite, color.Bold)
		f.StringQuoteColor = color.New(color.FgBlue, color.Bold)

		return jsoncolor.MarshalIndentWithFormatter(v, "", "  ", f)
	}
}
