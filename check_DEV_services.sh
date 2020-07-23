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

while read service_info; do
    #curl "${service_url}"
    # echo "${service_info}"
    instance_name=$(cut --fields 4 --delimiter=',' <<< "${service_info}")
    short_name=$(cut --fields 2 --delimiter=',' <<< "${service_info}")
    host=$(cut --fields 3 --delimiter=',' <<< "${service_info}")
    port=$(cut --fields 5 --delimiter=',' <<< "${service_info}")
    echo -n "${instance_name} ${short_name} "
    curl "${host}:${port}/monitoring/healthcheck"
    echo ""
done < <(sed 1d "${dev_csv}" | grep --invert-match --extended-regexp '^$')
#done < <(sed 1d "${dev_csv}" | awk -F , '{print $3 ":" $5 "/monitoring/healthcheck"}' | grep --invert-match --extended-regexp '^$')
echo

exit 0