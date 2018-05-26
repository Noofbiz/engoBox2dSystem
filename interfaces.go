package engoBox2dSystem

import "engo.io/engo/common"

// GetBox2dComponent gets the *Box2dComponent from anything with a one, so they can implement
// the interfaces and AddByInterface can work
func (b *Box2dComponent) GetBox2dComponent() *Box2dComponent {
	return b
}

// GetMouseComponent gets the *MouseComponent
func (m *MouseComponent) GetMouseComponent() *MouseComponent {
	return m
}

// Box2dFace is an interface for the Box2dComponent
type Box2dFace interface {
	GetBox2dComponent() *Box2dComponent
}

// MouseFace is an interface for the MouseComponent
type MouseFace interface {
	GetMouseComponent() *MouseComponent
}

// Collisionable is for the CollisionSystem's AddByInterface
type Collisionable interface {
	common.BasicFace
	common.SpaceFace
	Box2dFace
}

// Mouseable is for he MouseSystem's AddByInterface
type Mouseable interface {
	common.BasicFace
	MouseFace
	common.SpaceFace
	common.RenderFace
	Box2dFace
}

// Physicsable is for the PhysicsSystem's AddByInterface
type Physicsable interface {
	common.BasicFace
	common.SpaceFace
	Box2dFace
}
