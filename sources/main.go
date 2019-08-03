

package main


import "errors"
import "fmt"
import "log"
import "os"
import "regexp"




func main_0 () (error) {
	logf ('i', 0xfe68754c, "hello!")
	return nil
}




func main () () {
	
	log.SetFlags (0)
	
	if _error := main_0 (); _error == nil {
		os.Exit (0)
	} else {
		panic (abortError (_error))
	}
}




func logf (_slug rune, _code uint32, _format string, _arguments ... interface{}) () {
	_message := fmt.Sprintf (_format, _arguments ...)
	_prefix := fmt.Sprintf ("[%c%c]  [%08x]  ", _slug, _slug, _code)
	log.Print (_prefix + _message + "\n")
}

func logError (_slug rune, _error error) () {
	logErrorf (_slug, 0x55d59c80, _error, "unexpected error encountered!")
}

func logErrorf (_slug rune, _code uint32, _error error, _format string, _arguments ... interface{}) () {
	logf (_slug, _code, _format, _arguments ...)
	if _error != nil {
		_errorString := _error.Error ()
		_errorRegexp := regexp.MustCompile (`^\[[0-9a-f]{8}\]  [^\n]+$`)
		if _matches := _errorRegexp.MatchString (_errorString); _matches {
			log.Printf ("[%c%c]  %s\n", _slug, _slug, _errorString)
		} else {
			log.Printf ("[%c%c]  [%08x]  %s\n", _slug, _slug, 0xda900de1, _errorString)
			log.Printf ("[%c%c]  [%08x]  %#v\n", _slug, _slug, 0x4fb5d56d, _error)
		}
	}
}


func abortf (_code uint32, _format string, _arguments ... interface{}) (error) {
	return abortErrorf (nil, _code, _format, _arguments ...)
}

func abortError (_error error) (error) {
	return abortErrorf (_error, 0xe6ed2b0f, "unexpected error encountered!")
}

func abortErrorf (_error error, _code uint32, _format string, _arguments ... interface{}) (error) {
	logErrorf ('!', _code, _error, _format, _arguments ...)
	logf ('!', 0xb7a5fb86, "aborting!")
	os.Exit (1)
	panic (0xa235deea)
}


func errorf (_code uint32, _format string, _arguments ... interface{}) (error) {
	_message := fmt.Sprintf (_format, _arguments ...)
	_prefix := fmt.Sprintf ("[%08x]  ", _code)
	return errors.New (_prefix + _message)
}

