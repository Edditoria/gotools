// godate for Good-Oh-Date! A friendly command-line tool for date.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Edditoria/gotools/flags"
)

// Preset format of date and time.
type presetFormat struct {
	Layout  string // in Go time format.
	IsLocal bool   // to show local time or UTC time.
}

// @implement [fmt.Stringer].
func (p presetFormat) String() string {
	if p.IsLocal {
		return time.Now().Format(p.Layout)
	}
	return time.Now().UTC().Format(p.Layout)
}

// Collection of preset format.
var (
	preDateTimeFriendly = presetFormat{"2006-01-02 15:04:05 UTC-07", true}
	preDateTimeSerial   = presetFormat{"2006-01-02-150405", true}
	preDateTimeUTC      = presetFormat{time.RFC3339, false}
)

var (
	flagPreset     string
	flagPresetEnum *flags.StringFlagEnum
)

func init() {
	flagPresetEnum = flags.NewStringFlagEnum("p")
	flagPresetEnum.Append("u", preDateTimeUTC)
	flagPresetEnum.Append("utc", preDateTimeUTC)
	flagPresetEnum.Append("s", preDateTimeSerial)
	flagPresetEnum.Append("serial", preDateTimeSerial)
	flagPresetEnum.Append("f", preDateTimeFriendly)
	flagPresetEnum.Append("friendly", preDateTimeFriendly)
	flag.StringVar(&flagPreset, "p", "", "p for preset: "+flagPresetEnum.UsageLine())
}

func main() {
	// When cli: godate (without any arg)
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "%s\n", preDateTimeFriendly)
		os.Exit(0)
	}

	// When cli: godate [action] [-flag...]
	switch os.Args[1] {
	default:
		// When cli: godate [-flag...]
		flag.Parse()
		preOpt, err := flagPresetEnum.Record(flagPreset)
		if err != nil {
			os.Stderr.WriteString("flag needs a correct argument: -p\n")
			flag.Usage()
			os.Exit(2)
		}
		fmt.Fprintf(os.Stdout, "%s\n", preOpt)
		os.Exit(0)
	}
}
