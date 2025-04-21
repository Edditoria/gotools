/*
Add more ~~absraction~~ to Golang's [flag] package. Hehe...
*/
package flags

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Edditoria/gotools/lists"
)

// Create a new, empty [StringFlagEnum].
func NewStringFlagEnum(name string) *StringFlagEnum {
	enum := &StringFlagEnum{
		Name: name,
	}
	enum.OrderedMap = lists.NewOrderedMap()
	enum.UsageLine = enum.DefaultUsageLine
	return enum
}

// An enum-like struct to limit options for `yourcommand -stringflag=[option|option...]`.
// Options are stored in order.
// This struct leverage [lists.OrderedMap] and promotes its methods.
//
// You may want to start with [NewStringFlagEnum] function to initiate ordered map in this struct.
type StringFlagEnum struct {
	// Promotes [lists.OrderedMap] with all methods.
	*lists.OrderedMap
	// Expect original flag name in the string flag, e.g. "dev" in "-dev=edditoria".
	Name string
	// Set your own one-line usage. Default will be set when do [NewStringFlagEnum]
	UsageLine func() string
}

// One-line usage for the string flag.
func (enum StringFlagEnum) DefaultUsageLine() string {
	if enum.UsageLine == nil {
		flagValues := enum.OrderedMap.Keys()
		return fmt.Sprintf("-%s=[%s]", enum.Name, strings.Join(flagValues, "|"))
	} else {
		return enum.UsageLine()
	}
}

// Write better help message for your command.
// This function is the heart of this package.
//
// You can do something like:
//
//	flag.Usage = func() { flags.PrintHelp(args...) }
//
// Expect it prints to `stderr` or `stdout` like [flag.PrintDefaults].
func PrintHelp(name string, desc string, subcmdGroupList []*SubcmdGroup) {
	fmt.Fprintf(flag.CommandLine.Output(), "%s: ", name)
	fmt.Fprintf(flag.CommandLine.Output(), "%s\n\n", desc)
	fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\n")
	var sb strings.Builder
	for _, subcmdGroup := range subcmdGroupList {
		subcmdGroup.Help(&sb)
	}
	fmt.Fprintf(flag.CommandLine.Output(), "%s", sb.String())
}

// Check if a [flag.Flag] is passed by cli user.
//   - @return first bool: Did the flag parse by the program developer?
//   - @return second bool: Did user pass the flag argument? (Always false if first bool false)
//
// NOTE: It returns (true, true) if user do `yourcommand -u=` and is properly parsed.
func IsFlagPassed(flagName string) (bool, bool) {
	if flag.Lookup(flagName) == nil {
		return false, false
	}
	wasSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flagName {
			wasSet = true
		}
	})
	return true, wasSet
}
