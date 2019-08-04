

package lib


import "bufio"
import "bytes"
import "io"
import "os/exec"
import "sort"
import "strings"
import "sync"




func commandEnvironment (_context *Context, _overrides map[string]string) ([]string) {
	
	_environmentMap := make (map[string]string, len (_context.cleanEnvironment) + len (_overrides))
	
	for _name, _value := range _context.cleanEnvironment {
		_environmentMap[_name] = _value
	}
	for _name, _value := range _overrides {
		if _value != "" {
			_environmentMap[_name] = _value
		} else {
			delete (_environmentMap, _name)
		}
	}
	
	var _environmentArray []string = make ([]string, 0, len (_environmentMap))
	for _name, _value := range _environmentMap {
		_variable := _name + "=" + _value
		_environmentArray = append (_environmentArray, _variable)
	}
	
	sort.Strings (_environmentArray)
	
	return _environmentArray
}




func commandExecuteAndPipe (_command *exec.Cmd, _inputsChannel <-chan string, _outputsChannel chan<- string) (int, uint, uint, error) {
	
	var _stdin io.WriteCloser
	if _inputsChannel != nil {
		if _stream, _error := _command.StdinPipe (); _error == nil {
			// NOTE:  Due to race conditions within the goroutine, we leave this to be closed by the garbage collector.
			// defer _stream.Close ()
			_stdin = _stream
		} else {
			return -1, 0, 0, _error
		}
	}
	
	var _stdout io.ReadCloser
	if _outputsChannel != nil {
		if _stream, _error := _command.StdoutPipe (); _error == nil {
			// NOTE:  Due to race conditions within the goroutine, we leave this to be closed by the garbage collector.
			// defer _stream.Close ()
			_stdout = _stream
		} else {
			return -1, 0, 0, _error
		}
	}
	
	if _error := _command.Start (); _error != nil {
		return -1, 0, 0, _error
	}
	
	_waiter := & sync.WaitGroup {}
	
	var _stdinError error
	var _inputsCount uint
	if _inputsChannel != nil {
		_waiter.Add (1)
		go func () () {
//			logf ('d', 0x41785333, "starting stdin loop")
			_buffer := bytes.NewBuffer (nil)
			for {
				_input, _ok := <- _inputsChannel
//				logf ('d', 0xf997ad63, "writing to stdin: `%s`", _input)
				if _ok {
					_buffer.Reset ()
					_buffer.WriteString (_input)
					_buffer.WriteByte ('\n')
					if _, _error := _buffer.WriteTo (_stdin); _error != nil {
						_stdinError = _error
						break
					}
					_inputsCount += 1
				} else {
					break
				}
			}
			if _error := _stdin.Close (); _error != nil {
				_stdinError = _error
			}
//			logf ('d', 0xc6eca1ca, "ending stdin loop")
			_waiter.Done ()
		} ()
	}
	
	var _stdoutError error
	var _outputsCount uint
	if _outputsChannel != nil {
		_waiter.Add (1)
		go func () () {
//			logf ('d', 0x61503d28, "starting stdout loop")
			_buffer := bufio.NewReader (_stdout)
			for {
				if _line, _error := _buffer.ReadString ('\n'); _error == nil {
					_output := strings.TrimRight (_line, "\n")
//					logf ('d', 0xa6f11fbf, "read from stdout: `%s`", _output)
					_outputsChannel <- _output
					_outputsCount += 1
				} else if _error == io.EOF {
					if _line != "" {
						_stdoutError = errorf (0x1bc14ac4, "expected proper line")
					}
					break
				} else {
					_stdoutError = _error
					break
				}
			}
			if _error := _stdout.Close (); _error != nil {
				_stdoutError = _error
			}
//			logf ('d', 0x90515c65, "ending stdout loop")
			_waiter.Done ()
		} ()
	}
	
	var _waitError error
//	logf ('d', 0x7ce5281a, "starting wait")
	if _error := _command.Wait (); _error != nil {
		_waitError = _error
	}
//	logf ('d', 0xa36df40d, "ending wait")
	
	_waiter.Wait ()
	
	if _stdinError != nil {
		return -1, 0, 0, _stdinError
	}
	if _stdoutError != nil {
		return -1, 0, 0, _stdoutError
	}
	
	if _waitError != nil {
		if _command.ProcessState.Exited () {
			if _exitCode := _command.ProcessState.ExitCode (); _exitCode >= 0 {
				return _exitCode, _inputsCount, _outputsCount, nil
			} else {
				return -1, _inputsCount, _outputsCount, _waitError
			}
		} else {
			return -1, 0, 0, _waitError
		}
	} else {
		return 0, _inputsCount, _outputsCount, nil
	}
}




func commandExecuteGetStdout (_command *exec.Cmd) (int, []byte, error) {
	
	_stdout := bytes.NewBuffer (nil)
	_stdout.Grow (128 * 1024)
	
	if _command.Stdout == nil {
		_command.Stdout = _stdout
	} else {
		return -1, nil, errorf (0x7cd15552, "invalid state")
	}
	
	_waitError := _command.Run ()
	
	if _waitError != nil {
		if _command.ProcessState.Exited () {
			if _exitCode := _command.ProcessState.ExitCode (); _exitCode >= 0 {
				return _exitCode, _stdout.Bytes (), nil
			} else {
				return -1, _stdout.Bytes (), _waitError
			}
		} else {
			return -1, nil, _waitError
		}
	} else {
		return 0, _stdout.Bytes (), nil
	}
}

