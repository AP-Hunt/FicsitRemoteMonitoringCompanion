package exporter

import (
	"log"
	"time"
)

type VehicleCollector struct {
	FRMAddress      string
	TrackedVehicles map[string]*VehicleDetails
}

type VehicleDetails struct {
	Id           string   `json:"ID"`
	VehicleType  string   `json:"Name"`
	Location     Location `json:"location"`
	ForwardSpeed float64  `json:"ForwardSpeed"`
	AutoPilot    bool     `json:"AutoPilot"`
	Fuel         Fuel     `json:Fuel`
	PathName     string   `json:"PathName"`
	DepartTime   time.Time
	Departed     bool
}

type Fuel struct {
	Name   string  `json:Name`
	Amount float64 `json:Amount`
}

func (v *VehicleDetails) recordElapsedTime() {
	now := Clock.Now()
	tripSeconds := now.Sub(v.DepartTime).Seconds()
	VehicleRoundTrip.WithLabelValues(v.Id, v.VehicleType, v.PathName).Set(tripSeconds)
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

func (d *VehicleDetails) handleTimingUpdates(trackedVehicles map[string]*VehicleDetails) {
	if d.AutoPilot {
		vehicle, exists := trackedVehicles[d.Id]
		if exists && vehicle.isCompletingTrip(d.Location) {
			vehicle.recordElapsedTime()
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

func NewVehicleCollector(frmAddress string) *VehicleCollector {
	return &VehicleCollector{
		FRMAddress:      frmAddress,
		TrackedVehicles: make(map[string]*VehicleDetails),
	}
}

func (c *VehicleCollector) Collect() {
	details := []VehicleDetails{}
	err := retrieveData(c.FRMAddress, &details)
	if err != nil {
		log.Printf("error reading vehicle statistics from FRM: %s\n", err)
		return
	}

	for _, d := range details {
		VehicleFuel.WithLabelValues(d.Id, d.VehicleType, d.Fuel.Name).Set(d.Fuel.Amount)

		d.handleTimingUpdates(c.TrackedVehicles)
	}
}
