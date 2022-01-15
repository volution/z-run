

package premain


import "fmt"
import "log"
import "os"
import "os/exec"
import "path/filepath"
import "runtime"
import "runtime/debug"
import "sort"
import "strings"
import "syscall"
import "unicode"
import "unicode/utf8"

import isatty "github.com/mattn/go-isatty"

import embedded "github.com/cipriancraciun/z-run/embedded"

import . "github.com/cipriancraciun/z-run/lib/run"
import . "github.com/cipriancraciun/z-run/lib/common"




func PreMain () () {
	
	
	runtime.GOMAXPROCS (1)
	debug.SetMaxThreads (16)
	debug.SetMaxStack (128 * 1024)
	debug.SetGCPercent (500)
	
	
	log.SetFlags (0)
	
	
	var _executable0 string
	if _executable_0, _error := os.Executable (); _error == nil {
//		Logf ('d', 0xda96ed44, "%s", _executable_0)
		_executable0 = _executable_0
	} else if _executable_0 := os.Getenv ("ZRUN_EXECUTABLE"); _executable_0 != "" {
		_executable0 = _executable_0
	} else {
		panic (AbortError (Errorf (0x905621f4, "can't resolve `z-run` executable")))
	}
	
	
	var _executable string
	if _executable_0, _error := filepath.EvalSymlinks (_executable0); _error == nil {
		_executable = _executable_0
	} else {
		panic (AbortError (Errorw (0x127e013a, _error)))
	}
	
	
	for _, _variable := range os.Environ () {
		_accept := true
		if _splitIndex := strings.IndexByte (_variable, '='); _splitIndex >= 0 {
			_variable = _variable[:_splitIndex]
			if _variable == "" {
				_accept = false
			}
		} else {
			_accept = false
		}
		if _accept && strings.HasPrefix (_variable, "HYPERFINE_") {
			_accept = false
		}
		if !_accept {
			if _error := os.Unsetenv (_variable); _error != nil {
				panic (AbortError (Errorw (0x14024051, _error)))
			}
		}
	}
	
	
	
	
	if (len (os.Args) == 2) {
		_argument := os.Args[1]
		
		if (_argument == "--version") || (_argument == "-v") {
			
			fmt.Fprintf (os.Stdout, "* version       : %s\n", BUILD_VERSION)
			fmt.Fprintf (os.Stdout, "* executable    : %s\n", os.Args[0])
			fmt.Fprintf (os.Stdout, "* build target  : %s, %s-%s, %s, %s\n", BUILD_TARGET, BUILD_TARGET_OS, BUILD_TARGET_ARCH, BUILD_COMPILER_VERSION, BUILD_COMPILER_TYPE)
			fmt.Fprintf (os.Stdout, "* build number  : %s, %s\n", BUILD_NUMBER, BUILD_TIMESTAMP)
			fmt.Fprintf (os.Stdout, "* code & issues : %s\n", PROJECT_URL)
			fmt.Fprintf (os.Stdout, "* sources git   : %s\n", BUILD_GIT_HASH)
			fmt.Fprintf (os.Stdout, "* sources hash  : %s\n", BUILD_SOURCES_HASH)
			fmt.Fprintf (os.Stdout, "* uname node    : %s\n", UNAME_NODE)
			fmt.Fprintf (os.Stdout, "* uname system  : %s, %s, %s\n", UNAME_SYSTEM, UNAME_RELEASE, UNAME_MACHINE)
			
			os.Exit (0)
			panic (0x66203ba4)
		}
		
		if _argument == "--sources-md5" {
			if _, _error := os.Stdout.WriteString (embedded.BuildSourcesMd5); _error != nil {
				panic (AbortError (Errorw (0x7471032d, _error)))
			}
			os.Exit (0)
			panic (0x9aaf5529)
		}
		
		if _argument == "--sources-cpio" {
			if _, _error := os.Stdout.Write (embedded.BuildSourcesCpioGz); _error != nil {
				panic (AbortError (Errorw (0x8034bf3e, _error)))
			}
			os.Exit (0)
			panic (0x7dca7f0e)
		}
		
		if (_argument == "--manual") || (_argument == "--manual-text") || (_argument == "--manual-html") || (_argument == "--manual-man") {
			_manual := ""
			switch _argument {
				case "--manual", "--manual-text" :
					_manual = embedded.ManualTxt
				case "--manual-html" :
					_manual = embedded.ManualHtml
				case "--manual-man" :
					_manual = embedded.ManualMan
				default :
					panic (0x41b79a1d)
			}
			fmt.Fprint (os.Stdout, _manual)
			os.Exit (0)
			panic (0x1005584b)
		}
		
		if (_argument == "--help") || (_argument == "-h") {
			fmt.Fprint (os.Stdout, embedded.ManualTxt)
			os.Exit (0)
			panic (0xec70ce24)
		}
		
		
		if (_argument == "--shell") || (_argument == "--shell-untainted") {
			
			if _error := CheckTerminal (); _error != nil {
				Logf ('e', 0xf2f72641, "stdin, stdout or stderr are not a TTY;  aborting!")
				os.Exit (1)
				panic (0x67e03fa2)
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
					os.Exit (1)
					panic (0x5a0fe24e)
			}
			
			for _, _chunk := range _chunks {
				if _, _error := os.Stdout.Write ([]byte (_chunk)); _error != nil {
					panic (AbortError (Errorw (0xc08e9e5a, _error)))
				}
			}
			
			os.Exit (0)
			panic (0xda66de5d)
		}
	}
	
	
	
	
	_preMainContext := & PreMainContext {}
	
	_preMainContext.Executable0 = _executable0
	_preMainContext.Executable = _executable
	
	
	if os.Getenv ("ZRUN_EXECUTABLE") == "" {
		switch runtime.GOOS {
			
			case "linux" :
				// NOP
			
			case "darwin" :
				// NOP
			
			case "freebsd", "openbsd" :
				// NOP
			
			case "netbsd", "dragonfly" :
				Logf ('i', 0xc8f30933, "this tool was not tested on your OS;  please be cautions!")
			
			default :
				Logf ('e', 0xcdd5f570, "this tool was not tested on your OS;  it is highly unlikely that it will work;  aborting!")
				os.Exit (1)
				panic (0x9c080b95)
		}
	}
	
	
	_argument0 := os.Args[0]
	_arguments := append ([]string (nil), os.Args[1:] ...)
	
	if _argument0 == _executable0 {
		_argument0 = _executable
	} else if (_argument0 != _executable) && ! strings.HasPrefix (_argument0, "[") {
		if _executable_0, _error := filepath.EvalSymlinks (_argument0); _error == nil {
			if _executable == _executable_0 {
				_argument0 = _executable
			}
		}
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
	
	
	_environment := make (map[string]string, 128)
	_preMainContext.Environment = make ([]string, 0, 128)
	for _, _variable := range os.Environ () {
		_preMainContext.Environment = append (_preMainContext.Environment, _variable)
		
		if _splitIndex := strings.IndexByte (_variable, '='); _splitIndex >= 0 {
			
			_name := _variable[:_splitIndex]
			_value := _variable[_splitIndex + 1:]
			
			_nameTrimmed := strings.TrimSpace (_name)
			if _name != _nameTrimmed {
				Logf ('w', 0x1d362f26, "invalid environment variable (name has spaces):  `%s`", _name)
				_name = _nameTrimmed
			}
			if strings.IndexFunc (_name, func (r rune) (bool) { return unicode.IsSpace (r) || (r > unicode.MaxASCII) }) >= 0 {
				Logf ('w', 0x81ac6f2e, "invalid environment variable (name is not ASCII):  `%s`", _name)
			}
			
			if _name == "" {
				Logf ('w', 0x0ffb0031, "invalid environment variable (name empty):  `%s`", _variable)
			} else if ! utf8.Valid ([]byte (_name)) {
				Logf ('w', 0x54278534, "invalid environment variable (name invalid UTF-c):  `%s`", _name)
			} else if ! utf8.Valid ([]byte (_value)) {
				Logf ('w', 0x785ba004, "invalid environment variable (value invalid UTF-c):  `%s`", _name)
			} else if _value == "" {
//				Logf ('w', 0xfe658d34, "invalid environment variable (value empty):  `%s`", _name)
			} else if _, _exists := _environment[_name]; _exists {
				Logf ('w', 0x7e7e41a5, "invalid environment variable (name duplicate):  `%s`", _name)
			} else {
				_environment[_nameTrimmed] = _value
			}
			
		} else {
			Logf ('w', 0xe745517c, "invalid environment variable (missing `=`):  `%s`", _variable)
		}
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
				os.Exit (1)
			}
			if _argument00, _error := filepath.EvalSymlinks (_argument0); (_error != nil) || (_argument00 != _executable) {
				Logf ('e', 0xf1f1a024, "invalid argument0: `%s`, expected `%s`;  aborting!", _argument0, _executable)
				os.Exit (1)
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
			os.Exit (1)
	}
	
	if _error := Main (_executable, _argument0, _arguments, _environment, "", ""); _error == nil {
		os.Exit (0)
		panic (0xe0e1c1a1)
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




func CheckTerminal () (*Error) {
	if ! isatty.IsTerminal (os.Stdin.Fd ()) {
		return Errorf (0x05d60b72, "stdin is not a TTY")
	}
	if ! isatty.IsTerminal (os.Stdout.Fd ()) {
		return Errorf (0xc432630a, "stdout is not a TTY")
	}
	if ! isatty.IsTerminal (os.Stderr.Fd ()) {
		return Errorf (0x77924518, "stderr is not a TTY")
	}
	return nil
}

