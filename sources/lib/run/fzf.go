

package zrun


import "os"
import "unsafe"

import fzf "github.com/junegunn/fzf/src"
import fzf_tui "github.com/junegunn/fzf/src/tui"
import isatty "github.com/mattn/go-isatty"

import . "github.com/cipriancraciun/z-run/lib/common"




func FzfMain (_embedded bool, _arguments []string, _environment map[string]string) (*Error) {
	
	
	if _embedded {
		if len (_arguments) != 0 {
			return Errorf (0x68f8e127, "invalid arguments")
		}
	}
	
	
	_dryrun := false
	if !_embedded {
		for _, _argument := range _arguments {
			switch _argument {
				case "-h", "-help", "--help" :
					_dryrun = true
					_arguments = []string {"--help"}
					break
				case "-v", "-version", "--version" :
					_dryrun = true
					_arguments = []string {"--version"}
					break
			}
		}
	}
	
	
	if !_dryrun {
		if isatty.IsTerminal (os.Stdin.Fd ()) {
			return Errorf (0x34efe59c, "stdin is a TTY")
		}
		if isatty.IsTerminal (os.Stdout.Fd ()) {
			return Errorf (0xf12b8d81, "stdout is a TTY")
		}
		if ! isatty.IsTerminal (os.Stderr.Fd ()) {
			return Errorf (0x55a1298a, "stderr is not a TTY")
		}
	}
	
	
	os.Args = append ([]string {"z-run"}, _arguments ...)
	
	
	fzf.Init ()
	
	var _options *fzf.Options
	
	if _embedded {
		
		_options = fzf.DefaultOptions ()
		
		_options.Prompt = ": "
		
		_options.Fuzzy = false
		_options.Extended = true
		_options.Case = fzf.CaseIgnore
		_options.Normalize = true
		_options.Sort = 1
		_options.Multi = 0
		
		_options.Theme = fzf_tui.NoColorTheme ()
		_options.Bold = false
		_options.ClearOnExit = true
		_options.Mouse = false
		
		// NOTE:  Replace `accept` with `accept-non-empty` action!
		for _key, _actions := range _options.Keymap {
			if (_key.Type == fzf_tui.CtrlM) || (_key.Type == fzf_tui.DoubleClick) {
				_action := &_actions[0]
				* ((*int) (unsafe.Pointer (_action))) = 7
			}
		}
		
	} else {
		
		_options = fzf.ParseOptions ()
		
	}
	
	fzf.Run (_options, "z-run", "")
	
	panic (0x4716a580)
}

