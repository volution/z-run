

package execve


import "encoding/base64"
import "encoding/json"
import "os"
import "os/exec"
import "syscall"
import "sort"
import "strings"

import . "github.com/volution/z-run/lib/mainlib"
import . "github.com/volution/z-run/lib/common"




type ExecveMainFlags struct {
	
	Invoke ExecveMainInvokeFlags `group:"Invoke Options"`
}


type ExecveMainInvokeFlags struct {
	
	Json *string `long:"descriptor" value-name:"{json}" description:"JSON context;"`
	JsonEncoded *string `long:"descriptor-encoded" value-name:"{base64url-json}" description:"Base64-URL encoded JSON context;"`
}




type ExecveContext struct {
	Executable *string `json:"executable,omitempty"`
	Argument0 *string `json:"argument0,omitempty"`
	Arguments []string `json:"arguments,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	ExecutablePaths []string `json:"executable-paths,omitempty"`
	Terminal *string `json:"terminal,omitempty"`
	Workspace *string `json:"workspace,omitempty"`
	Stdin *string `json:"stdin,omitempty"`
	Stdout *string `json:"stdout,omitempty"`
	Stderr *string `json:"stderr,omitempty"`
}




func ExecveMain (_arguments []string, _environment map[string]string) (*Error) {
	
	_flags := & ExecveMainFlags {}
	
	if _error := ResolveMainFlags ("z-execve", _arguments, _environment, _flags, os.Stderr); _error != nil {
		return _error
	}
	
	return ExecveMainWithFlags (_flags, _environment)
}




func ExecveMainWithFlags (_flags *ExecveMainFlags, _environmentMap map[string]string) (*Error) {
	
	_contextJson := FlagStringOrDefault (_flags.Invoke.Json, "")
	_contextJsonEncoded := FlagStringOrDefault (_flags.Invoke.JsonEncoded, "")
	
	if (_contextJson != "") && (_contextJsonEncoded != "") {
		return Errorf (0xa99ff11a, "invalid arguments:  JSON and JSON-encoded are mutually exclusive!")
	}
	if (_contextJson == "") && (_contextJsonEncoded == "") {
		return Errorf (0xf0cc72cd, "invalid arguments:  expected context!")
	}
	
	var _context ExecveContext
	var _contextJsonData []byte
	if _contextJson != "" {
		_contextJsonData = []byte (_contextJson)
	} else {
		if _data, _error := base64.RawURLEncoding.DecodeString (_contextJsonEncoded); _error == nil {
			_contextJsonData = _data
		} else {
			return Errorw (0x1f614834, _error)
		}
	}
	if _error := json.Unmarshal (_contextJsonData, &_context); _error == nil {
		// NOP
	} else {
		return Errorw (0x8c9e0d38, _error)
	}
	
	if _context.Executable == nil {
		return Errorf (0x1e3aeced, "invalid context:  missing executable")
	}
	_executable := *_context.Executable
	if _context.ExecutablePaths != nil {
		if _executable_0, _error := ResolveExecutable (_executable, _context.ExecutablePaths); _error == nil {
			_executable = _executable_0
		} else {
			return _error
		}
	} else if _executable_0, _error := exec.LookPath (_executable); _error == nil {
		_executable = _executable_0
	} else if _executable_0, _error := ResolveExecutable (_executable, []string { "/usr/local/bin", "/usr/bin", "/bin" }); _error == nil {
		_executable = _executable_0
	} else {
		return Errorf (0xf2af1661, "unresolved executable `%s`", _executable)
	}
	
	_arguments := make ([]string, 0, len (_context.Arguments) + 1)
	if _context.Argument0 != nil {
		_arguments = append (_arguments, *_context.Argument0)
	} else {
		_arguments = append (_arguments, _executable)
	}
	_arguments = append (_arguments, _context.Arguments ...)
	
	var _environmentArray []string = make ([]string, 0, len (_context.Environment) + 2)
	_pathInjected := false
	if _context.ExecutablePaths != nil {
		_executablePaths := strings.Join (_context.ExecutablePaths, string (os.PathListSeparator))
		_environmentArray = append (_environmentArray, "PATH=" + _executablePaths)
		_pathInjected = true
	}
	_termInjected := false
	if _context.Terminal != nil {
		_environmentArray = append (_environmentArray, "TERM=" + *_context.Terminal)
		_termInjected = true
	}
	
	if _context.Environment != nil {
		_environmentMap = _context.Environment
	}
	for _name, _value := range _environmentMap {
		if (_name == "PATH") && _pathInjected {
			continue
		}
		if (_name == "TERM") && _termInjected {
			continue
		}
		_variable := _name + "=" + _value
		_environmentArray = append (_environmentArray, _variable)
	}
	sort.Strings (_environmentArray)
	
	_workspace := ""
	if _context.Workspace != nil {
		_workspace = *_context.Workspace
	}
	
	if _workspace != "" {
		if _error := syscall.Chdir (_workspace); _error != nil {
			panic (AbortError (Errorf (0xd75304e8, "failed to chdir `%s`  //  %v", _workspace, _error)))
		}
	}
	
	if _context.Stdin != nil {
		switch *_context.Stdin {
			case "/dev/stderr" :
				if _error := Syscall_Dup2or3 (int (os.Stderr.Fd ()), int (os.Stdin.Fd ())); _error != nil {
					return Errorw (0x06fee7e8, _error)
				}
			default :
				if _fd, _error := syscall.Open (*_context.Stdin, os.O_RDONLY, 0); _error == nil {
					if _error := Syscall_Dup2or3 (_fd, int (os.Stdin.Fd ())); _error != nil {
						return Errorw (0x5c57325c, _error)
					}
				} else {
					return Errorf (0x50fb349e, "failed opening `%s` for reading  //  `%v`", *_context.Stdin, _error)
				}
		}
	}
	if _context.Stdout != nil {
		switch *_context.Stdout {
			case "/dev/stderr" :
				if _error := Syscall_Dup2or3 (int (os.Stderr.Fd ()), int (os.Stdout.Fd ())); _error != nil {
					return Errorw (0xb0e35a05, _error)
				}
			default :
				if _fd, _error := syscall.Open (*_context.Stdout, os.O_WRONLY | os.O_APPEND, 0); _error == nil {
					if _error := Syscall_Dup2or3 (_fd, int (os.Stdout.Fd ())); _error != nil {
						return Errorw (0x7c9e7923, _error)
					}
				} else {
					return Errorf (0xec2731bc, "failed opening `%s` for writing  //  `%v`", *_context.Stdout, _error)
				}
		}
	}
	if _context.Stderr != nil {
		switch *_context.Stderr {
			default :
				if _fd, _error := syscall.Open (*_context.Stderr, os.O_WRONLY | os.O_APPEND, 0); _error == nil {
					if _error := Syscall_Dup2or3 (_fd, int (os.Stderr.Fd ())); _error != nil {
						return Errorw (0x1cbddf68, _error)
					}
				} else {
					return Errorf (0x3c801323, "failed opening `%s` for writing  //  `%v`", *_context.Stderr, _error)
				}
		}
	}
	
	if _error := syscall.Exec (_executable, _arguments, _environmentArray); _error != nil {
		panic (AbortError (Errorf (0xb2b6ff3e, "failed to exec `%s`  //  %v", _executable, _error)))
	}
	
	panic (AbortUnreachable (0x08d52120))
}


