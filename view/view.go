package view

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/selassje/githubProfiles/controller"
	"github.com/zserge/webview"
)

const (
	windowWidth  = 260
	windowHeight = 260
)

var indexHTML = `
<!doctype html>
<html>
	<head>
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
	</head>
	<body>
		<button onclick="external.invoke('searchUser:'+document.getElementById('searched-user').value)">
			Search User
		</button>
		<input id="searched-user" type="text" value="selassje"></input>
		<p><div id = "User">User:</div><img id="avatar" src="F:\Contests\go\githubProfiles\avatar.jpg" width="42" height="42" style="float:right"></p>
		<p id = "Repos Count">Repos Count:</p>
		<p><div id = "Followers">Followers:</div></p>
		<select id="followers-list" size="5" name ="FollowersList" hidden></select><select id="top-repos" size="4" name ="TopReposList" hidden></select>  
	</body>
</html>
`

func updateField(w webview.WebView, e string, v string) {
	jsCode := `document.getElementById("` + e + `").innerHTML ="` + e + `: ` + v + `"`
	w.Eval(jsCode)
}

func updateAvatar(w webview.WebView, image []byte) {
	jsCode := `var encoder = new JPEGEncoder(9);
	var jpgFile = encoder.encode(` + string(image) + `, 9);
	document.getElementById("avatar").src = jpgFile;	`
	w.Eval(jsCode)
}

func updateFollowers(w webview.WebView, followers []string) {
	var jsCode string
	if len(followers) == 0 {
		jsCode = `document.getElementById("followers-list").setAttribute("hidden");`
	} else {
		jsCode = `document.getElementById("followers-list").removeAttribute("hidden");`
		for _, follower := range followers {
			jsCode += ` var option = document.createElement("option");
			            option.innerHTML = "` + follower + `";
						document.getElementById("followers-list").appendChild(option);`
		}
	}
	w.Eval(jsCode)
}

func handleRPC(w webview.WebView, data string) {
	switch {
	case strings.HasPrefix(data, "searchUser:"):
		userName := strings.TrimPrefix(data, "searchUser:")
		user, err := controller.GetUserInfo(userName)
		var userStr, reposCountStr, followersCountStr string
		//var avatar []byte
		var followers []string
		if err == nil {
			//fmt.Println(user)
			userStr = user.Username
			reposCountStr = strconv.Itoa(user.ReposCount)
			followersCountStr = strconv.Itoa(len(user.Followers))
			followers = user.Followers
			//avatar   = user.Avatar
		} else {
			userStr = "User " + userName + " not found"
		}
		updateField(w, "Repos Count", reposCountStr)
		updateField(w, "User", userStr)
		updateField(w, "Followers",followersCountStr)
		updateFollowers(w, followers)
		//updateAvatar(w, avatar)
	}
}

func RunGui() {
	w := webview.New(webview.Settings{
		Width:                  windowWidth,
		Height:                 windowHeight,
		Title:                  "GitHubProfiles",
		Resizable:              false,
		URL:                    "data:text/html," + url.PathEscape(indexHTML),
		ExternalInvokeCallback: handleRPC,
	})
	w.SetColor(255, 255, 255, 255)
	defer w.Exit()
	w.Run()
}
