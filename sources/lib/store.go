

package zrun



type StoreOutput interface {
	Include (_namespace string, _key string, _value interface{}) (*Error)
	Commit () (*Error)
}

type StoreInput interface {
	Select (_namespace string, _key string, _value interface{}) (bool, *Error)
	Close () (*Error)
}

