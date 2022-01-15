

package zrun


import . "github.com/cipriancraciun/z-run/lib/store"
import . "github.com/cipriancraciun/z-run/lib/common"




type LibraryStore interface {
	
	SelectFingerprints () ([]string, *Error)
	ResolveFullByFingerprint (_fingerprint string) (*Scriptlet, *Error)
	ResolveMetaByFingerprint (_fingerprint string) (*Scriptlet, *Error)
	ResolveBodyByFingerprint (_fingerprint string) (string, bool, *Error)
	
	SelectLabels () ([]string, *Error)
	SelectLabelsAll () ([]string, *Error)
	ResolveFullByLabel (_label string) (*Scriptlet, *Error)
	ResolveMetaByLabel (_label string) (*Scriptlet, *Error)
	ResolveBodyByLabel (_label string) (string, bool, *Error)
	ResolveFingerprintByLabel (_label string) (string, bool, *Error)
	
	ResolveContextByIdentifier (_fingerprint string) (*ScriptletContext, bool, *Error)
	
	SelectLibrarySources () (LibrarySources, *Error)
	SelectLibraryContext () (*LibraryContext, *Error)
	
	Identifier () (string, *Error)
	Fingerprint () (string, *Error)
	
	Url () (string)
	
	Close () (*Error)
}


type LibraryStoreInput struct {
	store StoreInput
	url string
	instance string
}




func NewLibraryStoreInput (_store StoreInput, _url string, _instance string) (*LibraryStoreInput, *Error) {
	_library := & LibraryStoreInput {
			store : _store,
			url : _url,
			instance : _instance,
		}
	return _library, nil
}


func (_library *LibraryStoreInput) SelectFingerprints () ([]string, *Error) {
	var _value []string
	if _found, _error := _library.store.SelectObject (_library.instance, false, "scriptlets-indices", "fingerprints", &_value); _error == nil {
		if _found {
			return _value, nil
		} else {
			return nil, Errorf (0x7f976073, "invalid store")
		}
	} else {
		return nil, _error
	}
}

func (_library *LibraryStoreInput) SelectLabels () ([]string, *Error) {
	var _value []string
	if _found, _error := _library.store.SelectObject (_library.instance, false, "scriptlets-indices", "labels", &_value); _error == nil {
		if _found {
			return _value, nil
		} else {
			return nil, Errorf (0x64c3a996, "invalid store")
		}
	} else {
		return nil, _error
	}
}

func (_library *LibraryStoreInput) SelectLabelsAll () ([]string, *Error) {
	var _value []string
	if _found, _error := _library.store.SelectObject (_library.instance, false, "scriptlets-indices", "labels-all", &_value); _error == nil {
		if _found {
			return _value, nil
		} else {
			return nil, Errorf (0x4d9d3702, "invalid store")
		}
	} else {
		return nil, _error
	}
}


func (_library *LibraryStoreInput) ResolveFullByFingerprint (_fingerprint string) (*Scriptlet, *Error) {
	if _scriptlet, _error := _library.ResolveMetaByFingerprint (_fingerprint); _error == nil {
		if _scriptlet != nil {
			if _body, _found, _error := _library.ResolveBodyByFingerprint (_fingerprint); _error == nil {
				if _found {
					_scriptlet.Body = _body
				} else {
					return nil, Errorf (0x5c4c21e2, "invalid store")
				}
			} else {
				return nil, _error
			}
			if _scriptlet.ContextIdentifier != "" {
				if _context, _found, _error := _library.ResolveContextByIdentifier (_scriptlet.ContextIdentifier); _error == nil {
					if _found {
						_scriptlet.Context = _context
					} else {
						return nil, Errorf (0x656d6774, "invalid store")
					}
				}
			}
			return _scriptlet, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, _error
	}
}

func (_library *LibraryStoreInput) ResolveMetaByFingerprint (_fingerprint string) (*Scriptlet, *Error) {
	var _value *Scriptlet
	if _found, _error := _library.store.SelectObject (_library.instance, true, "scriptlets-meta", _fingerprint, &_value); _error == nil {
		if _found {
			return _value, nil
		} else {
			return nil, Errorf (0x008e4a04, "invalid store")
		}
	} else {
		return nil, _error
	}
}

func (_library *LibraryStoreInput) ResolveBodyByFingerprint (_fingerprint string) (string, bool, *Error) {
	if _found, _value, _error := _library.store.SelectRawString (_library.instance, true, "scriptlets-body", _fingerprint); _error == nil {
		if _found {
			return _value, _found, nil
		} else {
			return "", false, Errorf (0x4fd14583, "invalid store")
		}
	} else {
		return "", false, _error
	}
}


func (_library *LibraryStoreInput) ResolveFullByLabel (_label string) (*Scriptlet, *Error) {
	if _fingerprint, _found, _error := _library.ResolveFingerprintByLabel (_label); _error == nil {
		if _found {
			return _library.ResolveFullByFingerprint (_fingerprint)
		} else {
			return nil, nil
		}
	} else {
		return nil, _error
	}
}

func (_library *LibraryStoreInput) ResolveMetaByLabel (_label string) (*Scriptlet, *Error) {
	if _fingerprint, _found, _error := _library.ResolveFingerprintByLabel (_label); _error == nil {
		if _found {
			return _library.ResolveMetaByFingerprint (_fingerprint)
		} else {
			return nil, nil
		}
	} else {
		return nil, _error
	}
}

func (_library *LibraryStoreInput) ResolveBodyByLabel (_label string) (string, bool, *Error) {
	if _fingerprint, _found, _error := _library.ResolveFingerprintByLabel (_label); _error == nil {
		if _found {
			return _library.ResolveBodyByFingerprint (_fingerprint)
		} else {
			return "", false, nil
		}
	} else {
		return "", false, _error
	}
}

func (_library *LibraryStoreInput) ResolveFingerprintByLabel (_label string) (string, bool, *Error) {
	if _found, _value, _error := _library.store.SelectRawString (_library.instance, false, "scriptlets-fingerprint-by-label", _label); _error == nil {
		if _found {
			return _value, _found, nil
		} else {
			return "", false, nil
		}
	} else {
		return "", false, _error
	}
}

func (_library *LibraryStoreInput) ResolveContextByIdentifier (_identifier string) (*ScriptletContext, bool, *Error) {
	var _value *ScriptletContext
	if _found, _error := _library.store.SelectObject (_library.instance, false, "scriptlet-contexts-by-identifier", _identifier, &_value); _error == nil {
		if _found {
			return _value, _found, nil
		} else {
			return nil, false, nil
		}
	} else {
		return nil, false, _error
	}
}


func (_library *LibraryStoreInput) SelectLibrarySources () (LibrarySources, *Error) {
	var _value LibrarySources
	if _found, _error := _library.store.SelectObject (_library.instance, false, "library-meta", "library-sources", &_value); _error == nil {
		if _found {
			return _value, nil
		} else {
			return nil, Errorf (0x2986327f, "invalid store")
		}
	} else {
		return nil, _error
	}
}

func (_library *LibraryStoreInput) SelectLibraryContext () (*LibraryContext, *Error) {
	var _value *LibraryContext
	if _found, _error := _library.store.SelectObject (_library.instance, false, "library-meta", "library-context", &_value); _error == nil {
		if _found {
			return _value, nil
		} else {
			return nil, Errorf (0x2986327f, "9ddd4b2a store")
		}
	} else {
		return nil, _error
	}
}

func (_library *LibraryStoreInput) Identifier () (string, *Error) {
	if _found, _value, _error := _library.store.SelectRawString (_library.instance, false, "library-meta", "identifier"); _error == nil {
		if _found {
			return _value, nil
		} else {
			return "", Errorf (0x1b88b9d5, "invalid store")
		}
	} else {
		return "", _error
	}
}

func (_library *LibraryStoreInput) Fingerprint () (string, *Error) {
	if _found, _value, _error := _library.store.SelectRawString (_library.instance, false, "library-meta", "fingerprint"); _error == nil {
		if _found {
			return _value, nil
		} else {
			return "", Errorf (0x1b88b9d5, "invalid store")
		}
	} else {
		return "", _error
	}
}


func (_library *LibraryStoreInput) Url () (string) {
	return _library.url
}


func (_library *LibraryStoreInput) Close () (*Error) {
	return _library.store.Close ()
}

