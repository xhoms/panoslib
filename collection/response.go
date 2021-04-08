package collection

import "encoding/xml"

/*
<response status="success">
  <result>
    <uid-response>
      <version>2.0</version>
      <payload>
        <unregister></unregister>
        <register></register>
      </payload>
    </uid-response>
  </result>
</response>
*/

// APIResponse is the placeholder for a PAN-OS XML API Response
type APIResponse struct {
	XMLName xml.Name `xml:"response"`
	Status  string   `xml:"status,attr"`
	Result  *struct {
		UidResponse *UIDResponse `xml:"uid-response"`
	} `xml:"result,omitempty"`
	Msg *struct {
		Line []struct {
			UidResponse *UIDResponse `xml:"uid-response"`
		} `xml:"line"`
	} `xml:"msg,omitempty"`
}

// UIDResponse is the placeholder for a PAN-OS XML UserID API Response
type UIDResponse struct {
	Version string `xml:"version"`
	Payload struct {
		Unregister *struct {
			Entry []UIDResponseEntry `xml:"entry"`
		} `xml:"unregister"`
		Register *struct {
			Entry []UIDResponseEntry `xml:"entry"`
		} `xml:"register"`
	} `xml:"payload"`
}

// UIDResponseEntry is a PAN-OS XML UserID API Response item
type UIDResponseEntry struct {
	IP  string `xml:"ip,attr"`
	Msg string `xml:"message,attr"`
}
