#!/dev/null

set -e -E -u -o pipefail -o noclobber -o noglob +o braceexpand || exit -- 1
trap 'printf -- "[ee] failed: %%s\n" "${BASH_COMMAND}" >&2' ERR || exit -- 1

BASH_ARGV0='z-run'
ZRUN=( "${ZRUN_EXECUTABLE}" )
X_RUN=( "${ZRUN[@]}" )
