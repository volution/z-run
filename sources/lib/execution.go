

package zrun


import "bytes"
import "fmt"
import "os"
import "os/exec"
import "strings"
import "syscall"

import "golang.org/x/sys/unix"




func prepareExecution (_libraryUrl string, _libraryFingerprint string, _interpreter string, _scriptlet *Scriptlet, _includeArguments bool, _context *Context) (*exec.Cmd, []int, *Error) {
	
	var _interpreterExecutable string
	var _interpreterArguments []string = make ([]string, 0, len (_context.cleanArguments) + 16)
	var _interpreterAllowsArguments = false
	
	if _interpreter == "" {
		_interpreter = _scriptlet.Interpreter
	}
	
	
	switch _interpreter {
		
		case "<exec>", "<bash>" :
			_interpreterAllowsArguments = true
		
		case "<print>" :
			_interpreterAllowsArguments = false
		
		case "<template>" :
			_interpreterAllowsArguments = true
		
		case "<menu>" :
			_interpreterAllowsArguments = false
		
		default :
			return nil, nil, errorf (0x0873f2db, "unknown scriptlet interpreter `%s` for `%s`", _interpreter, _scriptlet.Label)
	}
	
	if _includeArguments && (len (_context.cleanArguments) > 0) && !_interpreterAllowsArguments {
		return nil, nil, errorf (0x4ef9e048, "unexpected arguments")
	}
	
	var _interpreterScriptInput int
	var _interpreterScriptOutput *os.File
	var _interpreterScriptDescriptors [2]int
	var _interpreterScriptUnused = false
	if _error := syscall.Pipe (_interpreterScriptDescriptors[:]); _error == nil {
		if _, _error := unix.FcntlInt (uintptr (_interpreterScriptDescriptors[1]), unix.F_SETPIPE_SZ, 1048576); _error != nil {
			logf ('w', 0x4d3414c8, "failed increasing pipe buffer size;  ignoring!")
		}
		_interpreterScriptInput = _interpreterScriptDescriptors[0]
		_interpreterScriptOutput = os.NewFile (uintptr (_interpreterScriptDescriptors[1]), "")
	} else {
		return nil, nil, errorw (0xece645ff, _error)
	}
	
	_interpreterScriptBuffer := bytes.NewBuffer (nil)
	_interpreterScriptBuffer.Grow (128 * 1024)
	
	switch _interpreter {
		
		case "<exec>" :
			_interpreterExecutable = _scriptlet.InterpreterExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					_scriptlet.InterpreterExecutable,
				)
			_interpreterArguments = append (
					_interpreterArguments,
					_scriptlet.InterpreterArguments ...
				)
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (_scriptlet.Body)
		
		case "<bash>" :
			_interpreterExecutable = "/bin/bash"
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:bash] [%s]", _scriptlet.Label),
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
		
		case "<template>" :
			_interpreterExecutable = _context.selfExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					"[z-run:template]",
					fmt.Sprintf (":: %s", _scriptlet.Label),
				)
			_interpreterScriptUnused = true
		
		case "<template-raw>" :
			_interpreterExecutable = _context.selfExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:template-raw] [%s]", _scriptlet.Label),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (_scriptlet.Body)
		
		case "<menu>" :
			_interpreterExecutable = _context.selfExecutable
			_interpreterArguments = append (
					_interpreterArguments,
					fmt.Sprintf ("[z-run:menu] [%s]", _scriptlet.Label),
					fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput),
				)
			_interpreterScriptBuffer.WriteString (_scriptlet.Body)
		
		default :
			panic (0xe95f68a0)
	}
	
	
	var _descriptors []int
	if ! _interpreterScriptUnused {
		
//		logf ('d', 0xedfcf88b, "\n----------\n%s----------\n", _interpreterScriptBuffer.Bytes ())
		
		if _, _error := _interpreterScriptBuffer.WriteTo (_interpreterScriptOutput); _error == nil {
			_interpreterScriptOutput.Close ()
		} else {
			syscall.Close (_interpreterScriptInput)
			_interpreterScriptOutput.Close ()
			return nil, nil, errorw (0xf789ed3f, _error)
		}
		
		_descriptors = []int {
				_interpreterScriptInput,
			}
		
	} else {
		syscall.Close (_interpreterScriptInput)
		_interpreterScriptOutput.Close ()
	}
	
	if _includeArguments {
		_interpreterArguments = append (_interpreterArguments, _context.cleanArguments ...)
	}
	
	_interpreterEnvironment := processEnvironment (_context, map[string]string {
			"ZRUN_EXECUTABLE" : _context.selfExecutable,
			"ZRUN_WORKSPACE" : _context.workspace,
			"ZRUN_LIBRARY_CACHE" : _libraryUrl,
			"ZRUN_FINGERPRINT" : _libraryFingerprint,
		})
	
	if strings.IndexByte (_interpreterExecutable, os.PathSeparator) < 0 {
		if _path, _error := exec.LookPath (_interpreterExecutable); _error == nil {
			_interpreterExecutable = _path
		} else {
			syscall.Close (_interpreterScriptInput)
			_interpreterScriptOutput.Close ()
			return nil, nil, errorw (0xd8f4497c, _error)
		}
	}
	
//	logf ('d', 0xcc6d38ba, "command: `%#v` `%#v`", _interpreterExecutable, _interpreterArguments)
	
	_command := & exec.Cmd {
			Path : _interpreterExecutable,
			Args : _interpreterArguments,
			Env : _interpreterEnvironment,
			Dir : _context.workspace,
		}
	
	return _command, _descriptors, nil
}

