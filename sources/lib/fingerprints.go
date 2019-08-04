

package lib


import "crypto/sha256"
import "encoding/binary"
import "encoding/hex"
import "hash"
import "io"




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


func (_fingerprinter Fingerprinter) Bytes (_value []byte) (Fingerprinter) {
	_fingerprinter.hasher.Write (_value)
	return _fingerprinter
}
func (_fingerprinter Fingerprinter) BytesWithLen (_value []byte) (Fingerprinter) {
	_fingerprinter.Uint64 (uint64 (len (_value)))
	_fingerprinter.hasher.Write (_value)
	return _fingerprinter
}

