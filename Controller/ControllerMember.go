package Controller

import (
	"THR/Database"
	"THR/Model"
	"THR/Node"
	"database/sql"
	"errors"
	"fmt"
	"time"
)


func IsMemberSame(username string, noTelp int) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM members WHERE username = ? OR noTelp = ?)"
    err := Database.DBConnect.QueryRow(query, username, noTelp).Scan(&exists)
    if err != nil {
        return false, err
    }
    return exists, nil
}

func ValidasiInsertMember(username string, noTelp int, point int) (bool, error) {
    // Cek apakah data member sudah ada
    isDataSama, err := IsMemberSame(username, noTelp)
    if err != nil {
        return false, err
    }
    if isDataSama {
        return false, fmt.Errorf("member dengan username atau noTelp tersebut sudah ada")
    }

    // Menggunakan waktu sekarang dan mengonversinya ke zona waktu Surabaya
    loc, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        return false, err
    }
    setCreateAt := time.Now().In(loc)

    newMember := Node.MemberNode{
        Username: username,
        NoTelp:   noTelp,
        Point:    point,
        CreateAt: setCreateAt,
    }

    query := "INSERT INTO members (username, noTelp, point, createAt) VALUES (?, ?, ?, ?)"
    _, err = Database.DBConnect.Exec(query, newMember.Username, newMember.NoTelp, newMember.Point, newMember.CreateAt)
    if err != nil {
        return false, err
    }

    return true, nil
}



func GetAllMembers() ([]Node.MemberNode, error) {
    query := "SELECT id, username, noTelp, point, createAt FROM members"
    rows, err := Database.DBConnect.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var members []Node.MemberNode
    for rows.Next() {
        //var penampung
        var id int
        var username string
        var noTelp int
        var point int
        var createAt string // Baca sebagai string terlebih dahulu
    
        //pindahkan isi query ke var
        if err := rows.Scan(&id, &username, &noTelp, &point, &createAt); err != nil {
            return nil, err
        }
    
        // Konversi string ke time.Time
        parsedTime, err := time.Parse("2006-01-02 15:04:05", createAt)
        if err != nil {
            return nil, err
        }

        // Tentukan zona waktu Surabaya
        loc, err := time.LoadLocation("Asia/Jakarta")
        if err != nil {
            return nil, err
        }

        // Konversi waktu ke zona waktu Surabaya
        parsedTime =  parsedTime.In(loc)
    
        member := Node.MemberNode{
            Id:       id,
            Username: username,
            NoTelp:   noTelp,
            Point:    point,
            CreateAt: parsedTime,
        }
    
        //masukkan data member ke slice members
        members = append(members, member)
    }
    if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func ValidasiDeleteMember(id int) error {
	query := "DELETE FROM members WHERE id = ?"
	_, err := Database.DBConnect.Exec(query, id)
	return err
}

func ValidasiSearchMember(id int)(string,*Node.MemberLL){
	cur := &Database.HeadMember
	if(cur.Next != nil){
		result := Model.SearchMember(id)
		if(result != nil){
			return "Data Member Ditemukan!",result.Next
		}
		return "Tidak ada Data Member dengan id Tersebut!",nil
	}
	return "Data Member Masih Kosong!",nil
}

func ValidasiUpdateMember(member Node.MemberNode) error {
	query := "UPDATE members SET username = ?, noTelp = ?  WHERE id = ?"
	_, err := Database.DBConnect.Exec(query, member.Username,member.NoTelp,member.Id)
	return err
}

func GetMemberById(id int) (Node.MemberNode, error) {
	query := "SELECT id, username, noTelp, point, createAt FROM members WHERE id = ?"
	row := Database.DBConnect.QueryRow(query, id)

    var createAt string 
	var member Node.MemberNode
	if err := row.Scan(&member.Id, &member.Username, &member.NoTelp, &member.Point, &createAt); err != nil {
		if err == sql.ErrNoRows {
			return member, errors.New("member not found")
		}
     // Konversi string ke time.Time
    parsedTime, _ := time.Parse("2006-01-02 15:04:05", createAt)
 
    // Tentukan zona waktu Surabaya
    loc, _ := time.LoadLocation("Asia/Jakarta")
    
    // Konversi waktu ke zona waktu Surabaya
    parsedTime =  parsedTime.In(loc)
    member.CreateAt = parsedTime
		return member, err
	}
	return member, nil
}

//WEB
func ValidasiTambahMemberPoints(id int, pointReward int) (string, bool) {
    member, memberErr := GetMemberById(id)
    if memberErr != nil {
        return memberErr.Error(), false
    }

    pointReward += member.Point

    query := "UPDATE members SET point = ?  WHERE id = ?"
     Database.DBConnect.Exec(query, pointReward,member.Id)


    return "Member points updated successfully", true
}

func ValidasiKurangiMemberPoints(id int,  pointUsed int) (string, bool) {
    member, memberErr := GetMemberById(id)
    if memberErr != nil {
        return memberErr.Error(), false
    }

    member.Point -= pointUsed

    query := "UPDATE members SET point = ?  WHERE id = ?"
     Database.DBConnect.Exec(query, member.Point,member.Id)


    return "Member points updated successfully", true
}
