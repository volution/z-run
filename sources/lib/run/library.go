

package zrun


import "strings"

import . "github.com/volution/z-run/lib/library"
import . "github.com/volution/z-run/lib/common"




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
