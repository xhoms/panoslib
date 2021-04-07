package uid

type IPTag struct {
	IP, Tag string
	Tout    *uint
}

type UserMap struct {
	IP, User string
	Tout     *uint
}

type UserGroup struct {
	User, Group string
	Tout        *uint
}
