package rememberthemilk

import (
	//"io/ioutil"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func (rtm *RememberTheMilk) Authenticate() error {
	if err := rtm.getFrob(); err != nil {
		return err
	}

	humanURL, err := rtm.humanURL()
	if err != nil {
		return errors.Wrap(err, "Failed to make human URL")
	}

	fmt.Printf("\n\n\nGotta human auth to:\n%s\nreturn when done... ", humanURL)
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	if err := rtm.getToken(); err != nil {
		return err
	}

	return nil
}

func (rtm *RememberTheMilk) getToken() error {
	mv, err := rtm.Req("rtm.auth.getToken", Param("frob", rtm.apiFrob))
	if err != nil {
		return errors.Wrap(err, "failed to get token")
	}

	//fmt.Println(mv)
	token, err := mv.ValueForPathString("rsp.auth.token")
	if err != nil {
		spew.Dump(mv)
		return errors.Wrap(err, "Failed to get value")
	}
	rtm.apiToken = token

	return nil

}
func (rtm *RememberTheMilk) humanURL() (string, error) {
	urlPath, err := url.Parse("services/auth")

	if err != nil {
		return "", errors.Wrap(err, "URL failure")
	}

	url := rtm.baseUrl.ResolveReference(urlPath)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create http req")
	}

	q := req.URL.Query()
	q.Add("perms", "delete")
	q.Add("api_key", rtm.apiKey)
	q.Add("frob", rtm.apiFrob)

	//Sign it. Signing is based on the URL param, so we add more, and re-encode
	req.URL.RawQuery = q.Encode()
	q.Add("api_sig", rtm.signAuthReq(req))
	req.URL.RawQuery = q.Encode()

	return req.URL.String(), nil
}

// signAuthReq generates a signature via the algorythem described in
// https://www.rememberthemilk.com/services/api/authentication.rtm
//
// This code has all kinds of assumptiopns baked in. But, they derive
// from the documented use case. So it's what it is.
func (rtm *RememberTheMilk) signAuthReq(req *http.Request) string {
	// We need to sort a bunch of stuff. So let's make an array. Leave
	// one empty spot for the secret to land in
	blobs := make([]string, len(req.URL.Query())+1)

	i := 0
	for k, v := range req.URL.Query() {
		blobs[i] = strings.Join([]string{k, v[0]}, "") //fmt.Sprintf("%s%s", k, v[0])
		i += 1
	}

	sort.Strings(blobs)

	// Prepend our secret
	//
	// This works, because we  left a "" in the array, which has been sorted to the beginning.
	blobs[0] = rtm.apiSecret

	concatSecret := strings.Join(blobs, "")

	hexSig := md5.Sum([]byte(concatSecret))

	return hex.EncodeToString(hexSig[:])
}

// getFrob calls the RTM api to get a frob.
// (eg: a token)
func (rtm *RememberTheMilk) getFrob() error {
	// Get the frob (we mean token)
	mv, err := rtm.Req("rtm.auth.getFrob")
	if err != nil {
		return errors.Wrap(err, "failed to get frob")
	}
	frob, err := mv.ValueForPathString("rsp.frob")
	if err != nil {
		return errors.Wrap(err, "Failed to get frob")
	}

	rtm.apiFrob = frob

	return nil
}
