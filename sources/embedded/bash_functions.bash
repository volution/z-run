#!/dev/null


function _ () {
	
	local -- _outcome=0
	
	local -- _terminal _print
	if test -t 2 ; then
		case "${TERM}" in
			( screen | screen.* )
				_terminal='OK'
			;;
		esac
	fi
	
	local -- _label _snippet
	if test "${#}" -eq 0 ; then
		if test -n "${_terminal}" ; then
			printf -- '\x1b]0;[z-run:select]\x07' >&2
		fi
		local -- _label_and_snippet
		_label_and_snippet="$(
				export -- SHLVL=0 OLDPWD="${PWD}"
				exec -- z-run select-export-scriptlet-label-and-body
			)" || _outcome="${?}"
		if test -n "${_terminal}" ; then
			printf -- '\x1b]0;\x07' >&2
		fi
		if test -z "${_label_and_snippet}" -a "${_outcome}" -eq 0 ; then
			printf -- '[z-run] [ii]  quit.\n' >&2
			return -- 0
		fi
		if test -z "${_label_and_snippet}" -o "${_outcome}" -ne 0 ; then
			printf -- '[z-run] [ee]  unexpected failure;  aborting!\n' >&2
			return -- "${_outcome}"
		fi
		_label="${_label_and_snippet%%$'\n'*}"
		_snippet="${_label_and_snippet#*$'\n'}"
		_print=true
		if test -z "${_label}" -o -z "${_snippet}" ; then
			printf -- '[z-run] [ee]  unexpected failure;  aborting!\n' >&2
			return -- 2
		fi
	else
		_label="${1}"
		_print=false
		shift -- 1
		_snippet="$(
				export -- SHLVL=0 OLDPWD="${PWD}"
				exec -- z-run select-export-scriptlet-body "${_label}"
			)" || _outcome="${?}"
		if test -z "${_snippet}" -a "${_outcome}" -eq 0 ; then
			return -- 0
		fi
		if test -z "${_snippet}" -o "${_outcome}" -ne 0 ; then
			printf -- '[z-run] [ee]  unexpected failure;  aborting!\n' >&2
			return -- "${_outcome}"
		fi
	fi
	
	local -- _history_tokens=( "${_label}" )
	if test "${#}" -ne 0 ; then
		_history_tokens+=( "${@}" )
	fi
	local -- _history_tokens_escaped=()
	local -- _history_token _history_token_escaped _history_token_quote=\'\\\'\'
	for _history_token in "${_history_tokens[@]}" ; do
		if ! [[ "${_history_token}" =~ ^[a-zA-Z0-9_./@-]+$ ]] ; then
			_history_token_escaped="'${_history_token//\'/${_history_token_quote}}'"
		else
			_history_token_escaped="${_history_token}"
		fi
		_history_tokens_escaped+=( "${_history_token_escaped}" )
	done
	history -s -- "_${_history_tokens_escaped[*]/#/ }"
	
	if test "${_print}" == true -a -n "${_terminal}" ; then
		printf -- '> _%s\n' "${_history_tokens_escaped[*]/#/ }" >&2
	fi
	
	if [[ "${_snippet} " =~ ^([[:space:]]*)(export([[:space:]]+--)?[[:space:]]+)?((([_a-zA-Z0-9]+)=([^[:space:][:cntrl:]]+|\'[^\']+\'|\"[^\"]+\")?)([[:space:]]+))+$ ]] || [[ "${_snippet} " =~ ^([[:space:]]*)(unset([[:space:]]+--)?[[:space:]]+)(([_a-zA-Z0-9]+)([[:space:]]+))+$ ]] ; then
		if test "${#}" -ne 0 ; then
			printf -- '[z-run] [ee]  unexpected arguments;  aborting!\n' >&2
			return -- 2
		else
			eval -- "${_snippet}"
			_outcome="${?}"
		fi
	else
		if test -n "${_terminal}" ; then
			printf -- '\x1b]0;[z-run] %s\x07' "${_label}" >&2
		fi
		(
			export -- SHLVL=0 OLDPWD="${PWD}"
			exec -- z-run "${_label}" "${@}"
		) || _outcome="${?}"
		if test -n "${_terminal}" ; then
			printf -- '\x1b]0;\x07' >&2
		fi
	fi
	return -- "${_outcome}"
}

