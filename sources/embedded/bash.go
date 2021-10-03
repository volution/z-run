

package embedded


import _ "embed"




//go:embed bash/prolog_0.bash
var BashProlog0 string

//go:embed bash/prolog.bash
var BashProlog string

//go:embed bash/shell_rc.bash
var BashShellRc string

//go:embed bash/shell_functions.bash
var BashShellFunctions string

