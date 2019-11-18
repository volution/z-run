

package zrun


import "fmt"
import "log"
import "os"




type Error struct {
	code uint32
	message string
	error error
}




func logf (_slug rune, _code uint32, _format string, _arguments ... interface{}) () {
	_pid := os.Getpid ()
	_message := fmt.Sprintf (_format, _arguments ...)
	_prefix := fmt.Sprintf ("[%08d] [%c%c] [%08x]  ", _pid, _slug, _slug, _code)
	log.Print (_prefix + _message + "\n")
}

func logError (_slug rune, _error *Error) () {
	logErrorf (_slug, 0x55d59c80, _error, "unexpected error encountered!")
}

func logErrorf (_slug rune, _code uint32, _error *Error, _format string, _arguments ... interface{}) () {
	_pid := os.Getpid ()
	if (_format != "") || (len (_arguments) != 0) {
		logf (_slug, _code, _format, _arguments ...)
	}
	if _error != nil {
		if _error.message != "" {
			log.Printf ("[%08d] [%c%c] [%08x]  %s\n", _pid, _slug, _slug, _error.code, _error.message)
		} else {
			log.Printf ("[%08d] [%c%c] [%08x]  %s\n", _pid, _slug, _slug, _error.code, "unexpected error encountered!")
		}
		if _error.error != nil {
			log.Printf ("[%08d] [%c%c] [%08x]  %s\n", _pid, _slug, _slug, _error.code, _error.error.Error ())
			log.Printf ("[%08d] [%c%c] [%08x]  %#v\n", _pid, _slug, _slug, _error.code, _error.error)
		}
	}
}


func abortError (_error *Error) (*Error) {
	return abortErrorf (_error, _error.code, "")
}

func abortErrorf (_error *Error, _code uint32, _format string, _arguments ... interface{}) (*Error) {
	logErrorf ('!', _code, _error, _format, _arguments ...)
//	logf ('!', 0xb7a5fb86, "aborting!")
	os.Exit (1)
	panic (0xa235deea)
}


func errorf (_code uint32, _format string, _arguments ... interface{}) (*Error) {
	_message := fmt.Sprintf (_format, _arguments ...)
	return & Error {
			code : _code,
			message : _message,
			error : nil,
		}
}

func errorw (_code uint32, _error error) (*Error) {
	if _code == 0 {
		panic (0xa4ddfd33)
	}
	return & Error {
			code : _code,
			message : "",
			error : _error,
		}
}

