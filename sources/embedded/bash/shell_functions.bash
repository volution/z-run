#!/dev/null


function _ () {
	
	local -- _zrun
	if test -n "${ZRUN_EXECUTABLE:-}" ; then
		_zrun="${ZRUN_EXECUTABLE}"
	else
		_zrun="$( type -P -- z-run || true )"
	fi
	if test -z "${_zrun}" ; then
		printf -- '[z-run] [ee]  missing `z-run`;  aborting!\n' >&2
		return -- 2
	fi
	_select_wrapper=()
	if test -n "${TMUX:-}" -a "${TERM:-dumb}" != dumb -a "$( type -P -- tmux-popup )" ; then
		_select_wrapper+=( env ZFZF_FULLSCREEN=true tmux-popup -- )
	fi
	
	local -- _history_command _history_save
	_history_command="$( history 1 )"
	_history_save=false
	if [[ "${_history_command}" =~ ^([[:space:]]*[0-9]+[[:space:]]+)(.*)$ ]] ; then
		_history_command="${BASH_REMATCH[2]}"
	else
		_history_command=''
	fi
	
	local -- _menu_trigger _menu_label
	if test "${#}" -eq 0 ; then
		_menu_trigger=true
		_menu_label=''
		if test "${_history_command}" == '_' ; then
			_history_save=true
		fi
	elif test "${#}" -eq 1 ; then
		if test "${1:0:2}" == '::' -a \( "${1}" == ':: *' -o "${1% / ...}" != "${1}" -o "${1% / \*}" != "${1}" \) ; then
			_menu_trigger=true
			_menu_label="${1}"
			shift -- 1
			if test "${_history_command}" == "_ '${_menu_label}'" ; then
				_history_save=true
			fi
		else
			_menu_trigger=false
			_menu_label=''
		fi
	else
		_menu_trigger=false
		_menu_label=''
	fi
	
	local -- _terminal=''
	if test -t 2 -a "${TERM:-dumb}" != dumb ; then
		_terminal='OK'
	fi
	
	local -- _label _snippet
	local -- _outcome=0
	if test "${_menu_trigger}" == true ; then
		if test -n "${_terminal}" ; then
			printf -- '\x1b]0;[z-run:menu]\x07' >&2
		fi
		local -- _label_and_snippet
		_label_and_snippet="$(
				export -- SHLVL=0 OLDPWD="${PWD}"
				if test -n "${_menu_label}" ; then
					exec -- "${_select_wrapper[@]}" "${_zrun}" select-export-scriptlet-label-and-body "${_menu_label}"
				else
					exec -- "${_select_wrapper[@]}" "${_zrun}" select-export-scriptlet-label-and-body
				fi
			)" || _outcome="${?}"
		if test -n "${_terminal}" ; then
			printf -- '\x1b]0;\x07' >&2
		fi
		if test -z "${_label_and_snippet}" -a "${_outcome}" -eq 0 ; then
			printf -- '[z-run] [ii]  menu quit.\n' >&2
			return -- 0
		fi
		if test -z "${_label_and_snippet}" -o "${_outcome}" -ne 0 ; then
			printf -- '[z-run] [ee]  unexpected failure;  aborting!\n' >&2
			return -- 2
		fi
		_label="${_label_and_snippet%%$'\n'*}"
		_snippet="${_label_and_snippet#*$'\n'}"
		if test -z "${_label}" -o -z "${_snippet}" ; then
			printf -- '[z-run] [ee]  unexpected failure;  aborting!\n' >&2
			return -- 2
		fi
	else
		_label="${1}"
		shift -- 1
		_snippet="$(
				export -- SHLVL=0 OLDPWD="${PWD}"
				exec -- "${_zrun}" select-export-scriptlet-body "${_label}"
			)" || _outcome="${?}"
		if test -z "${_snippet}" -a "${_outcome}" -eq 0 ; then
			printf -- '[z-run] [ee]  unexpected failure;  aborting!\n' >&2
			return -- 2
		fi
		if test -z "${_snippet}" -o "${_outcome}" -ne 0 ; then
			printf -- '[z-run] [ee]  unexpected failure;  aborting!\n' >&2
			return -- 2
		fi
	fi
	
	local -- _history_tokens=( "${_label}" )
	if test "${#}" -ne 0 ; then
		_history_tokens+=( "${@}" )
	fi
	local -- _history_tokens_escaped=()
	local -- _history_token _history_token_escaped _history_token_quote=\'\\\'\'
	for _history_token in "${_history_tokens[@]}" ; do
		if [[ "${_history_token}" =~ ^[a-zA-Z0-9_./@-]+$ ]] ; then
			_history_token_escaped="${_history_token}"
		else
			_history_token_escaped="'${_history_token//\'/${_history_token_quote}}'"
		fi
		_history_tokens_escaped+=( "${_history_token_escaped}" )
	done
	if test "${_history_save}" == true ; then
		if test -n "${_terminal}" ; then
			printf -- '> _%s\n' "${_history_tokens_escaped[*]/#/ }" >&2
		fi
		history -s -- "_${_history_tokens_escaped[*]/#/ }"
	fi
	
	if [[ "${_snippet} " =~ ^([[:space:]]*)(export([[:space:]]+--)?[[:space:]]+)?((([_a-zA-Z0-9]+)=([^[:space:][:cntrl:]]+|\'[^\']+\'|\"[^\"]+\")?)([[:space:]]+))+$ ]] || [[ "${_snippet} " =~ ^([[:space:]]*)(unset([[:space:]]+--)?[[:space:]]+)(([_a-zA-Z0-9]+)([[:space:]]+))+$ ]] ; then
		if test "${#}" -ne 0 ; then
			printf -- '[z-run] [ee]  unexpected arguments;  aborting!\n' >&2
			return -- 2
		else
			eval -- "${_snippet}" || _outcome="${?}"
		fi
	else
		if test -n "${_terminal}" ; then
			printf -- '\x1b]0;[z-run] %s\x07' "${_label}" >&2
		fi
		(
			export -- SHLVL=0 OLDPWD="${PWD}"
			exec -- "${_zrun}" "${_label}" "${@}"
		) || _outcome="${?}"
		if test -n "${_terminal}" ; then
			printf -- '\x1b]0;\x07' >&2
		fi
	fi
	
	return -- "${_outcome}"
}

