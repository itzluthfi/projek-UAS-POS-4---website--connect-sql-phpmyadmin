package route

import (
	"THR/Handler"
	"log"
	"net/http"
)

func RouteMember() {
    http.HandleFunc("/manageMember", Handler.ViewHandlerMember)
    http.HandleFunc("/insertMember", Handler.InsertMemberHandler)
    http.HandleFunc("/updateMember", Handler.UpdateMemberHandler)
    http.HandleFunc("/deleteMember", Handler.DeleteMemberHandler)
    http.HandleFunc("/getMemberDetails", Handler.GetMemberDetailsHandler)
    http.HandleFunc("/TambahMemberPoints", Handler.TambahMemberPointsHandler)
    http.HandleFunc("/KurangiMemberPoints", Handler.KurangiMemberPointsHandler)
}

func RouteItem() {
    http.HandleFunc("/manageItem", Handler.ViewHandlerItem)
    http.HandleFunc("/insertItem", Handler.InsertItemHandler)
    http.HandleFunc("/updateItem", Handler.UpdateItemHandler)
    http.HandleFunc("/deleteItem", Handler.DeleteItemHandler)
    http.HandleFunc("/getItemDetails", Handler.GetItemDetailsHandler)
    http.HandleFunc("/getAllItems", Handler.GetAllItemsHandler)
    http.HandleFunc("/TambahItemStock", Handler.TambahItemStockHandler)
    http.HandleFunc("/KurangiItemStock", Handler.KurangiItemStockHandler)
}

func RouteUser() {
    http.HandleFunc("/manageUser", Handler.ViewHandlerUser)
    http.HandleFunc("/insertUser", Handler.InsertUserHandler)
    http.HandleFunc("/updateUser", Handler.UpdateUserHandler)
    http.HandleFunc("/deleteUser", Handler.DeleteUserHandler)
}

func RouteKasir() {
    http.HandleFunc("/kasirNonMember", Handler.ViewKasirNonMemberHandler)
    http.HandleFunc("/kasirMember", Handler.ViewKasirMemberHandler)
    http.HandleFunc("/historyPenjualan", Handler.ViewHistoryPenjualanHandler)

}

func RoutePenjualan() {
    http.HandleFunc("/recordSale", Handler.RecordSaleHandler) 
    http.HandleFunc("/getSalesHistory", Handler.HandleGetSalesHistory)
	http.HandleFunc("/deletePenjualan", Handler.HandleDeletePenjualan)
	http.HandleFunc("/manageHistoryPenjualan", Handler.ViewManageHistoryPenjualanHandler)
    http.HandleFunc("/getDetailPenjualan", Handler.GetDetailPenjualanHandler)
}


func Auth(){
    http.HandleFunc("/", Handler.HomeHandler)
    http.HandleFunc("/login", Handler.LoginHandler)
    http.HandleFunc("/logout", Handler.LogoutHandler)
}

func RunServer() {
    // Serve static files
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    
    // Register routes 
    Auth()
    RouteMember()
    RouteItem()
    RouteUser()
    RouteKasir()
    RoutePenjualan()

    // Start the server on port 8080
    log.Println("Server berjalan pada http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Gagal menjalankan server: %v", err)
    }
}
