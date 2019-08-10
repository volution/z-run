

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

func (_store *JsonStreamStoreOutput) Include (_namespace string, _key string, _value interface{}) (error) {
	_record := & JsonStoreRecord {
			Namespace : _namespace,
			Key : _key,
			Value : _value,
		}
	return _store.encoder.Encode (_record)
}

func (_store *JsonStreamStoreOutput) Commit () (error) {
	var _error error
	if _store.closer != nil {
		_error = _store.closer.Close ()
	}
	_store.stream = nil
	_store.closer = nil
	_store.encoder = nil
	return _error
}

