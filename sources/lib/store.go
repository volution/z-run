

package zrun



type StoreOutput interface {
	IncludeObject (_instance string, _global bool, _namespace string, _key string, _value interface{}) (*Error)
	IncludeRawString (_instance string, _global bool, _namespace string, _key string, _value string) (*Error)
	IncludeRawBytes (_instance string, _global bool, _namespace string, _key string, _value []byte) (*Error)
	Commit () (*Error)
}

type StoreInput interface {
	SelectObject (_instance string, _global bool, _namespace string, _key string, _value interface{}) (bool, *Error)
	SelectRawString (_instance string, _global bool, _namespace string, _key string) (bool, string, *Error)
	SelectRawBytes (_instance string, _global bool, _namespace string, _key string) (bool, []byte, *Error)
	Close () (*Error)
}

