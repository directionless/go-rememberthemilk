package rememberthemilk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/clbanning/mxj"
	"github.com/pkg/errors"
)

type reqOptions struct {
	Params   *url.Values
	Unsigned bool
}

type reqOpt func(*reqOptions)

func Unsigned() reqOpt {
	return func(r *reqOptions) {
		r.Unsigned = true
	}
}

func Param(k, v string) reqOpt {
	return func(r *reqOptions) {
		r.Params.Add(k, v)
	}
}

// Response format
//
// While the docs [1] nominally claim to support json, there are
// issues. The json is generated by a naive xml2json routine, which
// can produce inconsist data structures. So, we'll need to stick with the xml.
// See [2] for more information
//
// It's additionally worth noting that the documented `format`
// variable, doesn't behave as documented. `format=json` is valid, but
// nothing else seems to be. So, omit it to get the xml
//
// [1] https://www.rememberthemilk.com/services/api/response.rtm
// [2] https://groups.google.com/forum/#!searchin/rememberthemilk-api/objects%7Csort:date/rememberthemilk-api/aNegBdRtw5E

func (rtm *RememberTheMilk) Req(method string, opts ...reqOpt) (mxj.Map, error) {
	reqOpts := &reqOptions{
		Params: &url.Values{},
	}

	for _, opt := range opts {
		opt(reqOpts)
	}

	urlPath, err := url.Parse("services/rest")
	if err != nil {
		return nil, errors.Wrap(err, "URL failure")
	}

	url := rtm.baseUrl.ResolveReference(urlPath)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http req")
	}

	reqOpts.Params.Add("method", method)

	// BUG: This logic sucks Need a way to override
	reqOpts.Params.Add("api_key", rtm.apiKey)
	if rtm.apiToken != "" {
		reqOpts.Params.Add("auth_token", rtm.apiToken)
	}

	//for k, v := range params {
	//	q.Add(k, v)
	//}

	// Add the params into the request
	req.URL.RawQuery = reqOpts.Params.Encode()

	// If this should be signed, then sign the request as-is, and add the sig param onto the end.
	if !reqOpts.Unsigned {
		reqOpts.Params.Add("api_sig", rtm.signAuthReq(req))
		req.URL.RawQuery = reqOpts.Params.Encode()
	}

	fmt.Println(req.URL.String())

	resp, err := rtm.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed fetch http req")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	// FIXME: Do we need something like io.CopyN here?
	// Probably depends on whether the xml decoder takes the stream or a string.

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("API request not successful. Got status %d", resp.StatusCode)
	}

	mv, err := mxj.NewMapXml(body) // unmarshal
	if err != nil {
		return nil, errors.Wrap(err, "failed xml unmarshal")
	}
	mv.Remove("rsp.api_key")

	status, _ := mv.ValueForPathString("rsp.-stat")
	if status != "ok" {
		return nil, errors.Errorf("API xml response not OK: %s\n", body)
	}

	return mv, nil
}