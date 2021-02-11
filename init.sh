#!/bin/bash

set -e
set -x

operator-sdk init --plugins helm --domain horusec.io --group install --version v1alpha1 --kind HorusecManager
