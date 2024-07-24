package Node

import "time"

type MemberNode struct {
	Id       int
	Username string
	NoTelp   int
	Point    int
	CreateAt time.Time
}

type MemberLL struct {
	Member MemberNode
	Next   *MemberLL
}
