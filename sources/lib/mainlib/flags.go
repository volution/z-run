

package mainlib


import "bytes"
import "io"
import "os"

import "github.com/jessevdk/go-flags"

import . "github.com/volution/z-run/lib/common"




func ResolveMainFlags (_executableName string, _arguments []string, _environment map[string]string, _flags interface{}, _helpStream io.Writer) (*Error) {
	
	// FIXME:  Because the parser always uses the actual environment variables and not `_environment` we need to force it...
	_overrideEnv := func (_env map[string]string) () {
		os.Clearenv ()
		for _name, _value := range _env {
			os.Setenv (_name, _value)
		}
	}
	_environmentOriginal := EnvironmentListToMap (os.Environ ())
	_overrideEnv (_environment)
	defer _overrideEnv (_environmentOriginal)
	
	var _parserFlags flags.Options = flags.PassDoubleDash
	if _helpStream != nil {
		_parserFlags |= flags.HelpFlag
	}
	
	_parser := flags.NewNamedParser (_executableName, _parserFlags)
	_parser.SubcommandsOptional = true
	
	if _, _error := _parser.AddGroup ("Main Options", "", _flags); _error != nil {
		return Errorw (0x39fd8298, _error)
	}
	
	if _restArguments, _error := _parser.ParseArgs (_arguments); _error != nil {
		_parserError, _ := _error.(*flags.Error)
		if (_parserError != nil) && (_parserError.Type == flags.ErrHelp) && (_helpStream != nil) {
			_buffer := bytes.NewBuffer (nil)
			_buffer.WriteByte ('\n')
			_parser.WriteHelp (_buffer)
			_buffer.WriteByte ('\n')
			if _, _error := _helpStream.Write (_buffer.Bytes ()); _error == nil {
				return ResolveMainFlagsHelpError
			} else {
				return Errorw (0xe312aa70, _error)
			}
		} else {
			return Errorw (0x4b242ec4, _error)
		}
	} else if len (_restArguments) != 0 {
		return Errorf (0x32c22da5, "invalid arguments!")
	}
	
	return nil
}


var ResolveMainFlagsHelpError = Errorf (0xf9a36b30, "help invoked;  aborting!")




func FlagBoolOrDefault (_value *bool, _default bool) (bool) {
	if _value != nil {
		return *_value
	}
	return _default
}

func FlagUint16OrDefault (_value *uint16, _default uint16) (uint16) {
	if _value != nil {
		return *_value
	}
	return _default
}

func FlagStringOrDefault (_value *string, _default string) (string) {
	if _value != nil {
		return *_value
	}
	return _default
}

func FlagStringsOrDefault (_value *[]string, _default []string) ([]string) {
	if _value != nil {
		return *_value
	}
	return _default
}




func Flag2BoolOrDefault (_value_1 *bool, _value_2 *bool, _default bool) (bool) {
	if _value_1 != nil {
		return *_value_1
	}
	if _value_2 != nil {
		return *_value_2
	}
	return _default
}

func Flag2Uint16OrDefault (_value_1 *uint16, _value_2 *uint16, _default uint16) (uint16) {
	if _value_1 != nil {
		return *_value_1
	}
	if _value_2 != nil {
		return *_value_2
	}
	return _default
}

func Flag2StringOrDefault (_value_1 *string, _value_2 *string, _default string) (string) {
	if _value_1 != nil {
		return *_value_1
	}
	if _value_2 != nil {
		return *_value_2
	}
	return _default
}

func Flag2StringsOrDefault (_value_1 *[]string, _value_2 *[]string, _default []string) ([]string) {
	if _value_1 != nil {
		return *_value_1
	}
	if _value_2 != nil {
		return *_value_2
	}
	return _default
}

