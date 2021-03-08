

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




func resolveRelativePath (_workspace string, _base string, _path string) (string, *Error) {
	
	if (_path == ".") || (_path == "..") || (_path == "_") {
		_path = _path + "/"
	}
	
	if _path == "" {
		// NOP
	} else if path.IsAbs (_path) {
		// NOP
	} else if strings.HasPrefix (_path, "." + string (os.PathSeparator)) {
		_path = path.Join (_workspace, _path)
	} else if strings.HasPrefix (_path, ".." + string (os.PathSeparator)) {
		_path = path.Join (_workspace, _path)
	} else if strings.HasPrefix (_path, "_" + string (os.PathSeparator)) {
		_path = path.Join (_base, _path[2:])
	} else {
		return "", errorf (0x3ca0a241, "invalid path syntax: `%s`", _path)
	}
	
	_path = path.Clean (_path)
	
	if _path == "" {
		return "", errorf (0xe971645a, "invalid empty path")
	}
	
	return _path, nil
}


func resolveAbsolutePath (_workspace string, _base string, _path string) (string, *Error) {
	
	if _path_0, _error := resolveRelativePath (_workspace, _base, _path); _error != nil {
		return "", _error
	} else {
		_path = _path_0
	}
	if _path_0, _error := filepath.Abs (_path); _error == nil {
		_path = _path_0
	} else {
		return "", errorw (0xb007b166, _error)
	}
	
	return _path, nil
}


func replaceVariables (_input string) (string, *Error) {
	
	// FIXME:  Implement this better!
	
	_input = strings.ReplaceAll (_input, "${ZRUN_OS}", BUILD_TARGET_OS)
	_input = strings.ReplaceAll (_input, "${ZRUN_ARCH}", BUILD_TARGET_ARCH)
	_input = strings.ReplaceAll (_input, "${ZRUN_VERSION}", BUILD_VERSION)
	
	_input = strings.ReplaceAll (_input, "${UNAME_NODE}", UNAME_NODE)
	_input = strings.ReplaceAll (_input, "${UNAME_SYSTEM}", UNAME_SYSTEM)
	_input = strings.ReplaceAll (_input, "${UNAME_RELEASE}", UNAME_RELEASE)
	_input = strings.ReplaceAll (_input, "${UNAME_VERSION}", UNAME_VERSION)
	_input = strings.ReplaceAll (_input, "${UNAME_MACHINE}", UNAME_MACHINE)
	_input = strings.ReplaceAll (_input, "${UNAME_FINGERPRINT}", UNAME_FINGERPRINT)
	
	return _input, nil
}

