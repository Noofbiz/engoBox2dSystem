package engoBox2dSystem

import "engo.io/engo/math"

// PxToMeters converts from RenderSystem's pixels to box2d's meters
func PxToMeters(px float32) float32 {
	return px / pixelsPerMeter
}

// MetersToPx converts from box2d's meters to RenderSystem's pixels
func MetersToPx(m float32) float32 {
	return m * pixelsPerMeter
}

// RadToDeg converts from radians to degrees
func RadToDeg(r float32) float32 {
	return r * 180 / math.Pi
}

// DegToRad converts from degrees to radians
func DegToRad(d float32) float32 {
	return d * math.Pi / 180
}
