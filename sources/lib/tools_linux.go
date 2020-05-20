

// +build linux


package zrun


import "os"
import "path"
import "syscall"

import "golang.org/x/sys/unix"




// FIXME:  Merge with non-Linux variant!

func createPipe (_size int, _cacheRoot string) (int, *os.File, *Error) {
	
	var _interpreterScriptInput int
	var _interpreterScriptOutput *os.File
	var _interpreterScriptDescriptors [2]int
	
	if _size < (1 * 1024 * 1024 - 128 * 1024) {
		if _error := syscall.Pipe (_interpreterScriptDescriptors[:]); _error == nil {
			if _, _error := unix.FcntlInt (uintptr (_interpreterScriptDescriptors[1]), unix.F_SETPIPE_SZ, 1 * 1024 * 1024); _error != nil {
				logf ('w', 0x4d3414c8, "failed increasing pipe buffer size;  ignoring!")
			}
			_interpreterScriptInput = _interpreterScriptDescriptors[0]
			_interpreterScriptOutput = os.NewFile (uintptr (_interpreterScriptDescriptors[1]), "")
		} else {
			return -1, nil, errorw (0xece645ff, _error)
		}
	} else {
		if _cacheRoot == "" {
			// FIXME:  We should make sure that the cache path is never empty!
			panic (0xd6f17610)
		}
		_temporaryPath := path.Join (_cacheRoot, generateRandomToken () + ".scriptlet")
		if _descriptor, _error := syscall.Open (_temporaryPath, syscall.O_CREAT | syscall.O_EXCL | syscall.O_WRONLY, 0600); _error == nil {
			_interpreterScriptOutput = os.NewFile (uintptr (_descriptor), "")
		} else {
			return -1, nil, errorw (0x2b19feaa, _error)
		}
		if _descriptor, _error := syscall.Open (_temporaryPath, syscall.O_RDONLY, 0600); _error == nil {
			_interpreterScriptInput = _descriptor
		} else {
			// FIXME:  Here we leak the first descriptor!
			return -1, nil, errorw (0x694ce572, _error)
		}
		if _error := syscall.Unlink (_temporaryPath); _error != nil {
			// FIXME:  Here we leak both descriptors!
			return -1, nil, errorw (0xc5afd6fd, _error)
		}
	}
	
	return _interpreterScriptInput, _interpreterScriptOutput, nil
}

