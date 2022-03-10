// Copyright 2021 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
	// mage:import
	_ "github.com/ZupIT/horusec-devkit/pkg/utils/mageutils"
)

const (
	replacePathAnalytic            = "this.components.analytic.container.image.tag=\"%s\""
	replacePathApi                 = "this.components.api.container.image.tag=\"%s\""
	replacePathAuth                = "this.components.auth.container.image.tag=\"%s\""
	replacePathCore                = "this.components.core.container.image.tag=\"%s\""
	replacePathManager             = "this.components.manager.container.image.tag=\"%s\""
	replacePathMessages            = "this.components.messages.container.image.tag=\"%s\""
	replacePathVulnerability       = "this.components.vulnerability.container.image.tag=\"%s\""
	replacePathWebhook             = "this.components.webhook.container.image.tag=\"%s\""
	replacePathDatabaseMigration   = "this.global.database.migration.image.tag=\"%s\""
	replacePathAnalyticDatabase    = "this.components.analytic.database.migration.image.tag=\"%s\""
	pathToReplaceSeedKustomization = "config/manager/kustomization.yaml"
	pathToReplaceSeedReadme        = "README.md"
	defaultJsonPath                = "api/v2alpha1/horusec_platform_defaults.json"
)

const (
	envPlatformVersion = "HORUSEC_PLATFORM_VERSION"
	envActualVersion   = "HORUSEC_ACTUAL_VERSION"
	envReleaseVersion  = "HORUSEC_RELEASE_VERSION"
)

// UpdateVersioningFiles update project version in all files
func UpdateVersioningFiles() error {
	if err := sh.RunV("npm", "install", "-g", "json"); err != nil {
		return err
	}

	for _, valueToReplace := range replaceValues() {
		if err := replacePlatformVersion(valueToReplace); err != nil {
			return err
		}
	}

	return updateOperatorVersions(getActualVersion(), getReleaseVersion())
}

func replacePlatformVersion(valueToReplace string) error {
	valueReplaced := fmt.Sprintf(valueToReplace, getPlatformVersion())

	return sh.RunV("json", "-I", "-f", defaultJsonPath, "-e", valueReplaced)
}

func replaceValues() []string {
	return []string{
		replacePathAnalytic,
		replacePathApi,
		replacePathAuth,
		replacePathCore,
		replacePathManager,
		replacePathMessages,
		replacePathVulnerability,
		replacePathWebhook,
		replacePathDatabaseMigration,
		replacePathAnalyticDatabase,
	}
}

func sedValues() []string {
	return []string{
		pathToReplaceSeedKustomization,
		pathToReplaceSeedReadme,
	}
}

func updateOperatorVersions(old, new string) error {
	for _, path := range sedValues() {
		sed := fmt.Sprintf("s/%s/%s/g", old, new)
		if err := sh.Run("sed", "-i", sed, path); err != nil {
			return err
		}
	}

	return nil
}

func getActualVersion() string {
	return os.Getenv(envActualVersion)
}

func getReleaseVersion() string {
	return os.Getenv(envReleaseVersion)
}

func getPlatformVersion() string {
	return os.Getenv(envPlatformVersion)
}

func SignImage(tag string) error {
	imageWithTag := fmt.Sprintf("horuszup/horusec-operator:%s", tag)

	if err := sh.Run("cosign", "sign", "-key",
		"$COSIGN_KEY_LOCATION", imageWithTag); err != nil {
		return err
	}

	return nil
}
