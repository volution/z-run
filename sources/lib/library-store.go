

package lib




type LibraryStore interface {
	
	SelectFingerprints () ([]string, error)
	ResolveFullByFingerprint (_label string) (*Scriptlet, error)
	ResolveMetaByFingerprint (_label string) (*Scriptlet, error)
	ResolveBodyByFingerprint (_label string) (string, bool, error)
	
	SelectLabels () ([]string, error)
	ResolveFullByLabel (_label string) (*Scriptlet, error)
	ResolveMetaByLabel (_label string) (*Scriptlet, error)
	ResolveBodyByLabel (_label string) (string, bool, error)
	ResolveFingerprintByLabel (_label string) (string, bool, error)
	
	SelectSources () (LibrarySources, error)
	
	Url () (string)
}


type LibraryStoreInput struct {
	store StoreInput
	url string
}




func NewLibraryStoreInput (_store StoreInput, _url string) (*LibraryStoreInput, error) {
	_library := & LibraryStoreInput {
			store : _store,
			url : _url,
		}
	return _library, nil
}


func (_library *LibraryStoreInput) SelectFingerprints () ([]string, error) {
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

func (_library *LibraryStoreInput) SelectLabels () ([]string, error) {
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


func (_library *LibraryStoreInput) ResolveFullByFingerprint (_fingerprint string) (*Scriptlet, error) {
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

func (_library *LibraryStoreInput) ResolveMetaByFingerprint (_fingerprint string) (*Scriptlet, error) {
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

func (_library *LibraryStoreInput) ResolveBodyByFingerprint (_fingerprint string) (string, bool, error) {
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


func (_library *LibraryStoreInput) ResolveFullByLabel (_label string) (*Scriptlet, error) {
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

func (_library *LibraryStoreInput) ResolveMetaByLabel (_label string) (*Scriptlet, error) {
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

func (_library *LibraryStoreInput) ResolveBodyByLabel (_label string) (string, bool, error) {
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

func (_library *LibraryStoreInput) ResolveFingerprintByLabel (_label string) (string, bool, error) {
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


func (_library *LibraryStoreInput) SelectSources () (LibrarySources, error) {
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

