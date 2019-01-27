package rememberthemilk

import (
	//"encoding/xml"
	//"io/ioutil"

	"net/http"
	"net/url"

	//"github.com/clbanning/mxj"
	"github.com/pkg/errors"
)

type RememberTheMilk struct {
	apiFrob   string
	apiKey    string
	apiSecret string
	apiToken  string
	baseUrl   *url.URL
	userAgent string

	httpClient *http.Client
}

func New(apiKey string, apiSecret string) (*RememberTheMilk, error) {
	baseUrl, err := url.Parse("https://api.rememberthemilk.com")
	if err != nil {
		return nil, errors.Wrap(err, "parsing api endpoint url")
	}

	httpClient := &http.Client{
		// CheckRedirect: redirectPolicyFunc,
	}

	rtm := &RememberTheMilk{
		baseUrl:    baseUrl,
		apiSecret:  apiSecret,
		apiKey:     apiKey,
		httpClient: httpClient,
		apiFrob:    "",
	}

	// Check connectivity
	alive, err := rtm.IsAlive()
	if err != nil {
		return nil, err
	}

	if !alive {
		return nil, errors.New("Failed initial connectivity check")
	}

	return rtm, nil
}

func (rtm *RememberTheMilk) IsAlive() (bool, error) {
	nonce := nonce(16)

	mv, err := rtm.Req("rtm.test.echo", Param("nonce", nonce), Unsigned())
	if err != nil {
		return false, errors.Wrap(err, "failed req")
	}

	v, err := mv.ValueForPathString("rsp.nonce")
	if err != nil {
		return false, errors.Wrap(err, "Failed to get value")
	}

	return v == nonce, nil
}

func (rtm *RememberTheMilk) EnsureAuth() error {
	// BUG needs a checktoken
	return rtm.Authenticate()

}
