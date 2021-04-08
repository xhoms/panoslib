package uid_test

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	x "github.com/xhoms/panoslib/collection"
	"github.com/xhoms/panoslib/uid"
)

func TestSingleRegister(t *testing.T) {
	dag := []uid.IPTag{{IP: "1.1.1.1", Tag: "foo"}}
	mp := uid.NewUIDBuilder().Register(dag)
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil, p.Register == nil, len(p.Register.Entry) == 0:
					err = errors.New("nil unmarshal")
				case p.Register.Entry[0].IP != "1.1.1.1",
					len(p.Register.Entry[0].Tag.Member) != 1,
					p.Register.Entry[0].Tag.Member[0].Timeout != nil,
					p.Register.Entry[0].Tag.Member[0].Member != "foo":
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestSingleRegisterAlternate(t *testing.T) {
	mp := uid.NewUIDBuilder().RegisterIP("1.1.1.1", "foo", nil)
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil, p.Register == nil, len(p.Register.Entry) == 0:
					err = errors.New("nil unmarshal")
				case p.Register.Entry[0].IP != "1.1.1.1",
					len(p.Register.Entry[0].Tag.Member) != 1,
					p.Register.Entry[0].Tag.Member[0].Timeout != nil,
					p.Register.Entry[0].Tag.Member[0].Member != "foo":
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestSingleRegisterTout(t *testing.T) {
	var tout uint = 60
	dag := []uid.IPTag{{IP: "1.1.1.1", Tag: "foo", Tout: &tout}}
	mp := uid.NewUIDBuilder().Register(dag)
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil, p.Register == nil, len(p.Register.Entry) == 0:
					err = errors.New("nil unmarshal")
				case p.Register.Entry[0].IP != "1.1.1.1",
					len(p.Register.Entry[0].Tag.Member) != 1,
					p.Register.Entry[0].Tag.Member[0].Member != "foo",
					*p.Register.Entry[0].Tag.Member[0].Timeout != "60":
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestSingleRegisterToutAlternate(t *testing.T) {
	var tout uint = 60
	mp := uid.NewUIDBuilder().RegisterIP("1.1.1.1", "foo", &tout)
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil, p.Register == nil, len(p.Register.Entry) == 0:
					err = errors.New("nil unmarshal")
				case p.Register.Entry[0].IP != "1.1.1.1",
					len(p.Register.Entry[0].Tag.Member) != 1,
					p.Register.Entry[0].Tag.Member[0].Member != "foo",
					*p.Register.Entry[0].Tag.Member[0].Timeout != "60":
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestSingleUnregister(t *testing.T) {
	dag := []uid.IPTag{{IP: "1.1.1.1", Tag: "foo"}}
	mp := uid.NewUIDBuilder().Unregister(dag)
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil, p.Unregister == nil, len(p.Unregister.Entry) == 0:
					err = errors.New("nil unmarshal")
				case p.Unregister.Entry[0].IP != "1.1.1.1",
					len(p.Unregister.Entry[0].Tag.Member) != 1,
					p.Unregister.Entry[0].Tag.Member[0].Timeout != nil,
					p.Unregister.Entry[0].Tag.Member[0].Member != "foo":
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestSingleRegUnreg(t *testing.T) {
	dagreg := []uid.IPTag{{IP: "1.1.1.1", Tag: "foo"}}
	dagunreg := []uid.IPTag{{IP: "2.2.2.2", Tag: "bar"}}
	mp := uid.NewUIDBuilder().Register(dagreg).Unregister(dagunreg)
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil,
					p.Unregister == nil,
					p.Register == nil,
					len(p.Unregister.Entry) == 0,
					len(p.Register.Entry) == 0:
					err = errors.New("nil unmarshal")
				case p.Register.Entry[0].IP != "1.1.1.1",
					p.Unregister.Entry[0].IP != "2.2.2.2",
					len(p.Register.Entry[0].Tag.Member) != 1,
					len(p.Unregister.Entry[0].Tag.Member) != 1,
					p.Register.Entry[0].Tag.Member[0].Member != "foo",
					p.Register.Entry[0].Tag.Member[0].Timeout != nil,
					p.Unregister.Entry[0].Tag.Member[0].Member != "bar",
					p.Unregister.Entry[0].Tag.Member[0].Timeout != nil:
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestMultiRegUnreg(t *testing.T) {
	dagreg := []uid.IPTag{{IP: "1.1.1.1", Tag: "foo"}, {IP: "1.1.2.2", Tag: "foo"}, {IP: "1.1.2.2", Tag: "bar"}}
	dagunreg := []uid.IPTag{{IP: "2.2.2.2", Tag: "bar"}, {IP: "2.2.1.1", Tag: "bar"}, {IP: "2.2.1.1", Tag: "foo"}}
	mp := uid.NewUIDBuilder().Register(dagreg).Unregister(dagunreg)
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil,
					p.Unregister == nil,
					p.Register == nil,
					len(p.Unregister.Entry) == 0,
					len(p.Register.Entry) == 0:
					err = errors.New("nil unmarshal")
				case len(p.Register.Entry) != 2,
					len(p.Register.Entry[0].Tag.Member)+len(p.Register.Entry[1].Tag.Member) != 3,
					len(p.Unregister.Entry) != 2,
					len(p.Unregister.Entry[0].Tag.Member)+len(p.Unregister.Entry[1].Tag.Member) != 3:
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestSingleLogin(t *testing.T) {
	var tout uint = 60
	login := []uid.UserMap{{IP: "1.1.1.1", User: "foo@test.local", Tout: &tout}}
	mp := uid.NewUIDBuilder().Login(login)
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil,
					p.Login == nil,
					len(p.Login.Entry) == 0:
					err = errors.New("nil unmarshal")
				case p.Login.Entry[0].Name != "foo@test.local",
					p.Login.Entry[0].IP != "1.1.1.1",
					*p.Login.Entry[0].Timeout != "60":
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestSingleLoginAlternate(t *testing.T) {
	var tout uint = 60
	mp := uid.NewUIDBuilder().LoginUser("foo@test.local", "1.1.1.1", &tout)
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil,
					p.Login == nil,
					len(p.Login.Entry) == 0:
					err = errors.New("nil unmarshal")
				case p.Login.Entry[0].Name != "foo@test.local",
					p.Login.Entry[0].IP != "1.1.1.1",
					*p.Login.Entry[0].Timeout != "60":
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestSingleLogout(t *testing.T) {
	login := []uid.UserMap{{IP: "1.1.1.1", User: "foo@test.local"}}
	mp := uid.NewUIDBuilder().Logout(login)
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil,
					p.Logout == nil,
					len(p.Logout.Entry) == 0:
					err = errors.New("nil unmarshal")
				case p.Logout.Entry[0].Name != "foo@test.local",
					p.Logout.Entry[0].IP != "1.1.1.1",
					p.Logout.Entry[0].Timeout != nil:
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestSingleLogoutAlternate(t *testing.T) {
	mp := uid.NewUIDBuilder().LogoutUser("foo@test.local", "1.1.1.1")
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil,
					p.Logout == nil,
					len(p.Logout.Entry) == 0:
					err = errors.New("nil unmarshal")
				case p.Logout.Entry[0].Name != "foo@test.local",
					p.Logout.Entry[0].IP != "1.1.1.1",
					p.Logout.Entry[0].Timeout != nil:
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestLogMulti(t *testing.T) {
	mp := uid.NewUIDBuilder().
		LoginUser("foo@test.local", "1.1.1.1", nil).
		LoginUser("foo@test.local", "1.1.2.2", nil).
		LoginUser("bar@test.local", "1.1.3.3", nil).
		LogoutUser("foo@test.local", "1.1.4.4").
		LogoutUser("foo@test.local", "1.1.5.5").
		LogoutUser("bar@test.local", "1.1.6.6")
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil,
					p.Login == nil,
					len(p.Login.Entry) == 0,
					p.Logout == nil,
					len(p.Logout.Entry) == 0:
					err = errors.New("nil unmarshal")
				case len(p.Login.Entry) != 3,
					len(p.Logout.Entry) != 3:
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestSingleGroup(t *testing.T) {
	dug := []uid.UserGroup{{User: "foo@test.local", Group: "admin"}}
	mp := uid.NewUIDBuilder().Group(dug)
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil, p.RegisterUser == nil, len(p.RegisterUser.Entry) == 0:
					err = errors.New("nil unmarshal")
				case p.RegisterUser.Entry[0].User != "foo@test.local",
					len(p.RegisterUser.Entry[0].Tag.Member) != 1,
					p.RegisterUser.Entry[0].Tag.Member[0].Member != "admin":
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestSingleGroupAlternate(t *testing.T) {
	mp := uid.NewUIDBuilder().GroupUser("foo@test.local", "admin", nil)
	var err error
	var p *x.UIDMsgPayload
	var b []byte
	if p, err = mp.Payload(nil); err == nil {
		if b, err = xml.Marshal(p); err == nil {
			p = &x.UIDMsgPayload{}
			if err = xml.Unmarshal(b, p); err == nil {
				switch {
				case p == nil, p.RegisterUser == nil, len(p.RegisterUser.Entry) == 0:
					err = errors.New("nil unmarshal")
				case p.RegisterUser.Entry[0].User != "foo@test.local",
					len(p.RegisterUser.Entry[0].Tag.Member) != 1,
					p.RegisterUser.Entry[0].Tag.Member[0].Member != "admin":
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

func TestNil(t *testing.T) {
	var err error
	if payload, err := uid.NewUIDBuilder().Payload(nil); err == nil {
		if _, err = xml.MarshalIndent(payload, "", "  "); err == nil {
			return
		}
	}
	t.Error(err)
}

func TestMultiRegUnregUid(t *testing.T) {
	dagreg := []uid.IPTag{{IP: "1.1.1.1", Tag: "foo"}, {IP: "1.1.2.2", Tag: "foo"}, {IP: "1.1.2.2", Tag: "bar"}}
	dagunreg := []uid.IPTag{{IP: "2.2.2.2", Tag: "bar"}, {IP: "2.2.1.1", Tag: "bar"}, {IP: "2.2.1.1", Tag: "foo"}}
	mp := uid.NewUIDBuilder().Register(dagreg).Unregister(dagunreg)
	var err error
	var u *x.UIDMessage
	var b []byte
	if u, err = mp.UIDMessage(nil); err == nil {
		if b, err = xml.Marshal(u); err == nil {
			u = &x.UIDMessage{}
			if err = xml.Unmarshal(b, u); err == nil {
				switch {
				case u == nil,
					u.Payload == nil,
					u.Payload.Unregister == nil,
					u.Payload.Register == nil,
					len(u.Payload.Unregister.Entry) == 0,
					len(u.Payload.Register.Entry) == 0:
					err = errors.New("nil unmarshal")
				case len(u.Payload.Register.Entry) != 2,
					len(u.Payload.Register.Entry[0].Tag.Member)+len(u.Payload.Register.Entry[1].Tag.Member) != 3,
					len(u.Payload.Unregister.Entry) != 2,
					len(u.Payload.Unregister.Entry[0].Tag.Member)+len(u.Payload.Unregister.Entry[1].Tag.Member) != 3:
					err = errors.New("recovery error")
				default:
					return
				}
			}
		}
	}
	t.Error(err)
}

type client string

func (c client) Do(req *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(c))),
		StatusCode: http.StatusOK,
	}
	return
}

func TestMultiRegUnregPush(t *testing.T) {
	var respBody client = `
<response status="success">
	<result>
		<uid-response>
			<version>2.0</version>
			<payload>
				<unregister> </unregister>
				<register> </register>
			</payload>
		</uid-response>
	</result>
</response>
	`
	dagreg := []uid.IPTag{{IP: "1.1.1.1", Tag: "foo"}, {IP: "1.1.2.2", Tag: "foo"}, {IP: "1.1.2.2", Tag: "bar"}}
	dagunreg := []uid.IPTag{{IP: "2.2.2.2", Tag: "bar"}, {IP: "2.2.1.1", Tag: "bar"}, {IP: "2.2.1.1", Tag: "foo"}}
	mp := uid.NewUIDBuilder().Register(dagreg).Unregister(dagunreg)
	var err error
	var resp *http.Response
	if resp, err = mp.Push("vm.test.local", "apikey", respBody, nil); err == nil {
		var apiresp *x.APIResponse
		if apiresp, err = uid.Validate(resp, err); err == nil {
			if apiresp.Result != nil &&
				apiresp.Result.UidResponse != nil &&
				apiresp.Result.UidResponse.Version == "2.0" &&
				apiresp.Result.UidResponse.Payload.Register != nil &&
				apiresp.Result.UidResponse.Payload.Unregister != nil {
				return
			} else {
				err = errors.New("recovery error")
			}
		}
	}
	t.Error(err)
}

func TestMultiRegUnregPushErr(t *testing.T) {
	var respBody client = `
<response status="error">
	<msg>
		<line>
			<uid-response>
				<version>2.0</version>
				<payload>
					<unregister> </unregister>
					<register>
						<entry ip="1.1.2.2" message="tag foo already exists, ignore"/>
						<entry ip="1.1.2.2" message="tag bar already exists, ignore"/>
						<entry ip="1.1.1.1" message="tag foo already exists, ignore"/>
					</register>
				</payload>
			</uid-response>
		</line>
	</msg>
</response>
`
	dagreg := []uid.IPTag{{IP: "1.1.1.1", Tag: "foo"}, {IP: "1.1.2.2", Tag: "foo"}, {IP: "1.1.2.2", Tag: "bar"}}
	dagunreg := []uid.IPTag{{IP: "2.2.2.2", Tag: "bar"}, {IP: "2.2.1.1", Tag: "bar"}, {IP: "2.2.1.1", Tag: "foo"}}
	mp := uid.NewUIDBuilder().Register(dagreg).Unregister(dagunreg)
	var err error
	var resp *http.Response
	if resp, err = mp.Push("vm.test.local", "apikey", respBody, nil); err == nil {
		var apiresp *x.APIResponse
		if apiresp, _ = uid.Validate(resp, err); apiresp != nil {
			if _, err = xml.Marshal(apiresp); err == nil {
				if apiresp.Msg != nil &&
					len(apiresp.Msg.Line) == 1 &&
					apiresp.Msg.Line[0].UidResponse != nil &&
					apiresp.Msg.Line[0].UidResponse.Version == "2.0" &&
					apiresp.Msg.Line[0].UidResponse.Payload.Register != nil &&
					len(apiresp.Msg.Line[0].UidResponse.Payload.Register.Entry) == 3 &&
					apiresp.Msg.Line[0].UidResponse.Payload.Unregister != nil {
					return
				} else {
					err = errors.New("recovery error")
				}
			}
		} else {
			err = errors.New("nil api response")
		}
	}
	t.Error(err)
}

func TestFromPayloadRegister(t *testing.T) {
	raw := `
<payload>
	<register>
	  <entry ip="1.1.1.1">
		<tag>
		  <member timeout="60">win</member>
		</tag>
	  </entry>
	</register>
	<register-user>
	  <entry user="a@test.local">
		<tag>
		  <member timeout="60">admin</member>
		</tag>
	  </entry>
	</register-user>
	<login>
	  <entry name="a@test.local" ip="1.1.1.1" timeout="60"></entry>
	</login>
  </payload>
`
	payload := &x.UIDMsgPayload{}
	var tout uint = 61
	var err error
	if err = xml.Unmarshal([]byte(raw), payload); err == nil {
		mp := uid.NewBuilderFromPayload(payload).
			LoginUser("b@test.local", "1.1.2.2", &tout).
			GroupUser("b@test.local", "admin", &tout).
			RegisterIP("1.1.2.2", "win", &tout)
		if payload, err = mp.Payload(nil); err == nil {
			if payload != nil &&
				payload.Login != nil && len(payload.Login.Entry) == 2 &&
				payload.Register != nil && len(payload.Register.Entry) == 2 &&
				payload.RegisterUser != nil && len(payload.RegisterUser.Entry) == 2 {
				return
			} else {
				err = errors.New("recovery error")
			}
		}
	}
	t.Error(err)
}

func TestFromPayloadRegisterT(t *testing.T) {
	raw := `
<payload>
	<unregister>
	  <entry ip="1.1.1.1">
		<tag>
		  <member>win</member>
		</tag>
	  </entry>
	</unregister>
	<unregister-user>
	  <entry user="a@test.local">
		<tag>
		  <member>admin</member>
		</tag>
	  </entry>
	</unregister-user>
	<logout>
	  <entry name="a@test.local" ip="1.1.1.1"></entry>
	</logout>
  </payload>
`
	payload := &x.UIDMsgPayload{}
	var err error
	if err = xml.Unmarshal([]byte(raw), payload); err == nil {
		mp := uid.NewBuilderFromPayload(payload).
			LogoutUser("b@test.local", "1.1.2.2").
			UngroupUser("b@test.local", "admin").
			UnregisterIP("1.1.2.2", "win")
		if payload, err = mp.Payload(nil); err == nil {
			if payload != nil &&
				payload.Logout != nil && len(payload.Logout.Entry) == 2 &&
				payload.Unregister != nil && len(payload.Unregister.Entry) == 2 &&
				payload.UnregisterUser != nil && len(payload.UnregisterUser.Entry) == 2 {
				return
			} else {
				err = errors.New("recovery error")
			}
		}
	}
	t.Error(err)
}
