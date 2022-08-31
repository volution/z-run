

package mainlib


import "fmt"
import "os"
import "os/exec"
import "strings"
import "syscall"

import . "github.com/volution/z-run/lib/common"




func ExecPagerPerhaps (_executable string, _arguments []string, _argumentsPathReplacement string, _environment map[string]string, _data []byte) (*Error) {
	
	if IsStdoutTerminal () && IsStderrTerminal () {
		return ExecPager (_executable, _arguments, _argumentsPathReplacement, _environment, _data)
	}
	
	if _, _error := os.Stdout.Write (_data); _error != nil {
		return Errorw (0x08681475, _error)
	}
	
	panic (ExitMainSucceeded ())
}




func ExecPager (_executable string, _arguments []string, _argumentsPathReplacement string, _environment map[string]string, _data []byte) (*Error) {
	
	_input, _output, _error := CreatePipe (len (_data), "/tmp")
	if _error != nil {
		panic (AbortError (_error))
	}
	
	if _, _error := _output.Write ([]byte (_data)); _error != nil {
		panic (AbortError (Errorw (0xa15cc434, _error)))
	}
	if _error := _output.Close (); _error != nil {
		panic (AbortError (Errorw (0x03fca643, _error)))
	}
	
	var _executableActual string
	if _executableActual_0, _error := exec.LookPath (_executable); _error == nil {
		_executableActual = _executableActual_0
	} else if _executableActual_0, _error := ResolveExecutable (_executable, []string { "/usr/local/bin", "/usr/bin", "/bin" }); _error == nil {
		_executableActual = _executableActual_0
	} else {
		if _, _error := os.Stdout.Write (_data); _error != nil {
			return Errorw (0xa86c9e4c, _error)
		}
		panic (ExitMainSucceeded ())
	}
	
	_inputPath := fmt.Sprintf ("/dev/fd/%d", _input)
	
	_argumentsActual := make ([]string, 0, 1 + len (_arguments))
	_argumentsActual = append (_argumentsActual, _executableActual)
	for _, _argument := range _arguments {
		if _argumentsPathReplacement != "" {
			_argument = strings.ReplaceAll (_argument, _argumentsPathReplacement, _inputPath)
		}
		_argumentsActual = append (_argumentsActual, _argument)
	}
	
	_environmentList := EnvironmentMapToList (_environment)
	
	if _error := Syscall_Dup2or3 (int (os.Stderr.Fd ()), int (os.Stdin.Fd ())); _error != nil {
		return Errorw (0x99bfeed4, _error)
	}
	if _error := Syscall_Dup2or3 (int (os.Stderr.Fd ()), int (os.Stdout.Fd ())); _error != nil {
		return Errorw (0x3cc9c54c, _error)
	}
	
	if _error := syscall.Exec (_executableActual, _argumentsActual, _environmentList); _error != nil {
		return Errorw (0xa5c95681, _error)
	}
	
	panic (0xf4813cc2)
}

