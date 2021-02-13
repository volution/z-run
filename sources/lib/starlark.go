

package zrun


import "os"

import "go.starlark.net/starlark"




func starlarkMain (_selfExecutable string, _arguments []string, _environment map[string]string) (*Error) {
	
	if len (_arguments) < 1 {
		return errorf (0xd776bfb0, "invalid arguments")
	}
	
	var _sourcePath = _arguments[0]
	_arguments = _arguments[1:]
	
	var _sourceBody string
	if _, _data, _error := loadFromFile (_sourcePath); _error == nil {
		_sourceBody = string (_data)
	} else {
		return _error
	}
	
	_error := executeStarlark_0 (
			_sourceBody,
			_arguments,
			_environment,
			_selfExecutable,
			"",
			"",
			"",
			nil,
		)
	if _error != nil {
		return _error
	}
	
	os.Exit (0)
	panic (0x340179ee)
}




func executeStarlark (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (*Error) {
	
	if _scriptlet.Interpreter != "<starlark>" {
		return errorf (0x8044a656, "invalid interpreter")
	}
	
	_libraryUrl := _library.Url ()
	_libraryFingerprint := ""
	if _libraryFingerprint_0, _error := _library.Fingerprint (); _error == nil {
		_libraryFingerprint = _libraryFingerprint_0
	} else {
		return _error
	}
	
	_extraFunctions := make (map[string]interface{}, 16)
	
	return executeStarlark_0 (
			_scriptlet.Body,
			_context.cleanArguments,
			_context.cleanEnvironment,
			_context.selfExecutable,
			_context.workspace,
			_libraryUrl,
			_libraryFingerprint,
			_extraFunctions,
		)
}




func executeStarlark_0 (
			_source string,
			_arguments []string,
			_environment map[string]string,
			_selfExecutable string,
			_workspace string,
			_libraryUrl string,
			_libraryFingerprint string,
			_extraFunctions map[string]interface{},
		) (*Error) {
	
	_thread := & starlark.Thread {
			Name : "scriptlet",
			Print : func (_thread *starlark.Thread, _message string) () {
					logf ('>', 0x98499c0a, "%s", _message)
				},
			Load : nil,
		}
	
	_builtins := starlark.StringDict {}
	
	if _, _error := starlark.ExecFile (_thread, "<scriptlet>", _source, _builtins); _error != nil {
		return errorw (0x4027d940, _error)
	}
	
	return nil
}

