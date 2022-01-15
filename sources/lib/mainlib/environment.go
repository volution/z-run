

package mainlib


import "os"
import "path/filepath"
import "runtime"
import "strings"
import "unicode"
import "unicode/utf8"

import . "github.com/cipriancraciun/z-run/lib/common"




func ResolveMainExecutable (_executableHint string, _environmentHint string) (string, string, *Error) {
	
	var _executable0 string
	var _environmentHinted bool
	if _executable_0, _error := os.Executable (); _error == nil {
//		Logf ('d', 0xda96ed44, "%s", _executable_0)
		_executable0 = _executable_0
	} else if _environmentHint != "" {
		if _executable_0 := os.Getenv (_environmentHint); _executable_0 != "" {
			_environmentHinted = true
			_executable0 = _executable_0
		}
	}
	
	if _executable0 == "" {
		return "", "", Errorf (0x905621f4, "can't resolve `%s` executable", _executableHint)
	}
	
	var _executable string
	if _executable_0, _error := filepath.EvalSymlinks (_executable0); _error == nil {
		_executable = _executable_0
	} else {
		return "", "", Errorw (0x127e013a, _error)
	}
	
	if !_environmentHinted {
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
				return "", "", Errorf (0xcdd5f570, "this tool was not tested on your OS;  it is highly unlikely that it will work;  aborting!")
		}
	}
	
	return _executable0, _executable, nil
}




func ResolveMainArguments (_executable0 string, _executable string) (string, []string, *Error) {
	
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
	
	return _argument0, _arguments, nil
}




func ResolveMainEnvironment () (map[string]string, []string, *Error) {
	
	_environmentMap := make (map[string]string, 128)
	_environmentList := make ([]string, 0, 128)
	
	for _, _variable := range os.Environ () {
		
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
				Logf ('w', 0x54278534, "invalid environment variable (name invalid UTF-8):  `%s`", _name)
			} else if ! utf8.Valid ([]byte (_value)) {
				Logf ('w', 0x785ba004, "invalid environment variable (value invalid UTF-8):  `%s`", _name)
			} else if _value == "" {
//				Logf ('w', 0xfe658d34, "invalid environment variable (value empty):  `%s`", _name)
			} else if _, _exists := _environmentMap[_name]; _exists {
				Logf ('w', 0x7e7e41a5, "invalid environment variable (name duplicate):  `%s`", _name)
			} else {
				_environmentMap[_nameTrimmed] = _value
				_environmentList = append (_environmentList, _name + "=" + _variable)
			}
			
		} else {
			Logf ('w', 0xe745517c, "invalid environment variable (missing `=`):  `%s`", _variable)
		}
	}
	
	return _environmentMap, _environmentList, nil
}




func CleanMainEnvironment () (*Error) {
	
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
				return Errorw (0x14024051, _error)
			}
		}
	}
	
	return nil
}




func ResetMainEnvironment (_argument0 string, _arguments []string, _environment map[string]string) (*Error) {
	
	os.Args = append ([]string {_argument0}, _arguments ...)
	
	os.Clearenv ()
	for _name, _value := range _environment {
		os.Setenv (_name, _value)
	}
	
	return nil
}

