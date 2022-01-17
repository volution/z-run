

package zrun


import "encoding/base64"
import "encoding/json"
import "fmt"
import "os"
import "path"
import "path/filepath"
import "strings"

import . "github.com/cipriancraciun/z-run/lib/library"
import . "github.com/cipriancraciun/z-run/lib/store"
import . "github.com/cipriancraciun/z-run/lib/common"
import . "github.com/cipriancraciun/z-run/embedded"




type Context struct {
	selfExecutable string
	selfArgument0 string
	selfArguments []string
	selfEnvironment map[string]string
	cleanArguments []string
	cleanEnvironment map[string]string
	executablePaths []string
	terminal string
	workspace string
	cacheRoot string
	cacheEnabled bool
	top bool
	preMainReExecute func (string) (*Error)
}


type SshContext struct {
	target string
	launcher string
	delegate string
	exportEnvironment []string
	executablePaths []string
	terminal string
	workspace string
	cache string
	libraryLocalSocket string
	libraryRemoteSocket string
	token string
}

type InvokeContext struct {
	Version string `json:"version,omitempty"`
	Library string `json:"library,omitempty"`
	Scriptlet string `json:"scriptlet,omitempty"`
	Arguments []string `json:"arguments,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	ExecutablePaths []string `json:"executable-paths,omitempty"`
	Terminal string `json:"terminal,omitempty"`
	Workspace string `json:"workspace,omitempty"`
	CacheRoot string `json:"cache-root,omitempty"`
}




func RunMain (_executable string, _argument0 string, _arguments []string, _environment map[string]string, _commandOverride string, _scriptletOverride string, _preMainReExecute func (string) (*Error)) (*Error) {
	
	var _command string = _commandOverride
	var _scriptlet string = _scriptletOverride
	
	var _librarySourcePath string
	var _libraryCacheUrl string
	var _libraryLookupPaths []string = make ([]string, 0, 128)
	
	var _cleanArguments []string
	var _cleanEnvironment map[string]string = make (map[string]string, len (_environment))
	
	var _executablePaths []string = make ([]string, 0, 32)
	var _terminal string
	
	var _workspace string
	var _cacheRoot string
	
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
				Logf ('w', 0x9bc8b3da, "environment variable does not have canonical name;  expected `%s`, encountered `%s`!", _nameCanonical, _name)
			}
			
			switch _nameCanonical {
				
				case "ZRUN_LIBRARY_SOURCE" :
					_librarySourcePath = _value
				case "ZRUN_LIBRARY_URL" :
					_libraryCacheUrl = _value
				case "ZRUN_LIBRARY_IDENTIFIER", "ZRUN_LIBRARY_FINGERPRINT" :
					// FIXME:  Validate that this value actually matches given library.
					_top = false
				
				case "ZRUN_WORKSPACE" :
					_workspace = _value
				
				case "ZRUN_EXECUTABLE" :
					if _executable != _value {
						Logf ('w', 0xfb1f0645, "environment variable mismatched:  `%s`;  expected `%s`, encountered `%s`!", _nameCanonical, _executable, _value)
					}
				
				case "ZRUN_CACHE" :
					_cacheRoot = _value
				
				default :
					Logf ('w', 0xafe247b0, "environment variable unknown:  `%s` with value `%s`", _nameCanonical, _value)
			}
			
		} else if strings.HasPrefix (_nameCanonical, "X_RUN_") || strings.HasPrefix (_nameCanonical, "_X_RUN_") {
			
			if _name != _nameCanonical {
				Logf ('w', 0x37850eb3, "environment variable does not have canonical name;  expected `%s`, encountered `%s`!", _nameCanonical, _name)
			}
			
			Logf ('w', 0xdf61b057, "environment variable unknown:  `%s` with value `%s`", _nameCanonical, _value)
			
		} else if _name == "_" {
			
			// NOTE:  For some reason `bash` always exports `_` containing the executable, thus we ignore it!
			
		} else if _name == "PATH" {
			
			for _, _path := range filepath.SplitList (_value) {
				if _path != "" {
					_executablePaths = append (_executablePaths, _path)
				}
			}
			
		} else if _name == "TERM" {
			
			_cleanEnvironment[_name] = _value
			
			if _terminal == "" {
				_terminal = _value
			}
			
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
				return Errorf (0xae04b5ff, "unexpected argument `%s`", _argument)
			}
			
			if _argument == "--exec" {
				if _index == 0 {
					_execMode = true
					_librarySourcePath = ""
					_libraryCacheUrl = ""
					_workspace = ""
					_top = true
				} else {
					return Errorf (0x12cdad05, "unexpected argument `--exec` (only first)")
				}
				
			} else if _argument == "--untainted" {
				if _index == 0 {
					_librarySourcePath = ""
					_libraryCacheUrl = ""
					_workspace = ""
					_top = true
				} else {
					return Errorf (0x7ce36db9, "unexpected argument `--untainted` (only first)")
				}
				
			} else if _argument == "--invoke" {
				if _index == 0 {
					_invokeMode = true
				} else {
					return Errorf (0x03da9932, "unexpected argument `--invoke` (only first)")
				}
				
			} else if _argument == "--ssh" {
				if _index == 0 {
					_sshMode = true
					_sshContext = & SshContext {}
				} else {
					return Errorf (0x6a4ba0b6, "unexpected argument `--ssh` (only first)")
				}
				
			} else if strings.HasPrefix (_argument, "--library-source=") {
				_librarySourcePath = _argument[len ("--library-source="):]
				_libraryCacheUrl = ""
				
			} else if strings.HasPrefix (_argument, "--library-url=") {
				_libraryCacheUrl = _argument[len ("--library-url="):]
				
			} else if strings.HasPrefix (_argument, "--workspace=") {
				_workspace = _argument[len ("--workspace="):]
				
			} else if strings.HasPrefix (_argument, "--ssh-") {
				if !_sshMode {
					return Errorf (0x0e1cdc68, "unexpected argument `%s` (only with `--ssh`)", _argument)
				}
				if strings.HasPrefix (_argument, "--ssh-target=") {
					_target := _argument[len ("--ssh-target="):]
					_sshContext.target = _target
				} else if strings.HasPrefix (_argument, "--ssh-export=") {
					_name := _argument[len ("--ssh-export="):]
					_sshContext.exportEnvironment = append (_sshContext.exportEnvironment, _name)
				} else if strings.HasPrefix (_argument, "--ssh-path=") {
					_paths := filepath.SplitList (_argument[len ("--ssh-path="):])
					for _, _path := range _paths {
						_sshContext.executablePaths = append (_sshContext.executablePaths, _path)
					}
				} else if strings.HasPrefix (_argument, "--ssh-terminal=") {
					_terminal := _argument[len ("--ssh-terminal="):]
					_sshContext.terminal = _terminal
				} else if strings.HasPrefix (_argument, "--ssh-workspace=") {
					_workspace := _argument[len ("--ssh-workspace="):]
					_sshContext.workspace = _workspace
				} else {
					return Errorf (0xeed6fa17, "invalid argument `%s`", _argument)
				}
				
			} else {
				return Errorf (0x33555ffb, "invalid argument `%s`", _argument)
			}
			
		} else if strings.HasPrefix (_argument, "::") {
			_scriptlet = _argument
			
		} else if _sshMode {
			return Errorf (0x7af6b31f, "unexpected argument `%s` (for `--ssh`)", _argument)
			
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
					
					case "select-execute-scriptlet-loop", "select-execute-loop", "loop" :
						_command = "select-execute-scriptlet-loop"
					
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
					
					case "parse-library-without-output" :
						_command = "parse-library-without-output"
					
					case "export-library-json" :
						_command = "export-library-json"
					
					case "export-library-cdb" :
						_command = "export-library-cdb"
					
					case "export-library-rpc" :
						_command = "export-library-rpc"
					
					case "export-library-url" :
						_command = "export-library-url"
					
					case "export-library-identifier" :
						_command = "export-library-identifier"
					
					case "export-library-fingerprint" :
						_command = "export-library-fingerprint"
					
					default :
						return Errorf (0x63de47f8, "unexpected argument `%s`", _argument)
				}
				
			} else if (_command == "execute-scriptlet-ssh") && (_sshContext.target == "") {
				_sshContext.target = _argument
				
			} else if (_command == "export-library-json") || (_command == "export-library-cdb") || (_command == "export-library-rpc") {
				_cleanArguments = _arguments[_index:]
				break
				
			} else {
				return Errorf (0x6a6a6cef, "unexpected argument `%s`", _argument)
			}
		}
	}
	
	if _execMode {
		if _librarySourcePath == "" {
			return Errorf (0xe13c6051, "invalid arguments:  expected source path")
		}
		if _scriptlet != "" {
			return Errorf (0xb2b83ca4, "invalid arguments:  unexpected scriptlet")
		}
		if _path, _error := filepath.EvalSymlinks (_librarySourcePath); _error == nil {
			_librarySourcePath = _path
		} else {
			return Errorw (0x8e43d808, _error)
		}
		if len (_cleanArguments) > 0 {
			_scriptlet = _cleanArguments[0]
			_cleanArguments = _cleanArguments[1:]
			if ! strings.HasPrefix (_scriptlet, "::") {
				_scriptlet = ":: " + _scriptlet
			}
		}
		_workspace = path.Dir (_librarySourcePath)
		for _, _subfolder := range ResolveSourceSubfolders {
			_suffix := string (filepath.Separator) + _subfolder
			if strings.HasSuffix (_workspace, _suffix) {
				_workspace = _workspace[: len (_workspace) - len (_suffix)]
				break
			}
		}
	}
	
	if _invokeMode {
		if _command != "" {
			return Errorf (0x996e4c2b, "invalid arguments:  unexpected command")
		}
		if _scriptlet != "" {
			return Errorf (0x3344b708, "invalid arguments:  unexpected scriptlet")
		}
		if len (_cleanArguments) != 0 {
			return Errorf (0x71d92eec, "invalid arguments:  unexpected arguments")
		}
		if (_librarySourcePath != "") || (_libraryCacheUrl != "") || (len (_libraryLookupPaths) != 0) {
			return Errorf (0x70d72d6d, "invalid arguments:  unexpected library source, cache or lookup")
		}
		if _workspace != "" {
			return Errorf (0x70d37d49, "invalid arguments:  unexpected workspace")
		}
		if _cacheRoot != "" {
			return Errorf (0xb8216677, "invalid arguments:  unexpected cache")
		}
		if _invokeContextEncoded == "" {
			return Errorf (0xd08f059a, "invalid arguments:  expected invoke context")
		}
		var _context InvokeContext
		if _data, _error := base64.RawURLEncoding.DecodeString (_invokeContextEncoded); _error == nil {
			if _error := json.Unmarshal (_data, &_context); _error == nil {
				// NOP
			} else {
				return Errorw (0x29dfbd1e, _error)
			}
		} else {
			return Errorw (0x7c9f1ada, _error)
		}
		if _context.Version != BUILD_VERSION {
			return Errorf (0xfe2f9709, "mismatched version, self `%s`, other `%s`", BUILD_VERSION, _context.Version)
		}
		_scriptlet = _context.Scriptlet
		if ! strings.HasPrefix (_scriptlet, "::") {
			_scriptlet = ":: " + _scriptlet
		}
		_cleanArguments = _context.Arguments
		for _name, _value := range _context.Environment {
			_cleanEnvironment[_name] = _value
		}
		_executablePaths = _context.ExecutablePaths
		_terminal = _context.Terminal
		_libraryCacheUrl = _context.Library
		_workspace = _context.Workspace
		_cacheRoot = _context.CacheRoot
		_top = false
	}
	
	if _sshMode {
		if _command != "" {
			return Errorf (0x8937413a, "invalid arguments:  unexpected command")
		}
		_command = "execute-scriptlet-ssh"
	}
	
	if _scriptlet != "" {
		if strings.HasPrefix (_scriptlet, ":: ") {
			_scriptlet = _scriptlet[3:]
		} else {
			return Errorf (0x72ad17f7, "invalid scriptlet label `%s`", _scriptlet)
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
	if (_command == "parse-library") || (_command == "parse-library-without-output")  {
		_cacheEnabled = false
	}
	if !_cacheEnabled {
		if _libraryCacheUrl != "" {
			Logf ('w', 0xdb80c4de, "library URL specified, but caching is disabled;  ignoring cached path!")
			_libraryCacheUrl = ""
		}
	}
	
	if _cacheRoot == "" {
		if _cacheRoot_0, _error := resolveCache (); _error == nil {
			_cacheRoot = _cacheRoot_0
		} else {
			return _error
		}
	}
	
	if _workspace == "" {
		if _path, _error := os.Getwd (); _error == nil {
			_workspace = _path
		} else {
			return Errorw (0x69daa060, _error)
		}
		var _insideVcs bool
		for _, _subfolder := range ResolveWorkspaceSubfolders {
			if _, _error := os.Lstat (path.Join (_workspace, _subfolder)); _error == nil {
				_insideVcs = true
				break
			} else if os.IsNotExist (_error) {
				// NOP
			} else {
				return Errorw (0x43bc330c, _error)
			}
		}
		if !_insideVcs {
			if _home, _error := os.UserHomeDir (); _error == nil {
				_libraryLookupPaths = append (_libraryLookupPaths, _home)
			} else {
				return Errorw (0x3884b718, _error)
			}
		}
	}
	if _path, _error := filepath.Abs (_workspace); _error == nil {
		_workspace = _path
	} else {
		return Errorw (0x9f5c1d2a, _error)
	}
	
	if _terminal == "dumb" {
		_terminal = ""
	}
	
	if _libraryCacheUrl != "" {
		if _librarySourcePath != "" {
			Logf ('w', 0x1fe0b572, "library URL specified, but also source path specified;  ignoring cached path!")
			_libraryCacheUrl = ""
		}
	}
	
	_context := & Context {
			selfExecutable : _executable,
			selfArguments : _arguments,
			selfEnvironment : _environment,
			cleanArguments : _cleanArguments,
			cleanEnvironment : _cleanEnvironment,
			executablePaths : _executablePaths,
			terminal : _terminal,
			workspace : _workspace,
			cacheRoot : _cacheRoot,
			cacheEnabled : _cacheEnabled,
			preMainReExecute : _preMainReExecute,
		}
	
	var _library LibraryStore
	if _libraryCacheUrl != "" {
		if strings.HasPrefix (_libraryCacheUrl, "unix:") || strings.HasPrefix (_libraryCacheUrl, "tcp:") {
			if _library_0, _error := NewLibraryRpcClient (_libraryCacheUrl); _error == nil {
				_library = _library_0
			} else {
				return _error
			}
		} else {
			if _library_0, _error := resolveLibraryCached (_libraryCacheUrl); _error == nil {
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
	
	if _top {
		if _libraryContext, _error := _library.SelectLibraryContext (); _error == nil {
			if (_libraryContext.SelfExecutable != "") && (_libraryContext.SelfExecutable != _context.selfExecutable) {
				_context.preMainReExecute (_libraryContext.SelfExecutable)
			}
		} else {
			return _error
		}
	}
	
//	Logf ('d', 0x66e8b16d, "%v `%s` `%s` %d %v", _top, _command, _scriptlet, len (_cleanArguments), _cleanArguments)
	
	switch _command {
		
		
		case "execute-scriptlet" :
			if _scriptlet == "" {
				return Errorf (0x39718e70, "execute:  expected scriptlet")
			}
			return doHandleWithLabel (_library, _scriptlet, doHandleExecuteScriptlet, _context)
		
		case "execute-scriptlet-ssh" :
			if _scriptlet == "" {
				return Errorf (0xcc3c2ea6, "execute-ssh:  expected scriptlet")
			}
			_handler := func (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (bool, *Error) {
					return doHandleExecuteScriptletSsh (_library, _scriptlet, _sshContext, _context)
				}
			return doHandleWithLabel (_library, _scriptlet, _handler, _context)
		
		case "export-scriptlet-body" :
			if _scriptlet == "" {
				return Errorf (0xf24640a2, "export:  expected scriptlet")
			}
			if len (_cleanArguments) != 0 {
				return Errorf (0xcf8db3c0, "export:  unexpected arguments")
			}
			_handler := func (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (bool, *Error) {
					return doHandleExportScriptletBody (_library, _scriptlet, os.Stdout, _context)
				}
			return doHandleWithLabel (_library, _scriptlet, _handler, _context)
		
		
		case "select-execute-scriptlet" :
			if len (_cleanArguments) != 0 {
				return Errorf (0x203e410a, "select:  unexpected arguments")
			}
			if _scriptlet != "" {
				return doSelectHandleWithLabel (_library, _scriptlet, doHandleExecuteScriptlet, _context)
			} else {
				return doSelectHandle (_library, doHandleExecuteScriptlet, _context)
			}
		
		case "select-execute-scriptlet-loop" :
			if len (_cleanArguments) != 0 {
				return Errorf (0x2c6bb7ce, "select:  unexpected arguments")
			}
			for {
				_handled := false
				_handler := func (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (bool, *Error) {
						if _error := executeScriptlet (_library, _scriptlet, true, _context); _error == nil {
							_handled = true
							return true, nil
						} else {
							return false, _error
						}
					}
				if _scriptlet != "" {
					if _error := doSelectHandleWithLabel (_library, _scriptlet, _handler, _context); _error != nil {
						return _error
					}
				} else {
					if _error := doSelectHandle (_library, _handler, _context); _error != nil {
						return _error
					}
				}
				_quit := false
				if ! _handled {
					if _quit_0, _error := menuQuit (_context); _error == nil {
						_quit = _quit_0
					} else {
						return _error
					}
				} else {
					if _quit_0, _error := menuPause (_context); _error == nil {
						_quit = _quit_0
					} else {
						return _error
					}
				}
				if _quit {
					return nil
				} else {
					continue
				}
			}
		
		case "select-export-scriptlet-label" :
			if len (_cleanArguments) != 0 {
				return Errorf (0x2d19b1bc, "select:  unexpected arguments")
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
				return Errorf (0x5f573713, "select:  unexpected arguments")
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
				return Errorf (0xe4f7e6f5, "export:  unexpected arguments")
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
				return Errorf (0xf7b9c7f3, "export:  unexpected scriptlet or arguments")
			}
			return doExportScriptletLabels (_library, os.Stdout, _context)
		
		
		case "parse-library", "parse-library-without-output" :
			if (_scriptlet != "") || (len (_cleanArguments) != 0) {
				return Errorf (0x400ec122, "export:  unexpected scriptlet or arguments")
			}
			if _command == "parse-library" {
				return doExportLibraryJson (_library, os.Stdout, _context)
			} else {
				return nil
			}
		
		case "export-library-json" :
			if (_scriptlet != "") || (len (_cleanArguments) != 0) {
				return Errorf (0xdd0752b8, "export:  unexpected scriptlet or arguments")
			}
			return doExportLibraryStore (_library, NewJsonStreamStoreOutput (os.Stdout, nil), _context)
		
		case "export-library-cdb" :
			if _scriptlet != "" {
				return Errorf (0x492ac50e, "export:  unexpected scriptlet")
			}
			if len (_cleanArguments) != 1 {
				return Errorf (0xf76f4459, "export:  expected database path")
			}
			return doExportLibraryCdb (_library, _cleanArguments[0], _context)
		
		case "export-library-rpc" :
			if _scriptlet != "" {
				return Errorf (0x04d71684, "export:  unexpected scriptlet")
			}
			if len (_cleanArguments) != 1 {
				return Errorf (0xe7886d74, "export:  expected RPC url")
			}
			return doExportLibraryRpc (_library, _cleanArguments[0], _context)
		
		case "export-library-url" :
			if _scriptlet != "" {
				return Errorf (0xff2905b6, "export:  unexpected scriptlet")
			}
			if len (_cleanArguments) != 0 {
				return Errorf (0x3cda6407, "export:  unexpected arguments")
			}
			fmt.Fprintln (os.Stdout, _library.Url ())
			return nil
		
		case "export-library-identifier", "export-library-fingerprint" :
			if _scriptlet != "" {
				return Errorf (0x3483242d, "export:  unexpected scriptlet")
			}
			if len (_cleanArguments) != 0 {
				return Errorf (0x2a741648, "export:  unexpected arguments")
			}
			switch _command {
				case "export-library-identifier" :
					if _identifier, _error := _library.Identifier (); _error == nil {
						fmt.Fprintln (os.Stdout, _identifier)
						return nil
					} else {
						return _error
					}
				case "export-library-fingerprint" :
					if _fingerprint, _error := _library.Fingerprint (); _error == nil {
						fmt.Fprintln (os.Stdout, _fingerprint)
						return nil
					} else {
						return _error
					}
				default :
					panic (0xddb85cc9)
			}
		
		case "" :
			return Errorf (0x5d2a4326, "expected command")
		
		default :
			return Errorf (0x66cf8700, "unexpected command `%s`", _command)
	}
}

