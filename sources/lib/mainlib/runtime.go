

package mainlib


import "log"
import "os"
import "runtime"
import "runtime/debug"

import . "github.com/cipriancraciun/z-run/lib/common"




func InitializeMainRuntime () (*Error) {
	
	runtime.GOMAXPROCS (1)
	debug.SetMaxThreads (16)
	debug.SetMaxStack (128 * 1024)
	debug.SetGCPercent (500)
	
	log.SetFlags (0)
	
	return nil
}




func ExitMainSucceeded () (*Error) {
	os.Exit (0)
	return Errorf (0x62b9e61e, "fallthrough")
}

func ExitMainFailed () (*Error) {
	os.Exit (1)
	return Errorf (0xaf3f9be2, "fallthrough")
}

