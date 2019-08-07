

package lib


import "bytes"
import "fmt"
import "os"
import "os/exec"
import "syscall"




func prepareExecution (_library LibraryStore, _scriptlet *Scriptlet, _context *Context) (*exec.Cmd, []int, error) {
	
	var _interpreterExecutable string
	var _interpreterArguments []string = make ([]string, 0, len (_context.cleanArguments) + 16)
	
	var _interpreterScriptInput int
	var _interpreterScriptOutput *os.File
	var _interpreterScriptDescriptors [2]int
	if _error := syscall.Pipe (_interpreterScriptDescriptors[:]); _error == nil {
		_interpreterScriptInput = _interpreterScriptDescriptors[0]
		_interpreterScriptOutput = os.NewFile (uintptr (_interpreterScriptDescriptors[1]), "")
	} else {
		return nil, nil, _error
	}
	
	_interpreterScriptBuffer := bytes.NewBuffer (nil)
	_interpreterScriptBuffer.Grow (128 * 1024)
	
	switch _scriptlet.Interpreter {
		
		case "<shell>" :
			_interpreterExecutable = "/bin/bash"
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:shell] [%s]", _scriptlet.Label),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (
					fmt.Sprintf (
`#!/dev/null
set -e -E -u -o pipefail -o noclobber -o noglob +o braceexpand || exit -- 1
trap 'printf -- "[ee] failed: %%s\n" "${BASH_COMMAND}" >&2' ERR || exit -- 1
BASH_ARGV0='z-run'
ZRUN=( %s )
X_RUN=( %s )
exec %d<&-

`,
							_context.selfExecutable,
							_context.selfExecutable,
							_interpreterScriptInput,
						))
			_interpreterScriptBuffer.WriteString (_scriptlet.Body)
		
		case "<print>" :
			_interpreterExecutable = "/bin/cat"
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:print] [%s]", _scriptlet.Label),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (_scriptlet.Body)
		
		default :
			syscall.Close (_interpreterScriptInput)
			_interpreterScriptOutput.Close ()
			return nil, nil, errorf (0x0873f2db, "unknown scriptlet interpreter `%s` for `%s`", _scriptlet.Interpreter, _scriptlet.Label)
	}
	
//	logf ('d', 0xedfcf88b, "\n----------\n%s----------\n", _interpreterScriptBuffer.Bytes ())
	
	if _, _error := _interpreterScriptBuffer.WriteTo (_interpreterScriptOutput); _error == nil {
		_interpreterScriptOutput.Close ()
	} else {
		syscall.Close (_interpreterScriptInput)
		_interpreterScriptOutput.Close ()
		return nil, nil, _error
	}
	
	_interpreterArguments = append (_interpreterArguments, _context.cleanArguments ...)
	_interpreterEnvironment := processEnvironment (_context, map[string]string {
			"ZRUN_EXECUTABLE" : _context.selfExecutable,
			"ZRUN_LIBRARY" : _library.Url (),
		})
	
	_descriptors := []int {
			_interpreterScriptInput,
		}
	
	_command := & exec.Cmd {
			Path : _interpreterExecutable,
			Args : _interpreterArguments,
			Env : _interpreterEnvironment,
		}
	
	return _command, _descriptors, nil
}

