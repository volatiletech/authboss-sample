{{define "pagetitle"}}Blogs - Index{{end}}
<h3>Dashboard</h3>

<a href="/blogs/new">New</a><br /><br />
<table style="border: 1px solid #000;">
	<thead>
		<tr>
			<td>Title</td>
			<td>Author</td>
			<td>Date</td>
		</tr>
	</thead>
	<tbody>
	{{range $i, $post := .posts}}
		<tr>
			<td><a href="/blogs/{{$i}}">{{$post.Title}}</td>
			<td>{{$post.AuthorID}}</td>
			<td>{{formatDate $post.Date}}</td>
		</tr>
	{{end}}
</table>