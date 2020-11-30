#!/dev/null

import sys
import os

def zrun () : raise Exception (("0bf229bc", "not-implemented"))
def zrun_panic (_code, _message, **_arguments) :
	sys.stderr.write (("[z-run:%08d] [!!] [%08x]  " % (os.getpid (), _code)) + (_message % _arguments) + "\n")
	sys.exit (1)

if __name__ != "__main__" :
	zrun_panic (0xdd55192b, "invalid state!")

