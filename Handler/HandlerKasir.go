package Handler

import (
	"THR/Controller"
	"THR/Node"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
)

var (
	kasirTmpl           = template.Must(template.ParseFiles("View/kasir/kasirNonMember.html"))
	kasirMemberTmpl     = template.Must(template.ParseFiles("View/kasir/kasirMember.html"))
	historyPenjualanTmpl = template.Must(template.ParseFiles("View/kasir/historyPenjualan.html"))
	
)

func ViewHistoryPenjualanHandler(w http.ResponseWriter, r *http.Request) {
	var username, role string
	if cookie, err := r.Cookie("username"); err == nil {
		username = cookie.Value
		user := Controller.GetUserByUsername(username)
		role = user.Role
	}

	// Ambil data penjualan dari linked list atau database
	salesHistory := Controller.GetSalesHistory() // Fungsi untuk mengambil riwayat penjualan

	data := struct {
		Username string
		Role     string
		Sales    []Node.NodePenjualan
	}{
		Username: username,
		Role:     role,
		Sales:    salesHistory,
	}

	historyPenjualanTmpl.Execute(w, data)
}

func ViewKasirMemberHandler(w http.ResponseWriter, r *http.Request) {
	var username, role string
	if cookie, err := r.Cookie("username"); err == nil {
		username = cookie.Value
		user := Controller.GetUserByUsername(username)
		role = user.Role
	}

	data := struct {
		Username string
		Role     string
	}{
		Username: username,
		Role:     role,
	}

	kasirMemberTmpl.Execute(w, data)
}



func ViewKasirNonMemberHandler(w http.ResponseWriter, r *http.Request) {
	var username, role string
	if cookie, err := r.Cookie("username"); err == nil {
		username = cookie.Value
		user := Controller.GetUserByUsername(username)
		role = user.Role
	}

	data := struct {
		Username string
		Role     string
	}{
		Username: username,
		Role:     role,
	}

	kasirTmpl.Execute(w, data)
}

func GetItemDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idItemStr := r.URL.Query().Get("id")
	idItem, err := strconv.Atoi(idItemStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	item, itemErr := Controller.GetItemById(idItem)
	if itemErr != nil {
		http.Error(w, itemErr.Error(), http.StatusNotFound)
		return
	}

	response := struct {
		ID          int    `json:"id"`
		Nama        string `json:"nama"`
		JmlStock    int    `json:"jmlStock"`
		Harga       int    `json:"harga"`
		Diskon      int    `json:"diskon"`
		HargaDiskon int    `json:"hargaDiskon"`
		JmlTerjual int    `json:"jmlTerjual"`
	}{
		ID:          item.Id,
		Nama:        item.Nama,
		JmlStock:    item.JmlStock,
		Harga:       item.Harga,
		Diskon:      item.Diskon,
		HargaDiskon: item.HargaDiskon,
		JmlTerjual:  item.JmlTerjual,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertPenjualanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		idItems := r.Form["idItem[]"]
		jmlPesanans := r.Form["jmlPesanan[]"]
		idMember, _ := strconv.Atoi(r.FormValue("idMember"))
		jmlTunai, _ := strconv.Atoi(r.FormValue("jmlTunai"))
		jmlPoint, _ := strconv.Atoi(r.FormValue("jmlPoint"))

		var detailItems []Node.NodeDetailPenjualan
		var totalDiskon int
		for i, idItemStr := range idItems {
			idItem, _ := strconv.Atoi(idItemStr)
			jmlPesanan, _ := strconv.Atoi(jmlPesanans[i])
			message, items := Controller.ValidasiPilihItem(idItem, jmlPesanan)
			if items == nil {
				http.Error(w, message, http.StatusBadRequest)
				return
			}
			detailItems = append(detailItems, items...)
			for _, item := range items {
				totalDiskon += item.Harga - item.HargaDiskon * item.JmlPesanan
			}
		}

		var addressMember *Node.MemberLL
		if idMember != 0 {
			addressMember = Controller.ValidasiIsMember(idMember)
			if addressMember == nil {
				http.Error(w, "Member tidak ditemukan", http.StatusBadRequest)
				return
			}
		}

		message, success := Controller.ValidasiInsertPenjualan(addressMember, detailItems, jmlTunai, 0, totalDiskon, jmlPoint)
		if !success {
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(message))
	} else {
		http.Redirect(w, r, "/kasirNonMember", http.StatusSeeOther)
	}
}

func TambahItemStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
    
	idItemStr := r.FormValue("id")
	idItem, err := strconv.Atoi(idItemStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	jumlahStr := r.FormValue("jumlah")
	jumlah, err := strconv.Atoi(jumlahStr)
	if err != nil || jumlah <= 0 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	err2 := Controller.ValidasiTambahStokItem(idItem, jumlah)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Stock updated successfully"}`))
}

func KurangiItemStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idItemStr := r.FormValue("id")
	idItem, err := strconv.Atoi(idItemStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	jumlahStr := r.FormValue("jumlah")
	jumlah, err := strconv.Atoi(jumlahStr)
	if err != nil || jumlah <= 0 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	err2 := Controller.ValidasiKurangiStokItem(idItem, jumlah)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Stock updated successfully"}`))
}
