

package lib


import "os"
import "path"
import "sort"




func resolveSources (_candidate string) ([]*Source, error) {
	
	_sources := make ([]*Source, 0, 128)
	
	_candidate, _stat, _error := resolveSourcesPath_0 (_candidate)
	if _error != nil {
		return nil, _error
	}
	
	_statMode := _stat.Mode ()
	switch {
		
		case _statMode.IsRegular () :
			_fingerprint := NewFingerprinter () .StringWithLen (_candidate) .Int64 (_stat.Size ()) .Int64 (_stat.ModTime () .Unix ()) .Build ()
			_source := & Source {
					Path : _candidate,
					Executable : (_statMode.Perm () & 0111) != 0,
					FingerprintMeta : _fingerprint,
				}
			_sources = append (_sources, _source)
		
		case _statMode.IsDir () :
			return nil, errorf (0x8a04b23b, "not-implemented")
		
		default :
			return nil, errorf (0xa35428a2, "invalid source `%s`", _candidate)
	}
	
	return _sources, nil
}


func resolveSourcesPath_0 (_candidate string) (string, os.FileInfo, error) {
	if _candidate != "" {
//		logf ('d', 0x16563f01, "using candidate `%s`...", _candidate)
		return resolveSourcesPath_2 (_candidate)
	} else {
//		logf ('d', 0xef5420f5, "searching candidate...")
		return resolveSourcesPath_1 ()
	}
}


func resolveSourcesPath_1 () (string, os.FileInfo, error) {
	
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
		return resolveSourcesPath_2 (_candidates[0])
	}
}


func resolveSourcesPath_2 (_candidate string) (string, os.FileInfo, error) {
	if _stat, _error := os.Stat (_candidate); _error == nil {
		return _candidate, _stat, nil
	} else if os.IsNotExist (_error) {
		return "", nil, errorf (0x4b0005de, "source does not exist `%s`", _candidate)
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
	
	var _sourcesFingerprint string
	{
		_fingerprints := make ([]string, 0, len (_sources) * 2)
		for _, _source := range _sources {
			_fingerprints = append (_fingerprints, _source.FingerprintMeta)
			if _source.FingerprintData != "" {
				_fingerprints = append (_fingerprints, _source.FingerprintData)
			}
		}
		for _name, _value := range _context.cleanEnvironment {
			_fingerprint := NewFingerprinter () .String ("d33041e6571901d0a5a6dfbde7c7a312") .StringWithLen (_name) .StringWithLen (_value) .Build ()
			_fingerprints = append (_fingerprints, _fingerprint)
		}
		sort.Strings (_fingerprints)
		_fingerprinter := NewFingerprinter ()
		for _, _fingerprint := range _fingerprints {
			_fingerprinter.StringWithLen (_fingerprint)
		}
		_sourcesFingerprint = _fingerprinter.Build ()
	}
	
	var _cacheLibrary string
	
	if _context.cacheEnabled && (_context.cacheRoot != "") {
		_cacheLibrary = path.Join (_context.cacheRoot, _sourcesFingerprint + ".cdb")
		if _store, _error := resolveLibraryCached (_cacheLibrary); _error == nil {
			return _store, nil
		} else if ! os.IsNotExist (_error) {
			return nil, _error
		}
	}
	
	var _library *Library
	logf ('i', 0xbd44916b, "parsing library from sources...")
	if _library_0, _error := parseLibrary (_sources, _sourcesFingerprint, _context); _error == nil {
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
//		logf ('d', 0x63ae360d, "opened library cachad at `%s`;", _cacheLibrary)
		return NewLibraryStoreInput (_store, _path)
	} else {
		return nil, _error
	}
}

