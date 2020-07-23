#!/bin/bash
# set -o xtrace
set -o nounset
set -o pipefail

readonly work_dir="$(dirname "$(readlink --canonicalize-existing "${0}")")"
readonly dev_csv="${work_dir}/dev.csv"
readonly error_csv_file_not_found=81

if [[ ! -f "${dev_csv}" ]]; then
    echo "error: ${dev_csv} not found" >&2
    exit ${error_csv_file_not_found}
fi

readonly number_of_lines=$(sed 1d "${dev_csv}" | wc --lines)
readonly instance_line=$(awk -v min=2 -v max=${number_of_lines} 'BEGIN{srand();print int(rand()*(max-min))+min}')

readonly instance_info=$(sed --quiet "${instance_line}p" "${dev_csv}")
readonly instance_port=$(awk --field-separator ',' '{print $5}' <<< "${instance_info}")
readonly instance_url=$(awk --field-separator ',' '{print $6}' <<< "${instance_info}")

curl "http://localhost:${instance_port}/monitoring/down"

exit 0
