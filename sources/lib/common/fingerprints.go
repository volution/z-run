

package common


import "crypto/sha256"
import "encoding/binary"
import "encoding/hex"
import "hash"
import "io"
import "sort"




type Fingerprinter struct {
	hasher hash.Hash
}




func NewFingerprinter () (Fingerprinter) {
	return Fingerprinter {
			hasher : sha256.New (),
		}
}


func (_fingerprinter Fingerprinter) Build () (string) {
	return hex.EncodeToString (_fingerprinter.hasher.Sum (nil))
}


func (_fingerprinter Fingerprinter) Uint64 (_value uint64) (Fingerprinter) {
	var _bytes [8]byte
	binary.BigEndian.PutUint64 (_bytes[:], _value)
	_fingerprinter.hasher.Write (_bytes[:])
	return _fingerprinter
}
func (_fingerprinter Fingerprinter) Int64 (_value int64) (Fingerprinter) {
	return _fingerprinter.Uint64 (uint64 (_value))
}


func (_fingerprinter Fingerprinter) String (_value string) (Fingerprinter) {
	io.WriteString (_fingerprinter.hasher, _value)
	return _fingerprinter
}
func (_fingerprinter Fingerprinter) StringWithLen (_value string) (Fingerprinter) {
	_fingerprinter.Uint64 (uint64 (len (_value)))
	io.WriteString (_fingerprinter.hasher, _value)
	return _fingerprinter
}
func (_fingerprinter Fingerprinter) StringsWithLen (_values []string) (Fingerprinter) {
	for _, _value := range _values {
		_fingerprinter.StringWithLen (_value)
	}
	return _fingerprinter
}


func (_fingerprinter Fingerprinter) Bytes (_value []byte) (Fingerprinter) {
	_fingerprinter.hasher.Write (_value)
	return _fingerprinter
}
func (_fingerprinter Fingerprinter) BytesWithLen (_value []byte) (Fingerprinter) {
	_fingerprinter.Uint64 (uint64 (len (_value)))
	_fingerprinter.hasher.Write (_value)
	return _fingerprinter
}


func (_fingerprinter Fingerprinter) StringsMap (_map map[string]string) (Fingerprinter) {
	_keys := make ([]string, 0, len (_map))
	for _key, _ := range _map {
		_keys = append (_keys, _key)
	}
	sort.Strings (_keys)
	for _, _key := range _keys {
		_value := _map[_key]
		_fingerprinter.StringWithLen (_key)
		_fingerprinter.StringWithLen (_value)
	}
	return _fingerprinter
}


func (_fingerprinter Fingerprinter) Bool (_value bool) (Fingerprinter) {
	var _bytes [1]byte
	if _value {
		_bytes[0] = 1
	} else {
		_bytes[0] = 0
	}
	_fingerprinter.hasher.Write (_bytes[:])
	return _fingerprinter
}


func FingerprintStringsQuick (_seed string, _values ... string) string {
	return NewFingerprinter () .String (_seed) .StringsWithLen (_values) .Build ()
}


func FingerprintString (_value string) string {
	return FingerprintBytes ([]byte (_value))
}

func FingerprintBytes (_value []byte) string {
	_hasher := sha256.New ()
	_hasher.Write (_value)
	return hex.EncodeToString (_hasher.Sum (nil))
}

