package outputwriter

import (
	"fmt"
	"os"
)

func WriteStdoutf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

func WriteStderrf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
