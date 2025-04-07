// godate for Good-Oh-Date! A friendly command-line tool for date.
package main

import (
	"fmt"
	"time"
)

// Collection of preset format.
const (
	preDateDefault = "2006-01-02"
	preTimeDefault = "15:04:05"
	preDateTimeDefault = preDateDefault + " " + preTimeDefault
)

func main() {
	now := time.Now().Format(preDateTimeDefault)
	fmt.Println(now)
}
