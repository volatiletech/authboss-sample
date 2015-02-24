<html>
	<head>
		<title>{{define "title"}}{{end}}</title>
	</head>
<body>
	<div style="float: right;">
		{{if .IsLoggedIn}}<span>Hello, {{.username}}</span>{{else}}<span>User is not logged in.</span>{{end}}
	</div>
	{{yield}}
	{{define "authboss"}}{{end}}
	{{template "authboss" .}}
</body>
</html>