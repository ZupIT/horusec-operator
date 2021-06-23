#!/bin/bash
# Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


UPDATE_TYPE=$1
IS_TO_UPDATE_LATEST=$2
DIRECTORY=$(pwd)
IMAGE_NAME="horuszup/horusec-operator"

installSemver () {
    mkdir -p bin
    curl -fsSL -o "$DIRECTORY/bin/install-semver.sh" https://raw.githubusercontent.com/ZupIT/horusec-devkit/main/scripts/install-semver.sh
    chmod +x "$DIRECTORY/bin/install-semver.sh"
    "$DIRECTORY/bin/install-semver.sh"
    if ! semver &> /dev/null
    then
        exit 1
    fi
}

validateUpdateType () {
    case "$UPDATE_TYPE" in
        "alpha") # Used to update an bugfix or an new feature in develop branch
            echo "Update type selected is alpha" ;;
        "rc") # Used when you finish development and start testing in the test environment and in develop branch
            echo "Update type selected is rc(release-candidate)" ;;
        "release") # Used when an correction was applied in master branch
            echo "Update type selected is release" ;;
        "minor") # Used when an new feature is enable in production environment and in master branch
            echo "Update type selected is minor" ;;
        "major") # Used when an big refactor is necessary to breaking changes in master branch
            echo "Update type selected is major" ;;
        *)
            echo "Param Update type is invalid, please use the examples bellow allowed and try again!"
            echo "Params Update type allowed: alpha, rc, release, minor, major"
            exit 1;;
    esac
}

validateIsToUpdateLatest () {
    if [[ "$IS_TO_UPDATE_LATEST" != "true" && "$IS_TO_UPDATE_LATEST" != "false" ]]
    then
        echo "Param \"is to update latest\" is invalid, please use the examples bellow allowed and try again!"
        echo "Param \"is to update latest\" allowed: true, false"
        exit 1
    fi
}

updateVersion () {
    resetAlphaRcToMaster

    semver up "$UPDATE_TYPE"

    LATEST_VERSION=""
    if [[ "$UPDATE_TYPE" == "alpha" || "$UPDATE_TYPE" == "rc" ]]
    then
        LATEST_VERSION=$(semver get "$UPDATE_TYPE")
    else
        LATEST_VERSION=$(semver get release)
    fi

    if [ "$IS_TO_UPDATE_LATEST" == "true" ]
    then
        if ! docker build -t "horuszup/horusec-operator:latest" -f "$DIRECTORY/deployments/dockerfiles/Dockerfile" .; then
            exit 1
        fi
        if ! docker push "horuszup/horusec-operator:latest"; then
            exit 1
        fi
    fi
    if ! docker build -t "horuszup/horusec-operator:$LATEST_VERSION" -f "$DIRECTORY/deployments/dockerfiles/Dockerfile" .; then
        exit 1
    fi
    if ! docker push "horuszup/horusec-operator:$LATEST_VERSION"; then
        exit 1
    fi
}

resetAlphaRcToMaster () {
    if [[ "$UPDATE_TYPE" == "release" || "$UPDATE_TYPE" == "minor" || "$UPDATE_TYPE" == "major" ]]
    then
        alpha_version=$(semver get alpha)
        rc_version=$(semver get rc)
        if [[ "${alpha_version: -2}" != ".0" || "${rc_version: -2}" != ".0" ]]
        then
            echo "Alpha or Release-Candidate found, reseting version to:"
            semver up release
        fi
    fi
}

validateUpdateType

validateIsToUpdateLatest

installSemver

updateVersion
