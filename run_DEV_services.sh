#!/bin/bash
# set -o xtrace
set -o nounset
set -o pipefail

readonly work_dir="$(dirname "$(readlink --canonicalize-existing "${0}")")"
readonly service_dir="${work_dir}/test-service"
readonly jar_file="${service_dir}/target/health-mock-0.0.1-SNAPSHOT.jar"
readonly error_jar_not_found=80
readonly error_csv_file_not_found=81
readonly dev_csv="${work_dir}/dev.csv"

if [[ ! -f "${jar_file}" ]]; then
    echo "error: ${jar_file} not found" >&2
    exit ${error_jar_not_found}
fi

if [[ ! -f "${dev_csv}" ]]; then
    echo "error: ${dev_csv} not found" >&2
    exit ${error_csv_file_not_found}
fi

while read port; do
    nohup java -jar "${jar_file}" --server.port="${port}" &
done < <(awk -F ',' '{print $5}' "${dev_csv}" | sed 1d | grep --invert-match --extended-regexp '^$')

exit 0