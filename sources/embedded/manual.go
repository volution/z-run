

package embedded


import _ "embed"




//go:embed manual/z-run.txt
var ManualTxt string

//go:embed manual/z-run.html
var ManualHtml string

//go:embed manual/z-run.man
var ManualMan string

