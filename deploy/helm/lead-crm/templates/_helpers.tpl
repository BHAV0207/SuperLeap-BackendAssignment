{{- define "lead-crm.labels" -}}
app.kubernetes.io/part-of: lead-crm
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
{{- end }}

{{- define "lead-crm.appSelector" -}}
app.kubernetes.io/name: lead-crm
app.kubernetes.io/component: api
{{- end }}

{{- define "lead-crm.postgresSelector" -}}
app.kubernetes.io/name: postgres
app.kubernetes.io/component: database
{{- end }}

{{- define "lead-crm.redisSelector" -}}
app.kubernetes.io/name: redis
app.kubernetes.io/component: cache
{{- end }}

{{- define "lead-crm.image" -}}
{{- $tag := required "app.image.tag is required (use the git commit sha)" .Values.app.image.tag -}}
{{ .Values.app.image.repository }}:{{ $tag }}
{{- end }}

{{/* hardened container settings shared by all pods */}}
{{- define "lead-crm.containerSecurity" -}}
allowPrivilegeEscalation: false
readOnlyRootFilesystem: true
capabilities:
  drop: ["ALL"]
{{- end }}
