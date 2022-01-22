

package mainlib


import "bytes"
import "fmt"
import "os"
import "strings"

import . "github.com/volution/z-run/lib/common"
import . "github.com/volution/z-run/embedded"




func InterceptMainSpecialFlags (_executableName string, _executable0 string, _executable string, _helpText string, _manualText string, _manualHtml string, _manualMan string, _environment map[string]string) (*Error) {
	
	if (len (os.Args) == 2) {
		_argument := os.Args[1]
		
		if (_argument == "--version") || (_argument == "-v") {
			_buffer := bytes.NewBuffer (nil)
			fmt.Fprintf (_buffer, "* tool          : %s\n", _executableName)
			fmt.Fprintf (_buffer, "* version       : %s\n", BUILD_VERSION)
			if _executable0 == _executable {
				fmt.Fprintf (_buffer, "* executable    : %s\n", _executable)
			} else {
				fmt.Fprintf (_buffer, "* executable    : %s\n", _executable)
				fmt.Fprintf (_buffer, "* executable-0  : %s\n", _executable0)
			}
			fmt.Fprintf (_buffer, "* build target  : %s, %s-%s, %s, %s\n", BUILD_TARGET, BUILD_TARGET_OS, BUILD_TARGET_ARCH, BUILD_COMPILER_VERSION, BUILD_COMPILER_TYPE)
			fmt.Fprintf (_buffer, "* build number  : %s, %s\n", BUILD_NUMBER, BUILD_TIMESTAMP)
			fmt.Fprintf (_buffer, "* code & issues : %s\n", PROJECT_URL)
			fmt.Fprintf (_buffer, "* sources git   : %s\n", BUILD_GIT_HASH)
			fmt.Fprintf (_buffer, "* sources hash  : %s\n", BUILD_SOURCES_HASH)
			fmt.Fprintf (_buffer, "* uname node    : %s\n", UNAME_NODE)
			fmt.Fprintf (_buffer, "* uname system  : %s, %s, %s\n", UNAME_SYSTEM, UNAME_RELEASE, UNAME_MACHINE)
			fmt.Fprintf (_buffer, "* uname hash    : %s\n", UNAME_FINGERPRINT)
			if _, _error := _buffer.WriteTo (os.Stdout); _error != nil {
				panic (AbortError (Errorw (0x36e1aa05, _error)))
			}
			panic (ExitMainSucceeded ())
		}
		
		if _argument == "--sources-md5" {
			if _, _error := os.Stdout.WriteString (BuildSourcesMd5); _error != nil {
				panic (AbortError (Errorw (0x7471032d, _error)))
			}
			panic (ExitMainSucceeded ())
		}
		
		if _argument == "--sources-cpio" {
			if _, _error := os.Stdout.Write (BuildSourcesCpioGz); _error != nil {
				panic (AbortError (Errorw (0x8034bf3e, _error)))
			}
			panic (ExitMainSucceeded ())
		}
		
		if (_argument == "-h") ||
				(_argument == "--help") ||
				(_argument == "--manual") ||
				(_argument == "--manual-text") || (_argument == "--manual-txt") ||
				(_argument == "--manual-html") ||
				(_argument == "--manual-man") ||
				(_argument == "--readme") ||
				(_argument == "--readme-text") || (_argument == "--readme-txt") ||
				(_argument == "--readme-html") {
			_replacements := map[string]string {
					"@{PROJECT_URL}" : PROJECT_URL,
					"@{BUILD_TARGET}" : BUILD_TARGET,
					"@{BUILD_TARGET_ARCH}" : BUILD_TARGET_ARCH,
					"@{BUILD_TARGET_OS}" : BUILD_TARGET_OS,
					"@{BUILD_COMPILER_TYPE}" : BUILD_COMPILER_TYPE,
					"@{BUILD_COMPILER_VERSION}" : BUILD_COMPILER_VERSION,
					"@{BUILD_DEVELOPMENT}" : fmt.Sprintf ("%s", BUILD_DEVELOPMENT),
					"@{BUILD_VERSION}" : BUILD_VERSION,
					"@{BUILD_NUMBER}" : BUILD_NUMBER,
					"@{BUILD_TIMESTAMP}" : BUILD_TIMESTAMP,
					"@{BUILD_GIT_HASH}" : BUILD_GIT_HASH,
					"@{BUILD_SOURCES_HASH}" : BUILD_SOURCES_HASH,
					"@{UNAME_NODE}" : UNAME_NODE,
					"@{UNAME_SYSTEM}" : UNAME_SYSTEM,
					"@{UNAME_RELEASE}" : UNAME_RELEASE,
					"@{UNAME_VERSION}" : UNAME_VERSION,
					"@{UNAME_MACHINE}" : UNAME_MACHINE,
					"@{UNAME_FINGERPRINT}" : UNAME_FINGERPRINT,
				}
			_manual := ""
			_useType := ""
			_useDecorations := false
			switch _argument {
				case "--help", "-h" :
					_manual = _helpText
					_useType = "text"
					_useDecorations = true
				case "--manual", "--manual-text", "--manual-txt" :
					_manual = _manualText
					_useType = "text"
					_useDecorations = true
				case "--manual-html" :
					_manual = _manualHtml
					_useType = "html"
				case "--manual-man" :
					_manual = _manualMan
					_useType = "man"
				case "--readme", "--readme-text", "--readme-txt" :
					_manual = ReadmeTxt
					_useType = "text"
					_useDecorations = true
				case "--readme-html" :
					_manual = ReadmeHtml
					_useType = "html"
				default :
					panic (0x41b79a1d)
			}
			if _manual != "__custom__" {
				if _manual == "" {
					panic (AbortError (Errorf (0x7f11c1ac, "manual not available")))
				}
				_chunks := make ([]string, 0, 8)
				if _useDecorations {
					_chunks = append (_chunks, HelpHeader, _manual, HelpFooter)
				} else {
					_chunks = append (_chunks, _manual)
				}
				for _index := range _chunks {
					for _key, _replacement := range _replacements {
						_chunks[_index] = strings.ReplaceAll (_chunks[_index], _key, _replacement)
					}
				}
				_buffer := bytes.NewBuffer (nil)
				for _, _chunk := range _chunks {
					_buffer.WriteString (_chunk)
				}
				switch _useType {
					case "text" :
						if _error := ExecPagerPerhaps ("less", []string {"-f", "@{PATH}"}, "@{PATH}", _environment, _buffer.Bytes ()); _error != nil {
							panic (AbortError (_error))
						}
					case "man" :
						if _error := ExecPagerPerhaps ("man", []string {"@{PATH}"}, "@{PATH}", _environment, _buffer.Bytes ()); _error != nil {
							panic (AbortError (_error))
						}
					case "html" :
						// if _error := ExecPagerPerhaps ("w3m", []string {"-T", "text/html", "-I", "UTF-8", "@{PATH}"}, "@{PATH}", _environment, _buffer.Bytes ()); _error != nil {
						// if _error := ExecPagerPerhaps ("links", []string {"-force-html", "@{PATH}"}, "@{PATH}", _environment, _buffer.Bytes ()); _error != nil {
						if _error := ExecPagerPerhaps ("lynx", []string {"-force-html", "@{PATH}"}, "@{PATH}", _environment, _buffer.Bytes ()); _error != nil {
							panic (AbortError (_error))
						}
				}
				panic (0xacd19d57)
			}
		}
	}
	
	return nil
}

