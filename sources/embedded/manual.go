

package embedded


import _ "embed"




//go:embed manual/readme.txt
var ReadmeTxt string

//go:embed manual/readme.html
var ReadmeHtml string


//go:embed manual/z-run--help.txt
var ZrunHelpTxt string

//go:embed manual/z-run--manual.txt
var ZrunManualTxt string

//go:embed manual/z-run--manual.html
var ZrunManualHtml string

//go:embed manual/z-run--manual.man
var ZrunManualMan string


//go:embed manual/help--header.txt
var HelpHeader string

//go:embed manual/help--footer.txt
var HelpFooter string

