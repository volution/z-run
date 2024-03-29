

package zrun


import "fmt"
import "os"
import "path"
import "path/filepath"
import "sort"
import "strings"

import . "github.com/volution/z-run/lib/library"
import . "github.com/volution/z-run/lib/store"
import . "github.com/volution/z-run/lib/common"
import . "github.com/volution/z-run/embedded"




func resolveSources (_candidate string, _workspace string, _lookupPaths []string, _execMode bool) ([]*Source, *Error) {
	
	_sources := make ([]*Source, 0, 128)
	
	_path, _stat, _error := resolveSourcePath_0 (_candidate, _workspace, _lookupPaths)
	if _error != nil {
		return nil, _error
	}
	
	_statMode := _stat.Mode ()
	switch {
		
		case _statMode.IsRegular () :
			if _source, _error := resolveSource_0 (_path, _stat); _error == nil {
				if _execMode {
					_source.Executable = false
				}
				_sources = append (_sources, _source)
			} else {
				return nil, _error
			}
		
		case _statMode.IsDir () :
			return nil, Errorf (0x8a04b23b, "not-implemented")
		
		default :
			return nil, Errorf (0xa35428a2, "invalid source `%s`", _path)
	}
	
	return _sources, nil
}


func resolveSource (_candidate string, _workspace string, _lookupPaths []string) (*Source, *Error) {
	if _path, _stat, _error := resolveSourcePath_0 (_candidate, _workspace, _lookupPaths); _error == nil {
		return resolveSource_0 (_path, _stat)
	} else {
		return nil, _error
	}
}


func resolveSource_0 (_path string, _stat os.FileInfo) (*Source, *Error) {
	_statMode := _stat.Mode ()
	if _statMode.IsRegular () {
		_source := & Source {
				Path : _path,
				Executable : (_statMode.Perm () & 0111) != 0,
				FingerprintMeta : fingerprintSource_1 (_path, _stat),
			}
		return _source, nil
	} else {
		return nil, Errorf (0x557961c4, "invalid source `%s`", _path)
	}
}




func fingerprintSource (_path string) (*Source, *Error) {
	if _path, _stat, _error := resolveSourcePath_2 (_path); _error == nil {
		_source := & Source {
				Path : _path,
				Executable : false,
				FingerprintMeta : fingerprintSource_1 (_path, _stat),
			}
		return _source, nil
	} else {
		return nil, _error
	}
}


func fingerprintSource_0 (_path string) (string, *Error) {
	if _stat, _error := os.Stat (_path); _error == nil {
		return fingerprintSource_1 (_path, _stat), nil
	} else {
		return "", Errorw (0x375f2514, _error)
	}
}


func fingerprintSource_1 (_path string, _stat os.FileInfo) (string) {
	return NewFingerprinter () .StringWithLen (_path) .Int64 (int64 (_stat.Mode ())) .Int64 (_stat.Size ()) .Int64 (_stat.ModTime () .Unix ()) .Build ()
}




func resolveSourcePath_0 (_candidate string, _workspace string, _lookupPaths []string) (string, os.FileInfo, *Error) {
	if _candidate != "" {
//		Logf ('d', 0x16563f01, "using candidate `%s`...", _candidate)
		return resolveSourcePath_2 (_candidate)
	} else {
//		Logf ('d', 0xef5420f5, "searching candidate...")
		return resolveSourcePath_1 (_workspace, _lookupPaths)
	}
}


func resolveSourcePath_1 (_workspace string, _lookupPaths []string) (string, os.FileInfo, *Error) {
	
	type folder struct {
		path string
		fallback bool
	}
	
	_folders := make ([]folder, 0, 128)
	
	if _workspace != "" {
		_folders = append (_folders, folder { _workspace, false })
		for _, _subfolder := range ResolveWorkspaceSubfolders {
			_folders = append (_folders, folder { path.Join (_workspace, _subfolder), false })
		}
		for _, _folder := range _folders {
			for _, _subfolder := range ResolveSourceSubfolders {
				_folders = append (_folders, folder { path.Join (_folder.path, _subfolder), _folder.fallback })
			}
		}
	}
	
	for _, _lookupPath := range _lookupPaths {
		_folders = append (_folders, folder { _lookupPath, true })
	}
	
	_candidates := make ([]string, 0, 16)
	
	for _, _folder := range _folders {
		if _folder.fallback && (len (_candidates) > 0) {
			continue
		}
		if _, _error := os.Lstat (_folder.path); _error == nil {
			if _stat, _error := os.Stat (_folder.path); _error == nil {
				if ! _stat.IsDir () {
					continue
				}
			} else {
				return "", nil, Errorw (0x8c151dc4, _error)
			}
		} else if os.IsNotExist (_error) {
			continue
		} else {
			return "", nil, Errorw (0x336d65f9, _error)
		}
		for _, _file := range ResolveSourceFiles {
			_path := path.Join (_folder.path, _file)
			if _, _error := os.Lstat (_path); _error == nil {
				if _, _error := os.Stat (_path); _error == nil {
					_candidates = append (_candidates, _path)
				} else {
					return "", nil, Errorw (0xe1cc83e6, _error)
				}
			} else if os.IsNotExist (_error) {
				// NOP
			} else {
				return "", nil, Errorw (0x49b2b24c, _error)
			}
		}
	}
	
	if len (_candidates) == 0 {
		return "", nil, Errorf (0x779f9a9f, "no sources found")
	} else if len (_candidates) > 1 {
		return "", nil, Errorf (0x519bb041, "too many sources found: `%s`", _candidates)
	} else {
		return resolveSourcePath_2 (_candidates[0])
	}
}


func resolveSourcePath_2 (_path string) (string, os.FileInfo, *Error) {
	if _stat, _error := os.Stat (_path); _error == nil {
		if _path, _error := filepath.Abs (_path); _error == nil {
			return _path, _stat, nil
		} else {
			return "", nil, Errorw (0x53c05222, _error)
		}
	} else if os.IsNotExist (_error) {
		return "", nil, Errorf (0x4b0005de, "source does not exist `%s`", _path)
	} else {
		return "", nil, Errorw (0x43066170, _error)
	}
}




func resolveCache () (string, *Error) {
	var _cache string
	if _cache_0, _error := os.UserCacheDir (); _error == nil {
		_cache = _cache_0
	} else {
		return "", Errorw (0x4d666a7f, _error)
	}
	if _error := os.MkdirAll (_cache, 0750); _error != nil {
		return "", Errorw (0xf214ed44, _error)
	}
	_cache = path.Join (_cache, "z-run")
	if _error := os.MkdirAll (_cache, 0750); _error != nil {
		return "", Errorw (0xa66a0341, _error)
	}
	return _cache, nil
}




func resolveLibrary (_candidate string, _context *Context, _lookupPaths []string, _execMode bool) (LibraryStore, *Error) {
	
	_sources, _error := resolveSources (_candidate, _context.workspace, _lookupPaths, _execMode)
	if _error != nil {
		return nil, _error
	}
	
	var _libraryIdentifier string
	{
		_fingerprints := make ([]string, 0, 1024)
		_fingerprints = append (_fingerprints, "build-version:" + BUILD_VERSION)
		_fingerprints = append (_fingerprints, "build-sources:" + BUILD_SOURCES_HASH)
		_fingerprints = append (_fingerprints, "workspace:" + _context.workspace)
		for _, _source := range _sources {
			_fingerprints = append (_fingerprints, "sources:" + _source.FingerprintMeta)
		}
		for _name, _value := range _context.cleanEnvironment {
			_fingerprint := NewFingerprinter () .StringWithLen (_name) .StringWithLen (_value) .Build ()
			_fingerprints = append (_fingerprints, "clean-environment:" + _fingerprint)
		}
		sort.Strings (_fingerprints)
		_libraryIdentifier = NewFingerprinter () .StringsWithLen (_fingerprints) .Build ()
//		Logf ('d', 0x2243d3c0, "%s", _libraryIdentifier)
//		for _, _fingerprint := range _fingerprints {
//			Logf ('d', 0x56370791, "%s", _fingerprint)
//		}
	}
	
	var _cacheLibrary string
	
	if _context.cacheEnabled {
		if _error := MakeCacheFolder (_context.cacheRoot, "libraries-cdb"); _error != nil {
			return nil, _error
		}
		_cacheLibrary = path.Join (_context.cacheRoot, "libraries-cdb", _libraryIdentifier + ".cdb")
		if _stat, _error := os.Lstat (_cacheLibrary); _error == nil {
			if _stat.Mode () .Type () == os.ModeSymlink {
				if _target, _error := filepath.EvalSymlinks (_cacheLibrary); _error == nil {
					_cacheLibrary = _target
				} else {
					return nil, Errorw (0x1d1aa28a, _error)
				}
			}
		} else if ! os.IsNotExist (_error) {
			return nil, Errorw (0x1dd01c9c, _error)
		}
		if _stat, _error := os.Lstat (_cacheLibrary); _error == nil {
			if ! _stat.Mode () .IsRegular () {
				return nil, Errorf (0x5b3ae1d5, "invalid library cached at `%s`;", _cacheLibrary)
			}
			if _library, _error := resolveLibraryCached (_cacheLibrary); _error == nil {
				if _fresh, _error := checkLibraryCached (_library); _error == nil {
					if _fresh {
//						Logf ('d', 0xa33ecc63, "using library cached at `%s`;", _cacheLibrary)
						return _library, nil
					} else {
//						Logf ('d', 0x8fc67fa1, "ignoring library cached at `%s`;", _cacheLibrary)
						_library.Close ()
					}
				} else {
					_library.Close ()
					return nil, _error
				}
			} else if (_error.Error != nil) && ! os.IsNotExist (_error.Error) {
				return nil, _error
			}
		} else if ! os.IsNotExist (_error) {
			return nil, Errorw (0x19404141, _error)
		}
	}
	
	var _library *Library
//	Logf ('i', 0xbd44916b, "parsing library from sources...")
	if _library_0, _error := parseLibrary (_sources, _libraryIdentifier, _context); _error == nil {
//		Logf ('d', 0x71b45ebc, "parsed library from sources;")
		_library = _library_0
	} else {
		return nil, _error
	}
	
	var _libraryFingerprint string
	if _fingerprint, _error := _library.Fingerprint (); _error == nil {
		_libraryFingerprint = _fingerprint
	} else {
		return nil, _error
	}
	
	if _context.cacheEnabled {
		_cacheLibraryLink := path.Join (_context.cacheRoot, "libraries-cdb", _libraryIdentifier + ".cdb")
		_cacheLibraryLinkTmp := fmt.Sprintf ("%s--%08x.tmp", _cacheLibraryLink, os.Getpid ())
		_cacheLibraryStable := path.Join (_context.cacheRoot, "libraries-cdb", _libraryFingerprint + ".cdb")
		if _error := doExportLibraryCdb (_library, _cacheLibraryStable, _context); _error == nil {
//			Logf ('d', 0xdf78377c, "created library cached link at `%s`;", _cacheLibraryLink)
//			Logf ('d', 0x43e263da, "created library cached stable at `%s`;", _cacheLibraryStable)
			_library.UrlSet (_cacheLibraryStable)
		} else {
			return nil, _error
		}
		if _error := os.Symlink (_libraryFingerprint + ".cdb", _cacheLibraryLinkTmp); _error != nil {
			return nil, Errorw (0xc8bbef7a, _error)
		}
		if _error := os.Rename (_cacheLibraryLinkTmp, _cacheLibraryLink); _error != nil {
			return nil, Errorw (0x39c785b9, _error)
		}
	}
	
	return _library, _error
}




func resolveLibraryCached (_path string) (LibraryStore, *Error) {
	_fileName := path.Base (_path)
	if ! strings.HasSuffix (_fileName, ".cdb") {
		return nil, Errorf (0x06574f0e, "invalid library cached file name `%s`", _path)
	}
	_fingerprint := _fileName[: len (_fileName) - 4]
	if _store, _error := NewCdbStoreInput (_path); _error == nil {
		if _library, _error := NewLibraryStoreInput (_store, _path, _fingerprint); _error == nil {
//			Logf ('d', 0x63ae360d, "opened library cached at `%s`;", _path)
			if _fingerprint_0, _error := _library.Fingerprint (); _error == nil {
				if _fingerprint_0 == _fingerprint {
					return _library, nil
				} else {
					return nil, Errorf (0xa0e14143, "invalid store")
				}
			} else {
				_store.Close ()
				return nil, _error
			}
		} else {
			_store.Close ()
			return nil, _error
		}
	} else {
		return nil, _error
	}
}


func checkLibraryCached (_library LibraryStore) (bool, *Error) {
	var _sources LibrarySources
	if _sources_0, _error := _library.SelectLibrarySources (); _error == nil {
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
			return false, Errorw (0x0f7764d9, _error)
		}
	}
	return true, nil
}




var ResolveWorkspaceSubfolders = []string {
		"__",
		".git",
		".hg",
		".svn",
		".bzr",
		".darcs",
	}

var ResolveSourceSubfolders = []string {
		"z-run", "_z-run", ".z-run",
		"zrun", "_zrun", ".zrun",
		"scriptlets", "_scriptlets", ".scriptlets",
		"scripts", "_scripts", ".scripts",
		"bin", "_bin", ".bin",
	}

var ResolveSourceFiles = []string {
			"z-run", "_z-run", ".z-run",
			"zrun", "_zrun", ".zrun",
	}

