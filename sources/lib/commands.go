

package zrun


import crand "crypto/rand"
import "bytes"
import "encoding/base64"
import "encoding/hex"
import "encoding/json"
import "fmt"
import "io"
import "os"
import "os/exec"
import "path"
import "sort"
import "strings"
import "syscall"




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
	
	if _sources, _error := _library.SelectSources (); _error == nil {
		if _error := _store.Include ("library-meta", "sources", _sources); _error != nil {
			return _error
		}
	}
	
	_fingerprints := make ([]string, 0, 1024)
	_fingerprintsByLabels := make (map[string]string, 1024)
	_labels := make ([]string, 0, 1024)
	_labelsAll := make ([]string, 0, 1024)
	_labelsByFingerprints := make (map[string]string, 1024)
	
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
			if _error := _store.Include ("scriptlets-fingerprint-by-label", _label, _fingerprint); _error != nil {
				return _error
			}
			if _error := _store.Include ("scriptlets-label-by-fingerprint", _fingerprint, _label); _error != nil {
				return _error
			}
			if _error := _store.Include ("scriptlets-meta", _fingerprint, _meta); _error != nil {
				return _error
			}
			_fingerprints = append (_fingerprints, _fingerprint)
			_labelsAll = append (_labelsAll, _label)
			if !_meta.Hidden || _meta.Visible {
				_labels = append (_labels, _label)
			}
			_fingerprintsByLabels[_label] = _fingerprint
			_labelsByFingerprints[_fingerprint] = _label
		}
		
		if _body, _found, _error := _library.ResolveBodyByFingerprint (_fingerprintFromStore); _error == nil {
			if !_found {
				return errorf (0xd80a265e, "invalid store")
			}
			if _error := _store.Include ("scriptlets-body", _fingerprintFromStore, _body); _error != nil {
				return _error
			}
		}
	}
	
	sort.Strings (_fingerprints)
	sort.Strings (_labels)
	sort.Strings (_labelsAll)
	
	if _error := _store.Include ("scriptlets-indices", "fingerprints", _fingerprints); _error != nil {
		return _error
	}
	if _error := _store.Include ("scriptlets-indices", "labels-by-fingerprints", _labelsByFingerprints); _error != nil {
		return _error
	}
	if _error := _store.Include ("scriptlets-indices", "labels", _labels); _error != nil {
		return _error
	}
	if _error := _store.Include ("scriptlets-indices", "labels-all", _labelsAll); _error != nil {
		return _error
	}
	if _error := _store.Include ("scriptlets-indices", "fingerprints-by-labels", _fingerprintsByLabels); _error != nil {
		return _error
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




func executeScriptlet (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (*Error) {
	
	if _scriptlet.Interpreter == "<template>" {
		return executeTemplate (_library, _scriptlet, _context, os.Stdout)
	}
	
	var _command *exec.Cmd
	var _descriptors []int
	if _command_0, _descriptors_0, _error := prepareExecution (_library.Url (), "", _scriptlet, true, _context); _error == nil {
		_command = _command_0
		_descriptors = _descriptors_0
	} else {
		return _error
	}
	
	_closeDescriptors := func () () {
		for _, _descriptor := range _descriptors {
			syscall.Close (_descriptor)
		}
	}
	
	if _command.Dir != "" {
		if _error := os.Chdir (_command.Dir); _error != nil {
			return errorw (0xe4bab179, _error)
		}
	}
	if _command.Stdin != nil {
		_closeDescriptors ()
		return errorf (0x78cfda21, "invalid state")
	}
	if _command.Stdout != nil {
		_closeDescriptors ()
		return errorf (0xf9a9dc74, "invalid state")
	}
	if _command.Stderr != nil {
		_closeDescriptors ()
		return errorf (0xf887025f, "invalid state")
	}
	if _command.ExtraFiles != nil {
		_closeDescriptors ()
		return errorf (0x50354e63, "invalid state")
	}
	if (_command.Process != nil) || (_command.ProcessState != nil) {
		_closeDescriptors ()
		return errorf (0x9d640d1e, "invalid state")
	}
	
	if _error := syscall.Exec (_command.Path, _command.Args, _command.Env); _error != nil {
		_closeDescriptors ()
		return errorw (0x99b54af1, _error)
	} else {
		panic (0xb6dfe17e)
	}
}




func doHandleExecuteScriptletSsh (_library LibraryStore, _scriptlet *Scriptlet, _sshContext *SshContext, _context *Context) (bool, *Error) {
	
	_sshTarget := _sshContext.target
	_sshLauncher := _sshContext.launcher
	_sshDelegate := _sshContext.delegate
	_sshWorkspace := _sshContext.workspace
	_sshCache := _sshContext.cache
	_sshTerminal := _sshContext.terminal
	_sshLibraryLocalSocket := _sshContext.libraryLocalSocket
	_sshLibraryRemoteSocket := _sshContext.libraryRemoteSocket
	_sshExportEnvironment := _sshContext.exportEnvironment
	_sshToken := _sshContext.token
	
	if _sshTarget == "" {
		return false, errorf (0xa552d948, "invalid target:  missing")
	}
	
	if _sshToken == "" {
		var _data [128 / 8]byte
		if _read, _error := crand.Read (_data[:]); _error == nil {
			if _read != (128 / 8) {
				return false, errorf (0xdc5228aa, "invalid state")
			}
		} else {
			return false, errorw (0xc2595acb, _error)
		}
		_sshToken = hex.EncodeToString (_data[:])
	}
	
	if _sshLauncher == "" {
		_sshLauncher = "ssh"
	}
	if strings.IndexByte (_sshLauncher, os.PathSeparator) < 0 {
		if _path, _error := exec.LookPath (_sshLauncher); _error == nil {
			_sshLauncher = _path
		} else {
			return false, errorw (0x9c296054, _error)
		}
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
	
	if _sshWorkspace == "" {
		_sshWorkspace = "/tmp"
	}
	if _sshCache == "" {
		_sshCache = "/tmp"
	}
	if _sshTerminal == "" {
		_sshTerminal = _context.terminal
	}
	
	if _sshLibraryLocalSocket == "" {
		_cacheRoot := _context.cacheRoot
		if _cacheRoot == "" {
			if _cacheRoot_0, _error := resolveCache (); _error == nil {
				_cacheRoot = _cacheRoot_0
			} else {
				return false, _error
			}
		}
		_sshLibraryLocalSocket = path.Join (_cacheRoot, fmt.Sprintf ("%s-%08x.sock", _sshToken, os.Getpid ()))
	}
	if _sshLibraryRemoteSocket == "" {
		_sshLibraryRemoteSocket = path.Join (_sshCache, fmt.Sprintf ("%s.sock", _sshToken))
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
	
	_invokeEnvironment := make (map[string]string, len (_sshExportEnvironment))
	for _, _name := range _sshExportEnvironment {
		if _value, _ok := _context.cleanEnvironment[_name]; _ok {
			_invokeEnvironment[_name] = _value
		}
	}
	
	_invokeContext := & InvokeContext {
			Library : "unix:" + _sshLibraryRemoteSocket,
			Scriptlet : _scriptlet.Label,
			Arguments : _context.cleanArguments,
			Environment : _invokeEnvironment,
			Workspace : _sshWorkspace,
			Cache : _sshCache,
			Terminal : _sshTerminal,
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
		return false, errorw (0x881a60b5, _error)
	}
	
	return true, nil
}




func doHandleExecuteScriptlet (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (bool, *Error) {
	if _error := executeScriptlet (_library, _scriptlet, _context); _error == nil {
		return true, nil
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
						if _, _outputData, _error := loadFromScriptlet (_library.Url (), "", _scriptlet, _context); _error == nil {
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
								return false, nil
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

