package engoBox2dSystem

import "github.com/ByteArena/box2d"

// World is the box2d World used to generate bodies, test bodies for collisions, and simulate physics
var World = box2d.MakeB2World(box2d.B2Vec2{X: 0, Y: 0})
