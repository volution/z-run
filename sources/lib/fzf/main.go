

package fzf


import . "github.com/cipriancraciun/z-run/lib/mainlib"
import . "github.com/cipriancraciun/z-run/lib/common"




func FzfPreMain (_embedded bool, _fullscreen bool) () {
	
	var _name string
	if _embedded {
		_name = "z-select"
	} else {
		_name = "z-fzf"
	}
	
	_descriptor := & MainDescriptor {
			Main : func (_context *MainContext) (*Error) {
					return FzfMain (_embedded, _fullscreen, _context.Arguments, _context.EnvironmentMap)
				},
			ExecutableName : _name,
			HelpTxt : "__custom__",
		}
	
	PreMainWith (_descriptor)
}

