{{/* vim: set filetype=mustache: */}}
{{/*
Return the proper Horusec Auth image name
*/}}
{{- define "auth.image" -}}
{{- $registryName := .Values.components.auth.container.image.registry -}}
{{- $repositoryName := .Values.components.auth.container.image.repository -}}
{{- $tag := .Values.components.auth.container.image.tag | toString -}}
{{- if .Values.global }}
    {{- if .Values.global.common.container.image.registry }}
        {{- printf "%s/%s:%s" .Values.global.common.container.image.registry $repositoryName $tag -}}
    {{- else -}}
        {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
    {{- end -}}
{{- else -}}
    {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec Manager image name
*/}}
{{- define "manager.image" -}}
{{- $registryName := .Values.components.manager.container.image.registry -}}
{{- $repositoryName := .Values.components.manager.container.image.repository -}}
{{- $tag := .Values.components.manager.container.image.tag | toString -}}
{{- if .Values.global }}
    {{- if .Values.global.common.container.image.registry }}
        {{- printf "%s/%s:%s" .Values.global.common.container.image.registry $repositoryName $tag -}}
    {{- else -}}
        {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
    {{- end -}}
{{- else -}}
    {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec Account image name
*/}}
{{- define "account.image" -}}
{{- $registryName := .Values.components.account.container.image.registry -}}
{{- $repositoryName := .Values.components.account.container.image.repository -}}
{{- $tag := .Values.components.account.container.image.tag | toString -}}
{{- if .Values.global }}
    {{- if .Values.global.common.container.image.registry }}
        {{- printf "%s/%s:%s" .Values.global.common.container.image.registry $repositoryName $tag -}}
    {{- else -}}
        {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
    {{- end -}}
{{- else -}}
    {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec API image name
*/}}
{{- define "api.image" -}}
{{- $registryName := .Values.components.api.container.image.registry -}}
{{- $repositoryName := .Values.components.api.container.image.repository -}}
{{- $tag := .Values.components.api.container.image.tag | toString -}}
{{- if .Values.global }}
    {{- if .Values.global.common.container.image.registry }}
        {{- printf "%s/%s:%s" .Values.global.common.container.image.registry $repositoryName $tag -}}
    {{- else -}}
        {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
    {{- end -}}
{{- else -}}
    {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec Analytic image name
*/}}
{{- define "analytic.image" -}}
{{- $registryName := .Values.components.analytic.container.image.registry -}}
{{- $repositoryName := .Values.components.analytic.container.image.repository -}}
{{- $tag := .Values.components.analytic.container.image.tag | toString -}}
{{- if .Values.global }}
    {{- if .Values.global.common.container.image.registry }}
        {{- printf "%s/%s:%s" .Values.global.common.container.image.registry $repositoryName $tag -}}
    {{- else -}}
        {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
    {{- end -}}
{{- else -}}
    {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec Account Image Registry Secret Names
*/}}
{{- define "account.imagePullSecrets" -}}
{{- if .Values.global }}
{{- if .Values.global.common.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.global.common.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- else if .Values.components.account.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.account.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if .Values.components.account.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.account.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec Analytic Image Registry Secret Names
*/}}
{{- define "analytic.imagePullSecrets" -}}
{{- if .Values.global }}
{{- if .Values.global.common.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.global.common.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- else if .Values.components.analytic.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.analytic.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if .Values.components.analytic.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.analytic.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec API Image Registry Secret Names
*/}}
{{- define "api.imagePullSecrets" -}}
{{- if .Values.global }}
{{- if .Values.global.common.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.global.common.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- else if .Values.components.api.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.api.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if .Values.components.api.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.api.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec Manager Image Registry Secret Names
*/}}
{{- define "manager.imagePullSecrets" -}}
{{- if .Values.global }}
{{- if .Values.global.common.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.global.common.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- else if .Values.components.manager.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.manager.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if .Values.components.manager.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.manager.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec Auth Image Registry Secret Names
*/}}
{{- define "auth.imagePullSecrets" -}}
{{- if .Values.global }}
{{- if .Values.global.common.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.global.common.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- else if .Values.components.auth.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.auth.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if .Values.components.auth.container.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.auth.container.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for deployment.
*/}}
{{- define "deployment.apiVersion" -}}
{{- if semverCompare "<1.14-0" .Capabilities.KubeVersion.GitVersion -}}
{{- print "extensions/v1beta1" -}}
{{- else -}}
{{- print "apps/v1" -}}
{{- end -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for Ingress.
*/}}
{{- define "ingress.apiVersion" -}}
{{- if semverCompare "<1.14-0" .Capabilities.KubeVersion.GitVersion -}}
{{- print "extensions/v1beta1" -}}
{{- else -}}
{{- print "networking.k8s.io/v1beta1" -}}
{{- end -}}
{{- end -}}

{{/*
True if Ingress is enabled for any of the components.
*/}}
{{- define "ingress.enabled" -}}
{{- if or .Values.components.auth.ingress.enabled .Values.components.manager.ingress.enabled .Values.components.api.ingress.enabled .Values.components.analytic.ingress.enabled .Values.components.account.ingress.enabled }}
    {{- true -}}
{{- end -}}
{{- end -}}
