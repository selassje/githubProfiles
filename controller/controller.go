package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
)

type User struct {
	Username   string
	Avatar     []byte
	Followers  []string
	ReposCount int
	TopRepos   []string
}

func getHttpResponseBody(url string) (body io.Reader, err error) {
	httpResp, err := http.Get(url)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	body = httpResp.Body
	return
}

func get(url string) (err error, image []byte) {
	httpResp, err := http.Get(url)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	image, err = ioutil.ReadAll(httpResp.Body)
	return
}

func performRESTJsonQuery(query string, queryResponse interface{}) (err error) {
	httpResp, err := http.Get(query)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	dec := json.NewDecoder(httpResp.Body)
	err = dec.Decode(queryResponse)
	if err == nil || err == io.EOF {
		err = nil
	}
	return
}

func GetUserInfo(username string) (user User, err error) {
	var userQueryResponse struct {
		Total_count int
		Items       []struct {
			Login         string
			Avatar_url    string
			Followers_url string
			Repos_url     string
		}
	}
	err = performRESTJsonQuery("https://api.github.com/search/users?q="+username, &userQueryResponse)
	if err != nil {
		return
	}

	if userQueryResponse.Total_count == 0 {
		err = fmt.Errorf("User not found")
		return
	}

	foundUser := userQueryResponse.Items[0]
	var reposQueryResponse []struct{ Name string }
	err = performRESTJsonQuery(foundUser.Repos_url, &reposQueryResponse)
	if err != nil {
		return
	}

	var followersQueryResponse []struct{ Login string }
	err = performRESTJsonQuery(foundUser.Followers_url, &followersQueryResponse)
	if err != nil {
		return
	}

	httpRespAvatar, err := http.Get(foundUser.Avatar_url)
	if err != nil {
		return
	}
	defer httpRespAvatar.Body.Close()
	avatar, err := ioutil.ReadAll(httpRespAvatar.Body)
	if err != nil {
		return
	}

	for _, follower := range followersQueryResponse {
		user.Followers = append(user.Followers, follower.Login)
	}
	topReposLimit := int(math.Min(4, float64(len(reposQueryResponse))))
	for _, repo := range reposQueryResponse[:topReposLimit] {
		user.TopRepos = append(user.TopRepos, repo.Name)
	}
	user.Username = foundUser.Login
	user.Avatar = avatar
	user.ReposCount = len(reposQueryResponse)
	return
}
