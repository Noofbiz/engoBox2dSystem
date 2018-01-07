package engoBox2dSystem

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"

	"github.com/ByteArena/box2d"
)

// Box2dComponent holds the box2d Body for use by Systems
type Box2dComponent struct {
	// Body is the box2d body
	Body *box2d.B2Body
}

var listOfBodiesToRemove []*box2d.B2Body

// DestroyBody destroys the box2d body from the World
// this does it safely at the end of an Update, so no bodies are removed during
// a simulation step, which can cause a crash
func (b *Box2dComponent) DestroyBody() {
	listOfBodiesToRemove = append(listOfBodiesToRemove, b.Body)
}

type box2dPhysicsEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*Box2dComponent
}

// Box2dPhysicsSystem provides a system that allows entites to follow the box2d
// physics engine calculations.
type Box2dPhysicsSystem struct {
	entities []box2dPhysicsEntity
}

// Add adds the entity to the physics system
// An entity needs a engo.io/ecs.BasicEntity, engo.io/engo/common.SpaceComponent, and a Box2dComponent in order to be added to the system
func (b *Box2dPhysicsSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent, box *Box2dComponent) {
	b.entities = append(b.entities, box2dPhysicsEntity{basic, space, box})
}

// Remove removes the entity from the physics system.
func (b *Box2dPhysicsSystem) Remove(basic ecs.BasicEntity) {
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
func (b *Box2dPhysicsSystem) Update(dt float32) {
	//Set World components to the Render/Space Components
	for _, e := range b.entities {
		position := box2d.B2Vec2{
			X: float64(PxToMeters(e.SpaceComponent.Center().X)),
			Y: float64(PxToMeters(e.SpaceComponent.Center().Y)),
		}
		e.Body.SetTransform(position, float64(DegToRad(e.SpaceComponent.Rotation)))
	}

	World.Step(float64(dt), velocityIterations, positionIterations)

	//Update Render/Space components to World components after simulation
	for _, e := range b.entities {
		position := e.Body.GetPosition()
		point := engo.Point{
			X: MetersToPx(float32(position.X)),
			Y: MetersToPx(float32(position.Y)),
		}
		e.SpaceComponent.Rotation = RadToDeg(float32(e.Body.GetAngle()))
		e.SpaceComponent.SetCenter(point)
	}

	//Remove all bodies on list for removal
	for _, bod := range listOfBodiesToRemove {
		World.DestroyBody(bod)
	}
	listOfBodiesToRemove = make([]*box2d.B2Body, 0)
}
