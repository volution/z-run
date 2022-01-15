

package zrun


import "strings"

import . "github.com/cipriancraciun/z-run/lib/common"




type Scriptlet struct {
	
	Label string `json:"label"`
	Kind string `json:"kind"`
	Interpreter string `json:"interpreter"`
	InterpreterExecutable string `json:"interpreter_executable,omitempty"`
	InterpreterArguments []string `json:"interpreter_arguments,omitempty"`
	InterpreterArgumentsExtraDash bool `json:"interpreter_arguments_extra_dash,omitempty"`
	InterpreterArgumentsExtraAllowed bool `json:"interpreter_arguments_extra_allowed,omitempty"`
	InterpreterEnvironment map[string]string `json:"interpreter_environment,omitempty"`
	Body string `json:"body,omitempty"`
	BodyFingerprint string `json:"body_fingerprint,omitempty"`
	Context *ScriptletContext `json:"-"`
	ContextIdentifier string `json:"context_identifier,omitempty"`
	Fingerprint string `json:"fingerprint"`
	
	// FIXME:  Should move out of scriptlet and into library!
	Source ScriptletSource `json:"source"`
	
	// FIXME:  Should move out of scriptlet and into library!
	Index uint `json:"-"`
	Visible bool `json:"visible"`
	Hidden bool `json:"hidden"`
	Menus []string `json:"menus,omitempty"`
}

type ScriptletSource struct {
	Path string `json:"path"`
	LineStart uint `json:"line_start"`
	LineEnd uint `json:"line_end"`
	BodyOffset uint `json:"body_offset"`
}

type ScriptletContext struct {
	Identifier string `json:"identifier"`
	ExecutablePaths []string `json:"executable_paths"`
	Environment map[string]string `json:"environment"`
}




type Library struct {
	
	Scriptlets LibraryScriptlets `json:"scriptlets"`
	
	ScriptletFingerprints []string `json:"scriptlet_fingerprints"`
	ScriptletsByFingerprint map[string]uint `json:"scriptlets_by_fingerprint"`
	
	ScriptletLabels []string `json:"scriptlet_labels"`
	ScriptletLabelsAll []string `json:"scriptlet_labels_all"`
	ScriptletsByLabel map[string]uint `json:"scriptlets_by_label"`
	
	ScriptletContexts map[string]*ScriptletContext `json:"scriptlet_contexts"`
	
	LibrarySources LibrarySources `json:"library_sources"`
	LibraryContext *LibraryContext `json:"library_context"`
	
	LibraryIdentifier string `json:"library_identifier"`
	LibraryFingerprint string `json:"library_fingerprint"`
	SourcesFingerprint string `json:"sources_fingerprint"`
	
	url string
}




type LibraryContext struct {
	SelfExecutable string `json:"self_executable"`
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
			ScriptletContexts : make (map[string]*ScriptletContext, 16),
			LibraryContext : & LibraryContext {},
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


func (_library *Library) ResolveContextByIdentifier (_fingerprint string) (*ScriptletContext, bool, *Error) {
	if _context, _exists := _library.ScriptletContexts[_fingerprint]; _exists {
		return _context, true, nil
	} else {
		return nil, false, Errorf (0x30d90869, "invalid scriptlet context fingerprint `%s`", _fingerprint)
	}
}


func (_library *Library) SelectLibrarySources () (LibrarySources, *Error) {
	return _library.LibrarySources, nil
}

func (_library *Library) SelectLibraryContext () (*LibraryContext, *Error) {
	return _library.LibraryContext, nil
}


func (_library *Library) Url () (string) {
	return _library.url
}

func (_library *Library) Identifier () (string, *Error) {
	_identifier := _library.LibraryIdentifier
	if _identifier != "" {
		return _identifier, nil
	} else {
		return "", Errorf (0x41ca16f5, "invalid state")
	}
}

func (_library *Library) Fingerprint () (string, *Error) {
	_fingerprint := _library.LibraryFingerprint
	if _fingerprint != "" {
		return _fingerprint, nil
	} else {
		return "", Errorf (0x7c26bcc2, "invalid state")
	}
}


func (_library *Library) Close () (*Error) {
	*_library = Library {}
	return nil
}




func includeScriptlet (_library *Library, _scriptlet *Scriptlet) (*Error) {
	
	if _scriptlet.Label != strings.TrimSpace (_scriptlet.Label) {
		return Errorf (0xd8797e9e, "invalid scriptlet label `%s`", _scriptlet.Label)
	}
	if _scriptlet.Label == "" {
		return Errorf (0xaede3d8c, "invalid scriptlet label `%s`", _scriptlet.Label)
	}
	if _, _exists := _library.ScriptletsByLabel[_scriptlet.Label]; _exists {
		return Errorf (0x883f9a7f, "duplicate scriptlet label `%s`", _scriptlet.Label)
	}
	
	if _scriptlet.ContextIdentifier != "" {
		if _, _exists := _library.ScriptletContexts[_scriptlet.ContextIdentifier]; !_exists {
			return Errorf (0xc9cc9f6e, "invalid scriptlet context identifier `%s`", _scriptlet.ContextIdentifier)
		}
	}
	
	switch _scriptlet.Interpreter {
		case "<detect>", "<print>", "<menu>" :
			// NOP
		default :
			return Errorf (0xbf289098, "invalid scriptlet interpreter `%s`", _scriptlet.Interpreter)
	}
	if (_scriptlet.InterpreterExecutable != "") || (_scriptlet.InterpreterArguments != nil) || (_scriptlet.InterpreterEnvironment != nil) || _scriptlet.InterpreterArgumentsExtraDash || _scriptlet.InterpreterArgumentsExtraAllowed {
		return Errorf (0x901675e8, "invalid scriptlet interpreter state")
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
			return Errorf (0x4b8aacf2, "invalid scriptlet kind `%s`", _scriptlet.Kind)
	}
	
	_fingerprint := NewFingerprinter () .
			StringWithLen (_scriptlet.Label) .
			StringWithLen (_scriptlet.Kind) .
			StringWithLen (_scriptlet.Interpreter) .
			StringWithLen (_scriptlet.InterpreterExecutable) .
			StringsWithLen (_scriptlet.InterpreterArguments) .
			Bool (_scriptlet.InterpreterArgumentsExtraDash) .
			Bool (_scriptlet.InterpreterArgumentsExtraAllowed) .
			StringsMap (_scriptlet.InterpreterEnvironment) .
			StringWithLen (_scriptlet.BodyFingerprint) .
			StringWithLen (_scriptlet.ContextIdentifier) .
			Build ()
	
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


func includeScriptletContext (_library *Library, _context *ScriptletContext) (*Error) {
	
	if _context.Identifier == "" {
		return Errorf (0x92fc0d53, "invalid scriptlet context identifier `%s`", _context.Identifier)
	}
	if _, _exists := _library.ScriptletContexts[_context.Identifier]; _exists {
		return Errorf (0xfe91d3ae, "invalid scriptlet context identifier `%s`", _context.Identifier)
	}
	
	_library.ScriptletContexts[_context.Identifier] = _context
	
	return nil
}




func includeSource (_library *Library, _source *Source) (*Error) {
	if _source.Path == "" {
		return Errorf (0x12bdc134, "invalid state")
	}
	if _source.FingerprintMeta == "" {
		return Errorf (0x152074de, "invalid state")
	}
//	if _source.FingerprintData == "" {
//		return Errorf (0x401d0c16, "invalid state")
//	}
	for _, _existing := range _library.LibrarySources {
		if _existing.Path == _source.Path {
			return Errorf (0xf01b93ea, "invalid state")
		}
		if _existing.FingerprintMeta == _source.FingerprintMeta {
			return Errorf (0x310f6193, "invalid state")
		}
		if (_existing.FingerprintData == _source.FingerprintData) && (_existing.FingerprintData != "") {
			return Errorf (0x00fb18a1, "invalid state %#v %#v", _existing, _source)
		}
	}
	_library.LibrarySources = append (_library.LibrarySources, _source)
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

