package debughelp

import "fmt"

// Peek prints values with a limit to 200 counts
func Peek(count int, a ...interface{}) {
	if count < 200 {
		fmt.Println(a...)
	}
}
