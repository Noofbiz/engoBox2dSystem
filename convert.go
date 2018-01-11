package engoBox2dSystem

import (
	"engo.io/engo/math"
        "engo.io/engo"
	
        "github.com/ByteArena/box2d"
)

//Convert handles conversion between the space component and the Box2d World
type Convert struct {
        // PixelsPerMeter is how many pixels in the space component are in one meter in the Box2d World
        PixelsPerMeter float32
}

// TheConverter is the internal converter used by the system. Use this rather than create your own,
// but you can change the pixels per meter here and it'll change it for all the systems too. 
var TheConverter = &Convert{20}

// UnitType is the type of unit to convert, either positional or angular
type UnitType uint8

const (
        // Positional represents the units for the X or Y coordinates
        Positional UnitType = iota
        // Angular represents a rotation
        Angular
)

// ToRender converts values from the Box2d World into values for the space component
func (c *Convert) ToRender(val float64, unit UnitType) float32 {
        var ret float32
	switch unit {
                case Positional:
                        ret = c.MetersToPx(val)
                case Angular:
                        ret = c.RadToDeg(val)
	}
        return ret
}

// ToBox2d converts values from the space component into values for the Box2d World 
func (c *Convert) ToBox2d(val float32, unit UnitType) float64 {
	var ret float64
	switch unit {
                case Positional:
                        ret = c.PxToMeters(val)
                case Angular:
                        ret = c.DegToRad(val)
	}
        return ret
}

// ToEngoPoint converts a box2d.B2Vec2 into an engo.Point
// 
// note that the units are converted, not just copying values
func (c *Convert) ToEngoPoint(vec box2d.B2Vec2) engo.Point {
        return engo.Point{
                X: c.ToRender(vec.X, Positional),
                Y: c.ToRender(vec.Y, Positional),
	}
}

// ToBox2d2Vec converts an engo.Point into a box2d.B2Vec2
//
// note that the units are converted, not just copying values
func (c *Convert) ToBox2d2Vec(pt engo.Point) box2d.B2Vec2 {
	return box2d.B2Vec2{
                X: c.ToRender(pt.X, Positional),
                Y: c.ToRender(pt.Y, Positional),
	}
}

// PxToMeters converts from the space component's px to Box2d's meters
func (c *Convert) PxToMeters(px float32) float64 {
	return float64(px / c.PixelsPerMeter)
}

// MetersToPx converts from Box2d's meters to the space component's pixels
func (c *Convert) MetersToPx(m float64) float32 {
	return float32(m) * c.PixelsPerMeter
}

// RadToDeg converts from radians to degrees
func (c *Convert) RadToDeg(r float64) float32 {
	return float32(r * 180) / math.Pi
}

// DegToRad converts from degrees to radians
func (c *Convert) DegToRad(d float32) float64 {
	return float64(d * math.Pi / 180)
}
