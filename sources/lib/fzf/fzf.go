

package fzf


import fzf "github.com/junegunn/fzf/src"

import . "github.com/volution/z-run/lib/common"
import . "github.com/volution/z-run/embedded"




func FzfMain (_embedded bool, _fullscreen bool, _arguments []string, _environment map[string]string) (*Error) {
	
	
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
					return Errorf (0xbf8543cb, "invalid arguments")
			}
		}
	}
	
	
	if !_dryrun {
		if IsStdinTerminal () {
			return Errorf (0x34efe59c, "stdin is a TTY")
		}
		if IsStdoutTerminal () {
			return Errorf (0xf12b8d81, "stdout is a TTY")
		}
		if ! IsStderrTerminal () {
			return Errorf (0x55a1298a, "stderr is not a TTY")
		}
	}
	
	
	switch _fullscreenEnv, _ := _environment["ZFZF_FULLSCREEN"]; _fullscreenEnv {
		case "true" :
			_fullscreen = true
		case "false" :
			_fullscreen = false
		default :
			// NOP
	}
	
	
	fzf.MinimalMain (_arguments, _fullscreen, BUILD_VERSION, "z-fzf")
	
	panic (AbortUnreachable (0x4716a580))
}

