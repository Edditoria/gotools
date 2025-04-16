// godate for Good-Oh-Date! A friendly command-line tool for date.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Edditoria/gotools/flags"
	"github.com/Edditoria/gotools/lists"
)

// Preset format of date and time.
type presetFormat struct {
	Layout  string // in Go time format.
	IsLocal bool   // to show local time or UTC time.
	Usage   string // One-line usage.
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
	preDateTimeFriendly = presetFormat{"2006-01-02 15:04:05 Mon UTC-07", true, "friendly local date-time"}
	preDateTimeUTC      = presetFormat{time.RFC3339, false, "RFC3339-compatible UTC date-time"}
	preDateTimeStrict   = presetFormat{time.RFC3339, true, "RFC3339-compatible local date-time"}
	preDateTimeSerial   = presetFormat{"2006-01-02-150405", true, "serial-like local date-time"}
	preDateOnly         = presetFormat{"2006-01-02", true, "local date without time"}
	preTimeOnly         = presetFormat{"15:04:05", true, "local time without date"}
)

var (
	flagPreset     string
	flagPresetEnum *flags.StringFlagEnum
	flagDateOnly   bool
	flagTimeOnly   bool
)

func init() {
	flag.BoolVar(&flagDateOnly, "d", false, "local date only")
	flag.BoolVar(&flagTimeOnly, "t", false, "local time only")
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
	case "utc":
		exitCode := handleQuickCommand("utc", &preDateTimeUTC)
		os.Exit(exitCode)
	case "strict":
		exitCode := handleQuickCommand("strict", &preDateTimeStrict)
		os.Exit(exitCode)
	case "serial":
		exitCode := handleQuickCommand("serial", &preDateTimeSerial)
		os.Exit(exitCode)
	default:
		// When cli: godate [-flag...] or godate wrongArg(s)
		flag.Parse()
		presetOpt, presetErr := handlePresetFlag()
		if presetErr != nil {
			os.Stderr.WriteString(presetErr.Error() + "\n")
			flag.PrintDefaults()
			os.Exit(2)
		}
		if flagDateOnly && flagTimeOnly {
			fmt.Fprintf(os.Stdout, "%s\n", presetOpt)
			os.Exit(0)
		}
		if flagDateOnly {
			fmt.Fprintf(os.Stdout, "%s\n", preDateOnly)
			os.Exit(0)
		}
		if flagTimeOnly {
			fmt.Fprintf(os.Stdout, "%s\n", preTimeOnly)
			os.Exit(0)
		}
		fmt.Fprintf(os.Stdout, "%s\n", presetOpt)
		os.Exit(0)
	}
}

// It will print messages in console.
//   - @return Exit code.
func handleQuickCommand(name string, preFmt *presetFormat) int {
	if len(os.Args) == 3 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		fmt.Fprintf(os.Stdout, "Print %s\n", preFmt.Usage)
		return 0
	}
	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "godate: %s does not take argument\nSee: godate %s -h\n", name, name)
		return 2
	}
	fmt.Fprintf(os.Stdout, "%s\n", preFmt)
	return 0
}

func handlePresetFlag() (*presetFormat, error) {
	_, pPassedByUser := flags.IsFlagPassed("p")
	if !pPassedByUser {
		return &preDateTimeFriendly, nil
	}
	if pPassedByUser && flagPreset == "" {
		return nil, errors.New("flag needs an argument: -p=")
	}
	preOptI, err := flagPresetEnum.Record(flagPreset)
	if errors.Is(err, lists.ErrKeyNotFound) {
		return nil, errors.New("bad argument: -p=" + flagPreset)
	} else if err != nil {
		return nil, fmt.Errorf("unexpected program error: %w", err)
	}
	preOpt, ok := preOptI.(presetFormat)
	if !ok {
		return nil, errors.New("unexpected program error: unexpected type assertion: preOpt.(*presetFormat) not ok")
	}
	return &preOpt, nil
}
