package flags_test

import (
	"testing"

	"github.com/Edditoria/gotools/flags"
)

func Test_LongestSubcmdName(t *testing.T) {
	var want, got int

	zero := flags.NewSubcmd("", "", "", 0)
	four := flags.NewSubcmd("four", "", "", 0)
	six := flags.NewSubcmd("sixxxx", "", "", 0)
	eight := flags.NewSubcmd("eiiiight", "", "", 0)

	list := make([]*flags.Subcmd, 0)
	list = append(list, four, six)

	want = 6
	got = flags.LongestSubcmdName(list)
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
	list = append(list, eight)
	want = 8
	got = flags.LongestSubcmdName(list)
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
	want = 4
	got = flags.LongestSubcmdName([]*flags.Subcmd{zero, four})
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
	want = 0
	got = flags.LongestSubcmdName(make([]*flags.Subcmd, 0))
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
}
