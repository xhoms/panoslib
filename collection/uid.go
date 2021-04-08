/*
package collection provides GO structs meant to be used to build PAN-OS XML API payloads
with xml.Marshal() as well as structs meant to be used to parse PAN-OS XML API responses
*/
package collection

import (
	"encoding/xml"
)

// UIDMessage is the main placeholder for a PAN-OS XML API UserID message
type UIDMessage struct {
	XMLName xml.Name       `xml:"uid-message"`
	Type    string         `xml:"type"`
	Payload *UIDMsgPayload `xml:"payload"`
	Version string         `xml:"version"`
}

// UIDMsgPayload contains data for different PAN-OS XML API UserID operations
type UIDMsgPayload struct {
	XMLName        xml.Name              `xml:"payload"`
	Register       *UIDMsgPldDAGRegUnreg `xml:"register,omitempty"`
	Unregister     *UIDMsgPldDAGRegUnreg `xml:"unregister,omitempty"`
	RegisterUser   *UIDMsgPldDUGRegUnreg `xml:"register-user,omitempty"`
	UnregisterUser *UIDMsgPldDUGRegUnreg `xml:"unregister-user,omitempty"`
	Login          *UIDMsgPldLogInOut    `xml:"login,omitempty"`
	Logout         *UIDMsgPldLogInOut    `xml:"logout,omitempty"`
}

// UIDMsgPldLogInOut is the list of entries for a UserID login or logout operation
type UIDMsgPldLogInOut struct {
	Entry []UIDMsgPldLogEntry `xml:"entry"`
}

// UIDMsgPldLogEntry contains data for a single login / logout operation
type UIDMsgPldLogEntry struct {
	Name    string  `xml:"name,attr"`
	IP      string  `xml:"ip,attr"`
	Timeout *string `xml:"timeout,attr,omitempty"`
}

/* UIDMsgPldDUGRegUnreg is the list of entries for a UserID Dynamic User Group (register / unregister) operation

Example User-ID DUG payload

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
type UIDMsgPldDUGRegUnreg struct {
	Entry []UIDMsgPldDUGEntry `xml:"entry"`
}

// UIDMsgPldDUGEntry contains data for a single Dynamic User Group operation
type UIDMsgPldDUGEntry struct {
	Tag  UIDMsgPldDxGEntryTag `xml:"tag"`
	User string               `xml:"user,attr"`
}

/* UIDMsgPldDAGRegUnreg is the list of entries for a UserID Dynamic Address Group (register / unregister) operation

Example User-ID DAG payload

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
type UIDMsgPldDAGRegUnreg struct {
	Entry []UIDMsgPldDAGEntry `xml:"entry"`
}

// UIDMsgPldDAGEntry contains data for a single Dynamic Address Group operation
type UIDMsgPldDAGEntry struct {
	Tag        UIDMsgPldDxGEntryTag `xml:"tag"`
	IP         string               `xml:"ip,attr"`
	Persistent string               `xml:"persistent,attr,omitempty"`
}

// UIDMsgPldDxGEntryTag contains the list of tags for a DAG or DUG operation
type UIDMsgPldDxGEntryTag struct {
	Member []UIDMsgPldDxGEnrtryTagMember `xml:"member"`
}

// UIDMsgPldDxGEnrtryTagMember data for a single tag operation
type UIDMsgPldDxGEnrtryTagMember struct {
	Member  string  `xml:",chardata"`
	Timeout *string `xml:"timeout,attr,omitempty"`
}
