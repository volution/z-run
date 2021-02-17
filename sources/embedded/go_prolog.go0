// ################################################################################
// ################################################################################


import (
	_Z_fmt "fmt"
	_Z_os "os"
)


type _Z struct {}

var Z *_Z = _Z_new ()


func _Z_new () (*_Z) {
	return & _Z {}
}

func (_z *_Z) Exit (_code uint8) () {
	_Z_os.Exit (int (_code))
}

func (_z *_Z) LogNotice (_code uint32, _format string, _arguments ... interface{}) () {
	_z._log ("ii", _code, _format, _arguments)
}

func (_z * _Z) _log (_slug string, _code uint32, _format string, _arguments []interface{}) () {
	_pid := _Z_os.Getpid ()
	_message := _Z_fmt.Sprintf (_format, _arguments ...)
	_Z_fmt.Fprintf (_Z_os.Stderr, "[z-run:%08d] [%s] [%08x]  %s", _pid, _slug, _code, _message)
}


// ################################################################################
// ################################################################################
