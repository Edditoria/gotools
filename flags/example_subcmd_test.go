package flags_test

import (
	"flag"
	"fmt"

	"github.com/Edditoria/gotools/flags"
)

func ExampleLongestSubcmdName() {
	// The name "long", 4 characters:
	subcmd1 := flags.NewSubcmd("long", "", "", flag.ExitOnError)
	// The name "loooooong", 9 characters:
	subcmd2 := flags.NewSubcmd("loooooong", "", "", flag.ExitOnError)
	longest := flags.LongestSubcmdName([]*flags.Subcmd{subcmd1, subcmd2})
	fmt.Printf("%d", longest)
	// Output: 9
}

func ExampleSubcmd_HelpLine() {
	add := flags.NewSubcmd("add", "Add file contents to the index", "", flag.ExitOnError)
	checkout := flags.NewSubcmd("checkout", "Switch branches or restore working tree files", "", flag.ExitOnError)
	minWidth := max(len(add.Name()), len(checkout.Name()))

	// Emulate printing subcommand session of help message in main command:
	fmt.Printf("1:|%s|\n", add.HelpLine(2, minWidth))
	fmt.Printf("2:|%s|\n", checkout.HelpLine(2, minWidth))
	// Output:
	// 1:|  add       Add file contents to the index|
	// 2:|  checkout  Switch branches or restore working tree files|
}
