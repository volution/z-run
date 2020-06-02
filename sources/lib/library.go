

package zrun


import "strings"




type Scriptlet struct {
	Index uint `json:"id"`
	Label string `json:"label"`
	Kind string `json:"kind"`
	Interpreter string `json:"interpreter"`
	InterpreterExecutable string `json:"interpreter-executable,omitempty"`
	InterpreterArguments []string `json:"interpreter-arguments,omitempty"`
	InterpreterEnvironment map[string]string `json:"interpreter-environment,omitempty"`
	Context *ScriptletContext `json:"-"`
	ContextFingerprint string `json:"context,omitempty"`
	Body string `json:"body,omitempty"`
	Fingerprint string `json:"fingerprint"`
	Source ScriptletSource `json:"source"`
	Visible bool `json:"visible"`
	Hidden bool `json:"hidden"`
	Menus []string `json:"menus"`
}

type ScriptletSource struct {
	Path string `json:"path"`
	LineStart uint `json:"line_start"`
	LineEnd uint `json:"line_end"`
}

type ScriptletContext struct {
	Fingerprint string `json:"fingerprint"`
	ExecutablePaths []string `json:"executable-paths"`
	Environment map[string]string `json:"environment"`
}




type Library struct {
	
	Scriptlets LibraryScriptlets `json:"scriptlets"`
	
	ScriptletFingerprints []string `json:"fingerprints"`
	ScriptletsByFingerprint map[string]uint `json:"index_by_fingerprint"`
	
	ScriptletLabels []string `json:"labels"`
	ScriptletLabelsAll []string `json:"labels"`
	ScriptletsByLabel map[string]uint `json:"index_by_label"`
	
	ScriptletsContexts map[string]*ScriptletContext `json:"scriptlets-contexts"`
	
	Sources LibrarySources `json:"sources"`
	
	SourcesFingerprint string `json:"sources_fingerprint"`
	EnvironmentFingerprint string `json:"environment_fingerprint"`
	LibraryFingerprint string `json:"library_fingerprint"`
	
	url string
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
			ScriptletFingerprints : make ([]string, 0, 1024),
			ScriptletsByFingerprint : make (map[string]uint, 1024),
			ScriptletLabels : make ([]string, 0, 1024),
			ScriptletLabelsAll : make ([]string, 0, 1024),
			ScriptletsByLabel : make (map[string]uint, 1024),
			ScriptletsContexts : make (map[string]*ScriptletContext, 16),
		}
}


func (_library *Library) SelectFingerprints () ([]string, *Error) {
	return _library.ScriptletFingerprints, nil
}

func (_library *Library) SelectLabels () ([]string, *Error) {
	return _library.ScriptletLabels, nil
}

func (_library *Library) SelectLabelsAll () ([]string, *Error) {
	return _library.ScriptletLabelsAll, nil
}


func (_library *Library) ResolveFullByFingerprint (_fingerprint string) (*Scriptlet, *Error) {
	if _index, _exists := _library.ScriptletsByFingerprint[_fingerprint]; _exists {
		_scriptlet := _library.Scriptlets[_index]
		return _scriptlet, nil
	} else {
		return nil, nil
	}
}

func (_library *Library) ResolveMetaByFingerprint (_fingerprint string) (*Scriptlet, *Error) {
	if _scriptlet, _error := _library.ResolveFullByFingerprint (_fingerprint); (_error == nil) && (_scriptlet != nil) {
		_meta := & Scriptlet {}
		*_meta = *_scriptlet
		_meta.Body = ""
		return _meta, nil
	} else {
		return nil, _error
	}
}

func (_library *Library) ResolveBodyByFingerprint (_fingerprint string) (string, bool, *Error) {
	if _scriptlet, _error := _library.ResolveFullByFingerprint (_fingerprint); (_error == nil) && (_scriptlet != nil) {
		return _scriptlet.Body, true, nil
	} else {
		return "", false, _error
	}
}


func (_library *Library) ResolveFullByLabel (_label string) (*Scriptlet, *Error) {
	if _index, _exists := _library.ScriptletsByLabel[_label]; _exists {
		_scriptlet := _library.Scriptlets[_index]
		return _scriptlet, nil
	} else {
		return nil, nil
	}
}

func (_library *Library) ResolveMetaByLabel (_label string) (*Scriptlet, *Error) {
	if _scriptlet, _error := _library.ResolveFullByLabel (_label); (_error == nil) && (_scriptlet != nil) {
		_meta := & Scriptlet {}
		*_meta = *_scriptlet
		_meta.Body = ""
		return _meta, nil
	} else {
		return nil, _error
	}
}

func (_library *Library) ResolveBodyByLabel (_label string) (string, bool, *Error) {
	if _scriptlet, _error := _library.ResolveFullByLabel (_label); (_error == nil) && (_scriptlet != nil) {
		return _scriptlet.Body, true, nil
	} else {
		return "", false, _error
	}
}

func (_library *Library) ResolveFingerprintByLabel (_label string) (string, bool, *Error) {
	if _scriptlet, _error := _library.ResolveFullByLabel (_label); (_error == nil) && (_scriptlet != nil) {
		return _scriptlet.Fingerprint, true, nil
	} else {
		return "", false, _error
	}
}


func (_library *Library) ResolveContextByFingerprint (_fingerprint string) (*ScriptletContext, bool, *Error) {
	if _context, _exists := _library.ScriptletsContexts[_fingerprint]; _exists {
		return _context, true, nil
	} else {
		return nil, false, errorf (0x30d90869, "invalid scriptlet context fingerprint `%s`", _fingerprint)
	}
}


func (_library *Library) SelectSources () (LibrarySources, *Error) {
	return _library.Sources, nil
}


func (_library *Library) Url () (string) {
	return _library.url
}

func (_library *Library) Fingerprint () (string, *Error) {
	_fingerprint := _library.LibraryFingerprint
	if _fingerprint != "" {
		return _fingerprint, nil
	} else {
		return "", errorf (0x7c26bcc2, "invalid state")
	}
}


func (_library *Library) Close () (*Error) {
	*_library = Library {}
	return nil
}




func includeScriptlet (_library *Library, _scriptlet *Scriptlet) (*Error) {
	
	if _scriptlet.Label != strings.TrimSpace (_scriptlet.Label) {
		return errorf (0xd8797e9e, "invalid scriptlet label `%s`", _scriptlet.Label)
	}
	if _scriptlet.Label == "" {
		return errorf (0xaede3d8c, "invalid scriptlet label `%s`", _scriptlet.Label)
	}
	if _, _exists := _library.ScriptletsByLabel[_scriptlet.Label]; _exists {
		return errorf (0x883f9a7f, "duplicate scriptlet label `%s`", _scriptlet.Label)
	}
	
	switch _scriptlet.Interpreter {
		case "<detect>", "<print>", "<menu>" :
			// NOP
		default :
			return errorf (0xbf289098, "invalid scriptlet interpreter `%s`", _scriptlet.Interpreter)
	}
	if (_scriptlet.InterpreterExecutable != "") || (_scriptlet.InterpreterArguments != nil) || (_scriptlet.InterpreterEnvironment != nil) {
		return errorf (0x901675e8, "invalid scriptlet interpreter state")
	}
	
	switch _scriptlet.Kind {
		case "executable" :
			_scriptlet.Kind = "executable-pending"
		case "generator" :
			_scriptlet.Kind = "generator-pending"
		case "script-replacer" :
			_scriptlet.Kind = "script-replacer-pending"
		case "print-replacer" :
			_scriptlet.Kind = "print-replacer-pending"
		case "menu" :
			_scriptlet.Kind = "menu-pending"
		default :
			return errorf (0x4b8aacf2, "invalid scriptlet kind `%s`", _scriptlet.Kind)
	}
	
	_fingerprint := NewFingerprinter () .StringWithLen (_scriptlet.Label) .StringWithLen (_scriptlet.Kind) .StringWithLen (_scriptlet.Interpreter) .StringWithLen (_scriptlet.Body) .Build ()
	
	if _, _exists := _library.ScriptletsByFingerprint[_fingerprint]; _exists {
		return nil
	}
	
	_scriptlet.Index = uint (len (_library.Scriptlets))
	_scriptlet.Fingerprint = _fingerprint
	
	_library.Scriptlets = append (_library.Scriptlets, _scriptlet)
	
	_library.ScriptletFingerprints = append (_library.ScriptletFingerprints, _scriptlet.Fingerprint)
	_library.ScriptletLabelsAll = append (_library.ScriptletLabelsAll, _scriptlet.Label)
	if !_scriptlet.Hidden || _scriptlet.Visible {
		_library.ScriptletLabels = append (_library.ScriptletLabels, _scriptlet.Label)
	}
	
	_library.ScriptletsByFingerprint[_scriptlet.Fingerprint] = _scriptlet.Index
	_library.ScriptletsByLabel[_scriptlet.Label] = _scriptlet.Index
	
	return nil
}




func includeSource (_library *Library, _source *Source) (*Error) {
	if _source.Path == "" {
		return errorf (0x12bdc134, "invalid state")
	}
	if _source.FingerprintMeta == "" {
		return errorf (0x152074de, "invalid state")
	}
//	if _source.FingerprintData == "" {
//		return errorf (0x401d0c16, "invalid state")
//	}
	for _, _existing := range _library.Sources {
		if _existing.Path == _source.Path {
			return errorf (0xf01b93ea, "invalid state")
		}
		if _existing.FingerprintMeta == _source.FingerprintMeta {
			return errorf (0x310f6193, "invalid state")
		}
		if (_existing.FingerprintData == _source.FingerprintData) && (_existing.FingerprintData != "") {
			return errorf (0x00fb18a1, "invalid state %#v %#v", _existing, _source)
		}
	}
	_library.Sources = append (_library.Sources, _source)
	return nil
}




type LibraryScriptlets []*Scriptlet

func (_scriptlets LibraryScriptlets) Len () (int) {
	return len (_scriptlets)
}
func (_scriptlets LibraryScriptlets) Less (_left int, _right int) (bool) {
	return (_scriptlets[_left].Label < _scriptlets[_right].Label)
}
func (_scriptlets LibraryScriptlets) Swap (_left int, _right int) () {
	_scriptlets[_left], _scriptlets[_right] = _scriptlets[_right], _scriptlets[_left]
}


type LibrarySources []*Source

func (_sources LibrarySources) Len () (int) {
	return len (_sources)
}
func (_sources LibrarySources) Less (_left int, _right int) (bool) {
	return (_sources[_left].Path < _sources[_right].Path)
}
func (_sources LibrarySources) Swap (_left int, _right int) () {
	_sources[_left], _sources[_right] = _sources[_right], _sources[_left]
}

