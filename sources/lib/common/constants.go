

package common


import "runtime"
import "strings"

import "golang.org/x/sys/unix"

import embedded "github.com/cipriancraciun/z-run/embedded"




var PROJECT_URL = "https://github.com/cipriancraciun/z-run"

var BUILD_TARGET string = "{unknown-target}"
var BUILD_TARGET_ARCH string = runtime.GOARCH
var BUILD_TARGET_OS string = runtime.GOOS
var BUILD_COMPILER_TYPE string = runtime.Compiler
var BUILD_COMPILER_VERSION string = runtime.Version ()
var BUILD_DEVELOPMENT bool = (BUILD_TARGET == "development")

var BUILD_VERSION string = strings.Trim (embedded.BuildVersion, "\n")
var BUILD_NUMBER string = strings.Trim (embedded.BuildNumber, "\n")
var BUILD_TIMESTAMP string = strings.Trim (embedded.BuildTimestamp, "\n")

var BUILD_GIT_HASH string = "{unknown-git-hash}"
var BUILD_SOURCES_HASH string = strings.Trim (embedded.BuildSourcesHash, "\n")

var UNAME_NODE string = "{unknown-node}"
var UNAME_SYSTEM string = "{unknown-system}"
var UNAME_RELEASE string = "{unknown-release}"
var UNAME_VERSION string = "{unknown-version}"
var UNAME_MACHINE string = "{unknown-machine}"
var UNAME_FINGERPRINT string = "{unknown-fingerprint}"




func init () () {
	
	var _uname unix.Utsname
	if _error := unix.Uname (&_uname); _error != nil {
		panic (_error)
	}
	
	_convert := func (_bytes []byte, _default string) (string) {
		_buffer := make ([]byte, 0, len (_bytes))
		for _, _byte := range _bytes {
			if _byte == 0 {
				break
			}
			_buffer = append (_buffer, byte (_byte))
		}
		if len (_buffer) > 0 {
			return string (_buffer)
		} else {
			return _default
		}
	}
	
	UNAME_NODE = _convert (_uname.Nodename[:], "{unknown-node}")
	UNAME_SYSTEM = _convert (_uname.Sysname[:], "{unknown-system}")
	UNAME_RELEASE = _convert (_uname.Release[:], "{unknown-release}")
	UNAME_VERSION = _convert (_uname.Version[:], "{unknown-version}")
	UNAME_MACHINE = _convert (_uname.Machine[:], "{unknown-machine}")
	
	if _index := strings.Index (UNAME_NODE, "."); _index != -1 {
		UNAME_NODE = UNAME_NODE[0 : _index]
	}
	if UNAME_NODE == "" {
		UNAME_NODE = "{unknown-node}"
	}
	
	UNAME_FINGERPRINT = FingerprintStringsQuick ("98ff673c677ffaeb481ce53a8deef977", UNAME_NODE, UNAME_SYSTEM, UNAME_RELEASE, UNAME_VERSION, UNAME_MACHINE) [0:32]
}

