package Node

type NodeDetailPenjualan struct {
	NodeItem
	JmlPesanan int
	JmlTerjual int
}

type NodePenjualan struct {
	IdPenjualan int
	Total       int
	Detail      []NodeDetailPenjualan
	JmlTunai    int
	Kembalian   int
	TotalDiskon int
	IsMember    bool
	PointReward int
	PointUsed   int
	NamaMember  string
	CreateAt    string
}

type PenjualanLL struct {
	Penjualan NodePenjualan
	Next      *PenjualanLL
}
