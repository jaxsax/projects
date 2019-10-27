#!/usr/bin/env bash
# https://stackoverflow.com/questions/192249/how-do-i-parse-command-line-arguments-in-bash
# set -o errexit -o pipefail -o noclobber -o nounset

MIN_MINUTES=59
CACHE_FOLDER="/tank/bazel/cache"

usage() {
    echo "$0: [cache_folder|$CACHE_FOLDER] [flags]"
    printf "\t-m, --min\tMinimum minutes before considering a file expired\n"
}

OPTS=$(getopt -o hm: --long help,min: -- "$@")

if [ $? != 0 ]; then echo "Invalid options" >&2; exit 1; fi

echo "$OPTS"
eval set -- "$OPTS"

while true; do
    case "$1" in
        -h | --help)
            usage; shift
            exit 1
            ;;
        -m | --min)
            MIN_MINUTES="$2"; shift; shift
            ;;
        --)
            shift; break;
            ;;
        *) break ;;
    esac
    shift
done

if [[ $# -ge 1 ]]; then
    CACHE_FOLDER="$1"
fi

find "$CACHE_FOLDER" -mmin "+$MIN_MINUTES" -type f