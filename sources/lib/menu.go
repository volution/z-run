

package zrun


import "os"
import "os/exec"


import isatty "github.com/mattn/go-isatty"




func menuSelect (_inputs []string, _context *Context) ([]string, error) {
	
	_inputsChannel := make (chan string, 1024)
	_outputsChannel := make (chan string, 1024)
	_outputs := make ([]string, 0, 1024)
	
	go func () () {
		for _, _input := range _inputs {
			_inputsChannel <- _input
		}
		close (_inputsChannel)
	} ()
	
	go func () () {
		for {
			_output, _ok := <- _outputsChannel
			if _ok {
				_outputs = append (_outputs, _output)
			} else {
				break
			}
		}
		close (_outputsChannel)
	} ()
	
	if _error := menuSelect_0 (_inputsChannel, _outputsChannel, _context); _error == nil {
		return _outputs, nil
	} else {
		return nil, _error
	}
}


func menuSelect_0 (_inputsChannel <-chan string, _outputsChannel chan<- string, _context *Context) (error) {
	
	if _context.terminal != "" {
		if ! isatty.IsTerminal (os.Stderr.Fd ()) {
			return errorf (0xfc026596, "stderr is not a TTY")
		}
	} else {
//		return errorf (0xbdbc268d, "expected `TERM`")
	}
	
	_command := & exec.Cmd {
			Stdin : nil,
			Stdout : nil,
			Stderr : os.Stderr,
			Dir : "",
		}
	
	_commandFzf := false
	if _context.terminal != "" {
		_commandFzf = true
		_command.Path = _context.selfExecutable
		_command.Args = []string {
				"[z-run:select]",
			}
		_command.Env = []string {
				"TERM=" + _context.terminal,
			}
	} else if _path, _error := exec.LookPath ("x-input"); _error == nil {
		_command.Path = _path
		_command.Args = []string {
				"[z-run:input]",
				"select",
				"run:",
			}
		_command.Env = processEnvironment (_context, nil)
	} else {
		return errorf (0xb91714f7, "expected `x-input`")
	}
	
	if _exitCode, _, _outputsCount, _error := processExecuteAndPipe (_command, _inputsChannel, _outputsChannel); _error == nil {
		if _commandFzf {
			switch _exitCode {
				case 0 :
					if _outputsCount == 0 {
						return errorf (0xbb7ff442, "invalid outputs")
					}
				case 1 :
					if _outputsCount != 0 {
						return errorf (0x6bd364da, "invalid outputs")
					}
				case 130 :
					if _outputsCount != 0 {
						return errorf (0xac4b1681, "invalid outputs")
					}
				case 2 :
					return errorf (0x85cabb2a, "failed")
				default :
					return errorf (0xef9908df, "failed")
			}
		} else {
			if _exitCode != 0 {
				return errorf (0xb156b11d, "failed")
			}
		}
	} else {
		return _error
	}
	
	return nil
}

