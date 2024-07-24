package Controller

import (
	"THR/Database"
	"THR/Node"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

//SERVICE
func ValidasiTambahStokItem(id int,addJml int)(error){
	item,itemErr := GetItemById(id)
	if itemErr != nil {
		fmt.Println("item not found")
		return itemErr
	}

	item.JmlStock += addJml
	query := "UPDATE items SET jmlStock = ? WHERE id = ?"
	_, err := Database.DBConnect.Exec(query, item.JmlStock,id)
	return err
}

func ValidasiKurangiStokItem(id int,addJml int)error{
	item,itemErr := GetItemById(id)
	if itemErr != nil {
		fmt.Println("item not found")

		return itemErr
	}

	item.JmlStock -= addJml
	query := "UPDATE items SET jmlStock = ? WHERE id = ?"
	_, err := Database.DBConnect.Exec(query, item.JmlStock,id)
	return err
}

func ValidasiInsertItem(nama string,jmlStock int, Harga int,diskon int) (bool,error) {

	//jika input tidak valid
    if(nama == "" || Harga == 0 || jmlStock == 0){
		fmt.Println("inputan tidak valid")
       return false,nil
	}	
	
    // Menggunakan waktu sekarang dan mengonversinya ke zona waktu Surabaya
    loc, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        return false, err
    }
    setCreateAt := time.Now().In(loc)

	hargaDiskon := Harga*diskon/100

    query := "INSERT INTO items (nama, jmlStock, Harga, HargaDiskon, diskon,createAt) VALUES (?, ?, ?, ?, ?, ?)"
    _, err = Database.DBConnect.Exec(query, nama, jmlStock, Harga,hargaDiskon, diskon,setCreateAt)
    if err == nil {
		return true,nil
    }
	fmt.Println(err)
	return false,err
}


func GetAllItems()([]Node.NodeItem,error){
	query := "SELECT id,nama,jmlStock,harga,hargaDiskon,diskon,createAt,jmlTerjual FROM items"
	rows,err := Database.DBConnect.Query(query)

	if err != nil {
        return nil, err
    }
    defer rows.Close()

    var AllItems []Node.NodeItem
    for rows.Next() {
        //var penampung
        var id int
        var nama string
        var harga int
        var hargaDiskon int
        var diskon int
        var jmlTerjual int
        var jmlStock int
        var createAt string // Baca sebagai string terlebih dahulu
    
        //pindahkan isi query ke var
        if err := rows.Scan(&id, &nama,&jmlStock ,&harga, &hargaDiskon,&diskon, &createAt,&jmlTerjual); err != nil {
            return nil, err
        }
    
        // Konversi string ke time.Time
        parsedTime, err := time.Parse("2006-01-02 15:04:05", createAt)
        if err != nil {
            return nil, err
        }

        // Tentukan zona waktu Surabaya
        loc, err := time.LoadLocation("Asia/Jakarta")
        if err != nil {
            return nil, err
        }

        // Konversi waktu ke zona waktu Surabaya
        parsedTime =  parsedTime.In(loc)
    
        item := Node.NodeItem{
            Id:       id,
            Nama: nama,
			JmlStock: jmlStock,
            Harga:   harga,
            HargaDiskon:    hargaDiskon,
			Diskon: diskon,
            CreateAt: parsedTime,
			JmlTerjual: jmlTerjual,
        }
    
        //masukkan data member ke slice AllItems
        AllItems = append(AllItems, item)
    }
    if err := rows.Err(); err != nil {
		return nil, err
	}

	return AllItems, nil
}

func GetItemById(id int) (Node.NodeItem, error) {
	query := "SELECT id, nama, jmlStock, harga, diskon, hargaDiskon, jmlTerjual, createAt FROM items WHERE id = ?"
	row := Database.DBConnect.QueryRow(query, id)

    var createAt string 
	var item Node.NodeItem
	if err := row.Scan(&item.Id, &item.Nama, &item.JmlStock, &item.Harga,&item.Diskon, &item.HargaDiskon,&item.JmlTerjual,&createAt); err != nil {
		if err == sql.ErrNoRows {
			return item, errors.New("item not found")
		}
     // Konversi string ke time.Time
    parsedTime, _ := time.Parse("2006-01-02 15:04:05", createAt)
 
    // Tentukan zona waktu Surabaya
    loc, _ := time.LoadLocation("Asia/Jakarta")
    
    // Konversi waktu ke zona waktu Surabaya
    parsedTime =  parsedTime.In(loc)
    item.CreateAt = parsedTime
		return item, err
	}
	return item, nil
}


func DeleteItem(id int) error {
	query := "DELETE FROM items WHERE id = ?"
	_, err := Database.DBConnect.Exec(query, id)
	return err
}


func UpdateItem(nama string,jmlStock int,Harga int,diskon int,idItem int)error{
	//update jml terjual
	// item,_ := GetItemById(idItem)
	// item.JmlTerjual += 1
	hargaDiskon := Harga*diskon/100

	query := "UPDATE items SET nama = ?, jmlStock = ?, harga = ?,diskon = ?,hargaDiskon = ? WHERE id = ?"
	_, err := Database.DBConnect.Exec(query, nama,jmlStock,Harga,diskon,hargaDiskon,idItem)
	return err
	
}