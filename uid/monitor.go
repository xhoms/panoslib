package uid

import (
	"encoding/json"
	"sort"
	"strconv"
	"time"
)

var MSIZE = ""
var MTOUT = ""

type item struct {
	subject, key string
	Valid        int64
}

type index map[string]map[string]*item

func (i index) add(im *item) {
	if s, exists := i[im.key]; exists {
		s[im.subject] = im
	} else {
		i[im.key] = map[string]*item{im.subject: im}
	}
}

func (i index) rm(subject, key string) {
	if s, exists := i[key]; exists {
		delete(s, subject)
		if len(s) == 0 {
			delete(i, key)
		}
	}
}

func (i index) list(key string) (sub []string) {
	submap := i[key]
	sub = make([]string, len(submap))
	idx := 0
	for s := range submap {
		sub[idx] = s
		idx++
	}
	return
}

func (i index) get(subject, key string) (im *item) {
	if k, ke := i[key]; ke {
		im = k[subject]
	}
	return
}

type db struct {
	items []*item
	index index
	size  int
}

func newDb() (out *db) {
	size := 100
	if s, e := strconv.Atoi(MSIZE); e == nil {
		size = s
	}
	out = &db{
		items: make([]*item, 0, size),
		index: make(map[string]map[string]*item),
		size:  size,
	}
	return
}

func (d *db) Len() int {
	return len(d.items)
}

func (d *db) Less(i, j int) bool {
	return d.items[i].Valid < d.items[j].Valid
}

func (d *db) Swap(i, j int) {
	d.items[i], d.items[j] = d.items[j], d.items[i]
}

func (d *db) gb(t time.Time) {
	idx := 0
	sort.Sort(d)
	for idx = range d.items {
		if d.items[idx].Valid >= t.UnixNano() {
			idx--
			break
		}
	}
	idx++
	var imindex index = make(map[string]map[string]*item, len(d.items)-idx)
	for idx2 := idx; idx2 < len(d.items); idx2++ {
		imindex.add(d.items[idx2])
	}
	d.items, d.index = d.items[idx:], imindex
}

func (d *db) append(subject, key string, valid int64) {
	if im := d.index.get(subject, key); im != nil {
		im.Valid = valid
	} else {
		im := item{subject: subject, key: key, Valid: valid}
		d.index.add(&im)
		d.items = append(d.items, &im)
	}
}

func (d *db) remove(subject, key string) {
	if im := d.index.get(subject, key); im != nil {
		idx := 0
		found := false
		for idx = range d.items {
			if d.items[idx].subject == subject && d.items[idx].key == key {
				found = true
				break
			}
		}
		if found {
			d.index.rm(subject, key)
			d.items = append(d.items[:idx], d.items[idx+1:]...)
		}
	}
}

func (d *db) list(key string) (out []string) {
	out = d.index.list(key)
	return
}

type Mon struct {
	maxtout   time.Duration
	userMap   *db
	userGroup *db
	ipTag     *db
}

func NewMonitor() (m *Mon) {
	maxtout := time.Hour * 720
	if t, e := strconv.Atoi(MTOUT); e == nil {
		maxtout = time.Minute * time.Duration(t)
	}
	m = &Mon{
		maxtout:   maxtout,
		userMap:   newDb(),
		userGroup: newDb(),
		ipTag:     newDb(),
	}
	return
}

func (m *Mon) Log(op operation, subject, value string, tout *uint) {
	valid := time.Now()
	if tout == nil {
		valid = valid.Add(m.maxtout)
	} else {
		valid = valid.Add(time.Minute * time.Duration(*tout))
	}
	switch op {
	case Login:
		m.userMap.append(value, subject, valid.UnixNano())
	case Logout:
		m.userMap.remove(value, subject)
	case Group:
		m.userGroup.append(subject, value, valid.UnixNano())
	case Ungroup:
		m.userGroup.remove(subject, value)
	case Register:
		m.ipTag.append(subject, value, valid.UnixNano())
	case Unregister:
		m.ipTag.remove(subject, value)
	}
}

func (m *Mon) UserIP(user string) []string {
	return m.userMap.list(user)
}

func (m *Mon) GroupIP(group string) (out []string) {
	out = make([]string, 0, m.userMap.size)
	for u := range m.userGroup.index[group] {
		out = append(out, m.userMap.list(u)...)
	}
	return out
}

func (m *Mon) TagIP(tag string) []string {
	return m.ipTag.list(tag)
}

func (m *Mon) CleanUp(t time.Time) {
	m.userMap.gb(t)
	m.userGroup.gb(t)
	m.ipTag.gb(t)
}

func (m *Mon) Dump() (out string) {
	j := &struct {
		UserIp    index
		GroupUser index
		TagIP     index
	}{m.userMap.index, m.userGroup.index, m.ipTag.index}
	if o, err := json.MarshalIndent(j, "", "  "); err == nil {
		out = string(o)
	}
	return
}
