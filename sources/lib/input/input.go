

package input


import "fmt"
import "io"
import "os"
import "syscall"
import "strings"

import "github.com/peterh/liner"

import . "github.com/cipriancraciun/z-run/lib/mainlib"
import . "github.com/cipriancraciun/z-run/lib/common"




type InputMainFlags struct {
	
	Message *string `long:"message" short:"m" value-name:"{message}" description:"message to be displayed before the prompt line;"`
	Prompt *string `long:"prompt" short:"p" value-name:"{prompt}" description:"message to de displayed on the prompt line, before the input; \n if spaces are desired between the message and the input, then include them in the message itself;"`
	Default *string `long:"default" short:"d" value-name:"{default}" description:"contents to be used as the default input;  (not allowed with sensitive or confirm;)"`
	
	Sensitive *bool `long:"sensitive" short:"s" description:"enables hiding the input;  useful for entering passwords and other sensitive information;"`
	
	Trim *bool `long:"trim" short:"t" description:"enables triming prefix and suffix spaces from the input;  useful for handling copy-pasted information;"`
	NotEmpty *bool `long:"not-empty" short:"n" description:"enables checking if the input is not empty, else the tool exits with an error;"`
	
	Retries *uint16 `long:"retry" short:"r" value-name:"{retries}" description:"enables retrying the input, in case of not-empty or confirm modes;"`
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
	_default := FlagStringOrDefault (_flags.Default, "")
	_sensitive := FlagBoolOrDefault (_flags.Sensitive, false)
	_trim := FlagBoolOrDefault (_flags.Trim, false)
	_notEmpty := FlagBoolOrDefault (_flags.NotEmpty, false)
	_retries := FlagUint16OrDefault (_flags.Retries, 0)
	
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
	
	var _output string
	
	_loop := uint16 (0)
	for _loop <= _retries {
		_loop += 1
		_prompt_0 := _prompt
		if _loop > 1 {
			_prompt_0 = fmt.Sprintf ("[%d] %s", _loop - 1, _prompt)
		}
		_output_0, _canceled, _error := input (_prompt_0, _default, _sensitive, _trim)
		if _canceled {
			panic (ExitMainFailed ())
		}
		if _error != nil {
			panic (AbortError (_error))
		}
		if _notEmpty && (_output_0 == "") {
			continue
		}
		_output = _output_0
		_loop = 0
		break
	}
	if _loop > 0 {
		panic (ExitMainFailed ())
	}
	
	fmt.Fprintln (_stdout, _output)
	
	panic (ExitMainSucceeded ())
}




func input (_prompt string, _default string, _sensitive bool, _trim bool) (string, bool, *Error) {
	
	var _output string
	var _error error
	
	_liner := liner.NewLiner ()
	_liner.SetCtrlCAborts (true)
	
	if _sensitive {
		_output, _error = _liner.PasswordPrompt (_prompt)
	} else {
		if _default != "" {
			_output, _error = _liner.PromptWithSuggestion (_prompt, _default, -1)
		} else {
			_output, _error = _liner.Prompt (_prompt)
		}
	}
	_liner.Close ()
	
	if _error == nil {
		if _trim {
			_output = strings.TrimSpace (_output)
		}
	} else {
		if _error == io.EOF {
			fmt.Fprintln (os.Stderr)
			return "", true, Errorf (0x4f6d6f8d, "canceled")
		} else if _error == liner.ErrPromptAborted {
			return "", true, Errorf (0x5e488998, "canceled")
		} else {
			return "", false, Errorw (0xa6e02efc, _error)
		}
	}
	
	return _output, false, nil
}

