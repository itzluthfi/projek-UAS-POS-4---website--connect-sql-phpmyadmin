package Model

import (
	"THR/Database"
	"THR/Node"
	"time"
)


func InsertPenjualanMember(addressMember *Node.MemberLL, detailItem []Node.NodeDetailPenjualan,jmlTunai,kembalian,totalDiskon,jmlPoint int) {	
	newNode := Node.PenjualanLL{}
	newNode.Penjualan.IdPenjualan = GetPenjualanLastId()
	newNode.Penjualan.JmlTunai = jmlTunai
	newNode.Penjualan.TotalDiskon = totalDiskon
	newNode.Penjualan.Kembalian = kembalian
	newNode.Penjualan.Detail = detailItem
	newNode.Penjualan.Total = GetTotalDetail(detailItem) - totalDiskon - jmlPoint
	newNode.Penjualan.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	newNode.Next = nil
	
	cur := &Database.HeadPenjualan
	//jika data pertama
	if cur.Next == nil {
		cur.Next = &newNode
		MemberUpdatePoint(addressMember, GetTotalDetail(detailItem)-totalDiskon-jmlPoint)

	} else {
		//geser ke data terakhir
		for cur.Next != nil {
			cur = cur.Next
			MemberUpdatePoint(addressMember, GetTotalDetail(detailItem)-totalDiskon-jmlPoint)
		}
		cur.Next = &newNode
	}
}

func InsertPenjualanNonMember(details []Node.NodeDetailPenjualan,jmlTunai,kembalian,totalDiskon int) {	
	
	newNode := Node.PenjualanLL{}
	newNode.Penjualan.IdPenjualan = GetPenjualanLastId()
	newNode.Penjualan.JmlTunai = jmlTunai
	newNode.Penjualan.TotalDiskon = totalDiskon
	newNode.Penjualan.Kembalian = kembalian
	newNode.Penjualan.Detail = details
	newNode.Penjualan.Total = GetTotalDetail(details) - totalDiskon
	newNode.Penjualan.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	newNode.Next = nil
	
	cur := &Database.HeadPenjualan
	//jika data pertama
	if cur.Next == nil {
		cur.Next = &newNode
	} else {
		//geser ke data terakhir
		for cur.Next != nil {
			cur = cur.Next
		}
		cur.Next = &newNode
	}
}

func GetTotalDetail(details []Node.NodeDetailPenjualan) int {
	var total int
	for _, detail := range details {
		total += detail.JmlPesanan * detail.NodeItem.Harga
	}
	return total
}


func GetPenjualanLastId() int {
	cur := &Database.HeadPenjualan
	if cur.Next == nil {
		//jika kosong
		return 1
	} else {
		//geser ke data terakhir
		for cur.Next != nil {
			cur = cur.Next
		}
		return cur.Penjualan.IdPenjualan + 1
	}
}

func ReadAllPenjualan() []Node.NodePenjualan{
	//cur head
	cur := &Database.HeadPenjualan
	var PenjualanTable []Node.NodePenjualan   
	//masukkan semua data ke slice
	for cur.Next != nil {
		cur = cur.Next
		PenjualanTable = append(PenjualanTable, cur.Penjualan)
	}
	return PenjualanTable
}

func SearchPenjualan(id int) (*Node.PenjualanLL){
	cur := &Database.HeadPenjualan
	   
		for cur.Next != nil {
			//cari data
			if cur.Next.Penjualan.IdPenjualan == id {
				//alamat data prev
				return cur
			}
			cur = cur.Next
	    }

	return nil
}


func KembalikanJmlStok(detailItems []Node.NodeDetailPenjualan){
	for _, detailItem := range detailItems {
		addressItem := SearchItem(detailItem.NodeItem.Id)
		addressItem.Next.Item.JmlStock += detailItem.JmlPesanan
    }    
}

