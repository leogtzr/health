#!/bin/bash
set -o xtrace
set -o nounset
set -o pipefail

readonly work_dir="$(dirname "$(readlink --canonicalize-existing "${0}")")"
readonly dev_csv="${work_dir}/dev.csv"
readonly error_csv_file_not_found=81
readonly error_missing_args=82

if [[ ! -f "${dev_csv}" ]]; then
    echo "error: ${dev_csv} not found" >&2
    exit ${error_csv_file_not_found}
fi

if ((${#} != 1)); then
    echo "error: missing argument" >&2
    exit ${error_missing_args}
fi

readonly instance_info=$(sed --quiet "${1}p" "${dev_csv}")
readonly instance_port=$(awk --field-separator ',' '{print $5}' <<< "${instance_info}")
readonly instance_url=$(awk --field-separator ',' '{print $6}' <<< "${instance_info}")

curl "http://localhost:${instance_port}/monitoring/up"

exit 0