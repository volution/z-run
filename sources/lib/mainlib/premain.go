

package mainlib


import "os"

import . "github.com/volution/z-run/lib/common"




type MainDescriptor struct {
	
	Main func (_context *MainContext) (*Error)
	
	ExecutableName string
	ExecutableEnvironmentHint string
	
	Flags interface{}
	
	HelpTxt string
	
	ManualTxt string
	ManualHtml string
	ManualMan string
}


type MainContext struct {
	
	Descriptor *MainDescriptor
	
	Executable0 string
	Executable string
	
	Argument0 string
	Arguments []string
	
	EnvironmentMap map[string]string
	EnvironmentList []string
	
	Flags interface{}
}




func PreMainWith (_descriptor *MainDescriptor) () {
	
	if _error := InitializeMainRuntime (); _error != nil {
		panic (AbortError (_error))
	}
	
	LogfTool = _descriptor.ExecutableName
	
	if _error := CleanMainEnvironment (); _error != nil {
		panic (AbortError (_error))
	}
	
	var _executable0, _executable string
	if _executable0_0, _executable_0, _error := ResolveMainExecutable (_descriptor.ExecutableName, _descriptor.ExecutableEnvironmentHint); _error == nil {
		_executable0 = _executable0_0
		_executable = _executable_0
	} else {
		panic (AbortError (_error))
	}
	
	var _environmentMap map[string]string
	var _environmentList []string
	if _environmentMap_0, _environmentList_0, _error := ResolveMainEnvironment (); _error == nil {
		_environmentMap = _environmentMap_0
		_environmentList = _environmentList_0
	} else {
		panic (AbortError (_error))
	}
	
	if _error := InterceptMainSpecialFlags (_descriptor.ExecutableName, _executable0, _executable, _descriptor.HelpTxt, _descriptor.ManualTxt, _descriptor.ManualHtml, _descriptor.ManualMan, _environmentMap); _error != nil {
		panic (AbortError (_error))
	}
	
	var _argument0 string
	var _arguments []string
	if _argument0_0, _arguments_0, _error := ResolveMainArguments (_executable0, _executable); _error == nil {
		_argument0 = _argument0_0
		_arguments = _arguments_0
	} else {
		panic (AbortError (_error))
	}
	
	if _descriptor.Flags != nil {
		if _error := ResolveMainFlags (_descriptor.ExecutableName, _arguments, _environmentMap, _descriptor.Flags, os.Stderr); _error != nil {
			panic (AbortError (_error))
		}
	}
	
	_context := & MainContext {
			
			Descriptor : _descriptor,
			
			Flags : _descriptor.Flags,
			
			Executable0 : _executable0,
			Executable : _executable,
			
			Argument0 : _argument0,
			Arguments : _arguments,
			
			EnvironmentMap : _environmentMap,
			EnvironmentList : _environmentList,
		}
	
	if _error := _descriptor.Main (_context); _error != nil {
		panic (AbortError (_error))
	} else {
		panic (ExitMainSucceeded ())
	}
}

