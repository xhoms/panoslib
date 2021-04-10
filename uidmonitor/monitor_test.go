package uidmonitor_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/xhoms/panoslib/uid"
	"github.com/xhoms/panoslib/uidmonitor"
)

func TestMonitor(t *testing.T) {
	now := time.Now()
	t1 := now.Add(10 * time.Minute)
	var tout uint = 9
	var err error
	c := uidmonitor.NewMemMonitor()
	if _, err = uid.NewUIDBuilder().
		LoginUser("a1@test.local", "1.1.1.1", nil).
		LoginUser("a1@test.local", "1.1.2.2", &tout).
		LoginUser("a2@test.local", "2.2.2.2", &tout).
		LoginUser("b1@test.local", "3.3.3.3", &tout).
		GroupUser("a1@test.local", "a", nil).
		GroupUser("a2@test.local", "a", &tout).
		GroupUser("b1@test.local", "b", &tout).
		RegisterIP("1.1.1.1", "good", nil).
		RegisterIP("2.2.2.2", "good", &tout).
		RegisterIP("3.3.3.3", "bad", &tout).
		Payload(c); err == nil {
		if len(c.UserIP("a1@test.local")) == 2 &&
			len(c.GroupIP("a")) == 3 &&
			len(c.TagIP("good")) == 2 {
			if _, err = uid.NewUIDBuilder().
				LogoutUser("a1@test.local", "1.1.2.2").
				UngroupUser("a1@test.local", "a").
				UnregisterIP("1.1.1.1", "good").
				Payload(c); err == nil {
				if len(c.UserIP("a1@test.local")) == 1 &&
					len(c.GroupIP("a")) == 1 &&
					len(c.TagIP("good")) == 1 &&
					len(c.UserIP("a2@test.local")) == 1 &&
					len(c.GroupIP("b")) == 1 &&
					len(c.TagIP("bad")) == 1 {
					c.CleanUp(t1)
					if len(c.UserIP("a2@test.local")) == 0 &&
						len(c.GroupIP("b")) == 0 &&
						len(c.TagIP("bad")) == 0 {
						return
					} else {
						err = errors.New("recovery error 3")
					}
				} else {
					err = errors.New("recovery error 2")
				}
			}
		} else {
			err = errors.New("recovery error 1")
		}
	}
	t.Error(err)
}

func ExampleNewMemMonitor() {
	now := time.Now()
	t1 := now.Add(10 * time.Second)
	var tout uint = 9
	var err error
	c := uidmonitor.NewMemMonitor()
	if _, err = uid.NewUIDBuilder().
		RegisterIP("1.1.1.1", "good", nil).
		RegisterIP("2.2.2.2", "good", &tout).
		RegisterIP("3.3.3.3", "good", &tout).
		Payload(c); err == nil {
		if len(c.TagIP("good")) == 3 { // at this point he have three entries
			if _, err = uid.NewUIDBuilder().
				UnregisterIP("3.3.3.3", "good").
				Payload(c); err == nil {
				if len(c.TagIP("good")) == 2 { // one entry lost due to explicit unregister
					/*
						purge at t1 (10 seconds after) will remove the entry
						that expires at 9 seconds
					*/
					c.CleanUp(t1)
					fmt.Println(len(c.TagIP("good")))
				}
			}
		}
	}
	// Output: 1
}

func ExampleMemMonitor_Dump() {
	var tout uint = 60
	c := uidmonitor.NewMemMonitor()
	if _, err := uid.NewUIDBuilder().
		RegisterIP("1.1.1.1", "good", &tout).
		LoginUser("foo@test.local", "1.1.1.1", &tout).
		GroupUser("admin", "foo@test.local", &tout).
		Payload(c); err == nil {
		fmt.Println(c.Dump())
	}
}

func ExampleMemMonitor_TagIP() {
	c := uidmonitor.NewMemMonitor()
	if _, err := uid.NewUIDBuilder().
		RegisterIP("1.1.1.1", "windows", nil).
		Payload(c); err == nil {
		fmt.Println(c.TagIP("windows"))
	}
	// Output: [1.1.1.1]
}

func ExampleMemMonitor_GroupIP() {
	c := uidmonitor.NewMemMonitor()
	if _, err := uid.NewUIDBuilder().
		LoginUser("foo@test.local", "1.1.1.1", nil).
		GroupUser("foo@test.local", "admin", nil).
		Payload(c); err == nil {
		fmt.Println(c.GroupIP("admin"))
	}
	// Output: [1.1.1.1]
}

func ExampleMemMonitor_UserIP() {
	c := uidmonitor.NewMemMonitor()
	if _, err := uid.NewUIDBuilder().
		LoginUser("foo@test.local", "1.1.1.1", nil).
		Payload(c); err == nil {
		fmt.Println(c.UserIP("foo@test.local"))
	}
	// Output: [1.1.1.1]
}
