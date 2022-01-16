
// +build !linux


package common


import "fmt"
import "os"
import "path"
import "runtime"
import "syscall"




// FIXME:  Merge with Linux variant!

func CreatePipe (_size int, _cacheRoot string) (int, *os.File, *Error) {
	
	var _interpreterScriptInput int
	var _interpreterScriptOutput *os.File
	var _interpreterScriptDescriptors [2]int
	
	_maxPipeSize := 0
	switch runtime.GOOS {
		case "darwin" :
			_maxPipeSize = 16 * 1024
		case "freebsd" :
			_maxPipeSize = 512
		case "openbsd" :
			_maxPipeSize = 16 * 1024
		default :
			_maxPipeSize = 0
	}
	
	if _size <= _maxPipeSize {
		if _error := syscall.Pipe (_interpreterScriptDescriptors[:]); _error == nil {
			_interpreterScriptInput = _interpreterScriptDescriptors[0]
			_interpreterScriptOutput = os.NewFile (uintptr (_interpreterScriptDescriptors[1]), "")
		} else {
			return -1, nil, Errorw (0xdf86c693, _error)
		}
	} else {
		if _cacheRoot == "" {
			// FIXME:  We should make sure that the cache path is never empty!
			panic (0xf273c23b)
		}
		if _error := MakeCacheFolder (_cacheRoot, "buffers"); _error != nil {
			return -1, nil, _error
		}
		_temporaryPath := path.Join (_cacheRoot, "buffers", GenerateRandomToken () + ".buffer")
		if _descriptor, _error := syscall.Open (_temporaryPath, syscall.O_CREAT | syscall.O_EXCL | syscall.O_WRONLY, 0600); _error == nil {
			_interpreterScriptOutput = os.NewFile (uintptr (_descriptor), "")
		} else {
			return -1, nil, Errorw (0x1a0f0ad1, _error)
		}
		if _descriptor, _error := syscall.Open (_temporaryPath, syscall.O_RDONLY, 0600); _error == nil {
			_interpreterScriptInput = _descriptor
		} else {
			// FIXME:  Here we leak the first descriptor!
			return -1, nil, Errorw (0x2e13a94f, _error)
		}
		if _error := syscall.Unlink (_temporaryPath); _error != nil {
			// FIXME:  Here we leak both descriptors!
			return -1, nil, Errorw (0x97df431d, _error)
		}
	}
	
	if _, _error := os.Stat (fmt.Sprintf ("/dev/fd/%d", _interpreterScriptInput)); _error != nil {
		// FIXME:  Here we leak both descriptors!
		return -1, nil, Errorw (0x5ea72831, _error)
	}
	
	return _interpreterScriptInput, _interpreterScriptOutput, nil
}


