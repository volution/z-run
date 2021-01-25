#!/dev/null
set -e -E -u -o pipefail -o noclobber -o noglob +o braceexpand || exit -- 1
trap -- 'printf -- "[ee]  failed:  %s\\n" "${BASH_COMMAND}" >&2 ; exit -- 1' ERR || exit -- 1
ZRUN=( "${ZRUN_EXECUTABLE}" )
