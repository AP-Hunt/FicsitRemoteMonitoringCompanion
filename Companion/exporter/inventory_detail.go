package exporter

type InventoryItem struct {
	Name      string  `json:"Name"`
	Amount    float64 `json:"Amount"`
	MaxAmount float64 `json:"MaxAmount"`
}

type ContainerDetail struct {
	Name      string          `json:"Name"`
	Location  Location        `json:"location"`
	Inventory []InventoryItem `json:"Inventory"`
}
