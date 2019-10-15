

package zrun


import "io"
import "io/ioutil"
import "os"
import "os/exec"
import "path"
import "path/filepath"
import "sort"
import "strings"
import "unicode"
import "unicode/utf8"




func parseLibrary (_sources []*Source, _environmentFingerprint string, _context *Context) (*Library, error) {
	
	_library := NewLibrary ()
	_library.EnvironmentFingerprint = _environmentFingerprint
	
	for _, _source := range _sources {
		if _error := parseFromSource (_library, _source, _context); _error != nil {
			return nil, _error
		}
	}
	
	for {
		_found := false
		for _, _scriptlet := range _library.Scriptlets {
			switch _scriptlet.Kind {
				case "executable", "generator" :
					// NOP
				case "generator-pending" :
					if _error := parseFromGenerator (_library, _scriptlet, _context); _error == nil {
						_scriptlet.Kind = "generator"
						_found = true
					} else {
						return nil, _error
					}
				case "replacer-pending" :
					if _error := parseFromReplacer (_library, _scriptlet, _context); _error == nil {
						_scriptlet.Kind = "executable"
						_found = true
					} else {
						return nil, _error
					}
				case "menu-pending" :
					// NOP
				default :
					return nil, errorf (0xd5f0c788, "invalid state `%s`", _scriptlet.Kind)
			}
		}
		if !_found {
			break
		}
	}
	
	{
		for _, _scriptlet := range _library.Scriptlets {
			switch _scriptlet.Kind {
				case "menu-pending" :
					if _error := parseFromMenu (_library, _scriptlet, _context); _error == nil {
						_scriptlet.Kind = "menu"
					} else {
						return nil, _error
					}
				default :
					// NOP
			}
		}
	}
	
	{
		for _, _scriptlet := range _library.Scriptlets {
			sort.Strings (_scriptlet.Menus)
			if len (_scriptlet.Menus) != 0 {
				_scriptlet.Hidden = true
			}
		}
	}
	
	{
		sort.Sort (_library.Scriptlets)
		sort.Strings (_library.ScriptletFingerprints)
		_library.ScriptletLabels = make ([]string, 0, len (_library.Scriptlets))
		for _index, _scriptlet := range _library.Scriptlets {
			_library.ScriptletsByFingerprint[_scriptlet.Fingerprint] = uint (_index)
			_library.ScriptletsByLabel[_scriptlet.Label] = uint (_index)
			if !_scriptlet.Hidden || _scriptlet.Visible {
				_library.ScriptletLabels = append (_library.ScriptletLabels, _scriptlet.Label)
			}
		}
	}
	
	{
		sort.Sort (_library.Sources)
		_fingerprints := make ([]string, 0, len (_library.Sources))
		for _, _source := range _library.Sources {
			_fingerprints = append (_fingerprints, _source.FingerprintData)
		}
		sort.Strings (_fingerprints)
		_library.SourcesFingerprint = NewFingerprinter () .StringsWithLen (_fingerprints) .Build ()
	}
	
	return _library, nil
}




func parseFromGenerator (_library *Library, _source *Scriptlet, _context *Context) (error) {
	if _, _data, _error := loadFromScriptlet (_library, "", _source, _context); _error == nil {
		return parseFromData (_library, _data, _source.Source.Path, _context)
	} else {
		return _error
	}
}

func parseFromReplacer (_library *Library, _source *Scriptlet, _context *Context) (error) {
	if _, _data, _error := loadFromScriptlet (_library, "<shell>", _source, _context); _error == nil {
		if utf8.Valid (_data) {
			_source.Body = string (_data)
			return nil
		} else {
			return errorf (0xdb3b92b7, "invalid UTF-8")
		}
	} else {
		return _error
	}
}

func parseFromMenu (_library *Library, _source *Scriptlet, _context *Context) (error) {
	_labels := make ([]string, 0, 1024)
	_matchers := strings.Split (_source.Body, "\n")
	for _, _scriptlet := range _library.Scriptlets {
		if _scriptlet == _source {
			continue
		}
		if _scriptlet.Hidden {
			continue
		}
		for _, _matcher := range _matchers {
			if _matcher == "" {
				continue
			}
			if _matcher == "*" {
				_labels = append (_labels, _scriptlet.Label)
			} else if strings.HasPrefix (_matcher, "+^ ") {
				_pattern := _matcher[3:]
				if strings.HasPrefix (_scriptlet.Label, _pattern) {
					_labels = append (_labels, _scriptlet.Label)
					_scriptlet.Menus = append (_scriptlet.Menus, _source.Label)
				}
			} else {
				return errorf (0xa068e934, "invalid menu matcher `%s`", _matcher)
			}
		}
	}
	sort.Strings (_labels)
	_buffer := strings.Builder {}
	_previousLabel := ""
	for _, _label := range _labels {
		if _label == _previousLabel {
			continue
		}
		_buffer.WriteString (_label)
		_buffer.WriteByte ('\n')
		_previousLabel = _label
	}
	_source.Body = _buffer.String ()
	return nil
}




func parseFromSource (_library *Library, _source *Source, _context *Context) (error) {
	if _data, _error := loadFromSource (_library, _source, _context); _error == nil {
		if _error := includeSource (_library, _source); _error != nil {
			return _error
		}
		return parseFromData (_library, _data, _source.Path, _context)
	} else {
		return _error
	}
}

func parseFromFile (_library *Library, _sourcePath string, _context *Context) (string, error) {
	if _fingerprint, _data, _error := loadFromFile (_sourcePath); _error == nil {
		if _error := parseFromData (_library, _data, _sourcePath, _context); _error == nil {
			return _fingerprint, nil
		} else {
			return "", _error
		}
	} else {
		return "", _error
	}
}

func parseFromStream (_library *Library, _stream io.Reader, _sourcePath string, _context *Context) (string, error) {
	if _fingerprint, _data, _error := loadFromStream (_stream); _error == nil {
		if _error := parseFromData (_library, _data, _sourcePath, _context); _error == nil {
			return _fingerprint, nil
		} else {
			return "", _error
		}
	} else {
		return "", _error
	}
}




func parseFromData (_library *Library, _sourceData []byte, _sourcePath string, _context *Context) (error) {
	
	var _source string
	if utf8.Valid (_sourceData) {
		_source = string (_sourceData)
	} else {
		return errorf (0x2a19cfc7, "invalid UTF-8 source")
	}
	
	const (
		WAITING = 1 + iota
		SCRIPTLET_BODY
		SCRIPTLET_PUSH
		SKIPPING
	)
	
	type scriptletState struct {
		label string
		kind string
		interpreter string
		disabled bool
		hidden bool
		visible bool
		body string
		bodyBuffer strings.Builder
		bodyStrip string
		bodyLines uint
		lineStart uint
		lineEnd uint
	}
	
	_trimCr := func (s string) (string) { return strings.TrimFunc (s, func (r rune) (bool) { return r == '\r' }) }
	_trimSpace := func (s string) (string) { return strings.TrimFunc (s, unicode.IsSpace) }
	_trimRightSpace := func (s string) (string) { return strings.TrimRightFunc (s, unicode.IsSpace) }
	
	_state := WAITING
	_remaining := _source
	_lineIndex := uint (0)
	var _scriptletState scriptletState
	
	for {
		
		_remaining = _trimCr (_remaining)
		_remainingLen := len (_remaining)
		if _remainingLen == 0 {
			break
		}
		
		var _line string
		if _lineEnd := strings.IndexByte (_remaining, '\n'); _lineEnd >= 0 {
			_line = _remaining[:_lineEnd]
			_remaining = _remaining[_lineEnd + 1:]
		} else {
			_line = _remaining
			_remaining = ""
		}
		
		_line = _trimCr (_line)
		_lineIndex += 1
		
//		logf ('d', 0xc2d2b73d, "processing line (%d):  %s", _lineIndex, _line)
		
		switch _state {
			
			case WAITING :
				_lineTrimmed := _trimRightSpace (_line)
				
				_disabled := false
				if strings.HasPrefix (_lineTrimmed, "##") && (_lineTrimmed != "##") {
					_lineTrimmed = _lineTrimmed[2:]
					_disabled = true
				}
				_visible := false
				if strings.HasPrefix (_lineTrimmed, "++") && (_lineTrimmed != "++") {
					_lineTrimmed = _lineTrimmed[2:]
					_visible = true
				}
				_hidden := false
				if strings.HasPrefix (_lineTrimmed, "--") && (_lineTrimmed != "--") {
					_lineTrimmed = _lineTrimmed[2:]
					_hidden = true
				}
				
				if _lineTrimmed == "" {
					// NOP
					
				} else if strings.HasPrefix (_lineTrimmed, ":: ") ||
						strings.HasPrefix (_lineTrimmed, "::.. ") ||
						strings.HasPrefix (_lineTrimmed, "::~~ ") ||
						strings.HasPrefix (_lineTrimmed, "::~~.. ") ||
						strings.HasPrefix (_lineTrimmed, "::&& ") ||
						strings.HasPrefix (_lineTrimmed, "::&&.. ") ||
						strings.HasPrefix (_lineTrimmed, "::== ") ||
						strings.HasPrefix (_lineTrimmed, "::// ") {
					
					_prefix := _lineTrimmed[: strings.IndexByte (_lineTrimmed, ' ')]
					_label := ""
					_body := ""
					{
						_text := _lineTrimmed[len (_prefix) + 1:]
						if _splitIndex := strings.Index (_text, " :: "); _splitIndex >= 0 {
							_label = _text[:_splitIndex]
							_body = _text[_splitIndex + 4:]
						} else if _prefix != ":://" {
							return errorf (0x53eafa1a, "invalid syntax (%d):  missing scriptlet separator `::` | %s", _lineIndex, _line)
						} else {
							_label = _text
						}
						_label = _trimSpace (_label)
						_body = _trimSpace (_body)
					}
					
					if _label == "" {
						return errorf (0xddec2340, "invalid syntax (%d):  empty scriptlet label | %s", _lineIndex, _line)
					}
					if (_body == "") && (_prefix != ":://") {
						return errorf (0xc1dc94cc, "invalid syntax (%d):  empty scriptlet body | %s", _lineIndex, _line)
					}
					
					_kind := ""
					_interpreter := ""
					_include := false
					switch _prefix[2:] {
						case "" :
							_kind = "executable"
							_interpreter = "<shell>"
						case ".." :
							_kind = "executable"
							_interpreter = "<print>"
						case "~~" :
							_kind = "replacer"
							_interpreter = "<shell>"
						case "~~.." :
							_kind = "replacer"
							_interpreter = "<print>"
						case "==" :
							_kind = "generator"
							_interpreter = "<shell>"
							_hidden = true
						case "&&" :
							_kind = "executable"
							_interpreter = "<shell>"
							_include = true
						case "&&.." :
							_kind = "executable"
							_interpreter = "<print>"
							_include = true
						case "//" :
							_kind = "menu"
							_interpreter = "<menu>"
							_include = false
							if _body == "" {
								if (_label == "*") || (_label == "* ...") {
									_body = "*"
								} else {
									_body = "+^ " + _label
									if strings.HasSuffix (_body, " ...") {
										_body = _body[:len (_body) - 3]
									}
								}
							}
						default :
							return errorf (0xfba805b9, "invalid syntax (%d):  unknown scriptlet type | %s", _lineIndex, _line)
					}
					
					if _include {
						_includePath := path.Join (path.Dir (_sourcePath), _body)
						if _includeSource, _error := resolveSource (_includePath); _error == nil {
							if _data, _error := loadFromSource (_library, _includeSource, _context); _error == nil {
								if _error := includeSource (_library, _includeSource); _error != nil {
									return _error
								}
								if utf8.Valid (_sourceData) {
									_body = string (_data)
								} else {
									return errorf (0x16010e20, "invalid UTF-8")
								}
							} else {
								return _error
							}
						}
					} else {
						_body = _body + "\n"
					}
					
					_scriptletState = scriptletState {
							label : _label,
							kind : _kind,
							interpreter : _interpreter,
							disabled : _disabled,
							visible : _visible,
							hidden : _hidden,
							body : _body,
							lineStart : _lineIndex,
							lineEnd : _lineIndex,
						}
					
					_state = SCRIPTLET_PUSH
					
				} else if strings.HasPrefix (_lineTrimmed, "<< ") ||
						strings.HasPrefix (_lineTrimmed, "<<.. ") ||
						strings.HasPrefix (_lineTrimmed, "<<~~ ") ||
						strings.HasPrefix (_lineTrimmed, "<<~~.. ") ||
						strings.HasPrefix (_lineTrimmed, "<<== ") {
					
					_prefix := _lineTrimmed[: strings.IndexByte (_lineTrimmed, ' ')]
					_label := ""
					{
						_label = _lineTrimmed[len (_prefix) + 1:]
						_label = _trimSpace (_label)
					}
					
					if _label == "" {
						return errorf (0x64c17a76, "invalid syntax (%d):  empty scriptlet label | %s", _lineIndex, _line)
					}
					
					_kind := ""
					_interpreter := ""
					switch _prefix[2:] {
						case "" :
							_kind = "executable"
							_interpreter = "<shell>"
						case ".." :
							_kind = "executable"
							_interpreter = "<print>"
						case "~~" :
							_kind = "replacer"
							_interpreter = "<shell>"
						case "~~.." :
							_kind = "replacer"
							_interpreter = "<print>"
						case "==" :
							_kind = "generator"
							_interpreter = "<shell>"
							_hidden = true
						default :
							return errorf (0xd08972fe, "invalid syntax (%d):  unknown scriptlet type | %s", _lineIndex, _line)
					}
					
					_scriptletState = scriptletState {
							label : _label,
							kind : _kind,
							interpreter : _interpreter,
							disabled : _disabled,
							visible : _visible,
							hidden : _hidden,
							lineStart : _lineIndex,
						}
					
					_state = SCRIPTLET_BODY
					
				} else if strings.HasPrefix (_lineTrimmed, "&& ") {
					
					_includePath := _lineTrimmed[strings.IndexByte (_lineTrimmed, ' ') + 1:]
					_includePath = path.Join (path.Dir (_sourcePath), _includePath)
					
					if !_disabled {
						if _includeSources, _error := resolveSources (_includePath); _error == nil {
							for _, _includeSource := range _includeSources {
								if _error := parseFromSource (_library, _includeSource, _context); _error != nil {
									return _error
								}
							}
						} else {
							return _error
						}
					}
					
				} else if strings.HasPrefix (_lineTrimmed, "&&__ ") {
					
					_includePath := _lineTrimmed[strings.IndexByte (_lineTrimmed, ' ') + 1:]
					_includePath = path.Join (path.Dir (_sourcePath), _includePath)
					
					if !_disabled {
						if _includeSource, _error := fingerprintSource (_includePath); _error == nil {
							if _error := includeSource (_library, _includeSource); _error != nil {
								return _error
							}
						} else {
							return _error
						}
					}
					
				} else if strings.HasPrefix (_lineTrimmed, "{{") {
					if _disabled {
						_state = SKIPPING
					} else {
						return errorf (0x79d4d781, "invalid syntax (%d):  unknown block type | %s", _lineIndex, _line)
					}
					
				} else if strings.HasPrefix (_line, "#!/") {
					if (_lineIndex == 1) {
						// NOP
					} else {
						// FIXME:  This should be a warning!
					}
					
				} else {
					return errorf (0x9f8daae4, "invalid syntax (%d):  unexpected statement | %s", _lineIndex, _line)
				}
			
			case SCRIPTLET_BODY :
				_lineTrimmed := _trimRightSpace (_line)
				
				if _lineTrimmed == "!!" {
					_scriptletState.body = _scriptletState.bodyBuffer.String ()
					_scriptletState.lineEnd = _lineIndex
					_state = SCRIPTLET_PUSH
					
				} else if strings.HasPrefix (_lineTrimmed, "!!") {
					return errorf (0xf9900c0c, "invalid syntax (%d):  unexpected statement | %s", _lineIndex, _line)
					
				} else if _lineTrimmed == "" {
					_scriptletState.bodyBuffer.WriteByte ('\n')
					
				} else {
					if _scriptletState.bodyLines == 0 {
						if _stripIndex := strings.IndexFunc (_line, func (r rune) (bool) { return ! unicode.IsSpace (r) }); _stripIndex > 0 {
							_scriptletState.bodyStrip = _line[:_stripIndex]
						}
					}
					if ! strings.HasPrefix (_line, _scriptletState.bodyStrip) {
						logf ('w', 0xc4e05443, "invalid syntax (%d):  unexpected indentation `%s`", _lineIndex, _line)
					}
					_bodyLine := _line[len (_scriptletState.bodyStrip):]
					_scriptletState.bodyBuffer.WriteString (_bodyLine)
					_scriptletState.bodyBuffer.WriteByte ('\n')
					_scriptletState.bodyLines += 1
				}
			
			case SKIPPING :
				_lineTrimmed := _trimRightSpace (_line)
				if _lineTrimmed == "##}}" {
					_state = WAITING
				} else if strings.HasPrefix (_lineTrimmed, "##}}") {
					return errorf (0x183de0fd, "invalid syntax (%d):  unexpected statement | %s", _lineIndex, _line)
				} else {
					// NOP
				}
		}
		
		if _state == SCRIPTLET_PUSH {
			if !_scriptletState.disabled {
				_scriptlet := & Scriptlet {
						Label : _scriptletState.label,
						Kind : _scriptletState.kind,
						Interpreter : _scriptletState.interpreter,
						Visible : _scriptletState.visible,
						Hidden : _scriptletState.hidden,
						Body : _scriptletState.body,
						Source : ScriptletSource {
								Path : _sourcePath,
								LineStart : _scriptletState.lineStart,
								LineEnd : _scriptletState.lineEnd,
							},
					}
				if _error := includeScriptlet (_library, _scriptlet); _error != nil {
					return _error
				}
			}
			_state = WAITING
		}
	}
	
	switch _state {
		case WAITING :
		case SCRIPTLET_BODY :
			return errorf (0x9d55df33, "invalid syntax (%d):  missing scriptlet body closing tag `!!` (and reached end of file)", _lineIndex)
		case SKIPPING :
			return errorf (0x357f15e1, "invalid syntax (%d):  missing comment body closing tag `##}}` (and reached end of file)", _lineIndex)
		default :
			return errorf (0xc0f78380, "invalid syntax (%d):  unexpected state `%s` (and reached end of file)", _lineIndex, _state)
	}
	
	return nil
}




func loadFromSource (_library *Library, _source *Source, _context *Context) ([]byte, error) {
	if _fingerprint, _data, _error := loadFromSource_0 (_library, _source, _context); _error == nil {
		if _source.FingerprintData == "" {
			_source.FingerprintData = _fingerprint
		} else if _source.FingerprintData != _fingerprint {
			return nil, errorf (0x6293d72d, "invalid state")
		}
		return _data, nil
	} else {
		return nil, _error
	}
}


func loadFromSource_0 (_library *Library, _source *Source, _context *Context) (string, []byte, error) {
	
	if !_source.Executable {
		return loadFromFile (_source.Path)
		
	} else {
		
		_executable := _source.Path
		if ! path.IsAbs (_executable) {
			if _executable_0, _error := filepath.Abs (_executable); _error == nil {
				_executable = _executable_0
			} else {
				return "", nil, _error
			}
		}
		
		// FIXME:  Hash the actual executable!
		_fingerprint := _source.FingerprintMeta
		
		_command := & exec.Cmd {
				Path : _executable,
				Args : []string {
						"[z-run:generator]",
					},
				Env : processEnvironment (_context, nil),
				Stdin : nil,
				Stdout : nil,
				Stderr : os.Stderr,
				Dir : "",
			}
		
		if _, _data, _error := loadFromCommand (_command); _error == nil {
			return _fingerprint, _data, _error
		} else {
			return "", nil, _error
		}
	}
}


func loadFromFile (_path string) (string, []byte, error) {
	if _stream, _error := os.Open (_path); _error == nil {
		defer _stream.Close ()
		return loadFromStream (_stream)
	} else {
		return "", nil, _error
	}
}


func loadFromStream (_stream io.Reader) (string, []byte, error) {
	if _data, _error := ioutil.ReadAll (_stream); _error == nil {
		_fingerprint := NewFingerprinter () .Bytes (_data) .Build ()
		return _fingerprint, _data, nil
	} else {
		return "", nil, _error
	}
}


func loadFromScriptlet (_library LibraryStore, _interpreter string, _scriptlet *Scriptlet, _context *Context) (string, []byte, error) {
	
	var _command *exec.Cmd
	var _descriptors []int
	if _command_0, _descriptors_0, _error := prepareExecution (_library, _interpreter, _scriptlet, false, _context); _error == nil {
		_command = _command_0
		_descriptors = _descriptors_0
	} else {
		return "", nil, _error
	}
	
	for _, _descriptor := range _descriptors {
		_command.ExtraFiles = append (_command.ExtraFiles, os.NewFile (uintptr (_descriptor), ""))
	}
	
	return loadFromCommand (_command)
}


func loadFromCommand (_command *exec.Cmd) (string, []byte, error) {
	
	if _command.Stderr == nil {
		_command.Stderr = os.Stderr
	}
	
	if _exitCode, _data, _error := processExecuteGetStdout (_command); _error == nil {
		if _exitCode == 0 {
			_fingerprint := NewFingerprinter () .Bytes (_data) .Build ()
			return _fingerprint, _data, nil
		} else {
			return "", nil, errorf (0x42669a76, "command failed with exit code `%d`", _exitCode)
		}
	} else {
		return "", nil, _error
	}
}

