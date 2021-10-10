

package zrun


import "bytes"
import "encoding/base64"
import "encoding/json"
import "fmt"
import "io"
import "os"
import "os/exec"
import "path"
import "sort"
import "strings"

import isatty "github.com/mattn/go-isatty"




func doExportScriptletLabels (_library LibraryStore, _stream io.Writer, _context *Context) (*Error) {
	if _labels, _error := _library.SelectLabelsAll (); _error == nil {
		_buffer := bytes.NewBuffer (nil)
		for _, _label := range _labels {
			_buffer.WriteString (_label)
			_buffer.WriteByte ('\n')
		}
		if _, _error := _stream.Write (_buffer.Bytes ()) ; _error == nil {
			return nil
		} else {
			return errorw (0x1215e523, _error)
		}
	} else {
		return _error
	}
}




func doExportLibraryJson (_library LibraryStore, _stream io.Writer, _context *Context) (*Error) {
	_library, _ok := _library.(*Library)
	if !_ok {
		return errorf (0x4f480517, "only works with in-memory library store")
	}
	_encoder := json.NewEncoder (_stream)
	_encoder.SetIndent ("", "    ")
	_encoder.SetEscapeHTML (false)
	if _error := _encoder.Encode (_library); _error == nil {
		return nil
	} else {
		return errorw (0xb3a1eb99, _error)
	}
}




func doExportLibraryStore (_library LibraryStore, _store StoreOutput, _context *Context) (*Error) {
	
	var _libraryIdentifier string
	if _identifier_0, _error := _library.Identifier (); _error == nil {
		_libraryIdentifier = _identifier_0
	} else {
		return _error
	}
	
	var _libraryFingerprint string
	if _fingerprint_0, _error := _library.Fingerprint (); _error == nil {
		_libraryFingerprint = _fingerprint_0
	} else {
		return _error
	}
	
	if _error := _store.IncludeRawString (_libraryIdentifier, false, "library-meta", "identifier", _libraryIdentifier); _error != nil {
		return _error
	}
	if _error := _store.IncludeRawString (_libraryIdentifier, false, "library-meta", "fingerprint", _libraryFingerprint); _error != nil {
		return _error
	}
	
	if _error := _store.IncludeRawString (_libraryFingerprint, false, "library-meta", "identifier", _libraryIdentifier); _error != nil {
		return _error
	}
	if _error := _store.IncludeRawString (_libraryFingerprint, false, "library-meta", "fingerprint", _libraryFingerprint); _error != nil {
		return _error
	}
	
	if _sources, _error := _library.SelectLibrarySources (); _error == nil {
		if _error := _store.IncludeObject (_libraryFingerprint, false, "library-meta", "library-sources", _sources); _error != nil {
			return _error
		}
	} else {
		return _error
	}
	
	if _context, _error := _library.SelectLibraryContext (); _error == nil {
		if _error := _store.IncludeObject (_libraryFingerprint, false, "library-meta", "library-context", _context); _error != nil {
			return _error
		}
	} else {
		return _error
	}
	
	_fingerprints := make ([]string, 0, 1024)
	_fingerprintsByLabels := make (map[string]string, 1024)
	_labels := make ([]string, 0, 1024)
	_labelsAll := make ([]string, 0, 1024)
	_labelsByFingerprints := make (map[string]string, 1024)
	_contextsIdentifiers := make (map[string]bool, 16)
	
	var _fingerprintsFromStore []string
	if _fingerprints_0, _error := _library.SelectFingerprints (); _error == nil {
		_fingerprintsFromStore = _fingerprints_0
	} else {
		return _error
	}
	
	for _, _fingerprintFromStore := range _fingerprintsFromStore {
		
		if _meta, _error := _library.ResolveMetaByFingerprint (_fingerprintFromStore); _error == nil {
			if _meta == nil {
				return errorf (0x20bc9d40, "invalid store")
			}
			_fingerprint := _meta.Fingerprint
			_label := _meta.Label
			if _error := _store.IncludeRawString (_libraryFingerprint, false, "scriptlets-fingerprint-by-label", _label, _fingerprint); _error != nil {
				return _error
			}
			if _error := _store.IncludeRawString (_libraryFingerprint, false, "scriptlets-label-by-fingerprint", _fingerprint, _label); _error != nil {
				return _error
			}
			if _error := _store.IncludeObject (_libraryFingerprint, true, "scriptlets-meta", _fingerprint, _meta); _error != nil {
				return _error
			}
			_fingerprints = append (_fingerprints, _fingerprint)
			_labelsAll = append (_labelsAll, _label)
			if !_meta.Hidden || _meta.Visible {
				_labels = append (_labels, _label)
			}
			_fingerprintsByLabels[_label] = _fingerprint
			_labelsByFingerprints[_fingerprint] = _label
			if _meta.ContextIdentifier != "" {
				_contextsIdentifiers[_meta.ContextIdentifier] = true
			}
		} else {
			return _error
		}
		
		if _body, _found, _error := _library.ResolveBodyByFingerprint (_fingerprintFromStore); _error == nil {
			if !_found {
				return errorf (0xd80a265e, "invalid store")
			}
			if _error := _store.IncludeRawString (_libraryFingerprint, true, "scriptlets-body", _fingerprintFromStore, _body); _error != nil {
				return _error
			}
		} else {
			return _error
		}
	}
	
	sort.Strings (_fingerprints)
	sort.Strings (_labels)
	sort.Strings (_labelsAll)
	
	if _error := _store.IncludeObject (_libraryFingerprint, false, "scriptlets-indices", "fingerprints", _fingerprints); _error != nil {
		return _error
	}
	if _error := _store.IncludeObject (_libraryFingerprint, false, "scriptlets-indices", "labels-by-fingerprints", _labelsByFingerprints); _error != nil {
		return _error
	}
	if _error := _store.IncludeObject (_libraryFingerprint, false, "scriptlets-indices", "labels", _labels); _error != nil {
		return _error
	}
	if _error := _store.IncludeObject (_libraryFingerprint, false, "scriptlets-indices", "labels-all", _labelsAll); _error != nil {
		return _error
	}
	if _error := _store.IncludeObject (_libraryFingerprint, false, "scriptlets-indices", "fingerprints-by-labels", _fingerprintsByLabels); _error != nil {
		return _error
	}
	
	for _contextIdentifier, _ := range _contextsIdentifiers {
		if _context, _found, _error := _library.ResolveContextByIdentifier (_contextIdentifier); _error == nil {
			if _found && _contextIdentifier == _context.Identifier {
				if _error := _store.IncludeObject (_libraryFingerprint, false, "scriptlet-contexts-by-identifier", _context.Identifier, _context); _error != nil {
					return _error
				}
			} else {
				return errorf (0x7c9d40d9, "invalid store")
			}
		} else {
			return _error
		}
	}
	
	if _error := _store.Commit (); _error != nil {
		return _error
	}
	
	return nil
}


func doExportLibraryCdb (_library LibraryStore, _path string, _context *Context) (*Error) {
	if _store, _error := NewCdbStoreOutput (_path); _error == nil {
		return doExportLibraryStore (_library, _store, _context)
	} else {
		return _error
	}
}


func doExportLibraryRpc (_library LibraryStore, _url string, _context *Context) (*Error) {
	if _server, _error := NewLibraryRpcServer (_library, _url); _error == nil {
		return _server.Serve ()
	} else {
		return _error
	}
}




func doHandleExecuteScriptletSsh (_library LibraryStore, _scriptlet *Scriptlet, _sshContext *SshContext, _context *Context) (bool, *Error) {
	
	_sshTarget := _sshContext.target
	_sshLauncher := _sshContext.launcher
	_sshDelegate := _sshContext.delegate
	_sshExportEnvironment := _sshContext.exportEnvironment
	_sshExecutablePaths := _sshContext.executablePaths
	_sshTerminal := _sshContext.terminal
	_sshWorkspace := _sshContext.workspace
	_sshCache := _sshContext.cache
	_sshLibraryLocalSocket := _sshContext.libraryLocalSocket
	_sshLibraryRemoteSocket := _sshContext.libraryRemoteSocket
	_sshToken := _sshContext.token
	
	if _sshTarget == "" {
		return false, errorf (0xa552d948, "invalid target:  missing")
	}
	
	if _sshToken == "" {
		_sshToken = generateRandomToken ()
	}
	
	if _sshLauncher == "" {
		_sshLauncher = "ssh"
	}
	if _sshLauncher_0, _error := resolveExecutable (_sshLauncher, _context.executablePaths); _error == nil {
		_sshLauncher = _sshLauncher_0
	} else {
		return false, _error
	}
	
	if _sshDelegate == "" {
		_sshDelegate = "z-run"
	}
	if fmt.Sprintf ("%+q", _sshDelegate) != ("\"" + _sshDelegate + "\"") {
		return false, errorf (0x230a2fc4, "invalid delegate:  non ASCII")
	}
	if strings.ContainsAny (_sshDelegate, " !\"#$%&'()*+,:;<=>?@[\\]^`{|}~") {
		return false, errorf (0x1e5c7a40, "invalid delegate:  contains disallowed special character")
	}
	
	var _invokeEnvironment map[string]string
	if _sshExportEnvironment != nil {
		_invokeEnvironment = make (map[string]string, len (_sshExportEnvironment))
		for _, _name := range _sshExportEnvironment {
			if _value, _ok := _context.cleanEnvironment[_name]; _ok {
				_invokeEnvironment[_name] = _value
			}
		}
	}
	
	var _invokeExecutablePaths []string
	if _sshExecutablePaths != nil {
		_invokeExecutablePaths = make ([]string, 0, len (_sshExecutablePaths))
		for _, _path := range _sshExecutablePaths {
			_invokeExecutablePaths = append (_invokeExecutablePaths, _path)
		}
	} else {
		_invokeExecutablePaths = []string {
				"/usr/local/bin",
				"/usr/local/sbin",
				"/usr/bin",
				"/usr/sbin",
				"/bin",
				"/sbin",
			}
	}
	
	if _sshTerminal == "" {
		_sshTerminal = _context.terminal
	}
	
	if _sshWorkspace == "" {
		_sshWorkspace = "/tmp"
	}
	if _sshCache == "" {
		_sshCache = "/tmp"
	}
	
	if _error := makeCacheFolder (_context.cacheRoot, "ssh-sockets"); _error != nil {
		return false, _error
	}
	if _sshLibraryLocalSocket == "" {
		_sshLibraryLocalSocket = path.Join (_context.cacheRoot, "ssh-sockets", fmt.Sprintf ("%s-%08x.sock", _sshToken, os.Getpid ()))
	}
	if _sshLibraryRemoteSocket == "" {
		_sshLibraryRemoteSocket = path.Join (_sshCache, fmt.Sprintf ("z-run--%s.sock", _sshToken))
	}
	
	var _rpc *LibraryRpcServer
	if _rpc_0, _error := NewLibraryRpcServer (_library, "unix:" + _sshLibraryLocalSocket); _error == nil {
		_rpc = _rpc_0
	} else {
		return false, _error
	}
	if _error := _rpc.ServeStart (); _error != nil {
		return false, _error
	}
	defer _rpc.ServeStop ()
	
	_invokeContext := & InvokeContext {
			Version : BUILD_VERSION,
			Library : "unix:" + _sshLibraryRemoteSocket,
			Scriptlet : _scriptlet.Label,
			Arguments : _context.cleanArguments,
			Environment : _invokeEnvironment,
			ExecutablePaths : _invokeExecutablePaths,
			Terminal : _sshTerminal,
			Workspace : _sshWorkspace,
			CacheRoot : _sshCache,
		}
	
	var _invokeContextEncoded string
	if _data, _error := json.Marshal (_invokeContext); _error == nil {
		_invokeContextEncoded = base64.RawURLEncoding.EncodeToString (_data)
	} else {
		return false, errorw (0x5bbc9bcc, _error)
	}
	
	_sshArguments := make ([]string, 0, 16)
	_sshArguments = append (_sshArguments, _sshLauncher)
	if _sshTerminal == "" {
		_sshArguments = append (_sshArguments, "-T")
	} else {
		_stdinIsTty := isatty.IsTerminal (os.Stdin.Fd ())
		_stdoutIsTty := isatty.IsTerminal (os.Stdout.Fd ())
		_stderrIsTty := isatty.IsTerminal (os.Stderr.Fd ())
		if !_stdinIsTty || !_stdoutIsTty || !_stderrIsTty {
			// NOP
		} else if _stdinIsTty && _stdoutIsTty && _stderrIsTty {
//			logf ('d', 0x93bc5a69, "SSH with TTY allowed;")
			_sshArguments = append (_sshArguments, "-t")
		}
	}
	
	_sshArguments = append (_sshArguments, "-R", _sshLibraryRemoteSocket + ":" + _sshLibraryLocalSocket)
	_sshArguments = append (_sshArguments, "--")
	_sshArguments = append (_sshArguments, _sshTarget)
	_sshArguments = append (_sshArguments, "exec", _sshDelegate, "--invoke", _invokeContextEncoded)
	
	_sshCommand := & exec.Cmd {
			Path : _sshLauncher,
			Args : _sshArguments,
			Env : nil,
			Dir : "",
			Stdin : os.Stdin,
			Stdout : os.Stdout,
			Stderr : os.Stderr,
		}
	
	if _error := _sshCommand.Run (); _error != nil {
		return false, errorf (0x881a60b5, "failed to spawn `%s`  //  %v", _sshCommand.Path, _error)
	}
	
	return true, nil
}




func doHandleExecuteScriptlet (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (bool, *Error) {
	if _error := executeScriptlet (_library, _scriptlet, false, _context); _error == nil {
		panic (0x64a96e7f)
	} else {
		return false, _error
	}
}


func doHandleExportScriptletLabel (_library LibraryStore, _scriptlet *Scriptlet, _stream io.Writer, _context *Context) (bool, *Error) {
	if _, _error := fmt.Fprintf (_stream, "%s\n", _scriptlet.Label); _error == nil {
		return true, nil
	} else {
		return false, errorw (0x75e2481b, _error)
	}
}


func doHandleExportScriptletBody (_library LibraryStore, _scriptlet *Scriptlet, _stream io.Writer, _context *Context) (bool, *Error) {
	if _body, _found, _error := _library.ResolveBodyByLabel (_scriptlet.Label); _error == nil {
		if _found {
			if _, _error := io.WriteString (_stream, _body); _error == nil {
				return true, nil
			} else {
				return false, errorw (0x313551b7, _error)
			}
		} else {
			return false, errorf (0x95e0b174, "undefined scriptlet body `%s`", _scriptlet.Label)
		}
	} else {
		return false, _error
	}
}


func doHandleExportScriptletLegacy (_library LibraryStore, _scriptlet *Scriptlet, _stream io.Writer, _context *Context) (bool, *Error) {
	if _, _error := fmt.Fprintf (_stream, ":: %s\n%s\n", _scriptlet.Label, _scriptlet.Body); _error == nil {
		return true, nil
	} else {
		return false, errorw (0x1827808c, _error)
	}
}




type doHandler func (LibraryStore, *Scriptlet, *Context) (bool, *Error)


func doHandleWithLabel (_library LibraryStore, _label string, _handler doHandler, _context *Context) (*Error) {
	if _scriptlet, _error := _library.ResolveFullByLabel (_label); _error == nil {
		if _scriptlet != nil {
			if _handled, _error := _handler (_library, _scriptlet, _context); _error == nil {
				if _handled {
					return nil
				} else {
					return errorf (0x832d3041, "unhandled scriptlet `%s`", _label)
				}
			} else {
				return _error
			}
		} else {
			return errorf (0x3be6dcd7, "unknown scriptlet for `%s`", _label)
		}
	} else {
		return _error
	}
}




func doSelectHandle (_library LibraryStore, _handler doHandler, _context *Context) (*Error) {
	return doSelectHandle_1 (_library, _handler, _context)
}

func doSelectHandleWithLabel (_library LibraryStore, _label string, _handler doHandler, _context *Context) (*Error) {
	if _label != "" {
		if _handled, _error := doSelectHandle_2 (_library, _label, _handler, _context); _error == nil {
			if _handled {
				return nil
			} else {
//				return errorf (0xb78ed426, "unhandled scriptlet `%s`", _label)
				return nil
			}
		} else {
			return _error
		}
	} else {
		return errorf (0xab61361e, "invalid label")
	}
}


func doSelectHandle_1 (_library LibraryStore, _handler doHandler, _context *Context) (*Error) {
	for {
		if _label, _error := doSelectLabel_0 (_library, _context); _error == nil {
			if _label != "" {
				if _handled, _error := doSelectHandle_2 (_library, _label, _handler, _context); _error == nil {
					if _handled {
						return nil
					} else {
						continue
					}
				} else {
					return _error
				}
			} else {
				return nil
			}
		} else {
			return _error
		}
	}
}


func doSelectHandle_2 (_library LibraryStore, _label string, _handler doHandler, _context *Context) (bool, *Error) {
	if _scriptlet, _error := _library.ResolveFullByLabel (_label); _error == nil {
		if _scriptlet != nil {
			if _scriptlet.Kind == "menu" {
				for {
					var _outputLines []string
					if _scriptlet.Interpreter != "<menu>" {
						_libraryFingerprint := ""
						if _libraryFingerprint_0, _error := _library.Fingerprint (); _error == nil {
							_libraryFingerprint = _libraryFingerprint_0
						} else {
							return false, _error
						}
						if _, _outputData, _error := loadFromScriptlet (_library.Url (), _libraryFingerprint, "", _scriptlet, _context); _error == nil {
							_outputText := string (_outputData)
							_outputText = strings.TrimSpace (_outputText)
							if _outputText != "" {
								_outputLines = strings.Split (_outputText, "\n")
							}
						} else {
							return false, _error
						}
					} else {
						_inputLines := strings.Split (_scriptlet.Body, "\n")
						if len (_inputLines) > 0 {
							if _inputLines[len (_inputLines) - 1] == "" {
								_inputLines = _inputLines[: len (_inputLines) - 1]
							}
						}
						if _outputLines_0, _error := menuSelect (_inputLines, _context); _error == nil {
							_outputLines = _outputLines_0
						} else {
							return false, _error
						}
					}
					if len (_outputLines) == 0 {
						return false, nil
					} else if len (_outputLines) == 1 {
						_outputLabel := strings.TrimSpace (_outputLines[0])
						if _handled, _error := doSelectHandle_2 (_library, _outputLabel, _handler, _context); _error == nil {
							if _handled {
								return true, nil
							} else {
								continue
							}
						} else {
							return false, _error
						}
					} else {
						return false, errorf (0xc8e9bf9b, "invalid output")
					}
				}
			} else {
				if _handled, _error := _handler (_library, _scriptlet, _context); _error == nil {
					if _handled {
						return true, nil
					} else {
						return false, errorf (0x4be0f083, "unhandled scriptlet `%s`", _label)
					}
				} else {
					return false, _error
				}
			}
		} else {
			return false, errorf (0x76c06125, "unknown scriptlet for `%s`", _label)
		}
	} else {
		return false, _error
	}
}




func doSelectScriptlet (_library LibraryStore, _label string, _context *Context) (*Scriptlet, *Error) {
	if _label == "" {
		if _label_0, _error := doSelectLabel_0 (_library, _context); _error == nil {
			_label = _label_0
		} else {
			return nil, _error
		}
	}
	if _label == "" {
		return nil, nil
	}
	if _scriptlet, _error := _library.ResolveMetaByLabel (_label); _error == nil {
		if _scriptlet != nil {
			return _scriptlet, nil
		} else {
			return nil, errorf (0x06ef8e1d, "unknown scriptlet for `%s`", _label)
		}
	} else {
		return nil, _error
	}
}

func doSelectLabel (_library LibraryStore, _stream io.Writer, _context *Context) (*Error) {
	if _label, _error := doSelectLabel_0 (_library, _context); _error == nil {
		if _label != "" {
			if _, _error := fmt.Fprintf (_stream, "%s\n", _label); _error != nil {
				return errorw (0x0b7f5d09, _error)
			}
		} else {
			return nil
		}
	} else {
		return _error
	}
	return nil
}

func doSelectLabels (_library LibraryStore, _stream io.Writer, _context *Context) (*Error) {
	if _labels, _error := doSelectLabels_0 (_library, _context); _error == nil {
		if len (_labels) != 0 {
			for _, _label := range _labels {
				if _, _error := fmt.Fprintf (_stream, "%s\n", _label); _error != nil {
					return errorw (0x7f0c1122, _error)
				}
			}
		} else {
			return nil
		}
	} else {
		return _error
	}
	return nil
}




func doSelectLabel_0 (_library LibraryStore, _context *Context) (string, *Error) {
	if _labels, _error := doSelectLabels_0 (_library, _context); _error == nil {
		if len (_labels) == 0 {
			return "", nil
		} else if len (_labels) == 1 {
			return _labels[0], nil
		} else {
			return "", errorf (0xa11d1022, "no scriptlet selected")
		}
	} else {
		return "", _error
	}
}

func doSelectLabel_1 (_inputs []string, _context *Context) (string, *Error) {
	if _labels, _error := doSelectLabels_1 (_inputs, _context); _error == nil {
		if len (_labels) == 0 {
			return "", nil
		} else if len (_labels) == 1 {
			return _labels[0], nil
		} else {
			return "", errorf (0xb7ae6203, "no scriptlet selected")
		}
	} else {
		return "", _error
	}
}


func doSelectLabels_0 (_library LibraryStore, _context *Context) ([]string, *Error) {
	var _inputs []string
	if _inputs_0, _error := _library.SelectLabels (); _error == nil {
		_inputs = _inputs_0
	} else {
		return nil, _error
	}
	return doSelectLabels_1 (_inputs, _context)
}

func doSelectLabels_1 (_inputs []string, _context *Context) ([]string, *Error) {
	var _outputs []string
	if _outputs_0, _error := menuSelect (_inputs, _context); _error == nil {
		_outputs = _outputs_0
	} else {
		return nil, _error
	}
	return _outputs, nil
}

