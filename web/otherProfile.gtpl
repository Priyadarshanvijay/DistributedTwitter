<html>
	<head>
	<title></title>
	</head>
	<body>
    <h1>Home</h1>
		{{.Username}} Following: {{.FollowingNum}} Followers: {{.FollowersNum}}.
		<form action="/home" method="get">
			<input type="submit" value="Home">
		</form>
    Following List:
    {{if .Following}}
			{{range .Following}}
        {{.}} <br>
			{{end}}
		{{else}}
			<p>{{.Username}}'s following list is empty</p>
		{{end}}
    Followers List:
    {{if .Followers}}
			{{range .Followers}}
        {{.}} <br>
			{{end}}
		{{else}}
			<p>{{.Username}}'s followers list is empty</p>
		{{end}}
		<div style="width:100%; height:10%">
		<h3> Feed </h3>
		{{if .Posts}}
			{{range .Posts}}
				<div style="border: thin solid black">
				<h3 style="display: inline-block;">{{.author}}</h3> <p style="display: inline-block;">{{.createdAt}}</p>
				<p>{{.content}}</p>
				</div>
			{{end}}
		{{else}}
			<p>{{.Username}}'s feed is empty</p>
		{{end}}
	</body>
</html>