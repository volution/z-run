

package lib


import "os"
import "path"
import "path/filepath"
import "sort"




func resolveSources (_candidate string) ([]*Source, error) {
	
	_sources := make ([]*Source, 0, 128)
	
	_path, _stat, _error := resolveSourcePath_0 (_candidate)
	if _error != nil {
		return nil, _error
	}
	
	_statMode := _stat.Mode ()
	switch {
		
		case _statMode.IsRegular () :
			if _source, _error := resolveSource_0 (_path, _stat); _error == nil {
				_sources = append (_sources, _source)
			} else {
				return nil, _error
			}
		
		case _statMode.IsDir () :
			return nil, errorf (0x8a04b23b, "not-implemented")
		
		default :
			return nil, errorf (0xa35428a2, "invalid source `%s`", _path)
	}
	
	return _sources, nil
}


func resolveSource (_candidate string) (*Source, error) {
	if _path, _stat, _error := resolveSourcePath_0 (_candidate); _error == nil {
		return resolveSource_0 (_path, _stat)
	} else {
		return nil, _error
	}
}




func resolveSource_0 (_path string, _stat os.FileInfo) (*Source, error) {
	_statMode := _stat.Mode ()
	if _statMode.IsRegular () {
		_source := & Source {
				Path : _path,
				Executable : (_statMode.Perm () & 0111) != 0,
				FingerprintMeta : fingerprintSource_1 (_path, _stat),
			}
		return _source, nil
	} else {
		return nil, errorf (0x557961c4, "invalid source `%s`", _path)
	}
}


func fingerprintSource_0 (_path string) (string, error) {
	if _stat, _error := os.Stat (_path); _error == nil {
		return fingerprintSource_1 (_path, _stat), nil
	} else {
		return "", _error
	}
}

func fingerprintSource_1 (_path string, _stat os.FileInfo) (string) {
	return NewFingerprinter () .StringWithLen (_path) .Int64 (int64 (_stat.Mode ())) .Int64 (_stat.Size ()) .Int64 (_stat.ModTime () .Unix ()) .Build ()
}




func resolveSourcePath_0 (_candidate string) (string, os.FileInfo, error) {
	if _candidate != "" {
//		logf ('d', 0x16563f01, "using candidate `%s`...", _candidate)
		return resolveSourcePath_2 (_candidate)
	} else {
//		logf ('d', 0xef5420f5, "searching candidate...")
		return resolveSourcePath_1 ()
	}
}


func resolveSourcePath_1 () (string, os.FileInfo, error) {
	
	_folders := make ([]string, 0, 128)
	_folders = append (_folders,
			".",
			path.Join (".", ".git"),
			path.Join (".", ".hg"),
			path.Join (".", ".svn"),
		)
	for _, _folder := range _folders {
		_folders = append (_folders, path.Join (_folder, "scripts"))
	}
	var _home string
	if _home_0, _error := os.UserHomeDir (); _error == nil {
		_home = _home_0
		_folders = append (_folders, _home)
	}
	
	_files := []string {
			
			"z-run",
			".z-run",
			"_z-run",
			
			"x-run",
			".x-run",
			"_x-run",
		}
	
	_candidates := make ([]string, 0, 16)
	
	for _, _folder := range _folders {
		if (_folder == _home) && (len (_candidates) > 0) {
			continue
		}
		for _, _file := range _files {
			_path := path.Join (_folder, _file)
			if _, _error := os.Lstat (_path); _error == nil {
				_candidates = append (_candidates, _path)
			} else if os.IsNotExist (_error) {
				// NOP
			} else {
				return "", nil, _error
			}
		}
	}
	
	if len (_candidates) == 0 {
		return "", nil, errorf (0x779f9a9f, "no sources found")
	} else if len (_candidates) > 1 {
		return "", nil, errorf (0x519bb041, "too many sources found: `%s`", _candidates)
	} else {
		return resolveSourcePath_2 (_candidates[0])
	}
}


func resolveSourcePath_2 (_path string) (string, os.FileInfo, error) {
	if _stat, _error := os.Stat (_path); _error == nil {
		if _path, _error := filepath.Abs (_path); _error == nil {
			return _path, _stat, nil
		} else {
			return "", nil, _error
		}
	} else if os.IsNotExist (_error) {
		return "", nil, errorf (0x4b0005de, "source does not exist `%s`", _path)
	} else {
		return "", nil, _error
	}
}




func resolveCache () (string, error) {
	var _cache string
	if _cache_0, _error := os.UserCacheDir (); _error == nil {
		_cache = _cache_0
	} else {
		return "", _error
	}
	if _error := os.MkdirAll (_cache, 0750); _error != nil {
		return "", _error
	}
	_cache = path.Join (_cache, "z-run")
	if _error := os.MkdirAll (_cache, 0750); _error != nil {
		return "", _error
	}
	return _cache, nil
}




func resolveLibrary (_candidate string, _context *Context) (LibraryStore, error) {
	
	_sources, _error := resolveSources (_candidate)
	if _error != nil {
		return nil, _error
	}
	
	var _environmentFingerprint string
	{
		_fingerprints := make ([]string, 0, len (_sources) * 2)
		if _fingerprint, _error := fingerprintSource_0 (_context.selfExecutable); _error == nil {
			_fingerprints = append (_fingerprints, "self-executable:" + _fingerprint)
		} else {
			return nil, _error
		}
		_fingerprints = append (_fingerprints, "workspace:" + _context.workspace)
		for _, _source := range _sources {
			_fingerprints = append (_fingerprints, "sources:" + _source.FingerprintMeta)
		}
		for _name, _value := range _context.cleanEnvironment {
			_fingerprint := NewFingerprinter () .StringWithLen (_name) .StringWithLen (_value) .Build ()
			_fingerprints = append (_fingerprints, "clean-environment:" + _fingerprint)
		}
		sort.Strings (_fingerprints)
		_environmentFingerprint = NewFingerprinter () .StringsWithLen (_fingerprints) .Build ()
	}
	
	var _cacheLibrary string
	
	if _context.cacheEnabled && (_context.cacheRoot != "") {
		_cacheLibrary = path.Join (_context.cacheRoot, _environmentFingerprint + ".cdb")
		if _library, _error := resolveLibraryCached (_cacheLibrary); _error == nil {
			if _fresh, _error := checkLibraryCached (_library); _error == nil {
				if _fresh {
//					logf ('d', 0xa33ecc63, "using library cached at `%s`;", _cacheLibrary)
					return _library, nil
				} else {
//					logf ('d', 0x8fc67fa1, "ignoring library cached at `%s`;", _cacheLibrary)
					_library.Close ()
				}
			} else {
				_library.Close ()
				return nil, _error
			}
		} else if ! os.IsNotExist (_error) {
			return nil, _error
		}
	}
	
	var _library *Library
//	logf ('i', 0xbd44916b, "parsing library from sources...")
	if _library_0, _error := parseLibrary (_sources, _environmentFingerprint, _context); _error == nil {
//		logf ('d', 0x71b45ebc, "parsed library from sources;")
		_library = _library_0
	} else {
		return nil, _error
	}
	
	if _context.cacheEnabled {
		if _error := doExportLibraryCdb (_library, _cacheLibrary, _context); _error == nil {
//			logf ('d', 0xdf78377c, "created library cached at `%s`;", _cacheLibrary)
			_library.url = _cacheLibrary
		} else {
			return nil, _error
		}
	}
	
	return _library, _error
}




func resolveLibraryCached (_path string) (LibraryStore, error) {
	if _store, _error := NewCdbStoreInput (_path); _error == nil {
		if _library, _error := NewLibraryStoreInput (_store, _path); _error == nil {
//			logf ('d', 0x63ae360d, "opened library cached at `%s`;", _cacheLibrary)
			return _library, nil
		} else {
			_store.Close ()
			return nil, _error
		}
	} else {
		return nil, _error
	}
}


func checkLibraryCached (_library LibraryStore) (bool, error) {
	var _sources LibrarySources
	if _sources_0, _error := _library.SelectSources (); _error == nil {
		_sources = _sources_0
	} else {
		return false, _error
	}
	for _, _source := range _sources {
		if _stat, _error := os.Stat (_source.Path); _error == nil {
			_fingerprint := fingerprintSource_1 (_source.Path, _stat)
			if _fingerprint != _source.FingerprintMeta {
				return false, nil
			}
		} else if os.IsNotExist (_error) {
			return false, nil
		} else {
			return false, _error
		}
	}
	return true, nil
}

