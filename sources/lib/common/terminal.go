

package common


import "os"

import isatty "github.com/mattn/go-isatty"




func CheckStdioTerminal () (*Error) {
	if ! IsStdinTerminal () {
		return Errorf (0x05d60b72, "stdin is not a TTY")
	}
	if ! IsStdoutTerminal () {
		return Errorf (0xc432630a, "stdout is not a TTY")
	}
	if ! IsStderrTerminal () {
		return Errorf (0x77924518, "stderr is not a TTY")
	}
	return nil
}




func IsStdinTerminal () (bool) {
	return isatty.IsTerminal (os.Stdin.Fd ())
}

func IsStdoutTerminal () (bool) {
	return isatty.IsTerminal (os.Stdout.Fd ())
}

func IsStderrTerminal () (bool) {
	return isatty.IsTerminal (os.Stderr.Fd ())
}

func IsFdTerminal (_descriptor uintptr) (bool) {
	return isatty.IsTerminal (_descriptor)
}

