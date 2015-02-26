{{define "pagetitle"}}Blogs - {{.post.Title}}{{end}}
{{with .post}}
<h3>Dashboard</h3>
<table style="border: 1px solid #000;">
	<tbody>
		<tr>
			<td>Title</td><td>{{.Title}}</td>
		<tr></tr>
			<td>Author</td><td>{{.AuthorID}}</td>
		<tr></tr>
			<td>Date</td><td>{{.Date}}</td>
		</tr>
	</tbody>
</table>
<div>
	<p>{{.Content}}</p>
</div>
<br />
<a href="/blogs">Index</a> | <a href="/blogs/{{.ID}}/edit">Edit</a> | <a href="/blogs/{{.ID}}/destroy">Destroy</a>
{{end}}