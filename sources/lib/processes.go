

package zrun


import "bufio"
import "bytes"
import "io"
import "os"
import "os/exec"
import "sort"
import "strings"
import "sync"




func processEnvironment_1 (_context *Context, _overrides ... map[string]string) ([]string) {
	
	_extraEnvironment := make (map[string]string, 16)
	
	_extraEnvironment["ZRUN_EXECUTABLE"] = _context.selfExecutable
	_extraEnvironment["ZRUN_WORKSPACE"] = _context.workspace
	_extraEnvironment["ZRUN_CACHE"] = _context.cacheRoot
	
	if _context.executablePaths != nil {
		_paths := strings.Join (_context.executablePaths, string (os.PathListSeparator))
		_paths = strings.Trim (_paths, string (os.PathListSeparator))
		if _paths == "" {
			_paths = "/dev/null"
		}
		_extraEnvironment["PATH"] = _paths
	} else {
		_extraEnvironment["PATH"] = "/dev/null"
	}
	
	if _context.terminal != "" {
		_extraEnvironment["TERM"] = _context.terminal
	} else {
		_extraEnvironment["TERM"] = "dumb"
	}
	
	_overrides_0 := make ([]map[string]string, 0, 1 + len (_overrides))
	_overrides_0 = append (_overrides_0, _extraEnvironment)
	_overrides_0 = append (_overrides_0, _overrides ...)
	return processEnvironment_0 (_context.cleanEnvironment, _overrides_0 ...)
}


func processEnvironment_0 (_environment map[string]string, _overrides ... map[string]string) ([]string) {
	
	_environmentMap := make (map[string]string, len (_environment))
	
	for _name, _value := range _environment {
		_environmentMap[_name] = _value
	}
	for _, _overrides := range _overrides {
		for _name, _value := range _overrides {
			if _value != "" {
				_environmentMap[_name] = _value
			} else {
				delete (_environmentMap, _name)
			}
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




func processExecuteAndPipe (_command *exec.Cmd, _inputsChannel <-chan string, _outputsChannel chan<- string, _ignoreMissingNewline bool) (int, uint, uint, *Error) {
	
	var _stdin io.WriteCloser
	if _inputsChannel != nil {
		if _stream, _error := _command.StdinPipe (); _error == nil {
			// NOTE:  Due to race conditions within the goroutine, we leave this to be closed by the garbage collector.
			// defer _stream.Close ()
			_stdin = _stream
		} else {
			return -1, 0, 0, errorw (0xb3c4228d, _error)
		}
	}
	
	var _stdout io.ReadCloser
	if _outputsChannel != nil {
		if _stream, _error := _command.StdoutPipe (); _error == nil {
			// NOTE:  Due to race conditions within the goroutine, we leave this to be closed by the garbage collector.
			// defer _stream.Close ()
			_stdout = _stream
		} else {
			return -1, 0, 0, errorw (0x1067469f, _error)
		}
	}
	
	if _error := _command.Start (); _error != nil {
		return -1, 0, 0, errorw (0x26e1988c, _error)
	}
	
	_waiter := & sync.WaitGroup {}
	
	var _stdinError *Error
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
						_stdinError = errorw (0xb5ca9a1c, _error)
						break
					}
					_inputsCount += 1
				} else {
					break
				}
			}
			if _error := _stdin.Close (); _error != nil {
				_stdinError = errorw (0x7e9a4f14, _error)
			}
//			logf ('d', 0xc6eca1ca, "ending stdin loop")
			_waiter.Done ()
		} ()
	}
	
	var _stdoutError *Error
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
					if _line == "" {
						break
					} else {
						if _ignoreMissingNewline {
							_output := _line
//							logf ('d', 0x369ccac9, "read from stdout (without newline): `%s`", _output)
							_outputsChannel <- _output
							_outputsCount += 1
							break
						} else {
							_stdoutError = errorf (0x1bc14ac4, "expected proper line")
							break
						}
					}
				} else {
					_stdoutError = errorw (0xb783c8c4, _error)
					break
				}
			}
			if _error := _stdout.Close (); _error != nil {
				_stdoutError = errorw (0xf185ae2a, _error)
			}
//			logf ('d', 0x90515c65, "ending stdout loop")
			_waiter.Done ()
		} ()
	}
	
	_waiter.Wait ()
	
	var _waitError *Error
//	logf ('d', 0x7ce5281a, "starting wait")
	if _error := _command.Wait (); _error != nil {
		_waitError = errorw (0x6f9dfa7d, _error)
	}
//	logf ('d', 0xa36df40d, "ending wait")
	
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




func processExecuteGetStdout (_command *exec.Cmd) (int, []byte, *Error) {
	
	_stdout := bytes.NewBuffer (nil)
	_stdout.Grow (128 * 1024)
	
	if _command.Stdout == nil {
		_command.Stdout = _stdout
	} else {
		return -1, nil, errorf (0x7cd15552, "invalid state")
	}
	
	_waitError := _command.Run ()
	
	for _, _descriptor := range _command.ExtraFiles {
		_descriptor.Close ()
	}
	
	if _waitError != nil {
		if _command.ProcessState.Exited () {
			if _exitCode := _command.ProcessState.ExitCode (); _exitCode >= 0 {
				return _exitCode, _stdout.Bytes (), nil
			} else {
				return -1, _stdout.Bytes (), errorw (0xc8553b48, _waitError)
			}
		} else {
			return -1, nil, errorw (0x4b785e1d, _waitError)
		}
	} else {
		return 0, _stdout.Bytes (), nil
	}
}

