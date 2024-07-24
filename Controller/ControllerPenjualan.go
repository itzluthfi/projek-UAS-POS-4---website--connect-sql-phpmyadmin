package Controller

import (
	"THR/Model"
	"THR/Node"
	"errors"
	"log"
	"sync"
)

func ValidasiPilihItem(id int, jmlPesanan int) (string, []Node.NodeDetailPenjualan) {
	var Items []Node.NodeDetailPenjualan

	cekId := Model.SearchItem(id)

	if cekId != nil && jmlPesanan > 0 {
		addressItem := validasiStokKosong(id)
		if addressItem == nil {
			return "stok Item masih kosong!", nil
		} else {
			if addressItem.Next.Item.JmlStock < jmlPesanan {
				return "jumlah Penjualan Yang anda masukkan lebih besar dari stok Item!", nil
			}
		}
		item := addressItem.Next.Item
		newItem := Node.NodeDetailPenjualan{
			NodeItem: Node.NodeItem{
				Id:     item.Id,
				Nama:   item.Nama,
				Harga:  item.Harga,
				Diskon: item.Diskon,
				HargaDiskon: item.HargaDiskon, // Menambahkan HargaDiskon ke struct NodeDetailPenjualan
			},
			JmlPesanan:  jmlPesanan,
		}
		Items = append(Items, newItem)
		addressItem.Next.Item.JmlStock -= jmlPesanan
		return "Data Item Baru Berhasil Ditambahkan!", Items
	}

	return "id atau jumlah yang anda inputkan tidak valid!", nil
}

func ValidasiIsMember(idMember int) *Node.MemberLL {
	addressMember := Model.SearchMember(idMember)
	if addressMember == nil {
		return nil
	}
	return addressMember.Next
}

func ValidasiInsertPenjualan(addressMember *Node.MemberLL, detailItems []Node.NodeDetailPenjualan, jmlTunai, kembalian, totalDiskon, jmlPoint int) (string, bool) {
	if jmlTunai == -1 {
		Model.KembalikanJmlStok(detailItems)
		return "Pembayaran kurang, silakan masukkan jumlah pembayaran yang mencukupi.", false
	}

	if addressMember != nil {
		Model.InsertPenjualanMember(addressMember, detailItems, jmlTunai, kembalian, totalDiskon, jmlPoint)
		return "Member Berhasil Dikonfirmasi", true
	}
	Model.InsertPenjualanNonMember(detailItems, jmlTunai, kembalian, totalDiskon)
	return "Data Berhasil Ditambahkan!", true
}

// func ValidasiSearchPenjualan(id int) (bool, *Node.PenjualanLL) {
// 	cur := &Database.HeadPenjualan

// 	if cur.Next != nil {
// 		hasil := Model.SearchPenjualan(id)
// 		if hasil != nil {
// 			return true, hasil.Next
// 		} else {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

func validasiStokKosong(id int) *Node.ItemLL {
	result := Model.SearchItem(id)

	if result.Next.Item.JmlStock > 0 {
		return result
	}

	return nil
}




//WEB
var (
	penjualanHead *Node.PenjualanLL
	penjualanLock sync.Mutex
	lastIDPenjualan int
)

// Inisialisasi linked list penjualan
func init() {
	penjualanHead = nil
	lastIDPenjualan = 0
}

// Generate ID Penjualan baru
func GenerateIdPenjualan() int {
	penjualanLock.Lock()
	defer penjualanLock.Unlock()
	lastIDPenjualan++
	return lastIDPenjualan
}

// // Cari penjualan berdasarkan ID
// func ValidasiSearchPenjualan(id int) (bool, *Node.NodePenjualan) {
// 	penjualanLock.Lock()
// 	defer penjualanLock.Unlock()

// 	current := penjualanHead
// 	for current != nil {
// 		if current.Penjualan.IdPenjualan == id {
// 			return true, &current.Penjualan
// 		}
// 		current = current.Next
// 	}

// 	return false, nil
// }

// Tambahkan penjualan baru ke linked list
func AddPenjualan(penjualan *Node.PenjualanLL) {
	penjualanLock.Lock()
	defer penjualanLock.Unlock()

	if penjualanHead == nil {
		penjualanHead = penjualan
	} else {
		current := penjualanHead
		for current.Next != nil {
			current = current.Next
		}
		current.Next = penjualan
	}
}

// Dapatkan riwayat penjualan dari linked list
func GetSalesHistory() []Node.NodePenjualan {
	penjualanLock.Lock()
	defer penjualanLock.Unlock()

	var history []Node.NodePenjualan
	current := penjualanHead

	for current != nil {
		history = append(history, current.Penjualan)
		current = current.Next
	}

	return history
}

// Hapus penjualan dari linked list
func DeletePenjualan(id int) error {
	penjualanLock.Lock()
	defer penjualanLock.Unlock()

	if penjualanHead == nil {
		return errors.New("Tidak ada penjualan yang ditemukan")
	}

	if penjualanHead.Penjualan.IdPenjualan == id {
		penjualanHead = penjualanHead.Next
		return nil
	}

	current := penjualanHead
	for current.Next != nil {
		if current.Next.Penjualan.IdPenjualan == id {
			current.Next = current.Next.Next
			return nil
		}
		current = current.Next
	}

	return errors.New("Penjualan tidak ditemukan")
}

// Dapatkan detail penjualan berdasarkan ID
func GetDetailPenjualan(id int) (*Node.NodePenjualan, []Node.NodeDetailPenjualan) {
    penjualanLock.Lock()
    defer penjualanLock.Unlock()

    log.Println("Searching for ID:", id)  // Debug log
    current := penjualanHead
    for current != nil {
        log.Println("Checking penjualan with ID:", current.Penjualan.IdPenjualan)  // Debug log
        if current.Penjualan.IdPenjualan == id {
            log.Println("Found penjualan for ID:", id)  // Debug log
            return &current.Penjualan, current.Penjualan.Detail
        }
        current = current.Next
    }
    log.Println("Penjualan not found for ID:", id)  // Debug log
    return nil, nil
}
