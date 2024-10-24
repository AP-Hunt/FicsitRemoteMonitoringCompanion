package exporter

var (
	SmelterPower             = 4.0
	ConstructorPower         = 4.0
	AssemblerPower           = 15.0
	ManufacturerPower        = 55.0
	BlenderPower             = 75.0
	RefineryPower            = 30.0
	ParticleAcceleratorPower = 1500.0
)

type BuildingDetail struct {
	Building       string       `json:"Name"`
	Location       Location     `json:"location"`
	Recipe         string       `json:"Recipe"`
	Production     []Production `json:"production"`
	Ingredients    []Ingredient `json:"ingredients"`
	ManuSpeed      float64      `json:"ManuSpeed"`
	IsConfigured   bool         `json:"IsConfigured"`
	IsProducing    bool         `json:"IsProducing"`
	IsPaused       bool         `json:"IsPaused"`
	CircuitGroupId int          `json:"CircuitGroupID"`
	PowerInfo      PowerInfo    `json:"PowerInfo"`
}

type Production struct {
	Name        string  `json:"Name"`
	CurrentProd float64 `json:"CurrentProd"`
	MaxProd     float64 `json:"MaxProd"`
	ProdPercent float64 `json:"ProdPercent"`
}

type Ingredient struct {
	Name            string  `json:"Name"`
	CurrentConsumed float64 `json:"CurrentConsumed"`
	MaxConsumed     float64 `json:"MaxConsumed"`
	ConsPercent     float64 `json:"ConsPercent"`
}
