

package input


import . "github.com/cipriancraciun/z-run/lib/mainlib"
import . "github.com/cipriancraciun/z-run/lib/common"




func InputPreMain () () {
	
	_descriptor := & MainDescriptor {
			Main : func (_context *MainContext) (*Error) {
					return InputMain (_context.Arguments, _context.EnvironmentMap)
				},
			ExecutableName : "z-input",
			HelpTxt : "__custom__",
		}
	
	PreMainWith (_descriptor)
}

