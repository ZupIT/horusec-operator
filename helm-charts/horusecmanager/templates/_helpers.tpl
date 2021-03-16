{{/* vim: set filetype=mustache: */}}
{{/*
Return the proper Horusec Auth image name
*/}}
{{- define "auth.image" -}}
{{- $registryName := .Values.components.auth.image.registry -}}
{{- $repositoryName := .Values.components.auth.image.repository -}}
{{- $tag := .Values.components.auth.image.tag | toString -}}
{{- if .Values.global }}
    {{- if .Values.global.image.registry }}
        {{- printf "%s/%s:%s" .Values.global.image.registry $repositoryName $tag -}}
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
{{- $registryName := .Values.components.manager.image.registry -}}
{{- $repositoryName := .Values.components.manager.image.repository -}}
{{- $tag := .Values.components.manager.image.tag | toString -}}
{{- if .Values.global }}
    {{- if .Values.global.image.registry }}
        {{- printf "%s/%s:%s" .Values.global.image.registry $repositoryName $tag -}}
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
{{- $registryName := .Values.components.account.image.registry -}}
{{- $repositoryName := .Values.components.account.image.repository -}}
{{- $tag := .Values.components.account.image.tag | toString -}}
{{- if .Values.global }}
    {{- if .Values.global.image.registry }}
        {{- printf "%s/%s:%s" .Values.global.image.registry $repositoryName $tag -}}
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
{{- $registryName := .Values.components.api.image.registry -}}
{{- $repositoryName := .Values.components.api.image.repository -}}
{{- $tag := .Values.components.api.image.tag | toString -}}
{{- if .Values.global }}
    {{- if .Values.global.image.registry }}
        {{- printf "%s/%s:%s" .Values.global.image.registry $repositoryName $tag -}}
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
{{- $registryName := .Values.components.analytic.image.registry -}}
{{- $repositoryName := .Values.components.analytic.image.repository -}}
{{- $tag := .Values.components.analytic.image.tag | toString -}}
{{- if .Values.global }}
    {{- if .Values.global.image.registry }}
        {{- printf "%s/%s:%s" .Values.global.image.registry $repositoryName $tag -}}
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
{{- if .Values.global.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.global.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- else if .Values.components.account.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.account.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if .Values.components.account.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.account.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec Analytic Image Registry Secret Names
*/}}
{{- define "analytic.imagePullSecrets" -}}
{{- if .Values.global }}
{{- if .Values.global.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.global.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- else if .Values.components.analytic.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.analytic.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if .Values.components.analytic.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.analytic.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec API Image Registry Secret Names
*/}}
{{- define "api.imagePullSecrets" -}}
{{- if .Values.global }}
{{- if .Values.global.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.global.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- else if .Values.components.api.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.api.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if .Values.components.api.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.api.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec Manager Image Registry Secret Names
*/}}
{{- define "manager.imagePullSecrets" -}}
{{- if .Values.global }}
{{- if .Values.global.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.global.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- else if .Values.components.manager.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.manager.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if .Values.components.manager.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.manager.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec Auth Image Registry Secret Names
*/}}
{{- define "auth.imagePullSecrets" -}}
{{- if .Values.global }}
{{- if .Values.global.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.global.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- else if .Values.components.auth.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.auth.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if .Values.components.auth.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.components.auth.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}
