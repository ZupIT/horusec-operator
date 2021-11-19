//go:build mage
// +build mage

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

func UpdateVersioningFiles() error {
	if err := sh.RunV("npm", "install", "-g", "json"); err != nil {
		return err
	}

	for _, valueToReplace := range replaceValues() {
		if err := replacePlatformVersion(valueToReplace); err != nil {
			return err
		}
	}

	return updateOperatorVersion()
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

func updateOperatorVersion() error {
	for _, path := range sedValues() {
		sed := fmt.Sprintf("s/%s/%s/g", getActualVersion(), getReleaseVersion())
		if err := sh.Run("sed", "-i", sed, path); err != nil {

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
