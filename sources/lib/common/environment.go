

package common


import "sort"




func EnvironmentMapToList (_map map[string]string) ([]string) {
	_list := make ([]string, 0, len (_map))
	for _key, _value := range _map {
		_pair := _key + "=" + _value
		_list = append (_list, _pair)
	}
	sort.Strings (_list)
	return _list
}

