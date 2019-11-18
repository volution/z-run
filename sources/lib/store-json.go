

package zrun


import "encoding/json"
import "io"




type JsonStreamStoreOutput struct {
	stream io.Writer
	closer io.Closer
	encoder *json.Encoder
}


type JsonStoreRecord struct {
	Namespace string `json:"namespace"`
	Key string `json:"key"`
	Value interface{} `json:"value"`
}




func NewJsonStreamStoreOutput (_stream io.Writer, _closer io.Closer) (*JsonStreamStoreOutput) {
	_encoder := json.NewEncoder (_stream)
	_encoder.SetIndent ("", "    ")
	_encoder.SetEscapeHTML (false)
	return & JsonStreamStoreOutput {
			stream : _stream,
			closer : _closer,
			encoder : _encoder,
		}
}

func (_store *JsonStreamStoreOutput) Include (_namespace string, _key string, _value interface{}) (*Error) {
	_record := & JsonStoreRecord {
			Namespace : _namespace,
			Key : _key,
			Value : _value,
		}
	if _error := _store.encoder.Encode (_record); _error == nil {
		return nil
	} else {
		return errorw (0x5435f95a, _error)
	}
}

func (_store *JsonStreamStoreOutput) Commit () (*Error) {
	var _error *Error
	if _store.closer != nil {
		if _error_0 := _store.closer.Close (); _error_0 != nil {
			_error = errorw (0x9f5565fc, _error_0)
		}
	}
	_store.stream = nil
	_store.closer = nil
	_store.encoder = nil
	return _error
}

