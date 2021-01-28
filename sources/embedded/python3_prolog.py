#!/dev/null


################################################################################
################################################################################


def __Z__create (*, Z = None, __import__ = __import__) :
	
	## --------------------------------------------------------------------------------
	
	PY = __import__ ("types") .ModuleType ("PY")
	PY.sys = __import__ ("sys")
	
	if PY.sys.version_info[0] != 3 or PY.sys.version_info[1] < 6 :
		PY.sys.stderr.write ("[z-run] [!!] [68c3c553]  requires Python3.6+;  aborting!\n")
		PY.sys.stderr.flush ()
		PY.sys.exit (1)
	
	## --------------------------------------------------------------------------------
	
	PY.binascii = __import__ ("binascii")
	PY.builtins = __import__ ("builtins")
	PY.errno = __import__ ("errno")
	PY.fcntl = __import__ ("fcntl")
	PY.fnmatch = __import__ ("fnmatch")
	PY.hashlib = __import__ ("hashlib")
	PY.io = __import__ ("io")
	PY.json = __import__ ("json")
	PY.os = __import__ ("os")
	PY.random = __import__ ("random") .SystemRandom ()
	PY.re = __import__ ("re")
	PY.signal = __import__ ("signal")
	PY.stat = __import__ ("stat")
	PY.subprocess = __import__ ("subprocess")
	PY.time = __import__ ("time")
	PY.traceback = __import__ ("traceback")
	PY.types = __import__ ("types")
	
	PY.path = PY.os.path
	
	PY.bytes = PY.builtins.bytes
	PY.str = PY.builtins.str
	
	PY.int = PY.builtins.int
	PY.float = PY.builtins.float
	
	PY.tuple = PY.builtins.tuple
	PY.list = PY.builtins.list
	PY.dict = PY.builtins.dict
	PY.range = PY.builtins.range
	PY.len = PY.builtins.len
	PY.sorted = PY.builtins.sorted
	PY.reversed = PY.builtins.reversed
	
	PY.isinstance = PY.builtins.isinstance
	PY.OSError = PY.builtins.OSError
	PY.SystemExit = PY.builtins.SystemExit
	
	## --------------------------------------------------------------------------------
	
	if Z is None :
		Z = __import__ ("types") .ModuleType ("Z")
	
	Z.py = PY
	
	## --------------------------------------------------------------------------------
	
	def _inject (_function) :
		_name = _function.__name__
		if _name.startswith ("__Z__") :
			_name = _name[5:]
		else :
			assert False, ("[83cec849]  invalid inject name: `%s`" % _name)
		def __Z__wrapper (*_arguments_list, **_arguments_map) :
			_arguments_map = {("_" + _name if _name[0] != "_" else _name) : _value for _name, _value in _arguments_map.items ()}
			try :
				return _function (*_arguments_list, **_arguments_map)
			except PY.SystemExit :
				raise
			except :
				_error = PY.sys.exc_info ()
				_traceback_error = PY.traceback.extract_tb (_error[2])
				_traceback_caller = PY.traceback.extract_stack ()
				_error = _error[1]
				Z._panic_with_traceback (0x63468d09, _error, _traceback_error, _traceback_caller)
		_function.__name__ = "Z." + _name
		__Z__wrapper.__name__ = "Z." + _name
		Z.__dict__[_name] = __Z__wrapper
		return None
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__zspawn (_scriptlet, *_arguments, _wait = True, _stdin_data = None, _stdout_data = None, _stderr_data = None, _fd_close = False, _enforce = True, **_options) :
		_descriptor = Z._zexec_prepare (_scriptlet, _arguments, **_options)
		return Z.spawn_0 (_descriptor, _wait = _wait, _stdin_data = _stdin_data, _stdout_data = _stdout_data, _stderr_data = _stderr_data, _fd_close = _fd_close, _enforce = _enforce)
	
	@_inject
	def __Z__zspawn_capture (_scriptlet, *_arguments, _stdin_data = False, _fd_close = False, _env = None, _env_overrides = None, _path = None, _path_prepend = None, _chdir = None, _sudo = None, _sudo_user = None, **_options) :
		_output = Z.zspawn (_scriptlet, *_arguments, _wait = True, _stdin_data = _stdin_data, _stdout_data = PY.str, _enforce = True, _fd_close = _fd_close, _env = _env, _env_overrides = _env_overrides, _path = _path, _path_prepend = _path_prepend, _chdir = _chdir, _sudo = _sudo, _sudo_user = _sudo_user)
		return Z._file_read_process (_output, **_options)
	
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
	def __Z__spawn (_executable, *_arguments, _wait = True, _stdin_data = None, _stdout_data = None, _stderr_data = None, _fd_close = False, _enforce = True, **_options) :
		_descriptor = Z._exec_prepare (_executable, _arguments, **_options)
		return Z.spawn_0 (_descriptor, _wait = _wait, _stdin_data = _stdin_data, _stdout_data = _stdout_data, _stderr_data = _stderr_data, _fd_close = _fd_close, _enforce = _enforce)
	
	@_inject
	def __Z__spawn_capture (_executable, *_arguments, _stdin_data = False, _fd_close = False, _env = None, _env_overrides = None, _path = None, _path_prepend = None, _chdir = None, _sudo = None, _sudo_user = None, **_options) :
		_output = Z.spawn (_executable, *_arguments, _wait = True, _stdin_data = _stdin_data, _stdout_data = PY.str, _enforce = True, _fd_close = _fd_close, _env = _env, _env_overrides = _env_overrides, _path = _path, _path_prepend = _path_prepend, _chdir = _chdir, _sudo = _sudo, _sudo_user = _sudo_user)
		return Z._file_read_process (_output, **_options)
	
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
	def __Z__spawn_0 (_descriptor, *, _wait = True, _stdin_data = None, _stdout_data = None, _stderr_data = None, _fd_close = False, _enforce = True) :
		# FIXME:  Handle lookup!
		_executable, _lookup, _arguments, _environment, _chdir, _files = _descriptor
		if _files is not None :
			_stdin, _stdout, _stderr = _files
			_stdin = Z._fd (_stdin) if _stdin is not False else PY.subprocess.DEVNULL
			_stdout = Z._fd (_stdout) if _stdout is not False else PY.subprocess.DEVNULL
			_stderr = Z._fd (_stderr) if _stderr is not False else PY.subprocess.DEVNULL
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
			elif PY.isinstance (_stdin_data, PY.str) or PY.isinstance (_stdin_data, PY.bytes) :
				if PY.isinstance (_stdin_data, PY.str) :
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
			elif _stdout_data is True or _stdout_data is PY.str or _stdout_data is PY.bytes :
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
			elif _stderr_data is True or _stderr_data is PY.str or _stderr_data is PY.bytes :
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
				if _stdout_data is True or _stdout_data is PY.str :
					_stdout_data_0 = _stdout_data_0.decode ("utf-8")
				elif _stdout_data is PY.bytes :
					pass
				elif _stdout_data is not None :
					Z.panic (0x70227d93, "invalid state")
				if _stderr_data is True or _stderr_data is PY.str :
					_stderr_data_0 = _stderr_data_0.decode ("utf-8")
				elif _stderr_data is PY.bytes :
					pass
				elif _stderr_data is not None :
					Z.panic (0x07342d09, "invalid state")
				_stdout_data = _stdout_data_0
				_stderr_data = _stderr_data_0
			else :
				_process.wait ()
			_outcome = _process.returncode
			if _enforce and _outcome != 0 :
				Z.panic ((_enforce, 0x7d3900c4), "spawn `%s` `%s` failed with status: %d", _arguments[0], _arguments[1:], _outcome)
			if _should_communicate :
				if _stderr_data is not None :
					if _enforce :
						_outcome = (_stdout_data, _stderr_data)
					else :
						_outcome = (_outcome, _stdout_data, _stderr_data)
				else :
					if _enforce :
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
		_dev_null = None
		if _files is not None :
			_stdin, _stdout, _stderr = _files
			if _stdin is False or _stdout is False or _stderr is False :
				_dev_null = Z.fd_open_null ()
			_stdin = Z._fd (_stdin) if _stdin is not False else _dev_null
			_stdout = Z._fd (_stdout) if _stdout is not False else _dev_null
			_stderr = Z._fd (_stderr) if _stderr is not False else _dev_null
		else :
			_stdin = None
			_stdout = None
			_stderr = None
		if _stdin is not None :
			PY.os.dup2 (_stdin, 0, True)
			if _fd_close and _stdin != 0 and _stdin != _dev_null :
				PY.os.close (_stdin)
		if _stdout is not None :
			PY.os.dup2 (_stdout, 1, True)
			if _fd_close and _stdout != 1 and _stdout != _dev_null :
				PY.os.close (_stdout)
		if _stderr is not None :
			PY.os.dup2 (_stderr, 2, True)
			if _fd_close and _stderr != 2 and _stderr != _dev_null :
				PY.os.close (_stderr)
		if _dev_null is not None :
			PY.os.close (_dev_null)
		_delegate (_executable, _arguments, _environment)
	
	@_inject
	def __Z___exec_prepare_0 (_executable, _lookup, _arguments, *, _env = None, _env_overrides = None, _path = None, _path_prepend = None, _chdir = None, _stdin = None, _stdout = None, _stderr = None, _sudo = None, _sudo_user = None) :
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
		if _sudo is not None :
			if _sudo_user is None :
				_sudo_user = "root"
			assert PY.isinstance (_sudo_user, PY.str), "[6846e3f1]"
			if _executable != _arguments[0] :
				Z.panic (0x922dd24b, "sudo can't set argument0")
			_executable = "sudo"
			_arguments = ["sudo", "-u", _sudo_user, "-H", "-n", "--"] + _arguments
		else :
			assert _sudo_user is None, "[dd439dde]"
		return _executable, _lookup, _arguments, _environment, _chdir, _files
	
	@_inject
	def __Z__sudo_prepare (_user = None) :
		if _user is None or _user is True :
			_user = "root"
		assert PY.isinstance (_user, PY.str), "[9d018081]"
		Z.spawn ("sudo", "-u", _user, "-v", "-p", "[z-run:%08d] [>>]  SUDO authentication:  password for user `%%p` is required to invoke commands as user `%%U`;  enter password: " % Z.pid)
	
	@_inject
	def __Z__process_wait (_pid, *, _enforce = None) :
		_pid_0 = Z._pid (_pid)
		_pid, _outcome = PY.os.waitpid (_pid_0, 0)
		if PY.os.WIFEXITED (_outcome) :
			_outcome = PY.os.WEXITSTATUS (_outcome)
		elif PY.os.WIFSIGNALED (_outcome) :
			_outcome = 0 - PY.os.WTERMSIG (_outcome)
		else :
			Z.panic (0xb4179b04, "waiting `%s` failed with unknown outcome: %r", _pid, _outcome)
		if _enforce and _outcome != 0 :
			Z.panic ((_enforce, 0xc0e8ec5d), "waiting `%s` failed with status: %d", _pid, _outcome)
		if _pid == _pid_0 :
			return _outcome
		else :
			return _pid, _outcome
	
	@_inject
	def __Z__process_signal (_pid, _signal, *, _wait = False, _enforce = None) :
		_pid = Z._pid (_pid)
		PY.os.kill (_pid, _signal)
		if _wait :
			return Z.process_wait (_pid, _enforce = _enforce)
	
	@_inject
	def __Z__process_terminate (_pid, *, _wait = False, _enforce = None) :
		return Z.process_signal (_pid, PY.signal.SIGTERM)
	
	@_inject
	def __Z__process_kill (_pid, *, _wait = False, _enforce = None) :
		return Z.process_signal (_pid, PY.signal.SIGKILL)
	
	@_inject
	def __Z___pid (_process) :
		if _process is None :
			_pid = None
		elif PY.isinstance (_process, PY.builtins.int) :
			if _process >= 1 :
				_pid = _process
			else :
				Z.panic (0x1fe23810, "invalid process id (negative)")
		elif PY.isinstance (_process, PY.subprocess.Popen) :
			_pid = _process.pid
		else :
			Z.panic (0x3960636f, "invalid process id (unknown): %r", _pid)
		return _pid
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__pipeline (_commands, *, _wait = True, _fd_close = False, _enforce = True) :
		_count = PY.len (_commands)
		if _count == 0 :
			Z.panic (0x1b1812d7, "pipeline empty")
		_pipes = []
		_pipes.append ((None, None))
		for _index in PY.range (_count - 1) :
			_pipes.append (PY.os.pipe ())
		_pipes.append ((None, None))
		_processes = []
		for _index in PY.range (_count) :
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
		_signal_handler_old = PY.signal.signal (PY.signal.SIGCHLD, lambda _1, _2 : None)
		while True :
			for _process in _processes :
				if _process is None :
					continue
				_index, _process, _arguments = _process
				if _terminated == (_count - 1) :
					_process.wait ()
				if _process.poll () is None :
					continue
				_terminated += 1
				if _process.returncode != 0 :
					_succeeded = False
					if _enforce :
						Z.log_warning (0x76d05a67, "spawn `%s` `%s` failed with status: %d", _arguments[0], _arguments[1:], _process.returncode)
				_processes[_index] = None
			if _terminated == _count :
				break
			PY.signal.sigtimedwait ([PY.signal.SIGCHLD], 6)
		PY.signal.signal (PY.signal.SIGCHLD, _signal_handler_old)
		if _enforce and not _succeeded :
			Z.panic ((_enforce, 0x1d6fad91), "pipeline failed")
		return _succeeded
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__log_error (_code, _message, *_arguments) :
		Z._log_write ("ee", _code, _message, _arguments)
	
	@_inject
	def __Z__log_warning (_code, _message, *_arguments) :
		if not Z.log_warning_enabled : return
		Z._log_write ("ww", _code, _message, _arguments)
	
	@_inject
	def __Z__log_notice (_code, _message, *_arguments) :
		if not Z.log_warning_enabled or not Z.log_notice_enabled : return
		Z._log_write ("ii", _code, _message, _arguments)
	
	@_inject
	def __Z__log_debug (_code, _message, *_arguments) :
		if not Z.log_warning_enabled or not Z.log_notice_enabled or not Z.log_debug_enabled : return
		Z._log_write ("dd", _code, _message, _arguments)
	
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
	def __Z__sleep (_interval) :
		PY.time.sleep (_interval)
	
	@_inject
	def __Z__time_as_seconds_since_epoch () :
		_time = PY.time.time ()
		return _time
	
	@_inject
	def __Z__time_as_token (_reference = None, *, _utc = False) :
		if _utc :
			_time = PY.time.gmtime (_reference)
		else :
			_time = PY.time.localtime (_reference)
		_time = PY.time.strftime ("%Y-%m-%d-%H-%M-%S", _time)
		return _time
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__panic (_code, _message, *_arguments) :
		_code = Z._panic_code (_code)
		Z._log_write ("!!", _code, _message, _arguments)
		Z.exit (1)
	
	@_inject
	def __Z__panic_not_implemented (_code) :
		Z.panic (_code, "not implemented")
	
	@_inject
	def __Z___panic_with_excepthook (_error_type, _error, _traceback) :
		_traceback_error = PY.traceback.extract_tb (_traceback)
		_traceback_caller = []
		Z._panic_with_traceback (0x338a0a5b, _error, _traceback_error, _traceback_caller)
	
	@_inject
	def __Z___panic_with_traceback (_code, _error, _traceback_error, _traceback_caller) :
		
		_code = Z._panic_code (_code)
		Z._log_write ("!!", _code, "unexpected error encountered:", ())
		Z._log_write ("!!", _code, "  %r", (_error,))
		
		def _traceback_frame_log (_tag, _frame) :
			_file = _frame.filename
			_line = _frame.lineno
			_context = _frame.name
			if _file == Z._scriptlet_begin_file :
				if _context == "<module>" :
					_context = "<scriptlet>"
				elif _context == "__Z__wrapper" :
					return
				elif _context.startswith ("__Z__") :
					_context = "Z." + _context[5:]
				if _line > Z._scriptlet_begin_line :
					_line -= Z._scriptlet_begin_line
					_line += Z._scriptlet_source_line_start
					_position = "`:: %s` @ %d @ `%s`" % (Z._scriptlet_label, _line, Z._scriptlet_source_path)
				else :
					_position = "<python3+> @ %d" % (_line)
			else :
				_position = "`%s` @ %d" % (_file, _line)
			Z._log_write ("!!", _code, "  [%s]  %-30s | %s", (_tag, _context, _position))
		
		for _traceback_frame in PY.reversed (_traceback_error) :
			_traceback_frame_log ("raised", _traceback_frame)
		for _traceback_frame in PY.reversed (_traceback_caller) :
			_traceback_frame_log ("caller", _traceback_frame)
		
		Z.panic (_code, "aborting!")
	
	@_inject
	def __Z___panic_code (_code) :
		_code_fallback = 0xee4006b2
		if PY.isinstance (_code, PY.int) :
			pass
		elif _code is None :
			_code = _code_fallback
		elif PY.isinstance (_code, PY.tuple) :
			_code_0 = _code_fallback
			for _code_0 in _code :
				if PY.isinstance (_code, PY.int) :
					break
				_code_0 = _code_fallback
			_code = _code_0
		else :
			_code = _code_fallback
		return _code
	
	@_inject
	def __Z___scriptlet_begin_from_fd (_fd, _label, _source_path, _source_line_start, _source_line_end) :
		PY.os.close (_fd)
		Z._scriptlet_label = _label
		Z._scriptlet_source_path = _source_path
		Z._scriptlet_source_line_start = _source_line_start
		Z._scriptlet_source_line_end = _source_line_end
		_traceback_frame = PY.traceback.extract_stack () [0]
		Z._scriptlet_begin_file = _traceback_frame.filename
		Z._scriptlet_begin_line = _traceback_frame.lineno
		PY.sys.excepthook = Z._panic_with_excepthook
	
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
		_actual = PY.len (Z.arguments)
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
			return PY.tuple (Z.arguments)
		else :
			if _rest is None or _rest is True :
				return PY.tuple (Z.arguments[:_min]) + (PY.list (Z.arguments[_min:]),)
			else :
				return PY.tuple (Z.arguments)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__enforce (_condition, *, _value = None, _enforce = None, _message = None) :
		if _message is None : _message = "enforcement failed"
		if not PY.isinstance (_condition, PY.builtins.bool) :
			if _enforce is None : _enforce = 0x97ee7cf1
			Z.panic (_enforce, _message)
		if not _condition :
			if _enforce is None : _enforce = 0x0e36cc55
			Z.panic (_enforce, _message)
		return _condition
	
	@_inject
	def __Z__enforce_prefix (_value, _prefix, **_options) :
		return Z.enforce (PY.isinstance (_value, PY.str) and _value.startswith (_prefix), value = _value, **_options)
	
	@_inject
	def __Z__enforce_suffix (_value, _suffix, **_options) :
		return Z.enforce (PY.isinstance (_value, PY.str) and _value.endswith (_suffix), value = _value, **_options)
	
	@_inject
	def __Z__enforce_regex (_value, _pattern, **_options) :
		_pattern = Z.regex (_pattern)
		return Z.enforce (PY.isinstance (_value, PY.str) and _pattern.match (_value) is not None, value = _value, **_options)
	
	@_inject
	def __Z__regex (_pattern) :
		return PY.re.compile (_pattern, PY.re.ASCII | PY.re.DOTALL)
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__path (_path, *, _absolute = False, _canonical = False, _relative = None) :
		if PY.isinstance (_path, PY.tuple) or PY.isinstance (_path, PY.list) :
			_path = Z.path_join (*_path)
		assert PY.isinstance (_path, PY.str) or PY.isinstance (_path, PY.bytes), "[cb7b97d3]"
		_path = Z.path_normalize (_path)
		if _absolute :
			_path = Z.path_absolute (_path)
		if _canonical :
			_path = Z.path_canonical (_path)
		if _relative is not None :
			_path = Z.path_relative (_path, _relative)
		return _path
	
	@_inject
	def __Z__path_join (*_parts) :
		_path = PY.path.join (*_parts)
		_path = Z.path_normalize (_path)
		return _path
	
	@_inject
	def __Z__path_dirname (_path) :
		return PY.path.dirname (_path)
	
	@_inject
	def __Z__path_basename (_path) :
		return PY.path.basename (_path)
	
	@_inject
	def __Z__path_split_last (_path) :
		_dirname, _basename = PY.path.split (_path)
		return _dirname, _basename
	
	@_inject
	def __Z__path_split_all (_path) :
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
		_components = PY.tuple (_components)
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
	def __Z__path_normalize (_path) :
		_path = PY.path.normpath (_path)
		if _path.startswith ("//") :
			_path = "/" + _path.lstrip ("/")
		return _path
	
	@_inject
	def __Z__path_absolute (_path) :
		return PY.path.abspath (_path)
	
	@_inject
	def __Z__path_canonical (_path) :
		return PY.path.realpath (_path)
	
	@_inject
	def __Z__path_relative (_path, _base) :
		_path = PY.path.relpath (_path, _base)
		if not _path.startswith (".") and not _path.startswith ("/") :
			_path = PY.path.join (".", _path)
		return _path
	
	@_inject
	def __Z__path_matches (_path, _pattern) :
		return PY.fnmatch.fnmatch (_path, _pattern)
	
	@_inject
	def __Z__path_temporary_for (_path, **_options) :
		_dirname, _basename = PY.path.split (_path)
		return Z.path_temporary_in (_dirname, _basename, **_options)
	
	@_inject
	def __Z__path_temporary_in (_path, _name, *, _prefix = ".tmp.", _infix = ".", _suffix = "", _token = 8, _pid = True) :
		if _path is None :
			_path = Z.environment_or_none.TMPDIR
		if _path is None :
			_path = "/tmp"
		_name = _name.strip (".")
		if _name == "" or "/" in _name :
			Z.panic (0x9fbff49a, "invalid path")
		_token = Z.random_token (_token)
		_name = _prefix + (PY.str (Z.pid) + "-" if _pid else "") + _token + _infix + _name + _suffix
		_path = PY.path.join (_path, _name)
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
	def __Z__exists (_path, *, _follow = True, _enforce = False) :
		_stat = Z.stat (_path, _follow = _follow)
		if _stat is None and _enforce :
			Z.panic ((_enforce, 0x383d3cc5), "file-system path does not exist `%s`", _path)
		return _stat is not None
	
	@_inject
	def __Z__not_exists (_path, *, _follow = True, _enforce = False) :
		_stat = Z.stat (_path, _follow = _follow)
		if _stat is not None and _enforce :
			Z.panic ((_enforce, 0x9064abfc), "file-system path already exists `%s`", _path)
		return _stat is None
	
	@_inject
	def __Z__is_file (_path, *, _follow = True, _enforce = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode)), _follow = _follow, _enforce = _enforce)
	
	@_inject
	def __Z__is_file_empty (_path, *, _follow = True, _enforce = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode) and _stat.st_size == 0), _follow = _follow, _enforce = _enforce)
	
	@_inject
	def __Z__is_file_not_empty (_path, *, _follow = True, _enforce = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode) and _stat.st_size > 0), _follow = _follow, _enforce = _enforce)
	
	@_inject
	def __Z__is_folder (_path, *, _follow = True, _enforce = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISDIR (_stat.st_mode)), _follow = _follow, _enforce = _enforce)
	
	@_inject
	def __Z__is_file_or_folder (_path, *, _follow = True, _enforce = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISREG (_stat.st_mode) or PY.stat.S_ISDIR (_stat.st_mode)), _follow = _follow, _enforce = _enforce)
	
	@_inject
	def __Z__is_symlink (_path, *, _enforce = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISLNK (_stat.st_mode)), _follow = False, _enforce = _enforce)
	
	@_inject
	def __Z__is_pipe (_path, *, _follow = True, _enforce = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISFIFO (_stat.st_mode)), _follow = _follow, _enforce = _enforce)
	
	@_inject
	def __Z__is_socket (_path, *, _follow = True, _enforce = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISSOCK (_stat.st_mode)), _follow = _follow, _enforce = _enforce)
	
	@_inject
	def __Z__is_dev_block (_path, *, _follow = True, _enforce = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISBLK (_stat.st_mode)), _follow = _follow, _enforce = _enforce)
	
	@_inject
	def __Z__is_dev_char (_path, *, _follow = True, _enforce = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISCHR (_stat.st_mode)), _follow = _follow, _enforce = _enforce)
	
	@_inject
	def __Z__is_special (_path, *, _follow = True, _enforce = False) :
		return Z._stat_check (_path, (lambda _stat : PY.stat.S_ISFIFO (_stat.st_mode) or PY.stat.S_ISSOCK (_stat.st_mode) or PY.stat.S_ISBLK (_stat.st_mode) or PY.stat.S_ISCHR (_stat.st_mode)), _follow = _follow, _enforce = _enforce)
	
	@_inject
	def __Z___stat_check (_path, _check, *, _follow = True, _enforce = False) :
		_stat = Z.stat (_path, _follow = _follow)
		if _stat is None :
			if _enforce :
				Z.panic ((_enforce, 0xa23c577e), "file-system path not found `%s`", _path)
			else :
				return None
		if _check (_stat) :
			return True
		else :
			if _enforce :
				Z.panic ((_enforce, 0xfdbdc9a5), "file-system stat check failed `%s`", _path)
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
		except PY.OSError as _error :
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
	def __Z__rmdir (_path) :
		PY.os.rmdir (_path)
	
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
	def __Z__unlink (_target) :
		PY.os.unlink (_target)
	
	@_inject
	def __Z__symlink (_source, _target) :
		PY.os.symlink (_source, _target)
	
	@_inject
	def __Z__file_read (_path, *, _enforce = True, **_options) :
		_fd = Z.fd_open_for_read (_path, _enforce = _enforce)
		if _fd is None :
			return None
		_buffers = []
		while True :
			_buffer = PY.os.read (_fd, 1024 * 1024)
			if PY.len (_buffer) == 0 :
				break
			_buffers.append (_buffer)
		PY.os.close (_fd)
		_data = b"".join (_buffers)
		_data = Z._file_read_process (_data, **_options)
		return _data
	
	@_inject
	def __Z___file_read_process (_output, *, _data = None, _line = None, _lines = None, _separator = None, _json = None) :
		assert PY.isinstance (_output, PY.bytes) or PY.isinstance (_output, PY.str), "[27429689]"
		assert _line is None or _line is True, "[ba79d691]"
		assert _lines is None or _lines is True, "[fdfc8771]"
		assert _json is None or _json is True, "[5ff197c3]"
		assert _separator is None or PY.isinstance (_separator, PY.str), "[4002d3ef]"
		if _data :
			assert not _line and not _lines and not _separator and not _json, "[0d9c84ca]"
			if _data is PY.str :
				if PY.isinstance (_output, PY.bytes) :
					_output = _output.decode ("utf-8")
			elif _data is PY.bytes :
				if PY.isinstance (_output, PY.str) :
					_output = _output.encode ("utf-8")
			else :
				assert False, "[d2c9b0f9]"
		if _line or _lines :
			assert not _data, "[be243d10]"
			if PY.isinstance (_output, PY.bytes) :
				_output = _output.decode ("utf-8")
			if _separator is None :
				_separator = "\n"
			else :
				assert _separator != "", "[357b50e2]"
		else :
			assert not _separator, "[f3bdf0eb]"
		if _json :
			assert not _data, "[be2949cc]"
		if _line :
			assert not _lines, "[cefc2173]"
			if _output == "" :
				_output = None
			else :
				if _output[0 - PY.len (_separator) :] == _separator :
					_output = _output[: 0 - PY.len (_separator)]
				_output = _output.split (_separator)
				if PY.len (_output) == 1 or (PY.len (_output) == 2 and _output[1] == "") :
					_output = _output[0]
				else :
					Z.panic (0x5a1b7a79, "output is made of multiple lines")
		elif _lines :
			assert not _line, "[ddd9a8bc]"
			if _output == "" :
				_output = None
			else :
				if _output[0 - PY.len (_separator) :] == _separator :
					_output = _output[: 0 - PY.len (_separator)]
				_output = _output.split (_separator)
		if _json :
			if _output is None :
				Z.panic (0xbd0b4c92, "output is empty")
			if _line :
				_output = PY.json.loads (_output)
			elif _lines :
				_output = [PY.json.loads (_output) for _output in _output]
			else :
				if PY.isinstance (_output, PY.bytes) :
					_output = _output.decode ("utf-8")
				_output = PY.json.loads (_output)
		return _output
	
	@_inject
	def __Z__file_write (_path, _data, *, _json = None, _mode = None, _replace = None, _create = None, _exclusive = None, _append = None, _truncate = None, _enforce = True) :
		if _create is None and _replace is None :
			_replace = True
		elif _replace is not None :
			assert _create is None and _exclusive is None and _append is None and _truncate is None, "[349cbb39]"
		elif _create is not None :
			assert _replace is None, "[8357ff43]"
			if _exclusive is None :
				_exclusive = True
		if _replace :
			_path_temporary = Z.path_temporary_for (_path)
			_fd = Z.fd_open_for_write (_path_temporary, _mode = _mode, _create = True, _exclusive = True, _enforce = _enforce)
		else :
			_fd = Z.fd_open_for_write (_path, _mode = _mode, _create = _create, _exclusive = _exclusive, _append = _append, _truncate = _truncate, _enforce = _enforce)
		if _fd is None :
			return False
		if _json :
			_data = PY.json.dumps (_data)
		if _data is None :
			_buffer = b""
		elif PY.isinstance (_data, PY.str) :
			_buffer = _data.encode ("utf-8")
		elif PY.isinstance (_data, PY.bytes) :
			_buffer = _data
		else :
			Z.panic (0x5fbb8542, "invalid data type")
		while PY.len (_buffer) > 0 :
			_offset = PY.os.write (_fd, _buffer)
			_buffer = _buffer[_offset:]
		PY.os.close (_fd)
		if _replace :
			Z.rename (_path_temporary, _path)
		return True
	
	## --------------------------------------------------------------------------------
	
	@_inject
	def __Z__fd_open_for_read (_path, *, _close_on_exec = True, _enforce = True) :
		_flags = PY.os.O_RDONLY | PY.os.O_NOCTTY
		if _close_on_exec :
			_flags |= PY.os.O_CLOEXEC
		try :
			_file = PY.os.open (_path, _flags)
		except PY.OSError as _error :
			if _error.errno == PY.errno.ENOENT :
				if _enforce :
					Z.panic ((_enforce, 0x9ca19f87), "open `%s` failed:  does not exist", _path)
			else :
				raise
		return _file
	
	@_inject
	def __Z__fd_open_for_write (_path, *, _mode = None, _create = True, _exclusive = None, _read = False, _append = False, _truncate = False, _close_on_exec = True, _enforce = True) :
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
		except PY.OSError as _error :
			if _error.errno == PY.errno.ENOENT :
				if _enforce :
					Z.panic ((_enforce, 0x008389df), "open `%s` failed:  does not exist", _path)
			elif _error.errno == PY.errno.EEXIST :
				if _enforce :
					Z.panic ((_enforce, 0x37020d90), "open `%s` failed:  already exists", _path)
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
	def __Z__is_fd (_file, *, _enforce = None) :
		if _file is None :
			if _enforce :
				Z.panic ((_enforce, 0x3c3ff99e), "invalid file descriptor (none)")
			return False
		elif PY.isinstance (_file, PY.builtins.int) :
			if _file >= 0 :
				return True
			else :
				Z.panic (0xf42752f0, "invalid file descriptor (negative)")
		elif PY.isinstance (_file, PY.io.IOBase) :
			try :
				_file.fileno ()
				return True
			except PY.OSError as _error :
				Z.panic (0x7f97a33f, "invalid file descriptor (not supported): %r  //  %s", _error)
		else :
			if _enforce :
				Z.panic ((_enforce, 0xe0d78a09), "invalid file descriptor (unknown): %r", _file)
			return False
	
	@_inject
	def __Z___fd (_file) :
		if _file is None :
			_fd = None
		elif PY.isinstance (_file, PY.builtins.int) :
			if _file >= 0 :
				_fd = _file
			else :
				Z.panic (0x4327e646, "invalid file descriptor (negative)")
		elif PY.isinstance (_file, PY.io.IOBase) :
			try :
				_fd = _file.fileno ()
			except PY.OSError as _error :
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
	def __Z__data_md5 (_value) :
		return Z._data_hash (PY.hashlib.md5 (), _value)
	
	@_inject
	def __Z__data_sha1 (_value) :
		return Z._data_hash (PY.hashlib.sha1 (), _value)
	
	@_inject
	def __Z__data_sha2_224 (_value) :
		return Z._data_hash (PY.hashlib.sha224 (), _value)
	
	@_inject
	def __Z__data_sha2_256 (_value) :
		return Z._data_hash (PY.hashlib.sha256 (), _value)
	
	@_inject
	def __Z__data_sha2_384 (_value) :
		return Z._data_hash (PY.hashlib.sha384 (), _value)
	
	@_inject
	def __Z__data_sha2_512 (_value) :
		return Z._data_hash (PY.hashlib.sha512 (), _value)
	
	@_inject
	def __Z__data_sha3_224 (_value) :
		return Z._data_hash (PY.hashlib.sha3_224 (), _value)
	
	@_inject
	def __Z__data_sha3_256 (_value) :
		return Z._data_hash (PY.hashlib.sha3_256 (), _value)
	
	@_inject
	def __Z__data_sha3_384 (_value) :
		return Z._data_hash (PY.hashlib.sha3_384 (), _value)
	
	@_inject
	def __Z__data_sha3_512 (_value) :
		return Z._data_hash (PY.hashlib.sha3_512 (), _value)
	
	@_inject
	def __Z__data_blake2b (_value, *, _size = 64) :
		return Z._data_hash (PY.hashlib.blake2b (digest_size = _size), _value)
	
	@_inject
	def __Z__data_blake2s (_value, *, _size = 32) :
		return Z._data_hash (PY.hashlib.blake2s (digest_size = _size), _value)
	
	@_inject
	def __Z___data_hash (_hasher, _value) :
		if PY.isinstance (_value, PY.str) :
			_value = _value.encode ("utf-8")
		else :
			assert PY.isinstance (_value, PY.bytes), "[c682f292]"
		_hasher.update (_value)
		_hash = _hasher.hexdigest ()
		return _hash
	
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
		_data = _data.decode ("utf-8")
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
		_sequence = PY.list (_sequence)
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
	Z.arguments = PY.tuple (PY.sys.argv[1:])
	Z.environment = Z_environment ()
	Z.environment_or_none = Z_environment_or_none ()
	
	Z.executable = Z.environment.ZRUN_EXECUTABLE
	
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

