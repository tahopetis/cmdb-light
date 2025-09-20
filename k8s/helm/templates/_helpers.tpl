{{/*
Expand the name of the chart.
*/}}
{{- define "cmdb-lite.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cmdb-lite.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cmdb-lite.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "cmdb-lite.labels" -}}
helm.sh/chart: {{ include "cmdb-lite.chart" . }}
{{ include "cmdb-lite.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "cmdb-lite.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cmdb-lite.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "cmdb-lite.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "cmdb-lite.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Backend labels
*/}}
{{- define "cmdb-lite.backend.labels" -}}
helm.sh/chart: {{ include "cmdb-lite.chart" . }}
{{ include "cmdb-lite.backend.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Backend selector labels
*/}}
{{- define "cmdb-lite.backend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cmdb-lite.name" . }}-backend
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Frontend labels
*/}}
{{- define "cmdb-lite.frontend.labels" -}}
helm.sh/chart: {{ include "cmdb-lite.chart" . }}
{{ include "cmdb-lite.frontend.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Frontend selector labels
*/}}
{{- define "cmdb-lite.frontend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cmdb-lite.name" . }}-frontend
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Adminer labels
*/}}
{{- define "cmdb-lite.adminer.labels" -}}
helm.sh/chart: {{ include "cmdb-lite.chart" . }}
{{ include "cmdb-lite.adminer.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Adminer selector labels
*/}}
{{- define "cmdb-lite.adminer.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cmdb-lite.name" . }}-adminer
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}