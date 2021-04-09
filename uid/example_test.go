package uid_test

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/xhoms/panoslib/uid"
)

/*
Create a PAN-OS XML User-ID compatible message from a list of user-to-ip maps, user-to-group maps and a list
of ip-to-tag maps. Notice timeouts are passed as pointers to distinguish between zero and not-present
*/
func ExampleUIDBuilder() {
	var tout uint = 60
	login := []uid.UserMap{
		{User: "foo@test.local", IP: "1.1.1.1", Tout: &tout},
		{User: "bar@test.local", IP: "2.2.2.2"},
	}
	group := []uid.UserGroup{
		{Group: "admin", User: "foo@test.local"},
		{Group: "devops", User: "bar@test.local", Tout: &tout},
	}
	tag := []uid.IPTag{
		{Tag: "windows", IP: "1.1.1.1"},
		{Tag: "linux", IP: "2.2.2.2"},
		{Tag: "avscanned", IP: "2.2.22", Tout: &tout},
	}
	if uidmsg, err := uid.NewUIDBuilder().
		Login(login).
		Group(group).
		Register(tag).
		UIDMessage(nil); err == nil {
		if msg, err := xml.MarshalIndent(uidmsg, "", " "); err == nil {
			log.Println(string(msg))
		}
	}
}

/*
Create a User-ID payload with a single IP-to-tag entry
*/
func ExampleUIDBuilder_RegisterIP() {
	var tout uint = 60
	if p, err := uid.NewUIDBuilder().
		RegisterIP("1.1.1.1", "windows", &tout).
		Payload(nil); err == nil {
		if b, err := xml.Marshal(p); err == nil {
			fmt.Println(string(b))
		}
	}
	// Output: <payload><register><entry ip="1.1.1.1"><tag><member timeout="60">windows</member></tag></entry></register></payload>
}

/*
Create a User-ID payload with a single user-to-IP entry
*/
func ExampleUIDBuilder_LoginUser() {
	var tout uint = 60
	if p, err := uid.NewUIDBuilder().
		LoginUser("foo@test.local", "1.1.1.1", &tout).
		Payload(nil); err == nil {
		if b, err := xml.Marshal(p); err == nil {
			fmt.Println(string(b))
		}
	}
	// Output: <payload><login><entry name="foo@test.local" ip="1.1.1.1" timeout="60"></entry></login></payload>
}

/*
Create a User-ID payload with a single user-to-group (DUG) entry
*/
func ExampleUIDBuilder_GroupUser() {
	var tout uint = 60
	if p, err := uid.NewUIDBuilder().
		GroupUser("foo@test.local", "admin", &tout).
		Payload(nil); err == nil {
		if b, err := xml.Marshal(p); err == nil {
			fmt.Println(string(b))
		}
	}
	// Output: <payload><register-user><entry user="foo@test.local"><tag><member timeout="60">admin</member></tag></entry></register-user></payload>
}

/*
Create a User-ID payload with a single IP-to-tag entry, push the message to the PAN-OS device using the
default http client and parse the response.
*/
func ExampleUIDBuilder_Push() {
	var tout uint = 60
	resp, err := uid.NewUIDBuilder().
		RegisterIP("1.1.1.1", "windows", &tout).
		Push("10.1.1.1:443", "<my-api-key>", http.DefaultClient, nil)
	if apiResp, err := uid.Validate(resp, err); err == nil {
		fmt.Println(apiResp.Status)
	} else {
		fmt.Println(err)
	}
}

/*
Create a User-ID payload with a single IP-to-tag entry, push the message to the PAN-OS device using the
default http client and parse the response.
*/
func ExampleValidate() {
	var tout uint = 60
	resp, err := uid.NewUIDBuilder().
		RegisterIP("1.1.1.1", "windows", &tout).
		Push("10.1.1.1:443", "<my-api-key>", http.DefaultClient, nil)
	if apiResp, err := uid.Validate(resp, err); err == nil {
		fmt.Println(apiResp.Status)
	} else {
		fmt.Println(err)
	}
}
