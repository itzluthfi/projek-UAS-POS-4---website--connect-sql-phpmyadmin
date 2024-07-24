package Controller

import (
	"THR/Database"
	"THR/Node"
	"database/sql"
	"errors"
	"fmt"
)

// var Users = []Node.NodeUser{
// 	{1,"Admin", "admin123", "admin"},
// 	{2,"Adam", "adam123", "admin"},
// 	{3,"Habib", "kasir123", "kasir"},
// 	{4,"Luthfi", "kasir321", "kasir"},
// }

func GetAllUsers() ([]Node.NodeUser, error) {
	query := "SELECT id, username, password, role FROM users"
	rows, err := Database.DBConnect.Query(query)
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

func GetUserById(id int) (Node.NodeUser, error) {
	query := "SELECT id, username, password, role FROM users WHERE id = ?"
	row := Database.DBConnect.QueryRow(query, id)

	var user Node.NodeUser
	if err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func InsertUser(username, password, role string)(int64, error) {
	newUser := Node.NodeUser{
		Username: username,
		Password: password,
		Role:     role,
	}

	// allUsers,getUserErr := GetAllUsers()
	// if getUserErr != nil {
	// 	fmt.Println("gagal get all users")
	// }
	//allUsers = append(allUsers, newUser)
	query := "INSERT INTO users (username, password, role) VALUES (?, ?, ?)"
	result, err := Database.DBConnect.Exec(query, newUser.Username, newUser.Password, newUser.Role)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}



func GetUserByUsername(username string) Node.NodeUser {
	allUsers,_ := GetAllUsers()
	for _, user := range allUsers {
		if user.Username == username {
			return user
		}
	}
	return Node.NodeUser{}
}



func UpdateUser(user Node.NodeUser) error {
	query := "UPDATE users SET username = ?, password = ?, role = ? WHERE id = ?"
	_, err := Database.DBConnect.Exec(query, user.Username, user.Password, user.Role, user.Id)
	return err
}



func DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := Database.DBConnect.Exec(query, id)
	return err
}

func VerifikasiUser(username, password string) (bool, string) {
	allUsers,err := GetAllUsers()
	if err != nil {
		fmt.Println("gagal mengambil users")
	}
	for _, user := range allUsers {
		if user.Username == username && user.Password == password {
			return true, user.Role
		}
	}
	return false, ""
}
