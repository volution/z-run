

package common


import "bufio"
import "io"
import "os"
import "strings"




func ReadFileLines (_path string, _separator byte, _trimSpaces bool, _skipEmpty bool) ([]string, *Error) {
	
	_stream, _error := os.Open (_path)
	if _error != nil {
		return nil, Errorw (0x8279d199, _error)
	}
	defer _stream.Close ()
	
	_lines := make ([]string, 0, 1024)
	
	_buffer := bufio.NewReader (_stream)
	for {
		
		_line, _error := _buffer.ReadString (_separator)
		if _error == io.EOF {
			break
		} else if _error != nil {
			return nil, Errorw (0xc0f8ce25, _error)
		}
		
		_lineLen := len (_line)
		if (_lineLen > 0) && (_line[_lineLen - 1] == _separator) {
			_line = _line[: _lineLen - 1]
			_lineLen -= 1
		}
		
		if _trimSpaces {
			_line = strings.TrimSpace (_line)
			_lineLen = len (_line)
		}
		if (_lineLen == 0) && _skipEmpty {
			continue
		}
		
		_lines = append (_lines, _line)
	}
	
	return _lines, nil
}

