package pkgparse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ReleaseInfo struct {
	Url     string
	Assets  []AssetInfo
	TagName string `json:"tag_name"`
	Date    string `json:"published_at"`
}

type AssetInfo struct {
	Name               string
	Size               uint32
	BrowserDownloadUrl string `json:"browser_download_url"`
}

type ReleaseTagInfo struct {
	TagName string `json:"tag_name"`
	Date    string `json:"published_at"`
}

func getLatestGithubReleaseTag(user string, repo string) (*ReleaseTagInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", user, repo)
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if !(r.StatusCode >= 200 && r.StatusCode < 300) {
		return nil, fmt.Errorf("bad HTTP Response: %s", r.Status)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var releases []ReleaseTagInfo
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return nil, fmt.Errorf("github releases JSON response not in expected format")
	}
	if len(releases) == 0 {
		return nil, fmt.Errorf("expected at least one release listed at %s, unable to resolve latest", url)
	}
	return &releases[0], nil
}
