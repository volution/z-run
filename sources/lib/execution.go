

package zrun


import "bytes"
import "fmt"
import "io/ioutil"
import "os"
import "os/exec"
import "path"
import "strings"
import "strconv"
import "syscall"




func prepareExecution (_libraryUrl string, _libraryFingerprint string, _interpreter string, _scriptlet *Scriptlet, _includeArguments bool, _context *Context) (*exec.Cmd, []int, *Error) {
	
	if _interpreter == "" {
		_interpreter = _scriptlet.Interpreter
	}
	
	if _scriptlet.ContextFingerprint != "" {
		if _scriptlet.Context == nil {
			return nil, nil, errorf (0x93547e3b, "invalid store")
		}
	}
	
	var _scriptletExecutablePaths []string
	var _scriptletEnvironment map[string]string
	if _scriptlet.Context != nil {
		_scriptletExecutablePaths = _scriptlet.Context.ExecutablePaths
		_scriptletEnvironment = _scriptlet.Context.Environment
	}
	
	return prepareExecution_0 (
			
			_libraryUrl,
			_libraryFingerprint,
			
			_interpreter,
			_scriptlet.InterpreterExecutable,
			_scriptlet.InterpreterArguments,
			_scriptlet.InterpreterEnvironment,
			
			_scriptlet.Fingerprint,
			_scriptlet.Label,
			_scriptlet.Body,
			
			_scriptlet.Source.Path,
			_scriptlet.Source.LineStart,
			_scriptlet.Source.LineEnd,
			
			_scriptletExecutablePaths,
			_scriptletEnvironment,
			
			_includeArguments,
			
			_context.selfExecutable,
			_context.cleanArguments,
			_context.cleanEnvironment,
			
			_context.workspace,
			_context.executablePaths,
			_context.cacheRoot,
		)
}




func prepareExecution_0 (
			
			_libraryUrl string,
			_libraryFingerprint string,
			
			_interpreter string,
			_scriptletInterpreterExecutable string,
			_scriptletInterpreterArguments []string,
			_scriptletInterpreterEnvironment map[string]string,
			
			_scriptletFingerprint string,
			_scriptletLabel string,
			_scriptletBody string,
			
			_scriptletSourcePath string,
			_scriptletSourceLineStart uint,
			_scriptletSourceLineEnd uint,
			
			_scriptletExecutablePaths []string,
			_scriptletEnvironment map[string]string,
			
			_includeArguments bool,
			
			_selfExecutable string,
			_cleanArguments []string,
			_cleanEnvironment map[string]string,
			
			_contextWorkspace string,
			_contextExecutablePaths []string,
			_contextCacheRoot string,
			
		)
		(*exec.Cmd, []int, *Error)
{
	
	var _interpreterExecutable string
	var _interpreterArguments []string = make ([]string, 0, len (_cleanArguments) + 16)
	var _interpreterEnvironment map[string]string
	var _interpreterAllowsArguments = false
	
	var _executablePaths []string = make ([]string, 0, 128)
	var _environment map[string]string = make (map[string]string, 128)
	
	switch _interpreter {
		
		case "<exec>", "<bash+>", "<python3+>" :
			_interpreterAllowsArguments = true
		
		case "<print>" :
			_interpreterAllowsArguments = false
		
		case "<template>" :
			_interpreterAllowsArguments = true
		
		case "<menu>" :
			_interpreterAllowsArguments = false
		
		case "<go>", "<go+>" :
			_interpreterAllowsArguments = true
		
		default :
			return nil, nil, errorf (0x0873f2db, "unknown scriptlet interpreter `%s` for `%s`", _interpreter, _scriptletLabel)
	}
	
	if _includeArguments && (len (_cleanArguments) > 0) && !_interpreterAllowsArguments {
		return nil, nil, errorf (0x4ef9e048, "unexpected arguments")
	}
	
	var _interpreterScriptInput int
	var _interpreterScriptOutput *os.File
	var _interpreterScriptUnused = false
	
	if _interpreterScriptInput_0, _interpreterScriptOutput_0, _error := createPipe (len (_scriptletBody), _contextCacheRoot); _error == nil {
		_interpreterScriptInput = _interpreterScriptInput_0
		_interpreterScriptOutput = _interpreterScriptOutput_0
	} else {
		return nil, nil, _error
	}
	
	_interpreterScriptBuffer := bytes.NewBuffer (nil)
	_interpreterScriptBuffer.Grow (128 * 1024)
	
	switch _interpreter {
		
		case "<exec>" :
			_interpreterExecutable = _scriptletInterpreterExecutable
			_interpreterEnvironment = _scriptletInterpreterEnvironment
			_interpreterArguments = append (
					_interpreterArguments,
					_scriptletInterpreterExecutable,
				)
			_interpreterArguments = append (
					_interpreterArguments,
					_scriptletInterpreterArguments ...
				)
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (_scriptletBody)
		
		case "<bash+>" :
			_interpreterExecutable = _scriptletInterpreterExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:bash] [%s]", _scriptletLabel),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (embeddedBashProlog)
			_interpreterScriptBuffer.WriteString (fmt.Sprintf ("exec %d<&-\n", _interpreterScriptInput))
			_interpreterScriptBuffer.WriteString (_scriptletBody)
		
		case "<python3+>" :
			_interpreterExecutable = _scriptletInterpreterExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:python3] [%s]", _scriptletLabel),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (embeddedPython3Prolog)
			_interpreterScriptBuffer.WriteString (fmt.Sprintf (
					"Z._scriptlet_begin_from_fd (%d, %s, %s, %d, %d)\n",
					_interpreterScriptInput,
					strconv.QuoteToASCII (_scriptletLabel),
					strconv.QuoteToASCII (_scriptletSourcePath),
					_scriptletSourceLineStart,
					_scriptletSourceLineEnd,
				))
			_interpreterScriptBuffer.WriteString (_scriptletBody)
		
		case "<print>" :
			_interpreterExecutable = "cat"
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:print] [%s]", _scriptletLabel),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (_scriptletBody)
		
		case "<template>" :
			_interpreterExecutable = _selfExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					"[z-run:template]",
					fmt.Sprintf (":: %s", _scriptletLabel),
				)
			_interpreterScriptUnused = true
		
		case "<go>", "<go+>" :
			
			_goFingerprint := _scriptletFingerprint
			_goSource := path.Join (_contextCacheRoot, _goFingerprint + ".go")
			_goExecutable := path.Join (_contextCacheRoot, _goFingerprint + ".exec")
			
			if _, _error := os.Stat (_goExecutable); _error == nil {
				// PASS
			} else if os.IsNotExist (_error) {
				
				_interpreterScriptBuffer.WriteString ("package main\n");
				if _interpreter == "<go+>" {
					_lines := strings.Split (_scriptletBody, "\n")
					for _index, _line := range _lines {
						if strings.HasPrefix (_line, "import ") {
							_interpreterScriptBuffer.WriteString (_line)
							_interpreterScriptBuffer.WriteString ("\n")
						} else {
							_lines = _lines[_index:]
							break
						}
					}
					_interpreterScriptBuffer.WriteString (embeddedGoProlog)
					_interpreterScriptBuffer.WriteString ("\nfunc main () () {\n")
					for _, _line := range _lines {
						_interpreterScriptBuffer.WriteString (_line)
						_interpreterScriptBuffer.WriteString ("\n")
					}
					_interpreterScriptBuffer.WriteString ("\n}\n")
				} else {
					_interpreterScriptBuffer.WriteString (_scriptletBody)
				}
				
				_goSourceTmp := path.Join (_contextCacheRoot, generateRandomToken () + ".tmp")
				if _error := ioutil.WriteFile (_goSourceTmp, _interpreterScriptBuffer.Bytes (), 0600); _error != nil {
					return nil, nil, errorw (0x55976c12, _error)
				}
				if _error := os.Rename (_goSourceTmp, _goSource); _error != nil {
					return nil, nil, errorw (0x5367f11a, _error)
				}
				
				_goExecutableTmp := path.Join (_contextCacheRoot, generateRandomToken () + ".tmp")
				
				_goRoot := path.Join (_contextCacheRoot, "go")
				_goCache := path.Join (_goRoot, "cache")
				_goTmp := path.Join (_goRoot, "tmp")
				for _, _mkdirPath := range []string { _goRoot, _goCache, _goTmp } {
					if _error := os.Mkdir (_mkdirPath, 0700); _error != nil && ! os.IsExist (_error) {
						return nil, nil, errorw (0x5097b00d, _error)
					}
				}
				
				_goExec := ""
				if _goExec_0, _error := resolveExecutable ("go", _contextExecutablePaths); _error == nil {
					_goExec = _goExec_0
				} else {
					return nil, nil, _error
				}
				
				_goBuild := & exec.Cmd {
						Path : _goExec,
						Args : []string {
								_goExec, "build",
								"-o", _goExecutableTmp,
								"-ldflags", "-s -w",
								"--",
								_goSource,
							},
						Dir : _goRoot,
						Env : []string {
								"GO111MODULE=off",
								"GOPATH=" + _goRoot,
								"GOCACHE=" + _goCache,
								"GOTMPDIR=" + _goTmp,
								"TMPDIR=" + _goTmp,
							},
						Stdin : nil,
						Stdout : nil,
						Stderr : os.Stderr,
					}
				
				if _error := _goBuild.Run (); _error == nil {
					if _error := os.Rename (_goExecutableTmp, _goExecutable); _error != nil {
						return nil, nil, errorw (0xbeffd67b, _error)
					}
				} else {
					_ = os.Remove (_goExecutableTmp)
					return nil, nil, errorw (0x72eb9cad, _error)
				}
				
			} else {
				return nil, nil, errorw (0x46248f88, _error)
			}
			
			_interpreterExecutable = _goExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:go] [%s]", _scriptletLabel),
				)
			_interpreterScriptUnused = true
		
		case "<template-raw>" :
			_interpreterExecutable = _selfExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:template-raw] [%s]", _scriptletLabel),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (_scriptletBody)
		
		case "<menu>" :
			_interpreterExecutable = _selfExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:menu] [%s]", _scriptletLabel),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (_scriptletBody)
		
		default :
			panic (0xe95f68a0)
	}
	
	
	var _descriptors []int
	if ! _interpreterScriptUnused {
		
//		logf ('d', 0xedfcf88b, "\n----------\n%s----------\n", _interpreterScriptBuffer.Bytes ())
		
		if _, _error := _interpreterScriptBuffer.WriteTo (_interpreterScriptOutput); _error == nil {
			_interpreterScriptOutput.Close ()
		} else {
			syscall.Close (_interpreterScriptInput)
			_interpreterScriptOutput.Close ()
			return nil, nil, errorw (0xf789ed3f, _error)
		}
		
		_descriptors = []int {
				_interpreterScriptInput,
			}
		
	} else {
		syscall.Close (_interpreterScriptInput)
		_interpreterScriptOutput.Close ()
	}
	
	if _includeArguments {
		_interpreterArguments = append (_interpreterArguments, _cleanArguments ...)
	}
	
	_executablePaths = append (_executablePaths, _scriptletExecutablePaths ...)
	_executablePaths = append (_executablePaths, _contextExecutablePaths ...)
	
	_environment["PATH"] = strings.Join (_executablePaths, string (os.PathListSeparator))
	_environment["ZRUN_WORKSPACE"] = _contextWorkspace
	_environment["ZRUN_LIBRARY_CACHE"] = _libraryUrl
	_environment["ZRUN_FINGERPRINT"] = _libraryFingerprint
	
	_interpreterEnvironment_0 := processEnvironment_0 (
			_selfExecutable,
			_cleanEnvironment,
			[]map[string]string {
					_interpreterEnvironment,
					_scriptletEnvironment,
					_environment,
				},
		)
	
	if _interpreterExecutable_0, _error := resolveExecutable (_interpreterExecutable, _executablePaths); _error == nil {
		_interpreterExecutable = _interpreterExecutable_0
	} else {
		syscall.Close (_interpreterScriptInput)
		_interpreterScriptOutput.Close ()
		return nil, nil, _error
	}
	
	_command := & exec.Cmd {
			Path : _interpreterExecutable,
			Args : _interpreterArguments,
			Env : _interpreterEnvironment_0,
			Dir : _contextWorkspace,
		}
	
//	logf ('d', 0xcc6d38ba, "%v", _command.Path)
//	logf ('d', 0xdb26cbac, "%v", _command.Args[0])
//	logf ('d', 0x7b0c717d, "%v", _command.Args[1:])
//	logf ('d', 0xaa0b151d, "%v", _command.Env)
	
	return _command, _descriptors, nil
}




func executeScriptlet (_library LibraryStore, _scriptlet *Scriptlet, _fork bool, _context *Context) (*Error) {
	
	if _scriptlet.Interpreter == "<template>" {
		return executeTemplate (_library, _scriptlet, _context, os.Stdout)
	}
	
	var _libraryFingerprint string
	if _libraryFingerprint_0, _error := _library.Fingerprint (); _error == nil {
		_libraryFingerprint = _libraryFingerprint_0
	} else {
		return _error
	}
	
	return executeScriptlet_0 (_library.Url (), _libraryFingerprint, _scriptlet, _fork, _context)
}




func executeScriptlet_0 (_libraryUrl string, _libraryFingerprint string, _scriptlet *Scriptlet, _fork bool, _context *Context) (*Error) {
	
	var _command *exec.Cmd
	var _descriptors []int
	if _command_0, _descriptors_0, _error := prepareExecution (_libraryUrl, _libraryFingerprint, "", _scriptlet, true, _context); _error == nil {
		_command = _command_0
		_descriptors = _descriptors_0
	} else {
		return _error
	}
	
	_closeDescriptors := func () () {
		for _, _descriptor := range _descriptors {
			syscall.Close (_descriptor)
		}
	}
	
	if _command.Dir != "" {
		if _error := os.Chdir (_command.Dir); _error != nil {
			return errorw (0xe4bab179, _error)
		}
	}
	if _command.Stdin != nil {
		_closeDescriptors ()
		return errorf (0x78cfda21, "invalid state")
	}
	if _command.Stdout != nil {
		_closeDescriptors ()
		return errorf (0xf9a9dc74, "invalid state")
	}
	if _command.Stderr != nil {
		_closeDescriptors ()
		return errorf (0xf887025f, "invalid state")
	}
	if _command.ExtraFiles != nil {
		_closeDescriptors ()
		return errorf (0x50354e63, "invalid state")
	}
	if (_command.Process != nil) || (_command.ProcessState != nil) {
		_closeDescriptors ()
		return errorf (0x9d640d1e, "invalid state")
	}
	
	if ! _fork {
		
		if _error := syscall.Exec (_command.Path, _command.Args, _command.Env); _error != nil {
			_closeDescriptors ()
			return errorw (0x99b54af1, _error)
		} else {
			panic (0xb6dfe17e)
		}
		
	} else {
		
		for _, _descriptor := range _descriptors {
			_command.ExtraFiles = append (_command.ExtraFiles, os.NewFile (uintptr (_descriptor), ""))
		}
		
		_command.Stdin = os.Stdin
		_command.Stdout = os.Stdout
		_command.Stderr = os.Stderr
		
		_waitError := _command.Run ()
		
		_closeDescriptors ()
		
		if _waitError != nil {
			if _command.ProcessState.Exited () {
				if _exitCode := _command.ProcessState.ExitCode (); _exitCode >= 0 {
					return errorf (0xa10d5811, "spawn `%s` failed with status `%d`", _scriptlet.Label, _exitCode)
				} else {
					return errorf (0x9cfebeaf, "invalid state")
				}
			} else {
				return errorw (0x07b37e04, _waitError)
			}
		} else {
			return nil
		}
	}
}

