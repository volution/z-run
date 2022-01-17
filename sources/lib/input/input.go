

package input


import "fmt"
import "io"
import "os"
import "math/rand"
import "syscall"
import "strings"
import "time"

import "github.com/peterh/liner"

import . "github.com/cipriancraciun/z-run/lib/mainlib"
import . "github.com/cipriancraciun/z-run/lib/common"




type InputMainFlags struct {
	
	Message *string `long:"message" short:"m" value-name:"{message}" description:"message to be displayed on the first line, before the prompt line;"`
	Prompt *string `long:"prompt" short:"p" value-name:"{prompt}" description:"message to de displayed on the prompt line, before the input contents; \n if spaces are desired between the message and the input, then include them in the message; \n if confirm or repeat modes are enabled then '&{EXPECTED}' is replaced by the expected input; \n if retry mode is enabled then '&{RETRY}' is replaced by the retry index;"`
	PromptRepeat *string `long:"prompt-repeat" short:"P" value-name:"{prompt}" description:"(see the previous option;)"`
	
	Default *string `long:"default" short:"d" value-name:"{default}" description:"contents to be used as the default input; \n (not allowed with sensitive, repeat or confirm modes;)"`
	Options *[]string `long:"option" short:"o" value-name:"{option}" description:"contents to be used for auto-completing the input; \n if the default input is desired for auto-completion, then include it in the options; \n (not allowed with sensitive or confirm modes;)"`
	OptionsFiles *[]string `long:"options-file" value-name:"{path}" description:"contents to be used for auto-completion the input, read from the given file, each separated by newline; (multiple allowed;)"`
	OptionsFilesZero *[]string `long:"options-file-zero" value-name:"{path}" description:"contents to be used for auto-completion the input, read from the given file, each separated by '\\0'; (multiple allowed;)"`
	
	Sensitive *bool `long:"sensitive" short:"s" description:"enables a mode that hides the input; \n useful for entering passwords and other sensitive information;"`
	Repeat *bool `long:"repeat" short:"r" description:"enables asking the user to renter the input; \n (not allowed with default or confirm modes;)"`
	
	Retries *uint16 `long:"retry" short:"R" value-name:"{retries}" description:"enables retrying the input, in case of not-empty, repeat or confirm modes;"`
	
	Trim *bool `long:"trim" short:"t" description:"enables triming prefix and suffix spaces from the input; \n useful for handling copy-pasted information;"`
	NotEmpty *bool `long:"not-empty" short:"n" description:"enables checking if the input is not empty;"`
	
	Confirm *bool `long:"confirm" short:"c" description:"enables a mode that displays a token (random or given), and asks the user to re-enter it correctly;"`
	ConfirmToken *string `long:"confirm-token" short:"C" value-name:"{confirm}" description:"contents to be used as the confirm token; \n (will automatically enable confirm mode;  the contents will be automatically trimmed;)"`
	
	OutputSeparator *string `long:"output-separator" value-name:"{separator}" description:"overries the newline separator, written after each input contents;"`
	OutputSeparatorNone *bool `long:"output-separator-none" short:"N" description:"overrides the newline separator with nothing, written after each input contents;"`
	OutputSeparatorZero *bool `long:"output-separator-null" short:"Z" description:"overrides the newline separator with '\\0', written after each input contents;"`
	
	OutputFd *uint16 `long:"output-fd" value-name:"{fd}" description:"overrides input contents writing to the given file-descriptor;"`
	TtyInputFd *uint16 `long:"tty-input-fd" value-name:"{fd}" description:"overrides terminal input from the given file-descriptor;"`
	TtyOutputFd *uint16 `long:"tty-output-fd" value-name:"{fd}" description:"overrides terminal output to the given file-descriptor;"`
	TtyIgnoreChecks *bool `long:"tty-ignore-checks" short:"T" description:"disable checking for a TTY on stderr, and a non-TTY on stdout;"`
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
	_promptRepeat := FlagStringOrDefault (_flags.PromptRepeat, "")
	_default := FlagStringOrDefault (_flags.Default, "")
	_optionsValues := FlagStringsOrDefault (_flags.Options, nil)
	_optionsFiles := FlagStringsOrDefault (_flags.OptionsFiles, nil)
	_optionsFilesZero := FlagStringsOrDefault (_flags.OptionsFilesZero, nil)
	_sensitive := FlagBoolOrDefault (_flags.Sensitive, false)
	_repeat := FlagBoolOrDefault (_flags.Repeat, false)
	_retries := FlagUint16OrDefault (_flags.Retries, 0)
	_trim := FlagBoolOrDefault (_flags.Trim, false)
	_notEmpty := FlagBoolOrDefault (_flags.NotEmpty, false)
	_confirm := FlagBoolOrDefault (_flags.Confirm, false)
	_confirmToken := FlagStringOrDefault (_flags.ConfirmToken, "")
	_outputSeparator := FlagStringOrDefault (_flags.OutputSeparator, "\n")
	_outputSeparatorNone := FlagBoolOrDefault (_flags.OutputSeparatorNone, false)
	_outputSeparatorZero := FlagBoolOrDefault (_flags.OutputSeparatorZero, false)
	_outputFd := uintptr (FlagUint16OrDefault (_flags.OutputFd, 1))
	_ttyInputFd := uintptr (FlagUint16OrDefault (_flags.TtyInputFd, 2))
	_ttyOutputFd := uintptr (FlagUint16OrDefault (_flags.TtyOutputFd, 2))
	_ttyIgnoreChecks := FlagBoolOrDefault (_flags.TtyIgnoreChecks, false)
	
	
	if (_flags.Default != nil) && (_sensitive || _repeat || _confirm) {
		return Errorf (0x64a90a9f, "`--default` not allowed with `--sensitive`, `--repeat` or `--confirm`!")
	}
	if ((_flags.Options != nil) || (_flags.OptionsFiles != nil) || (_flags.OptionsFilesZero != nil)) && (_sensitive || _confirm) {
		return Errorf (0xe06e39d2, "`--option`, `--option-file`, and `--option-file-zero` not allowed with `--sensitive` or `--confirm`!")
	}
	
	_options := []string (nil)
	_optionsSeen := make (map[string]bool, 1024)
	_optionsAppend := func (_values []string) {
			for _, _value := range _values {
				if _, _exists := _optionsSeen[_value]; _exists {
					continue
				}
				_options = append (_options, _value)
				_optionsSeen[_value] = true
			}
		}
	if _optionsValues != nil {
		_optionsAppend (_optionsValues)
	}
	for _, _path := range _optionsFiles {
		if _lines, _error := ReadFileLines (_path, '\n', true, true); _error == nil {
			_optionsAppend (_lines)
		} else {
			return _error
		}
	}
	for _, _path := range _optionsFilesZero {
		if _lines, _error := ReadFileLines (_path, '\000', true, true); _error == nil {
			_optionsAppend (_lines)
		} else {
			return _error
		}
	}
	
	if (_flags.ConfirmToken != nil) && !_confirm {
		_confirm = true
	}
	if _confirm {
		_confirmToken = strings.TrimSpace (_confirmToken)
	}
	if _confirm && (_sensitive || _repeat) {
		return Errorf (0xba914320, "`--confirm` not allowed with `--sensitive` or `--repeat`!")
	}
	
	if (_flags.PromptRepeat != nil) && (_flags.Prompt == nil) {
		return Errorf (0xf3b77f85, "`--prompt-repeat` not allowed without `--prompt`!")
	}
	if (_promptRepeat != "") && !_repeat {
		return Errorf (0x9dcc574c, "`--prompt-repeat` not allowed without `--repeat`!")
	}
	if _promptRepeat == "" {
		_promptRepeat = _prompt
	}
	
	if (_flags.OutputSeparator != nil) || (_flags.OutputSeparatorNone != nil) || (_flags.OutputSeparatorZero != nil) {
		if (_flags.OutputSeparator != nil) && (_flags.OutputSeparatorZero != nil) {
			return Errorf (0x86622151, "`--output-separator` and `--output-separator-null` are mutually exclusive!")
		}
		if (_flags.OutputSeparator != nil) && (_flags.OutputSeparatorNone != nil) {
			return Errorf (0x86622151, "`--output-separator` and `--output-separator-none` are mutually exclusive!")
		}
		if (_flags.OutputSeparatorNone != nil) && (_flags.OutputSeparatorZero != nil) {
			return Errorf (0x86622151, "`--output-separator-none` and `--output-separator-zero` are mutually exclusive!")
		}
		if _confirm {
			return Errorf (0xee6fbcd7, "`--output-separator`, `--output-separator-none`, and `--output-separator-null` not allowed with `--confirm`!")
		}
		if _outputSeparatorNone {
			_outputSeparator = ""
		}
		if _outputSeparatorZero {
			_outputSeparator = "\000"
		}
	}
	
	
	
	
	if !_ttyIgnoreChecks {
		if IsFdTerminal (_outputFd) {
			return Errorf (0xbddf576d, "stdout is a TTY")
		}
		if ! IsFdTerminal (_ttyInputFd) {
			return Errorf (0xf33f2d91, "stderr is not a TTY")
		}
		if ! IsFdTerminal (_ttyOutputFd) {
			return Errorf (0xe8c5f8bc, "stderr is not a TTY")
		}
	}
	
	{
		if _fd_0, _error := syscall.Dup (int (_outputFd)); _error == nil {
			_outputFd = uintptr (_fd_0)
		} else {
			return Errorw (0x59a1994e, _error)
		}
		if _fd_0, _error := syscall.Dup (int (_ttyInputFd)); _error == nil {
			_ttyInputFd = uintptr (_fd_0)
		} else {
			return Errorw (0x0ceb87ec, _error)
		}
		if _fd_0, _error := syscall.Dup (int (_ttyOutputFd)); _error == nil {
			_ttyOutputFd = uintptr (_fd_0)
		} else {
			return Errorw (0x8dc54e20, _error)
		}
	}
	
	// FIXME:  Make `liner` work without `stdin` or `stdout`!
	
	{
		if _error := syscall.Dup2 (int (_ttyInputFd), 0); _error != nil {
			return Errorw (0x180f62b3, _error)
		}
		if _error := syscall.Dup2 (int (_ttyOutputFd), 1); _error != nil {
			return Errorw (0xe252bec9, _error)
		}
	}
	
	_outputStream := os.NewFile (uintptr (_outputFd), "/dev/null")
	os.Stdin = os.NewFile (uintptr (_ttyInputFd), "/dev/stdin")
	os.Stdout = os.NewFile (uintptr (_ttyOutputFd), "/dev/stdout")
	
	
	
	
	if _message != "" {
		fmt.Fprintln (os.Stdout, _message)
	}
	
	var _input string
	var _inputEntered bool
	
	_loop := uint16 (0)
	for _loop <= _retries {
		_loop += 1
		
		_confirmEnabled := false
		_confirmTokenOutput := ""
		_confirmTokenExpected := ""
		if _confirm {
			_confirmEnabled = true
			if _confirmToken != "" {
				_confirmTokenExpected = _confirmToken
			} else {
				_confirmTokenExpected = token (0)
			}
			_confirmTokenOutput = _confirmTokenExpected
		} else if _repeat && _inputEntered {
			_confirmEnabled = true
			_confirmTokenExpected = _input
			if _sensitive {
				_confirmTokenOutput = "****"
			} else {
				_confirmTokenOutput = _confirmTokenExpected
			}
		}
		
		_prompt_0 := _prompt
		if _confirmEnabled {
			_prompt_1 := _promptRepeat
			if _prompt_1 == "" {
				_prompt_1 = _prompt
			}
			_prompt_2 := strings.ReplaceAll (_prompt_1, "&{EXPECTED}", _confirmTokenOutput)
			if _prompt_2 != _prompt_1 {
				_prompt_0 = _prompt_2
			} else {
				if _confirm {
					_prompt_0 = fmt.Sprintf ("[confirm `%s`] %s", _confirmTokenOutput, _prompt_1)
				} else {
					_prompt_0 = fmt.Sprintf ("[repeat `%s`] %s", _confirmTokenOutput, _prompt_1)
				}
			}
		}
		if _retries > 0 {
			_retryToken := fmt.Sprintf ("%d", _loop - 1)
			_prompt_1 := strings.ReplaceAll (_prompt_0, "&{RETRY}", _retryToken)
			if _prompt_1 != _prompt_0 {
				_prompt_0 = _prompt_1
			} else {
				if _loop > 1 {
					_prompt_0 = fmt.Sprintf ("[retry %d] %s", _loop - 1, _prompt_0)
				}
			}
		}
		
		_input_0, _canceled, _error := input (_prompt_0, _default, _options, _sensitive, _trim)
		
		if _canceled {
			panic (ExitMainFailed ())
		}
		if _error != nil {
			panic (AbortError (_error))
		}
		if _notEmpty && (_input_0 == "") {
			continue
		}
		
		if _confirmEnabled && (_input_0 != _confirmTokenExpected) {
			continue
		}
		
		if _inputEntered {
			_loop = 0
			break
		} else {
			_input = _input_0
			_inputEntered = true
			_loop = 0
			if !_repeat {
				break
			}
		}
	}
	
	if _loop > 0 {
		panic (ExitMainFailed ())
	}
	
	if _confirm {
		_input = ""
	}
	
	if _input != "" {
		var _buffer strings.Builder
		_buffer.WriteString (_input)
		_buffer.WriteString (_outputSeparator)
		if _, _error := io.WriteString (_outputStream, _buffer.String ()); _error != nil {
			panic (AbortError (Errorw (0x2970f683, _error)))
		}
	}
	
	panic (ExitMainSucceeded ())
}




func input (_prompt string, _default string, _options []string, _sensitive bool, _trim bool) (string, bool, *Error) {
	
	var _input string
	var _error error
	
	_liner := liner.NewLiner ()
	_liner.SetCtrlCAborts (true)
	defer _liner.Close ()
	
	if len (_options) > 0 {
		_completer := func (_prefix string) ([]string) {
				_filtered := make ([]string, 0, len (_options))
				for _, _option := range _options {
					if strings.HasPrefix (_option, _prefix) {
						_filtered = append (_filtered, _option)
					}
				}
				return _filtered
			}
		_liner.SetCompleter (_completer)
		_liner.SetTabCompletionStyle (liner.TabPrints)
	}
	
	if _sensitive {
		_input, _error = _liner.PasswordPrompt (_prompt)
	} else {
		if _default != "" {
			_input, _error = _liner.PromptWithSuggestion (_prompt, _default, -1)
		} else {
			_input, _error = _liner.Prompt (_prompt)
		}
	}
	_liner.Close ()
	
	if _error == nil {
		if _trim {
			_input = strings.TrimSpace (_input)
		}
	} else {
		if _error == io.EOF {
			fmt.Fprintln (os.Stdout)
			return "", true, Errorf (0x4f6d6f8d, "canceled")
		} else if _error == liner.ErrPromptAborted {
			return "", true, Errorf (0x5e488998, "canceled")
		} else {
			return "", false, Errorw (0xa6e02efc, _error)
		}
	}
	
	return _input, false, nil
}




func token (_length uint) (string) {
	
	if _length == 0 {
		_length = 4
	}
	
	// NOTE:  This token doesn't need to be cryptographically secure.
	_source := rand.NewSource (time.Now () .UnixNano ())
	_rand := rand.New (_source)
	
	var _buffer strings.Builder
	
	if _length <= 10 {
		_digits := _rand.Perm (10) [:_length]
		for _, _digit := range _digits {
			_buffer.WriteByte ('0' + byte (_digit))
		}
	} else {
		var _index uint
		for _index = 0; _index < _length; _index += 1 {
			_digit := _rand.Intn (10)
			_buffer.WriteByte ('0' + byte (_digit))
		}
	}
	
	return _buffer.String ()
}

