#!/dev/null


################################################################################
################################################################################


def __Z__create (*, Z = None, __import__ = __import__) :
	
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
	PY.fcntl = __import__ ("fcntl")
	PY.io = __import__ ("io")
	PY.re = __import__ ("re")
	PY.json = __import__ ("json")
	PY.fnmatch = __import__ ("fnmatch")
	PY.random = __import__ ("random") .SystemRandom ()
	PY.binascii = __import__ ("binascii")
	
	PY.path = PY.os.path
	
	if PY.sys.version_info[0] > 2 :
		PY.basestring = str
		PY.str = str
		PY.unicode = str
		PY.bytes = bytes
		PY.builtins = __import__ ("builtins")
	else :
		PY.basestring = basestring
		PY.str = str
		PY.unicode = unicode
		PY.bytes = str
		PY.builtins = __builtins__
	
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
	def __Z__zspawn (_scriptlet, *_arguments, _wait = True, _stdin_data = None, _stdout_data = None, _stderr_data = None, _fd_close = False, _panic = True, **_options) :
		_descriptor = Z._zexec_prepare (_scriptlet, _arguments, **_options)
		return Z.spawn_0 (_descriptor, _wait = _wait, _stdin_data = _stdin_data, _stdout_data = _stdout_data, _stderr_data = _stderr_data, _fd_close = _fd_close, _panic = _panic)
	
	@_inject
	def __Z__zspawn_capture (_scriptlet, *_arguments, **_options) :
		_output = Z.zspawn (_scriptlet, *_arguments, _wait = True, _stdin_data = False, _stdout_data = str, _panic = True)
		return Z._spawn_capture_output (_output, **_options)
	
	@_inject
	def __Z__zexec (_scriptlet, *_arguments, _fd_close = True, **_options) :
		_descriptor = Z._zexec_prepare (_scriptlet, _arguments, **_options)
		return Z.exec_0 (_descriptor, _fd_close = _fd_close)
	
	@_inject
	def __Z__zcmd (_scriptlet, *_arguments, **_options) :
		return Z._zexec_prepare (_scriptlet, _arguments, **_options)
	
	@_inject
	def __Z___zexec_prepare (_scriptlet, _arguments, **_options) :
		_executable = Z.executable
		if not _scriptlet.startswith ("::") :
			Z.panic (0xbd1641c7, "invalid scriptlet: `%s`", _scriptlet)
		_arguments_all = ["[z-run]", _scriptlet]
		_arguments_all.extend (_arguments)
		return Z._exec_prepare_0 (_executable, False, _arguments_all, **_options)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__spawn (_executable, *_arguments, _wait = True, _stdin_data = None, _stdout_data = None, _stderr_data = None, _fd_close = False, _panic = True, **_options) :
		_descriptor = Z._exec_prepare (_executable, _arguments, **_options)
		return Z.spawn_0 (_descriptor, _wait = _wait, _stdin_data = _stdin_data, _stdout_data = _stdout_data, _stderr_data = _stderr_data, _fd_close = _fd_close, _panic = _panic)
	
	@_inject
	def __Z__spawn_capture (_executable, *_arguments, **_options) :
		_output = Z.spawn (_executable, *_arguments, _wait = True, _stdin_data = False, _stdout_data = str, _panic = True)
		return Z._spawn_capture_output (_output, **_options)
	
	@_inject
	def __Z__exec (_executable, *_arguments, _fd_close = True, **_options) :
		_descriptor = Z._exec_prepare (_executable, _arguments, **_options)
		return Z.exec_0 (_descriptor, _fd_close = _fd_close)
	
	@_inject
	def __Z__cmd (_executable, *_arguments, **_options) :
		return Z._exec_prepare (_executable, _arguments, **_options)
	
	@_inject
	def __Z___exec_prepare (_executable, _arguments, **_options) :
		_arguments_all = [_executable]
		_arguments_all.extend (_arguments)
		return Z._exec_prepare_0 (_executable, True, _arguments_all, **_options)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__spawn_0 (_descriptor, *, _wait = True, _stdin_data = None, _stdout_data = None, _stderr_data = None, _fd_close = False, _panic = True) :
		# FIXME:  Handle lookup!
		_executable, _lookup, _arguments, _environment, _chdir, _files = _descriptor
		if _files is not None :
			_stdin, _stdout, _stderr = _files
			_stdin = Z._fd (_stdin)
			_stdout = Z._fd (_stdout)
			_stderr = Z._fd (_stderr)
		else :
			_stdin = None
			_stdout = None
			_stderr = None
		_should_communicate = False
		if _stdin_data is not None :
			if _stdin is not None :
				Z.panic (0x5fcf5035, "stdin unexpected")
			if _stdin_data is False :
				_stdin = PY.subprocess.DEVNULL
				_stdin_data = None
			elif isinstance (_stdin_data, PY.basestring) or isinstance (_stdin_data, PY.bytes) :
				if isinstance (_stdin_data, PY.basestring) :
					_stdin_data = _stdin_data.encode ("utf-8")
				_stdin = PY.subprocess.PIPE
				_should_communicate = True
			else :
				Z.panic (0x5566ac86, "stdin data invalid")
		if _stdout_data is not None :
			if _stdout is not None :
				Z.panic (0xab0c2481, "stdout unexpected")
			if _stdout_data is False :
				_stdout = PY.subprocess.DEVNULL
				_stdout_data = None
			elif _stdout_data is True or _stdout_data is PY.str or _stdout_data is PY.unicode or _stdout_data is PY.bytes :
				_stdout = PY.subprocess.PIPE
				_should_communicate = True
			else :
				Z.panic (0xff42dc1e, "stdout data invalid")
		if _stderr_data is not None :
			if _stderr is not None :
				Z.panic (0xa42d1da3, "stderr unexpected")
			if _stderr_data is False :
				_stderr = PY.subprocess.DEVNULL
				_stderr_data = None
			elif _stderr_data is True or _stderr_data is PY.str or _stderr_data is PY.unicode or _stderr_data is PY.bytes :
				_stderr = PY.subprocess.PIPE
				_should_communicate = True
			else :
				Z.panic (0x2d0ed281, "stderr data invalid")
		if _should_communicate and not _wait :
			Z.panic (0xe20e7d58, "data arguments require waiting")
		_process = PY.subprocess.Popen (
				_arguments,
				executable = _executable,
				env = _environment,
				cwd = _chdir,
				stdin = _stdin,
				stdout = _stdout,
				stderr = _stderr,
				close_fds = False,
				shell = False,
			)
		if _fd_close :
			if _stdin is not None :
				PY.os.close (_stdin)
			if _stdout is not None :
				PY.os.close (_stdout)
			if _stderr is not None :
				PY.os.close (_stderr)
		if _wait :
			if _should_communicate :
				_stdout_data_0, _stderr_data_0 = _process.communicate (_stdin_data)
				if _stdout_data is True or _stdout_data is PY.unicode :
					_stdout_data_0 = _stdout_data_0.decode ("utf-8")
				elif _stdout_data is PY.str :
					_stdout_data_0 = _stdout_data_0.decode ("ascii")
				elif _stdout_data is PY.bytes :
					pass
				elif _stdout_data is not None :
					Z.panic (0x70227d93, "invalid state")
				if _stderr_data is True or _stderr_data is PY.unicode :
					_stderr_data_0 = _stderr_data_0.decode ("utf-8")
				elif _stderr_data is PY.str :
					_stderr_data_0 = _stderr_data_0.decode ("ascii")
				elif _stderr_data is PY.bytes :
					pass
				elif _stderr_data is not None :
					Z.panic (0x07342d09, "invalid state")
				_stdout_data = _stdout_data_0
				_stderr_data = _stderr_data_0
			else :
				_process.wait ()
			_outcome = _process.returncode
			if _panic and _outcome != 0 :
				Z.panic ((_panic, 0x7d3900c4), "spawn `%s` `%s` failed with status: %d", _arguments[0], _arguments[1:], _outcome)
			if _should_communicate :
				if _stderr_data is not None :
					if _panic :
						_outcome = (_stdout_data, _stderr_data)
					else :
						_outcome = (_outcome, _stdout_data, _stderr_data)
				else :
					if _panic :
						_outcome = _stdout_data
					else :
						_outcome = (_outcome, _stdout_data)
		else :
			_outcome = _process.pid
		return _outcome
	
	@_inject
	def __Z__exec_0 (_descriptor, *, _fd_close = True) :
		_executable, _lookup, _arguments, _environment, _chdir, _files = _descriptor
		if _chdir is not None :
			PY.os.chdir (_chdir)
		if _lookup :
			_delegate = PY.os.execvpe
		else :
			_delegate = PY.os.execve
		if _files is not None :
			_stdin, _stdout, _stderr = _files
			_stdin = Z._fd (_stdin)
			_stdout = Z._fd (_stdout)
			_stderr = Z._fd (_stderr)
		else :
			_stdin = None
			_stdout = None
			_stderr = None
		if _stdin is not None :
			PY.os.dup2 (_stdin, 0, True)
			if _fd_close and _stdin != 0 :
				PY.os.close (_stdin)
		if _stdout is not None :
			PY.os.dup2 (_stdout, 1, True)
			if _fd_close and _stdout != 1 :
				PY.os.close (_stdout)
		if _stderr is not None :
			PY.os.dup2 (_stderr, 2, True)
			if _fd_close and _stderr != 2 :
				PY.os.close (_stderr)
		_delegate (_executable, _arguments, _environment)
	
	@_inject
	def __Z___exec_prepare_0 (_executable, _lookup, _arguments, *, _env = None, _env_overrides = None, _path = None, _path_prepend = None, _chdir = None, _stdin = None, _stdout = None, _stderr = None) :
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
		if _stdin is not None or _stdout is not None or _stderr is not None :
			_files = (_stdin, _stdout, _stderr)
		else :
			_files = None
		return _executable, _lookup, _arguments, _environment, _chdir, _files
	
	@_inject
	def __Z__process_wait (_pid, *, _panic = None) :
		_pid_0 = Z._pid (_pid)
		_pid, _outcome = PY.os.waitpid (_pid_0, 0)
		if PY.os.WIFEXITED (_outcome) :
			_outcome = PY.os.WEXITSTATUS (_outcome)
		elif PY.os.WIFSIGNALED (_outcome) :
			_outcome = 0 - PY.os.WTERMSIG (_outcome)
		else :
			Z.panic (0xb4179b04, "waiting `%s` failed with unknown outcome: %r", _pid, _outcome)
		if _panic and _outcome != 0 :
			Z.panic ((_panic, 0xc0e8ec5d), "waiting `%s` failed with status: %d", _pid, _outcome)
		if _pid == _pid_0 :
			return _outcome
		else :
			return _pid, _outcome
	
	@_inject
	def __Z__process_signal (_pid, _signal, *, _wait = False, _panic = None) :
		_pid = Z._pid (_pid)
		PY.os.kill (_pid, _signal)
		if _wait :
			return Z.process_wait (_pid, _panic = _panic)
	
	@_inject
	def __Z__process_terminate (_pid, *, _wait = False, _panic = None) :
		return Z.process_signal (_pid, PY.signal.SIGTERM)
	
	@_inject
	def __Z__process_kill (_pid, *, _wait = False, _panic = None) :
		return Z.process_signal (_pid, PY.signal.SIGKILL)
	
	@_inject
	def __Z___pid (_process) :
		if _process is None :
			_pid = None
		elif isinstance (_process, PY.builtins.int) :
			if _process >= 1 :
				_pid = _process
			else :
				Z.panic (0x1fe23810, "invalid process id (negative)")
		elif isinstance (_process, PY.subprocess.Popen) :
			_pid = _process.pid
		else :
			Z.panic (0x3960636f, "invalid process id (unknown): %r", _pid)
		return _pid
	
	@_inject
	def __Z___spawn_capture_output (_output, *, _line = False, _lines = False, _separator = None, _json = False) :
		if _line or _lines :
			if _separator is None :
				_separator = "\n"
			else :
				assert _separator != "", "[357b50e2]"
		if _line :
			assert not _lines, "[cefc2173]"
			if _output == "" :
				_output = None
			else :
				if _output[: 0 - len (_separator)] == _separator :
					_output = _output[: 0 - len (_separator)]
				_output = _output.split (_separator)
				if len (_output) == 1 or (len (_output) == 2 and _output[1] == "") :
					_output = _output[0]
				else :
					Z.panic (0x5a1b7a79, "output is made of multiple lines")
			if _json :
				_output = PY.json.loads (_output)
		elif _lines :
			assert not _line, "[ddd9a8bc]"
			if _output == "" :
				_output = None
			else :
				if _output[: 0 - len (_separator)] == _separator :
					_output = _output[: 0 - len (_separator)]
				_output = _output.split (_separator)
			if _json :
				_output = [PY.json.loads (_output) for _output in _output]
		elif _json :
			_output = PY.json.loads (_output)
		else :
			pass
		return _output
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__pipeline (_commands, *, _wait = True, _fd_close = False, _panic = True) :
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
			_executable, _lookup, _arguments, _environment, _chdir, _files = _commands[_index]
			if _files is not None :
				_stdin, _stdout, _stderr = _files
				_stdin = Z._fd (_stdin)
				_stdout = Z._fd (_stdout)
				_stderr = Z._fd (_stderr)
			else :
				_stdin = None
				_stdout = None
				_stderr = None
			_pipe_previous = _pipes[_index]
			_pipe_next = _pipes[_index + 1]
			_pipe_stdin = _pipe_previous[0]
			_pipe_stdout = _pipe_next[1]
			if _stdin is not None and _pipe_stdin is not None :
				Z.panic (0xd2333746, "stdin unexpected")
			if _stdout is not None and _pipe_stdout is not None :
				Z.panic (0xf9ec0dff, "stdout unexpected")
			_process = PY.subprocess.Popen (
					_arguments,
					executable = _executable,
					env = _environment,
					cwd = _chdir,
					stdin = _stdin if _stdin is not None else _pipe_stdin,
					stdout = _stdout if _stdout is not None else _pipe_stdout,
					stderr = _stderr,
					close_fds = False,
					shell = False,
				)
			if _wait :
				_processes.append ((_index, _process, _arguments))
			else :
				_processes.append (_process.pid)
			if _fd_close :
				if _stdin is not None :
					PY.os.close (_stdin)
				if _stdout is not None :
					PY.os.close (_stdout)
				if _stderr is not None :
					PY.os.close (_stderr)
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
	def __Z__log_cut (*, _important = True) :
		if not Z.log_notice_enabled : return
		if _important :
			Z.stderr.write (("\n[z-run:%08d] [%s]  " % (Z.pid, "--")) + ("-" * 60) + "\n\n")
			Z.stderr.flush ()
		else :
			Z.stderr.write (("[z-run:%08d] [%s]" % (Z.pid, "--")) + "\n")
			Z.stderr.flush ()
	
	@_inject
	def __Z___log_write (_slug, _code, _message, _arguments) :
		Z.stderr.write (("[z-run:%08d] [%s] [%08x]  " % (Z.pid, _slug, _code)) + (_message % _arguments) + "\n")
		Z.stderr.flush ()
	
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
	def __Z__expect_arguments (*, _exact = None, _min = None, _max = None, _rest = None) :
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
	def __Z__enforce (_condition, *, _code = None, _message = None) :
		if _message is None : _message = "enforcement failed"
		if not isinstance (_condition, PY.builtins.bool) :
			if _code is None : _code = 0x97ee7cf1
			Z.panic (_code, _message)
		if not _condition :
			if _code is None : _code = 0x0e36cc55
			Z.panic (_code, _message)
		return _condition
	
	@_inject
	def __Z__enforce_regex (_value, _pattern, *, _code = None, _message = None) :
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
	def __Z__path (_path, *, _absolute = False, _canonical = False, _relative = None) :
		if not isinstance (_path, PY.basestring) and not isinstance (_path, PY.bytes) :
			_path = PY.path.join (*_path)
		_path = Z.path_normalize (_path)
		if _absolute :
			_path = Z.path_absolute (_path)
		if _canonical :
			_path = Z.path_canonical (_path)
		if _relative is not None :
			_path = Z.path_relative (_path, _relative)
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
	
	@_inject
	def __Z__path_normalize (_path) :
		_path = PY.path.normpath (_path)
		if _path.startswith ("//") :
			_path = "/" + _path.lstrip ("/")
		return _path
	
	@_inject
	def __Z__path_dirname (_path) :
		return PY.path.dirname (_path)
	
	@_inject
	def __Z__path_basename (_path) :
		return PY.path.basename (_path)
	
	@_inject
	def __Z__path_canonical (_path) :
		return PY.path.realpath (_path)
	
	@_inject
	def __Z__path_absolute (_path) :
		return PY.path.abspath (_path)
	
	@_inject
	def __Z__path_relative (_path) :
		return PY.path.relpath (_path, _relative)
	
	@_inject
	def __Z__path_split (_path) :
		_components = []
		_dirname = Z.path_normalize (_path)
		while _dirname != "" :
			_dirname, _basename = PY.path.split (_dirname)
			if _basename != "" :
				_components.append (_basename)
			if _dirname == "/" :
				_components.append (_dirname)
				_dirname = ""
		_components = _components.reverse ()
		_components = tuple (_components)
		return _components
	
	@_inject
	def __Z__path_extension (_path) :
		_path, _extension = PY.path.splitext (_path)
		if _extension != "" :
			_extension = _extension[1:]
		else :
			_extension = None
		return _extension
	
	@_inject
	def __Z__path_without_extension (_path) :
		_path, _extension = PY.path.splitext (_path)
		return _path
	
	@_inject
	def __Z__path_matches (_path, _pattern) :
		return PY.fnmatch.fnmatch (_path, _pattern)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__exists (_path, *, _follow = True, _panic = False) :
		_stat = Z.stat (_path, _follow = _follow)
		if _stat is None and _panic :
			Z.panic ((_panic, 0x383d3cc5), "file-system path does not exist `%s`", _path)
		return _stat is not None
	
	@_inject
	def __Z__not_exists (_path, *, _follow = True, _panic = False) :
		_stat = Z.stat (_path, _follow = _follow)
		if _stat is not None and _panic :
			Z.panic ((_panic, 0x9064abfc), "file-system path already exists `%s`", _path)
		return _stat is None
	
	@_inject
	def __Z__is_file (_path, *, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode)), _follow = _follow, _panic = _panic)
	
	@_inject
	def __Z__is_file_empty (_path, *, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode) and _stat.st_size == 0), _follow = _follow, _panic = _panic)
	
	@_inject
	def __Z__is_file_not_empty (_path, *, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode) and _stat.st_size > 0), _follow = _follow, _panic = _panic)
	
	@_inject
	def __Z__is_folder (_path, *, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISDIR (_stat.st_mode)), _follow = _follow, _panic = _panic)
	
	@_inject
	def __Z__is_file_or_folder (_path, *, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode) or PY.stat.S_ISDIR (_stat.st_mode)), _follow = _follow, _panic = _panic)
	
	@_inject
	def __Z__is_symlink (_path, *, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISLNK (_stat.st_mode)), _follow = False, _panic = _panic)
	
	@_inject
	def __Z__is_pipe (_path, *, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISFIFO (_stat.st_mode)), _follow = _follow, _panic = _panic)
	
	@_inject
	def __Z__is_socket (_path, *, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISSOCK (_stat.st_mode)), _follow = _follow, _panic = _panic)
	
	@_inject
	def __Z__is_dev_block (_path, *, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISBLK (_stat.st_mode)), _follow = _follow, _panic = _panic)
	
	@_inject
	def __Z__is_dev_char (_path, *, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISCHR (_stat.st_mode)), _follow = _follow, _panic = _panic)
	
	@_inject
	def __Z__is_special (_path, *, _follow = True, _panic = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISFIFO (_stat.st_mode) or PY.stat.S_ISSOCK (_stat.st_mode) or PY.stat.S_ISBLK (_stat.st_mode) or PY.stat.S_ISCHR (_stat.st_mode)), _follow = _follow, _panic = _panic)
	
	@_inject
	def __Z___stat_check (_path, _check, *, _follow = True, _panic = False) :
		_stat = Z.stat (_path, _follow = _follow)
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
	def __Z__stat (_path, *, _follow = True) :
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
	def __Z__mkdir (_path, *, _mode = None, _recurse = False) :
		if _mode is None : _mode = 0o777
		if _recurse :
			PY.os.makedirs (_path, _mode, True)
		else :
			PY.os.mkdir (_path, _mode)
	
	@_inject
	def __Z__touch (_path, *, _mode = None, _create = True, _exclusive = None) :
		if _mode is None : _mode = 0o666
		_flags = PY.os.O_WRONLY | PY.os.O_NOCTTY
		if _create :
			_flags |= PY.os.O_CREAT
			if _exclusive is None :
				_exclusive = False
		else :
			assert not _exclusive, "[eae26c74]"
		if _exclusive :
			assert _create, "[89a753b5]"
			_flags |= PY.os.O_EXCL
		_file = PY.os.open (_path, _flags, _mode)
		PY.os.utime (_file, None)
		PY.os.close (_file)
	
	@_inject
	def __Z__rename (_source, _target) :
		PY.os.rename (_source, _target)
	
	@_inject
	def __Z__symlink (_source, _target) :
		PY.os.symlink (_source, _target)
	
	@_inject
	def __Z__file_read (_path, *, _data = None, _json = None, _panic = True) :
		_fd = Z.fd_open_for_read (_path, _panic = _panic)
		if _fd is None :
			return None
		_buffers = []
		while True :
			_buffer = PY.os.read (_fd, 1024 * 1024)
			if len (_buffer) == 0 :
				break
			_buffers.append (_buffer)
		PY.os.close (_fd)
		if _json :
			assert _data is None, "[be2949cc]"
		if _data is None or _data is PY.unicode :
			_buffers = [_buffer.decode ("utf-8") for _buffer in _buffers]
			_buffers = u"".join (_buffers)
		elif _data is PY.str :
			_buffers = [_buffer.decode ("ascii") for _buffer in _buffers]
			_buffers = "".join (_buffers)
		elif _data is PY.bytes :
			_buffers = "".join (_buffers)
		else :
			Z.panic (0x764fba9e, "invalid data type")
		if _json :
			_data = PY.json.loads (_buffers)
		else :
			_data = _buffers
		return _data
	
	@_inject
	def __Z__file_write (_path, _data, *, _json = None, _mode = None, _create = True, _exclusive = None, _append = False, _truncate = False, _panic = True) :
		_fd = Z.fd_open_for_write (_path, _create = _create, _exclusive = _exclusive, _append = _append, _truncate = _truncate, _panic = _panic)
		if _fd is None :
			return False
		if _json :
			_data = PY.json.dumps (_data)
		if _data is None :
			_buffer = b""
		elif isinstance (_data, PY.unicode) :
			_buffer = _data.encode ("utf-8")
		elif isinstance (_data, PY.str) :
			_buffer = _data.encode ("ascii")
		elif isinstance (_data, PY.bytes) :
			_buffer = _data
		else :
			Z.panic (0x5fbb8542, "invalid data type")
		while len (_buffer) > 0 :
			_offset = PY.os.write (_fd, _buffer)
			_buffer = _buffer[_offset:]
		PY.os.close (_fd)
		return True
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__fd_open_for_read (_path, *, _close_on_exec = True, _panic = True) :
		_flags = PY.os.O_RDONLY | PY.os.O_NOCTTY
		if _close_on_exec :
			_flags |= PY.os.O_CLOEXEC
		try :
			_file = PY.os.open (_path, _flags)
		except OSError as _error :
			if _error.errno == PY.errno.ENOENT :
				if _panic :
					Z.panic ((_panic, 0x9ca19f87), "open `%s` failed:  does not exist", _path)
			else :
				raise
		return _file
	
	@_inject
	def __Z__fd_open_for_write (_path, *, _mode = None, _create = True, _exclusive = None, _read = False, _append = False, _truncate = False, _close_on_exec = True, _panic = True) :
		if _mode is None : _mode = 0o666
		_flags = PY.os.O_WRONLY | PY.os.O_NOCTTY
		if _create :
			_flags |= PY.os.O_CREAT
			if _exclusive is None :
				_exclusive = True
		else :
			assert not _exclusive, "[0617f007]"
		if _exclusive :
			assert _create, "[45af3e8d]"
			_flags |= PY.os.O_EXCL
		if _read :
			_flags |= PY.os.O_RDWR
		if _append :
			_flags |= PY.os.O_APPEND
		if _truncate :
			_flags |= PY.os.O_TRUNC
		if _close_on_exec :
			_flags |= PY.os.O_CLOEXEC
		try :
			_file = PY.os.open (_path, _flags, _mode)
		except OSError as _error :
			if _error.errno == PY.errno.ENOENT :
				if _panic :
					Z.panic ((_panic, 0x008389df), "open `%s` failed:  does not exist", _path)
			elif _error.errno == PY.errno.EEXIST :
				if _panic :
					Z.panic ((_panic, 0x37020d90), "open `%s` failed:  already exists", _path)
			else :
				raise
		return _file
	
	@_inject
	def __Z__fd_open_null (*, _close_on_exec = True) :
		_flags = PY.os.O_RDWR | PY.os.O_NOCTTY
		if _close_on_exec :
			_flags |= PY.os.O_CLOEXEC
		_file = PY.os.open (PY.os.devnull, _flags)
		return _file
	
	@_inject
	def __Z__fd_open_pipes (*, _non_blocking = False, _close_on_exec = True) :
		_flags = 0
		if _non_blocking :
			_flags |= PY.os.O_NONBLOCK
		if _close_on_exec :
			_flags |= PY.os.O_CLOEXEC
		_input, _output = PY.os.pipe2 (_flags)
		return _input, _output
	
	@_inject
	def __Z__fd_clone (_file, *, _close_on_exec = True) :
		_file = Z._fd (_file)
		if _close_on_exec :
			_command = PY.fcntl.F_DUPFD_CLOEXEC
		else :
			_command = PY.fcntl.F_DUPFD
		_fd = PY.fcntl.fcntl (_file, _command, 3)
		return _fd
	
	@_inject
	def __Z__fd_flush (_file) :
		_file = Z._fd (_file)
		PY.os.fsync (_file)
	
	@_inject
	def __Z__fd_close (_file) :
		_file = Z._fd (_file)
		PY.os.close (_file)
	
	@_inject
	def __Z__is_fd (_file, *, _panic = None) :
		if _file is None :
			if _panic :
				Z.panic ((_panic, 0x3c3ff99e), "invalid file descriptor (none)")
			return False
		elif isinstance (_file, PY.builtins.int) :
			if _file >= 0 :
				return True
			else :
				Z.panic (0xf42752f0, "invalid file descriptor (negative)")
		elif isinstance (_file, PY.io.IOBase) :
			try :
				_file.fileno ()
				return True
			except OSError as _error :
				Z.panic (0x7f97a33f, "invalid file descriptor (not supported): %r  //  %s", _error)
		else :
			if _panic :
				Z.panic ((_panic, 0xe0d78a09), "invalid file descriptor (unknown): %r", _file)
			return False
	
	@_inject
	def __Z___fd (_file) :
		if _file is None :
			_fd = None
		elif isinstance (_file, PY.builtins.int) :
			if _file >= 0 :
				_fd = _file
			else :
				Z.panic (0x4327e646, "invalid file descriptor (negative)")
		elif isinstance (_file, PY.io.IOBase) :
			try :
				_fd = _file.fileno ()
			except OSError as _error :
				Z.panic (0x32026a93, "invalid file descriptor (not supported): %r  //  %s", _error)
		else :
			Z.panic (0xe0d78a09, "invalid file descriptor (unknown): %r", _file)
		return _fd
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__chdir (_path) :
		PY.os.chdir (_path)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__random_bytes (_bytes) :
		if _bytes == 0 :
			return b""
		return PY.os.urandom (_bytes)
	
	@_inject
	def __Z__random_token (_bytes) :
		if _bytes == 0 :
			return ""
		_data = Z.random_bytes (_bytes)
		_data = PY.binascii.b2a_hex (_data)
		_data = _data.decode ("ascii")
		return _data
	
	@_inject
	def __Z__random_integer (_bits) :
		if _bits == 0 :
			return 0
		return PY.random.getrandbits (_bits)
	
	@_inject
	def __Z__random_range (_start, _stop, _step = 1) :
		return PY.random.randrange (_start, _stop, _step)
	
	@_inject
	def __Z__random_float (_a, _b) :
		return PY.random.uniform (_a, _b)
	
	@_inject
	def __Z__random_select (_sequence) :
		return PY.random.choice (_sequence)
	
	@_inject
	def __Z__random_sample (_sequence, _count, *, _repeats = False) :
		if not _repeats :
			return PY.random.sample (_sequence, _count)
		else :
			return PY.random.choices (_sequence, k = _count)
	
	@_inject
	def __Z__random_shuffle (_sequence) :
		_sequence = list (_sequence)
		PY.random.shuffle (_sequence)
		return _sequence
	
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

