
package zrun


import "io"
import "os"
import "strings"
import "text/template"




func templateMain () (*Error) {
	
	if len (os.Args) < 2 {
		return errorf (0x47c3f9f1, "invalid arguments")
	}
	
	_source := strings.Builder {}
	if _stream, _error := os.Open (os.Args[1]); _error == nil {
		defer _stream.Close ()
		if _, _error := io.Copy (&_source, _stream); _error != nil {
			return errorw (0x693128fe, _error)
		}
	}
	
	_template := template.New ("z-run")
	if _, _error := _template.Parse (_source.String ()); _error != nil {
		return errorw (0xfd33768b, _error)
	}
	
	_data := map[string]interface{} {
			"arguments" : os.Args[2:],
		}
	
	if _error := _template.Execute (os.Stdout, _data); _error != nil {
		return errorw (0x23ce8919, _error)
	}
	
	os.Exit (0)
	panic (0x8e448279)
}

