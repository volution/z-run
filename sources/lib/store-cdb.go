

package zrun


import "bytes"
import "encoding/json"
import "fmt"
import "os"


import cdb "github.com/colinmarc/cdb"




type CdbStoreInput struct {
	reader *cdb.CDB
	path string
}


func NewCdbStoreInput (_path string) (*CdbStoreInput, *Error) {
	if _file, _error := os.Open (_path); _error == nil {
		defer _file.Close ()
		if _reader, _error := cdb.NewFromMappedWithHasher (_file, nil); _error == nil {
			_store := & CdbStoreInput {
					reader : _reader,
					path : _path,
				}
			return _store, nil
		} else {
			return nil, errorw (0xb4697ac2, _error)
		}
	} else {
		return nil, errorw (0x9566422e, _error)
	}
}


func (_store *CdbStoreInput) Select (_namespace string, _key string, _value interface{}) (bool, *Error) {
	
	_keyBuffer := bytes.NewBuffer (nil)
	_keyBuffer.WriteString (_namespace)
	_keyBuffer.WriteString (" // ")
	_keyBuffer.WriteString (_key)
	
	var _valueData []byte
	if _valueData_0, _error := _store.reader.Get (_keyBuffer.Bytes ()); _error == nil {
		_valueData = _valueData_0
	} else {
		return false, errorw (0x811cd704, _error)
	}
	
	if _valueData == nil {
		return false, nil
	}
	
	switch _value := _value.(type) {
		case []byte :
			return false, errorf (0xed6ab84e, "unexpected type")
		case *[]byte :
			if *_value == nil {
				*_value = _valueData
			} else {
				*_value = append ((*_value)[:0], _valueData ...)
			}
		case string :
			return false, errorf (0x36de066a, "unexpected type")
		case *string :
			*_value = string (_valueData)
		default :
			if _error := json.Unmarshal (_valueData, _value); _error != nil {
				return false, errorw (0x8ebfc830, _error)
			}
	}
	
	return true, nil
}


func (_store *CdbStoreInput) Close () (*Error) {
	if _error := _store.reader.Close (); _error == nil {
		_store.reader = nil
		return errorw (0x10dd4299, _error)
	} else {
		return nil
	}
}




type CdbStoreOutput struct {
	writer *cdb.Writer
	pathFinal string
	pathTemporary string
}


func NewCdbStoreOutput (_path string) (*CdbStoreOutput, *Error) {
	if _path == "" {
		return nil, errorf (0x6917ab7d, "invalid path")
	}
	_pathFinal := _path
	_pathTemporary := fmt.Sprintf ("%s--%08x.tmp", _pathFinal, os.Getpid ())
	if _writer, _error := cdb.Create (_pathTemporary); _error == nil {
		_store := & CdbStoreOutput {
				writer : _writer,
				pathFinal : _pathFinal,
				pathTemporary : _pathTemporary,
			}
		return _store, nil
	} else {
		return nil, errorw (0x064e6aec, _error)
	}
}


func (_store *CdbStoreOutput) Include (_namespace string, _key string, _value interface{}) (*Error) {
	
	_keyBuffer := bytes.NewBuffer (nil)
	_keyBuffer.WriteString (_namespace)
	_keyBuffer.WriteString (" // ")
	_keyBuffer.WriteString (_key)
	
	_valueBuffer := bytes.NewBuffer (nil)
	switch _value := _value.(type) {
		case []byte :
			_valueBuffer.Write (_value)
		case *[]byte :
			_valueBuffer.Write (*_value)
		case string :
			_valueBuffer.WriteString (_value)
		case *string :
			_valueBuffer.WriteString (*_value)
		default :
			if _error := json.NewEncoder (_valueBuffer) .Encode (_value); _error != nil {
				return errorw (0x02600262, _error)
			}
	}
	
	if _error := _store.writer.Put (_keyBuffer.Bytes (), _valueBuffer.Bytes ()); _error == nil {
		return nil
	} else {
		return errorw (0x28b9a333, _error)
	}
}

func (_store *CdbStoreOutput) Commit () (*Error) {
	if _error := _store.writer.Close (); _error != nil {
		return errorw (0xfc3f2b2b, _error)
	}
	_store.writer = nil
	if _error := os.Rename (_store.pathTemporary, _store.pathFinal); _error != nil {
		return errorw (0x62423da2, _error)
	}
	return nil
}

