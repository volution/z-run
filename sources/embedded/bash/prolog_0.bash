#!/dev/null
set -e -E -u -o pipefail -o noclobber -o noglob +o braceexpand || exit -- 1
trap -- 'printf -- "[z-run:%08d] [%s] [%08x]  aborting because of failed command (with exit code %d):  %s\\n" "${$}" "!!" "0xee49d661" "${?}" "${BASH_COMMAND}" >&2 ; exit -- 1' ERR || exit -- 1
ZRUN=( "${ZRUN_EXECUTABLE:-z-run}" )
