#!/dev/null


################################################################################
################################################################################


def __Z__create (Z = None, __import__ = __import__) :
	
	## --------------------------------------------------------------------------------
	
	if Z is None :
		Z = __import__ ("types") .ModuleType ("Z")
	
	## --------------------------------------------------------------------------------
	
	PY = __import__ ("types") .ModuleType ("PY")
	PY.os = __import__ ("os")
	PY.sys = __import__ ("sys")
	PY.signal = __import__ ("signal")
	PY.errno = __import__ ("errno")
	PY.subprocess = __import__ ("subprocess")
	PY.time = __import__ ("time")
	PY.stat = __import__ ("stat")
	PY.re = __import__ ("re")
	
	PY.path = PY.os.path
	
	if PY.sys.version_info[0] > 2 :
		PY.basestring = str
		PY.str = str
		PY.unicode = str
		PY.bytes = bytes
	else :
		PY.basestring = basestring
		PY.str = str
		PY.unicode = unicode
		PY.bytes = str
	
	Z.py = PY
	
	## --------------------------------------------------------------------------------
	
	def _inject (_function) :
		_name = _function.__name__
		if _name.startswith ("__Z__") :
			_name = _name[5:]
		else :
			assert False, ("[83cec849]  invalid inject name: `%s`" % _name)
		Z.__dict__[_name] = _function
		return _function
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__zspawn (_scriptlet, *_arguments, _wait = True, _panic = True, **_options) :
		_descriptor = Z._zexec_prepare (_scriptlet, _arguments, **_options)
		return Z.spawn_0 (_descriptor, _wait, _panic)
	
	@_inject
	def __Z__zexec (_scriptlet, *_arguments, **_options) :
		_descriptor = Z._zexec_prepare (_scriptlet, _arguments, **_options)
		return Z.exec_0 (_descriptor)
	
	@_inject
	def __Z__zcmd (_scriptlet, *_arguments, **_options) :
		return Z._zexec_prepare (_scriptlet, _arguments, **_options)
	
	@_inject
	def __Z___zexec_prepare (_scriptlet, _arguments, _env = None, _env_overrides = None, _path = None, _path_prepend = None, _chdir = None) :
		_executable = Z.executable
		if not _scriptlet.startswith ("::") :
			Z.panic (0xbd1641c7, "invalid scriptlet: `%s`", _scriptlet)
		_arguments_all = ["[z-run]", _scriptlet]
		_arguments_all.extend (_arguments)
		if _env is not None :
			_environment = { _name : _env[_name] for _name in _env }
		else :
			_environment = { _name : Z.environment[_name] for _name in Z.environment }
		if _env_overrides is not None :
			for _name in _env_overrides :
				_environment[_name] = _env_overrides[_name]
		if _path is not None :
			_environment["PATH"] = _path
		if _path_prepend is not None :
			_environment["PATH"] = Z.paths_prepend (_environment.get ("PATH"), *_path_prepend)
		return _executable, False, _arguments_all, _environment, _chdir
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__spawn (_scriptlet, *_arguments, _wait = True, _panic = True, **_options) :
		_descriptor = Z._exec_prepare (_scriptlet, _arguments, **_options)
		return Z.spawn_0 (_descriptor, _wait, _panic)
	
	@_inject
	def __Z__exec (_executable, *_arguments, **_options) :
		_descriptor = Z._exec_prepare (_executable, _arguments, **_options)
		return Z.exec_0 (_descriptor)
	
	@_inject
	def __Z__cmd (_scriptlet, *_arguments, **_options) :
		return Z._exec_prepare (_scriptlet, _arguments, **_options)
	
	@_inject
	def __Z___exec_prepare (_executable, _arguments, _env = None, _env_overrides = None, _path = None, _path_prepend = None, _chdir = None) :
		_arguments_all = [_executable]
		_arguments_all.extend (_arguments)
		if _env is not None :
			_environment = { _name : _env[_name] for _name in _env }
		else :
			_environment = { _name : Z.environment[_name] for _name in Z.environment }
		if _env_overrides is not None :
			for _name in _env_overrides :
				_environment[_name] = _env_overrides[_name]
		if _path is not None :
			_environment["PATH"] = _path
		if _path_prepend is not None :
			_environment["PATH"] = Z.paths_prepend (_environment.get ("PATH"), *_path_prepend)
		return _executable, True, _arguments_all, _environment, _chdir
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__spawn_0 (_descriptor, _wait = True, _panic = True) :
		# FIXME:  Handle lookup!
		_executable, _lookup, _arguments, _environment, _chdir = _descriptor
		_process = PY.subprocess.Popen (
				_arguments,
				executable = _executable,
				env = _environment,
				cwd = _chdir,
				stdin = None,
				stdout = None,
				stderr = None,
				close_fds = False,
				shell = False,
			)
		if _wait :
			_outcome = _process.wait ()
			if _panic and _outcome != 0 :
				Z.panic ((_panic, 0x7d3900c4), "spawn `%s` `%s` failed with status: %d", _arguments[0], _arguments[1:], _outcome)
		else :
			_outcome = _process.pid
		return _outcome
	
	@_inject
	def __Z__exec_0 (_descriptor) :
		_executable, _lookup, _arguments, _environment, _chdir = _descriptor
		if _chdir is not None :
			PY.os.chdir (_chdir)
		if _lookup :
			_delegate = PY.os.execvpe
		else :
			_delegate = PY.os.execve
		_delegate (_executable, _arguments, _environment)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__pipeline (_commands, _wait = True, _panic = True) :
		_count = len (_commands)
		if _count == 0 :
			Z.panic (0x1b1812d7, "pipeline empty")
		_pipes = []
		_pipes.append ((None, None))
		for _index in range (_count - 1) :
			_pipes.append (PY.os.pipe ())
		_pipes.append ((None, None))
		_processes = []
		for _index in range (_count) :
			# FIXME:  Handle lookup!
			_executable, _lookup, _arguments, _environment, _chdir = _commands[_index]
			_pipe_previous = _pipes[_index]
			_pipe_next = _pipes[_index + 1]
			_pipe_stdin = _pipe_previous[0]
			_pipe_stdout = _pipe_next[1]
			_process = PY.subprocess.Popen (
					_arguments,
					executable = _executable,
					env = _environment,
					cwd = _chdir,
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
	def __Z__log_error (_code, _message, *_arguments) :
		Z._log_write ('ee', _code, _message, _arguments)
	
	@_inject
	def __Z__log_warning (_code, _message, *_arguments) :
		if not Z.log_warning_enabled : return
		Z._log_write ('ww', _code, _message, _arguments)
	
	@_inject
	def __Z__log_notice (_code, _message, *_arguments) :
		if not Z.log_warning_enabled or not Z.log_notice_enabled : return
		Z._log_write ('ii', _code, _message, _arguments)
	
	@_inject
	def __Z__log_debug (_code, _message, *_arguments) :
		if not Z.log_warning_enabled or not Z.log_notice_enabled or not Z.log_debug_enabled : return
		Z._log_write ('dd', _code, _message, _arguments)
	
	@_inject
	def __Z__log_cut (_important = True) :
		if not Z.log_notice_enabled : return
		if _important :
			Z.stderr.write (("\n[z-run:%08d] [%s]  " % (Z.pid, "--")) + ("-" * 60) + "\n\n")
		else :
			Z.stderr.write (("[z-run:%08d] [%s]" % (Z.pid, "--")) + "\n")
	
	@_inject
	def __Z___log_write (_slug, _code, _message, _arguments) :
		Z.stderr.write (("[z-run:%08d] [%s] [%08x]  " % (Z.pid, _slug, _code)) + (_message % _arguments) + "\n")
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__exit (_status) :
		PY.sys.exit (_status)
	
	@_inject
	def __Z__panic (_code, _message, *_arguments) :
		if isinstance (_code, tuple) :
			_code_0 = 0xee4006b2
			for _code_0 in _code :
				if _code_0 is not False and _code_0 is not None :
					break
			_code = _code_0
		Z._log_write ('!!', _code, _message, _arguments)
		Z.exit (1)
	
	@_inject
	def __Z__sleep (_interval) :
		PY.time.sleep (_interval)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__expect_no_arguments () :
		return Z.expect_arguments (_exact = 0)
	
	@_inject
	def __Z__expect_arguments (_exact = None, _min = None, _max = None, _rest = None) :
		if _exact is not None :
			assert _min is None and _max is None, "[6ace19bf]"
			assert _rest is None, "[3465f981]"
			assert _exact >= 0, "[bfdc6527]"
		else :
			assert _min is not None or _max is not None, "[91787dc8]"
			assert _min is None or _min >= 0, "[19c3b0ad]"
			assert _max is None or _max >= 0, "[236d9471]"
			assert _min is None or _max is None or _min < _max, "[e43f6498]"
			assert _max is None or _max is None or _min < _max, "[0662ca21]"
			assert _rest is None or _rest is True or _rest is False, "[f4353518]"
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
			if _rest is None or _rest is True :
				return tuple (Z.arguments[:_min]) + (list (Z.arguments[_min:]),)
			else :
				return tuple (Z.arguments)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__enforce_regex (_value, _pattern, _code = None, _message = None) :
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
	def __Z__regex (_pattern) :
		return PY.re.compile (_pattern, PY.re.ASCII | PY.re.DOTALL)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__path (_path, _absolute = False, _canonical = False, _relative = None) :
		if not isinstance (_path, PY.basestring) and not isinstance (_path, PY.bytes) :
			_path = PY.path.join (*_path)
		_path = PY.path.normpath (_path)
		if _path.startswith ("//") :
			_path = "/" + _path.lstrip ("/")
		if _absolute :
			_path = PY.path.abspath (_path)
		if _canonical :
			_path = PY.path.realpath (_path)
		if _relative is not None :
			_path = PY.path.relpath (_relative)
		return _path
	
	@_inject
	def __Z__paths_append (_list, *_paths) :
		if _list is None :
			_list = []
		else :
			_list = [_path for _path in _list.split (":") if _path != ""]
		for _path in _paths :
			if ":" in _path :
				Z.panic (0x1908e707, "invalid path `%s` (contains `:`)", _path)
		_new = []
		_new.extend (_list)
		_new.extend (_paths)
		_new = ":".join (_new)
		return _new
	
	@_inject
	def __Z__paths_prepend (_list, *_paths) :
		if _list is None :
			_list = []
		else :
			_list = [_path for _path in _list.split (":") if _path != ""]
		for _path in _paths :
			if ":" in _path :
				Z.panic (0x1908e707, "invalid path `%s` (contains `:`)", _path)
		_new = []
		_new.extend (_paths)
		_new.extend (_list)
		_new = ":".join (_new)
		return _new
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__exists (_path, _follow = True, _panic = False) :
		_stat = Z.stat (_path, _follow)
		if _stat is None and _panic :
			Z.panic ((_panic, 0x383d3cc5), "file-system path does not exist `%s`", _path)
		return _stat is not None
	
	@_inject
	def __Z__not_exists (_path, _follow = True, _panic = False) :
		_stat = Z.stat (_path, _follow)
		if _stat is not None and _panic :
			Z.panic ((_panic, 0x9064abfc), "file-system path already exists `%s`", _path)
		return _stat is None
	
	@_inject
	def __Z__is_file (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __Z__is_file_empty (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode) and _stat.st_size == 0), _follow, _panic)
	
	@_inject
	def __Z__is_file_not_empty (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode) and _stat.st_size > 0), _follow, _panic)
	
	@_inject
	def __Z__is_folder (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISDIR (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __Z__is_file_or_folder (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode) or PY.stat.S_ISDIR (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __Z__is_symlink (_path, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISLNK (_stat.st_mode)), False, _panic)
	
	@_inject
	def __Z__is_pipe (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISFIFO (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __Z__is_socket (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISSOCK (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __Z__is_dev_block (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISBLK (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __Z__is_dev_char (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISCHR (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __Z__is_special (_path, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISFIFO (_stat.st_mode) or PY.stat.S_ISSOCK (_stat.st_mode) or PY.stat.S_ISBLK (_stat.st_mode) or PY.stat.S_ISCHR (_stat.st_mode)), _follow, _panic)
	
	@_inject
	def __Z___stat_check (_path, _check, _follow = True, _panic = False) :
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
	def __Z__stat (_path, _follow = True) :
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
	def __Z__mkdir (_path, _mode = None, _recurse = False) :
		if _mode is None : _mode = 0o777
		if _recurse :
			PY.os.makedirs (_path, _mode, True)
		else :
			PY.os.mkdir (_path, _mode)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__chdir (_path) :
		PY.os.chdir (_path)
	
	## --------------------------------------------------------------------------------
	
	class Z_environment :
		def __str__ (self) : return PY.os.environ.__repr__ ()
		def __repr__ (self) : return PY.os.environ.__repr__ ()
		def __getattr__ (self, _k) : return PY.os.environ.__getitem__ (_k)
		def __setattr__ (self, _k, _v) : return PY.os.environ.__setitem__ (_k, _v)
		def __len__ (self) : return PY.os.environ.__len__ ()
		def __getitem__ (self, _k) : return PY.os.environ.__getitem__ (_k)
		def __setitem__ (self, _k, _v) : return PY.os.environ.__setitem__ (_k, _v)
		def __delitem__ (self, _k) : return PY.os.environ.__delitem__ (_k)
		def __contains__ (self, _k) : return PY.os.environ.__contains__ (_k)
		def __iter__ (self) : return PY.os.environ.__iter__ ()
	
	class Z_environment_or_none :
		def __str__ (self) : return PY.os.environ.__repr__ ()
		def __repr__ (self) : return PY.os.environ.__repr__ ()
		def __getattr__ (self, _k) : return PY.os.environ.__getitem__ (_k) if PY.os.environ.__contains__ (_k) else None
		def __setattr__ (self, _k, _v) : return PY.os.environ.__setitem__ (_k, _v)
		def __len__ (self) : return PY.os.environ.__len__ ()
		def __getitem__ (self, _k) : return PY.os.environ.__getitem__ (_k) if PY.os.environ.__contains__ (_k) else None
		def __setitem__ (self, _k, _v) : return PY.os.environ.__setitem__ (_k, _v)
		def __delitem__ (self, _k) : return PY.os.environ.__delitem__ (_k)
		def __contains__ (self, _k) : return PY.os.environ.__contains__ (_k)
		def __iter__ (self) : return PY.os.environ.__iter__ ()
	
	## --------------------------------------------------------------------------------
	
	Z.pid = PY.os.getpid ()
	Z.arguments = tuple (PY.sys.argv[1:])
	Z.environment = Z_environment ()
	Z.environment_or_none = Z_environment_or_none ()
	
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
	
	(lambda Z : Z.py.sys.modules.__setitem__ ("Z", Z)) (__Z__create ())
	import Z
	
	Z.py.signal.signal (Z.py.signal.SIGINT, (lambda _0, _1 : Z.panic (0x6c751732, "scriptlet interrupted with SIGINT;  aborting!")))
	Z.py.signal.signal (Z.py.signal.SIGTERM, (lambda _0, _1 : Z.panic (0xb5067479, "scriptlet interrupted with SIGTERM;  aborting!")))
	Z.py.signal.signal (Z.py.signal.SIGQUIT, (lambda _0, _1 : Z.panic (0x921de146, "scriptlet interrupted with SIGQUIT;  aborting!")))
	Z.py.signal.signal (Z.py.signal.SIGHUP, (lambda _0, _1 : Z.panic (0xe2e4c7c5, "scriptlet interrupted with SIGHUP;  aborting!")))
	Z.py.signal.signal (Z.py.signal.SIGPIPE, (lambda _0, _1 : Z.panic (0xed8191f4, "scriptlet interrupted with SIGPIPE;  aborting!")))
	Z.py.signal.signal (Z.py.signal.SIGABRT, (lambda _0, _1 : Z.panic (0xd6af6d5b, "scriptlet interrupted with SIGABRT;  aborting!")))
	
	sys = Z.py.sys
	os = Z.py.os


################################################################################
################################################################################

