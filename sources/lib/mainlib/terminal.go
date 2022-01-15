

package mainlib


import "os"

import isatty "github.com/mattn/go-isatty"

import . "github.com/cipriancraciun/z-run/lib/common"




func CheckMainTerminal () (*Error) {
	if ! isatty.IsTerminal (os.Stdin.Fd ()) {
		return Errorf (0x05d60b72, "stdin is not a TTY")
	}
	if ! isatty.IsTerminal (os.Stdout.Fd ()) {
		return Errorf (0xc432630a, "stdout is not a TTY")
	}
	if ! isatty.IsTerminal (os.Stderr.Fd ()) {
		return Errorf (0x77924518, "stderr is not a TTY")
	}
	return nil
}

