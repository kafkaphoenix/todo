package main

import "github.com/calvinmclean/babyapi/html"

const (
	allTODOs         html.Template = "allTODOs"
	allTODOsTemplate string        = `<!doctype html>
<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>TODOs</title>
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/uikit@3.17.11/dist/css/uikit.min.css" />
		<script src="https://unpkg.com/htmx.org@1.9.8"></script>
		<script src="https://unpkg.com/htmx.org/dist/ext/sse.js"></script>
	</head>

	<style>
	tr.htmx-swapping td {
		opacity: 0;
		transition: opacity 1s ease-out;
	}
	</style>

	<body>
		<table class="uk-table uk-table-divider uk-margin-left uk-margin-right">
			<colgroup>
				<col>
				<col>
				<col style="width: 300px;">
			</colgroup>

			<thead>
				<tr>
					<th>Title</th>
					<th>Description</th>
					<th></th>
				</tr>
			</thead>

			<tbody hx-ext="sse" sse-connect="/todos/listen" sse-swap="newTODO" hx-swap="beforeend">
				<form hx-post="/todos" hx-swap="none" hx-on::after-request="this.reset()">
					<td>
						<input class="uk-input" name="Title" type="text">
					</td>
					<td>
						<input class="uk-input" name="Description" type="text">
					</td>
					<td>
						<button type="submit" class="uk-button uk-button-primary">Add TODO</button>
					</td>
				</form>

				{{ range . }}
				{{ template "todoRow" . }}
				{{ end }}
			</tbody>
		</table>
	</body>
</html>`

	todoRow         html.Template = "todoRow"
	todoRowTemplate string        = `<tr hx-target="this" hx-swap="outerHTML">
	<td>{{ .Title }}</td>
	<td>{{ .Description }}</td>
	<td>
		{{- $color := "primary" }}
		{{- $disabled := "" }}
		{{- if .Completed }}
			{{- $color = "secondary" }}
			{{- $disabled = "disabled" }}
		{{- end -}}

		<button class="uk-button uk-button-{{ $color }}"
			hx-put="/todos/{{ .ID }}"
			hx-headers='{"Accept": "text/html"}'
			hx-include="this"
			{{ $disabled }}>

			<input type="hidden" name="Completed" value="true">
			<input type="hidden" name="Title" value="{{ .Title }}">
			<input type="hidden" name="Description" value="{{ .Description }}">
			<input type="hidden" name="ID" value="{{ .ID }}">
			Complete
		</button>

		<button class="uk-button uk-button-danger" hx-delete="/todos/{{ .ID }}" hx-swap="swap:1s">
			Delete
		</button>
	</td>
</tr>`
)
