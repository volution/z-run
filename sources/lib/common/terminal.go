

package common


import "os"

import "golang.org/x/term"




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
	return IsFileTerminal (os.Stdin)
}

func IsStdoutTerminal () (bool) {
	return IsFileTerminal (os.Stdout)
}

func IsStderrTerminal () (bool) {
	return IsFileTerminal (os.Stderr)
}

func IsFileTerminal (_file *os.File) (bool) {
	return IsFdTerminal (_file.Fd ())
}

func IsFdTerminal (_descriptor uintptr) (bool) {
	return term.IsTerminal (int (_descriptor))
}

