#!/dev/null


if test -z "${BASH:-}" ; then
	echo "[ee]  failed:  expected \`bash\`, found \`${BASH:-none}\`;  aborting!" >&2
	exit 1
fi
if test -z "${BASH_VERSION:-}" ; then
	echo "[ee]  failed:  expected \`bash\` version minimum 5, found \`${BASH_VERSION:-none}\`;  aborting!" >&2
	exit 1
fi
if test "${BASH_VERSINFO[0]:-0}" -lt 5 ; then
	echo "[ee]  failed:  expected \`bash\` version minimum 5, found \`${BASH_VERSION:-none}\`;  aborting!" >&2
	exit 1
fi


set -e -E -u -o pipefail -o noclobber -o noglob +o braceexpand || exit -- 1
trap -- 'printf -- "[ee]  failed:  %s\\n" "${BASH_COMMAND}" >&2' ERR || exit -- 1


## sanity

if test "${PATH}" == '/usr/local/bin:/usr/bin:/bin:.' ; then
	PATH=''
fi
export -- PATH="${PATH:-/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin}"
export -- HOME="${HOME:-/tmp/user-${UID}--home}"
export -- USER="${USER:-user-${UID}}"
export -- TMPDIR="${TMPDIR:-/tmp}"
export -- SHELL="${SHELL:-${BASH}}"
export -- TERM="${TERM:-dumb}"
export -- SHLVL=1


## environment

if test -e "${HOME}/.bash/environment" ; then
	source -- "${HOME}/.bash/environment"
fi


## interactive

case "${TERM}" in
	( dumb )
		PS1='> '
	;;
	( screen | screen.* )
		case "${USER}" in
			( root )
				PS1='\n\n\[\033]0;\w\007\033kshell (\u@\h)\033\\\][\[\033[31;01m\]\u\[\033[00m\]@\[\033[35m\]\h\[\033[00m\]][\[\033[34m\]\w\[\033[00m\]]\n> '
			;;
			( * )
				PS1='\n\n\[\033]0;\w\007\033kshell (\u@\h)\033\\\][\[\033[32m\]\u\[\033[00m\]@\[\033[35m\]\h\[\033[00m\]][\[\033[34m\]\w\[\033[00m\]]\n> '
			;;
		esac
	;;
	( * )
		case "${USER}" in
			( root )
				PS1='\n\n[\[\033[31;01m\]\u\[\033[00m\]@\[\033[35m\]\h\[\033[00m\]][\[\033[34m\]\w\[\033[00m\]]\n> '
			;;
			( * )
				PS1='\n\n[\[\033[32m\]\u\[\033[00m\]@\[\033[35m\]\h\[\033[00m\]][\[\033[34m\]\w\[\033[00m\]]\n> '
			;;
		esac
	;;
esac
PS2='   > '
PS3=' ? > '
PS4='-  > '

PROMPT_DIRTRIM=0
PROMPT_COMMAND=''

INPUTRC=''
HISTFILE=''

HISTCONTROL=ignoredups:ignorespace:erasedups
HISTSIZE=16384
HISTFILESIZE=16384
HISTIGNORE=''

FIGNORE=''
GLOBIGNORE=''

unset -- MAILCHECK


## behaviour

IFS=''


## options

_set_options=(
		
		# expansion
		braceexpand=true
		histexpand=false
		nounset=true
		noglob=false
		
		# environment
		allexport=false
		
		# execution
		pipefail=true
		errexit=false
		errtrace=true
		functrace=true
		monitor=true
		hashall=true
		keyword=false
		noclobber=true
		physical=true
		notify=true
		verbose=false
		xtrace=false
		noexec=false
		onecmd=false
		
		# history
		## history=false
		
		# editing
		emacs=true
		vi=false
		ignoreeof=false
		
		# compatibility
		posix=false
		privileged=true
)

_shopt_options=(
		
		# globbing
		dotglob=true
		failglob=true
		nullglob=true
		nocaseglob=false
		extglob=true
		globstar=true
		
		# expansion
		extquote=true
		globasciiranges=true
		expand_aliases=true
		
		# history
		histappend=true
		histreedit=true
		histverify=true
		cmdhist=true
		lithist=true
		
		# execution
		inherit_errexit=true
		checkhash=true
		checkjobs=true
		huponexit=true
		execfail=false
		extdebug=false
		lastpipe=false
		sourcepath=false
		xpg_echo=true
		shift_verbose=true
		localvar_inherit=false
		localvar_unset=false
		nocasematch=false
		assoc_expand_once=false
		restricted_shell=false
		
		# miscellaneous
		checkwinsize=true
		promptvars=true
		gnu_errfmt=true
		
		# auto-complete and auto-magic
		no_empty_cmd_completion=true
		hostcomplete=false
		complete_fullquote=true
		direxpand=false
		dirspell=false
		force_fignore=true
		autocd=false
		cdable_vars=false
		cdspell=false
		interactive_comments=true
		mailwarn=false
		progcomp=false
		progcomp_alias=false
		
		# compatibility
		compat31=false
		compat32=false
		compat40=false
		compat41=false
		compat42=false
		compat43=false
		compat44=false
)

for _set_option in "${_set_options[@]}" ; do
	case "${_set_option}" in
		( *=true )
			set -o "${_set_option%=true}"
		;;
		( *=false )
			set +o "${_set_option%=false}"
		;;
		( * )
			printf -- '[ee]  invalid option: `%s`!\n' "${_set_option}" >&2
		;;
	esac
done
unset -- _set_option
unset -- _set_options

for _shopt_option in "${_shopt_options[@]}" ; do
	case "${_shopt_option}" in
		( *=true )
			shopt -s "${_shopt_option%=true}"
		;;
		( *=false )
			shopt -u "${_shopt_option%=false}"
		;;
		( * )
			printf -- '[ee]  invalid option: `%s`!\n' "${_shopt_option}" >&2
		;;
	esac
done
unset -- _shopt_option
unset -- _shopt_options


## bindings

if test -f "${HOME}/.bash/input" ; then
	bind -f "${HOME}/.bash/input"
fi


## functions

if test -f "${HOME}/.bash/functions" ; then
	source -- "${HOME}/.bash/functions"
fi

if test -f "${HOME}/.bash/aliases" ; then
	source -- "${HOME}/.bash/aliases"
fi


## finalize

if shopt -q -- login_shell ; then
	if test -e "${HOME}/Desktop" ; then
		cd -- "${HOME}/Desktop"
	elif test -e "${HOME}" ; then
		cd -- "${HOME}"
	fi
elif test -e "${HOME}" -a . -ef "${HOME}" ; then
	cd -- "${HOME}"
fi

set +e -E -u -o pipefail -o noclobber -o noglob -o braceexpand || exit -- 1
trap -- 'printf -- "\\n\\033[31;01m[ee] command failed (with exit code %s):\\n%s\\033[00m\\n" "${?}" "${BASH_COMMAND}" >&2' ERR || exit -- 1

