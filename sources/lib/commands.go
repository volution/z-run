

package zrun


import "bytes"
import "encoding/json"
import "fmt"
import "io"
import "os/exec"
import "sort"
import "strings"
import "syscall"




func doExportLabelsList (_library LibraryStore, _stream io.Writer, _context *Context) (error) {
	if _labels, _error := _library.SelectLabels (); _error == nil {
		_buffer := bytes.NewBuffer (nil)
		for _, _label := range _labels {
			_buffer.WriteString (_label)
			_buffer.WriteByte ('\n')
		}
		_, _error := _stream.Write (_buffer.Bytes ())
		return _error
	} else {
		return _error
	}
}


func doExportScript (_library LibraryStore, _label string, _stream io.Writer, _context *Context) (error) {
	if _body, _found, _error := _library.ResolveBodyByLabel (_label); _error == nil {
		if _found {
			_, _error := io.WriteString (_stream, _body)
			return _error
		} else {
			return errorf (0x95e0b174, "undefined scriptlet `%s`", _label)
		}
	} else {
		return _error
	}
}


func doExportLibraryJson (_library LibraryStore, _stream io.Writer, _context *Context) (error) {
	_library, _ok := _library.(*Library)
	if !_ok {
		return errorf (0x4f480517, "only works with in-memory library store")
	}
	_encoder := json.NewEncoder (_stream)
	_encoder.SetIndent ("", "    ")
	_encoder.SetEscapeHTML (false)
	return _encoder.Encode (_library)
}


func doExportLibraryStore (_library LibraryStore, _store StoreOutput, _context *Context) (error) {
	
	if _sources, _error := _library.SelectSources (); _error == nil {
		if _error := _store.Include ("library-meta", "sources", _sources); _error != nil {
			return _error
		}
	}
	
	_fingerprints := make ([]string, 0, 1024)
	_fingerprintsByLabels := make (map[string]string, 1024)
	_labels := make ([]string, 0, 1024)
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
			if !_meta.Hidden {
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
	
	if _error := _store.Include ("scriptlets-indices", "fingerprints", _fingerprints); _error != nil {
		return _error
	}
	if _error := _store.Include ("scriptlets-indices", "labels-by-fingerprints", _labelsByFingerprints); _error != nil {
		return _error
	}
	if _error := _store.Include ("scriptlets-indices", "labels", _labels); _error != nil {
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


func doExportLibraryCdb (_library LibraryStore, _path string, _context *Context) (error) {
	if _store, _error := NewCdbStoreOutput (_path); _error == nil {
		return doExportLibraryStore (_library, _store, _context)
	} else {
		return _error
	}
}




func doExecute (_library LibraryStore, _scriptletLabel string, _context *Context) (error) {
	if _scriptlet, _error := _library.ResolveFullByLabel (_scriptletLabel); _error == nil {
		if _scriptlet != nil {
			return doExecuteScriptlet (_library, _scriptlet, _context)
		} else {
			return errorf (0x3be6dcd7, "unknown scriptlet for `%s`", _scriptletLabel)
		}
	} else {
		return _error
	}
}


func doExecuteScriptlet (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (error) {
	
	var _command *exec.Cmd
	var _descriptors []int
	if _command_0, _descriptors_0, _error := prepareExecution (_library, "", _scriptlet, _context); _error == nil {
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
		_closeDescriptors ()
		return errorf (0xe4bab179, "invalid state")
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
		return _error
	} else {
		panic (0xb6dfe17e)
	}
}




func doSelectExecute (_library LibraryStore, _context *Context) (error) {
	for {
		if _label, _error := doSelectLabel_0 (_library, _context); _error == nil {
			if _label != "" {
				if _error := doSelectExecute_0 (_library, _label, _context); _error == nil {
					continue
				} else {
					return _error
				}
			} else {
//				logf ('d', 0x899f4d5a, "no scriptlet selected!")
				return nil
			}
		} else {
			return _error
		}
	}
}

func doSelectExecute_0 (_library LibraryStore, _label string, _context *Context) (error) {
	if _scriptlet, _error := _library.ResolveFullByLabel (_label); _error == nil {
		if _scriptlet != nil {
			if _scriptlet.Kind == "menu" {
				for {
					if _, _outputData, _error := loadFromScriptlet (_library, "", _scriptlet, _context); _error == nil {
						_outputText := string (_outputData)
						_outputText = strings.TrimSpace (_outputText)
						var _outputLines []string
						if _outputText != "" {
							_outputLines = strings.Split (_outputText, "\n")
						}
						if len (_outputLines) == 0 {
							return nil
						} else if len (_outputLines) == 1 {
							_outputLabel := strings.TrimSpace (_outputLines[0])
							if _error := doSelectExecute_0 (_library, _outputLabel, _context); _error == nil {
								continue
							} else {
								return _error
							}
						} else {
							return errorf (0xc8e9bf9b, "invalid output")
						}
					} else {
						return _error
					}
				}
			} else {
				return doExecute (_library, _label, _context)
			}
		} else {
			return errorf (0x3be6dcd7, "unknown scriptlet for `%s`", _label)
		}
	} else {
		return _error
	}
}




func doSelectLegacyOutput (_library LibraryStore, _label string, _stream io.Writer, _context *Context) (error) {
	if _scriptlet, _error := doSelectScriptlet (_library, _label, _context); _error == nil {
		if _scriptlet != nil {
			if _, _error := fmt.Fprintf (_stream, ":: %s\n%s\n", _scriptlet.Label, _scriptlet.Body); _error != nil {
				return _error
			}
			return nil
		} else {
//			logf ('d', 0x13ba2d3b, "no scriptlet selected!")
			return nil
		}
	} else {
		return _error
	}
}




func doSelectScriptlet (_library LibraryStore, _label string, _context *Context) (*Scriptlet, error) {
	if _label == "" {
		if _label_0, _error := doSelectLabel_0 (_library, _context); _error == nil {
			_label = _label_0
		} else {
			return nil, _error
		}
	}
	if _label == "" {
//		logf ('d', 0x16a84cd1, "no scriptlet selected!")
		return nil, nil
	}
	if _scriptlet, _error := _library.ResolveFullByLabel (_label); _error == nil {
		if _scriptlet != nil {
			return _scriptlet, nil
		} else {
			return nil, errorf (0x06ef8e1d, "unknown scriptlet for `%s`", _label)
		}
	} else {
		return nil, _error
	}
}

func doSelectLabel (_library LibraryStore, _stream io.Writer, _context *Context) (error) {
	if _label, _error := doSelectLabel_0 (_library, _context); _error == nil {
		if _label != "" {
			if _, _error := fmt.Fprintf (_stream, "%s\n", _label); _error != nil {
				return _error
			}
		} else {
//			logf ('d', 0x2f4b716a, "no scriptlet selected!")
			return nil
		}
	} else {
		return _error
	}
	return nil
}

func doSelectLabels (_library LibraryStore, _stream io.Writer, _context *Context) (error) {
	if _labels, _error := doSelectLabels_0 (_library, _context); _error == nil {
		if len (_labels) != 0 {
			for _, _label := range _labels {
				if _, _error := fmt.Fprintf (_stream, "%s\n", _label); _error != nil {
					return _error
				}
			}
		} else {
//			logf ('d', 0x9e431309, "no scriptlet selected!")
			return nil
		}
	} else {
		return _error
	}
	return nil
}


func doSelectLabel_0 (_library LibraryStore, _context *Context) (string, error) {
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

func doSelectLabel_1 (_inputs []string, _context *Context) (string, error) {
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

func doSelectLabels_0 (_library LibraryStore, _context *Context) ([]string, error) {
	var _inputs []string
	if _inputs_0, _error := _library.SelectLabels (); _error == nil {
		_inputs = _inputs_0
	} else {
		return nil, _error
	}
	return doSelectLabels_1 (_inputs, _context)
}

func doSelectLabels_1 (_inputs []string, _context *Context) ([]string, error) {
	var _outputs []string
	if _outputs_0, _error := menuSelect (_inputs, _context); _error == nil {
		_outputs = _outputs_0
	} else {
		return nil, _error
	}
	return _outputs, nil
}
