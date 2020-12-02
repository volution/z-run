

package zrun


import "fmt"
import "io"
import "io/ioutil"
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


import isatty "github.com/mattn/go-isatty"
import mpb "github.com/vbauerster/mpb/v5"
import mpb_decor "github.com/vbauerster/mpb/v5/decor"




type parseContext struct {
	scriptletContext *ScriptletContext
}




func parseLibrary (_sources []*Source, _environmentFingerprint string, _context *Context) (*Library, *Error) {
	
	_parseContext := & parseContext {}
	
	_parseContext.scriptletContext = & ScriptletContext {
			Fingerprint : generateRandomToken (),
			ExecutablePaths : make ([]string, 0, 16),
			Environment : make (map[string]string, 128),
		}
	
	_library := NewLibrary ()
	_library.EnvironmentFingerprint = _environmentFingerprint
	_library.LibraryFingerprint = _environmentFingerprint
	
	if _error := includeScriptletContext (_library, _parseContext.scriptletContext); _error != nil {
		return nil, _error
	}
	
	_libraryUrl := ""
	if true {
		var _cacheRoot = _context.cacheRoot
		if _cacheRoot == "" {
			_cacheRoot = "/tmp"
		}
		var _socketToken = generateRandomToken ()
		var _socketPath = path.Join (_cacheRoot, fmt.Sprintf ("%s-%08x.sock", _socketToken, os.Getpid ()))
		_libraryUrl = fmt.Sprintf ("unix:%s", _socketPath)
	} else {
		_libraryUrl = fmt.Sprintf ("unix:@%s-%08x", _environmentFingerprint, os.Getpid ())
	}
	
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
	
	if isatty.IsTerminal (os.Stderr.Fd ()) {
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
//			logf ('s', 0x1dcfbf6a, "parsing `./%s`...", _path)
//		} else {
//			return nil, errorw (0xfce841bd, _error)
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
				if _error := parseInterpreter (_library, _scriptlet, _context); _error != nil {
					return nil, _error
				}
//				logf ('d', 0xb76fe00e, "parsed interpreter for `%s` (`%s` / `%s`)...", _scriptlet.Label, _scriptlet.Kind, _scriptlet.Interpreter)
		}
		
		switch _scriptlet.Kind {
			case "executable-pending" :
				_scriptlet.Kind = "executable"
//				logf ('d', 0xaa6320d9, "parsed interpreter for `%s` (`%s` / `%s`)...", _scriptlet.Label, _scriptlet.Kind, _scriptlet.Interpreter)
		}
		
		switch _scriptlet.Kind {
			case "generator-pending" :
				if _error := parseFromGenerator (_library, _rpc.Url (), _library.LibraryFingerprint, _scriptlet, _context, _parseContext); _error == nil {
					_scriptlet.Kind = "generator"
//					logf ('d', 0x84787ea0, "parsed generator `%s` (`%s` / `%s`)...", _scriptlet.Label, _scriptlet.Kind, _scriptlet.Interpreter)
					continue _loop
				} else {
					return nil, _error
				}
			case "script-replacer-pending", "print-replacer-pending" :
				if _error := parseFromReplacer (_library, _rpc.Url (), _library.LibraryFingerprint, _scriptlet, _context); _error == nil {
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
					_scriptlet.InterpreterEnvironment = nil
//					logf ('d', 0x50ec5e25, "parsed replacer `%s` (`%s` / `%s`)...", _scriptlet.Label, _scriptlet.Kind, _scriptlet.Interpreter)
					continue _loop
				} else {
					return nil, _error
				}
		}
		
//		logf ('d', 0xd28b083f, "validating `%s` (`%s` / `%s`)...", _scriptlet.Label, _scriptlet.Kind, _scriptlet.Interpreter)
		switch _scriptlet.Kind {
			case "executable", "generator" :
				// NOP
			case "menu-pending" :
				// NOP
			default :
				return nil, errorf (0xd5f0c788, "invalid state `%s`", _scriptlet.Kind)
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
						return nil, errorf (0x51045db3, "invalid state `%s`", _menu)
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
		sort.Sort (_library.Sources)
		_fingerprints := make ([]string, 0, len (_library.Sources))
		for _, _source := range _library.Sources {
			_fingerprints = append (_fingerprints, _source.FingerprintData)
		}
		sort.Strings (_fingerprints)
		_library.SourcesFingerprint = NewFingerprinter () .StringsWithLen (_fingerprints) .Build ()
	}
	
	_library.LibraryFingerprint = NewFingerprinter () .StringWithLen (_library.EnvironmentFingerprint) .StringWithLen (_library.SourcesFingerprint) .Build ()
	
	if _progress != nil {
		_progress.SetTotal (_progress.Current () + 1, true)
		_progress.Increment ()
		_progress_0.Wait ()
		os.Stderr.Sync ()
	}
	
	return _library, nil
}


func parseInterpreter (_library *Library, _scriptlet *Scriptlet, _context *Context) (*Error) {
	
	switch _scriptlet.Interpreter {
		case "<detect>" :
			// NOP
		case "<exec>", "<print>", "<template>", "<menu>" :
			return nil
		case "<bash*>", "<python*>", "<python2*>", "<python3*>" :
			return nil
		default :
			return errorf (0xf65704dd, "invalid state `%s`", _scriptlet.Interpreter)
	}
	
	_headerLine := ""
	
	if ! strings.HasPrefix (_scriptlet.Body, "#!") {
		_headerLine = "<bash*>"
	} else {
		_headerLimit := strings.IndexByte (_scriptlet.Body, '\n')
		if _headerLimit < 0 {
			return errorf (0x42f372b7, "invalid header for `%s` (`\n` not found)", _scriptlet.Label)
		}
		_headerLine = _scriptlet.Body[2:_headerLimit]
		_scriptlet.Body = _scriptlet.Body[_headerLimit + 1 :]
	}
	
	_headerLine = strings.Trim (_headerLine, " ")
	
	if strings.HasPrefix (_headerLine, "{{}}") {
		_headerLine = "<template>" + _headerLine[4:]
	}
	
	if strings.HasPrefix (_headerLine, "<template>") {
		
		_headerLine = _headerLine[10:]
		_headerLine = strings.Trim (_headerLine, " ")
		
		if _headerLine != "" {
			return errorf (0x071546cf, "invalid header for `%s` (template with arguments)", _scriptlet.Label)
		}
		
		_scriptlet.Interpreter = "<template>"
		_scriptlet.InterpreterExecutable = ""
		_scriptlet.InterpreterArguments = nil
		_scriptlet.InterpreterEnvironment = nil
		
	} else if strings.HasPrefix (_headerLine, "<print>") {
		
		_headerLine = _headerLine[7:]
		_headerLine = strings.Trim (_headerLine, " ")
		
		if _headerLine != "" {
			return errorf (0x6a068d74, "invalid header for `%s` (print with arguments)", _scriptlet.Label)
		}
		
		_scriptlet.Interpreter = "<print>"
		_scriptlet.InterpreterExecutable = ""
		_scriptlet.InterpreterArguments = nil
		_scriptlet.InterpreterEnvironment = nil
		
	} else if strings.HasPrefix (_headerLine, "<menu>") {
		
		_headerLine = _headerLine[6:]
		_headerLine = strings.Trim (_headerLine, " ")
		
		if _headerLine != "" {
			return errorf (0xf542516e, "invalid header for `%s` (menu with arguments)", _scriptlet.Label)
		}
		
		_scriptlet.Interpreter = "<menu>"
		_scriptlet.InterpreterExecutable = ""
		_scriptlet.InterpreterArguments = nil
		_scriptlet.InterpreterEnvironment = nil
		
	} else {
		
		_header := make ([]string, 0, 16)
		for _, _part := range strings.Split (_headerLine, " ") {
			if _part == "" {
				continue
			}
			_header = append (_header, _part)
		}
		if len (_header) == 0 {
			return errorf (0x15d0485f, "invalid header for `%s` (empty)", _scriptlet.Label)
		}
		
		_interpreter := "<exec>"
		_executable := _header[0]
		_arguments := make ([]string, 0, len (_header) + 16)
		
		if strings.HasPrefix (_executable, "<") && strings.HasSuffix (_executable, ">") {
			_executable = _executable[1 : len (_executable) - 1]
			switch _executable {
				case "bash*", "python*", "python2*", "python3*" :
					_interpreter = "<" + _executable + ">"
					_executable = _executable[0 : len (_executable) - 1]
				case "bash", "python", "python2", "python3" :
					// NOP
				case "lua", "node", "perl", "ruby", "php", "php5", "php7", "php8" :
					// NOP
				case "tcl" :
					_executable = "tclsh"
				case "awk" :
					// FIXME:  Append an `--` after the script to suppress interpreting other arguments!
					_arguments = append (_arguments, "-f")
				case "sed", "grep" :
					// FIXME:  Make it so that it doesn't accept other arguments!
					_arguments = append (_arguments, "-E", "-f")
				case "jq" :
					// FIXME:  Make it so that it doesn't accept other arguments!
					_arguments = append (_arguments, "-f")
				case "make" :
					// FIXME:  Append an `--` after the script to suppress interpreting other arguments!
					_arguments = append (_arguments, "-s", "-r", "-R", "-f")
				case "ninja" :
					// FIXME:  Append an `--` after the script to suppress interpreting other arguments!
					_arguments = append (_arguments, "-f")
				default :
					return errorf (0x505f52c6, "invalid interpreter for `%s`", _scriptlet.Label)
			}
		}
		
		if strings.IndexByte (_executable, os.PathSeparator) >= 0 {
			if _executable_0, _error := filepath.Abs (_executable); _error == nil {
				_executable = _executable_0
			} else {
				return errorw (0xda17c780, _error)
			}
		}
		
		_arguments = append (_arguments, _header[1:] ...)
		
		_scriptlet.Interpreter = _interpreter
		_scriptlet.InterpreterExecutable = _executable
		_scriptlet.InterpreterArguments = _arguments
		_scriptlet.InterpreterEnvironment = nil
	}
	
	return nil
}




func parseFromGenerator (_library *Library, _libraryUrl string, _libraryFingerprint string, _source *Scriptlet, _context *Context, _parseContext *parseContext) (*Error) {
//	logf ('s', 0xf75b04b5, "parsing `:: %s`...", _source.Label)
	if _, _data, _error := loadFromScriptlet (_libraryUrl, _libraryFingerprint, "", _source, _context); _error == nil {
		return parseFromData (_library, _data, _source.Source.Path, _context, _parseContext)
	} else {
		return _error
	}
}

func parseFromReplacer (_library *Library, _libraryUrl string, _libraryFingerprint string, _source *Scriptlet, _context *Context) (*Error) {
//	logf ('s', 0xc336e8bd, "parsing `:: %s`...", _source.Label)
	if _, _data, _error := loadFromScriptlet (_libraryUrl, _libraryFingerprint, _source.Interpreter, _source, _context); _error == nil {
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

func parseFromMenu (_library *Library, _source *Scriptlet, _context *Context) (*Error) {
	if _source.Kind != "menu-pending" {
		return errorf (0x6834ac93, "invalid state")
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
				case '^', '~' :
					_matcher = "+" + _matcher
				default :
					return errorf (0x11ca1466, "invalid menu mode `%s`", _matcher)
			}
			_mode := _matcher[:1]
			_matcher = _matcher[1:]
			if _matcher == "" {
				for _, _scriptlet := range _scriptlets {
					_labels = append (_labels, _scriptlet.Label)
					_scriptlet.Menus = append (_scriptlet.Menus, _mode + " " + _source.Label)
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
					return errorf (0xabf68b41, "invalid menu matcher `%s`", _matcher)
				}
				for _, _scriptlet := range _scriptlets {
					if _pattern.MatchString (_scriptlet.Label) {
						_labels = append (_labels, _scriptlet.Label)
						_scriptlet.Menus = append (_scriptlet.Menus, _mode + " " + _source.Label)
					}
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
		_library.LibraryFingerprint = NewFingerprinter () .StringWithLen (_library.LibraryFingerprint) .StringWithLen (_source.FingerprintData) .Build ()
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
							_include = false
							if _body == "" {
								if _label == "*" {
									_body = "*"
								} else {
									if strings.HasSuffix (_label, " ...") {
										_body = "+^ " + _label[:len (_label) - 4]
									} else if strings.HasSuffix (_label, " *") {
										_body = "*^ " + _label[:len (_label) - 2]
									} else {
										return errorf (0x84fa71b1, "invalid syntax (%d):  invalid menu label | %s", _lineIndex, _line)
									}
								}
							}
						default :
							return errorf (0xfba805b9, "invalid syntax (%d):  unknown scriptlet type | %s", _lineIndex, _line)
					}
					
					if _include {
						_includePath := ""
						if _includePath_0, _error := resolveRelativePath (_context.workspace, path.Dir (_sourcePath), _body); _error != nil {
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
									return errorf (0x16010e20, "invalid UTF-8")
								}
								_library.LibraryFingerprint = NewFingerprinter () .StringWithLen (_library.LibraryFingerprint) .StringWithLen (_includeSource.FingerprintData) .Build ()
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
					if _includePath_0, _error := resolveRelativePath (_context.workspace, path.Dir (_sourcePath), _includePath); _error != nil {
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
							return _error
						}
					}
					
				} else if strings.HasPrefix (_lineTrimmed, "&&__ ") {
					
					_includePath := _lineTrimmed[strings.IndexByte (_lineTrimmed, ' ') + 1:]
					if _includePath_0, _error := resolveRelativePath (_context.workspace, path.Dir (_sourcePath), _includePath); _error != nil {
						return _error
					} else {
						_includePath = _includePath_0
					}
					
					if !_disabled {
						if _includeSource, _error := fingerprintSource (_includePath); _error == nil {
							if _error := includeSource (_library, _includeSource); _error != nil {
								return _error
							}
							_library.LibraryFingerprint = NewFingerprinter () .StringWithLen (_library.LibraryFingerprint) .StringWithLen (_includeSource.FingerprintData) .Build ()
						} else {
							return _error
						}
					}
					
				} else if strings.HasPrefix (_lineTrimmed, "&&== ") {
					
					_descriptor := _lineTrimmed[strings.IndexByte (_lineTrimmed, ' ') + 1:]
					_kind := _descriptor[: strings.IndexByte (_descriptor, ' ')]
					_descriptor = _descriptor[len (_kind) + 1 :]
					_kind = strings.TrimSpace (_kind)
					_descriptor = strings.TrimSpace (_descriptor)
					
					if _kind == "" {
						return errorf (0xed68e4c3, "invalid syntax (%d):  empty statement | %s", _lineIndex, _line)
					}
					
					switch _kind {
						
						case "path" :
							
							if _descriptor == "" {
								return errorf (0x8c2e1bf8, "invalid syntax (%d):  empty statement path descriptor | %s", _lineIndex, _line)
							}
							_path := ""
							if _path_0, _error := resolveAbsolutePath (_context.workspace, path.Dir (_sourcePath), _descriptor); _error != nil {
								return _error
							} else {
								_path = _path_0
							}
							if _stat, _error := os.Stat (_path); _error == nil {
								if ! _stat.IsDir () {
									return errorf (0xae85ad7e, "invalid syntax (%d):  invalid statement path value | %s", _lineIndex, _line)
								}
							} else {
								return errorw (0x79069c08, _error)
							}
							_parseContext.scriptletContext.ExecutablePaths = append (_parseContext.scriptletContext.ExecutablePaths, _path)
							
						case "environment", "env", "environment-path", "env-path" :
							
							if _descriptor == "" {
								return errorf (0x7f049882, "invalid syntax (%d):  empty statement environment descriptor | %s", _lineIndex, _line)
							}
							_descriptor := strings.SplitN (_descriptor, " ", 2)
							if len (_descriptor) != 2 {
								return errorf (0xfc01ef6a, "invalid syntax (%d):  invalid statement environment descriptor | %s", _lineIndex, _line)
							}
							_name := strings.TrimSpace (_descriptor[0])
							_value := strings.TrimSpace (_descriptor[1])
							if _name == "" {
								return errorf (0x6bce31bb, "invalid syntax (%d):  empty statement environment key | %s", _lineIndex, _line)
							}
							if _, _exists := _parseContext.scriptletContext.Environment[_name]; _exists {
								return errorf (0x774b50de, "invalid syntax (%d):  duplicate statement environment key | %s", _lineIndex, _line)
							}
							if (_kind == "environment-path") || (_kind == "env-path") {
								if _value == "" {
									return errorf (0x2124d511, "invalid syntax (%d):  empty statement environment descriptor | %s", _lineIndex, _line)
								}
								if _path_0, _error := resolveAbsolutePath (_context.workspace, path.Dir (_sourcePath), _value); _error != nil {
									return _error
								} else {
									_value = _path_0
								}
							}
							_parseContext.scriptletContext.Environment[_name] = _value
							
						case "environment-exclude", "env-exclude" :
							if _descriptor == "" {
								return errorf (0x7f049882, "invalid syntax (%d):  empty statement environment descriptor | %s", _lineIndex, _line)
							}
							_descriptor := strings.Split (_descriptor, " ")
							for _, _name := range _descriptor {
								if _name == "" {
									continue
								}
								if _, _exists := _parseContext.scriptletContext.Environment[_name]; _exists {
									return errorf (0x0ed5990f, "invalid syntax (%d):  duplicate statement environment key | %s", _lineIndex, _line)
								}
								_parseContext.scriptletContext.Environment[_name] = ""
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
						Context : _parseContext.scriptletContext,
						ContextFingerprint : _parseContext.scriptletContext.Fingerprint,
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
				if _error := parseInterpreter (_library, _scriptlet, _context); _error != nil {
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




func loadFromSource (_library *Library, _source *Source, _context *Context) ([]byte, *Error) {
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


func loadFromSource_0 (_library *Library, _source *Source, _context *Context) (string, []byte, *Error) {
	
	if !_source.Executable {
		return loadFromFile (_source.Path)
		
	} else {
		
		_executable := _source.Path
		if ! path.IsAbs (_executable) {
			if _executable_0, _error := filepath.Abs (_executable); _error == nil {
				_executable = _executable_0
			} else {
				return "", nil, errorw (0x5d248eeb, _error)
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
		return "", nil, errorw (0xe1528040, _error)
	}
}


func loadFromStream (_stream io.Reader) (string, []byte, *Error) {
	if _data, _error := ioutil.ReadAll (_stream); _error == nil {
		_fingerprint := NewFingerprinter () .Bytes (_data) .Build ()
		return _fingerprint, _data, nil
	} else {
		return "", nil, errorw (0x42be0790, _error)
	}
}


func loadFromScriptlet (_libraryUrl string, _libraryFingerprint string, _interpreter string, _scriptlet *Scriptlet, _context *Context) (string, []byte, *Error) {
	
	var _command *exec.Cmd
	var _descriptors []int
	if _command_0, _descriptors_0, _error := prepareExecution (_libraryUrl, _libraryFingerprint, _interpreter, _scriptlet, false, _context); _error == nil {
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

