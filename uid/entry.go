package uid

// IPTag is a convenience struct to create a list of ip-to-tag UserID entries
type IPTag struct {
	IP, Tag string
	Tout    *uint
}

// UserMap is a convenience struct to create a list of user-to-ip UserID entries
type UserMap struct {
	IP, User string
	Tout     *uint
}

// UserGroup is a convenience struct to create a list of user-to-group UserID entries
type UserGroup struct {
	User, Group string
	Tout        *uint
}
