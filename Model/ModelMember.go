package Model

import (
	"THR/Database"
	"THR/Node"
)

// func MemberInsert(username string,  noTelp string,point int) {
// 	tempLL := &Database.HeadMember

// 	member := Node.MemberNode{
// 		Id:       getMemberLastId(),
// 		Username: username,
// 		NoTelp:   noTelp,
// 		Point: point,
// 		CreateAt: time.Now().Format("2006-01-02 15:04:05"),
// 	}

// 	newLL := &Node.MemberLL{
// 		Member: member,
// 		Next:   nil,
// 	}

// 	if tempLL.Next == nil {
// 		tempLL.Next = newLL
// 	} else {
// 		for tempLL.Next != nil {
// 			tempLL = tempLL.Next
// 		}
// 		tempLL.Next = newLL
// 	}
// }


func getMemberLastId() int {
	var tempLL *Node.MemberLL
	tempLL = &Database.HeadMember

	if tempLL.Next == nil {
		return 1
	} else {
		for tempLL.Next != nil {
			tempLL = tempLL.Next
		}
		return tempLL.Member.Id + 1
	}
}

func MemberReadAll() []Node.MemberNode {
	var tempLL *Node.MemberLL
	tempLL = &Database.HeadMember
	var memberTable []Node.MemberNode
	for tempLL.Next != nil {
		tempLL = tempLL.Next
		memberTable = append(memberTable, tempLL.Member)
	}
	return memberTable
}

func SearchMember(id int) *Node.MemberLL {
	var tempLL *Node.MemberLL
	tempLL = &Database.HeadMember
	if tempLL.Next != nil {
		for tempLL.Next != nil {
			if tempLL.Next.Member.Id == id {
				//return tempLL 
				return tempLL
			}
			tempLL = tempLL.Next
		}
	}
	return nil
}

func SearchMemberWeb(id int) *Node.MemberNode {
	var tempLL *Node.MemberLL
	tempLL = &Database.HeadMember
	if tempLL.Next != nil {
		for tempLL.Next != nil {
			if tempLL.Next.Member.Id == id {
				return &tempLL.Next.Member // Return the actual member object
			}
			tempLL = tempLL.Next
		}
	}
	return nil
}

func MemberDelete(address *Node.MemberLL) {
	address.Next = address.Next.Next
}

// func MemberUpdate(address *Node.MemberLL,username string,noTelp string)  {	
// 		address.Member.Username = username
// 		address.Member.NoTelp = noTelp
// }

// //koin = 5% dari total Harga
// func MemberUpdatePoint(addressMember *Node.MemberLL, totalHarga int) {
// 	addressMember.Member.Point += (totalHarga *  5/100)	
// }

func MemberUpdatePoint(member *Node.MemberLL, newPoints int) {
    member.Member.Point += newPoints
}

func TambahMemberPoint(member *Node.MemberLL, newPoints int) {
    member.Member.Point += newPoints
}
func KurangiMemberPoint(member *Node.MemberLL, newPoints int) {
    member.Member.Point -= newPoints
}

