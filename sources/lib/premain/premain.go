

package premain


import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "sort"
import "strings"
import "syscall"

import embedded "github.com/cipriancraciun/z-run/embedded"

import . "github.com/cipriancraciun/z-run/lib/run"
import . "github.com/cipriancraciun/z-run/lib/fzf"
import . "github.com/cipriancraciun/z-run/lib/input"
import . "github.com/cipriancraciun/z-run/lib/mainlib"
import . "github.com/cipriancraciun/z-run/lib/common"




func PreMain () () {
	
	
	
	
	if _error := InitializeMainRuntime (); _error != nil {
		panic (AbortError (_error))
	}
	
	if _error := CleanMainEnvironment (); _error != nil {
		panic (AbortError (_error))
	}
	
	var _executable0, _executable string
	if _executable0_0, _executable_0, _error := ResolveMainExecutable ("z-run", "ZRUN_EXECUTABLE"); _error == nil {
		_executable0 = _executable0_0
		_executable = _executable_0
	} else {
		panic (AbortError (_error))
	}
	
	if _error := InterceptMainSpecialFlags ("z-run", _executable0, _executable, embedded.ManualTxt, embedded.ManualHtml, embedded.ManualMan); _error != nil {
		panic (AbortError (_error))
	}
	
	
	
	
	if (len (os.Args) == 2) {
		_argument := os.Args[1]
		
		
		if (_argument == "--shell") || (_argument == "--shell-untainted") {
			
			if _error := CheckMainTerminal (); _error != nil {
				Logf ('e', 0xf2f72641, "stdin, stdout or stderr are not a TTY;  aborting!")
				panic (ExitMainFailed ())
			}
			
			if _argument == "--shell-untainted" {
				os.Unsetenv ("ZRUN_WORKSPACE")
				os.Unsetenv ("ZRUN_LIBRARY_SOURCE")
				os.Unsetenv ("ZRUN_LIBRARY_URL")
				os.Unsetenv ("ZRUN_LIBRARY_IDENTIFIER")
				os.Unsetenv ("ZRUN_LIBRARY_FINGERPRINT")
			}
			
			_rc := "\n" + embedded.BashShellRc + "\n" + embedded.BashShellFunctions + "\n"
			
			_input, _output, _error := CreatePipe (len (_rc) + 256, "/tmp")
			if _error != nil {
				panic (AbortError (_error))
			}
			_rc += fmt.Sprintf ("exec %d<&-\n", _input)
			_rc += "printf -- '\n%s\n' '---- [z-run:shell] -------------------------------------------------------------' >&2\n\n"
			if _, _error := _output.Write ([]byte (_rc)); _error != nil {
				panic (AbortError (Errorw (0xc58e3fe6, _error)))
			}
			if _error := _output.Close (); _error != nil {
				panic (AbortError (Errorw (0x8741d077, _error)))
			}
			
			var _bash string
			if _bash_0, _error := exec.LookPath ("bash"); _error == nil {
				_bash = _bash_0
			} else if _bash_0, _error := ResolveExecutable ("bash", []string { "/usr/local/bin", "/usr/bin", "/bin" }); _error == nil {
				_bash = _bash_0
			} else {
				_bash = "/bin/bash"
			}
			
			_arguments := []string {
					"[z-run:shell]",
					"--noprofile",
					"--rcfile", fmt.Sprintf ("/dev/fd/%d", _input),
				}
			
			_environment := append ([]string (nil),  os.Environ () ...)
			_environment = append (_environment, "ZRUN_EXECUTABLE=" + _executable)
			
			if _error := syscall.Exec (_bash, _arguments, _environment); _error != nil {
				panic (AbortError (Errorf (0x8598d4c0, "failed to exec `%s`  //  %v", _bash, _error)))
			}
			panic (0xf4813cc2)
		}
		
		
		if strings.HasPrefix (_argument, "--export=") {
			
			_what := _argument[len ("--export=") :]
			var _chunks []string
			
			switch _what {
				
				case "shell-rc", "shell-rc-only", "shell-functions", "bash-prolog", "bash+-prolog", "python3+-prolog" :
					_chunks = append (_chunks,
							"################################################################################\n",
							"######## z-run v" + BUILD_VERSION + "\n",
						)
					switch _what {
						case "shell-rc" :
							_chunks = append (_chunks,
									"################################################################################\n\n",
									strings.Trim (strings.TrimPrefix (strings.Trim (embedded.BashShellRc, "\n"), "#!/dev/null\n"), "\n"),
									"\n\n################################################################################\n",
									"################################################################################\n\n",
									strings.Trim (strings.TrimPrefix (strings.Trim (embedded.BashShellFunctions, "\n"), "#!/dev/null\n"), "\n"),
									"\n\n################################################################################\n",
								)
						case "shell-rc-only" :
							_chunks = append (_chunks,
									"################################################################################\n\n",
									strings.Trim (strings.TrimPrefix (strings.Trim (embedded.BashShellRc, "\n"), "#!/dev/null\n"), "\n"),
									"\n\n################################################################################\n",
								)
						case "shell-functions" :
							_chunks = append (_chunks,
									"################################################################################\n\n",
									strings.Trim (strings.TrimPrefix (strings.Trim (embedded.BashShellFunctions, "\n"), "#!/dev/null\n"), "\n"),
									"\n\n################################################################################\n",
								)
						case "bash-prolog" :
							_chunks = append (_chunks,
									"################################################################################\n\n",
									strings.Trim (strings.TrimPrefix (strings.Trim (embedded.BashProlog0, "\n"), "#!/dev/null\n"), "\n"),
									"\n\n################################################################################\n",
								)
						case "bash+-prolog" :
							_chunks = append (_chunks,
									"################################################################################\n\n",
									strings.Trim (strings.TrimPrefix (strings.Trim (embedded.BashProlog, "\n"), "#!/dev/null\n"), "\n"),
									"\n\n################################################################################\n",
								)
						case "python3+-prolog" :
							_chunks = append (_chunks,
									strings.Trim (strings.TrimPrefix (strings.Trim (embedded.Python3Prolog, "\n"), "#!/dev/null\n"), "\n"),
								)
						default :
							panic (0xaec5d2dd)
					}
				
				case "go+-prolog" :
					_chunks = append (_chunks, embedded.GoProlog)
				
				default :
					Logf ('e', 0xa269e851, "invalid export `%s`;  aborting!", _what)
					panic (ExitMainFailed ())
			}
			
			for _, _chunk := range _chunks {
				if _, _error := os.Stdout.Write ([]byte (_chunk)); _error != nil {
					panic (AbortError (Errorw (0xc08e9e5a, _error)))
				}
			}
			
			panic (ExitMainSucceeded ())
		}
	}
	
	
	
	
	_preMainContext := & PreMainContext {}
	
	_preMainContext.Executable0 = _executable0
	_preMainContext.Executable = _executable
	
	var _argument0 string
	var _arguments []string
	if _argument0_0, _arguments_0, _error := ResolveMainArguments (_executable0, _executable); _error == nil {
		_argument0 = _argument0_0
		_arguments = _arguments_0
	} else {
		panic (AbortError (_error))
	}
	
	_preMainContext.Argument0 = _argument0
	_preMainContext.Arguments = append ([]string (nil), _arguments ...)
	
	if strings.HasPrefix (_argument0, "[z-run:menu] ") {
		_argument0 = "[z-run:menu]"
	} else if strings.HasPrefix (_argument0, "[z-run:print] ") {
		_argument0 = "[z-run:print]"
	} else if strings.HasPrefix (_argument0, "[z-run:template] ") {
		_argument0 = "[z-run:template]"
	} else if strings.HasPrefix (_argument0, "[z-run:starlark] ") {
		_argument0 = "[z-run:starlark]"
	}
	
	
	var _environment map[string]string
	if _environmentMap_0, _environmentList_0, _error := ResolveMainEnvironment (); _error == nil {
		_environment = _environmentMap_0
		_preMainContext.Environment = _environmentList_0
	} else {
		panic (AbortError (_error))
	}
	
	// FIXME:  This is for OpenBSD which doesn't have a way to find `os.Executable` outside of `arg0`...
	if _, _exists := _environment["ZRUN_EXECUTABLE"]; !_exists {
		_environment["ZRUN_EXECUTABLE"] = _executable
		_preMainContext.Environment = append (_preMainContext.Environment, "ZRUN_EXECUTABLE=" + _executable)
	}
	
//	Logf ('d', 0x06cd45f9, "self-executable0: %s", _executable0)
//	Logf ('d', 0x256b2c94, "self-executable: %s", _executable)
//	Logf ('d', 0xb59e4f73, "self-argument0: %s", _argument0)
//	Logf ('d', 0xf7d65090, "self-arguments: %s", _arguments)
//	Logf ('d', 0x7a411846, "self-environment: %s", _environment)
	
	PreMainContextGlobal = _preMainContext
	
	os.Args = append ([]string {"z-run"}, _arguments ...)
	os.Clearenv ()
	for _name, _value := range _environment {
		os.Setenv (_name, _value)
	}
	
	
	_argument0IsTool := true
	
	switch _argument0 {
		
		case "[z-run]" :
			_argument0 = "[z-run]"
		
		case "[z-run:library]" :
			_argument0 = "[z-run]"
		
		case "[z-run:menu]" :
			// NOP
		
		case "[z-run:print]" :
			// NOP
		
		case "[z-run:template]" :
			// NOP
		
		case "[z-run:starlark]" :
			// NOP
		
		case "[z-run:scriptlet]" :
			// NOP
		
		case "[z-run:scriptlet-exec]" :
			// NOP
		
		case "[z-run:input]" :
			// NOP
		
		case "[z-run:select]" :
			// NOP
		
		case "[z-run:fzf]" :
			// NOP
		
		case "z-run", "zrun", "x-run", "xrun", "_" :
			_argument0IsTool = false
		
		default :
			if strings.HasPrefix (_argument0, "[z-run:") {
				Logf ('e', 0xf6274ed5, "invalid argument0: `%s`;  aborting!", _argument0)
				panic (ExitMainFailed ())
			}
			if _argument00, _error := filepath.EvalSymlinks (_argument0); (_error != nil) || (_argument00 != _executable) {
				Logf ('e', 0xf1f1a024, "invalid argument0: `%s`, expected `%s`;  aborting!", _argument0, _executable)
				panic (ExitMainFailed ())
			}
			_argument0IsTool = false
	}
	
	
	if !_argument0IsTool && len (_arguments) == 0 {
		
		_argument0 = "[z-run]"
		
	} else if !_argument0IsTool && len (_arguments) >= 1 {
		
		_delegateExecutable := ""
		_delegateArgument0 := ""
		var _delegateArguments []string
		var _delegateEnvironment []string
		
		switch _arguments[0] {
			
			case "--scriptlet" :
				_argument0 = "[z-run:scriptlet]"
				_arguments = _arguments[1:]
			
			case "--scriptlet-exec" :
				_argument0 = "[z-run:scriptlet-exec]"
				_arguments = _arguments[1:]
			
			case "--fzf" :
				_delegateExecutable = _executable
				_delegateArgument0 = "[z-run:fzf]"
				_delegateArguments = _arguments[1:]
			
			case "--input" :
				_argument0 = "[z-run:input]"
				_arguments = _arguments[1:]
			
			case "--select" :
				_argument0 = "[z-run:select]"
				_arguments = _arguments[1:]
			
			case "--print" :
				_argument0 = "[z-run:print]"
				_arguments = _arguments[1:]
			
			case "--template" :
				_argument0 = "[z-run:template]"
				_arguments = _arguments[1:]
			
			case "--starlark" :
				_argument0 = "[z-run:starlark]"
				_arguments = _arguments[1:]
			
			default :
				_argument0 = "[z-run]"
		}
		
		if _delegateExecutable != "" {
			
			for _name, _value := range _environment {
				_delegateEnvironment = append (_delegateEnvironment, _name + "=" + _value)
			}
			sort.Strings (_delegateEnvironment)
			
			_delegateArguments := append ([]string {_delegateArgument0}, _delegateArguments ...)
			
			if _error := syscall.Exec (_delegateExecutable, _delegateArguments, _delegateEnvironment); _error != nil {
				panic (AbortError (Errorf (0x05bd220d, "failed to exec `%s`  //  %v", _delegateExecutable, _error)))
			} else {
				panic (0xe13aab5f)
			}
		}
	}
	
	
	os.Args = append ([]string {"z-run"}, _arguments ...)
	
	
	switch _argument0 {
		
		case "[z-run:scriptlet]" :
			if _error := ScriptletMain (_executable, _arguments, _environment, false); _error != nil {
				panic (AbortError (_error))
			} else {
				panic (0xb305aa74)
			}
		
		case "[z-run:scriptlet-exec]" :
			if _error := ScriptletMain (_executable, _arguments, _environment, true); _error != nil {
				panic (AbortError (_error))
			} else {
				panic (0x8f827319)
			}
		
		case "[z-run:input]" :
			if _error := InputMain (_arguments, _environment); _error != nil {
				panic (AbortError (_error))
			} else {
				panic (0xe62a9355)
			}
		
		case "[z-run:print]" :
			if _error := PrintMain (_executable, _arguments, _environment); _error != nil {
				panic (AbortError (_error))
			} else {
				panic (0xf2084070)
			}
		
		case "[z-run:template]" :
			if _error := TemplateMain (_executable, _arguments, _environment); _error != nil {
				panic (AbortError (_error))
			} else {
				panic (0x32241835)
			}
		
		case "[z-run:starlark]" :
			if _error := StarlarkMain (_executable, _arguments, _environment); _error != nil {
				panic (AbortError (_error))
			} else {
				panic (0xd6f5b038)
			}
		
		case "[z-run:menu]" :
			if _error := MenuMain (_executable, _arguments, _environment); _error != nil {
				panic (AbortError (_error))
			} else {
				panic (0x6b21e0ab)
			}
		
		case "[z-run:select]" :
			if _error := FzfMain (true, _arguments, _environment); _error != nil {
				panic (AbortError (_error))
			} else {
				panic (0x2346ca3f)
			}
		
		case "[z-run:fzf]" :
			if _error := FzfMain (false, _arguments, _environment); _error != nil {
				panic (AbortError (_error))
			} else {
				panic (0xfae3720e)
			}
		
		case "[z-run]" :
			// NOP
		
		default :
			Logf ('e', 0xf965e92e, "invalid argument0: `%s`;  aborting!", _argument0)
			panic (ExitMainFailed ())
	}
	
	if _error := Main (_executable, _argument0, _arguments, _environment, "", ""); _error == nil {
		panic (ExitMainSucceeded ())
	} else {
		panic (AbortError (_error))
	}
}




type PreMainContext struct {
	Executable0 string
	Executable string
	Argument0 string
	Arguments []string
	Environment []string
}

var PreMainContextGlobal *PreMainContext = nil




func PreMainReExecute (_executable string) (*Error) {
	if PreMainContextGlobal == nil {
		return Errorf (0x3d126cd2, "can't switch `z-run`")
	}
	_arguments := make ([]string, 0, len (PreMainContextGlobal.Arguments) + 1)
	if PreMainContextGlobal.Argument0 == PreMainContextGlobal.Executable {
		_arguments = append (_arguments, _executable)
	} else {
		_arguments = append (_arguments, PreMainContextGlobal.Argument0)
	}
	_arguments = append (_arguments, PreMainContextGlobal.Arguments ...)
//	Logf ('i', 0x91038b92, "switching `z-run` to: `%s`...", _executable)
	for _index, _pair := range PreMainContextGlobal.Environment {
		if strings.HasPrefix (_pair, "ZRUN_EXECUTABLE=") {
			PreMainContextGlobal.Environment[_index] = "ZRUN_EXECUTABLE=" + _executable
			break
		}
	}
	_error := syscall.Exec (
			_executable,
			_arguments,
			PreMainContextGlobal.Environment,
		)
	return Errorf (0x3d993836, "failed to exec `%s`  //  %v", _executable, _error)
}

