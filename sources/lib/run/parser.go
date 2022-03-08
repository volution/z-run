

package zrun


import "fmt"
import "io"
import "os"
import "os/exec"
import "path"
import "path/filepath"
import "regexp"
import "sort"
import "strings"
import "time"
import "unicode"
import "unicode/utf8"

import mpb "github.com/vbauerster/mpb/v5"
import mpb_decor "github.com/vbauerster/mpb/v5/decor"

import . "github.com/volution/z-run/lib/library"
import . "github.com/volution/z-run/lib/common"




type parseContext struct {
	scriptletContext *ScriptletContext
}




func parseLibrary (_sources []*Source, _libraryIdentifier string, _context *Context) (*Library, *Error) {
	
	_parseContext := & parseContext {}
	
	_parseContext.scriptletContext = & ScriptletContext {
			Identifier : "00000000000000000000000000000000",
			ExecutablePaths : make ([]string, 0, 16),
			EnvironmentOverrides : make (map[string]string, 128),
			EnvironmentFallbacks : make (map[string]string, 128),
		}
	
	_library := NewLibrary ()
	_library.LibraryIdentifier = _libraryIdentifier
	_library.LibraryFingerprint = _libraryIdentifier + "--parsing"
	
	if _error := includeScriptletContext (_library, _parseContext.scriptletContext); _error != nil {
		return nil, _error
	}
	
	if _error := MakeCacheFolder (_context.cacheRoot, "parse-sockets"); _error != nil {
		return nil, _error
	}
	_libraryUrl := fmt.Sprintf ("unix:%s", path.Join (_context.cacheRoot, "parse-sockets", fmt.Sprintf ("%s-%08x.sock", GenerateRandomToken (), os.Getpid ())))
	
	var _rpc *LibraryRpcServer
	if _rpc_0, _error := NewLibraryRpcServer (_library, _libraryUrl); _error == nil {
		_rpc = _rpc_0
	} else {
		return nil, _error
	}
	if _error := _rpc.ServeStart (); _error != nil {
		return nil, _error
	}
	defer _rpc.ServeStop ()
	
	var _progress_0 *mpb.Progress
	var _progress *mpb.Bar
	
	if IsStderrTerminal () {
		_progress_0 = mpb.New (
				mpb.WithOutput (os.Stderr),
				mpb.WithRefreshRate (250 * time.Millisecond),
				mpb.WithWidth (40), // FIXME:  Detect console width!
			)
		_progress = _progress_0.AddBar (
				int64 (len (_sources)),
				mpb.PrependDecorators (
						mpb_decor.Name ("[z-run]"),
					),
				mpb.AppendDecorators (
						mpb_decor.NewPercentage ("%d", mpb_decor.WCSyncWidthR),
					),
				mpb.BarRemoveOnComplete (),
			)
	}
	
	for _, _source := range _sources {
//		if _path, _error := filepath.Rel (_context.workspace, _source.Path); _error == nil {
//			Logf ('s', 0x1dcfbf6a, "parsing `./%s`...", _path)
//		} else {
//			return nil, Errorw (0xfce841bd, _error)
//		}
		if _error := parseFromSource (_library, _source, _context, _parseContext); _error != nil {
			return nil, _error
		}
		if _progress != nil {
			_progress.SetTotal (int64 (len (_sources) + len (_library.Scriptlets) + 1), false)
			_progress.Increment ()
		}
	}
	
	_loop : for _index := 0; _index < len (_library.Scriptlets); {
		
		if _progress != nil {
			_progress.SetTotal (int64 (len (_sources) + len (_library.Scriptlets) + 1), false)
			_progress.SetCurrent (int64 (len (_sources) + _index + 1))
		}
		
		_scriptlet := _library.Scriptlets[_index]
		
		switch _scriptlet.Kind {
			case "executable-pending", "generator-pending", "script-replacer-pending", "print-replacer-pending" :
				if _error := parseInterpreter (_scriptlet); _error != nil {
					return nil, _error
				}
//				Logf ('d', 0xb76fe00e, "parsed interpreter for `%s` (`%s` / `%s`)...", _scriptlet.Label, _scriptlet.Kind, _scriptlet.Interpreter)
		}
		
		switch _scriptlet.Kind {
			case "executable-pending" :
				_scriptlet.Kind = "executable"
//				Logf ('d', 0xaa6320d9, "parsed interpreter for `%s` (`%s` / `%s`)...", _scriptlet.Label, _scriptlet.Kind, _scriptlet.Interpreter)
		}
		
		switch _scriptlet.Kind {
			case "generator-pending" :
				if _error := parseFromGenerator (_library, _rpc.Url (), _library.LibraryIdentifier, _library.LibraryFingerprint, _scriptlet, _context, _parseContext); _error == nil {
					_scriptlet.Kind = "generator"
//					Logf ('d', 0x84787ea0, "parsed generator `%s` (`%s` / `%s`)...", _scriptlet.Label, _scriptlet.Kind, _scriptlet.Interpreter)
					continue _loop
				} else {
					return nil, _error
				}
			case "script-replacer-pending", "print-replacer-pending" :
				if _error := parseFromReplacer (_library, _rpc.Url (), _library.LibraryIdentifier, _library.LibraryFingerprint, _scriptlet, _context); _error == nil {
					switch _scriptlet.Kind {
						case "script-replacer-pending" :
							_scriptlet.Kind = "executable-pending"
							_scriptlet.Interpreter = "<detect>"
						case "print-replacer-pending" :
							_scriptlet.Kind = "executable-pending"
							_scriptlet.Interpreter = "<print>"
						default :
							panic (0x6ff57d12)
					}
					_scriptlet.InterpreterExecutable = ""
					_scriptlet.InterpreterArguments = nil
					_scriptlet.InterpreterEnvironmentOverrides = nil
					_scriptlet.InterpreterEnvironmentFallbacks = nil
//					Logf ('d', 0x50ec5e25, "parsed replacer `%s` (`%s` / `%s`)...", _scriptlet.Label, _scriptlet.Kind, _scriptlet.Interpreter)
					continue _loop
				} else {
					return nil, _error
				}
		}
		
//		Logf ('d', 0xd28b083f, "validating `%s` (`%s` / `%s`)...", _scriptlet.Label, _scriptlet.Kind, _scriptlet.Interpreter)
		switch _scriptlet.Kind {
			case "executable", "generator" :
				// NOP
			case "menu-pending" :
				// NOP
			default :
				return nil, Errorf (0xd5f0c788, "invalid state `%s`", _scriptlet.Kind)
		}
		
		_index += 1
	}
	
	if _progress != nil {
		_progress.SetTotal (int64 (len (_sources) + len (_library.Scriptlets) + 1), false)
		_progress.SetCurrent (int64 (len (_sources) + len (_library.Scriptlets)))
	}
	
	
	{
		_menus := make ([]*Scriptlet, 0, 1024)
		for _, _scriptlet := range _library.Scriptlets {
			if _scriptlet.Kind == "menu-pending" {
				_menus = append (_menus, _scriptlet)
			}
		}
		for _, _scriptlet := range _menus {
			if _error := parseFromMenu (_library, _scriptlet, _context); _error != nil {
				return nil, _error
			}
		}
		for _, _scriptlet := range _menus {
			_scriptlet.Kind = "menu"
		}
	}
	
	{
		for _, _scriptlet := range _library.Scriptlets {
			sort.Strings (_scriptlet.Menus)
			for _, _menu := range _scriptlet.Menus {
				switch _menu[0] {
					case '*' :
						// NOP
					case '+' :
						_scriptlet.Hidden = true
					default :
						return nil, Errorf (0x51045db3, "invalid state `%s`", _menu)
				}
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
		sort.Sort (_library.LibrarySources)
		_fingerprints := make ([]string, 0, len (_library.LibrarySources))
		for _, _source := range _library.LibrarySources {
			_fingerprints = append (_fingerprints, _source.FingerprintData)
		}
		sort.Strings (_fingerprints)
		_library.SourcesFingerprint = NewFingerprinter () .StringsWithLen (_fingerprints) .Build ()
	}
	
	_library.LibraryFingerprint = NewFingerprinter () .StringWithLen (_library.LibraryIdentifier) .StringWithLen (_library.SourcesFingerprint) .Build ()
	
	if _progress != nil {
		_progress.SetTotal (_progress.Current () + 1, true)
		_progress.Increment ()
		_progress_0.Wait ()
		os.Stderr.Sync ()
	}
	
	return _library, nil
}




func parseInterpreter (_scriptlet *Scriptlet) (*Error) {
	
	switch _scriptlet.Interpreter {
		case "<detect>" :
			// NOP
		case "<exec>", "<print>", "<template>", "<starlark>", "<menu>", "<go>", "<go+>" :
			return nil
		case "<bash>", "<bash+>", "<python3+>" :
			return nil
		default :
			return Errorf (0xf65704dd, "invalid state `%s`", _scriptlet.Interpreter)
	}
	
	_scriptletBody, _scriptletBodyOffset, _interpreter, _interpreterExecutable, _interpreterArguments, _interpreterArgumentsExtraDash, _interpreterArgumentsExtraAllowed, _interpreterEnvironmentOverrides, _interpreterEnvironmentFallbacks, _error := parseInterpreter_0 (_scriptlet.Label, _scriptlet.Body, "", "<bash>")
	if _error == nil {
		
		_scriptlet.Body = _scriptletBody
		_scriptlet.Source.BodyOffset = _scriptletBodyOffset
		_scriptlet.Interpreter = _interpreter
		_scriptlet.InterpreterExecutable = _interpreterExecutable
		_scriptlet.InterpreterArguments = _interpreterArguments
		_scriptlet.InterpreterArgumentsExtraDash = _interpreterArgumentsExtraDash
		_scriptlet.InterpreterArgumentsExtraAllowed = _interpreterArgumentsExtraAllowed
		_scriptlet.InterpreterEnvironmentOverrides = _interpreterEnvironmentOverrides
		_scriptlet.InterpreterEnvironmentFallbacks = _interpreterEnvironmentFallbacks
		
		return nil
		
	} else {
		return _error
	}
}


func parseInterpreter_0 (_scriptletLabel string, _scriptletBody_0 string, _scriptletHeader string, _interpreterFallback string) (
			_scriptletBody string,
			_scriptletBodyOffset uint,
			_interpreter string,
			_interpreterExecutable string,
			_interpreterArguments []string,
			_interpreterArgumentsExtraDash bool,
			_interpreterArgumentsExtraAllowed bool,
			_interpreterEnvironmentOverrides map[string]string,
			_interpreterEnvironmentFallbacks map[string]string,
			_errorReturn *Error,
		) {
	
	if _scriptletHeader == "" {
		if strings.HasPrefix (_scriptletBody_0, "#!:") {
			_headerLimit := strings.IndexByte (_scriptletBody_0, ' ')
			if _headerLimit < 0 {
				_errorReturn = Errorf (0x48970925, "invalid header for `%s` (space not found)", _scriptletLabel)
				return
			}
			_scriptletHeader = _scriptletBody_0[3:_headerLimit]
			_scriptletBody = _scriptletBody_0[_headerLimit + 1 :]
		} else if strings.HasPrefix (_scriptletBody_0, "#!") {
			_headerLimit := strings.IndexByte (_scriptletBody_0, '\n')
			if _headerLimit < 0 {
				_errorReturn = Errorf (0x42f372b7, "invalid header for `%s` (`\\n` not found)", _scriptletLabel)
				return
			}
			_scriptletHeader = _scriptletBody_0[2:_headerLimit]
			_scriptletBody = _scriptletBody_0[_headerLimit + 1 :]
			_scriptletBodyOffset = 1
		} else {
			if _interpreterFallback != "" {
				_scriptletHeader = _interpreterFallback
				_scriptletBody = _scriptletBody_0
			} else {
				_errorReturn = Errorf (0x273c8b2e, "missing header for `%s` (`#!` not found)", _scriptletLabel)
				return
			}
		}
	} else {
		_scriptletBody = _scriptletBody_0
	}
	
	_scriptletHeader = strings.Trim (_scriptletHeader, " ")
	
	if strings.HasPrefix (_scriptletHeader, "{{}}") {
		_scriptletHeader = "<template>" + _scriptletHeader[4:]
	}
	
	if strings.HasPrefix (_scriptletHeader, "<template>") {
		
		_scriptletHeader = _scriptletHeader[10:]
		_scriptletHeader = strings.Trim (_scriptletHeader, " ")
		
		if _scriptletHeader != "" {
			_errorReturn = Errorf (0x071546cf, "invalid header for `%s` (template with arguments)", _scriptletLabel)
			return
		}
		
		_interpreter = "<template>"
		_interpreterExecutable = ""
		_interpreterArguments = nil
		_interpreterArgumentsExtraDash = false
		_interpreterArgumentsExtraAllowed = true
		_interpreterEnvironmentOverrides = nil
		_interpreterEnvironmentFallbacks = nil
		
	} else if strings.HasPrefix (_scriptletHeader, "<starlark>") {
		
		_scriptletHeader = _scriptletHeader[10:]
		_scriptletHeader = strings.Trim (_scriptletHeader, " ")
		
		if _scriptletHeader != "" {
			_errorReturn = Errorf (0xc7d23157, "invalid header for `%s` (starlark with arguments)", _scriptletLabel)
			return
		}
		
		_interpreter = "<starlark>"
		_interpreterExecutable = ""
		_interpreterArguments = nil
		_interpreterArgumentsExtraDash = false
		_interpreterArgumentsExtraAllowed = true
		_interpreterEnvironmentOverrides = nil
		_interpreterEnvironmentFallbacks = nil
		
	} else if strings.HasPrefix (_scriptletHeader, "<print>") {
		
		_scriptletHeader = _scriptletHeader[7:]
		_scriptletHeader = strings.Trim (_scriptletHeader, " ")
		
		if _scriptletHeader != "" {
			_errorReturn = Errorf (0x6a068d74, "invalid header for `%s` (print with arguments)", _scriptletLabel)
			return
		}
		
		_interpreter = "<print>"
		_interpreterExecutable = ""
		_interpreterArguments = nil
		_interpreterArgumentsExtraDash = false
		_interpreterArgumentsExtraAllowed = false
		_interpreterEnvironmentOverrides = nil
		_interpreterEnvironmentFallbacks = nil
		
	} else if strings.HasPrefix (_scriptletHeader, "<menu>") {
		
		_scriptletHeader = _scriptletHeader[6:]
		_scriptletHeader = strings.Trim (_scriptletHeader, " ")
		
		if _scriptletHeader != "" {
			_errorReturn = Errorf (0xf542516e, "invalid header for `%s` (menu with arguments)", _scriptletLabel)
			return
		}
		
		_interpreter = "<menu>"
		_interpreterExecutable = ""
		_interpreterArguments = nil
		_interpreterArgumentsExtraDash = false
		_interpreterArgumentsExtraAllowed = true
		_interpreterEnvironmentOverrides = nil
		_interpreterEnvironmentFallbacks = nil
		
	} else if strings.HasPrefix (_scriptletHeader, "<go>") || strings.HasPrefix (_scriptletHeader, "<go+>") {
		
		if strings.HasPrefix (_scriptletHeader, "<go>") {
			_interpreter = "<go>"
			_scriptletHeader = _scriptletHeader[4:]
		} else {
			_interpreter = "<go+>"
			_scriptletHeader = _scriptletHeader[5:]
		}
		_scriptletHeader = strings.Trim (_scriptletHeader, " ")
		
		if _scriptletHeader != "" {
			_errorReturn = Errorf (0x80bc049a, "invalid header for `%s` (go with arguments)", _scriptletLabel)
			return
		}
		
		_interpreter = _interpreter
		_interpreterExecutable = ""
		_interpreterArguments = nil
		_interpreterArgumentsExtraDash = false
		_interpreterArgumentsExtraAllowed = true
		_interpreterEnvironmentOverrides = nil
		_interpreterEnvironmentFallbacks = nil
		
	} else {
		
		_header := make ([]string, 0, 16)
		for _, _part := range strings.Split (_scriptletHeader, " ") {
			if _part == "" {
				continue
			}
			_header = append (_header, _part)
		}
		if len (_header) == 0 {
			_errorReturn = Errorf (0x15d0485f, "invalid header for `%s` (empty)", _scriptletLabel)
			return
		}
		
		_interpreter = "<exec>"
		_interpreterExecutable = _header[0]
		_interpreterArguments = make ([]string, 0, len (_header) + 16)
		_interpreterArgumentsExtraDash = false
		_interpreterArgumentsExtraAllowed = false
		_interpreterArgumentsExtraDash = false
		_interpreterArgumentsExtraDashNow := false
		_interpreterEnvironmentOverrides = nil
		_interpreterEnvironmentFallbacks = nil
		
		if strings.HasPrefix (_interpreterExecutable, "<") && strings.HasSuffix (_interpreterExecutable, ">") {
			_interpreterExecutable = _interpreterExecutable[1 : len (_interpreterExecutable) - 1]
			switch _interpreterExecutable {
				
				case "bash" :
					_interpreter = "<" + _interpreterExecutable + ">"
					_interpreterArguments = append (_interpreterArguments)
					_interpreterArgumentsExtraDash = false
					_interpreterArgumentsExtraDashNow = true
					_interpreterArgumentsExtraAllowed = true
				
				case "bash+" :
					_interpreter = "<" + _interpreterExecutable + ">"
					_interpreterExecutable = _interpreterExecutable[0 : len (_interpreterExecutable) - 1]
					_interpreterArgumentsExtraDash = false
					_interpreterArgumentsExtraDashNow = true
					_interpreterArgumentsExtraAllowed = true
				
				case "python", "python2", "python2.7", "python3", "python3.6", "python3.7", "python3.8", "python3.9" :
					_interpreterArguments = append (_interpreterArguments, "-E", "-s", "-S", "-u")
//!					_interpreterArguments = append (_interpreterArguments, "-O", "-O")
					_interpreterArgumentsExtraDash = false
					_interpreterArgumentsExtraDashNow = true
					_interpreterArgumentsExtraAllowed = true
				
				case "python3+" :
					_interpreter = "<" + _interpreterExecutable + ">"
					_interpreterExecutable = _interpreterExecutable[0 : len (_interpreterExecutable) - 1]
					_interpreterArguments = append (_interpreterArguments, "-E", "-s", "-S", "-u")
//!					_interpreterArguments = append (_interpreterArguments, "-O", "-O")
					_interpreterArgumentsExtraDash = false
					_interpreterArgumentsExtraDashNow = true
					_interpreterArgumentsExtraAllowed = true
				
				case "lua", "node", "perl", "ruby", "php", "php5", "php7", "php8" :
					// FIXME: _interpreterArgumentsExtraDash = true
					_interpreterArgumentsExtraAllowed = true
				
				case "tcl" :
					_interpreterExecutable = "tclsh"
					// FIXME: _interpreterArgumentsExtraDash = true
					_interpreterArgumentsExtraAllowed = true
				
				case "awk" :
					_interpreterArguments = append (_interpreterArguments, "-f")
					_interpreterArgumentsExtraDash = true
					_interpreterArgumentsExtraAllowed = false
				
				case "sed", "grep" :
					_interpreterArguments = append (_interpreterArguments, "-E", "-f")
					_interpreterArgumentsExtraDash = false
					_interpreterArgumentsExtraAllowed = false
				
				case "jq" :
					_interpreterArguments = append (_interpreterArguments, "-f")
					_interpreterArgumentsExtraDash = false
					_interpreterArgumentsExtraAllowed = false
				
				case "make" :
					_interpreterArguments = append (_interpreterArguments, "-s", "-r", "-R", "-f")
					_interpreterArgumentsExtraDash = true
					_interpreterArgumentsExtraAllowed = true
				
				case "ninja" :
					_interpreterArguments = append (_interpreterArguments, "-f")
					_interpreterArgumentsExtraDash = true
					_interpreterArgumentsExtraAllowed = true
				
				default :
					_errorReturn = Errorf (0x505f52c6, "invalid interpreter for `%s`", _scriptletLabel)
					return
			}
		} else {
			
			_interpreterArgumentsExtraAllowed = true
		}
		
		if strings.IndexByte (_interpreterExecutable, os.PathSeparator) >= 0 {
			if _interpreterExecutable_0, _error := filepath.Abs (_interpreterExecutable); _error == nil {
				_interpreterExecutable = _interpreterExecutable_0
			} else {
				_errorReturn = Errorw (0xda17c780, _error)
				return
			}
		}
		
		if len (_header) > 1 {
			_interpreterArguments = append (_interpreterArguments, _header[1:] ...)
		}
		
		if _interpreterArgumentsExtraDashNow {
			_interpreterArguments = append (_interpreterArguments, "--")
		}
	}
	
	return
}




func parseFromGenerator (_library *Library, _libraryUrl string, _libraryIdentifier string, _libraryFingerprint string, _source *Scriptlet, _context *Context, _parseContext *parseContext) (*Error) {
//	Logf ('s', 0xf75b04b5, "parsing `:: %s`...", _source.Label)
	if _, _data, _error := loadFromScriptlet (_libraryUrl, _libraryIdentifier, _libraryFingerprint, "", _source, _context); _error == nil {
		return parseFromData (_library, _data, _source.Source.Path, _context, _parseContext)
	} else {
		return _error
	}
}

func parseFromReplacer (_library *Library, _libraryUrl string, _libraryIdentifier string, _libraryFingerprint string, _source *Scriptlet, _context *Context) (*Error) {
//	Logf ('s', 0xc336e8bd, "parsing `:: %s`...", _source.Label)
	if _, _data, _error := loadFromScriptlet (_libraryUrl, _libraryIdentifier, _libraryFingerprint, _source.Interpreter, _source, _context); _error == nil {
		if utf8.Valid (_data) {
			_source.Body = string (_data)
			return nil
		} else {
			return Errorf (0xdb3b92b7, "invalid UTF-8")
		}
	} else {
		return _error
	}
}

func parseFromMenu (_library *Library, _source *Scriptlet, _context *Context) (*Error) {
	if _source.Kind != "menu-pending" {
		return Errorf (0x6834ac93, "invalid state")
	}
	_labels := make ([]string, 0, 1024)
	_matchers := strings.Split (_source.Body, "\n")
	_scriptlets := make ([]*Scriptlet, 0, len (_library.Scriptlets))
	for _, _scriptlet := range _library.Scriptlets {
		if _scriptlet == _source {
			continue
		}
		if _scriptlet.Hidden {
			continue
		}
		_scriptlets = append (_scriptlets, _scriptlet)
	}
	for _, _matcher := range _matchers {
		if (_matcher == "") || (_matcher[0] == '#') {
			continue
		} else {
			switch _matcher[0] {
				case '*', '+' :
					// NOP
				case '=', '^', '~' :
					_matcher = "+" + _matcher
				default :
					return Errorf (0x11ca1466, "invalid menu mode `%s`", _matcher)
			}
			_mode := _matcher[:1]
			_matcher = _matcher[1:]
			if _matcher == "" {
				for _, _scriptlet := range _scriptlets {
					_labels = append (_labels, _scriptlet.Label)
					_scriptlet.Menus = append (_scriptlet.Menus, _mode + " " + _source.Label)
				}
			} else if strings.HasPrefix (_matcher, "= ") {
				_pattern := _matcher[2:]
				for _, _scriptlet := range _scriptlets {
					if _scriptlet.Label == _pattern {
						_labels = append (_labels, _scriptlet.Label)
						_scriptlet.Menus = append (_scriptlet.Menus, _mode + " " + _source.Label)
					}
				}
			} else if strings.HasPrefix (_matcher, "^ ") {
				_pattern := _matcher[2:]
				for _, _scriptlet := range _scriptlets {
					if strings.HasPrefix (_scriptlet.Label, _pattern) {
						_labels = append (_labels, _scriptlet.Label)
						_scriptlet.Menus = append (_scriptlet.Menus, _mode + " " + _source.Label)
					}
				}
			} else if strings.HasPrefix (_matcher, "~ ") {
				var _pattern *regexp.Regexp
				if _pattern_0, _error := regexp.Compile (_matcher[2:]); _error == nil {
					_pattern = _pattern_0
				} else {
					return Errorf (0xabf68b41, "invalid menu matcher `%s`", _matcher)
				}
				for _, _scriptlet := range _scriptlets {
					if _pattern.MatchString (_scriptlet.Label) {
						_labels = append (_labels, _scriptlet.Label)
						_scriptlet.Menus = append (_scriptlet.Menus, _mode + " " + _source.Label)
					}
				}
			} else {
				return Errorf (0xa068e934, "invalid menu matcher `%s`", _matcher)
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
	if len (_labels) == 0 {
		_source.Visible = false
		_source.Hidden = true
	}
	return nil
}




func parseFromSource (_library *Library, _source *Source, _context *Context, _parseContext *parseContext) (*Error) {
	if _data, _error := loadFromSource (_library, _source, _context); _error == nil {
		if _error := includeSource (_library, _source); _error != nil {
			return _error
		}
//?		_library.LibraryFingerprint = NewFingerprinter () .StringWithLen (_library.LibraryFingerprint) .StringWithLen (_source.FingerprintData) .Build ()
		return parseFromData (_library, _data, _source.Path, _context, _parseContext)
	} else {
		return _error
	}
}

func parseFromFile (_library *Library, _sourcePath string, _context *Context, _parseContext *parseContext) (string, *Error) {
	if _fingerprint, _data, _error := loadFromFile (_sourcePath); _error == nil {
		if _error := parseFromData (_library, _data, _sourcePath, _context, _parseContext); _error == nil {
			return _fingerprint, nil
		} else {
			return "", _error
		}
	} else {
		return "", _error
	}
}

func parseFromStream (_library *Library, _stream io.Reader, _sourcePath string, _context *Context, _parseContext *parseContext) (string, *Error) {
	if _fingerprint, _data, _error := loadFromStream (_stream); _error == nil {
		if _error := parseFromData (_library, _data, _sourcePath, _context, _parseContext); _error == nil {
			return _fingerprint, nil
		} else {
			return "", _error
		}
	} else {
		return "", _error
	}
}




func parseFromData (_library *Library, _sourceData []byte, _sourcePath string, _context *Context, _parseContext *parseContext) (*Error) {
	
	var _source string
	if utf8.Valid (_sourceData) {
		_source = string (_sourceData)
	} else {
		return Errorf (0x2a19cfc7, "invalid UTF-8 source")
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
		
//		Logf ('d', 0xc2d2b73d, "processing line (%d):  %s", _lineIndex, _line)
		
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
							return Errorf (0x53eafa1a, "invalid syntax (%d):  missing scriptlet separator `::` | %s", _lineIndex, _line)
						} else {
							_label = _text
						}
						_label = _trimSpace (_label)
						_body = _trimSpace (_body)
					}
					
					if _label == "" {
						return Errorf (0xddec2340, "invalid syntax (%d):  empty scriptlet label | %s", _lineIndex, _line)
					}
					if (_body == "") && (_prefix != ":://") {
						return Errorf (0xc1dc94cc, "invalid syntax (%d):  empty scriptlet body | %s", _lineIndex, _line)
					}
					
					_kind := ""
					_interpreter := ""
					_include := false
					switch _prefix[2:] {
						case "" :
							_kind = "executable"
							_interpreter = "<detect>"
						case ".." :
							_kind = "executable"
							_interpreter = "<print>"
						case "~~" :
							_kind = "script-replacer"
							_interpreter = "<detect>"
						case "~~.." :
							_kind = "print-replacer"
							_interpreter = "<detect>"
						case "==" :
							_kind = "generator"
							_interpreter = "<detect>"
							_hidden = true
						case "&&" :
							_kind = "executable"
							_interpreter = "<detect>"
							_include = true
						case "&&.." :
							_kind = "executable"
							_interpreter = "<print>"
							_include = true
						case "//" :
							_kind = "menu"
							_interpreter = "<menu>"
							if _body == "" {
								if _label == "*" {
									_body = "*"
								} else {
									if strings.HasSuffix (_label, " ...") {
										_body = "+^ " + _label[:len (_label) - 4]
									} else if strings.HasSuffix (_label, " *") {
										_body = "*^ " + _label[:len (_label) - 2]
									} else {
										return Errorf (0x84fa71b1, "invalid syntax (%d):  invalid menu label | %s", _lineIndex, _line)
									}
								}
							}
						default :
							return Errorf (0xfba805b9, "invalid syntax (%d):  unknown scriptlet type | %s", _lineIndex, _line)
					}
					
					if _include {
						_includePath := _body
						if _includePath_0, _error := ReplaceVariables (_includePath); _error != nil {
							return _error
						} else {
							_includePath = _includePath_0
						}
						if _includePath_0, _error := ResolveRelativePath (_context.workspace, path.Dir (_sourcePath), _includePath); _error != nil {
							return _error
						} else {
							_includePath = _includePath_0
						}
						if _includeSource, _error := resolveSource (_includePath, "", nil); _error == nil {
							if _data, _error := loadFromSource (_library, _includeSource, _context); _error == nil {
								if _error := includeSource (_library, _includeSource); _error != nil {
									return _error
								}
								if utf8.Valid (_sourceData) {
									_body = string (_data)
								} else {
									return Errorf (0x16010e20, "invalid UTF-8")
								}
//?								_library.LibraryFingerprint = NewFingerprinter () .StringWithLen (_library.LibraryFingerprint) .StringWithLen (_includeSource.FingerprintData) .Build ()
							} else {
								return _error
							}
						} else {
							return _error
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
						strings.HasPrefix (_lineTrimmed, "<<// ") ||
						strings.HasPrefix (_lineTrimmed, "<<== ") {
					
					_prefix := _lineTrimmed[: strings.IndexByte (_lineTrimmed, ' ')]
					_label := ""
					{
						_label = _lineTrimmed[len (_prefix) + 1:]
						_label = _trimSpace (_label)
					}
					
					if _label == "" {
						return Errorf (0x64c17a76, "invalid syntax (%d):  empty scriptlet label | %s", _lineIndex, _line)
					}
					
					_kind := ""
					_interpreter := ""
					switch _prefix[2:] {
						case "" :
							_kind = "executable"
							_interpreter = "<detect>"
						case ".." :
							_kind = "executable"
							_interpreter = "<print>"
						case "~~" :
							_kind = "script-replacer"
							_interpreter = "<detect>"
						case "~~.." :
							_kind = "print-replacer"
							_interpreter = "<detect>"
						case "==" :
							_kind = "generator"
							_interpreter = "<detect>"
							_hidden = true
						case "//" :
							_kind = "menu"
							_interpreter = "<menu>"
						default :
							return Errorf (0xd08972fe, "invalid syntax (%d):  unknown scriptlet type | %s", _lineIndex, _line)
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
					
				} else if strings.HasPrefix (_lineTrimmed, "&& ") || strings.HasPrefix (_lineTrimmed, "&&?? ") {
					
					_ignoreError := strings.HasPrefix (_lineTrimmed, "&&?? ")
					
					_includePath := _lineTrimmed[strings.IndexByte (_lineTrimmed, ' ') + 1:]
					if _includePath_0, _error := ReplaceVariables (_includePath); _error != nil {
						return _error
					} else {
						_includePath = _includePath_0
					}
					if _includePath_0, _error := ResolveRelativePath (_context.workspace, path.Dir (_sourcePath), _includePath); _error != nil {
						return _error
					} else {
						_includePath = _includePath_0
					}
					
					if !_disabled {
						if _includeSources, _error := resolveSources (_includePath, "", nil, false); _error == nil {
							for _, _includeSource := range _includeSources {
								if _error := parseFromSource (_library, _includeSource, _context, _parseContext); _error != nil {
									return _error
								}
							}
						} else {
							if !_ignoreError {
								return _error
							}
						}
					}
					
				} else if strings.HasPrefix (_lineTrimmed, "&&__ ") || strings.HasPrefix (_lineTrimmed, "&&??__ ") {
					
					_ignoreError := strings.HasPrefix (_lineTrimmed, "&&??__ ")
					
					_includePath := _lineTrimmed[strings.IndexByte (_lineTrimmed, ' ') + 1:]
					if _includePath_0, _error := ReplaceVariables (_includePath); _error != nil {
						return _error
					} else {
						_includePath = _includePath_0
					}
					if _includePath_0, _error := ResolveRelativePath (_context.workspace, path.Dir (_sourcePath), _includePath); _error != nil {
						return _error
					} else {
						_includePath = _includePath_0
					}
					
					if !_disabled {
						if _includeSource, _error := fingerprintSource (_includePath); _error == nil {
							if _error := includeSource (_library, _includeSource); _error != nil {
								return _error
							}
//?							_library.LibraryFingerprint = NewFingerprinter () .StringWithLen (_library.LibraryFingerprint) .StringWithLen (_includeSource.FingerprintData) .Build ()
						} else {
							if !_ignoreError {
								return _error
							}
						}
					}
					
				} else if strings.HasPrefix (_lineTrimmed, "&&== ") {
					
					_descriptor := _lineTrimmed[strings.IndexByte (_lineTrimmed, ' ') + 1:]
					_kind := _descriptor[: strings.IndexByte (_descriptor, ' ')]
					_descriptor = _descriptor[len (_kind) + 1 :]
					_kind = strings.TrimSpace (_kind)
					_descriptor = strings.TrimSpace (_descriptor)
					
					if _kind == "" {
						return Errorf (0xed68e4c3, "invalid syntax (%d):  empty statement | %s", _lineIndex, _line)
					}
					
					switch _kind {
						
						case "path" :
							_kind = "path"
						
						case "environment", "env" :
							_kind = "environment-override"
						case "environment-override", "env-override" :
							_kind = "environment-override"
						case "environment-fallback", "env-fallback" :
							_kind = "environment-fallback"
						
						case "environment-path", "env-path" :
							_kind = "environment-override-path"
						case "environment-override-path", "env-override-path", "environment-path-override", "env-path-override" :
							_kind = "environment-override-path"
						case "environment-fallback-path", "env-fallback-path", "environment-path-fallback", "env-path-fallback" :
							_kind = "environment-fallback-path"
						
						case "environment-append", "env-append" :
							_kind = "environment-override-append"
						case "environment-override-append", "env-override-append" :
							_kind = "environment-override-append"
						case "environment-fallback-append", "env-fallback-append" :
							_kind = "environment-fallback-append"
						
						case "environment-append-path", "env-append-path", "env-path-append" :
							_kind = "environment-override-append-path"
						case "environment-override-append-path", "env-override-append-path", "env-path-override-append" :
							_kind = "environment-override-append-path"
						case "environment-fallback-append-path", "env-fallback-append-path", "env-path-fallback-append" :
							_kind = "environment-fallback-append-path"
						
						case "environment-exclude", "env-exclude" :
							_kind = "environment-exclude"
						
						case "z-run" :
							_kind = "z-run"
						
						default :
							return Errorf (0x5b8c6ee2, "invalid syntax (%d):  invalid statement | %s", _lineIndex, _line)
					}
					
					switch _kind {
						
						case "path" :
							
							if _descriptor == "" {
								return Errorf (0x8c2e1bf8, "invalid syntax (%d):  empty statement path descriptor | %s", _lineIndex, _line)
							}
							_path := _descriptor
							if _path_0, _error := ReplaceVariables (_path); _error != nil {
								return _error
							} else {
								_path = _path_0
							}
							if _path_0, _error := ResolveAbsolutePath (_context.workspace, path.Dir (_sourcePath), _path); _error != nil {
								return _error
							} else {
								_path = _path_0
							}
							if _stat, _error := os.Stat (_path); _error == nil {
								if ! _stat.IsDir () {
									return Errorf (0xae85ad7e, "invalid syntax (%d):  invalid statement path value | %s", _lineIndex, _line)
								}
							} else {
								return Errorw (0x79069c08, _error)
							}
							if !_disabled {
								_parseContext.scriptletContext.ExecutablePaths = append (_parseContext.scriptletContext.ExecutablePaths, _path)
							}
							
						case "environment-override", "environment-fallback",
								"environment-override-path", "environment-fallback-path",
								"environment-override-append-path", "environment-fallback-append-path" :
							
							if _descriptor == "" {
								return Errorf (0xa36b04fa, "invalid syntax (%d):  empty statement environment descriptor | %s", _lineIndex, _line)
							}
							_descriptor := strings.SplitN (_descriptor, " ", 2)
							if len (_descriptor) != 2 {
								return Errorf (0xfc01ef6a, "invalid syntax (%d):  invalid statement environment descriptor | %s", _lineIndex, _line)
							}
							_name := strings.TrimSpace (_descriptor[0])
							_value := strings.TrimSpace (_descriptor[1])
							if _name == "" {
								return Errorf (0x6bce31bb, "invalid syntax (%d):  empty statement environment key | %s", _lineIndex, _line)
							}
							if _value_0, _error := ReplaceVariables (_value); _error != nil {
								return _error
							} else {
								_value = _value_0
							}
							if (_kind == "environment-override-path") || (_kind == "environment-fallback-path") || (_kind == "environment-override-append-path") || (_kind == "environment-fallback-append-path") {
								if _value == "" {
									return Errorf (0x2124d511, "invalid syntax (%d):  empty statement environment descriptor | %s", _lineIndex, _line)
								}
								if _path_0, _error := ResolveAbsolutePath (_context.workspace, path.Dir (_sourcePath), _value); _error != nil {
									return _error
								} else {
									_value = _path_0
								}
							}
							if !_disabled {
								var _environment map[string]string
								switch _kind {
									case "environment-override", "environment-override-path", "environment-override-append", "environment-override-append-path" :
										_environment = _parseContext.scriptletContext.EnvironmentOverrides
									case "environment-fallback", "environment-fallback-path", "environment-fallback-append", "environment-fallback-append-path" :
										_environment = _parseContext.scriptletContext.EnvironmentFallbacks
									default :
										panic (0xb8623b38)
								}
								_valueExisting, _exists := _environment[_name]
								if (_kind == "environment-override-append") || (_kind == "environment-override-append-path") || (_kind == "environment-fallback-append") || (_kind == "environment-fallback-append-path") {
									if _exists {
										_value = _valueExisting + string (os.PathListSeparator) + _value
									}
								} else {
									if _exists {
										return Errorf (0x774b50de, "invalid syntax (%d):  duplicate statement environment key | %s", _lineIndex, _line)
									}
								}
								_environment[_name] = _value
							}
							
						case "environment-exclude" :
							if _descriptor == "" {
								return Errorf (0x7f049882, "invalid syntax (%d):  empty statement environment descriptor | %s", _lineIndex, _line)
							}
							_descriptor := strings.Split (_descriptor, " ")
							for _, _name := range _descriptor {
								if _name == "" {
									continue
								}
								if _, _exists := _parseContext.scriptletContext.EnvironmentOverrides[_name]; _exists && !_disabled {
									return Errorf (0x0ed5990f, "invalid syntax (%d):  duplicate statement environment key | %s", _lineIndex, _line)
								}
								if !_disabled {
									_parseContext.scriptletContext.EnvironmentOverrides[_name] = ""
								}
							}
						
						case "z-run" :
							if _descriptor == "" {
								return Errorf (0x2abc5316, "invalid syntax (%d):  empty statement `z-run` executable descriptor | %s", _lineIndex, _line)
							}
							_newExecutable := _descriptor
							if _newExecutable_0, _error := ReplaceVariables (_newExecutable); _error != nil {
								return _error
							} else {
								_newExecutable = _newExecutable_0
							}
							if _newExecutable_0, _error := ResolveAbsolutePath (_context.workspace, path.Dir (_sourcePath), _newExecutable); _error == nil {
								_newExecutable = _newExecutable_0
							} else {
								return _error
							}
							if _newExecutable_0, _error := filepath.EvalSymlinks (_newExecutable); _error == nil {
								_newExecutable = _newExecutable_0
							} else {
								return Errorw (0x349ba402, _error)
							}
							var _statNewExecutable os.FileInfo
							if _stat_0, _error := os.Stat (_newExecutable); _error == nil {
								if ! _stat_0.Mode () .IsRegular () {
									return Errorf (0x60329e98, "invalid syntax (%d):  invalid statement `z-run` executable path (not a file) | %s", _lineIndex, _line)
								}
								_statNewExecutable = _stat_0
							} else {
								return Errorw (0xc9b0eb17, _error)
							}
							var _statSelfExecutable os.FileInfo
							if _stat_0, _error := os.Stat (_context.selfExecutable); _error == nil {
								_statSelfExecutable = _stat_0
							} else {
								return Errorw (0x4266df29, _error)
							}
							if !_disabled {
								if ! os.SameFile (_statNewExecutable, _statSelfExecutable) {
									// FIXME:  A better solution would be nice!
									return _context.preMainReExecute (_newExecutable)
								} else {
									if _library.LibraryContext.SelfExecutable == "" {
										_library.LibraryContext.SelfExecutable = _newExecutable
									} else if _library.LibraryContext.SelfExecutable != _newExecutable {
										return Errorf (0xf0be06d3, "invalid state (%d):  mismatched `z-run` executable | %s", _lineIndex, _line)
									}
								}
							}
						
						default :
							return Errorf (0xd901dd89, "invalid syntax (%d):  invalid statement | %s", _lineIndex, _line)
					}
					
				} else if strings.HasPrefix (_lineTrimmed, "{{") {
					if _disabled {
						_state = SKIPPING
					} else {
						return Errorf (0x79d4d781, "invalid syntax (%d):  unknown block type | %s", _lineIndex, _line)
					}
					
				} else if strings.HasPrefix (_line, "#!/") {
					if (_lineIndex == 1) {
						// NOP
					} else {
						// FIXME:  This should be a warning!
					}
					
				} else {
					return Errorf (0x9f8daae4, "invalid syntax (%d):  unexpected statement | %s", _lineIndex, _line)
				}
			
			case SCRIPTLET_BODY :
				_lineTrimmed := _trimRightSpace (_line)
				
				if _lineTrimmed == "!!" {
					_scriptletState.body = _scriptletState.bodyBuffer.String ()
					_scriptletState.lineEnd = _lineIndex
					_state = SCRIPTLET_PUSH
					
				} else if strings.HasPrefix (_lineTrimmed, "!!") {
					return Errorf (0xf9900c0c, "invalid syntax (%d):  unexpected statement | %s", _lineIndex, _line)
					
				} else if _lineTrimmed == "" {
					_scriptletState.bodyBuffer.WriteByte ('\n')
					
				} else {
					if _scriptletState.bodyLines == 0 {
						if _stripIndex := strings.IndexFunc (_line, func (r rune) (bool) { return ! unicode.IsSpace (r) }); _stripIndex > 0 {
							_scriptletState.bodyStrip = _line[:_stripIndex]
						}
					}
					if ! strings.HasPrefix (_line, _scriptletState.bodyStrip) {
						Logf ('w', 0xc4e05443, "invalid syntax (%d):  unexpected indentation `%s`", _lineIndex, _line)
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
					return Errorf (0x183de0fd, "invalid syntax (%d):  unexpected statement | %s", _lineIndex, _line)
				} else {
					// NOP
				}
		}
		
		if _state == SCRIPTLET_PUSH {
			if !_scriptletState.disabled {
				_bodyFingerprint := FingerprintString (_scriptletState.body)
				_scriptlet := & Scriptlet {
						Label : _scriptletState.label,
						Kind : _scriptletState.kind,
						Interpreter : _scriptletState.interpreter,
						Context : _parseContext.scriptletContext,
						ContextIdentifier : _parseContext.scriptletContext.Identifier,
						Visible : _scriptletState.visible,
						Hidden : _scriptletState.hidden,
						BodyFingerprint : _bodyFingerprint,
						Body : _scriptletState.body,
						Source : ScriptletSource {
								Path : _sourcePath,
								LineStart : _scriptletState.lineStart,
								LineEnd : _scriptletState.lineEnd,
								BodyOffset : 0,
							},
					}
				if _error := includeScriptlet (_library, _scriptlet); _error != nil {
					return _error
				}
				if _error := parseInterpreter (_scriptlet); _error != nil {
					return _error
				}
			}
			_state = WAITING
		}
	}
	
	switch _state {
		case WAITING :
		case SCRIPTLET_BODY :
			return Errorf (0x9d55df33, "invalid syntax (%d):  missing scriptlet body closing tag `!!` (and reached end of file)", _lineIndex)
		case SKIPPING :
			return Errorf (0x357f15e1, "invalid syntax (%d):  missing comment body closing tag `##}}` (and reached end of file)", _lineIndex)
		default :
			return Errorf (0xc0f78380, "invalid syntax (%d):  unexpected state `%s` (and reached end of file)", _lineIndex, _state)
	}
	
	return nil
}




func loadFromSource (_library *Library, _source *Source, _context *Context) ([]byte, *Error) {
	if _fingerprint, _data, _error := loadFromSource_0 (_library, _source, _context); _error == nil {
		if _source.FingerprintData == "" {
			_source.FingerprintData = _fingerprint
		} else if _source.FingerprintData != _fingerprint {
			return nil, Errorf (0x6293d72d, "invalid state")
		}
		return _data, nil
	} else {
		return nil, _error
	}
}


func loadFromSource_0 (_library *Library, _source *Source, _context *Context) (string, []byte, *Error) {
	
	if !_source.Executable {
		return loadFromFile (_source.Path)
		
	} else {
		
		_executable := _source.Path
		if ! path.IsAbs (_executable) {
			if _executable_0, _error := filepath.Abs (_executable); _error == nil {
				_executable = _executable_0
			} else {
				return "", nil, Errorw (0x5d248eeb, _error)
			}
		}
		
		// FIXME:  Hash the actual executable!
		_fingerprint := _source.FingerprintMeta
		
		_command := & exec.Cmd {
				Path : _executable,
				Args : []string {
						"[z-run:generator]",
					},
				Env : prepareEnvironment (_context, nil, nil),
				Stdin : nil,
				Stdout : nil,
				Stderr : os.Stderr,
				Dir : _context.workspace,
			}
		
		if _, _data, _error := loadFromCommand (_command); _error == nil {
			return _fingerprint, _data, _error
		} else {
			return "", nil, _error
		}
	}
}


func loadFromFile (_path string) (string, []byte, *Error) {
	if _stream, _error := os.Open (_path); _error == nil {
		defer _stream.Close ()
		return loadFromStream (_stream)
	} else {
		return "", nil, Errorw (0xe1528040, _error)
	}
}


func loadFromStream (_stream io.Reader) (string, []byte, *Error) {
	if _data, _error := io.ReadAll (_stream); _error == nil {
		_fingerprint := NewFingerprinter () .Bytes (_data) .Build ()
		return _fingerprint, _data, nil
	} else {
		return "", nil, Errorw (0x42be0790, _error)
	}
}


func loadFromScriptlet (_libraryUrl string, _libraryIdentifier string, _libraryFingerprint string, _interpreter string, _scriptlet *Scriptlet, _context *Context) (string, []byte, *Error) {
	
	var _command *exec.Cmd
	var _descriptors []int
	if _command_0, _descriptors_0, _error := prepareExecution (_libraryUrl, _libraryIdentifier, _libraryFingerprint, _interpreter, _scriptlet, false, _context); _error == nil {
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


func loadFromCommand (_command *exec.Cmd) (string, []byte, *Error) {
	
	if _command.Stderr == nil {
		_command.Stderr = os.Stderr
	}
	
	if _exitCode, _data, _error := ProcessExecuteGetStdout (_command); _error == nil {
		if _exitCode == 0 {
			_fingerprint := NewFingerprinter () .Bytes (_data) .Build ()
			return _fingerprint, _data, nil
		} else {
			return "", nil, Errorf (0x42669a76, "command failed with exit code `%d`", _exitCode)
		}
	} else {
		return "", nil, _error
	}
}

