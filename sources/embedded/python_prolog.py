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
	PY.errno = __import__ ("errno")
	PY.subprocess = __import__ ("subprocess")
	PY.time = __import__ ("time")
	PY.path = __import__ ("os.path")
	PY.stat = __import__ ("stat")
	PY.re = __import__ ("re")
	
	if PY.sys.version_info[0] > 2 :
		PY.basestring = str
	
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
		_environment = { _name : Z.environment[_name] for _name in Z.environment }
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
		_environment = { _name : Z.environment[_name] for _name in Z.environment }
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
				Z.panic ((_panic, 0x7d3900c4), "spawn `%s` `%s` failed with status: %d", _arguments[0], _arguments[1:], _outcome)
		else :
			_outcome = _delegate (PY.os.P_NOWAIT, _executable, _arguments, _environment)
			if _panic and _outcome <= 0 :
				Z.panic ((_panic, 0x56e47a07), "spawn `%s` `%s` failed with error (%d): %s", _arguments[0], _arguments[1:], _outcome, PY.os.strerror (_outcome))
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
			Z.panic ((_panic, 0x1d6fad91), "pipeline failed")
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
		if isinstance (_code, tuple) :
			_code_0 = 0xee4006b2
			for _code_0 in _code :
				if _code_0 is not False and _code_0 is not None :
					break
			_code = _code_0
		Z._log_write ('!!', _code, _message, _arguments)
		Z.exit (1)
	
	@_inject
	def __zrun__sleep (_interval) :
		PY.time.sleep (_interval)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __zrun__expect_no_arguments () :
		return Z.expect_arguments (_exact = 0)
	
	@_inject
	def __zrun__expect_arguments (_exact = None, _min = None, _max = None) :
		if _exact is not None :
			assert _min is None and _max is None, "[6ace19bf]"
			assert _exact >= 0, "[bfdc6527]"
		else :
			assert _min is not None or _max is not None, "[91787dc8]"
			assert _min is None or _min >= 0, "[19c3b0ad]"
			assert _max is None or _max >= 0, "[236d9471]"
			assert _min is None or _min < _max, "[e43f6498]"
			assert _max is None or _min < _max, "[0662ca21]"
		_actual = len (Z.arguments)
		if _exact is not None and _actual != _exact :
			if _exact == 0 :
				Z.panic (0x2aa5c1b7, "invalid arguments:  expected none, received %d!", _actual)
			else :
				Z.panic (0x2aa5c1b7, "invalid arguments:  expected exactly %d, received %d!", _exact, _actual)
		elif _min is not None and _actual < _min :
			Z.panic (0x5d1441a5, "invalid arguments:  expected at least %d, received %d!", _min, _actual)
		elif _max is not None and _actual > _max :
			Z.panic (0x6499e7a6, "invalid arguments:  expected at most %d, received %d!", _max, _actual)
		if _exact is not None :
			return tuple (Z.arguments)
		else :
			return tuple (Z.arguments[:_min]) + tuple (list (Z.arguments[_min:]))
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __zrun__enforce_regex (_value, _pattern, _code = None, _message = None) :
		if _message is None : _message = "enforcement failed"
		_pattern = Z.regex (_pattern)
		if not isinstance (_value, PY.basestring) :
			if _code is None : _code = 0x00a780ed
			Z.panic (_code, _message)
		if _pattern.match (_value) is None :
			if _code is None : _code = 0x9c922f7e
			Z.panic (_code, _message)
		return _value
	
	@_inject
	def __zrun__regex (_pattern) :
		return PY.re.compile (_pattern, PY.re.ASCII | PY.re.DOTALL)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __zrun__path (_path, _absolute = False, _canonical = False, _relative = None) :
		if not isinstance (_path, PY.basestring) and not isinstance (_path, bytes) :
			_path = PY.os.path.join (*_path)
		_path = PY.os.path.normpath (_path)
		if _path.startswith ("//") :
			_path = "/" + _path.lstrip ("/")
		if _absolute :
			_path = PY.os.path.abspath (_path)
		if _canonical :
			_path = PY.os.path.realpath (_path)
		if _relative is not None :
			_path = PY.os.path.relpath (_relative)
		return _path
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __zrun__exists (_path, _follow = True, _panic = False) :
		_stat = Z.stat (_path, _follow)
		if _stat is None and _panic :
			Z.panic ((_panic, 0x383d3cc5), "file-system path does not exist `%s`", _path)
		return _stat is not None
	
	@_inject
	def __zrun__not_exists (_path, _follow = True, _panic = False) :
		_stat = Z.stat (_path, _follow)
		if _stat is not None and _panic :
			Z.panic ((_panic, 0x9064abfc), "file-system path already exists `%s`", _path)
		return _stat is None
	
	@_inject
	def __zrun__is_file (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __zrun__is_file_empty (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode) and _stat.st_size == 0), _follow, _panic)
	
	@_inject
	def __zrun__is_file_not_empty (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode) and _stat.st_size > 0), _follow, _panic)
	
	@_inject
	def __zrun__is_folder (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISDIR (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __zrun__is_file_or_folder (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode) or PY.stat.S_ISDIR (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __zrun__is_symlink (_path, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISLNK (_stat.st_mode)), False, _panic)
	
	@_inject
	def __zrun__is_pipe (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISFIFO (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __zrun__is_socket (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISSOCK (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __zrun__is_dev_block (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISBLK (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __zrun__is_dev_char (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISCHR (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __zrun__is_special (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISFIFO (_stat.st_mode) or PY.stat.S_ISSOCK (_stat.st_mode) or PY.stat.S_ISBLK (_stat.st_mode) or PY.stat.S_ISCHR (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __zrun___stat_check (_path, _check, _follow = True, _panic = False) :
		_stat = Z.stat (_path, _follow)
		if _stat is None :
			if _panic :
				Z.panic ((_panic, 0xa23c577e), "file-system path not found `%s`", _path)
			else :
				return None
		if _check (_stat) :
			return True
		else :
			if _panic :
				Z.panic ((_panic, 0xfdbdc9a5), "file-system stat check failed `%s`", _path)
			else :
				return False
	
	@_inject
	def __zrun__stat (_path, _follow = True) :
		if _follow :
			_delegate = PY.os.stat
		else :
			_delegate = PY.os.lstat
		try :
			_stat = _delegate (_path)
		except OSError as _error :
			if _error.errno == PY.errno.ENOENT :
				_stat = None
			else :
				raise
		return _stat
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __zrun__mkdir (_path, _mode = None, _recurse = False) :
		if _mode is None : _mode = 0o777
		if _recurse :
			PY.os.makedirs (_path, _mode, True)
		else :
			PY.os.mkdir (_path, _mode)
	
	## --------------------------------------------------------------------------------
	
	class Z_environment :
		__str__ = PY.os.environ.__str__
		__repr__ = PY.os.environ.__repr__
		__getattr__ = PY.os.environ.__getitem__
		__setattr__ = PY.os.environ.__setitem__
		__len__ = PY.os.environ.__len__
		__getitem__ = PY.os.environ.__getitem__
		__setitem__ = PY.os.environ.__setitem__
		__delitem__ = PY.os.environ.__delitem__
		__contains__ = PY.os.environ.__contains__
		__iter__ = PY.os.environ.__iter__
	
	## --------------------------------------------------------------------------------
	
	Z.pid = PY.os.getpid ()
	Z.arguments = tuple (PY.sys.argv[1:])
	Z.environment = Z_environment ()
	
	Z.executable = Z.environment.ZRUN_EXECUTABLE
	Z.workspace = Z.environment.ZRUN_WORKSPACE
	Z.fingerprint = Z.environment.ZRUN_FINGERPRINT
	
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
	
	zrun.py.signal.signal (zrun.py.signal.SIGINT, (lambda _0, _1 : zrun.panic (0x6c751732, "scriptlet interrupted with SIGINT;  aborting!")))
	zrun.py.signal.signal (zrun.py.signal.SIGTERM, (lambda _0, _1 : zrun.panic (0xb5067479, "scriptlet interrupted with SIGTERM;  aborting!")))
	zrun.py.signal.signal (zrun.py.signal.SIGQUIT, (lambda _0, _1 : zrun.panic (0x921de146, "scriptlet interrupted with SIGQUIT;  aborting!")))
	zrun.py.signal.signal (zrun.py.signal.SIGHUP, (lambda _0, _1 : zrun.panic (0xe2e4c7c5, "scriptlet interrupted with SIGHUP;  aborting!")))
	zrun.py.signal.signal (zrun.py.signal.SIGPIPE, (lambda _0, _1 : zrun.panic (0xed8191f4, "scriptlet interrupted with SIGPIPE;  aborting!")))
	zrun.py.signal.signal (zrun.py.signal.SIGABRT, (lambda _0, _1 : zrun.panic (0xd6af6d5b, "scriptlet interrupted with SIGABRT;  aborting!")))
	
	sys = zrun.py.sys
	os = zrun.py.os


################################################################################
################################################################################

