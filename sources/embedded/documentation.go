

package embedded


import _ "embed"




//go:embed documentation/readme.txt
var ReadmeTxt string

//go:embed documentation/readme.html
var ReadmeHtml string


//go:embed documentation/z-run--help.txt
var ZrunHelpTxt string

//go:embed documentation/z-run--manual.txt
var ZrunManualTxt string

//go:embed documentation/z-run--manual.html
var ZrunManualHtml string

//go:embed documentation/z-run--manual.man
var ZrunManualMan string


//go:embed documentation/help--header.txt
var HelpHeader string

//go:embed documentation/help--footer.txt
var HelpFooter string




//go:embed documentation/sbom.txt
var SbomTxt string

//go:embed documentation/sbom.html
var SbomHtml string

//go:embed documentation/sbom.json
var SbomJson string

