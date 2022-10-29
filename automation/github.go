package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

func updateAutomationHealthComment(check_updates_failed, manifest_update_failed []string, client github.Client, func_ctx context.Context) error {
	issue_number := 900
	// get comments from the automation health check issue
	comments, _, err := client.Issues.ListComments(func_ctx, "vedantmgoyal2009", "winget-manifests-manager", issue_number, nil)
	if err != nil {
		return fmt.Errorf("error getting comments on issue %d: %v\n", issue_number, err)
	}

	// delete all comments from the bot (GitHub app)
	for _, comment := range comments {
		if comment.User.GetLogin() == "vedantmgoyal2009[bot]" {
			_, err := client.Issues.DeleteComment(func_ctx, "vedantmgoyal2009", "winget-manifests-manager", comment.GetID())
			if err != nil {
				return fmt.Errorf("error deleting comment: %v\n", err)
			}
		}
	}

	// create a comment body with the results of the automation run
	var commentBody string = "### Results of Automation run [" + os.Getenv("GITHUB_RUN_NUMBER") + "](https://github.com/vedantmgoyal2009/vedantmgoyal2009/actions/runs/" + os.Getenv("GITHUB_RUN_ID") + ")\r\n"
	commentBody += "**Error while checking for updates for packages:** "
	if len(check_updates_failed) != 0 {
		commentBody += strings.Join(check_updates_failed, "\r\n") + "\r\n"
	} else {
		commentBody += "No errors while checking for updates for packages :tada:\r\n"

	}
	commentBody += "**Error while upgrading packages:** "
	if len(manifest_update_failed) != 0 {
		commentBody += strings.Join(manifest_update_failed, "\r\n")
	} else {
		commentBody += "All packages were updated successfully :tada:"
	}

	// create a new comment with the comment body
	comment, res, err := client.Issues.CreateComment(func_ctx, "vedantmgoyal2009", "winget-manifests-manager", issue_number, &github.IssueComment{Body: &commentBody})
	log.Printf("comment: %v\n", comment)
	log.Printf("response: %v\n", res)
	if err != nil {
		return fmt.Errorf("error creating comment: %v\n", err)
	}
	return nil
}

func refreshGithubAppTokenIfExpired(token_var *string, wasRefreshed *bool) error {
	// get rate limit api response
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/rate_limit", nil)
	if err != nil {
		return fmt.Errorf("error creating request to get rate limit: %v\n", err)
	}
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + *token_var},
		"Accept":        []string{"application/vnd.github.v3+json"},
	}
	if _, err := httpClient.Do(req); err != nil {
		token, err := getGithubAppAuthToken()
		if err != nil {
			return fmt.Errorf("error refreshing github app auth token: %v\n", err)
		}
		*token_var = token
		*wasRefreshed = true
	} else {
		*wasRefreshed = false
	}
	return nil
}

func getGithubAppAuthToken() (string, error) {
	// get the private key from the environment variable
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("BOT_PRIVATE_KEY")))
	if err != nil {
		return "", fmt.Errorf("error parsing private key: %v\n", err)
	}

	// create a new jwt token
	iat := time.Now().Add(-30 * time.Second).Truncate(time.Second)
	exp := iat.Add(2 * time.Minute)
	jwt_token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": iat.Unix(),
		"exp": exp.Unix(),
		"iss": os.Getenv("BOT_APP_ID"),
	})

	// sign the jwt token with the private key
	signed_jwt_token, err := jwt_token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing jwt token with private key: %v\n", err)
	}

	// get github app installation token
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.github.com/app/installations/"+os.Getenv("BOT_INSTALLATION_ID")+"/access_tokens", nil)
	if err != nil {
		return "", fmt.Errorf("error creating request to get github app installation token: %v\n", err)
	}
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + signed_jwt_token},
		"Accept":        []string{"application/vnd.github.v3+json"},
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting github app installation token: %v\n", err)
	}

	// parse the response body
	var res_json map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&res_json)
	if err != nil {
		log.Fatalf("error parsing response body: %v\n", err)
	}

	return res_json["token"].(string), nil
}

func getGithubClient(token string, func_ctx context.Context) github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(func_ctx, ts)

	client := github.NewClient(tc)
	return *client
}
