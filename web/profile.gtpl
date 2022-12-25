<html>
	<head>
	<title></title>
	</head>
	<body>
    <h1>Home</h1>
		Hello {{.Username}} Following: {{.FollowingNum}} Followers: {{.FollowersNum}}.
		<form action="/logout" method="post">
			<input type="submit" value="Logout">
		</form>
		<form action="/home" method="get">
			<input type="submit" value="Home">
		</form>
    Following List:
    {{if .Following}}
			{{range .Following}}
        {{.}} 
        <form action="/unFollowUser" method="post">
          <input hidden type="text" name="username" value={{.}}>
          <input type="submit" value="Unfollow">
        </form>
			{{end}}
		{{else}}
			<p>Your following list is empty</p>
		{{end}}
    Followers List:
    {{if .Followers}}
			{{range .Followers}}
        {{.}} 
			{{end}}
		{{else}}
			<p>Your followers list is empty</p>
		{{end}}
		<div style="width:100%; height:10%">
		<h3> My Posts </h3>
		{{if .Posts}}
			{{range .Posts}}
				<div style="border: thin solid black">
				<h3 style="display: inline-block;">{{.author}}</h3> <p style="display: inline-block;">{{.createdAt}}</p>
				<p>{{.content}}</p>
        <form action="/deletePost" method="post">
          <input hidden type="text" name="postId" value={{.postId}}>
          <input type="submit" value="Delete Post">
        </form>
				</div>
			{{end}}
		{{else}}
			<p>Your posts are empty</p>
		{{end}}
	</body>
</html>