

package zrun


import "path/filepath"
import "strings"
import "unicode/utf8"




func scriptletMain (_selfExecutable string, _arguments []string, _environment map[string]string, _shabang bool) (*Error) {
	
	_label := ""
	_header := ""
	_path := ""
	
	if ! _shabang {
		for len (_arguments) == 0 {
			_argument := _arguments[0]
			if strings.HasPrefix (_argument, "--label=") {
				_label = _argument[len ("--label=") :]
				_arguments = _arguments[1:]
			} else if strings.HasPrefix (_argument, "--header=") {
				_header = _argument[len ("--header=") :]
				_arguments = _arguments[1:]
			} else if strings.HasPrefix (_argument, "--path=") {
				_path = _argument[len ("--path=") :]
				_arguments = _arguments[1:]
			} else if _argument == "--" {
				_arguments = _arguments[1:]
				break
			} else if strings.HasPrefix (_argument, "--") {
				return errorf (0xc2b4c6d5, "invalid arguments:  unknown argument `%s`", _argument)
			} else {
				return errorf (0xd1e079a1, "invalid arguments:  unknown argument `%s`, expected `--`", _argument)
			}
		}
	}
	
	if _label == "" {
		_label = "<scriptlet>"
	}
	
	if _path == "" {
		if (len (_arguments) == 0) {
			return errorf (0xc146acb9, "invalid arguments:  expected path")
		}
		_path = _arguments[0]
		_arguments = _arguments[1:]
	}
	
	var _fingerprint string
	var _body string
	if _fingerprint_0, _body_0, _error := loadFromFile (_path); _error == nil {
		if utf8.Valid (_body_0) {
			_fingerprint = _fingerprint_0
			_body = string (_body_0)
		} else {
			return errorf (0x03c1bbdb, "invalid scriptlet:  invalid UTF-8 body")
		}
	} else {
		return _error
	}
	
	var _cacheRoot string
	if _path, _ok := _environment["ZRUN_CACHE"]; _ok {
		_cacheRoot = _path
	} else {
		if _path, _error := resolveCache (); _error == nil {
			_cacheRoot = _path
		} else {
			return _error
		}
	}
	
	_executablePaths := make ([]string, 0, 32)
	if _path, _ok := _environment["PATH"]; _ok {
		for _, _path := range filepath.SplitList (_path) {
			if _path != "" {
				_executablePaths = append (_executablePaths, _path)
			}
		}
	}
	
	var _bodyOffset uint
	
	if _shabang {
		if strings.HasPrefix (_body, "#!") {
			_shabangLimit := strings.IndexByte (_body, '\n')
			if _shabangLimit < 0 {
				return errorf (0xa9c82426, "invalid scriptlet:  missing `#!...\\n`")
			}
			_body = _body[_shabangLimit + 1 :]
			_bodyOffset += 1
		} else {
			return errorf (0xbd838bd0, "invalid scriptlet:  missing `#!`")
		}
	}
	
	_body, _bodyOffsetParse, _interpreter, _interpreterExecutable, _interpreterArguments, _interpreterArgumentsExtraDash, _interpreterArgumentsExtraAllowed, _interpreterEnvironment, _errorParse := parseInterpreter_0 (_label, _body, _header, "")
	if _errorParse != nil {
		return _errorParse
	}
	_bodyOffset += _bodyOffsetParse
	
	_command, _descriptors, _errorPrepare := prepareExecution_0 (
			
			"",
			"",
			
			_interpreter,
			_interpreterExecutable,
			_interpreterArguments,
			_interpreterArgumentsExtraDash,
			_interpreterArgumentsExtraAllowed,
			_interpreterEnvironment,
			
			_fingerprint,
			_label,
			_body,
			
			_path,
			_bodyOffset,
			0,
			
			nil,
			nil,
			
			_selfExecutable,
			_arguments,
			_environment,
			
			"",
			_executablePaths,
			_cacheRoot,
		)
	if _errorPrepare != nil {
		return _errorPrepare
	}
	
	_errorExecute := executeScriptlet_0 (_label, _command, _descriptors, false)
	if _errorExecute != nil {
		return _errorExecute
	}
	
	panic (0xf5fd6b17)
}

