

package zrun


import "encoding/base64"
import "encoding/json"
import "log"
import "os"
import "path"
import "path/filepath"
import "strings"
import "syscall"
import "unicode"
import "unicode/utf8"




type Context struct {
	selfExecutable string
	selfArgument0 string
	selfArguments []string
	selfEnvironment map[string]string
	cleanArguments []string
	cleanEnvironment map[string]string
	workspace string
	cacheRoot string
	cacheEnabled bool
	terminal string
	top bool
}


type SshContext struct {
	target string
	launcher string
	delegate string
	workspace string
	cache string
	terminal string
	libraryLocalSocket string
	libraryRemoteSocket string
	token string
	exportEnvironment []string
}

type InvokeContext struct {
	Library string `json:"library"`
	Scriptlet string `json:"scriptlet"`
	Arguments []string `json:"arguments"`
	Environment map[string]string `json:"environment"`
	Workspace string `json:"workspace"`
	Cache string `json:"cache"`
	Terminal string `json:"terminal"`
}




func main_0 (_executable string, _argument0 string, _arguments []string, _environment map[string]string, _commandOverride string, _scriptletOverride string) (*Error) {
	
	var _command string = _commandOverride
	var _scriptlet string = _scriptletOverride
	
	var _librarySourcePath string
	var _libraryCachePath string
	var _libraryLookupPaths []string = make ([]string, 0, 128)
	
	var _cleanArguments []string
	var _cleanEnvironment map[string]string = make (map[string]string, len (_environment))
	
	var _workspace string
	var _cacheRoot string
	var _terminal string
	var _execMode bool
	var _invokeMode bool
	var _invokeContextEncoded string
	var _sshMode bool
	var _sshContext *SshContext
	var _top bool
	
	_top = true
	
	for _name, _value := range _environment {
		
		var _nameCanonical string
		{
			_nameCanonical = strings.ToUpper (_name)
			_nameCanonical = strings.ReplaceAll (_nameCanonical, "-", "_")
			for {
				_nameCanonical_0 := strings.ReplaceAll (_nameCanonical, "__", "_")
				if _nameCanonical != _nameCanonical_0 {
					_nameCanonical = _nameCanonical_0
				} else {
					break
				}
			}
			_nameCanonical = strings.Replace (_nameCanonical, "Z_RUN", "ZRUN", 1)
		}
		
		if strings.HasPrefix (_nameCanonical, "ZRUN_") || strings.HasPrefix (_nameCanonical, "_ZRUN_") {
			
			if _name != _nameCanonical {
				logf ('w', 0x9bc8b3da, "environment variable does not have canonical name;  expected `%s`, encountered `%s`!", _nameCanonical, _name)
			}
			
			switch _nameCanonical {
				case "ZRUN_LIBRARY_SOURCE" :
					_librarySourcePath = _value
				case "ZRUN_LIBRARY_CACHE" :
					_libraryCachePath = _value
				case "ZRUN_WORKSPACE" :
					_workspace = _value
				case "ZRUN_EXECUTABLE" :
					if _executable != _value {
						logf ('w', 0xfb1f0645, "environment variable mismatched:  `%s`;  expected `%s`, encountered `%s`!", _nameCanonical, _executable, _value)
					}
					// FIXME:  Find a better way to handle this!
					_top = false
				case "ZRUN_CACHE" :
					_cacheRoot = _value
				case "ZRUN_TERM" :
					_terminal = _value
				case "ZRUN_FINGERPRINT" :
					// FIXME:  Validate that this value actually matches given library.
				default :
					logf ('w', 0xafe247b0, "environment variable unknown:  `%s` with value `%s`", _nameCanonical, _value)
			}
			
		} else if strings.HasPrefix (_nameCanonical, "X_RUN_") || strings.HasPrefix (_nameCanonical, "_X_RUN_") {
			
			if _name != _nameCanonical {
				logf ('w', 0x37850eb3, "environment variable does not have canonical name;  expected `%s`, encountered `%s`!", _nameCanonical, _name)
			}
			
			logf ('w', 0xdf61b057, "environment variable unknown:  `%s` with value `%s`", _nameCanonical, _value)
			
		} else if _name == "_" {
			
			// NOTE:  For some reason `bash` always exports `_` containing the executable, thus we ignore it!
			
		} else {
			_cleanEnvironment[_name] = _value
		}
	}
	
	for _index, _argument := range _arguments {
		
		if _scriptlet != "" {
			_cleanArguments = _arguments[_index:]
			break
			
		} else if _execMode {
			_librarySourcePath = _arguments[_index]
			_cleanArguments = _arguments[_index + 1:]
			break
			
		} else if _invokeMode {
			_invokeContextEncoded = _arguments[_index]
			_cleanArguments = _arguments[_index + 1:]
			break
			
		} else if _argument == "--" {
			_cleanArguments = _arguments[_index + 1:]
			break
			
		} else if strings.HasPrefix (_argument, "-") {
			
			if _command != "" {
				return errorf (0xae04b5ff, "unexpected argument `%s`", _argument)
			}
			
			if _argument == "--exec" {
				if _index == 0 {
					_execMode = true
				} else {
					return errorf (0x12cdad05, "unexpected argument `--exec` (only first)")
				}
				
			} else if _argument == "--invoke" {
				if _index == 0 {
					_invokeMode = true
				} else {
					return errorf (0x03da9932, "unexpected argument `--invoke` (only first)")
				}
				
			} else if _argument == "--ssh" {
				if _index == 0 {
					_sshMode = true
					_sshContext = & SshContext {}
				} else {
					return errorf (0x6a4ba0b6, "unexpected argument `--ssh` (only first)")
				}
				
			} else if strings.HasPrefix (_argument, "--library-source=") {
				_librarySourcePath = _argument[len ("--library-source="):]
				_libraryCachePath = ""
				
			} else if strings.HasPrefix (_argument, "--library-cache=") {
				_libraryCachePath = _argument[len ("--library-cache="):]
				
			} else if strings.HasPrefix (_argument, "--workspace=") {
				_workspace = _argument[len ("--workspace="):]
				
			} else if strings.HasPrefix (_argument, "--ssh-") {
				if !_sshMode {
					return errorf (0x0e1cdc68, "unexpected argument `%s` (only with `--ssh`)", _argument)
				}
				if strings.HasPrefix (_argument, "--ssh-target=") {
					_target := _argument[len ("--ssh-target="):]
					_sshContext.target = _target
				} else if strings.HasPrefix (_argument, "--ssh-export=") {
					_name := _argument[len ("--ssh-export="):]
					_sshContext.exportEnvironment = append (_sshContext.exportEnvironment, _name)
				} else if strings.HasPrefix (_argument, "--ssh-workspace=") {
					_workspace := _argument[len ("--ssh-workspace="):]
					_sshContext.workspace = _workspace
				} else {
					return errorf (0x33555ffb, "invalid argument `%s`", _argument)
				}
				
			} else {
				return errorf (0x33555ffb, "invalid argument `%s`", _argument)
			}
			
		} else if strings.HasPrefix (_argument, "::") {
			_scriptlet = _argument
			
		} else if _sshMode {
			return errorf (0x7af6b31f, "unexpected argument `%s` (for `--ssh`)", _argument)
			
		} else {
			if _command == "" {
				switch _argument {
					
					case "execute-scriptlet", "execute" :
						_command = "execute-scriptlet"
					
					case "execute-scriptlet-ssh", "execute-ssh", "ssh" :
						_command = "execute-scriptlet-ssh"
						_sshContext = & SshContext {}
					
					case "export-scriptlet-body", "export-body" :
						_command = "export-scriptlet-body"
					
					case "select-execute-scriptlet", "select-execute" :
						_command = "select-execute-scriptlet"
					
					case "select-export-scriptlet-label", "select-label", "select" :
						_command = "select-export-scriptlet-label"
					
					case "select-export-scriptlet-body", "select-body" :
						_command = "select-export-scriptlet-body"
					
					case "select-export-scriptlet-label-and-body" :
						_command = "select-export-scriptlet-label-and-body"
					
					case "export-scriptlet-labels", "export-labels", "list" :
						_command = "export-scriptlet-labels"
					
					case "parse-library" :
						_command = "parse-library"
					
					case "export-library-json" :
						_command = "export-library-json"
					
					case "export-library-cdb" :
						_command = "export-library-cdb"
					
					case "export-library-rpc" :
						_command = "export-library-rpc"
					
					case "export-fingerprint" :
						_command = "export-fingerprint"
				}
				
			} else if (_command == "execute-scriptlet-ssh") && (_sshContext.target == "") {
				_sshContext.target = _argument
				
			} else {
				return errorf (0x6a6a6cef, "unexpected argument `%s`", _argument)
			}
		}
	}
	
	if _execMode {
		if _librarySourcePath == "" {
			return errorf (0xe13c6051, "invalid arguments:  expected source path")
		}
		if _scriptlet != "" {
			return errorf (0xb2b83ca4, "invalid arguments:  unexpected scriptlet")
		}
		if len (_cleanArguments) > 0 {
			_scriptlet = _cleanArguments[0]
			_cleanArguments = _cleanArguments[1:]
			if ! strings.HasPrefix (_scriptlet, "::") {
				_scriptlet = ":: " + _scriptlet
			}
		}
		_workspace = path.Dir (_librarySourcePath)
	}
	
	if _invokeMode {
		if _command != "" {
			return errorf (0x996e4c2b, "invalid arguments:  unexpected command")
		}
		if _scriptlet != "" {
			return errorf (0x3344b708, "invalid arguments:  unexpected scriptlet")
		}
		if len (_cleanArguments) != 0 {
			return errorf (0x71d92eec, "invalid arguments:  unexpected arguments")
		}
		if (_librarySourcePath != "") || (_libraryCachePath != "") || (len (_libraryLookupPaths) != 0) {
			return errorf (0x70d72d6d, "invalid arguments:  unexpected library source, cache or lookup")
		}
		if _workspace != "" {
			return errorf (0x70d37d49, "invalid arguments:  unexpected workspace")
		}
		if _cacheRoot != "" {
			return errorf (0xb8216677, "invalid arguments:  unexpected cache")
		}
		if _terminal != "" {
			return errorf (0xc99ea71b, "invalid arguments:  unexpected terminal")
		}
		if _invokeContextEncoded == "" {
			return errorf (0xd08f059a, "invalid arguments:  expected invoke context")
		}
		var _context InvokeContext
		if _data, _error := base64.RawURLEncoding.DecodeString (_invokeContextEncoded); _error == nil {
			if _error := json.Unmarshal (_data, &_context); _error == nil {
				// NOP
			} else {
				return errorw (0x29dfbd1e, _error)
			}
		} else {
			return errorw (0x7c9f1ada, _error)
		}
		_scriptlet = _context.Scriptlet
		if ! strings.HasPrefix (_scriptlet, "::") {
			_scriptlet = ":: " + _scriptlet
		}
		_cleanArguments = _context.Arguments
		for _name, _value := range _context.Environment {
			_cleanEnvironment[_name] = _value
		}
		_libraryCachePath = _context.Library
		_workspace = _context.Workspace
		_cacheRoot = _context.Cache
		_terminal = _context.Terminal
		_top = false
	}
	
	if _sshMode {
		if _command != "" {
			return errorf (0x8937413a, "invalid arguments:  unexpected command")
		}
		_command = "execute-scriptlet-ssh"
	}
	
	if _scriptlet != "" {
		if strings.HasPrefix (_scriptlet, ":: ") {
			_scriptlet = _scriptlet[3:]
		} else {
			return errorf (0x72ad17f7, "invalid scriptlet label `%s`", _scriptlet)
		}
	}
	
	if (_command == "") {
		if ((_scriptlet == "") || _top) && (len (_cleanArguments) == 0) {
			_command = "select-execute-scriptlet"
		} else {
			_command = "execute-scriptlet"
		}
	}
	
	_cacheEnabled := true
	if _command == "parse-library" {
		_cacheEnabled = false
	}
	if _cacheEnabled {
		if _cacheRoot == "" {
			if _cacheRoot_0, _error := resolveCache (); _error == nil {
				_cacheRoot = _cacheRoot_0
			} else {
				return _error
			}
		}
	} else {
		if _libraryCachePath != "" {
			logf ('w', 0xdb80c4de, "cached library path specified, but caching is disabled;  ignoring cached path!")
			_libraryCachePath = ""
		}
		_cacheRoot = ""
	}
	
	if _workspace == "" {
		if _path, _error := os.Getwd (); _error == nil {
			_workspace = _path
		} else {
			return errorw (0x69daa060, _error)
		}
		var _insideVcs bool
		for _, _subfolder := range resolveWorkspaceSubfolders {
			if _, _error := os.Lstat (path.Join (_workspace, _subfolder)); _error == nil {
				_insideVcs = true
				break
			} else if os.IsNotExist (_error) {
				// NOP
			} else {
				return errorw (0x43bc330c, _error)
			}
		}
		if !_insideVcs {
			if _home, _error := os.UserHomeDir (); _error == nil {
				_libraryLookupPaths = append (_libraryLookupPaths, _home)
			} else {
				return errorw (0x3884b718, _error)
			}
		}
	}
	if _path, _error := filepath.Abs (_workspace); _error == nil {
		_workspace = _path
	} else {
		return errorw (0x9f5c1d2a, _error)
	}
	
	if _terminal == "" {
		_terminal, _ = _cleanEnvironment["TERM"]
	}
	if _terminal == "dumb" {
		_terminal = ""
	}
	
	if _libraryCachePath != "" {
		if _librarySourcePath != "" {
			logf ('w', 0x1fe0b572, "cached library path specified, but also source path specified;  ignoring cached path!")
			_libraryCachePath = ""
		}
	}
	
	_context := & Context {
			selfExecutable : _executable,
			selfArguments : _arguments,
			selfEnvironment : _environment,
			cleanArguments : _cleanArguments,
			cleanEnvironment : _cleanEnvironment,
			workspace : _workspace,
			cacheRoot : _cacheRoot,
			cacheEnabled : _cacheEnabled,
			terminal : _terminal,
		}
	
	var _library LibraryStore
	if _libraryCachePath != "" {
		if strings.HasPrefix (_libraryCachePath, "unix:") || strings.HasPrefix (_libraryCachePath, "tcp:") {
			if _library_0, _error := NewLibraryRpcClient (_libraryCachePath); _error == nil {
				_library = _library_0
			} else {
				return _error
			}
		} else {
			if _library_0, _error := resolveLibraryCached (_libraryCachePath); _error == nil {
				_library = _library_0
			} else {
				return _error
			}
		}
	} else {
		if _library_0, _error := resolveLibrary (_librarySourcePath, _context, _libraryLookupPaths, _execMode); _error == nil {
			_library = _library_0
		} else {
			return _error
		}
	}
	
	switch _command {
		
		
		case "execute-scriptlet" :
			if _scriptlet == "" {
				return errorf (0x39718e70, "execute:  expected scriptlet")
			}
			return doHandleWithLabel (_library, _scriptlet, doHandleExecuteScriptlet, _context)
		
		case "execute-scriptlet-ssh" :
			if _scriptlet == "" {
				return errorf (0xcc3c2ea6, "execute-ssh:  expected scriptlet")
			}
			_handler := func (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (bool, *Error) {
					return doHandleExecuteScriptletSsh (_library, _scriptlet, _sshContext, _context)
				}
			return doHandleWithLabel (_library, _scriptlet, _handler, _context)
		
		case "export-scriptlet-body" :
			if _scriptlet == "" {
				return errorf (0xf24640a2, "export:  expected scriptlet")
			}
			if len (_cleanArguments) != 0 {
				return errorf (0xcf8db3c0, "export:  unexpected arguments")
			}
			_handler := func (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (bool, *Error) {
					return doHandleExportScriptletBody (_library, _scriptlet, os.Stdout, _context)
				}
			return doHandleWithLabel (_library, _scriptlet, _handler, _context)
		
		
		case "select-execute-scriptlet" :
			if len (_cleanArguments) != 0 {
				return errorf (0x203e410a, "select:  unexpected arguments")
			}
			if _scriptlet != "" {
				return doSelectHandleWithLabel (_library, _scriptlet, doHandleExecuteScriptlet, _context)
			} else {
				return doSelectHandle (_library, doHandleExecuteScriptlet, _context)
			}
		
		case "select-export-scriptlet-label" :
			if len (_cleanArguments) != 0 {
				return errorf (0x2d19b1bc, "select:  unexpected arguments")
			}
			_handler := func (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (bool, *Error) {
					return doHandleExportScriptletLabel (_library, _scriptlet, os.Stdout, _context)
				}
			if _scriptlet != "" {
				return doSelectHandleWithLabel (_library, _scriptlet, _handler, _context)
			} else {
				return doSelectHandle (_library, _handler, _context)
			}
		
		case "select-export-scriptlet-body" :
			if len (_cleanArguments) != 0 {
				return errorf (0x5f573713, "select:  unexpected arguments")
			}
			_handler := func (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (bool, *Error) {
					return doHandleExportScriptletBody (_library, _scriptlet, os.Stdout, _context)
				}
			if _scriptlet != "" {
				return doSelectHandleWithLabel (_library, _scriptlet, _handler, _context)
			} else {
				return doSelectHandle (_library, _handler, _context)
			}
		
		case "select-export-scriptlet-label-and-body" :
			if len (_cleanArguments) != 0 {
				return errorf (0xe4f7e6f5, "export:  unexpected arguments")
			}
			_handler := func (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (bool, *Error) {
					return doHandleExportScriptletLegacy (_library, _scriptlet, os.Stdout, _context)
				}
			if _scriptlet != "" {
				return doSelectHandleWithLabel (_library, _scriptlet, _handler, _context)
			} else {
				return doSelectHandle (_library, _handler, _context)
			}
		
		
		case "export-scriptlet-labels" :
			if (_scriptlet != "") || (len (_cleanArguments) != 0) {
				return errorf (0xf7b9c7f3, "export:  unexpected scriptlet or arguments")
			}
			return doExportScriptletLabels (_library, os.Stdout, _context)
		
		
		case "parse-library" :
			if (_scriptlet != "") || (len (_cleanArguments) != 0) {
				return errorf (0x400ec122, "export:  unexpected scriptlet or arguments")
			}
			return doExportLibraryJson (_library, os.Stdout, _context)
		
		case "export-library-json" :
			if (_scriptlet != "") || (len (_cleanArguments) != 0) {
				return errorf (0xdd0752b8, "export:  unexpected scriptlet or arguments")
			}
			return doExportLibraryStore (_library, NewJsonStreamStoreOutput (os.Stdout, nil), _context)
		
		case "export-library-cdb" :
			if _scriptlet != "" {
				return errorf (0x492ac50e, "export:  unexpected scriptlet")
			}
			if len (_cleanArguments) != 1 {
				return errorf (0xf76f4459, "export:  expected database path")
			}
			return doExportLibraryCdb (_library, _cleanArguments[0], _context)
		
		case "export-library-rpc" :
			if _scriptlet != "" {
				return errorf (0x04d71684, "export:  unexpected scriptlet")
			}
			if len (_cleanArguments) != 1 {
				return errorf (0xe7886d74, "export:  expected RPC url")
			}
			return doExportLibraryRpc (_library, _cleanArguments[0], _context)
		
		case "export-fingerprint" :
			if _scriptlet != "" {
				return errorf (0x3483242d, "export:  unexpected scriptlet")
			}
			if len (_cleanArguments) != 0 {
				return errorf (0x2a741648, "export:  unexpected arguments")
			}
			if _fingerprint, _error := _library.Fingerprint (); _error == nil {
				log.Println (_fingerprint)
				return nil
			} else {
				return _error
			}
		
		case "" :
			return errorf (0x5d2a4326, "expected command")
		
		default :
			return errorf (0x66cf8700, "unexpected command `%s`", _command)
	}
}




func Main () () {
	
	log.SetFlags (0)
	
	var _executable string
	if _executable_0, _error := os.Executable (); _error == nil {
		_executable = _executable_0
	} else {
		panic (abortError (errorw (0x75f2db30, _error)))
	}
	
	_argument0 := os.Args[0]
	
	if strings.HasPrefix (_argument0, "[z-run:menu] ") {
		_argument0 = "[z-run:menu]"
	} else if strings.HasPrefix (_argument0, "[z-run:select] ") {
		_argument0 = "[z-run:select]"
	} else if strings.HasPrefix (_argument0, "[z-run:template-raw] ") {
		_argument0 = "[z-run:template-raw]"
	}
	
	switch _argument0 {
		
		case "[z-run:menu]" :
			if _error := menuMain (); _error != nil {
				panic (abortError (_error))
			} else {
				panic (0x6b21e0ab)
			}
		
		case "[z-run:select]" :
			if _error := fzfMain (); _error != nil {
				panic (abortError (_error))
			} else {
				panic (0x2346ca3f)
			}
		
		case "[z-run:template-raw]" :
			if _error := templateMain (); _error != nil {
				panic (abortError (_error))
			} else {
				panic (0x32241835)
			}
		
		case "[z-run]" :
			// NOP
		
		default :
			_arguments := os.Args
			_arguments[0] = "[z-run]"
			_environment := os.Environ ()
			if _error := syscall.Exec (_executable, _arguments, _environment); _error != nil {
				panic (abortError (errorw (0x05bd220d, _error)))
			} else {
				panic (0xe13aab5f)
			}
	}
	
	_arguments := append ([]string (nil), os.Args[1:] ...)
	
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
	
	if _error := main_0 (_executable, _argument0, _arguments, _environment, "", ""); _error == nil {
		os.Exit (0)
		panic (0xe0e1c1a1)
	} else {
		panic (abortError (_error))
	}
}

