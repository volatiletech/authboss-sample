<html>
	<head>
		<title>{{template "pagetitle" .}}</title>
	</head>
<body>
	{{with .flash_success}}<div style="color: green;"><strong>{{.}}</strong></div>{{end}}
	{{with .flash_error}}<div style="color: red;"><strong>{{.}}</strong></div>{{end}}
	<div style="float: right;">
		{{if .loggedin}}<a href="auth/logout">Logout</a>{{else}}<a href="auth/login">Login</a>{{end}}
		<a href="auth/register">Register</a>
		{{if .loggedin}}<span>Hello, {{.current_user_name}}</span>{{else}}<span>User is not logged in.</span>{{end}}
	</div>
	{{template "yield" .}}
	{{template "authboss" .}}
</body>
</html>
{{define "pagetitle"}}{{end}}
{{define "yield"}}{{end}}
{{define "authboss"}}{{end}}