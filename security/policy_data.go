package security

type PolicyData struct {
	Role string
	Perm Permission
	Prio int64
	Cond any
}
