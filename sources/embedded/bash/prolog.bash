#!/dev/null

set -e -E -u -o pipefail -o noclobber -o noglob +o braceexpand || exit -- 1
trap -- 'printf -- "[z-run:%08d] [%s] [%08x]  abording because of failed command (with exit code %d):  %s\\n" "${ZRUN_PID:-${$}}" "!!" "0x5239a746" "${?}" "${BASH_COMMAND}" >&2 ; exit -- 1' ERR || exit -- 1

if test -z "${ZRUN_EXECUTABLE:-}" ; then
	export -- ZRUN_EXECUTABLE="$( type -P z-run || true )"
fi
if test -z "${ZRUN_EXECUTABLE:-}" ; then
	printf -- '[z-run:%08d] [%s] [%08x]  %s\n' "${ZRUN_PID:-${$}}" '!!' 0xd1fd3ac8 'missing `z-run`;  aborting!' >&2
	exit -- 1
fi
ZRUN=( "${ZRUN_EXECUTABLE}" )
ZRUN_ARGUMENTS=( "${@}" )
ZRUN_PID="${$}"
BASH_ARGV0='z-run'

function Z_log_error () {
	Z_log_write 'ee' "${@}"
}
function Z_log_warning () {
	Z_log_write 'ww' "${@}"
}
function Z_log_notice () {
	Z_log_write 'ii' "${@}"
}
function Z_log_debug () {
	Z_log_write 'dd' "${@}"
}
function Z_log_cut () {
	printf -- '\n[z-run:%08d] [--]             --------------------------------------------------------------------------------\n\n' "${ZRUN_PID:-${$}}" >&2
}

function Z_log_write () {
	if test "${#}" -lt 3 ; then
		printf -- '[z-run:%08d] [%s] [%08x]  %s\n' "${ZRUN_PID:-${$}}" '!!' 0x83e15bdc 'invalid log syntax!' >&2
	fi
	printf -- '[z-run:%08d] [%s] [%08x]  '"${3}"'\n' "${ZRUN_PID:-${$}}" "${1}" "${2}" "${@:4}" >&2
}

function Z_panic () {
	Z_log_write '!!' "${@}"
	if test "${$}" -ne "${ZRUN_PID:-${$}}" ; then
		kill -s SIGTERM -- "${ZRUN_PID:-${$}}"
	fi
	exit -- 1
}

function Z_zspawn () {
	if test "${#}" -eq 0 ; then
		Z_panic 0x85324aa7 'failed to spawn `z-run`:  missing scriptlet;  aborting!'
	fi
	if test "${1:0:2}" != '::' ; then
		Z_panic 0xf2d62e78 'failed to spawn `z-run`:  invalid scriptlet;  aborting!'
	fi
	if "${ZRUN}" "${@}" ; then
		return -- 0
	else
		Z_panic 0x64aa44ac 'failed to spawn `%s` (with exit code `%d`);  aborting!' "${1}" "${?}"
	fi
}

function Z_zspawn_return () {
	if test "${#}" -eq 0 ; then
		Z_panic 0xa53618e4 'failed to spawn `z-run`:  missing scriptlet;  aborting!'
	fi
	if test "${1:0:2}" != '::' ; then
		Z_panic 0x122f147e 'failed to spawn `z-run`:  invalid scriptlet;  aborting!'
	fi
	if "${ZRUN}" "${@}" ; then
		return -- 0
	else
		return -- "${?}"
	fi
}

function Z_zexec () {
	if test "${#}" -eq 0 ; then
		Z_panic 0x83b4d099 'failed to exec `z-run`:  missing scriptlet;  aborting!'
	fi
	if test "${1:0:2}" != '::' ; then
		Z_panic 0x185a6b16 'failed to exec `z-run`:  invalid scriptlet;  aborting!'
	fi
	if exec -- "${ZRUN}" "${@}" ; then
		return -- 0
	else
		Z_panic 0x143023c2 'failed to exec `z-run` (with exit code `%d`);  aborting!' "${?}"
	fi
}

function Z_enforce () {
	if test "${#}" -lt 3 ; then
		Z_panic 0x8a60b175 'invalid enforceme syntax;  aborting!'
	fi
	if "${@:3}" ; then
		return -- 0
	else
		if test "${2}" != '' ; then
			Z_panic "${1}" "${2}"
		else
			Z_panic "${1}" 'enforcement failed;  aborting!'
		fi
	fi
}

function Z_expect_no_arguments () {
	if test "${#ZRUN_ARGUMENTS[@]}" -ne 0 ; then
		Z_panic 0xc6f3aeb6 'no arguments expected;  aborting!'
	fi
}

