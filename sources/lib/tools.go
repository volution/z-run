

package zrun


import crand "crypto/rand"
import "encoding/hex"
import "strings"
import "path"
import "path/filepath"
import "os"




func generateRandomToken () (string) {
	var _data [128 / 8]byte
	if _read, _error := crand.Read (_data[:]); _error == nil {
		if _read != (128 / 8) {
			panic (0xe5f33271)
		}
	} else {
		panic (0x417cda5f)
	}
	_token := hex.EncodeToString (_data[:])
	return _token
}




func resolveExecutable (_executable string, _paths []string) (string, *Error) {
	
	// NOTE:  Based on `os.exec.LookPath` implementation:
	//        https://golang.org/src/os/exec/lp_unix.go
	
	_executable = path.Clean (_executable)
	
	if strings.ContainsRune (_executable, os.PathSeparator) {
		if _executable, _error := filepath.Abs (_executable); _error == nil {
			return _executable, nil
		} else {
			return "", errorw (0x352dbfef, _error)
		}
	}
	
	for _, _folder := range _paths {
		if _folder == "" {
			continue
		}
		if _stat, _error := os.Stat (_folder); _error == nil {
			if ! _stat.Mode () .IsDir () {
				continue
			}
		} else {
			continue
		}
		_file := path.Join (_folder, _executable)
		if _stat, _error := os.Stat (_file); _error == nil {
			if ! _stat.Mode () .IsRegular () {
				continue
			}
			if (_stat.Mode () & 0111) == 0 {
				continue
			}
		} else {
			continue
		}
		if _file, _error := filepath.Abs (_file); _error == nil {
			return _file, nil
		} else {
			return "", errorw (0xac1cdf5d, _error)
		}
	}
	
	return "", errorf (0x9db5ca84, "unresolved executable `%s`", _executable)
}

