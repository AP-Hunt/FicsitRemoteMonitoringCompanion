package exporter

import (
	"log"
	"time"
)

type VehicleCollector struct {
	endpoint        string
	TrackedVehicles map[string]*VehicleDetails
}

type VehicleDetails struct {
	Id           string   `json:"ID"`
	VehicleType  string   `json:"Name"`
	Location     Location `json:"location"`
	ForwardSpeed float64  `json:"ForwardSpeed"`
	AutoPilot    bool     `json:"Autopilot"`
	Fuel         []Fuel   `json:Fuel`
	PathName     string   `json:"PathName"`
	DepartTime   time.Time
	Departed     bool
}

type Fuel struct {
	Name   string  `json:Name`
	Amount float64 `json:Amount`
}

func (v *VehicleDetails) recordElapsedTime(frmAddress string, saveName string) {
	now := Clock.Now()
	tripSeconds := now.Sub(v.DepartTime).Seconds()
	VehicleRoundTrip.WithLabelValues(v.Id, v.VehicleType, v.PathName, frmAddress, saveName).Set(tripSeconds)
	v.Departed = false
}

func (v *VehicleDetails) isCompletingTrip(updatedLocation Location) bool {
	// vehicle near first tracked location facing roughly the same way
	return v.Departed && v.Location.isNearby(updatedLocation) && v.Location.isSameDirection(updatedLocation)
}

func (v *VehicleDetails) isStartingTrip(updatedLocation Location) bool {
	// vehicle departed from first tracked location
	return !v.Departed && !v.Location.isNearby(updatedLocation)
}

func (v *VehicleDetails) startTracking(trackedVehicles map[string]*VehicleDetails) {
	// Only start tracking the vehicle at low speeds so it's
	// likely at a station or somewhere easier to track.
	if v.ForwardSpeed < 10 {
		trackedVehicle := VehicleDetails{
			Id:          v.Id,
			Location:    v.Location,
			VehicleType: v.VehicleType,
			PathName:    v.PathName,
			Departed:    false,
		}
		trackedVehicles[v.Id] = &trackedVehicle
	}
}

func (d *VehicleDetails) handleTimingUpdates(trackedVehicles map[string]*VehicleDetails, frmAddress string, saveName string) {
	if d.AutoPilot {
		vehicle, exists := trackedVehicles[d.Id]
		if exists && vehicle.isCompletingTrip(d.Location) {
			vehicle.recordElapsedTime(frmAddress, saveName)
		} else if exists && vehicle.isStartingTrip(d.Location) {
			vehicle.Departed = true
			vehicle.DepartTime = Clock.Now()
		} else if !exists {
			d.startTracking(trackedVehicles)
		}
	} else {
		//remove manual vehicles, nothing to mark
		_, exists := trackedVehicles[d.Id]
		if exists {
			delete(trackedVehicles, d.Id)
		}
	}
}

func NewVehicleCollector(endpoint string) *VehicleCollector {
	return &VehicleCollector{
		endpoint:        endpoint,
		TrackedVehicles: make(map[string]*VehicleDetails),
	}
}

func (c *VehicleCollector) Collect(frmAddress string, saveName string) {
	details := []VehicleDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading vehicle statistics from FRM: %s\n", err)
		return
	}

	for _, d := range details {
		if len(d.Fuel) > 0 {
			VehicleFuel.WithLabelValues(d.Id, d.VehicleType, d.Fuel[0].Name, frmAddress, saveName).Set(d.Fuel[0].Amount)
		}

		d.handleTimingUpdates(c.TrackedVehicles, frmAddress, saveName)
	}
}
