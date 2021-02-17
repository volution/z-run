

package zrun


import "runtime"




var BUILD_TARGET string = "{unknown-target}"
var BUILD_TARGET_ARCH string = runtime.GOARCH
var BUILD_TARGET_OS string = runtime.GOOS
var BUILD_COMPILER_TYPE string = runtime.Compiler
var BUILD_COMPILER_VERSION string = runtime.Version ()

var BUILD_VERSION string = "{unknown-version}"
var BUILD_NUMBER string = "{unknown-build}"

var BUILD_COMMIT string = "{unknown-commit}"
var BUILD_TIMESTAMP string = "{unknown-timestamp}"
var BUILD_SOURCES_MD5 string = "{unknown-sources-md5}"
var BUILD_GIT_HASH string = "${unknown-git-hash}"

