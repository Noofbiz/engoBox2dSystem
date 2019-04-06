package engoBox2dSystem

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type physicsEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*Box2dComponent
}

// PhysicsSystem provides a system that allows entites to follow the box2d
// physics engine calculations.
type PhysicsSystem struct {
	entities []physicsEntity

	VelocityIterations, PositionIterations int
}

// Add adds the entity to the physics system
// An entity needs a github.com/EngoEngine/ecs.BasicEntity, github.com/EngoEngine/engo/common.SpaceComponent, and a Box2dComponent in order to be added to the system
func (b *PhysicsSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent, box *Box2dComponent) {
	b.entities = append(b.entities, physicsEntity{basic, space, box})
}

// AddByInterface adds an entity to the Physics system
func (b *PhysicsSystem) AddByInterface(o Physicsable) {
	b.Add(o.GetBasicEntity(), o.GetSpaceComponent(), o.GetBox2dComponent())
}

// Remove removes the entity from the physics system.
func (b *PhysicsSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range b.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		b.entities = append(b.entities[:delete], b.entities[delete+1:]...)
	}
}

// Update runs every time the systems update. Updates the box2d world and simulates
// physics based on the timestep, positions, and forces on the bodies.
func (b *PhysicsSystem) Update(dt float32) {
	//Set World components to the Render/Space Components
	for _, e := range b.entities {
		e.Body.SetTransform(Conv.ToBox2d2Vec(e.Center()), Conv.DegToRad(e.Rotation))
	}

	World.Step(float64(dt), b.VelocityIterations, b.PositionIterations)

	//Update Render/Space components to World components after simulation
	for _, e := range b.entities {
		e.SpaceComponent.Rotation = Conv.RadToDeg(e.Body.GetAngle())
		e.SpaceComponent.SetCenter(Conv.ToEngoPoint(e.Body.GetPosition()))
	}

	removeBodies()
}
