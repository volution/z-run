

package common


import crand "crypto/rand"
import "encoding/hex"




func GenerateRandomToken () (string) {
	var _data [128 / 8]byte
	if _read, _error := crand.Read (_data[:]); _error == nil {
		if _read != (128 / 8) {
			panic (AbortUnreachable (0xe5f33271))
		}
	} else {
		panic (AbortUnreachable (0x417cda5f))
	}
	_token := hex.EncodeToString (_data[:])
	return _token
}

