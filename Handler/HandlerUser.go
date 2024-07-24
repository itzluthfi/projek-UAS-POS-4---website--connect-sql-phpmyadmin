package Handler

import (
	"THR/Controller"
	"THR/Node"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var (
	tmplUser        = template.Must(template.ParseFiles("View/manageUser/manageUser.html"))
	loginTmpl       = template.Must(template.ParseFiles("View/login/login.html"))
	insertUserTmpl  = template.Must(template.ParseFiles("View/manageUser/insert.html"))
	updateUserTmpl  = template.Must(template.ParseFiles("View/manageUser/update.html"))	
)

func ViewHandlerUser(w http.ResponseWriter, r *http.Request) {
    users,userErr := Controller.GetAllUsers()
	if userErr != nil {
		fmt.Println("gagal mengambil data users")
		return
	}

    var username, role string
	var userId int
    if cookie, err := r.Cookie("username"); err == nil {
        username = cookie.Value
        user := Controller.GetUserByUsername(username)
        role = user.Role
		userId = user.Id
    }

    data := struct {
        Users    []Node.NodeUser
		Id       int
        Username string
        Role     string
        Success  bool
    }{
        Users:    users,
		Id:       userId,
        Username: username,
        Role:     role,
        Success:  r.URL.Query().Get("login") == "success",
    }

    tmplUser.ExecuteTemplate(w, "manageUser.html", data)
}

func InsertUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Proses form submission untuk insert user
		username := r.FormValue("username")
		password := r.FormValue("password")
		role := r.FormValue("role")

		// Panggil fungsi controller untuk insert user
		 Controller.InsertUser(username, password, role)
		

		// Redirect ke halaman manageUser setelah berhasil insert
		http.Redirect(w, r, "/manageUser", http.StatusSeeOther)
		return
	}

	// Jika bukan method POST, tampilkan form insert user
	insertUserTmpl.Execute(w, nil)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("Id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	user,_ := Controller.GetUserById(userId)

	
	
	if r.Method == http.MethodGet {
		if err := updateUserTmpl.Execute(w, user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		
		// Handle submission form update
		oldUsername := r.FormValue("oldUsername")
		newUsername := r.FormValue("username")
		newPassword := r.FormValue("password")
		newRole := r.FormValue("role")

		// Debugging: Log values
		log.Printf("Updating user: oldUsername=%s, newUsername=%s, password=%s, role=%s\n", oldUsername, newUsername, newPassword, newRole)

		newUser := Node.NodeUser{
			Id: userId, 
			Username: newUsername,
			Password: newPassword,
			Role: newRole,
		}
		// Panggil fungsi controller untuk update data pengguna
		Controller.UpdateUser(newUser)

		// Redirect kembali ke halaman manageUser setelah update
		http.Redirect(w, r, "/manageUser", http.StatusSeeOther)
	}
}


func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("Id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	Controller.DeleteUser(userId)
	http.Redirect(w, r, "/manageUser", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")
        if valid, role := Controller.VerifikasiUser(username, password); valid {
            // Simpan username di cookie
            http.SetCookie(w, &http.Cookie{
                Name:  "username",
                Value: username,
                Path:  "/",
            })
            
            // Redirect based on role
            if role == "admin" {
                http.Redirect(w, r, "/manageUser?login=success", http.StatusSeeOther)
            } else if role == "kasir" {
                http.Redirect(w, r, "/kasirNonMember?login=success", http.StatusSeeOther)
            }
            return
        } else {
            loginTmpl.Execute(w, map[string]string{"Error": "Username atau password salah"})
            return
        }
    }
    loginTmpl.Execute(w, nil)
}




func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Render template atau tampilkan halaman home
	homeTmpl := template.Must(template.ParseFiles("View/login/login.html"))
	homeTmpl.Execute(w, nil)
}


func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Logika logout seperti menghapus sesi atau token di sini jika diperlukan
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
