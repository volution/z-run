

package embedded


import _ "embed"




//go:embed go_prolog.go0
var GoProlog string

//go:embed python3_prolog.py
var Python3Prolog string

//go:embed bash_prolog_0.bash
var BashProlog0 string

//go:embed bash_prolog.bash
var BashProlog string

//go:embed bash_shell_rc.bash
var BashShellRc string

//go:embed bash_shell_functions.bash
var BashShellFunctions string


//go:embed manual.txt
var ManualTxt string


//go:embed build-version.txt
var BuildVersion string

//go:embed build-number.txt
var BuildNumber string

//go:embed build-timestamp.txt
var BuildTimestamp string

//go:embed build-sources-md5.txt
var BuildSourcesMd5 string

