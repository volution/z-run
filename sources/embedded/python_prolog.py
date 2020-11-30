#!/dev/null


################################################################################
################################################################################


def __zrun__inject (Z, __import__ = __import__) :
	
	## --------------------------------------------------------------------------------
	
	Z.os = __import__ ("os")
	Z.sys = __import__ ("sys")
	Z.shutil = __import__ ("shutil")
	Z.time = __import__ ("time")
	
	## --------------------------------------------------------------------------------
	
	def _inject (_function) :
		_name = _function.__name__
		if _name.startswith ("__zrun__") :
			_name = _name[8:]
		else :
			assert False, ("[83cec849]  invalid inject name: `%s`" % _name)
		Z.__dict__[_name] = _function
		return _function
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __zrun__zspawn (_scriptlet, *_arguments, **_options) :
		_descriptor = Z._zexec_prepare (_scriptlet, _arguments)
		return Z.spawn_0 (_descriptor, **_options)
	
	@_inject
	def __zrun__zexec (_scriptlet, *_arguments, **_options) :
		_descriptor = Z._zexec_prepare (_scriptlet, _arguments)
		return Z.exec_0 (_descriptor, **_options)
	
	@_inject
	def __zrun___zexec_prepare (_scriptlet, _arguments) :
		_executable = Z.executable
		if not _scriptlet.startswith ("::") :
			Z.panic (0xbd1641c7, "invalid scriptlet: `%s`", _scriptlet)
		_arguments_all = ["[z-run]", _scriptlet]
		_arguments_all.extend (_arguments)
		_environment = {}
		_environment.update (Z.environment)
		return _executable, False, _arguments_all, _environment
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __zrun__spawn (_scriptlet, *_arguments, **_options) :
		_descriptor = Z._exec_prepare (_scriptlet, _arguments)
		return Z.spawn_0 (_descriptor, **_options)
	
	@_inject
	def __zrun__exec (_executable, *_arguments, **_options) :
		_descriptor = Z._exec_prepare (_executable, _arguments)
		return Z.exec_0 (_descriptor, **_options)
	
	@_inject
	def __zrun___exec_prepare (_executable, _arguments) :
		_arguments_all = [_executable]
		_arguments_all.extend (_arguments)
		_environment = {}
		_environment.update (Z.environment)
		return _executable, True, _arguments_all, _environment
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __zrun__spawn_0 (_descriptor, _wait = True, _panic = True) :
		_executable, _lookup, _arguments, _environment = _descriptor
		if _lookup :
			_delegate = Z.os.spawnvpe
		else :
			_delegate = Z.os.spawnve
		if _wait :
			_outcome = _delegate (Z.os.P_WAIT, _executable, _arguments, _environment)
			if _panic and _outcome != 0 :
				Z.panic (0x3c14b9a0, "spawn `%s` `%s` failed with status: %d", _arguments[0], _arguments[1:], _outcome)
		else :
			_outcome = _delegate (Z.os.P_NOWAIT, _executable, _arguments, _environment)
			if _panic and _outcome <= 0 :
				Z.panic (0x36737d48, "spawn `%s` `%s` failed with error (%d): %s", _arguments[0], _arguments[1:], _outcome, Z.os.strerror (_outcome))
		return _outcome
	
	@_inject
	def __zrun__exec_0 (_descriptor) :
		_executable, _lookup, _arguments, _environment = _descriptor
		if _lookup :
			_delegate = Z.os.execvpe
		else :
			_delegate = Z.os.execve
		_delegate (_executable, _arguments, _environment)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __zrun__log_error (_code, _message, *_arguments) :
		Z._log_write ('ee', _code, _message, _arguments)
	
	@_inject
	def __zrun__log_warning (_code, _message, *_arguments) :
		if not Z.log_warning_enabled : return
		Z._log_write ('ww', _code, _message, _arguments)
	
	@_inject
	def __zrun__log_notice (_code, _message, *_arguments) :
		if not Z.log_warning_enabled or not Z.log_notice_enabled : return
		Z._log_write ('ii', _code, _message, _arguments)
	
	@_inject
	def __zrun__log_debug (_code, _message, *_arguments) :
		if not Z.log_warning_enabled or not Z.log_notice_enabled or not Z.log_debug_enabled : return
		Z._log_write ('dd', _code, _message, _arguments)
	
	@_inject
	def __zrun___log_write (_slug, _code, _message, _arguments) :
		Z.stderr.write (("[z-run:%08d] [%s] [%08x]  " % (Z.pid, _slug, _code)) + (_message % _arguments) + "\n")
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __zrun__exit (_status) :
		Z.sys.exit (_status)
	
	@_inject
	def __zrun__panic (_code, _message, *_arguments) :
		Z._log_write ('!!', _code, _message, _arguments)
		Z.exit (1)
	
	@_inject
	def __zrun__sleep (_interval) :
		Z.time.sleep (_interval)
	
	## --------------------------------------------------------------------------------
	
	Z.pid = Z.os.getpid ()
	Z.environment = Z.os.environ
	Z.executable = Z.environment["ZRUN_EXECUTABLE"]
	Z.workspace = Z.environment["ZRUN_WORKSPACE"]
	Z.fingerprint = Z.environment["ZRUN_FINGERPRINT"]
	
	Z.stdin = Z.sys.stdin
	Z.stdout = Z.sys.stdout
	Z.stderr = Z.sys.stderr
	
	Z.log_warning_enabled = True
	Z.log_notice_enabled = True
	Z.log_debug_enabled = False
	
	## --------------------------------------------------------------------------------
	
	return Z


################################################################################
################################################################################


if __name__ == "__main__" :
	
#!	import sys
#!	import os
	
	def __zrun () :
		assert False, "[c92bb585]"
	
	zrun = __zrun__inject (__zrun)
	
else :
	
	assert False, ("[b11458a7]  invalid module: `%s`" % __name__)


################################################################################
################################################################################

