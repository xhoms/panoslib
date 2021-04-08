package uid

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	x "github.com/xhoms/panoslib/collection"
)

// Validate provides PAN-OS XML User-ID API response validation. Error will be raised either
// by underlying http/net errors of because the PAN-OS User-ID response contains a non "success"
// status code
func Validate(resp *http.Response, resperr error) (apiResp *x.APIResponse, err error) {
	if resperr != nil {
		err = resperr
		return
	}
	if scode := resp.StatusCode; scode == 200 {
		if body, readErr := ioutil.ReadAll(resp.Body); readErr == nil {
			apiResp = &x.APIResponse{}
			if xmlerr := xml.Unmarshal(body, apiResp); xmlerr == nil {
				if apiResp.Status != "success" {
					err = fmt.Errorf("returned a non-sucess response (status: '%v')", apiResp.Status)
				}
			} else {
				err = fmt.Errorf("error unmarshaling xml body")
			}
		} else {
			err = fmt.Errorf("error reading body (%v)", readErr.Error())
		}
	} else {
		err = fmt.Errorf("replied with status code '%v'", scode)
	}
	return
}
