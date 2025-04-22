package flags

import (
	"flag"
	"strings"
)

// Subcommand group for building help message.
type SubcmdGroup struct {
	Name       string // Group name.
	Summary    string // Short description for main help.
	SubcmdList []*Subcmd
}

// Lookup to find subcommand using its name, that is exactly [flag.FlagSet.Name].
//   - @return *Subcmd : If not found, it returns nil.
//   - @return bool    : Represents "found" or not.
func (group *SubcmdGroup) Lookup(subcmdName string) (*Subcmd, bool) {
	for _, subcmd := range group.SubcmdList {
		if subcmd.FlagSet.Name() == subcmdName {
			return subcmd, true
		}
	}
	return nil, false
}

// Generate help message of the [SubcmdGroup] session.
// It is a part of full [PrintHelp].
// The message will look like:
//
//	[Name]: [Heading to brief]
//
//	  [subcmd1]  [subcmd description]
//	  [subcmd2]  [subcmd description]
//	  (and so on...)
//
// It already takes care of indentation and space alignment for you.
//   - @return string : The output always ends with one empty line.
//   - @return error  : Always nil because of [strings.Builder], but better to check error for future-proof.
func (group *SubcmdGroup) Help() (string, error) {
	// The implementation of WriteString and WriteByte always returns the error nil.
	// So, it does not need to check error. Just return nil finally.
	// See: https://pkg.go.dev/strings#Builder.Write
	var sb strings.Builder
	sb.WriteString(group.Name)
	sb.WriteString(": ")
	sb.WriteString(group.Summary)
	sb.WriteString("\n\n")
	longest := LongestSubcmdName(group.SubcmdList)
	for _, subcmd := range group.SubcmdList {
		sb.WriteString(subcmd.HelpLine(2, longest))
		sb.WriteByte('\n')
	}
	sb.WriteByte('\n')
	return sb.String(), nil
}

// Subcommand to expand [flag.FlagSet] features.
// Recommend to start with [NewSubcmd] function.
type Subcmd struct {
	flag.FlagSet        // Embed [flag.FlagSet] with all methods.
	Summary      string // Short description for main help and subcommand help. Optional.
	Details      string // TODO: Long description for subcommand help. Optional.
}

// Create a new [Subcmd]. It will create and embed a [flag.FlagSet] for you.
// To skip printing summary and details, pass an empty string.
//   - @param name          : Use for [flag.NewFlagSet].
//   - @param summary       : Short description for main help and subcommand help, or empty string.
//   - @param details       : Long description for subcommand help, or empty string.
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
	sb.WriteString(name)
	if subcmd.Summary != "" {
		nameLen := len(name)
		var margin int
		if nameLen > minWidth {
			margin = 2
		} else {
			margin = minWidth - nameLen + 2
		}
		sb.WriteString(strings.Repeat(" ", margin))
		sb.WriteString(subcmd.Summary)
	}
	return sb.String()
}

// Who's the longest?
// If the list is empty, it returns 0.
func LongestSubcmdName(subcmdList []*Subcmd) int {
	longest := 0
	if len(subcmdList) == 0 {
		return 0
	}
	for _, subcmd := range subcmdList {
		thisSize := len(subcmd.FlagSet.Name())
		if thisSize > longest {
			longest = thisSize
		}
	}
	return longest
}
