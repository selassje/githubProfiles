package view

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/selassje/githubProfiles/controller"
	"github.com/zserge/webview"
)

const (
	windowWidth  = 270
	windowHeight = 250
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
		<input id="searched-user" type="text" value="tom" style="width: 140px;"></input>
		<p><div id = "User">User:</div><img id="avatar" width="42" height="42" style="float:right" src=""></p>
		<p id = "Repos Count">Repos Count:</p>
		<div style="display:inline-block">
			<p><div id = "Followers">Followers:</div></p>
			<select id="followers-list" size="5" name ="FollowersList" style="width: 120px"></select>
		</div>
		<div style="display:inline-block; float:right;">
			<p>Top Repos:</p>
			<select id="top-repos" size="5" name ="TopReposList" style="width: 120px"></select>
		</div> 
	</body>
</html>
`

func updateField(w webview.WebView, e string, v string) {
	jsCode := `document.getElementById("` + e + `").innerHTML ="` + e + `: ` + v + `"`
	w.Eval(jsCode)
}

func updateAvatar(w webview.WebView, image []byte) {
	var jsCode string
	if len(image) == 0 {
		jsCode = `document.getElementById("avatar").setAttribute("hidden");`
	} else {
		cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		avatarPath := filepath.Join(cwd, "avatar.jpg")
		err = ioutil.WriteFile(avatarPath, image, 0644)
		if err != nil {
			log.Fatal(err)
		}
		jsCode = `document.getElementById("avatar").removeAttribute("hidden");`
		jsCode += fmt.Sprintf(`document.getElementById("avatar").src = "%s";`, avatarPath)
	} 
	w.Eval(jsCode)

}

func updateListBox(w webview.WebView, listName string, items []string) {
	jsCode := fmt.Sprintf(`document.getElementById("%s").innerHTML="";`, listName)
	if len(items) != 0 {
		for _, item := range items {
			jsCode += fmt.Sprintf(`var option = document.createElement("option");
			                       option.innerHTML = "%s";
					               document.getElementById("%s").appendChild(option);`, item, listName)
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
		var avatar []byte
		var followers, topRepos []string
		if err == nil {
			userStr = user.Username
			reposCountStr = strconv.Itoa(user.ReposCount)
			followersCountStr = strconv.Itoa(len(user.Followers))
			followers = user.Followers
			topRepos = user.TopRepos
			avatar   = user.Avatar
		} else {
			//userStr = err.Error()
		}
		updateField(w, "Repos Count", reposCountStr)
		updateField(w, "User", userStr)
		updateField(w, "Followers", followersCountStr)
		updateListBox(w, "followers-list", followers)
		updateListBox(w, "top-repos", topRepos)
		updateAvatar(w, avatar)
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
		Debug: true,
	})
	w.SetColor(255, 255, 255, 255)
	defer w.Exit()
	w.Run()
	
}
