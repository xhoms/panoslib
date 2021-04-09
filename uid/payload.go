package uid

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	x "github.com/xhoms/panoslib/collection"
)

type payload struct {
	register       *IPTag
	unregister_ip  *string
	unregister_tag *string
	login          *UserMap
	logout_user    *string
	logout_ip      *string
	group          *UserGroup
	ungroup_user   *string
	ungroup_group  *string
}

// UIDBuilder provides a "functional programming"-like constructor to build a PAN-OS XML User-ID API Payload.
// Methods for UIDBuilder are not thread safe. All operations between NewUIDBuilder() and the final action
// (Payload(), UIDMessage() or Push()) must happen inside the same goroutine
type UIDBuilder struct {
	entries []payload
	err     error
}

// NewUIDBuilder returns an uninitialized UIDBuilder struct. Functional equivalent to UIDBuilder{}
func NewUIDBuilder() (mp UIDBuilder) {
	mp = UIDBuilder{}
	return
}

// NewBuilderFromPayload returns an initialized UIDBuilder struct with data contained in the provided message payload.
// Its common use case is to provide augmentation to an existing message of for "man-in-the-middle" applications.
// For the latter see additional details in the MemMonitor type
func NewBuilderFromPayload(p *x.UIDMsgPayload) (mp UIDBuilder) {
	mp = NewUIDBuilder()
	if p != nil {
		if p.Logout != nil {
			for _, e := range p.Logout.Entry {
				mp = mp.LogoutUser(e.Name, e.IP)
			}
		}
		if p.Login != nil {
			for _, e := range p.Login.Entry {
				if tout, err := ptrstr2uint(e.Timeout); err == nil {
					mp = mp.LoginUser(e.Name, e.IP, tout)
				}
			}
		}
		if p.UnregisterUser != nil {
			for _, e := range p.UnregisterUser.Entry {
				for _, t := range e.Tag.Member {
					mp = mp.UngroupUser(e.User, t.Member)
				}
			}
		}
		if p.RegisterUser != nil {
			for _, e := range p.RegisterUser.Entry {
				for _, t := range e.Tag.Member {
					if tout, err := ptrstr2uint(t.Timeout); err == nil {
						mp = mp.GroupUser(e.User, t.Member, tout)
					}
				}
			}
		}
		if p.Unregister != nil {
			for _, e := range p.Unregister.Entry {
				for _, t := range e.Tag.Member {
					mp = mp.UnregisterIP(e.IP, t.Member)
				}
			}
		}
		if p.Register != nil {
			for _, e := range p.Register.Entry {
				for _, t := range e.Tag.Member {
					if tout, err := ptrstr2uint(t.Timeout); err == nil {
						mp = mp.RegisterIP(e.IP, t.Member, tout)
					}
				}
			}
		}
	}
	return
}

// Register is used to add a list of ip-to-tag entries into the User-ID payload
func (mp UIDBuilder) Register(dag []IPTag) (mpB UIDBuilder) {
	if mp.err != nil {
		mpB = UIDBuilder{err: mp.err}
		return
	}
	mpC := UIDBuilder{entries: make([]payload, len(dag))}

	for idx := range dag {
		mpC.entries[idx] = payload{
			register: &dag[idx],
		}
	}
	mpB = mp.Add(mpC)
	return
}

// RegisterIP is used to add as single ip-to-tag entry in the User-ID payload
func (mp UIDBuilder) RegisterIP(ip, tag string, tout *uint) (mpB UIDBuilder) {
	if mp.err != nil {
		mpB = UIDBuilder{err: mp.err}
		return
	}
	mpB = mp.Register([]IPTag{{IP: ip, Tag: tag, Tout: tout}})
	return
}

// Unregister is used to add a list of ip-to-tag entries in the "unregister" section into the User-ID payload
func (mp UIDBuilder) Unregister(dag []IPTag) (mpB UIDBuilder) {
	if mp.err != nil {
		mpB = UIDBuilder{err: mp.err}
		return
	}
	mpC := UIDBuilder{entries: make([]payload, len(dag))}

	for idx := range dag {
		mpC.entries[idx] = payload{
			unregister_ip:  &dag[idx].IP,
			unregister_tag: &dag[idx].Tag,
		}
	}
	mpB = mp.Add(mpC)
	return
}

// UnregisterIP is used to add as single ip-to-tag entry in the "unregister" section into the User-ID payload
func (mp UIDBuilder) UnregisterIP(ip, tag string) (mpB UIDBuilder) {
	if mp.err != nil {
		mpB = UIDBuilder{err: mp.err}
		return
	}
	mpB = mp.Unregister([]IPTag{{IP: ip, Tag: tag}})
	return
}

// Login is used to add a list of user-to-ip entries into the User-ID payload
func (mp UIDBuilder) Login(uid []UserMap) (mpB UIDBuilder) {
	if mp.err != nil {
		mpB = UIDBuilder{err: mp.err}
		return
	}
	mpC := UIDBuilder{entries: make([]payload, len(uid))}

	for idx := range uid {
		mpC.entries[idx] = payload{
			login: &uid[idx],
		}
	}
	mpB = mp.Add(mpC)
	return
}

// LoginUser is used to add as single user-to-ip entry in the User-ID payload
func (mp UIDBuilder) LoginUser(user, ip string, tout *uint) (mpB UIDBuilder) {
	if mp.err != nil {
		mpB = UIDBuilder{err: mp.err}
		return
	}
	mpB = mp.Login([]UserMap{{IP: ip, User: user, Tout: tout}})
	return
}

// Logout is used to add a list of user-to-ip entries in the "logout" section into the User-ID payload
func (mp UIDBuilder) Logout(uid []UserMap) (mpB UIDBuilder) {
	if mp.err != nil {
		mpB = UIDBuilder{err: mp.err}
		return
	}
	mpC := UIDBuilder{entries: make([]payload, len(uid))}

	for idx := range uid {
		mpC.entries[idx] = payload{
			logout_user: &uid[idx].User,
			logout_ip:   &uid[idx].IP,
		}
	}
	mpB = mp.Add(mpC)
	return
}

// LogoutUser is used to add a single user-to-ip entry in the "logout" section into the User-ID payload
func (mp UIDBuilder) LogoutUser(user, ip string) (mpB UIDBuilder) {
	if mp.err != nil {
		mpB = UIDBuilder{err: mp.err}
		return
	}
	mpB = mp.Logout([]UserMap{{IP: ip, User: user}})
	return
}

// Group is used to add a list of user-to-group (DUG) entries into the User-ID payload
func (mp UIDBuilder) Group(dug []UserGroup) (mpB UIDBuilder) {
	if mp.err != nil {
		mpB = UIDBuilder{err: mp.err}
		return
	}
	mpC := UIDBuilder{entries: make([]payload, len(dug))}

	for idx := range dug {
		mpC.entries[idx] = payload{
			group: &dug[idx],
		}
	}
	mpB = mp.Add(mpC)
	return
}

// GroupUser is used to add a single user-to-group (DUG) entry into the User-ID payload
func (mp UIDBuilder) GroupUser(user, group string, tout *uint) (mpB UIDBuilder) {
	if mp.err != nil {
		mpB = UIDBuilder{err: mp.err}
		return
	}
	mpB = mp.Group([]UserGroup{{User: user, Group: group, Tout: tout}})
	return
}

// Ungroup is used to add a list of user-to-group entries in the "unregister-user" section into the User-ID payload
func (mp UIDBuilder) Ungroup(dug []UserGroup) (mpB UIDBuilder) {
	if mp.err != nil {
		mpB = UIDBuilder{err: mp.err}
		return
	}
	mpC := UIDBuilder{entries: make([]payload, len(dug))}

	for idx := range dug {
		mpC.entries[idx] = payload{
			ungroup_user:  &dug[idx].User,
			ungroup_group: &dug[idx].Group,
		}
	}
	mpB = mp.Add(mpC)
	return
}

// UngroupUser is used to add a single of user-to-group entry in the "unregister-user" section into the User-ID payload
func (mp UIDBuilder) UngroupUser(user, group string) (mpB UIDBuilder) {
	if mp.err != nil {
		mpB = UIDBuilder{err: mp.err}
		return
	}
	mpB = mp.Ungroup([]UserGroup{{User: user, Group: group}})
	return
}

// Add merges data from mpB builder into this builder
func (mp UIDBuilder) Add(mpB UIDBuilder) (mpC UIDBuilder) {
	if mp.err != nil {
		mpC = UIDBuilder{
			err: mp.err,
		}
		return
	}
	mpC = UIDBuilder{
		entries: append(mp.entries, mpB.entries...),
	}
	return
}

/*
Payload is a final action. It merges all accumulated data into a PAN-OS XML
User-ID API payload

If a variable implementing the Monitor interface is provided then a log entry
will be issued to it for every entry in the payload. Order of log entries will
be unregister > unregister-user > logout > login > register-user > register
*/
func (mp UIDBuilder) Payload(m Monitor) (p *x.UIDMsgPayload, err error) {
	if mp.err != nil {
		return nil, mp.err
	}
	size := len(mp.entries)
	reg := make(map[string]map[string]*uint, size)     // IP.Tag.Tout
	unreg := make(map[string]map[string]interface{})   // IP.Tag
	login := make(map[string]map[string]*uint, size)   // User.IP.Tout
	logout := make(map[string]map[string]interface{})  // User.IP
	group := make(map[string]map[string]*uint, size)   // User.Group.Tout
	ungroup := make(map[string]map[string]interface{}) // User.Group
	for _, e := range mp.entries {
		if e.unregister_ip != nil && e.unregister_tag != nil {
			if unrege, exists := unreg[*e.unregister_ip]; exists {
				unrege[*e.unregister_tag] = nil
			} else {
				unreg[*e.unregister_ip] = map[string]interface{}{*e.unregister_tag: nil}
			}
		}
		if e.ungroup_user != nil && e.ungroup_group != nil {
			if ungrpe, exists := ungroup[*e.ungroup_user]; exists {
				ungrpe[*e.ungroup_group] = nil
			} else {
				ungroup[*e.ungroup_user] = map[string]interface{}{*e.ungroup_group: nil}
			}
		}
		if e.logout_ip != nil && e.logout_user != nil {
			if logoute, exists := logout[*e.logout_user]; exists {
				logoute[*e.logout_ip] = nil
			} else {
				logout[*e.logout_user] = map[string]interface{}{*e.logout_ip: nil}
			}
		}
		if e.login != nil {
			if loge, exists := login[e.login.User]; exists {
				loge[e.login.IP] = e.login.Tout
			} else {
				login[e.login.User] = map[string]*uint{e.login.IP: e.login.Tout}
			}
		}
		if e.group != nil {
			if grpe, exists := group[e.group.User]; exists {
				grpe[e.group.Group] = e.group.Tout
			} else {
				group[e.group.User] = map[string]*uint{e.group.Group: e.group.Tout}
			}
		}
		if e.register != nil {
			if rege, exists := reg[e.register.IP]; exists {
				rege[e.register.Tag] = e.register.Tout
			} else {
				reg[e.register.IP] = map[string]*uint{e.register.Tag: e.register.Tout}
			}
		}
	}
	p = &x.UIDMsgPayload{}
	if len(unreg) > 0 {
		p.Unregister = &x.UIDMsgPldDAGRegUnreg{
			Entry: make([]x.UIDMsgPldDAGEntry, len(unreg)),
		}
		entryidx := 0
		for ip, tagmap := range unreg {
			dagentry := x.UIDMsgPldDAGEntry{
				IP: ip,
				Tag: x.UIDMsgPldDxGEntryTag{
					Member: make([]x.UIDMsgPldDxGEnrtryTagMember, len(tagmap)),
				},
			}
			tagidx := 0
			for tag := range tagmap {
				member := x.UIDMsgPldDxGEnrtryTagMember{
					Member: tag,
				}
				dagentry.Tag.Member[tagidx] = member
				tagidx++
				mp.log(m, Unregister, ip, tag, nil)
			}
			p.Unregister.Entry[entryidx] = dagentry
			entryidx++
		}
	}
	if len(ungroup) > 0 {
		p.UnregisterUser = &x.UIDMsgPldDUGRegUnreg{
			Entry: make([]x.UIDMsgPldDUGEntry, len(ungroup)),
		}
		entryidx := 0
		for user, groupmap := range ungroup {
			dugentry := x.UIDMsgPldDUGEntry{
				User: user,
				Tag: x.UIDMsgPldDxGEntryTag{
					Member: make([]x.UIDMsgPldDxGEnrtryTagMember, len(groupmap)),
				},
			}
			tagidx := 0
			for tag := range groupmap {
				member := x.UIDMsgPldDxGEnrtryTagMember{
					Member: tag,
				}
				dugentry.Tag.Member[tagidx] = member
				tagidx++
				mp.log(m, Ungroup, user, tag, nil)
			}
			p.UnregisterUser.Entry[entryidx] = dugentry
			entryidx++
		}
	}
	if len(logout) > 0 {
		p.Logout = &x.UIDMsgPldLogInOut{
			Entry: make([]x.UIDMsgPldLogEntry, 0, len(logout)*2),
		}
		for user, ipmap := range logout {
			for ip := range ipmap {
				logoutentry := x.UIDMsgPldLogEntry{
					Name: user,
					IP:   ip,
				}
				p.Logout.Entry = append(p.Logout.Entry, logoutentry)
				mp.log(m, Logout, user, ip, nil)
			}
		}
	}
	if len(login) > 0 {
		p.Login = &x.UIDMsgPldLogInOut{
			Entry: make([]x.UIDMsgPldLogEntry, 0, len(login)*2),
		}
		for user, ipmap := range login {
			for ip, tout := range ipmap {
				loginentry := x.UIDMsgPldLogEntry{
					Name: user,
					IP:   ip,
				}
				if tout != nil {
					toutstr := fmt.Sprint(*tout)
					loginentry.Timeout = &toutstr
				}
				p.Login.Entry = append(p.Login.Entry, loginentry)
				mp.log(m, Login, user, ip, tout)
			}
		}
	}
	if len(group) > 0 {
		p.RegisterUser = &x.UIDMsgPldDUGRegUnreg{
			Entry: make([]x.UIDMsgPldDUGEntry, len(group)),
		}
		entryidx := 0
		for user, grp := range group {
			dugentry := x.UIDMsgPldDUGEntry{
				User: user,
				Tag: x.UIDMsgPldDxGEntryTag{
					Member: make([]x.UIDMsgPldDxGEnrtryTagMember, len(grp)),
				},
			}
			tagidx := 0
			for tag, tout := range grp {
				member := x.UIDMsgPldDxGEnrtryTagMember{
					Member: tag,
				}
				if tout != nil {
					toutstr := fmt.Sprint(*tout)
					member.Timeout = &toutstr
				}
				dugentry.Tag.Member[tagidx] = member
				tagidx++
				mp.log(m, Group, user, tag, tout)
			}
			p.RegisterUser.Entry[entryidx] = dugentry
			entryidx++
		}
	}
	if len(reg) > 0 {
		p.Register = &x.UIDMsgPldDAGRegUnreg{
			Entry: make([]x.UIDMsgPldDAGEntry, len(reg)),
		}
		entryidx := 0
		for ip, tag := range reg {
			dagentry := x.UIDMsgPldDAGEntry{
				IP: ip,
				Tag: x.UIDMsgPldDxGEntryTag{
					Member: make([]x.UIDMsgPldDxGEnrtryTagMember, len(tag)),
				},
			}
			tagidx := 0
			for tag, tout := range tag {
				member := x.UIDMsgPldDxGEnrtryTagMember{
					Member: tag,
				}
				if tout != nil {
					toutstr := fmt.Sprint(*tout)
					member.Timeout = &toutstr
				}
				dagentry.Tag.Member[tagidx] = member
				tagidx++
				mp.log(m, Register, ip, tag, tout)
			}
			p.Register.Entry[entryidx] = dagentry
			entryidx++
		}
	}
	return
}

/*
UIDMessage is a final action. It merges all accumulated data into a ready-to
use PAN-OS XML User-ID API message

If a variable implementing the Monitor interface is provided then a log entry
will be issued to it for every entry in the payload. Order of log entries will
be unregister > unregister-user > logout > login > register-user > register
*/
func (mp UIDBuilder) UIDMessage(m Monitor) (u *x.UIDMessage, err error) {
	var p *x.UIDMsgPayload
	if p, err = mp.Payload(m); err == nil {
		u = &x.UIDMessage{
			Type:    "update",
			Version: "2.0",
			Payload: p,
		}
	}
	return
}

/*
Push is a final action. It merges all accumulated data into a ready-to use
PAN-OS XML User-ID API message and sends it to the device leveraging a
provided http.Client.

If a variable implementing the Monitor interface is provided then a log entry
will be issued to it for every entry in the payload. Order of log entries will
be unregister > unregister-user > logout > login > register-user > register
*/
func (mp UIDBuilder) Push(
	hostport, apikey string,
	c Client,
	m Monitor) (resp *http.Response, err error) {
	var u *x.UIDMessage
	if u, err = mp.UIDMessage(m); err == nil {
		target := "https://" + hostport + "/api/?"
		values := url.Values{
			"key":  []string{apikey},
			"type": []string{"user-id"},
		}
		var cmd []byte
		if cmd, err = xml.Marshal(u); err == nil {
			values["cmd"] = []string{string(cmd)}
			var req *http.Request
			if req, err = http.NewRequest(http.MethodPost, target, strings.NewReader(values.Encode())); err == nil {
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				resp, err = c.Do(req)
			}
		}
	}
	return
}

func (mp UIDBuilder) log(m Monitor, op Operation, subject string, value string, tout *uint) {
	if m != nil {
		m.Log(op, subject, value, tout)
	}
}

func ptrstr2uint(in *string) (out *uint, err error) {
	if in != nil {
		var ui uint64
		if ui, err = strconv.ParseUint(*in, 10, 32); err == nil {
			ui32 := uint(ui)
			out = &ui32
		}
	}
	return
}
