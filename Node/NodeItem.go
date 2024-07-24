package Node

import "time"

type NodeItem struct {
	Id          int
	Nama        string
	JmlStock    int
	Harga       int
	HargaDiskon int
	Diskon      int
	CreateAt    time.Time
	JmlTerjual     int
}

type ItemLL struct {
	Item NodeItem
	Next *ItemLL
}
