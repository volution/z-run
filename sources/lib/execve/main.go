

package execve


import . "github.com/volution/z-run/lib/mainlib"
import . "github.com/volution/z-run/lib/common"




func ExecvePreMain () () {
	
	_descriptor := & MainDescriptor {
			Main : func (_context *MainContext) (*Error) {
					return ExecveMainWithFlags (_context.Flags.(*ExecveMainFlags), _context.EnvironmentMap)
				},
			Flags : & ExecveMainFlags {},
			ExecutableName : "z-execve",
			HelpTxt : "__custom__",
		}
	
	PreMainWith (_descriptor)
}

