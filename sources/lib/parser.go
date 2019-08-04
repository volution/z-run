

package lib


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




func parseLibrary (_sources []*Source, _sourcesFingerprint string, _context *Context) (*Library, error) {
	
	_library := NewLibrary ()
	_library.Sources = _sources
	_library.SourcesFingerprint = _sourcesFingerprint
	
	for _, _source := range _sources {
		if _fingerprint, _error := parseFromSource (_library, _source, _context); _error == nil {
			_source.FingerprintData = _fingerprint
		} else {
			return nil, _error
		}
	}
	
	sort.Strings (_library.ScriptletFingerprints)
	sort.Strings (_library.ScriptletLabels)
	
	return _library, nil
}



func parseFromSource (_library *Library, _source *Source, _context *Context) (string, error) {
	if _source.Executable {
		
		_executable := _source.Path
		if ! path.IsAbs (_executable) {
			if _executable_0, _error := filepath.Abs (_executable); _error == nil {
				_executable = _executable_0
			} else {
				return "", _error
			}
		}
		
		_command := & exec.Cmd {
				Path : _executable,
				Args : []string {
						"[x-run:generator]",
					},
				Env : commandEnvironment (_context, nil),
				Stdin : nil,
				Stdout : nil,
				Stderr : os.Stderr,
				Dir : "",
			}
		
		if _exitCode, _data, _error := commandExecuteGetStdout (_command); _error == nil {
			if _exitCode == 0 {
				return parseFromData (_library, string (_data), _source.Path)
			} else {
				return "", errorf (0x5fd5c9b7, "generated failed with exit code `%d`", _exitCode)
			}
		} else {
			return "", _error
		}
		
		return "", errorf (0x566095fc, "not-implemented")
		
	} else {
		return parseFromFile (_library, _source.Path)
	}
}


func parseFromFile (_library *Library, _path string) (string, error) {
	if _stream, _error := os.Open (_path); _error == nil {
		defer _stream.Close ()
		return parseFromStream (_library, _stream, _path)
	} else {
		return "", _error
	}
}


func parseFromStream (_library *Library, _stream io.Reader, _sourcePath string) (string, error) {
	if _data, _error := ioutil.ReadAll (_stream); _error == nil {
		if utf8.Valid (_data) {
			return parseFromData (_library, string (_data), _sourcePath)
		} else {
			return "", errorf (0x2a19cfc7, "invalid UTF-8 source")
		}
	} else {
		return "", _error
	}
}


func parseFromData (_library *Library, _source string, _sourcePath string) (string, error) {
	
	const (
		WAITING = 1 + iota
		SCRIPTLET_BODY
		SCRIPTLET_PUSH
		SKIPPING
	)
	
	type scriptletState struct {
		label string
		body string
		bodyBuffer strings.Builder
		bodyStrip string
		bodyLines uint
		lineStart uint
		lineEnd uint
	}
	
	_fingerprint := NewFingerprinter () .String (_source) .Build ()
	
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
				_line = _trimRightSpace (_line)
				_lineTrimmed := _trimSpace (_line)
				
				if _lineTrimmed == "" {
					// NOP
					
				} else if strings.HasPrefix (_line, ":: ") {
					
					_text := _line[3:]
					var _label string
					var _body string
					if _splitIndex := strings.Index (_text, " :: "); _splitIndex >= 0 {
						_label = _text[:_splitIndex]
						_body = _text[_splitIndex + 4:]
					} else {
						return "", errorf (0x53eafa1a, "invalid syntax (%d):  missing scriptlet separator `::`", _lineIndex, _line)
					}
					_label = _trimSpace (_label)
					_body = _trimSpace (_body)
					
					if _label == "" {
						return "", errorf (0xddec2340, "invalid syntax (%d):  empty scriptlet label", _lineIndex, _line)
					}
					if _body == "" {
						return "", errorf (0xc1dc94cc, "invalid syntax (%d):  empty scriptlet body", _lineIndex, _line)
					}
					
					_scriptletState = scriptletState {
							label : _label,
							body : _body + "\n",
							lineStart : _lineIndex,
							lineEnd : _lineIndex,
						}
					_state = SCRIPTLET_PUSH
					
				} else if strings.HasPrefix (_line, "<< ") {
					
					_label := _line[3:]
					_label = _trimSpace (_label)
					if _label == "" {
						return "", errorf (0x64c17a76, "invalid syntax (%d):  empty scriptlet label", _lineIndex, _line)
					}
					
					_scriptletState = scriptletState {
							label : _label,
							lineStart : _lineIndex,
						}
					_state = SCRIPTLET_BODY
					
				} else if strings.HasPrefix (_line, "##<< ") || (_lineTrimmed == "##<<") {
					_state = SKIPPING
					
				} else if strings.HasPrefix (_line, "#:: ") {
					// NOP
					
				} else if strings.HasPrefix (_line, "# ") || (_lineTrimmed == "#") {
					// NOP
					
				} else if (_lineIndex == 1) && strings.HasPrefix (_line, "#!/") {
					// NOP
					
				} else if false ||
						(_line == "##== sort = false") ||
						(_line == "##== sort = true") ||
						false {
					// NOP
					
				} else {
					return "", errorf (0x9f8daae4, "invalid syntax (%d):  unexpected statement `%s`", _lineIndex, _line)
				}
			
			case SCRIPTLET_BODY :
				_lineTrimmed := _trimSpace (_line)
				
				if _lineTrimmed == "!!" {
					_scriptletState.body = _scriptletState.bodyBuffer.String ()
					_scriptletState.lineEnd = _lineIndex
					_state = SCRIPTLET_PUSH
					
				} else if strings.HasPrefix (_line, "!!") {
					return "", errorf (0xf9900c0c, "invalid syntax (%d):  unexpected statement `%s`", _lineIndex, _line)
					
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
				_line = _trimRightSpace (_line)
				_lineTrimmed := _trimSpace (_line)
				if _lineTrimmed == "##!!" {
					_state = WAITING
				} else if strings.HasPrefix (_line, "##!!") {
					return "", errorf (0x183de0fd, "invalid syntax (%d):  unexpected statement `%s`", _lineIndex, _line)
				} else {
					// NOP
				}
		}
		
		if _state == SCRIPTLET_PUSH {
			_scriptlet := & Scriptlet {
					Label : _scriptletState.label,
					Body : _scriptletState.body,
					Source : ScriptletSource {
							Path : _sourcePath,
							LineStart : _scriptletState.lineStart,
							LineEnd : _scriptletState.lineEnd,
						},
				}
			if _error := includeScriptlet (_library, _scriptlet); _error != nil {
				return "", _error
			}
			_state = WAITING
		}
	}
	
	switch _state {
		case WAITING :
		case SCRIPTLET_BODY :
			return "", errorf (0x9d55df33, "invalid syntax (%d):  missing scriptlet body closing tag `!!` (and reached end of file)", _lineIndex)
		case SKIPPING :
			return "", errorf (0x357f15e1, "invalid syntax (%d):  missing comment body closing tag `##!!` (and reached end of file)", _lineIndex)
		default :
			return "", errorf (0xc0f78380, "invalid syntax (%d):  unexpected state `%s` (and reached end of file)", _lineIndex, _state)
	}
	
	return _fingerprint, nil
}

