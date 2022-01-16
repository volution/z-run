

package embedded


import _ "embed"




//go:embed manual/z-run--help.txt
var ZrunHelpTxt string

//go:embed manual/z-run--manual.txt
var ZrunManualTxt string

//go:embed manual/z-run--manual.html
var ZrunManualHtml string

//go:embed manual/z-run--manual.man
var ZrunManualMan string

