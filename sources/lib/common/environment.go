

package common


import "sort"
import "strings"




func EnvironmentMapToList (_map map[string]string) ([]string) {
	_list := make ([]string, 0, len (_map))
	for _name, _value := range _map {
		_pair := _name + "=" + _value
		_list = append (_list, _pair)
	}
	sort.Strings (_list)
	return _list
}


func EnvironmentListToMap (_list []string) (map[string]string) {
	_map := make (map[string]string, len (_list))
	for _, _pair := range _list {
		_split := strings.IndexByte (_pair, '=')
		if _split < 0 {
			continue
		}
		_name := _pair[:_split]
		_value := _pair[_split + 1:]
		_map[_name] = _value
	}
	return _map
}

