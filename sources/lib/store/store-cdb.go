

package store


import "bytes"
import "encoding/json"
import "fmt"
import "os"

import cdb "github.com/colinmarc/cdb"

import . "github.com/volution/z-run/lib/common"




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
			return nil, Errorw (0xb4697ac2, _error)
		}
	} else {
		return nil, Errorw (0x9566422e, _error)
	}
}


func (_store *CdbStoreInput) SelectObject (_instance string, _global bool, _namespace string, _key string, _value interface{}) (bool, *Error) {
	
	var _valueData []byte
	if _found, _data, _error := _store.SelectRawBytes (_instance, _global, _namespace, _key); _error == nil {
		if _found {
			_valueData = _data
		} else {
			return false, nil
		}
	} else {
		return false, _error
	}
	
	switch _value := _value.(type) {
		case []byte :
			return false, Errorf (0xed6ab84e, "unexpected type")
		case *[]byte :
			if *_value == nil {
				*_value = _valueData
			} else {
				*_value = append ((*_value)[:0], _valueData ...)
			}
		case string :
			return false, Errorf (0x36de066a, "unexpected type")
		case *string :
			*_value = string (_valueData)
		default :
			if _error := json.Unmarshal (_valueData, _value); _error != nil {
				return false, Errorw (0x8ebfc830, _error)
			}
	}
	
	return true, nil
}


func (_store *CdbStoreInput) SelectRawString (_instance string, _global bool, _namespace string, _key string) (bool, string, *Error) {
	
	var _valueData []byte
	if _found, _data, _error := _store.SelectRawBytes (_instance, _global, _namespace, _key); _error == nil {
		if _found {
			_valueData = _data
		} else {
			return false, "", nil
		}
	} else {
		return false, "", _error
	}
	
	_value := string (_valueData)
	
	return true, _value, nil
}


func (_store *CdbStoreInput) SelectRawBytes (_instance string, _global bool, _namespace string, _key string) (bool, []byte, *Error) {
	
	_keyBuffer := bytes.NewBuffer (nil)
	_keyBuffer.WriteString (_namespace)
	_keyBuffer.WriteString (" // ")
	_keyBuffer.WriteString (_key)
	
	var _valueData []byte
	if _valueData_0, _error := _store.reader.Get (_keyBuffer.Bytes ()); _error == nil {
		_valueData = _valueData_0
	} else {
		return false, nil, Errorw (0x811cd704, _error)
	}
	
	if _valueData == nil {
		return false, nil, nil
	}
	
	return true, _valueData, nil
}


func (_store *CdbStoreInput) Close () (*Error) {
	if _error := _store.reader.Close (); _error == nil {
		_store.reader = nil
		return Errorw (0x10dd4299, _error)
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
		return nil, Errorf (0x6917ab7d, "invalid path")
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
		return nil, Errorw (0x064e6aec, _error)
	}
}


func (_store *CdbStoreOutput) IncludeObject (_instance string, _global bool, _namespace string, _key string, _value interface{}) (*Error) {
	
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
				return Errorw (0x02600262, _error)
			}
	}
	
	return _store.IncludeRawBytes (_instance, _global, _namespace, _key, _valueBuffer.Bytes ())
}


func (_store *CdbStoreOutput) IncludeRawString (_instance string, _global bool, _namespace string, _key string, _value string) (*Error) {
	
	_valueData := []byte (_value)
	
	return _store.IncludeRawBytes (_instance, _global, _namespace, _key, _valueData)
}


func (_store *CdbStoreOutput) IncludeRawBytes (_instance string, _global bool, _namespace string, _key string, _value []byte) (*Error) {
	
	_keyBuffer := bytes.NewBuffer (nil)
	_keyBuffer.WriteString (_namespace)
	_keyBuffer.WriteString (" // ")
	_keyBuffer.WriteString (_key)
	
	if _error := _store.writer.Put (_keyBuffer.Bytes (), _value); _error == nil {
		return nil
	} else {
		return Errorw (0x28b9a333, _error)
	}
}


func (_store *CdbStoreOutput) Commit () (*Error) {
	if _error := _store.writer.Close (); _error != nil {
		return Errorw (0xfc3f2b2b, _error)
	}
	_store.writer = nil
	if _error := os.Rename (_store.pathTemporary, _store.pathFinal); _error != nil {
		return Errorw (0x62423da2, _error)
	}
	return nil
}

