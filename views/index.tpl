<h3>Dashboard</h3>
<table style="border: 1px solid #000;">
	<thead>
		<tr>
			<td>Title</td>
			<td>Author</td>
			<td>Date</td>
		</tr>
	</thead>
	<tbody>
	{{range .posts}}
		<tr>
			<td>{{.Name}}</td>
			<td>{{.Author}}</td>
			<td>{{formatDate .Date}}</td>
		</tr>
	{{end}}
</table>