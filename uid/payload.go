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

type mpayload struct {
	entries []payload
	err     error
}

func NewPayload() (mp mpayload) {
	mp = mpayload{}
	return
}

func NewFromPayload(p *x.Payload) (mp mpayload) {
	mp = NewPayload()
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

func (mp mpayload) Register(dag []IPTag) (mpB mpayload) {
	if mp.err != nil {
		mpB = mpayload{err: mp.err}
		return
	}
	mpC := mpayload{entries: make([]payload, len(dag))}

	for idx := range dag {
		mpC.entries[idx] = payload{
			register: &dag[idx],
		}
	}
	mpB = mp.Add(mpC)
	return
}

func (mp mpayload) RegisterIP(ip, tag string, tout *uint) (mpB mpayload) {
	if mp.err != nil {
		mpB = mpayload{err: mp.err}
		return
	}
	mpB = mp.Register([]IPTag{{IP: ip, Tag: tag, Tout: tout}})
	return
}

func (mp mpayload) Unregister(dag []IPTag) (mpB mpayload) {
	if mp.err != nil {
		mpB = mpayload{err: mp.err}
		return
	}
	mpC := mpayload{entries: make([]payload, len(dag))}

	for idx := range dag {
		mpC.entries[idx] = payload{
			unregister_ip:  &dag[idx].IP,
			unregister_tag: &dag[idx].Tag,
		}
	}
	mpB = mp.Add(mpC)
	return
}

func (mp mpayload) UnregisterIP(ip, tag string) (mpB mpayload) {
	if mp.err != nil {
		mpB = mpayload{err: mp.err}
		return
	}
	mpB = mp.Unregister([]IPTag{{IP: ip, Tag: tag}})
	return
}

func (mp mpayload) Login(uid []UserMap) (mpB mpayload) {
	if mp.err != nil {
		mpB = mpayload{err: mp.err}
		return
	}
	mpC := mpayload{entries: make([]payload, len(uid))}

	for idx := range uid {
		mpC.entries[idx] = payload{
			login: &uid[idx],
		}
	}
	mpB = mp.Add(mpC)
	return
}

func (mp mpayload) LoginUser(user, ip string, tout *uint) (mpB mpayload) {
	if mp.err != nil {
		mpB = mpayload{err: mp.err}
		return
	}
	mpB = mp.Login([]UserMap{{IP: ip, User: user, Tout: tout}})
	return
}

func (mp mpayload) Logout(uid []UserMap) (mpB mpayload) {
	if mp.err != nil {
		mpB = mpayload{err: mp.err}
		return
	}
	mpC := mpayload{entries: make([]payload, len(uid))}

	for idx := range uid {
		mpC.entries[idx] = payload{
			logout_user: &uid[idx].User,
			logout_ip:   &uid[idx].IP,
		}
	}
	mpB = mp.Add(mpC)
	return
}

func (mp mpayload) LogoutUser(user, ip string) (mpB mpayload) {
	if mp.err != nil {
		mpB = mpayload{err: mp.err}
		return
	}
	mpB = mp.Logout([]UserMap{{IP: ip, User: user}})
	return
}

func (mp mpayload) Group(dug []UserGroup) (mpB mpayload) {
	if mp.err != nil {
		mpB = mpayload{err: mp.err}
		return
	}
	mpC := mpayload{entries: make([]payload, len(dug))}

	for idx := range dug {
		mpC.entries[idx] = payload{
			group: &dug[idx],
		}
	}
	mpB = mp.Add(mpC)
	return
}

func (mp mpayload) GroupUser(user, group string, tout *uint) (mpB mpayload) {
	if mp.err != nil {
		mpB = mpayload{err: mp.err}
		return
	}
	mpB = mp.Group([]UserGroup{{User: user, Group: group, Tout: tout}})
	return
}

func (mp mpayload) Ungroup(dug []UserGroup) (mpB mpayload) {
	if mp.err != nil {
		mpB = mpayload{err: mp.err}
		return
	}
	mpC := mpayload{entries: make([]payload, len(dug))}

	for idx := range dug {
		mpC.entries[idx] = payload{
			ungroup_user:  &dug[idx].User,
			ungroup_group: &dug[idx].Group,
		}
	}
	mpB = mp.Add(mpC)
	return
}

func (mp mpayload) UngroupUser(user, group string) (mpB mpayload) {
	if mp.err != nil {
		mpB = mpayload{err: mp.err}
		return
	}
	mpB = mp.Ungroup([]UserGroup{{User: user, Group: group}})
	return
}

func (mp mpayload) Add(mpB mpayload) (mpC mpayload) {
	if mp.err != nil {
		mpC = mpayload{
			err: mp.err,
		}
		return
	}
	mpC = mpayload{
		entries: append(mp.entries, mpB.entries...),
	}
	return
}

func (mp mpayload) Payload(m Monitor) (p *x.Payload, err error) {
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
	p = &x.Payload{}
	if len(unreg) > 0 {
		p.Unregister = &x.DAGRegUnreg{
			Entry: make([]x.DAGEntry, len(unreg)),
		}
		entryidx := 0
		for ip, tagmap := range unreg {
			dagentry := x.DAGEntry{
				IP: ip,
				Tag: x.Tag{
					Member: make([]x.Member, len(tagmap)),
				},
			}
			tagidx := 0
			for tag := range tagmap {
				member := x.Member{
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
		p.UnregisterUser = &x.DUGRegUnreg{
			Entry: make([]x.DUGEntry, len(ungroup)),
		}
		entryidx := 0
		for user, groupmap := range ungroup {
			dugentry := x.DUGEntry{
				User: user,
				Tag: x.Tag{
					Member: make([]x.Member, len(groupmap)),
				},
			}
			tagidx := 0
			for tag := range groupmap {
				member := x.Member{
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
		p.Logout = &x.LogInOut{
			Entry: make([]x.LogEntry, 0, len(logout)*2),
		}
		for user, ipmap := range logout {
			for ip := range ipmap {
				logoutentry := x.LogEntry{
					Name: user,
					IP:   ip,
				}
				p.Logout.Entry = append(p.Logout.Entry, logoutentry)
				mp.log(m, Logout, user, ip, nil)
			}
		}
	}
	if len(login) > 0 {
		p.Login = &x.LogInOut{
			Entry: make([]x.LogEntry, 0, len(login)*2),
		}
		for user, ipmap := range login {
			for ip, tout := range ipmap {
				loginentry := x.LogEntry{
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
		p.RegisterUser = &x.DUGRegUnreg{
			Entry: make([]x.DUGEntry, len(group)),
		}
		entryidx := 0
		for user, grp := range group {
			dugentry := x.DUGEntry{
				User: user,
				Tag: x.Tag{
					Member: make([]x.Member, len(grp)),
				},
			}
			tagidx := 0
			for tag, tout := range grp {
				member := x.Member{
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
		p.Register = &x.DAGRegUnreg{
			Entry: make([]x.DAGEntry, len(reg)),
		}
		entryidx := 0
		for ip, tag := range reg {
			dagentry := x.DAGEntry{
				IP: ip,
				Tag: x.Tag{
					Member: make([]x.Member, len(tag)),
				},
			}
			tagidx := 0
			for tag, tout := range tag {
				member := x.Member{
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

func (mp mpayload) UIDMessage(m Monitor) (u *x.UIDMessage, err error) {
	var p *x.Payload
	if p, err = mp.Payload(m); err == nil {
		u = &x.UIDMessage{
			Type:    "update",
			Version: "2.0",
			Payload: p,
		}
	}
	return
}

func (mp mpayload) Push(
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

func (mp mpayload) log(m Monitor, op operation, subject string, value string, tout *uint) {
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
