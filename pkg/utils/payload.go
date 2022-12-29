package utils

import (
	"fmt"
	"strings"
)

func GetRepoNameFromGithubURL(url string) string {
	return strings.TrimSuffix(strings.TrimPrefix(url, "https://github.com/"), ".git")
}

func GetUsernameFromGithubURL(url string) string {
	return strings.Split(strings.TrimPrefix(url, "https://github.com/"), "/")[0]
}

func MakeContainerTag(name, branch, head string) string {
	return fmt.Sprintf("%s:%s-%s", name, branch, head[:7])
}

func MakeRepoUrlWithAccessToken(username, token, url string) string {
	return fmt.Sprintf("https://%s:%s@%s", username, token, strings.TrimPrefix(url, "https://"))
}
