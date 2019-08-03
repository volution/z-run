

package main


import "crypto/sha256"
import "encoding/binary"
import "encoding/hex"
import "encoding/json"
import "errors"
import "fmt"
import "hash"
import "io"
import "io/ioutil"
import "log"
import "os"
import "path"
import "regexp"
import "sort"
import "strings"
import "unicode"
import "unicode/utf8"




type Scriptlet struct {
	Index uint `json:"id"`
	Label string `json:"label"`
	Interpreter string `json:"interpreter"`
	Body string `json:"body"`
	Fingerprint string `json:"fingerprint"`
	Source ScriptletSource `json:"source"`
}

type ScriptletSource struct {
	Path string `json:"path"`
	LineStart uint `json:"line_start"`
	LineEnd uint `json:"line_end"`
}


type Library struct {
	Scriptlets []*Scriptlet `json:"scriptlets"`
	ScriptletsByLabel map[string]uint `json:"scriptlets_by_label"`
	ScriptletsByFingerprint map[string]uint `json:"scriptlets_by_fingerprint"`
	ScriptletLabels []string `json:"scriptlet_labels"`
	Sources []*Source `json:"sources"`
}


type Source struct {
	Path string `json:"path"`
	Executable bool `json:"executable"`
	FingerprintMeta string `json:"fingerprint_meta"`
	FingerprintData string `json:"fingerprint_data"`
}




func NewLibrary () (*Library) {
	return & Library {
			Scriptlets : make ([]*Scriptlet, 0, 1024),
			ScriptletsByLabel : make (map[string]uint, 1024),
			ScriptletsByFingerprint : make (map[string]uint, 1024),
			ScriptletLabels : make ([]string, 0, 1024),
		}
}


func (_library *Library) Labels () ([]string, error) {
	return _library.ScriptletLabels, nil
}


func (_library *Library) Resolve (_label string) (*Scriptlet, error) {
	if _index, _exists := _library.ScriptletsByLabel[_label]; _exists {
		_scriptlet := _library.Scriptlets[_index]
		return _scriptlet, nil
	} else {
		return nil, nil
	}
}




func parseFromSource (_library *Library, _source *Source) (string, error) {
	if _source.Executable {
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
		if ! utf8.Valid (_data) {
			return "", errorf (0x2a19cfc7, "invalid UTF-8 source")
		}
		_data := string (_data)
		if _error := parseFromData (_library, _data, _sourcePath); _error == nil {
			_fingerprint := NewFingerprinter () .String (_data) .Build ()
			return _fingerprint, _error
		} else {
			return "", _error
		}
	} else {
		return "", _error
	}
}


func parseFromData (_library *Library, _source string, _sourcePath string) (error) {
	
	
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
						return errorf (0x53eafa1a, "invalid syntax (%d):  missing scriptlet separator `::`", _lineIndex, _line)
					}
					_label = _trimSpace (_label)
					_body = _trimSpace (_body)
					
					if _label == "" {
						return errorf (0xddec2340, "invalid syntax (%d):  empty scriptlet label", _lineIndex, _line)
					}
					if _body == "" {
						return errorf (0xc1dc94cc, "invalid syntax (%d):  empty scriptlet body", _lineIndex, _line)
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
						return errorf (0x64c17a76, "invalid syntax (%d):  empty scriptlet label", _lineIndex, _line)
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
					return errorf (0x183de0fd, "invalid syntax (%d):  unexpected statement `%s`", _lineIndex, _line)
				}
			
			case SCRIPTLET_BODY :
				_lineTrimmed := _trimSpace (_line)
				
				if _lineTrimmed == "!!" {
					_scriptletState.body = _scriptletState.bodyBuffer.String ()
					_scriptletState.lineEnd = _lineIndex
					_state = SCRIPTLET_PUSH
					
				} else if strings.HasPrefix (_line, "!!") {
					return errorf (0xf9900c0c, "invalid syntax (%d):  unexpected statement `%s`", _lineIndex, _line)
					
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
					return errorf (0x183de0fd, "invalid syntax (%d):  unexpected statement `%s`", _lineIndex, _line)
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
				return _error
			}
			_state = WAITING
		}
	}
	
	switch _state {
		case WAITING :
		case SCRIPTLET_BODY :
			return errorf (0x9d55df33, "invalid syntax (%d):  missing scriptlet body closing tag `!!` (and reached end of file)", _lineIndex)
		case SKIPPING :
			return errorf (0x357f15e1, "invalid syntax (%d):  missing comment body closing tag `##!!` (and reached end of file)", _lineIndex)
		default :
			return errorf (0xc0f78380, "invalid syntax (%d):  unexpected state `%s` (and reached end of file)", _lineIndex, _state)
	}
	
	return nil
}




func includeScriptlet (_library *Library, _scriptlet *Scriptlet) (error) {
	
	if _scriptlet.Label != strings.TrimSpace (_scriptlet.Label) {
		return errorf (0xd8797e9e, "invalid scriptlet label `%s`", _scriptlet.Label)
	}
	if _scriptlet.Label == "" {
		return errorf (0xaede3d8c, "invalid scriptlet label `%s`", _scriptlet.Label)
	}
	if _, _exists := _library.ScriptletsByLabel[_scriptlet.Label]; _exists {
		return errorf (0x883f9a7f, "duplicate scriptlet label `%s`", _scriptlet.Label)
	}
	
	if _scriptlet.Interpreter == "" {
		_scriptlet.Interpreter = "<shell>"
	}
	
	_fingerprint := NewFingerprinter () .StringWithLen (_scriptlet.Label) .StringWithLen (_scriptlet.Interpreter) .StringWithLen (_scriptlet.Body) .Build ()
	
	if _, _exists := _library.ScriptletsByFingerprint[_fingerprint]; _exists {
		return nil
	}
	
	_scriptlet.Index = uint (len (_library.Scriptlets))
	_scriptlet.Fingerprint = _fingerprint
	
	_library.Scriptlets = append (_library.Scriptlets, _scriptlet)
	_library.ScriptletsByLabel[_scriptlet.Label] = _scriptlet.Index
	_library.ScriptletsByFingerprint[_scriptlet.Fingerprint] = _scriptlet.Index
	_library.ScriptletLabels = append (_library.ScriptletLabels, _scriptlet.Label)
	
	return nil
}




func resolveSources (_candidate string) ([]*Source, error) {
	
	_sources := make ([]*Source, 0, 128)
	
	_candidate, _stat, _error := resolveSourcesPath_0 (_candidate)
	if _error != nil {
		return nil, _error
	}
	
	_statMode := _stat.Mode ()
	switch {
		
		case _statMode.IsRegular () :
			_fingerprint := NewFingerprinter () .StringWithLen (_candidate) .Int64 (_stat.Size ()) .Int64 (_stat.ModTime () .Unix ()) .Build ()
			_source := & Source {
					Path : _candidate,
					Executable : (_statMode.Perm () & 0111) != 0,
					FingerprintMeta : _fingerprint,
				}
			_sources = append (_sources, _source)
		
		case _statMode.IsDir () :
			return nil, errorf (0x8a04b23b, "not-implemented")
		
		default :
			return nil, errorf (0xa35428a2, "invalid source `%s`", _candidate)
	}
	
	return _sources, nil
}


func resolveSourcesPath_0 (_candidate string) (string, os.FileInfo, error) {
	if _candidate != "" {
		return resolveSourcesPath_2 (_candidate)
	} else {
		return resolveSourcesPath_1 ()
	}
}


func resolveSourcesPath_1 () (string, os.FileInfo, error) {
	
	_folders := make ([]string, 0, 128)
	_folders = append (_folders,
			".",
			path.Join (".", ".git"),
			path.Join (".", ".hg"),
			path.Join (".", ".svn"),
		)
	for _, _folder := range _folders {
		_folders = append (_folders, path.Join (_folder, "scripts"))
	}
	
	_files := []string {
			"x-run",
			".x-run",
			"_x-run",
		}
	
	_candidates := make ([]string, 0, 16)
	
	for _, _folder := range _folders {
		for _, _file := range _files {
			_path := path.Join (_folder, _file)
			if _, _error := os.Lstat (_path); _error == nil {
				_candidates = append (_candidates, _path)
			} else if os.IsNotExist (_error) {
				// NOP
			} else {
				return "", nil, _error
			}
		}
	}
	
	if len (_candidates) == 0 {
		return "", nil, errorf (0x779f9a9f, "no sources found")
	} else if len (_candidates) > 1 {
		return "", nil, errorf (0x519bb041, "too many sources found: `%s`", _candidates)
	} else {
		return resolveSourcesPath_2 (_candidates[0])
	}
}


func resolveSourcesPath_2 (_candidate string) (string, os.FileInfo, error) {
	if _stat, _error := os.Stat (_candidate); _error == nil {
		return _candidate, _stat, nil
	} else if os.IsNotExist (_error) {
		return "", nil, errorf (0x4b0005de, "source does not exist `%s`", _candidate)
	} else {
		return "", nil, _error
	}
}




func loadLibrary (_candidate string) (*Library, error) {
	
	_sources, _error := resolveSources (_candidate)
	if _error != nil {
		return nil, _error
	}
	
	_library := NewLibrary ()
	_library.Sources = _sources
	
	for _, _source := range _sources {
		if _fingerprint, _error := parseFromSource (_library, _source); _error == nil {
			_source.FingerprintData = _fingerprint
		} else {
			return nil, _error
		}
	}
	
	sort.Strings (_library.ScriptletLabels)
	
	return _library, nil
}




func doExecute (_library *Library, _executable string, _scriptlet string, _arguments []string, _environment map[string]string) (error) {
	return errorf (0x4f41e9bd, "not-implemented")
}


func doExportLabelsList (_library *Library, _stream io.Writer) (error) {
	if _labels, _error := _library.Labels (); _error == nil {
		for _, _label := range _labels {
			if _, _error := fmt.Fprintf (_stream, "%s\n", _label); _error != nil {
				return _error
			}
		}
		return nil
	} else {
		return _error
	}
}


func doExportScript (_library *Library, _label string, _stream io.Writer) (error) {
	if _scriptlet, _error := _library.Resolve (_label); _error == nil {
		if _scriptlet != nil {
			_, _error := io.WriteString (_stream, _scriptlet.Body)
			return _error
		} else {
			return errorf (0x95e0b174, "undefined scriptlet `%s`", _label)
		}
	} else {
		return _error
	}
}


func doExportLibraryJson (_library *Library, _stream io.Writer) (error) {
	_encoder := json.NewEncoder (_stream)
	_encoder.SetIndent ("", "    ")
	_encoder.SetEscapeHTML (false)
	return _encoder.Encode (_library)
}




func main_0 (_executable string, _argument0 string, _arguments []string, _environment map[string]string) (error) {
	
	var _cachePath string
	var _sourcePath string
	var _command string
	var _scriptlet string
	
	var _cleanArguments []string
	var _cleanEnvironment map[string]string = make (map[string]string, len (_environment))
	
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
			_nameCanonical = strings.Replace (_nameCanonical, "X_RUN", "XRUN", 1)
		}
		
		if strings.HasPrefix (_nameCanonical, "XRUN") || strings.HasPrefix (_nameCanonical, "_XRUN") {
			
			if _name != _nameCanonical {
				logf ('w', 0x37850eb3, "environment variable does not have canonical name;  expected `%s`, encountered `%s`!", _nameCanonical, _name)
			}
			
			switch _nameCanonical {
				case "XRUN_SOURCE" :
					_sourcePath = _value
				case "XRUN_LIBRARY" :
					_cachePath = _value
				default :
					logf ('w', 0xdf61b057, "environment variable unknown: `%s`", _nameCanonical)
			}
			
		} else {
			_cleanEnvironment[_name] = _value
		}
	}
	
	for _index, _argument := range _arguments {
		
		if _argument == "--" {
			_cleanArguments = _arguments[_index + 1:]
			break
			
		} else if strings.HasPrefix (_argument, "-") {
			if _command != "" {
				return errorf (0xae04b5ff, "unexpected argument `%s`", _argument)
			}
			if strings.HasPrefix (_argument, "--source=") {
				_sourcePath = _argument[len ("--source="):]
			} else if strings.HasPrefix (_argument, "--library=") {
				_cachePath = _argument[len ("--library="):]
			} else {
				return errorf (0x33555ffb, "invalid argument `%s`", _argument)
			}
			
		} else if strings.HasPrefix (_argument, "::") {
			if _command == "" {
				_command = "execute"
			}
			_scriptlet = _argument
			_cleanArguments = _arguments[_index + 1:]
			break
			
		} else {
			if _command == "" {
				switch _argument {
					
					case "execute" :
						_command = "execute"
						continue
					
					case "export-script" :
						_command = "export-script"
						continue
					
					case "export-labels-list", "export-labels", "list" :
						_command = "export-labels-list"
						_index += 1
					
					case "export-library-json", "export-library", "export" :
						_command = "export-library-json"
						_index += 1
				}
			} else {
				_scriptlet = _argument
				_index += 1
			}
			_cleanArguments = _arguments[_index:]
			break
		}
	}
	
	if _cachePath != "" {
		if _sourcePath != "" {
			logf ('w', 0x1fe0b572, "both library and source path specified;  using library!")
			_sourcePath = ""
		}
	}
	
	if (_command == "") && (_scriptlet == "") {
		_command = "export-labels-list"
	}
	
	_library, _error := loadLibrary (_sourcePath)
	if _error != nil {
		return _error
	}
	
	switch _command {
		
		case "execute" :
			if _scriptlet == "" {
				return errorf (0x39718e70, "execute:  expected scriptlet")
			}
			return doExecute (_library, _executable, _scriptlet, _cleanArguments, _cleanEnvironment)
		
		case "export-script" :
			if _scriptlet == "" {
				return errorf (0xf24640a2, "export:  expected scriptlet")
			}
			if len (_cleanArguments) != 0 {
				return errorf (0xcf8db3c0, "export:  unexpected arguments")
			}
			return doExportScript (_library, _scriptlet, os.Stdout)
		
		case "export-labels-list" :
			if (_scriptlet != "") || (len (_cleanArguments) != 0) {
				return errorf (0xf7b9c7f3, "list:  unexpected scriptlet or arguments")
			}
			return doExportLabelsList (_library, os.Stdout)
		
		case "export-library-json" :
			if (_scriptlet != "") || (len (_cleanArguments) != 0) {
				return errorf (0x400ec122, "export:  unexpected scriptlet or arguments")
			}
			return doExportLibraryJson (_library, os.Stdout)
		
		case "" :
			return errorf (0x5d2a4326, "expected command")
		
		default :
			return errorf (0x66cf8700, "unexpected command `%s`", _command)
	}
}




func main () () {
	
	log.SetFlags (0)
	
	var _executable string
	if _executable_0, _error := os.Executable (); _error == nil {
		_executable = _executable_0
	} else {
		panic (abortError (_error))
	}
	
	_argument0 := os.Args[0]
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
	
	if _error := main_0 (_executable, _argument0, _arguments, _environment); _error == nil {
		os.Exit (0)
	} else {
		panic (abortError (_error))
	}
}




func logf (_slug rune, _code uint32, _format string, _arguments ... interface{}) () {
	_message := fmt.Sprintf (_format, _arguments ...)
	_prefix := fmt.Sprintf ("[%c%c]  [%08x]  ", _slug, _slug, _code)
	log.Print (_prefix + _message + "\n")
}

func logError (_slug rune, _error error) () {
	logErrorf (_slug, 0x55d59c80, _error, "unexpected error encountered!")
}

func logErrorf (_slug rune, _code uint32, _error error, _format string, _arguments ... interface{}) () {
	logf (_slug, _code, _format, _arguments ...)
	if _error != nil {
		_errorString := _error.Error ()
		_errorRegexp := regexp.MustCompile (`^\[[0-9a-f]{8}\]  [^\n]+$`)
		if _matches := _errorRegexp.MatchString (_errorString); _matches {
			log.Printf ("[%c%c]  %s\n", _slug, _slug, _errorString)
		} else {
			log.Printf ("[%c%c]  [%08x]  %s\n", _slug, _slug, 0xda900de1, _errorString)
			log.Printf ("[%c%c]  [%08x]  %#v\n", _slug, _slug, 0x4fb5d56d, _error)
		}
	}
}


func abortf (_code uint32, _format string, _arguments ... interface{}) (error) {
	return abortErrorf (nil, _code, _format, _arguments ...)
}

func abortError (_error error) (error) {
	return abortErrorf (_error, 0xe6ed2b0f, "unexpected error encountered!")
}

func abortErrorf (_error error, _code uint32, _format string, _arguments ... interface{}) (error) {
	logErrorf ('!', _code, _error, _format, _arguments ...)
	logf ('!', 0xb7a5fb86, "aborting!")
	os.Exit (1)
	panic (0xa235deea)
}


func errorf (_code uint32, _format string, _arguments ... interface{}) (error) {
	_message := fmt.Sprintf (_format, _arguments ...)
	_prefix := fmt.Sprintf ("[%08x]  ", _code)
	return errors.New (_prefix + _message)
}




type Fingerprinter struct {
	hasher hash.Hash
}

func NewFingerprinter () (Fingerprinter) {
	return Fingerprinter {
			hasher : sha256.New (),
		}
}

func (_fingerprinter Fingerprinter) Uint64 (_value uint64) (Fingerprinter) {
	var _bytes [8]byte
	binary.BigEndian.PutUint64 (_bytes[:], _value)
	_fingerprinter.hasher.Write (_bytes[:])
	return _fingerprinter
}
func (_fingerprinter Fingerprinter) Int64 (_value int64) (Fingerprinter) {
	return _fingerprinter.Uint64 (uint64 (_value))
}

func (_fingerprinter Fingerprinter) String (_value string) (Fingerprinter) {
	io.WriteString (_fingerprinter.hasher, _value)
	return _fingerprinter
}
func (_fingerprinter Fingerprinter) StringWithLen (_value string) (Fingerprinter) {
	_fingerprinter.Uint64 (uint64 (len (_value)))
	io.WriteString (_fingerprinter.hasher, _value)
	return _fingerprinter
}

func (_fingerprinter Fingerprinter) Bytes (_value []byte) (Fingerprinter) {
	_fingerprinter.hasher.Write (_value)
	return _fingerprinter
}
func (_fingerprinter Fingerprinter) BytesWithLen (_value []byte) (Fingerprinter) {
	_fingerprinter.Uint64 (uint64 (len (_value)))
	_fingerprinter.hasher.Write (_value)
	return _fingerprinter
}

func (_fingerprinter Fingerprinter) Build () (string) {
	return hex.EncodeToString (_fingerprinter.hasher.Sum (nil))
}

