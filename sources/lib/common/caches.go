

package common


import "path"
import "os"




func MakeCacheFolder (_cacheRoot string, _cacheFolder string) (*Error) {
	_cache := path.Join (_cacheRoot, _cacheFolder)
	if _error := os.MkdirAll (_cache, 0750); _error != nil {
		return Errorw (0x6f530744, _error)
	}
	return nil
}

