package twitter

import (
	"net/url"
	"regexp"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var statusRegex = regexp.MustCompile(`status\/([0-9]+)`)

type Twitter struct {
	client *twitter.Client
}

func New(accessToken string) *Twitter {
	config := &oauth2.Config{}
	token := &oauth2.Token{AccessToken: accessToken}
	// OAuth2 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	return &Twitter{client}
}

func (t *Twitter) Get(tweetURL string) (string, error) {
	// status show
	u, err := url.Parse(tweetURL)
	if err != nil {
		return "", errors.Wrap(err, "url.Parse")
	}
	parts := statusRegex.FindStringSubmatch(u.Path)
	if len(parts) == 0 {
		return "", errors.Errorf("invalid url: %q", tweetURL)
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", errors.Wrap(err, "strconv.ParseInt")
	}
	tweet, _, err := t.client.Statuses.Show(id, nil)
	if err != nil {
		return "", errors.Wrap(err, "t.client.Statuses.Show")
	}

	return tweet.Text, nil
}
