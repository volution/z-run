

package zrun


import embedded "github.com/cipriancraciun/z-run/embedded"

import "runtime"
import "strings"




var BUILD_TARGET string = "{unknown-target}"
var BUILD_TARGET_ARCH string = runtime.GOARCH
var BUILD_TARGET_OS string = runtime.GOOS
var BUILD_COMPILER_TYPE string = runtime.Compiler
var BUILD_COMPILER_VERSION string = runtime.Version ()

var BUILD_VERSION string = strings.Trim (embedded.BuildVersion, "\n")
var BUILD_NUMBER string = strings.Trim (embedded.BuildNumber, "\n")
var BUILD_TIMESTAMP string = strings.Trim (embedded.BuildTimestamp, "\n")
var BUILD_SOURCES_MD5 string = strings.Trim (embedded.BuildSourcesMd5, "\n")
var BUILD_GIT_HASH string = "${unknown-git-hash}"

