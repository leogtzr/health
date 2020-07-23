#!/bin/bash
# set -o xtrace
set -o nounset
set -o pipefail

readonly work_dir="$(dirname "$(readlink --canonicalize-existing "${0}")")"
readonly target_dir="${work_dir}/target"
readonly error_missing_jar_file=80
readonly mvn_default_goals="clean package"
readonly opts="-Dmaven.test.skip=true -DskipTests --threads 10"

mvn --quiet ${mvn_default_goals} ${opts} 2> /dev/null

exit
