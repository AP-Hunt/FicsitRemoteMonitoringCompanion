package exporter

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type VehicleCollector struct {
	endpoint        string
	TrackedVehicles map[string]*VehicleDetails
	metricsDropper *MetricsDropper
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

func (v *VehicleDetails) recordElapsedTime(frmAddress string, sessionName string) {
	now := Clock.Now()
	tripSeconds := now.Sub(v.DepartTime).Seconds()
	VehicleRoundTrip.WithLabelValues(v.Id, v.VehicleType, v.PathName, frmAddress, sessionName).Set(tripSeconds)
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

func (d *VehicleDetails) handleTimingUpdates(trackedVehicles map[string]*VehicleDetails, frmAddress string, sessionName string) {
	if d.AutoPilot {
		vehicle, exists := trackedVehicles[d.Id]
		if exists && vehicle.isCompletingTrip(d.Location) {
			vehicle.recordElapsedTime(frmAddress, sessionName)
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
		metricsDropper: NewMetricsDropper(
			VehicleRoundTrip,
			VehicleFuel,
		),
	}
}

func (c *VehicleCollector) Collect(frmAddress string, sessionName string) {
	details := []VehicleDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		c.metricsDropper.DropStaleMetricLabels()
		log.Printf("error reading vehicle statistics from FRM: %s\n", err)
		return
	}

	for _, d := range details {
		c.metricsDropper.CacheFreshMetricLabel(prometheus.Labels{"url": frmAddress, "session_name": sessionName, "id": d.Id})
		if len(d.Fuel) > 0 {
			VehicleFuel.WithLabelValues(d.Id, d.VehicleType, d.Fuel[0].Name, frmAddress, sessionName).Set(d.Fuel[0].Amount)
		}

		d.handleTimingUpdates(c.TrackedVehicles, frmAddress, sessionName)
	}

	c.metricsDropper.DropStaleMetricLabels()
}
