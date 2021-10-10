

package zrun



type StoreOutput interface {
	IncludeObject (_namespace string, _key string, _value interface{}) (*Error)
	IncludeRawString (_namespace string, _key string, _value string) (*Error)
	IncludeRawBytes (_namespace string, _key string, _value []byte) (*Error)
	Commit () (*Error)
}

type StoreInput interface {
	SelectObject (_namespace string, _key string, _value interface{}) (bool, *Error)
	SelectRawString (_namespace string, _key string) (bool, string, *Error)
	SelectRawBytes (_namespace string, _key string) (bool, []byte, *Error)
	Close () (*Error)
}

