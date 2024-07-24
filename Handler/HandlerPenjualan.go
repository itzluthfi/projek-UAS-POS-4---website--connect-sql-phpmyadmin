package Handler

import (
	"THR/Controller"
	"THR/Node"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type SaleDetail struct {
	NamaItem    string `json:"namaItem"`
	Jumlah      int    `json:"jumlah"`
	Harga       int    `json:"harga"`
	Diskon      int    `json:"diskon"`
	HargaDiskon int    `json:"hargaDiskon"`
	Subtotal    int    `json:"subtotal"`
}

type SaleData struct {
	TotalHarga  int         `json:"totalHarga"`
	TotalDiskon int         `json:"totalDiskon"`
	JumlahTunai int         `json:"jumlahTunai"`
	Kembalian   int         `json:"kembalian"`
	Details     []SaleDetail `json:"details"`
	IsMember    bool        `json:"isMember"`
	Tanggal     string      `json:"tanggal"`
    PointUsed   int         `json:"pointUsed"`
	PointReward int         `json:"pointReward"`  
    NamaMember  string      `json:"namaMember"`

}


var (
	manageHistoryPenjualanTmpl = template.Must(template.ParseFiles("View/managePenjualan/manageHistoryPenjualan.html"))
)

func ReadAllPenjualan() bool {
	AllPenjualan := Controller.GetSalesHistory()
	if AllPenjualan == nil {
		fmt.Println("Data Penjualan Masih Kosong!")
		return false
	} else {
		fmt.Println("=============================================================")
		fmt.Println("                  Laporan Data Semua Penjualan               ")
		fmt.Println("=============================================================")
		for _, Penjualan := range AllPenjualan {
			fmt.Printf("ID Penjualan   : %d\n", Penjualan.IdPenjualan)
			fmt.Printf("Is Member      : %t\n", Penjualan.IsMember)
			fmt.Printf("Create At      : %s\n", Penjualan.CreateAt)
			fmt.Println("------------------------------------------------------------")
			fmt.Println("                     Detail Item :                        ")
			fmt.Println("------------------------------------------------------------")
			fmt.Printf("| %-10s | %-12s | %-10s | %-12s | %-15s\n","Nama Item", "Harga satuan", "jumlah ","Diskon","sub total ")
			fmt.Println("------------------------------------------------------------")
			for _, item := range Penjualan.Detail {
				fmt.Printf("| %-10s | %-12d | %-10d | %-12d | %-15d\n", item.Nama, item.Harga, item.JmlPesanan,item.HargaDiskon,item.JmlPesanan*(item.Harga-(item.Harga*item.Diskon/100)))
			}
			fmt.Println("------------------------------------------------------------")
			fmt.Printf("Total Harga   \t\t\t\t : RP %d\n", Penjualan.Total)
			fmt.Println("------------------------------------------------------------")
			//detail Point
			fmt.Println("------------------------------------------------------------")
			fmt.Printf("Nama Member  \t\t\t\t : %s\n", Penjualan.NamaMember)
			fmt.Printf("Point Reward  \t\t\t\t : +RP %d\n", Penjualan.PointReward)
			fmt.Printf("Point Used  \t\t\t\t : RP %d\n", Penjualan.PointUsed)
			fmt.Println("------------------------------------------------------------")

			// Tampilkan detail pembayaran
				fmt.Printf("Jumlah Tunai  \t\t\t\t : RP %d\n", Penjualan.JmlTunai)
				fmt.Printf("Kembalian     \t\t\t\t : RP %d\n", Penjualan.Kembalian)
				fmt.Printf("Hemat          \t\t\t\t :(RP %d)\n", Penjualan.TotalDiskon)
			fmt.Println("============================================================")
		}
		return true
	}
}

func ViewManageHistoryPenjualanHandler(w http.ResponseWriter, r *http.Request) {
	var username, role string
	if cookie, err := r.Cookie("username"); err == nil {
		username = cookie.Value
		user := Controller.GetUserByUsername(username)
		role = user.Role
	}

	salesHistory := Controller.GetSalesHistory()

	data := struct {
		Username string
		Role     string
		Sales    []Node.NodePenjualan
	}{
		Username: username,
		Role:     role,
		Sales:    salesHistory,
	}

	if err := manageHistoryPenjualanTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func RecordSaleHandler(w http.ResponseWriter, r *http.Request) {
	var saleData SaleData

	err := json.NewDecoder(r.Body).Decode(&saleData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Simpan data penjualan ke dalam linked list atau database
	newSale := &Node.PenjualanLL{
		Penjualan: Node.NodePenjualan{
			IdPenjualan: Controller.GenerateIdPenjualan(), // Fungsi untuk menghasilkan ID penjualan
			Total:       saleData.TotalHarga,
			Detail:      convertSaleDetailsToNodeDetails(saleData.Details),
			JmlTunai:    saleData.JumlahTunai,
			Kembalian:   saleData.Kembalian,
			TotalDiskon: calculateTotalDiskon(saleData.Details),
			IsMember:    saleData.IsMember,
			NamaMember: saleData.NamaMember,
			PointReward: saleData.PointReward,
			PointUsed: saleData.PointUsed,
			CreateAt:    time.Now().Format("2006-01-02 15:04:05"),
		},
		Next: nil,
	}

	Controller.AddPenjualan(newSale) 
	ReadAllPenjualan()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Sale recorded successfully"})
}

func convertSaleDetailsToNodeDetails(details []SaleDetail) []Node.NodeDetailPenjualan {
	nodeDetails := make([]Node.NodeDetailPenjualan, len(details))
	for i, detail := range details {
		nodeDetails[i] = Node.NodeDetailPenjualan{
			NodeItem: Node.NodeItem{
				Nama: detail.NamaItem,
				Harga:    detail.Harga,
				HargaDiskon : detail.HargaDiskon,
				Diskon: detail.Diskon,
			},
			JmlPesanan: detail.Jumlah,
			JmlTerjual:  detail.Jumlah,
		}
	}
	return nodeDetails
}

func calculateTotalDiskon(details []SaleDetail) int {
	totalDiskon := 0
	for _, detail := range details {
		totalDiskon += detail.HargaDiskon * detail.Jumlah
	}
	return totalDiskon
}

func HandleGetSalesHistory(w http.ResponseWriter, r *http.Request) {
	history := Controller.GetSalesHistory()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func HandleDeletePenjualan(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	err = Controller.DeletePenjualan(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetDetailPenjualanHandler(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    penjualan, detailPenjualan := Controller.GetDetailPenjualan(id)
    if penjualan == nil {
        http.Error(w, "Penjualan not found", http.StatusNotFound)
        return
    }

    response := struct {
        IdPenjualan  int                     `json:"idPenjualan"`
        CreateAt     string                  `json:"createAt"`
        Total        int                     `json:"total"`
        JmlTunai     int                     `json:"jmlTunai"`
        Kembalian    int                     `json:"kembalian"`
        TotalDiskon  int                     `json:"totalDiskon"`
        IsMember     bool                    `json:"isMember"`
        Detail       []Node.NodeDetailPenjualan `json:"detail"`
        PointReward  int                    `json:"pointReward"`
        PointUsed    int                    `json:"pointUsed"`
        NamaMember   string                 `json:"namaMember"`
    }{
        IdPenjualan:  penjualan.IdPenjualan,
        CreateAt:     penjualan.CreateAt,
        Total:        penjualan.Total,
        JmlTunai:     penjualan.JmlTunai,
        Kembalian:    penjualan.Kembalian,
        TotalDiskon:  penjualan.TotalDiskon,
        IsMember:     penjualan.IsMember,
        Detail:       detailPenjualan,
        PointReward:  penjualan.PointReward,
        PointUsed:  penjualan.PointUsed,
		NamaMember: penjualan.NamaMember,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
