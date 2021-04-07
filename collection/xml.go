package collection

import (
	"encoding/xml"
)

/* Example User-ID DAG payload
<uid-message>
  <type>update</type>
  <payload>
    <register>
      <entry ip="10.10.10.10">
        <tag>
          <member timeout="10">tag10</member>
        </tag>
      </entry>
    </register>
  </payload>
</uid-message>
*/

type Member struct {
	Member  string  `xml:",chardata"`
	Timeout *string `xml:"timeout,attr,omitempty"`
}

type Tag struct {
	Member []Member `xml:"member"`
}

type DAGEntry struct {
	Tag        Tag    `xml:"tag"`
	IP         string `xml:"ip,attr"`
	Persistent string `xml:"persistent,attr,omitempty"`
}

type DAGRegUnreg struct {
	Entry []DAGEntry `xml:"entry"`
}

/* Example User-ID DUG payload
<uid-message>
     <type>update</type>
     <payload>
          <login>
               <entry name="armis-1.1.1.1" ip="1.1.1.1" timeout="60">
               </entry>
          </login>
          <register-user>
            <entry user="armis-1.1.1.1">
              <tag>
                <member>tag30</member>
              </tag>
            </entry>
          </register-user>
     </payload>
</uid-message>
*/

type LogEntry struct {
	Name    string  `xml:"name,attr"`
	IP      string  `xml:"ip,attr"`
	Timeout *string `xml:"timeout,attr,omitempty"`
}

type LogInOut struct {
	Entry []LogEntry `xml:"entry"`
}

type DUGEntry struct {
	Tag  Tag    `xml:"tag"`
	User string `xml:"user,attr"`
}

type DUGRegUnreg struct {
	Entry []DUGEntry `xml:"entry"`
}

type Payload struct {
	XMLName        xml.Name     `xml:"payload"`
	Register       *DAGRegUnreg `xml:"register,omitempty"`
	Unregister     *DAGRegUnreg `xml:"unregister,omitempty"`
	RegisterUser   *DUGRegUnreg `xml:"register-user,omitempty"`
	UnregisterUser *DUGRegUnreg `xml:"unregister-user,omitempty"`
	Login          *LogInOut    `xml:"login,omitempty"`
	Logout         *LogInOut    `xml:"logout,omitempty"`
}

type UIDMessage struct {
	XMLName xml.Name `xml:"uid-message"`
	Type    string   `xml:"type"`
	Payload *Payload `xml:"payload"`
	Version string   `xml:"version"`
}

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

type ResponseEntry struct {
	IP  string `xml:"ip,attr"`
	Msg string `xml:"message,attr"`
}

type UIDResponse struct {
	Version string `xml:"version"`
	Payload struct {
		Unregister *struct {
			Entry []ResponseEntry `xml:"entry"`
		} `xml:"unregister"`
		Register *struct {
			Entry []ResponseEntry `xml:"entry"`
		} `xml:"register"`
	} `xml:"payload"`
}

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
