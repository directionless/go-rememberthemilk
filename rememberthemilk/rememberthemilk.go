package rememberthemilk

import (
	//"encoding/xml"
	//"io/ioutil"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	//"github.com/clbanning/mxj"

	"github.com/pkg/errors"
)

type RememberTheMilk struct {
	baseUrl   *url.URL
	userAgent string
	authFile  string
	auth      *authData

	httpClient *http.Client
}

type authData struct {
	Frob   string
	Key    string
	Secret string
	Token  string
}

func New() (*RememberTheMilk, error) {
	baseUrl, err := url.Parse("https://api.rememberthemilk.com")
	if err != nil {
		return nil, errors.Wrap(err, "parsing api endpoint url")
	}

	httpClient := &http.Client{
		// CheckRedirect: redirectPolicyFunc,
	}

	rtm := &RememberTheMilk{
		authFile:   "/tmp/rtm.creds",
		baseUrl:    baseUrl,
		auth:       &authData{Frob: ""},
		httpClient: httpClient,
	}

	return rtm, nil
}

// SetAuth accepts API credentials and saves them into the object
func (rtm *RememberTheMilk) SetAuth(apiKey string, apiSecret string) error {
	rtm.auth.Secret = apiSecret
	rtm.auth.Key = apiKey

	// Check connectivity
	alive, err := rtm.IsAlive()
	if err != nil {
		return err
	}

	if !alive {
		return errors.New("Failed initial connectivity check")
	}

	return rtm.SaveAuth()
}

// LoadAuth loads credentials from disk
//
// TODO look at keychain based things
func (rtm *RememberTheMilk) LoadAuth() error {
	savedJson, err := ioutil.ReadFile(rtm.authFile)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return nil
		default:
			return errors.Wrapf(err, "error loading file %s", rtm.authFile)
		}
	}

	return json.Unmarshal(savedJson, rtm.auth)
}

// SaveAuth saves credentials to disk
func (rtm *RememberTheMilk) SaveAuth() error {
	authJson, err := json.MarshalIndent(rtm.auth, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshalling auth")
	}

	err = ioutil.WriteFile(rtm.authFile, authJson, 0600)
	return nil
}

func (rtm *RememberTheMilk) IsAlive() (bool, error) {
	nonce := nonce(16)

	var nonceRespose struct {
		*ResponseObject
		Nonce string `xml:"nonce"`
	}

	if err := rtm.Req("rtm.test.echo", &nonceRespose, Param("nonce", nonce), Unsigned()); err != nil {
		return false, errors.Wrap(err, "failed req")
	}

	return nonceRespose.Nonce == nonce, nil
}

func (rtm *RememberTheMilk) EnsureAuth() error {
	loggedIn, err := rtm.CheckToken()
	if err != nil {
		return errors.Wrap(err, "checktoken failed")
	}

	if loggedIn {
		return rtm.SaveAuth()
	}

	if err := rtm.Authenticate(); err != nil {
		return errors.Wrap(err, "authenticating")
	}

	return rtm.SaveAuth()
}

func (rtm *RememberTheMilk) CheckToken() (bool, error) {
	r := &checkTokenResponse{}
	err := rtm.Req("rtm.auth.checkToken", r)

	switch {
	case r.Auth.Token != "":
		rtm.auth.Token = r.Auth.Token
		return true, nil
	case r.Error.Code == "98":
		return false, nil
	case err != nil:
		return false, errors.Wrap(err, "unknown error checking token")
	}

	return false, errors.New("what")
}
