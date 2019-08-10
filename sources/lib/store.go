

package zrun



type StoreOutput interface {
	Include (_namespace string, _key string, _value interface{}) (error)
	Commit () (error)
}

type StoreInput interface {
	Select (_namespace string, _key string, _value interface{}) (bool, error)
	Close () (error)
}

