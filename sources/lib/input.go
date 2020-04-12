

package zrun


import "flag"
import "fmt"
import "io"
import "os"
import "syscall"


import "github.com/peterh/liner"
import isatty "github.com/mattn/go-isatty"




func inputMain (_arguments []string, _environment map[string]string) (*Error) {
	
	
	_message := ""
	_prompt := ">> "
	_sensitive := false
	
	
	_flags := flag.NewFlagSet ("[z-run:input]", flag.ContinueOnError)
	_flags.StringVar (&_message, "message", _message, "message (line or lines) to display before the prompt")
	_flags.StringVar (&_prompt, "prompt", _prompt, "prompt to display before the input")
	_flags.BoolVar (&_sensitive, "sensitive", _sensitive, "treat input as sensitive data and hide the echo")
	
	
	if _error := _flags.Parse (_arguments); _error != nil {
		if _error == flag.ErrHelp {
			os.Exit (0)
		} else {
			return errorw (0xfe27c070, _error)
		}
	}
	if _flags.NArg () > 0 {
		return errorf (0xdc26a939, "unexpected arguments")
	}
	
	
	if isatty.IsTerminal (os.Stdout.Fd ()) {
		return errorf (0xbddf576d, "stdout is a TTY")
	}
	if ! isatty.IsTerminal (os.Stderr.Fd ()) {
		return errorf (0xf33f2d91, "stderr is not a TTY")
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
				return errorf (0x4f6d6f8d, "canceled")
			} else if _error == liner.ErrPromptAborted {
				return errorf (0x5e488998, "canceled")
			} else {
				return (errorw (0xa6e02efc, _error))
			}
		}
		
		_output = _output_0
	}
	
	
	fmt.Fprintln (_stdout, _output)
	
	os.Exit (0)
	panic (0x4fd8aaa0)
}

