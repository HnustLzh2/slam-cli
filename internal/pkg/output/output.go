package output

import (
	"fmt"
	"io"
)

func Println(w io.Writer, msg string) {
	fmt.Fprintln(w, msg)
}
