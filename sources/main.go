

package main


import "bytes"
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


import cdb "github.com/cipriancraciun/go-cdb-lib"




type Scriptlet struct {
	Index uint `json:"id"`
	Label string `json:"label"`
	Interpreter string `json:"interpreter"`
	Body string `json:"body,omitempty"`
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
	
	ScriptletFingerprints []string `json:"fingerprints"`
	ScriptletsByFingerprint map[string]uint `json:"index_by_fingerprint"`
	
	ScriptletLabels []string `json:"labels"`
	ScriptletsByLabel map[string]uint `json:"index_by_label"`
	
	Sources []*Source `json:"sources"`
	SourcesFingerprint string `json:"sources_fingerprint"`
}


type Source struct {
	Path string `json:"path"`
	Executable bool `json:"executable"`
	FingerprintMeta string `json:"fingerprint_meta"`
	FingerprintData string `json:"fingerprint_data"`
}


type LibraryStore interface {
	
	SelectFingerprints () ([]string, error)
	ResolveFullByFingerprint (_label string) (*Scriptlet, error)
	ResolveMetaByFingerprint (_label string) (*Scriptlet, error)
	ResolveBodyByFingerprint (_label string) (string, bool, error)
	
	SelectLabels () ([]string, error)
	ResolveFullByLabel (_label string) (*Scriptlet, error)
	ResolveMetaByLabel (_label string) (*Scriptlet, error)
	ResolveBodyByLabel (_label string) (string, bool, error)
	ResolveFingerprintByLabel (_label string) (string, bool, error)
}




type StoreOutput interface {
	Include (_namespace string, _key string, _value interface{}) (error)
	Commit () (error)
}

type StoreInput interface {
	Select (_namespace string, _key string, _value interface{}) (error)
	Close () (error)
}




func NewLibrary () (*Library) {
	return & Library {
			Scriptlets : make ([]*Scriptlet, 0, 1024),
			ScriptletFingerprints : make ([]string, 0, 1024),
			ScriptletsByFingerprint : make (map[string]uint, 1024),
			ScriptletLabels : make ([]string, 0, 1024),
			ScriptletsByLabel : make (map[string]uint, 1024),
		}
}

func (_library *Library) SelectFingerprints () ([]string, error) {
	return _library.ScriptletFingerprints, nil
}

func (_library *Library) SelectLabels () ([]string, error) {
	return _library.ScriptletLabels, nil
}


func (_library *Library) ResolveFullByFingerprint (_fingerprint string) (*Scriptlet, error) {
	if _index, _exists := _library.ScriptletsByFingerprint[_fingerprint]; _exists {
		_scriptlet := _library.Scriptlets[_index]
		return _scriptlet, nil
	} else {
		return nil, nil
	}
}

func (_library *Library) ResolveMetaByFingerprint (_fingerprint string) (*Scriptlet, error) {
	if _scriptlet, _error := _library.ResolveFullByFingerprint (_fingerprint); (_error == nil) && (_scriptlet != nil) {
		_meta := & Scriptlet {}
		*_meta = *_scriptlet
		_meta.Body = ""
		return _meta, nil
	} else {
		return nil, _error
	}
}

func (_library *Library) ResolveBodyByFingerprint (_fingerprint string) (string, bool, error) {
	if _scriptlet, _error := _library.ResolveFullByFingerprint (_fingerprint); (_error == nil) && (_scriptlet != nil) {
		return _scriptlet.Body, true, nil
	} else {
		return "", false, _error
	}
}


func (_library *Library) ResolveFullByLabel (_label string) (*Scriptlet, error) {
	if _index, _exists := _library.ScriptletsByLabel[_label]; _exists {
		_scriptlet := _library.Scriptlets[_index]
		return _scriptlet, nil
	} else {
		return nil, nil
	}
}

func (_library *Library) ResolveMetaByLabel (_label string) (*Scriptlet, error) {
	if _scriptlet, _error := _library.ResolveFullByLabel (_label); (_error == nil) && (_scriptlet != nil) {
		_meta := & Scriptlet {}
		*_meta = *_scriptlet
		_meta.Body = ""
		return _meta, nil
	} else {
		return nil, _error
	}
}

func (_library *Library) ResolveBodyByLabel (_label string) (string, bool, error) {
	if _scriptlet, _error := _library.ResolveFullByLabel (_label); (_error == nil) && (_scriptlet != nil) {
		return _scriptlet.Body, true, nil
	} else {
		return "", false, _error
	}
}

func (_library *Library) ResolveFingerprintByLabel (_label string) (string, bool, error) {
	if _scriptlet, _error := _library.ResolveFullByLabel (_label); (_error == nil) && (_scriptlet != nil) {
		return _scriptlet.Fingerprint, true, nil
	} else {
		return "", false, _error
	}
}




func parseLibrary (_sources []*Source, _sourcesFingerprint string) (*Library, error) {
	
	_library := NewLibrary ()
	_library.Sources = _sources
	_library.SourcesFingerprint = _sourcesFingerprint
	
	for _, _source := range _sources {
		if _fingerprint, _error := parseFromSource (_library, _source); _error == nil {
			_source.FingerprintData = _fingerprint
		} else {
			return nil, _error
		}
	}
	
	sort.Strings (_library.ScriptletFingerprints)
	sort.Strings (_library.ScriptletLabels)
	
	return _library, nil
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
	_library.ScriptletFingerprints = append (_library.ScriptletFingerprints, _scriptlet.Fingerprint)
	_library.ScriptletsByFingerprint[_scriptlet.Fingerprint] = _scriptlet.Index
	_library.ScriptletLabels = append (_library.ScriptletLabels, _scriptlet.Label)
	_library.ScriptletsByLabel[_scriptlet.Label] = _scriptlet.Index
	
	return nil
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


func resolveCache () (string, error) {
	var _cache string
	if _cache_0, _error := os.UserCacheDir (); _error == nil {
		_cache = _cache_0
	} else {
		return "", _error
	}
	if _error := os.MkdirAll (_cache, 0750); _error != nil {
		return "", _error
	}
	_cache = path.Join (_cache, "x-run")
	if _error := os.MkdirAll (_cache, 0750); _error != nil {
		return "", _error
	}
	return _cache, nil
}




func resolveLibrary (_candidate string) (LibraryStore, error) {
	
	_sources, _error := resolveSources (_candidate)
	if _error != nil {
		return nil, _error
	}
	
	var _sourcesFingerprint string
	{
		_fingerprints := make ([]string, 0, len (_sources) * 2)
		for _, _source := range _sources {
			_fingerprints = append (_fingerprints, _source.FingerprintMeta)
			if _source.FingerprintData != "" {
				_fingerprints = append (_fingerprints, _source.FingerprintData)
			}
		}
		sort.Strings (_fingerprints)
		_fingerprinter := NewFingerprinter ()
		for _, _fingerprint := range _fingerprints {
			_fingerprinter.StringWithLen (_fingerprint)
		}
		_sourcesFingerprint = _fingerprinter.Build ()
	}
	
	if _library, _error := parseLibrary (_sources, _sourcesFingerprint); _error == nil {
		return _library, _error
	} else {
		return nil, _error
	}
}




func doExecute (_library LibraryStore, _executable string, _scriptlet string, _arguments []string, _environment map[string]string) (error) {
	return errorf (0x4f41e9bd, "not-implemented")
}


func doExportLabelsList (_library LibraryStore, _stream io.Writer) (error) {
	if _labels, _error := _library.SelectLabels (); _error == nil {
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


func doExportScript (_library LibraryStore, _label string, _stream io.Writer) (error) {
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


func doExportLibraryJson (_library LibraryStore, _stream io.Writer) (error) {
	_encoder := json.NewEncoder (_stream)
	_encoder.SetIndent ("", "    ")
	_encoder.SetEscapeHTML (false)
	return _encoder.Encode (_library)
}


func doExportLibraryStore (_library LibraryStore, _store StoreOutput) (error) {
	
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
	
	for _, _fingerprint := range _fingerprintsFromStore {
		
		if _meta, _error := _library.ResolveMetaByFingerprint (_fingerprint); _error == nil {
			if _meta == nil {
				return errorf (0x20bc9d40, "invalid store")
			}
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
			_fingerprintsByLabels[_label] = _fingerprint
			_labels = append (_labels, _label)
			_labelsByFingerprints[_fingerprint] = _label
		}
		
		if _body, _found, _error := _library.ResolveBodyByFingerprint (_fingerprint); _error == nil {
			if !_found {
				return errorf (0xd80a265e, "invalid store")
			}
			if _error := _store.Include ("scriptlets-body", _fingerprint, _body); _error != nil {
				return _error
			}
		}
	}
	
	sort.Strings (_fingerprints)
	sort.Strings (_labels)
	
	if _error := _store.Include ("scriptlets-indices", "fingerprints", _fingerprints); _error != nil {
		return _error
	}
	if _error := _store.Include ("scriptlets-indices", "fingerprints-by-labels", _fingerprintsByLabels); _error != nil {
		return _error
	}
	if _error := _store.Include ("scriptlets-indices", "labels", _labels); _error != nil {
		return _error
	}
	if _error := _store.Include ("scriptlets-indices", "labels-by-fingerprints", _labelsByFingerprints); _error != nil {
		return _error
	}
	
	if _error := _store.Commit (); _error != nil {
		return _error
	}
	
	return nil
}


func doExportLibraryCdb (_library LibraryStore, _path string) (error) {
	if _store, _error := NewCdbStoreOutput (_path); _error == nil {
		return doExportLibraryStore (_library, _store)
	} else {
		return _error
	}
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
					
					case "compile-library-json", "compile-library" :
						_command = "compile-library-json"
						_index += 1
					
					case "compile-library-cdb" :
						_command = "compile-library-cdb"
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
	
	_library, _error := resolveLibrary (_sourcePath)
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
		
		case "export-library-json", "compile-library-json" :
			if (_scriptlet != "") || (len (_cleanArguments) != 0) {
				return errorf (0x400ec122, "export:  unexpected scriptlet or arguments")
			}
			switch _command {
				case "export-library-json" :
					return doExportLibraryJson (_library, os.Stdout)
				case "compile-library-json" :
					return doExportLibraryStore (_library, NewJsonStreamStoreOutput (os.Stdout, nil))
				default :
					panic (0xda7243ef)
			}
		
		case "compile-library-cdb" :
			if _scriptlet != "" {
				return errorf (0x492ac50e, "export:  unexpected scriptlet")
			}
			if len (_cleanArguments) != 1 {
				return errorf (0xf76f4459, "export:  expected database path")
			}
			return doExportLibraryCdb (_library, _cleanArguments[0])
		
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




type JsonStreamStoreOutput struct {
	stream io.Writer
	closer io.Closer
	encoder *json.Encoder
}

type JsonStoreRecord struct {
	Namespace string `json:"namespace"`
	Key string `json:"key"`
	Value interface{} `json:"value"`
}

func NewJsonStreamStoreOutput (_stream io.Writer, _closer io.Closer) (*JsonStreamStoreOutput) {
	_encoder := json.NewEncoder (_stream)
	_encoder.SetIndent ("", "    ")
	_encoder.SetEscapeHTML (false)
	return & JsonStreamStoreOutput {
			stream : _stream,
			closer : _closer,
			encoder : _encoder,
		}
}

func (_store *JsonStreamStoreOutput) Include (_namespace string, _key string, _value interface{}) (error) {
	_record := & JsonStoreRecord {
			Namespace : _namespace,
			Key : _key,
			Value : _value,
		}
	return _store.encoder.Encode (_record)
}

func (_store *JsonStreamStoreOutput) Commit () (error) {
	var _error error
	if _store.closer != nil {
		_error = _store.closer.Close ()
	}
	_store.stream = nil
	_store.closer = nil
	_store.encoder = nil
	return _error
}




type CdbStoreInput struct {
	reader *cdb.CDB
	path string
}


func NewCdbStoreInput (_path string) (*CdbStoreInput, error) {
	if _file, _error := os.Open (_path); _error == nil {
		defer _file.Close ()
		if _reader, _error := cdb.NewFromMappedWithHasher (_file, nil); _error == nil {
			_store := & CdbStoreInput {
					reader : _reader,
					path : _path,
				}
			return _store, nil
		} else {
			return nil, _error
		}
	} else {
		return nil, _error
	}
}


func (_store *CdbStoreInput) Select (_namespace string, _key string, _value interface{}) (error) {
	
	_keyBuffer := bytes.NewBuffer (nil)
	_keyBuffer.WriteString (_namespace)
	_keyBuffer.WriteString (" // ")
	_keyBuffer.WriteString (_key)
	
	var _valueData []byte
	if _valueData_0, _error := _store.reader.Get (_keyBuffer.Bytes ()); _error == nil {
		_valueData = _valueData_0
	} else {
		return _error
	}
	
	switch _value := _value.(type) {
		case []byte :
			return errorf (0x36de066a, "unexpected type")
		case *[]byte :
			if *_value == nil {
				*_value = _valueData
			} else {
				*_value = append ((*_value)[:0], _valueData ...)
			}
		case string :
			return errorf (0x36de066a, "unexpected type")
		case *string :
			*_value = string (_valueData)
		default :
			if _error := json.Unmarshal (_valueData, _value); _error != nil {
				return _error
			}
	}
	
	return nil
}


func (_store *CdbStoreInput) Close () (error) {
	if _error := _store.reader.Close (); _error == nil {
		_store.reader = nil
		return _error
	} else {
		return nil
	}
}




type CdbStoreOutput struct {
	writer *cdb.Writer
	pathFinal string
	pathTemporary string
}


func NewCdbStoreOutput (_path string) (*CdbStoreOutput, error) {
	if _path == "" {
		return nil, errorf (0x6917ab7d, "invalid path")
	}
	_pathFinal := _path
	_pathTemporary := fmt.Sprintf ("%s--%08x.tmp", _pathFinal, os.Getpid ())
	if _writer, _error := cdb.Create (_pathTemporary); _error == nil {
		_store := & CdbStoreOutput {
				writer : _writer,
				pathFinal : _pathFinal,
				pathTemporary : _pathTemporary,
			}
		return _store, nil
	} else {
		return nil, _error
	}
}


func (_store *CdbStoreOutput) Include (_namespace string, _key string, _value interface{}) (error) {
	
	_keyBuffer := bytes.NewBuffer (nil)
	_keyBuffer.WriteString (_namespace)
	_keyBuffer.WriteString (" // ")
	_keyBuffer.WriteString (_key)
	
	_valueBuffer := bytes.NewBuffer (nil)
	switch _value := _value.(type) {
		case []byte :
			_valueBuffer.Write (_value)
		case *[]byte :
			_valueBuffer.Write (*_value)
		case string :
			_valueBuffer.WriteString (_value)
		case *string :
			_valueBuffer.WriteString (*_value)
		default :
			if _error := json.NewEncoder (_valueBuffer) .Encode (_value); _error != nil {
				return _error
			}
	}
	
	return _store.writer.Put (_keyBuffer.Bytes (), _valueBuffer.Bytes ())
}

func (_store *CdbStoreOutput) Commit () (error) {
	if _error := _store.writer.Close (); _error != nil {
		return _error
	}
	_store.writer = nil
	if _error := os.Rename (_store.pathTemporary, _store.pathFinal); _error != nil {
		return _error
	}
	return nil
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

