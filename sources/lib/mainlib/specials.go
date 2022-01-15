

package mainlib


import "fmt"
import "os"

import embedded "github.com/cipriancraciun/z-run/embedded"

import . "github.com/cipriancraciun/z-run/lib/common"




func InterceptMainSpecialFlags (_executableName string, _executable0 string, _executable string, _helpText string, _manualText string, _manualHtml string, _manualMan string) (*Error) {
	
	if (len (os.Args) == 2) {
		_argument := os.Args[1]
		
		if (_argument == "--version") || (_argument == "-v") {
			
			fmt.Fprintf (os.Stdout, "* tool          : %s\n", _executableName)
			fmt.Fprintf (os.Stdout, "* version       : %s\n", BUILD_VERSION)
			if _executable0 == _executable {
				fmt.Fprintf (os.Stdout, "* executable    : %s\n", _executable)
			} else {
				fmt.Fprintf (os.Stdout, "* executable    : %s\n", _executable)
				fmt.Fprintf (os.Stdout, "* executable-0  : %s\n", _executable0)
			}
			fmt.Fprintf (os.Stdout, "* build target  : %s, %s-%s, %s, %s\n", BUILD_TARGET, BUILD_TARGET_OS, BUILD_TARGET_ARCH, BUILD_COMPILER_VERSION, BUILD_COMPILER_TYPE)
			fmt.Fprintf (os.Stdout, "* build number  : %s, %s\n", BUILD_NUMBER, BUILD_TIMESTAMP)
			fmt.Fprintf (os.Stdout, "* code & issues : %s\n", PROJECT_URL)
			fmt.Fprintf (os.Stdout, "* sources git   : %s\n", BUILD_GIT_HASH)
			fmt.Fprintf (os.Stdout, "* sources hash  : %s\n", BUILD_SOURCES_HASH)
			fmt.Fprintf (os.Stdout, "* uname node    : %s\n", UNAME_NODE)
			fmt.Fprintf (os.Stdout, "* uname system  : %s, %s, %s\n", UNAME_SYSTEM, UNAME_RELEASE, UNAME_MACHINE)
			
			panic (ExitMainSucceeded ())
		}
		
		if _argument == "--sources-md5" {
			if _, _error := os.Stdout.WriteString (embedded.BuildSourcesMd5); _error != nil {
				panic (AbortError (Errorw (0x7471032d, _error)))
			}
			panic (ExitMainSucceeded ())
		}
		
		if _argument == "--sources-cpio" {
			if _, _error := os.Stdout.Write (embedded.BuildSourcesCpioGz); _error != nil {
				panic (AbortError (Errorw (0x8034bf3e, _error)))
			}
			panic (ExitMainSucceeded ())
		}
		
		if (_argument == "--help") || (_argument == "-h") || (_argument == "--manual") || (_argument == "--manual-text") || (_argument == "--manual-html") || (_argument == "--manual-man") {
			_manual := ""
			switch _argument {
				case "--help", "-h" :
					_manual = _helpText
				case "--manual", "--manual-text" :
					_manual = _manualText
				case "--manual-html" :
					_manual = _manualHtml
				case "--manual-man" :
					_manual = _manualMan
				default :
					panic (0x41b79a1d)
			}
			if _manual != "__custom__" {
				if _manual == "" {
					panic (AbortError (Errorf (0x7f11c1ac, "manual not available")))
				}
				fmt.Fprint (os.Stdout, _manual)
				panic (ExitMainSucceeded ())
			}
		}
	}
	
	return nil
}

