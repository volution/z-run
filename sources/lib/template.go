
package zrun


import "encoding/base64"
import "encoding/hex"
import "encoding/json"
import "io"
import "os"
import "strings"
import "text/template"




func templateMain () (*Error) {
	
	if len (os.Args) < 2 {
		return errorf (0x47c3f9f1, "invalid arguments")
	}
	
	var _source string
	if _stream, _error := os.Open (os.Args[1]); _error == nil {
		defer _stream.Close ()
		_buffer := strings.Builder {}
		if _, _error := io.Copy (&_buffer, _stream); _error == nil {
			_source = _buffer.String ()
		} else {
			return errorw (0x693128fe, _error)
		}
	}
	
	_functions := templateFunctions ()
	
	_template := template.New ("z-run")
	_template.Funcs (_functions)
	if _, _error := _template.Parse (_source); _error != nil {
		return errorw (0xfd33768b, _error)
	}
	
	_data := map[string]interface{} {
			"arguments" : os.Args[2:],
		}
	
	if _error := _template.Execute (os.Stdout, _data); _error != nil {
		return errorw (0x23ce8919, _error)
	}
	
	os.Exit (0)
	panic (0x8e448279)
}




func executeTemplate (_library LibraryStore, _scriptlet *Scriptlet, _context *Context, _output io.Writer) (*Error) {
	
	if _scriptlet.Interpreter != "<template>" {
		return errorf (0xa18a5ca9, "invalid interpreter")
	}
	
	_source := _scriptlet.Body
	
	_functions := templateFunctions ()
	_functions["ZRUN"] = func (_scriptlet string, _arguments ... string) (string, error) {
			return templateFuncZrun (_library, _context, _scriptlet, _arguments)
		}
	
	_template := template.New ("z-run")
	_template.Funcs (_functions)
	if _, _error := _template.Parse (_source); _error != nil {
		return errorw (0xad3804cc, _error)
	}
	
	_data := map[string]interface{} {
			"arguments" : _context.cleanArguments,
			"environment" : _context.cleanEnvironment,
			"workspace" : _context.workspace,
			"terminal" : _context.terminal,
			"executable" : _context.selfExecutable,
			"library" : _library,
		}
	
	if _error := _template.Execute (_output, _data); _error != nil {
		return errorw (0x0d6d4b96, _error)
	}
	
	return nil
}




func templateFuncZrun (_library LibraryStore, _context *Context, _scriptletLabel string, _arguments []string) (string, error) {
	if strings.HasPrefix (_scriptletLabel, ":: ") {
		_scriptletLabel = _scriptletLabel[3:]
	}
	if _scriptlet, _error := _library.ResolveFullByLabel (_scriptletLabel); _error == nil {
		if _scriptlet != nil {
			if _, _output, _error := loadFromScriptlet (_library.Url (), "", _scriptlet, _context); _error == nil {
				return string (_output), nil
			} else {
				return "", _error.ToError ()
			}
		} else {
			return "", errorf (0x944c3172, "unknown scriptlet `%s`", _scriptletLabel) .ToError ()
		}
	} else {
		return "", _error.ToError ()
	}
}




func templateFunctions () (map[string]interface{}) {
	return map[string]interface{} {
			
			"json_encode" : func (_input interface{}) (string, error) {
					_output, _error := json.Marshal (_input)
					return string (_output), _error
				},
			"json_decode" : func (_input string) (interface{}, error) {
					var _output interface{}
					_error := json.Unmarshal ([]byte (_input), &_output)
					return _output, _error
				},
			
			"hex_encode" : func (_input string) (string) {
					return hex.EncodeToString ([]byte (_input))
				},
			"hex_decode" : func (_input string) (string, error) {
					_output, _error := hex.DecodeString (_input)
					return string (_output), _error
				},
			
			"base64_encode" : func (_input string) (string) {
					return base64.StdEncoding.EncodeToString ([]byte (_input))
				},
			"base64_decode" : func (_input string) (string, error) {
					_output, _error := base64.StdEncoding.DecodeString (_input)
					return string (_output), _error
				},
			
			"split" : func (_separator string, _input string) ([]string) {
					return strings.Split (_input, _separator)
				},
			"join" : func (_separator string, _input []string) (string) {
					return strings.Join (_input, _separator)
				},
			
		}
}
