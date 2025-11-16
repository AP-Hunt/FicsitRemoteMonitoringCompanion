package exporter

type InventoryItem struct {
	Name      string `json:"Name"`
	Amount    int    `json:"Amount"`
	MaxAmount int    `json:"MaxAmount"`
}

type ContainerDetail struct {
	Name      string          `json:"Name"`
	Location  Location        `json:"location"`
	Inventory []InventoryItem `json:"Inventory"`
}
