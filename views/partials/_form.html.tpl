<div class="row">
	<div class="col-md-offset-1 col-md-10">
		<form action="/blogs{{with .post.ID}}/{{.}}/edit{{else}}/new{{end}}" method="post">
			<div class="form-group">
				<label for="title">Title</label>
				<input class="form-control" name="title" type="text" value="{{with .post}}{{.Title}}{{end}}"></input>
			</div>
			<div class="form-group">
				<label for="content">Content</label><br />
				<textarea class="form-control" name="content" cols="100" rows="5">{{with .post}}{{.Content}}{{end}}</textarea>
			</div>
			<input type="hidden" name="csrf_token" value="{{.csrf_token}}" />
			<div class="row text-right">
				{{with .post.ID}}
				<button class="btn btn-success" type="submit">Edit</button>
				{{else}}
				<button class="btn btn-success" type="submit">Create</button>
				{{end}}
				<a class="btn btn-link" href="/">Cancel</button>
			</form>
		</div>
	</div>
</div>	
