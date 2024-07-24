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
	tmplItem       = template.Must(template.ParseFiles("View/manageItem/manageItem.html"))
	insertItemTmpl = template.Must(template.ParseFiles("View/manageItem/insert.html"))
	updateItemTmpl = template.Must(template.ParseFiles("View/manageItem/update.html"))
)

func GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	AllItems,itemsErr := Controller.GetAllItems()
	if itemsErr != nil {
		fmt.Println(itemsErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AllItems)
}


func ViewHandlerItem(w http.ResponseWriter, r *http.Request) {
	items,itemsErr := Controller.GetAllItems()
	if itemsErr != nil {
		fmt.Println(itemsErr.Error())
		return 
	}

	var username, role string
	if cookie, err := r.Cookie("username"); err == nil {
		username = cookie.Value
		user := Controller.GetUserByUsername(username)
		role = user.Role
	}

	data := struct {
		Items    []Node.NodeItem
		Username string
		Role     string
	}{
		Items:    items,
		Username: username,
		Role:     role,
	}

	tmplItem.ExecuteTemplate(w, "manageItem.html", data)
}


func InsertItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nama := r.FormValue("nama")
		jmlStock, _ := strconv.Atoi(r.FormValue("jmlStock"))
		harga, _ := strconv.Atoi(r.FormValue("harga"))
		diskon, _ := strconv.Atoi(r.FormValue("diskon"))

		cek,err := Controller.ValidasiInsertItem(nama,jmlStock,harga,diskon)
		if cek {
			http.Redirect(w, r, "/manageItem", http.StatusSeeOther)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	insertItemTmpl.Execute(w, nil)
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	// Parse URL parameters
	r.ParseForm()
	id := r.Form.Get("id")
	idInt, _ := strconv.Atoi(id)
	item,itemErr := Controller.GetItemById(idInt)
	if itemErr != nil{
		http.Error(w, itemErr.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		// Display the update form with item data
		if err := updateItemTmpl.Execute(w, item); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		// Handle form submission
		r.ParseForm()
		nama := r.Form.Get("nama")
		jmlStock, _ := strconv.Atoi(r.Form.Get("jmlStock"))
		harga, _ := strconv.Atoi(r.Form.Get("harga"))
		diskon, _ := strconv.Atoi(r.Form.Get("diskon"))

		// Call the controller to update data
		err := Controller.UpdateItem(nama, jmlStock, harga, diskon, idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/manageItem", http.StatusSeeOther)
	}
}

func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := Controller.DeleteItem(id); err == nil {
		http.Redirect(w, r, "/manageItem", http.StatusSeeOther)
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}



