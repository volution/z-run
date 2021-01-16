

package zrun


import "fmt"
import "log"
import "os"
import "path/filepath"
import "runtime"
import "runtime/debug"
import "sort"
import "strings"
import "syscall"
import "unicode"
import "unicode/utf8"




func PreMain () () {
	
	
	if (len (os.Args) == 2) {
		if os.Args[1] == "--version" {
			fmt.Fprintf (os.Stdout, "* version       : %s\n", BUILD_VERSION)
			fmt.Fprintf (os.Stdout, "* executable    : %s\n", os.Args[0])
			fmt.Fprintf (os.Stdout, "* build target  : %s, %s-%s, %s\n", BUILD_TARGET, BUILD_TARGET_OS, BUILD_TARGET_ARCH, BUILD_COMPILER)
			fmt.Fprintf (os.Stdout, "* build number  : %s, %s\n", BUILD_NUMBER, BUILD_TIMESTAMP)
			fmt.Fprintf (os.Stdout, "* sources md5   : %s\n", BUILD_SOURCES_MD5)
			fmt.Fprintf (os.Stdout, "* sources git   : %s\n", BUILD_GIT_HASH)
			fmt.Fprintf (os.Stdout, "* code & issues : %s\n", "https://github.com/cipriancraciun/z-run")
			os.Exit (0)
			panic (0x66203ba4)
		}
	}
	
	
	log.SetFlags (0)
	
	
	runtime.GOMAXPROCS (1)
	debug.SetMaxThreads (16)
	debug.SetMaxStack (128 * 1024)
	debug.SetGCPercent (500)
	
	
	var _executable string
	if _executable_0, _error := os.Executable (); _error == nil {
		_executable = _executable_0
	} else {
		panic (abortError (errorw (0x75f2db30, _error)))
	}
	if _executable_0, _error := filepath.EvalSymlinks (_executable); _error == nil {
		_executable = _executable_0
	} else {
		panic (abortError (errorw (0x127e013a, _error)))
	}
	
	
	if os.Getenv ("ZRUN_EXECUTABLE") == "" {
		switch runtime.GOOS {
			
			case "linux" :
				// NOP
			
			case "darwin" :
				// NOP
			
			case "freebsd", "netbsd", "openbsd", "dragonfly" :
				logf ('i', 0xc8f30933, "this tool was not tested on your OS;  please be cautions!")
			
			default :
				logf ('e', 0xcdd5f570, "this tool was not tested on your OS;  it is highly unlikely that it will work;  aborting!")
				os.Exit (1)
				panic (0x9c080b95)
		}
	}
	
	
	_argument0 := os.Args[0]
	_arguments := append ([]string (nil), os.Args[1:] ...)
	
	if strings.HasPrefix (_argument0, "[z-run:menu] ") {
		_argument0 = "[z-run:menu]"
	} else if strings.HasPrefix (_argument0, "[z-run:select] ") {
		_argument0 = "[z-run:select]"
	} else if strings.HasPrefix (_argument0, "[z-run:print] ") {
		_argument0 = "[z-run:print]"
	} else if strings.HasPrefix (_argument0, "[z-run:template] ") {
		_argument0 = "[z-run:template]"
	} else if strings.HasPrefix (_argument0, "[z-run:template-raw] ") {
		_argument0 = "[z-run:template-raw]"
	}
	
	
	_environment := make (map[string]string, 128)
	for _, _variable := range os.Environ () {
		if _splitIndex := strings.IndexByte (_variable, '='); _splitIndex >= 0 {
			
			_name := _variable[:_splitIndex]
			_value := _variable[_splitIndex + 1:]
			
			_nameTrimmed := strings.TrimSpace (_name)
			if _name != _nameTrimmed {
				logf ('w', 0x1d362f26, "invalid environment variable (name has spaces):  `%s`", _name)
				_name = _nameTrimmed
			}
			if strings.IndexFunc (_name, func (r rune) (bool) { return unicode.IsSpace (r) || (r > unicode.MaxASCII) }) >= 0 {
				logf ('w', 0x81ac6f2e, "invalid environment variable (name is not ASCII):  `%s`", _name)
			}
			
			if _name == "" {
				logf ('w', 0x0ffb0031, "invalid environment variable (name empty):  `%s`", _variable)
			} else if ! utf8.Valid ([]byte (_name)) {
				logf ('w', 0x54278534, "invalid environment variable (name invalid UTF-c):  `%s`", _name)
			} else if ! utf8.Valid ([]byte (_value)) {
				logf ('w', 0x785ba004, "invalid environment variable (value invalid UTF-c):  `%s`", _name)
			} else if _value == "" {
//				logf ('w', 0xfe658d34, "invalid environment variable (value empty):  `%s`", _name)
			} else if _, _exists := _environment[_name]; _exists {
				logf ('w', 0x7e7e41a5, "invalid environment variable (name duplicate):  `%s`", _name)
			} else {
				_environment[_nameTrimmed] = _value
			}
			
		} else {
			logf ('w', 0xe745517c, "invalid environment variable (missing `=`):  `%s`", _variable)
		}
	}
	
//	logf ('d', 0x256b2c94, "self-executable: %s", _executable)
//	logf ('d', 0xb59e4f73, "self-argument0: %s", _argument0)
//	logf ('d', 0xf7d65090, "self-arguments: %s", _arguments)
//	logf ('d', 0x7a411846, "self-environment: %s", _environment)
	
	
	os.Args = append ([]string {"z-run"}, _arguments ...)
	os.Clearenv ()
	for _name, _value := range _environment {
		os.Setenv (_name, _value)
	}
	
	
	_argument0IsTool := true
	
	switch _argument0 {
		
		case "[z-run]" :
			_argument0 = "[z-run]"
		
		case "[z-run:print]" :
			_argument0 = "[z-run]"
		
		case "[z-run:template]" :
			_argument0 = "[z-run]"
		
		case "[z-run:template-raw]" :
			// NOP
		
		case "[z-run:menu]" :
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
				logf ('e', 0xf6274ed5, "invalid argument0: `%s`;  aborting!", _argument0)
				os.Exit (1)
			}
			if _argument00, _error := filepath.EvalSymlinks (_argument0); (_error != nil) || (_argument00 != _executable) {
				logf ('e', 0xf1f1a024, "invalid argument0: `%s`;  aborting!", _argument0)
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
				panic (abortError (errorw (0x05bd220d, _error)))
			} else {
				panic (0xe13aab5f)
			}
		}
	}
	
	
	os.Args = append ([]string {"z-run"}, _arguments ...)
	
	
	switch _argument0 {
		
		case "[z-run:input]" :
			if _error := inputMain (_arguments, _environment); _error != nil {
				panic (abortError (_error))
			} else {
				panic (0xe62a9355)
			}
		
		case "[z-run:template-raw]" :
			if _error := templateMain (_arguments, _environment); _error != nil {
				panic (abortError (_error))
			} else {
				panic (0x32241835)
			}
		
		case "[z-run:menu]" :
			if _error := menuMain (_executable, _arguments, _environment); _error != nil {
				panic (abortError (_error))
			} else {
				panic (0x6b21e0ab)
			}
		
		case "[z-run:select]" :
			if _error := fzfMain (true, _arguments, _environment); _error != nil {
				panic (abortError (_error))
			} else {
				panic (0x2346ca3f)
			}
		
		case "[z-run:fzf]" :
			if _error := fzfMain (false, _arguments, _environment); _error != nil {
				panic (abortError (_error))
			} else {
				panic (0xfae3720e)
			}
		
		case "[z-run]" :
			// NOP
		
		default :
			logf ('e', 0xf965e92e, "invalid argument0: `%s`;  aborting!", _argument0)
			os.Exit (1)
	}
	
	if _error := Main (_executable, _argument0, _arguments, _environment, "", ""); _error == nil {
		os.Exit (0)
		panic (0xe0e1c1a1)
	} else {
		panic (abortError (_error))
	}
}

