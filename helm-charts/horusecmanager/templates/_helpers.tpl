{{/* vim: set filetype=mustache: */}}
{{/*
Return the proper Horusec Auth image name
*/}}
{{- define "auth.image" -}}
{{- $registryName := .Values.components.auth.container.image.registry -}}
{{- $repositoryName := .Values.components.auth.container.image.repository -}}
{{- $tag := .Values.components.auth.container.image.tag | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}

{{/*
Return the proper Horusec Manager image name
*/}}
{{- define "manager.image" -}}
{{- $registryName := .Values.components.manager.container.image.registry -}}
{{- $repositoryName := .Values.components.manager.container.image.repository -}}
{{- $tag := .Values.components.manager.container.image.tag | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}

{{/*
Return the proper Horusec Account image name
*/}}
{{- define "account.image" -}}
{{- $registryName := .Values.components.account.container.image.registry -}}
{{- $repositoryName := .Values.components.account.container.image.repository -}}
{{- $tag := .Values.components.account.container.image.tag | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}

{{/*
Return the proper Horusec API image name
*/}}
{{- define "api.image" -}}
{{- $registryName := .Values.components.api.container.image.registry -}}
{{- $repositoryName := .Values.components.api.container.image.repository -}}
{{- $tag := .Values.components.api.container.image.tag | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}

{{/*
Return the proper Horusec Analytic image name
*/}}
{{- define "analytic.image" -}}
{{- $registryName := .Values.components.analytic.container.image.registry -}}
{{- $repositoryName := .Values.components.analytic.container.image.repository -}}
{{- $tag := .Values.components.analytic.container.image.tag | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}

{{/*
Return the proper Horusec Messages image name
*/}}
{{- define "messages.image" -}}
{{- $registryName := .Values.components.messages.container.image.registry -}}
{{- $repositoryName := .Values.components.messages.container.image.repository -}}
{{- $tag := .Values.components.messages.container.image.tag | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}

{{/*
Return the proper Horusec Webhook image name
*/}}
{{- define "webhook.image" -}}
{{- $registryName := .Values.components.webhook.container.image.registry -}}
{{- $repositoryName := .Values.components.webhook.container.image.repository -}}
{{- $tag := .Values.components.webhook.container.image.tag | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}

{{/*
Return the proper Horusec Account Image Registry Secret Names
*/}}
{{- define "account.imagePullSecrets" -}}
{{- if .Values.components.account.container.image.pullSecrets }}
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
{{- if .Values.components.analytic.container.image.pullSecrets }}
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
{{- if .Values.components.api.container.image.pullSecrets }}
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
{{- if .Values.components.manager.container.image.pullSecrets }}
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
{{- if .Values.components.auth.container.image.pullSecrets }}
imagePullSecrets:
    {{- range .Values.components.auth.container.image.pullSecrets }}
    - name: {{ . }}
    {{- end }}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec Webhook Image Registry Secret Names
*/}}
{{- define "webhook.imagePullSecrets" -}}
{{- if .Values.components.webhook.container.image.pullSecrets }}
imagePullSecrets:
    {{- range .Values.components.webhook.container.image.pullSecrets }}
    - name: {{ . }}
    {{- end }}
{{- end -}}
{{- end -}}

{{/*
Return the proper Horusec Messages Image Registry Secret Names
*/}}
{{- define "messages.imagePullSecrets" -}}
{{- if .Values.components.messages.container.image.pullSecrets }}
imagePullSecrets:
    {{- range .Values.components.messages.container.image.pullSecrets }}
    - name: {{ . }}
    {{- end }}
{{- end -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for deployment.
*/}}
{{- define "deployment.apiVersion" -}}
{{- if semverCompare "<1.14-0" .Capabilities.KubeVersion.Version -}}
{{- print "extensions/v1beta1" -}}
{{- else -}}
{{- print "apps/v1" -}}
{{- end -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for Ingress.
*/}}
{{- define "ingress.apiVersion" -}}
{{- if semverCompare "<1.14-0" .Capabilities.KubeVersion.Version -}}
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

{{/*
Print "true" if the API pathType field is supported.
*/}}
{{- define "ingress.supportsPathType" -}}
{{- if not .Capabilities -}}
{{- print "false" -}}
{{- else if semverCompare "<1.18-0" .Capabilities.KubeVersion.Version -}}
{{- print "false" -}}
{{- else -}}
{{- print "true" -}}
{{- end -}}
{{- end -}}

{{/*
If enabled, return Ingress Rules.
*/}}
{{- define "ingress.rules" -}}
{{- $components := list -}}
{{- range $_, $component := .Values.components -}}
    {{- if $component.ingress -}}
        {{- $components = append $components $component -}}
    {{- end -}}
{{- end -}}

{{- $hosts := dict -}}
{{- range $_, $component := $components -}}
    {{- if $component.ingress -}}
    {{ if not (hasKey $hosts $component.ingress.host) }}
        {{- $ingresses := list -}}
        {{- range $_, $otherComponent := $components -}}
            {{- if eq $component.ingress.host $otherComponent.ingress.host -}}
                {{- $ingresses = append $ingresses $otherComponent -}}
            {{- end -}}
        {{- end -}}
        {{- $_ := set $hosts $component.ingress.host (compact $ingresses) -}}
    {{- end -}}
    {{- end -}}
{{- end -}}

rules:
{{- range $host, $components := $hosts }}
  - host: {{ $host }}
    http:
      paths:
        {{- range $component := $components }}
        - path: {{ $component.ingress.path }}
          {{- if eq "true" (include "ingress.supportsPathType" .) }}
          pathType: Prefix
          {{- end }}
          backend:
            serviceName: {{ $component.name }}
            servicePort: {{ $component.port.http }}
        {{- end }}
{{- end -}}
{{- end -}}

{{/*
If enabled, return SSL/TLS Ingress YAML configuration.
*/}}
{{- define "ingress.tls" -}}
{{- $ingresses := list -}}
{{- range $_, $component := .Values.components -}}
    {{- if $component.ingress -}}
        {{- $ingresses = append $ingresses $component.ingress -}}
    {{- end -}}
{{- end -}}

{{- $secrets := dict -}}
{{- range $_, $ingress := $ingresses -}}
    {{- if and (eq "https" $ingress.scheme) $ingress.tls -}}
        {{ if not (hasKey $secrets $ingress.tls.secretName) }}
            {{- $hosts := list -}}
            {{- range $_, $otherIngress := $ingresses -}}
                {{- if $otherIngress.tls -}}
                {{- if eq $ingress.tls.secretName $otherIngress.tls.secretName -}}
                    {{- $hosts = append $hosts $otherIngress.host -}}
                {{- end -}}
                {{- end -}}
            {{- end -}}
            {{- $_ := set $secrets $ingress.tls.secretName (compact $hosts) -}}
        {{- end -}}
    {{- end -}}
{{- end -}}

{{- if $secrets -}}
tls:
  {{- range $secret, $hosts := $secrets }}
  {{- if $secret }}
  - hosts:
      {{- range $host := $hosts }}
      - {{ $host }}
      {{- end }}
    secretName: {{ $secret }}
  {{- end -}}
  {{- end -}}
{{- end -}}
{{- end -}}
