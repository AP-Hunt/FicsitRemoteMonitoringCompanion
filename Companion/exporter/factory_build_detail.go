package exporter

type BuildingDetail struct {
	Building     string       `json:"building"`
	Location     Location     `json:"location"`
	Recipe       string       `json:"Recipe"`
	Production   []Production `json:"production"`
	Ingredients  []Ingredient `json:"ingredients"`
	ManuSpeed    float64      `json:"ManuSpeed"`
	IsConfigured bool         `json:"IsConfigured"`
	IsProducing  bool         `json:"IsProducing"`
	IsPaused     bool         `json:"IsPaused"`
	CircuitID    int          `json:"CircuitID"`
}

type Location struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
	Rotation int     `json:"rotation"`
}

type Production struct {
	Name        string  `json:"Name"`
	CurrentProd float64 `json:"CurrentProd"`
	MaxProd     float64 `json:"MaxProd"`
	ProdPercent string  `json:"ProdPercent"`
}

type Ingredient struct {
	Name            string  `json:"Name"`
	CurrentConsumed float64 `json:"CurrentConsumed"`
	MaxConsumed     float64 `json:"MaxConsumed"`
	ConsPercent     string  `json:"ConsPercent"`
}
