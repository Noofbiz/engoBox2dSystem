package engoBox2dSystem

import (
	"github.com/ByteArena/box2d"
)

var (
	pixelsPerMeter                         float32
	velocityIterations, positionIterations int
	// World is the box2d World used to generate bodies, test bodies for collisions, and simulate physics
	World box2d.B2World
)

// InitBox2dSystem initalizes the box2d system for usage.
// Creates the box2d World and sets pixels per meter and iterations for the physics system.
// Gravity is the B2Vec2 used to create the box2d World
// PixelsPerMeter is the scale factor between the box2d system and the render system.
// VelocityIterations is how many times the physics engine calculates velocity per Update
// PositionIterations is how many times the physics engine calculates position per update
func InitBox2dSystem(Gravity box2d.B2Vec2, PixelsPerMeter float32, VelocityIterations, PositionIterations int) {
	pixelsPerMeter = PixelsPerMeter
	velocityIterations = VelocityIterations
	positionIterations = PositionIterations
	World = box2d.MakeB2World(Gravity)
}
