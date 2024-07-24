package Database

import (
	"THR/Node"
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var HeadItem = Node.ItemLL{}
var HeadPenjualan = Node.PenjualanLL{}
var HeadMember = Node.MemberLL{}


//connect sql
var DBConnect *sql.DB
var DBerr error

func Initialize() error {
	DBConnect, DBerr = sql.Open("mysql", "root:@tcp(localhost)/go")
	if DBerr != nil {
		return DBerr
	}

	if DBerr = DBConnect.Ping(); DBerr != nil {
		return DBerr
	}

	return nil
}

func CreateUser(user Node.NodeUser) (int64, error) {
	query := "INSERT INTO users (username, password, role) VALUES (?, ?, ?)"
	result, err := DBConnect.Exec(query, user.Username, user.Password, user.Role)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetUser(id int64) (Node.NodeUser, error) {
	query := "SELECT id, username, password, role FROM users WHERE id = ?"
	row := DBConnect.QueryRow(query, id)

	var user Node.NodeUser
	if err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func UpdateUser(user Node.NodeUser) error {
	query := "UPDATE users SET username = ?, password = ?, role = ? WHERE id = ?"
	_, err := DBConnect.Exec(query, user.Username, user.Password, user.Role, user.Id)
	return err
}

func DeleteUser(id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := DBConnect.Exec(query, id)
	return err
}

func GetAllUsers() ([]Node.NodeUser, error) {
	query := "SELECT id, username, password, role FROM users"
	rows, err := DBConnect.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []Node.NodeUser
	for rows.Next() {
		var user Node.NodeUser
		if err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}


