

package zrun




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
	
	SelectSources () (LibrarySources, *Error)
	
	Url () (string)
	
	Close () (*Error)
}


type LibraryStoreInput struct {
	store StoreInput
	url string
}




func NewLibraryStoreInput (_store StoreInput, _url string) (*LibraryStoreInput, *Error) {
	_library := & LibraryStoreInput {
			store : _store,
			url : _url,
		}
	return _library, nil
}


func (_library *LibraryStoreInput) SelectFingerprints () ([]string, *Error) {
	var _value []string
	if _found, _error := _library.store.Select ("scriptlets-indices", "fingerprints", &_value); _error == nil {
		if _found {
			return _value, nil
		} else {
			return nil, errorf (0x7f976073, "invalid store")
		}
	} else {
		return nil, _error
	}
}

func (_library *LibraryStoreInput) SelectLabels () ([]string, *Error) {
	var _value []string
	if _found, _error := _library.store.Select ("scriptlets-indices", "labels", &_value); _error == nil {
		if _found {
			return _value, nil
		} else {
			return nil, errorf (0x64c3a996, "invalid store")
		}
	} else {
		return nil, _error
	}
}

func (_library *LibraryStoreInput) SelectLabelsAll () ([]string, *Error) {
	var _value []string
	if _found, _error := _library.store.Select ("scriptlets-indices", "labels-all", &_value); _error == nil {
		if _found {
			return _value, nil
		} else {
			return nil, errorf (0x4d9d3702, "invalid store")
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
					return _scriptlet, nil
				} else {
					return nil, errorf (0x5c4c21e2, "invalid store")
				}
			} else {
				return nil, _error
			}
		} else {
			return nil, nil
		}
	} else {
		return nil, _error
	}
}

func (_library *LibraryStoreInput) ResolveMetaByFingerprint (_fingerprint string) (*Scriptlet, *Error) {
	var _value *Scriptlet
	if _found, _error := _library.store.Select ("scriptlets-meta", _fingerprint, &_value); _error == nil {
		if _found {
			return _value, nil
		} else {
			return nil, errorf (0x008e4a04, "invalid store")
		}
	} else {
		return nil, _error
	}
}

func (_library *LibraryStoreInput) ResolveBodyByFingerprint (_fingerprint string) (string, bool, *Error) {
	var _value string
	if _found, _error := _library.store.Select ("scriptlets-body", _fingerprint, &_value); _error == nil {
		if _found {
			return _value, _found, nil
		} else {
			return "", false, errorf (0x4fd14583, "invalid store")
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
	var _value string
	if _found, _error := _library.store.Select ("scriptlets-fingerprint-by-label", _label, &_value); _error == nil {
		if _found {
			return _value, _found, nil
		} else {
			return "", false, nil
		}
	} else {
		return "", false, _error
	}
}


func (_library *LibraryStoreInput) SelectSources () (LibrarySources, *Error) {
	var _value LibrarySources
	if _found, _error := _library.store.Select ("library-meta", "sources", &_value); _error == nil {
		if _found {
			return _value, nil
		} else {
			return nil, errorf (0x2986327f, "invalid store")
		}
	} else {
		return nil, _error
	}
}


func (_library *LibraryStoreInput) Url () (string) {
	return _library.url
}


func (_library *LibraryStoreInput) Close () (*Error) {
	return _library.store.Close ()
}

