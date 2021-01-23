

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
	
	var _interpreterExecutable string
	var _interpreterArguments []string = make ([]string, 0, len (_context.cleanArguments) + 16)
	var _interpreterEnvironment map[string]string
	var _interpreterAllowsArguments = false
	
	var _executablePaths []string = make ([]string, 0, 128)
	var _environment map[string]string = make (map[string]string, 128)
	
	if _interpreter == "" {
		_interpreter = _scriptlet.Interpreter
	}
	
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
			return nil, nil, errorf (0x0873f2db, "unknown scriptlet interpreter `%s` for `%s`", _interpreter, _scriptlet.Label)
	}
	
	if _includeArguments && (len (_context.cleanArguments) > 0) && !_interpreterAllowsArguments {
		return nil, nil, errorf (0x4ef9e048, "unexpected arguments")
	}
	
	if _scriptlet.ContextFingerprint != "" {
		if _scriptlet.Context == nil {
			return nil, nil, errorf (0x93547e3b, "invalid store")
		}
	}
	
	var _interpreterScriptInput int
	var _interpreterScriptOutput *os.File
	var _interpreterScriptUnused = false
	
	if _interpreterScriptInput_0, _interpreterScriptOutput_0, _error := createPipe (len (_scriptlet.Body), _context.cacheRoot); _error == nil {
		_interpreterScriptInput = _interpreterScriptInput_0
		_interpreterScriptOutput = _interpreterScriptOutput_0
	} else {
		return nil, nil, _error
	}
	
	_interpreterScriptBuffer := bytes.NewBuffer (nil)
	_interpreterScriptBuffer.Grow (128 * 1024)
	
	switch _interpreter {
		
		case "<exec>" :
			_interpreterExecutable = _scriptlet.InterpreterExecutable
			_interpreterEnvironment = _scriptlet.InterpreterEnvironment
			_interpreterArguments = append (
					_interpreterArguments,
					_scriptlet.InterpreterExecutable,
				)
			_interpreterArguments = append (
					_interpreterArguments,
					_scriptlet.InterpreterArguments ...
				)
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (_scriptlet.Body)
		
		case "<bash+>" :
			_interpreterExecutable = _scriptlet.InterpreterExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:bash] [%s]", _scriptlet.Label),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (embeddedBashProlog)
			_interpreterScriptBuffer.WriteString (fmt.Sprintf ("exec %d<&-\n", _interpreterScriptInput))
			_interpreterScriptBuffer.WriteString (_scriptlet.Body)
		
		case "<python3+>" :
			_interpreterExecutable = _scriptlet.InterpreterExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:python3] [%s]", _scriptlet.Label),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (embeddedPython3Prolog)
			_interpreterScriptBuffer.WriteString (fmt.Sprintf (
					"Z._scriptlet_begin_from_fd (%d, %s, %s, %d, %d)\n",
					_interpreterScriptInput,
					strconv.QuoteToASCII (_scriptlet.Label),
					strconv.QuoteToASCII (_scriptlet.Source.Path),
					_scriptlet.Source.LineStart,
					_scriptlet.Source.LineEnd,
				))
			_interpreterScriptBuffer.WriteString (_scriptlet.Body)
		
		case "<print>" :
			_interpreterExecutable = "cat"
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:print] [%s]", _scriptlet.Label),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (_scriptlet.Body)
		
		case "<template>" :
			_interpreterExecutable = _context.selfExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					"[z-run:template]",
					fmt.Sprintf (":: %s", _scriptlet.Label),
				)
			_interpreterScriptUnused = true
		
		case "<go>", "<go+>" :
			
			_goFingerprint := _scriptlet.Fingerprint
			_goSource := path.Join (_context.cacheRoot, _goFingerprint + ".go")
			_goExecutable := path.Join (_context.cacheRoot, _goFingerprint + ".exec")
			
			if _, _error := os.Stat (_goExecutable); _error == nil {
				// PASS
			} else if os.IsNotExist (_error) {
				
				_interpreterScriptBuffer.WriteString ("package main\n");
				if _interpreter == "<go+>" {
					_lines := strings.Split (_scriptlet.Body, "\n")
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
					_interpreterScriptBuffer.WriteString (_scriptlet.Body)
				}
				
				_goSourceTmp := path.Join (_context.cacheRoot, generateRandomToken () + ".tmp")
				if _error := ioutil.WriteFile (_goSourceTmp, _interpreterScriptBuffer.Bytes (), 0600); _error != nil {
					return nil, nil, errorw (0x55976c12, _error)
				}
				if _error := os.Rename (_goSourceTmp, _goSource); _error != nil {
					return nil, nil, errorw (0x5367f11a, _error)
				}
				
				_goExecutableTmp := path.Join (_context.cacheRoot, generateRandomToken () + ".tmp")
				
				_goRoot := path.Join (_context.cacheRoot, "go")
				_goCache := path.Join (_goRoot, "cache")
				_goTmp := path.Join (_goRoot, "tmp")
				for _, _mkdirPath := range []string { _goRoot, _goCache, _goTmp } {
					if _error := os.Mkdir (_mkdirPath, 0700); _error != nil && ! os.IsExist (_error) {
						return nil, nil, errorw (0x5097b00d, _error)
					}
				}
				
				_goExec := ""
				if _goExec_0, _error := resolveExecutable ("go", _context.executablePaths); _error == nil {
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
					fmt.Sprintf ("[z-run:go] [%s]", _scriptlet.Label),
				)
			_interpreterScriptUnused = true
		
		case "<template-raw>" :
			_interpreterExecutable = _context.selfExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:template-raw] [%s]", _scriptlet.Label),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (_scriptlet.Body)
		
		case "<menu>" :
			_interpreterExecutable = _context.selfExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:menu] [%s]", _scriptlet.Label),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (_scriptlet.Body)
		
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
		_interpreterArguments = append (_interpreterArguments, _context.cleanArguments ...)
	}
	
	var _scriptletEnvironment map[string]string
	if _scriptlet.Context != nil {
		_executablePaths = append (_executablePaths, _scriptlet.Context.ExecutablePaths ...)
		_scriptletEnvironment = _scriptlet.Context.Environment
	}
	
	_executablePaths = append (_executablePaths, _context.executablePaths ...)
	
	_environment["PATH"] = strings.Join (_executablePaths, string (os.PathListSeparator))
	_environment["ZRUN_WORKSPACE"] = _context.workspace
	_environment["ZRUN_LIBRARY_CACHE"] = _libraryUrl
	_environment["ZRUN_FINGERPRINT"] = _libraryFingerprint
	
	_interpreterEnvironment_0 := processEnvironment (
			_context,
			_interpreterEnvironment,
			_scriptletEnvironment,
			_environment,
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
			Dir : _context.workspace,
		}
	
//	logf ('d', 0xcc6d38ba, "%v", _command.Path)
//	logf ('d', 0xdb26cbac, "%v", _command.Args[0])
//	logf ('d', 0x7b0c717d, "%v", _command.Args[1:])
//	logf ('d', 0xaa0b151d, "%v", _command.Env)
	
	return _command, _descriptors, nil
}

