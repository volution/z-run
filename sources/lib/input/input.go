

package input


import "fmt"
import "io"
import "os"
import "syscall"

import "github.com/peterh/liner"

import . "github.com/cipriancraciun/z-run/lib/mainlib"
import . "github.com/cipriancraciun/z-run/lib/common"




type InputMainFlags struct {
	
	Message *string `long:"message" short:"m" value-name:"{message}"`
	Prompt *string `long:"prompt" short:"p" value-name:"{prompt}"`
	Sensitive *bool `long:"sensitive" short:"s"`
}




func InputMain (_arguments []string, _environment map[string]string) (*Error) {
	
	_flags := & InputMainFlags {}
	
	if _error := ResolveMainFlags ("z-input", _arguments, _environment, _flags, os.Stderr); _error != nil {
		return _error
	}
	
	return InputMainWithFlags (_flags)
}




func InputMainWithFlags (_flags *InputMainFlags) (*Error) {
	
	_message := FlagStringOrDefault (_flags.Message, "")
	_prompt := FlagStringOrDefault (_flags.Prompt, ">> ")
	_sensitive := FlagBoolOrDefault (_flags.Sensitive, false)
	
	if IsStdoutTerminal () {
		return Errorf (0xbddf576d, "stdout is a TTY")
	}
	if ! IsStderrTerminal () {
		return Errorf (0xf33f2d91, "stderr is not a TTY")
	}
	
	// FIXME:  Make `liner` work without `stdin` or `stdout`
	_stdout := os.Stdout
	os.Stdin = os.Stderr
	os.Stdout = os.Stderr
	syscall.Stdin = int (os.Stderr.Fd ())
	syscall.Stdout = int (os.Stderr.Fd ())
	syscall.Stderr = int (os.Stderr.Fd ())
	
	
	if _message != "" {
		fmt.Fprintln (os.Stderr, _message)
	}
	
	
	_output := ""
	{
		var _output_0 string
		var _error error
		
		_liner := liner.NewLiner ()
		_liner.SetCtrlCAborts (true)
		
		if _sensitive {
			_output_0, _error = _liner.PasswordPrompt (_prompt)
		} else {
			_output_0, _error = _liner.Prompt (_prompt)
		}
		_liner.Close ()
		
		if _error != nil {
			if _error == io.EOF {
				fmt.Fprintln (os.Stderr)
				return Errorf (0x4f6d6f8d, "canceled")
			} else if _error == liner.ErrPromptAborted {
				return Errorf (0x5e488998, "canceled")
			} else {
				return (Errorw (0xa6e02efc, _error))
			}
		}
		
		_output = _output_0
	}
	
	
	fmt.Fprintln (_stdout, _output)
	
	panic (ExitMainSucceeded ())
}

