<html>
<body>
<h3>Dashboard</h3>
{{if .IsLoggedIn}}
<p>Username: {{.Username}}</p>
{{else}}
boom headshot
{{end}}
</body>
</html>