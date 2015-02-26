<form action="/blogs{{with .id}}/{{.}}{{end}}" method="post">
	<label for="title">Title:</label>
	<input id="title" name="title" type="text" value="{{with .post}}{{.Title}}{{end}}"></input><br /><br />
	<label for="content">Content:</label><br />
	<textarea id="content" name="content" cols="100" rows="20">{{with .post}}{{.Content}}{{end}}</textarea><br /><br />
	{{with .id}}
		<input type="submit" value="Edit" />
	{{else}}
		<input type="submit" value="Create" />
	{{end}}

	<input type="hidden" name="crsf_token" value="{{.csrf_token}}" />
</form>