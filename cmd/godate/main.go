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

// CLI descriptions.
const (
	cliName     = "godate"
	cliDesc     = "godate for Good-Oh Date! A friendly command-line tool for date and time."
	cliFootnote = "Source code:\n  https://github.com/Edditoria/gotools"
)

// Preset format of date and time.
type presetFormat struct {
	Name    string
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

// Skipped details field.
func (p presetFormat) CreateSubcmd() *flags.Subcmd {
	subcmd := flags.NewSubcmd(p.Name, "Print "+p.Usage, "", flag.ExitOnError)
	return subcmd
}

// Collection of preset format.
var (
	preDateTimeFriendly = presetFormat{
		Name:    "friendly",
		Layout:  "2006-01-02 15:04:05 Mon UTC-07",
		IsLocal: true,
		Usage:   "friendly local date-time"}
	preDateTimeUTC = presetFormat{
		Name:    "utc",
		Layout:  time.RFC3339,
		IsLocal: false,
		Usage:   "RFC3339-compatible UTC date-time",
	}
	preDateTimeStrict = presetFormat{
		Name:    "strict",
		Layout:  time.RFC3339,
		IsLocal: true,
		Usage:   "RFC3339-compatible local date-time",
	}
	preDateTimeSerial = presetFormat{
		Name:    "serial",
		Layout:  "2006-01-02-150405",
		IsLocal: true,
		Usage:   "serial-like local date-time",
	}
	preDateOnly = presetFormat{
		Name:    "date-only",
		Layout:  "2006-01-02",
		IsLocal: true,
		Usage:   "local date without time",
	}
	preTimeOnly = presetFormat{
		Name:    "time-only",
		Layout:  "15:04:05",
		IsLocal: true,
		Usage:   "local time without date",
	}
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

	// Prepare "Quick actions":
	subcmdUTC := preDateTimeUTC.CreateSubcmd()
	subcmdStrict := preDateTimeStrict.CreateSubcmd()
	subcmdSerial := preDateTimeSerial.CreateSubcmd()

	// When cli: godate (without any arg)
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "%s\n", preDateTimeFriendly)
		os.Exit(0)
	}

	// When cli: godate [action] [-flag...]
	switch os.Args[1] {
	case subcmdUTC.Name():
		exitCode := handleQuickActions(subcmdUTC, &preDateTimeUTC)
		os.Exit(exitCode)
	case subcmdStrict.Name():
		exitCode := handleQuickActions(subcmdStrict, &preDateTimeStrict)
		os.Exit(exitCode)
	case subcmdSerial.Name():
		exitCode := handleQuickActions(subcmdSerial, &preDateTimeSerial)
		os.Exit(exitCode)
	default:
		// When cli: godate [-flag...] or godate wrongArg(s)
		quickActions := &flags.SubcmdGroup{
			Name:       "Quick actions",
			Summary:    "small commands for daily life",
			SubcmdList: []*flags.Subcmd{subcmdUTC, subcmdStrict, subcmdSerial},
		}
		subcmdGroups := []*flags.SubcmdGroup{quickActions, quickActions}
		flag.Usage = func() { flags.PrintHelp(cliName, cliDesc, subcmdGroups, cliFootnote) }
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
func handleQuickActions(subcmd *flags.Subcmd, preFmt *presetFormat) int {
	if len(os.Args) == 3 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		fmt.Fprintf(os.Stdout, "%s\n", subcmd.Summary)
		return 0
	}
	if len(os.Args) > 2 {
		name := subcmd.Name()
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
