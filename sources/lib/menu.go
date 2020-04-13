

package zrun


import "bufio"
import "io"
import "os"
import "os/exec"
import "strings"


import isatty "github.com/mattn/go-isatty"




func menuMain (_arguments []string, _environment map[string]string) (*Error) {
	
	if len (_arguments) != 1 {
		return errorf (0x6b439ede, "invalid arguments")
	}
	
	_inputs := make ([]string, 0, 1024)
	if _stream, _error := os.Open (_arguments[0]); _error == nil {
		defer _stream.Close ()
		_reader := bufio.NewReader (_stream)
		for {
			if _line, _error := _reader.ReadString ('\n'); _error == nil {
				_input := strings.TrimRight (_line, "\n")
				_inputs = append (_inputs, _input)
			} else if _error == io.EOF {
				if _line == "" {
					break
				} else {
					errorf (0x1f57b1db, "expected proper line")
				}
			} else {
				return errorw (0x3dd692c8, _error)
			}
		}
	}
	
	// FIXME:  Should refactor this!
	_executable, _ := os.Executable ()
	_context := & Context {
			selfExecutable : _executable,
			// FIXME:  Handle this!
			cleanEnvironment : nil,
			terminal : os.Getenv ("TERM"),
		}
	
	if _outputs, _error := menuSelect (_inputs, _context); _error == nil {
		for _, _output := range _outputs {
			if _, _error := io.WriteString (os.Stdout, _output + "\n"); _error != nil {
				return errorw (0xeb4af0b7, _error)
			}
		}
		os.Exit (0)
		panic (0x1ad77faa)
	} else {
		return _error
	}
}




func menuSelect (_inputs []string, _context *Context) ([]string, *Error) {
	
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


func menuSelect_0 (_inputsChannel <-chan string, _outputsChannel chan<- string, _context *Context) (*Error) {
	
	_hasTerminal := (_context.terminal != "") && (_context.terminal != "dumb")
	
	if _hasTerminal {
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
			Dir : _context.workspace,
		}
	
	_commandFzf := false
	if _hasTerminal {
		_commandFzf = true
		_command.Path = _context.selfExecutable
		_command.Args = []string {
				"[z-run:select]",
			}
		_command.Env = []string {
				"TERM=" + _context.terminal,
			}
	} else if _path, _error := exec.LookPath ("rofi"); _error == nil {
		_command.Path = _path
		_command.Args = []string {
				"[z-run:select]",
				"-dmenu",
				"-p", "z-run",
				"-l", "16",
				"-i",
				"-no-custom",
			}
	} else if _path, _error := exec.LookPath ("dmenu"); _error == nil {
		_command.Path = _path
		_command.Args = []string {
				"[z-run:select]",
				"-p", "z-run",
				"-l", "16",
				"-i",
			}
	} else {
		return errorf (0xb91714f7, "unresolved `x-input`")
	}
	
	if _command.Env == nil {
		if _context.cleanEnvironment != nil {
			_command.Env = processEnvironment (_context, nil)
		} else {
			// FIXME:  Handle this!
			// _command.Env = []string {}
		}
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
			switch _exitCode {
				case 0 :
					if _outputsCount == 0 {
						return errorf (0x4e0abce6, "invalid outputs")
					}
				case 1 :
					if _outputsCount != 0 {
						return errorf (0x6ad0fdcd, "invalid outputs")
					}
				default :
					return errorf (0xb156b11d, "failed")
			}
		}
	} else {
		return _error
	}
	
	return nil
}

