package exporter

import (
	"math"
)

// HACK: unfortunately max power reporting in the game is not reliable
// we do it ourselves.
// Variable rate buildings:
// portal, converter, quantum encoder, particle accelerator do not report correctly.
//
// pull static usage numbers from the wiki here
//
// additionally, for the particle accelerator we also need to query the recipe
// as different recipes have different power requirements
//
// Formulas for shards:
// https://satisfactory.wiki.gg/wiki/Tutorial:Advanced_clock_speed#Clock_speed_for_production_buildings
// formula for sloops:
// https://satisfactory.wiki.gg/wiki/Production_amplifier#Power_usage
//
// Sloop Power Multiplier: (1 + filled slots / total slots)^2
// Shard power: (clock speed / 100)^1.321928
//
// See also https://questions.satisfactorygame.com/post/67050eedddb9d97e071f63a2

var MaxPortalPower = 1000.0
var MaxConverterPower = 400.0
var MaxQuantumEncoderPower = 2000.0
var ClockspeedExponent = 1.321928

func powerMultiplier(clockspeed float64, sloops float64, slots float64) float64 {
	sloopPowerMultiplier := math.Pow(1.0+(sloops/slots), 2)
	shardPowerMultiplier := math.Pow((clockspeed / 100), ClockspeedExponent)
	return sloopPowerMultiplier * shardPowerMultiplier
}

func MaxParticleAcceleratorPower(recipe string) float64 {
	recipes := map[string]float64{
		"Dark Matter Crystal":         1500.0,
		"Diamonds":                    750.0,
		"Ficsonium":                   1500.0,
		"Nuclear Pasta":               1500.0,
		"Plutonium Pellet":            750.0,
		"Alternate: Cloudy Diamonds":             750.0,
		"Alternate: Dark Matter Crystallization": 1500.0,
		"Alternate: Dark Matter Trap":            1500.0,
		"Alternate: Instant Plutonium Cell":      750.0,
		"Alternate: Oil-Based Diamonds":          750.0,
		"Alternate: Petroleum Diamonds":          750.0,
		"Alternate: Turbo Diamonds":              750.0,
	}
	maxPower, ok := recipes[recipe]
	if ok {
		return maxPower
	} else {
		return 1500.0
	}
}
