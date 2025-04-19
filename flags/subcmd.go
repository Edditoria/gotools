package flags

import (
	"flag"
	"strings"
)

// Subcommand to expand [flag.FlagSet] features.
// Recommend to start with [NewSubcmd] function.
type Subcmd struct {
	flag.FlagSet        // Embed [flag.FlagSet] with all methods.
	Summary      string // Short description for main help and subcommand help.
	Details      string // Long description for subcommand help.
}

// Create a new [Subcmd]. It will create and embed a [flag.FlagSet] for you.
//   - @param name          : Use for [flag.NewFlagSet].
//   - @param summary       : Short description for main help and subcommand help.
//   - @param details       : Long description for subcommand help.
//   - @param errorHandling : Use for [flag.NewFlagSet].
//
// NOTE: You still need to parse the `FlagSet` in your program.
func NewSubcmd(name, summary, details string, errorHandling flag.ErrorHandling) *Subcmd {
	flagSet := flag.NewFlagSet(name, errorHandling)
	subcmd := &Subcmd{FlagSet: *flagSet, Summary: summary, Details: details}
	return subcmd
}

// Generate one-line help in format of:
// "subcmd  Short description for the subcommand".
// It does not end with a new-line, tho.
//   - @param indent   : Size of space indentation.
//   - @param minWidth : Minimum width to display subcommand name.
//
// Min-width is helpful to align a group of subcommands.
// It does not count the margin space (always 2-space) between name and summary.
// If min-width is not enough, this function will expand the space to `len(subcmd.Name())`.
func (subcmd *Subcmd) HelpLine(indent, minWidth int) string {
	var sb strings.Builder
	sb.WriteString(strings.Repeat(" ", indent))
	name := subcmd.FlagSet.Name()
	nameLen := len(name)
	var margin int
	if nameLen > minWidth {
		margin = 2
	} else {
		margin = minWidth - nameLen + 2
	}
	sb.WriteString(name)
	sb.WriteString(strings.Repeat(" ", margin))
	sb.WriteString(subcmd.Summary)
	return sb.String()
}
