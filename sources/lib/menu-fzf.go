

package zrun


import "os"


import fzf "github.com/junegunn/fzf/src"
import fzf_tui "github.com/junegunn/fzf/src/tui"
import isatty "github.com/mattn/go-isatty"




func fzfSelectMain () (error) {
	
	if len (os.Args) == 1 {
		// NOP
	} else if len (os.Args) == 2 {
		if _stream, _error := os.Open (os.Args[1]); _error == nil {
			os.Stdin.Close ()
			os.Stdin = _stream
		}
	} else {
		return errorf (0x68f8e127, "invalid arguments")
	}
	
	if isatty.IsTerminal (os.Stdin.Fd ()) {
		return errorf (0x34efe59c, "stdin is a TTY")
	}
//	if isatty.IsTerminal (os.Stdout.Fd ()) {
//		return errorf (0xf12b8d81, "stdout is a TTY")
//	}
	if ! isatty.IsTerminal (os.Stderr.Fd ()) {
		return errorf (0x55a1298a, "stderr is not a TTY")
	}
	
	_options := fzf.DefaultOptions ()
	
	_options.Fuzzy = false
	_options.Extended = true
	_options.Case = fzf.CaseIgnore
	_options.Normalize = true
	_options.Sort = 1
	_options.Multi = false
	
	_options.Theme = fzf_tui.Default16
	_options.Theme = nil
	_options.Bold = false
	_options.ClearOnExit = true
	_options.Mouse = false
	
	fzf.Run (_options, "z-run")
	panic (0x4716a580)
}

