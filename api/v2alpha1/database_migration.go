// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
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

package v2alpha1

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
)

func (in *HorusecPlatform) GetDatabaseUserSecretKeyRef() *corev1.SecretKeySelector {
	secretKeyRef := in.Spec.Global.Database.User.SecretKeyRef
	secretKey := secretKeyRef.Key
	if secretKey == "" {
		secretKey = "username"
	}
	secretName := secretKeyRef.Name
	if secretName == "" {
		secretName = "horusec-database"
	}
	return &corev1.SecretKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
		Key:                  secretKey,
	}
}

func (in *HorusecPlatform) GetDatabasePasswordSecretKeyRef() *corev1.SecretKeySelector {
	secretKeyRef := in.Spec.Global.Database.Password.SecretKeyRef
	secretKey := secretKeyRef.Key
	if secretKey == "" {
		secretKey = "password"
	}
	secretName := secretKeyRef.Name
	if secretName == "" {
		secretName = "horusec-database"
	}
	return &corev1.SecretKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
		Key:                  secretKey,
	}
}

func (in *HorusecPlatform) GetDatabaseMigrationImage() string {
	registry := in.GetDatabaseMigrationImageRegistry()
	repository := in.GetGlobalDatabaseMigrationImageRepository()
	tag := in.GetGlobalDatabaseMigrationImageTag()
	if registry != "" {
		return fmt.Sprintf("%v/%v:%v", registry, repository, tag)
	}
	return fmt.Sprintf("%v:%v", repository, tag)
}

func (in *HorusecPlatform) GetDatabaseMigrationImageRegistry() string {
	registry := in.Spec.Global.Database.Migration.Image.Registry
	if registry != "" {
		return registry
	}
	return "docker.io/horuszup"
}

func (in *HorusecPlatform) GetGlobalDatabaseMigrationImageRepository() string {
	repository := in.Spec.Global.Database.Migration.Image.Repository
	if repository != "" {
		return repository
	}
	return "horusec-migrations"
}

func (in *HorusecPlatform) GetGlobalDatabaseMigrationImageTag() string {
	tag := in.Spec.Global.Database.Migration.Image.Tag
	if tag != "" {
		return tag
	}
	return "v2.12.1"
}
