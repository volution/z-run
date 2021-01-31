#!/dev/null

set -e -E -u -o pipefail -o noclobber -o noglob +o braceexpand || exit -- 1
trap -- 'printf -- "[ee]  failed:  %s\\n" "${BASH_COMMAND}" >&2 ; exit -- 1' ERR || exit -- 1

if test -z "${ZRUN_EXECUTABLE:-}" ; then
	export -- ZRUN_EXECUTABLE="$( type -P z-run )"
fi
ZRUN=( "${ZRUN_EXECUTABLE}" )
ZRUN_ARGUMENTS=( "${@}" )
ZRUN_PID="${$}"

BASH_ARGV0='z-run'
X_RUN=( "${ZRUN[@]}" )

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

function Z_log_write () {
	if test "${#}" -lt 3 ; then
		printf -- '[z-run:%08d] [%s] [%08x]  invalid syntax!\n' "${ZRUN_PID}" '!!' 0x83e15bdc
	fi
	printf -- '[z-run:%08d] [%s] [%08x]  '"${3}"'\n' "${ZRUN_PID}" "${1}" "${2}" "${@:4}"
}

function Z_panic () {
	Z_log_write '!!' "${@}"
	kill -s SIGTERM -- "${ZRUN_PID}"
	exit -- 1
}

function Z_zspawn () {
	if test "${#}" -eq 0 ; then
		Z_panic 0x85324aa7 'failed to spawn z-run:  missing scriptlet!'
	fi
	if test "${1:0:2}" != '::' ; then
		Z_panic 0xf2d62e78 'failed to spawn z-run:  invalid scriptlet!'
	fi
	if "${ZRUN[@]}" "${@}" ; then
		return -- 0
	else
		Z_panic 0x64aa44ac 'failed to spawn z-run:  exited with status `%d`!' "${?}"
	fi
}

function Z_zspawn_return () {
	if test "${#}" -eq 0 ; then
		Z_panic 0xa53618e4 'failed to spawn z-run:  missing scriptlet!'
	fi
	if test "${1:0:2}" != '::' ; then
		Z_panic 0x122f147e 'failed to spawn z-run:  invalid scriptlet!'
	fi
	if "${ZRUN[@]}" "${@}" ; then
		return -- 0
	else
		return -- "${?}"
	fi
}

function Z_zexec () {
	if test "${#}" -eq 0 ; then
		Z_panic 0x83b4d099 'failed to exec z-run:  missing scriptlet!'
	fi
	if test "${1:0:2}" != '::' ; then
		Z_panic 0x185a6b16 'failed to exec z-run:  invalid scriptlet!'
	fi
	if exec -- "${ZRUN[@]}" "${@}" ; then
		return -- 0
	else
		Z_panic 0x143023c2 'failed to exec z-run:  exited with status `%d`!' "${?}"
	fi
}

function Z_enforce () {
	if test "${#}" -lt 3 ; then
		Z_panic 0x8a60b175 'failed enforcement:  invalid syntax!'
	fi
	if "${@:3}" ; then
		return -- 0
	else
		if test "${2}" != '' ; then
			Z_panic "${1}" "${2}"
		else
			Z_panic "${1}" 'enforcement failed!'
		fi
	fi
}

