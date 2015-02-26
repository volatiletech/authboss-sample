<html>
	<head>
		<title>{{template "pagetitle" .}}</title>
	</head>
<body>
	<div style="float: right;">
		{{if .loggedin}}<span>Hello, {{.username}}</span>{{else}}<span>User is not logged in.</span>{{end}}
	</div>
	{{template "yield" .}}
	{{template "authboss" .}}
</body>
</html>
{{define "pagetitle"}}{{end}}
{{define "authboss"}}{{end}}