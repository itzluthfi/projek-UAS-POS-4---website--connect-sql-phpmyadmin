package Handler

import (
	"THR/Controller"
	"THR/Node"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var (
    tmpl       = template.Must(template.ParseFiles("View/manageMember/manageMember.html"))
    insertTmpl = template.Must(template.ParseFiles("View/manageMember/insert.html"))
    updateTmpl = template.Must(template.ParseFiles("View/manageMember/update.html"))
    //loginForm  = template.Must(template.ParseFiles("View/login/login.html"))
)


func ViewHandlerMember(w http.ResponseWriter, r *http.Request) {
    members, err := Controller.GetAllMembers()
    if err != nil {
        fmt.Println("Gagal mendapatkan semua anggota:", err)
        http.Error(w, "Gagal mendapatkan semua anggota", http.StatusInternalServerError)
        return
    }

    var username, role string
    if cookie, err := r.Cookie("username"); err == nil {
        username = cookie.Value
        user := Controller.GetUserByUsername(username)
        role = user.Role
    }

    data := struct {
        Members  []Node.MemberNode
        Username string
        Role     string
    }{
        Members:  members,
        Username: username,
        Role:     role,
    }

    if err := tmpl.ExecuteTemplate(w, "manageMember.html", data); err != nil {
        fmt.Println("Gagal mengeksekusi template:", err)
        http.Error(w, "Gagal mengeksekusi template", http.StatusInternalServerError)
    }
}



func InsertMemberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		noTelpStr := r.FormValue("noTelp")
        noTelp, err := strconv.Atoi(noTelpStr)
        if err != nil {
			http.Error(w, "No Telp harus berupa angka", http.StatusBadRequest)
			return
		}
		pointStr := r.FormValue("point")
		point, err := strconv.Atoi(pointStr)
        if err != nil {
			http.Error(w, "Point harus berupa angka", http.StatusBadRequest)
			return
		}

        
		if cek,_ := Controller.ValidasiInsertMember(username, noTelp, point); cek {
			http.Redirect(w, r, "/manageMember", http.StatusSeeOther)
		} else {
            fmt.Println(cek)
			http.Error(w, "Data tidak valid atau pengguna sudah ada", http.StatusBadRequest)
		}
		return
	}
	insertTmpl.Execute(w, nil)
}

func UpdateMemberHandler(w http.ResponseWriter, r *http.Request) {
    // Parse URL parameters
    r.ParseForm()
    id := r.Form.Get("id")
    idInt, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "ID tidak valid", http.StatusBadRequest)
        return
    }

    member, errGetmember := Controller.GetMemberById(idInt)
    if errGetmember != nil {
        fmt.Println("Gagal mendapatkan member!")
        http.Error(w, "Gagal mendapatkan member", http.StatusInternalServerError)
        return
    }

    if r.Method == "GET" {
        // Display the update form with member data
        if err := updateTmpl.Execute(w, member); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    } else if r.Method == "POST" {
        // Handle form submission
        r.ParseForm()
        username := r.Form.Get("username")
        noTelpStr := r.Form.Get("noTelp")
        noTelp, err := strconv.Atoi(noTelpStr)
        if err != nil {
            http.Error(w, "noTelp harus berupa angka", http.StatusBadRequest)
            return
        }

        member := Node.MemberNode{
            Id:       idInt,
            Username: username,
            NoTelp:   noTelp,
        }

        // Call the controller to update data
        err2 := Controller.ValidasiUpdateMember(member)
        if err2 != nil {
            http.Error(w, err2.Error(), http.StatusBadRequest)
            return
        }
        http.Redirect(w, r, "/manageMember", http.StatusSeeOther)
    }
}



func DeleteMemberHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err2 := Controller.ValidasiDeleteMember(id); err2 == nil {
		http.Redirect(w, r, "/manageMember", http.StatusSeeOther)
		
	} else {
		fmt.Println("gagal menghapus data member")
	}
}

//WEB
func GetMemberDetailsHandler(w http.ResponseWriter, r *http.Request) {
    memberIdStr := r.URL.Query().Get("id")
    memberId, err := strconv.Atoi(memberIdStr)
    if err != nil {
        http.Error(w, "Invalid member ID", http.StatusBadRequest)
        return
    }

    member, memberErr := Controller.GetMemberById(memberId)
    if memberErr != nil {
        http.Error(w, memberErr.Error(), http.StatusNotFound)
        return
    }

    jsonResponse, err2 := json.Marshal(member)
    if err2 != nil {
        http.Error(w, "Error converting member data to JSON", http.StatusInternalServerError)
        return
    }

    fmt.Println("JSON Response:", string(jsonResponse))  // Logging data JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

func TambahMemberPointsHandler(w http.ResponseWriter, r *http.Request) {
    // Mendapatkan ID anggota dari parameter URL
    memberIdStr := r.URL.Query().Get("id")
    memberId, err := strconv.Atoi(memberIdStr)
    if err != nil {
        http.Error(w, "ID anggota tidak valid", http.StatusBadRequest)
        return
    }

    // Mendapatkan jumlah poin dari parameter URL
    poinStr := r.URL.Query().Get("poin")
    pointReward, err := strconv.Atoi(poinStr)
    if err != nil {
        http.Error(w, "Nilai poin tidak valid", http.StatusBadRequest)
        return
    }

    // Validasi dan tambahkan poin anggota
    message, success := Controller.ValidasiTambahMemberPoints(memberId, pointReward)
    if !success {
        http.Error(w, message, http.StatusBadRequest)
        return
    }

    // Menyiapkan respons JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    response := map[string]string{"message": "Poin berhasil diperbarui"}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Error saat mengkodekan respons JSON", http.StatusInternalServerError)
        return
    }
    w.Write(jsonResponse)
}

func KurangiMemberPointsHandler(w http.ResponseWriter, r *http.Request) {
    // Mendapatkan ID anggota dari parameter URL
    memberIdStr := r.URL.Query().Get("id")
    memberId, err := strconv.Atoi(memberIdStr)
    if err != nil {
        http.Error(w, "ID anggota tidak valid", http.StatusBadRequest)
        return
    }

    // Mendapatkan jumlah poin yang akan dikurangi dari parameter URL
    poinStr := r.URL.Query().Get("poin")
    pointUsed, err := strconv.Atoi(poinStr)
    if err != nil {
        http.Error(w, "Nilai poin tidak valid", http.StatusBadRequest)
        return
    }

    // Validasi dan kurangi poin anggota
    message, success := Controller.ValidasiKurangiMemberPoints(memberId, pointUsed)
    if !success {
        http.Error(w, message, http.StatusBadRequest)
        return
    }

    // Menyiapkan respons JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    response := map[string]string{"message": "Poin berhasil dikurangi"}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Error saat mengkodekan respons JSON", http.StatusInternalServerError)
        return
    }
    w.Write(jsonResponse)
}


