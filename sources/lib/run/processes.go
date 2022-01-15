

package zrun


import "bufio"
import "bytes"
import "io"
import "os/exec"
import "strings"
import "sync"

import . "github.com/cipriancraciun/z-run/lib/common"




func processExecuteAndPipe (_command *exec.Cmd, _inputsChannel <-chan string, _outputsChannel chan<- string, _ignoreMissingNewline bool) (int, uint, uint, *Error) {
	
	var _stdin io.WriteCloser
	if _inputsChannel != nil {
		if _stream, _error := _command.StdinPipe (); _error == nil {
			// NOTE:  Due to race conditions within the goroutine, we leave this to be closed by the garbage collector.
			// defer _stream.Close ()
			_stdin = _stream
		} else {
			return -1, 0, 0, Errorw (0xb3c4228d, _error)
		}
	}
	
	var _stdout io.ReadCloser
	if _outputsChannel != nil {
		if _stream, _error := _command.StdoutPipe (); _error == nil {
			// NOTE:  Due to race conditions within the goroutine, we leave this to be closed by the garbage collector.
			// defer _stream.Close ()
			_stdout = _stream
		} else {
			return -1, 0, 0, Errorw (0x1067469f, _error)
		}
	}
	
	if _error := _command.Start (); _error != nil {
		return -1, 0, 0, Errorf (0x26e1988c, "failed to spawn `%s`  //  %v", _command.Path, _error)
	}
	
	_waiter := & sync.WaitGroup {}
	
	var _stdinError *Error
	var _inputsCount uint
	if _inputsChannel != nil {
		_waiter.Add (1)
		go func () () {
//			Logf ('d', 0x41785333, "starting stdin loop")
			_buffer := bytes.NewBuffer (nil)
			for {
				_input, _ok := <- _inputsChannel
//				Logf ('d', 0xf997ad63, "writing to stdin: `%s`", _input)
				if _ok {
					_buffer.Reset ()
					_buffer.WriteString (_input)
					_buffer.WriteByte ('\n')
					if _, _error := _buffer.WriteTo (_stdin); _error != nil {
						_stdinError = Errorw (0xb5ca9a1c, _error)
						break
					}
					_inputsCount += 1
				} else {
					break
				}
			}
			if _error := _stdin.Close (); _error != nil {
				_stdinError = Errorw (0x7e9a4f14, _error)
			}
//			Logf ('d', 0xc6eca1ca, "ending stdin loop")
			_waiter.Done ()
		} ()
	}
	
	var _stdoutError *Error
	var _outputsCount uint
	if _outputsChannel != nil {
		_waiter.Add (1)
		go func () () {
//			Logf ('d', 0x61503d28, "starting stdout loop")
			_buffer := bufio.NewReader (_stdout)
			for {
				if _line, _error := _buffer.ReadString ('\n'); _error == nil {
					_output := strings.TrimRight (_line, "\n")
//					Logf ('d', 0xa6f11fbf, "read from stdout: `%s`", _output)
					_outputsChannel <- _output
					_outputsCount += 1
				} else if _error == io.EOF {
					if _line == "" {
						break
					} else {
						if _ignoreMissingNewline {
							_output := _line
//							Logf ('d', 0x369ccac9, "read from stdout (without newline): `%s`", _output)
							_outputsChannel <- _output
							_outputsCount += 1
							break
						} else {
							_stdoutError = Errorf (0x1bc14ac4, "expected proper line")
							break
						}
					}
				} else {
					_stdoutError = Errorw (0xb783c8c4, _error)
					break
				}
			}
			if _error := _stdout.Close (); _error != nil {
				_stdoutError = Errorw (0xf185ae2a, _error)
			}
//			Logf ('d', 0x90515c65, "ending stdout loop")
			_waiter.Done ()
		} ()
	}
	
	_waiter.Wait ()
	
	var _waitError *Error
//	Logf ('d', 0x7ce5281a, "starting wait")
	if _error := _command.Wait (); _error != nil {
		_waitError = Errorw (0x6f9dfa7d, _error)
	}
//	Logf ('d', 0xa36df40d, "ending wait")
	
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
		return -1, nil, Errorf (0x7cd15552, "invalid state")
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
				return -1, _stdout.Bytes (), Errorw (0xc8553b48, _waitError)
			}
		} else {
			return -1, nil, Errorw (0x4b785e1d, _waitError)
		}
	} else {
		return 0, _stdout.Bytes (), nil
	}
}

