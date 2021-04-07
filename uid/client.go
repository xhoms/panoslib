package uid

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	x "github.com/xhoms/panoslib/collection"
)

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
