package Model

import (
	"THR/Database"
	"THR/Node"
)

func TambahStokItem(cur *Node.ItemLL,addJml int){
	cur.Item.JmlStock += addJml
}

func KurangiStokItem(cur *Node.ItemLL,addJml int){
	cur.Item.JmlStock -= addJml
}

// func InsertItem(nama string, jmlStock int, harga int,diskon int) {	
// 	cur := &Database.HeadItem

// 	newNode := Node.ItemLL{}
// 	newNode.Item.Id = GetItemLastId()
// 	newNode.Item.Nama = nama
// 	newNode.Item.JmlStock = jmlStock
// 	newNode.Item.Harga = harga 
// 	newNode.Item.Diskon = diskon
// 	newNode.Item.HargaDiskon = harga*diskon/100
// 	newNode.Item.CreateAt = time.Now().Format("2006-01-02")
// 	newNode.Next = nil

// 	//jika data pertama
// 	if cur.Next == nil {
// 		cur.Next = &newNode
// 	} else {
// 		//geser ke data terakhir
// 		for cur.Next != nil {
// 			cur = cur.Next
// 		}
// 		cur.Next = &newNode
// 	}
// }

// func GetItemLastId() int { 
// 	cur := &Database.HeadItem

// 	if cur.Next == nil {
// 		return 1
// 	} else {
// 		for cur.Next != nil {
// 			cur = cur.Next
// 		}
// 		return cur.Item.Id + 1
// 	}
// }


func ReadAllItem() []Node.NodeItem{
	cur := &Database.HeadItem

	var ItemTable []Node.NodeItem 

	for cur.Next != nil {
		cur = cur.Next
		ItemTable = append(ItemTable, cur.Item)
	}

	return ItemTable
}


func SearchItem(id int) (*Node.ItemLL){
	cur := &Database.HeadItem
	   
		for cur.Next != nil {
			//cari data
			if cur.Next.Item.Id == id {
				//alamat data prev
				return cur
			}
			cur = cur.Next
	    }

	return nil
}	

func SearchItemWeb(id int) (*Node.NodeItem){
	cur := &Database.HeadItem
	   
		for cur.Next != nil {
			//cari data
			if cur.Next.Item.Id == id {
				//alamat data prev
				return &cur.Next.Item
			}
			cur = cur.Next
	    }

	return nil
}	

func UpdateItem(cur *Node.ItemLL,nama string,jmlStock int,harga int,diskon int){
	cur.Item.Nama = nama
	cur.Item.JmlStock = jmlStock
	cur.Item.Harga = harga 
	cur.Item.Diskon = diskon
	cur.Item.HargaDiskon = harga * diskon/100
}

func DeleteItem(cur *Node.ItemLL){
	cur.Next = cur.Next.Next
}

// func DiskonItem(cur *Node.ItemLL){

// }