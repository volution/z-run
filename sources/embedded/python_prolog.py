#!/dev/null


################################################################################
################################################################################


def __zrun__create (Z = None, __import__ = __import__) :
	
	## --------------------------------------------------------------------------------
	
	if Z is None :
		Z = __import__ ("types") .ModuleType ("zrun")
	
	## --------------------------------------------------------------------------------
	
	PY = __import__ ("types") .ModuleType ("zrun_PY")
	PY.os = __import__ ("os")
	PY.sys = __import__ ("sys")
	PY.signal = __import__ ("signal")
	PY.subprocess = __import__ ("subprocess")
	PY.time = __import__ ("time")
	Z.py = PY
	
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
	def __zrun__zcmd (_scriptlet, *_arguments) :
		return Z._zexec_prepare (_scriptlet, _arguments)
	
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
	def __zrun__cmd (_scriptlet, *_arguments) :
		return Z._exec_prepare (_scriptlet, _arguments)
	
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
			_delegate = PY.os.spawnvpe
		else :
			_delegate = PY.os.spawnve
		if _wait :
			_outcome = _delegate (PY.os.P_WAIT, _executable, _arguments, _environment)
			if _panic and _outcome != 0 :
				Z.panic (0x3c14b9a0, "spawn `%s` `%s` failed with status: %d", _arguments[0], _arguments[1:], _outcome)
		else :
			_outcome = _delegate (PY.os.P_NOWAIT, _executable, _arguments, _environment)
			if _panic and _outcome <= 0 :
				Z.panic (0x36737d48, "spawn `%s` `%s` failed with error (%d): %s", _arguments[0], _arguments[1:], _outcome, PY.os.strerror (_outcome))
		return _outcome
	
	@_inject
	def __zrun__exec_0 (_descriptor) :
		_executable, _lookup, _arguments, _environment = _descriptor
		if _lookup :
			_delegate = PY.os.execvpe
		else :
			_delegate = PY.os.execve
		_delegate (_executable, _arguments, _environment)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __zrun__pipeline (_commands, _wait = True, _panic = True) :
		_count = len (_commands)
		if _count == 0 :
			z.panic (0x1b1812d7, "pipeline empty")
		_pipes = []
		_pipes.append ((None, None))
		for _index in range (_count - 1) :
			_pipes.append (PY.os.pipe ())
		_pipes.append ((None, None))
		_processes = []
		for _index in range (_count) :
			_executable, _lookup, _arguments, _environment = _commands[_index]
			_pipe_previous = _pipes[_index]
			_pipe_next = _pipes[_index + 1]
			_pipe_stdin = _pipe_previous[0]
			_pipe_stdout = _pipe_next[1]
			_process = PY.subprocess.Popen (
					_arguments,
					executable = _executable,
					env = _environment,
					stdin = _pipe_stdin,
					stdout = _pipe_stdout,
					stderr = None,
					close_fds = False,
					shell = False,
				)
			if _wait :
				_processes.append ((_index, _process, _arguments))
			else :
				_processes.append (_process.pid)
			if _pipe_stdin is not None :
				PY.os.close (_pipe_stdin)
			if _pipe_stdout is not None :
				PY.os.close (_pipe_stdout)
		if not _wait :
			return _processes
		_succeeded = True
		_terminated = 0
		if Z.python_version >= 303 :
			_signal_handler_old = PY.signal.signal (PY.signal.SIGCHLD, lambda _1, _2 : None)
		while True :
			for _process in _processes :
				if _process is None :
					continue
				_index, _process, _arguments = _process
				if _terminated == (_count - 1) or not Z.python_version >= 303 :
					_process.wait ()
				if _process.poll () is None :
					continue
				_terminated += 1
				if _process.returncode != 0 :
					_succeeded = False
					if _panic :
						Z.log_warning (0x76d05a67, "spawn `%s` `%s` failed with status: %d", _arguments[0], _arguments[1:], _process.returncode)
				_processes[_index] = None
			if _terminated == _count :
				break
			if Z.python_version >= 303 :
				PY.signal.sigtimedwait ([PY.signal.SIGCHLD], 6)
		if Z.python_version >= 303 :
			PY.signal.signal (PY.signal.SIGCHLD, _signal_handler_old)
		if _panic and not _succeeded :
			Z.panic (0x1d6fad91, "pipeline failed")
		return _succeeded
	
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
		PY.sys.exit (_status)
	
	@_inject
	def __zrun__panic (_code, _message, *_arguments) :
		Z._log_write ('!!', _code, _message, _arguments)
		Z.exit (1)
	
	@_inject
	def __zrun__sleep (_interval) :
		PY.time.sleep (_interval)
	
	## --------------------------------------------------------------------------------
	
	Z.pid = PY.os.getpid ()
	Z.environment = PY.os.environ
	
	Z.executable = Z.environment["ZRUN_EXECUTABLE"]
	Z.workspace = Z.environment["ZRUN_WORKSPACE"]
	Z.fingerprint = Z.environment["ZRUN_FINGERPRINT"]
	
	Z.stdin = PY.sys.stdin
	Z.stdout = PY.sys.stdout
	Z.stderr = PY.sys.stderr
	
	Z.log_warning_enabled = True
	Z.log_notice_enabled = True
	Z.log_debug_enabled = False
	
	Z.python_version = PY.sys.version_info[0] * 100 + PY.sys.version_info[1]
	
	## --------------------------------------------------------------------------------
	
	return Z


################################################################################
################################################################################


if __name__ == "__main__" :
	
	(lambda Z : Z.py.sys.modules.__setitem__ ("zrun", Z)) (__zrun__create ())
	import zrun
	sys = zrun.py.sys
	os = zrun.py.os
	
else :
	
	assert False, ("[b11458a7]  invalid module: `%s`" % __name__)


################################################################################
################################################################################

