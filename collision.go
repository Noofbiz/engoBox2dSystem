package engoBox2dSystem

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/ByteArena/box2d"
)

// Box2dCollisionStartMessage is sent out for the box2d collision callback
// CollisionStart
type Box2dCollisionStartMessage struct {
	Contact box2d.B2ContactInterface
}

// Type implements the engo.Message interface
func (Box2dCollisionStartMessage) Type() string { return "Box2dCollisionStartMessage" }

// Box2dCollisionEndMessage is sent out for the box2d collision callback
// CollisionEnd
type Box2dCollisionEndMessage struct {
	Contact box2d.B2ContactInterface
}

// Type implements the engo.Message interface
func (Box2dCollisionEndMessage) Type() string { return "Box2dCollisionEndMessage" }

// Box2dPreSolveMessage is sent out before a step of the physics engine
type Box2dPreSolveMessage struct {
	Contact     box2d.B2ContactInterface
	OldManifold box2d.B2Manifold
}

// Type implements the engo.Message interface
func (Box2dPreSolveMessage) Type() string { return "Box2dPreSolveMessage" }

// Box2dPostSolveMessage is sent out after a step of the physics engine
type Box2dPostSolveMessage struct {
	Contact box2d.B2ContactInterface
	Impulse *box2d.B2ContactImpulse
}

// Type implements the engo.Message interface
func (Box2dPostSolveMessage) Type() string { return "Box2dPostSolveMessage" }

type box2dCollisionEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*Box2dComponent
}

// Box2dCollisionSystem is a system that handles the callbacks for box2d's
// collision system. This system does not require the physics system, but a
// box2d World must be initalized, and the entities need the box2d bodies to work.
type Box2dCollisionSystem struct {
	entities []box2dCollisionEntity
}

// New sets the system to the contact listener for box2d, which allows the collision
// messages to be sent out.
func (c *Box2dCollisionSystem) New(w *ecs.World) {
	World.SetContactListener(c)
}

// Add adds the entity to the collision system.
// It also adds the body's user data to the BasicEntity's ID, which makes it
// easy to figure out which entities are which when comparing in the messages / callbacks
func (c *Box2dCollisionSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent, box *Box2dComponent) {
	box.Body.SetUserData(basic.ID())
	c.entities = append(c.entities, box2dCollisionEntity{basic, space, box})
}

// Remove removes the entity from the system
func (c *Box2dCollisionSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

// Update just syncs the space components with the box2d bodies, so the collisions
// are based on the SpaceComponent position and rotation
func (c *Box2dCollisionSystem) Update(dt float32) {
	//Set World components to the Render/Space Components
	for _, e := range c.entities {
		position := box2d.B2Vec2{
			X: float64(PxToMeters(e.SpaceComponent.Center().X)),
			Y: float64(PxToMeters(e.SpaceComponent.Center().Y)),
		}
		e.Body.SetTransform(position, float64(DegToRad(e.SpaceComponent.Rotation)))
	}

	//Remove all bodies on list for removal
	for _, bod := range listOfBodiesToRemove {
		World.DestroyBody(bod)
	}
	listOfBodiesToRemove = make([]*box2d.B2Body, 0)
}

// BeginContact implements the B2ContactListener interface.
// when a BeginContact callback is made by box2d, it sends a message containing
// the information from the callback.
func (c *Box2dCollisionSystem) BeginContact(contact box2d.B2ContactInterface) {
	engo.Mailbox.Dispatch(Box2dCollisionStartMessage{
		Contact: contact,
	})
}

// EndContact implements the B2ContactListener interface.
// when a EndContact callback is made by box2d, it sends a message containing
// the information from the callback.
func (c *Box2dCollisionSystem) EndContact(contact box2d.B2ContactInterface) {
	engo.Mailbox.Dispatch(Box2dCollisionEndMessage{
		Contact: contact,
	})
}

// PreSolve implements the B2ContactListener interface.
// this is called after a contact is updated but before it goes to the solver.
// When it is called, a message is sent containing the information from the callback
func (c *Box2dCollisionSystem) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
	engo.Mailbox.Dispatch(Box2dPreSolveMessage{
		Contact:     contact,
		OldManifold: oldManifold,
	})
}

// PostSolve implements the B2ContactListener interface.
// this is called after the solver is finished.
// When it is called, a message is sent containing the information from the callback
func (c *Box2dCollisionSystem) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
	engo.Mailbox.Dispatch(Box2dPostSolveMessage{
		Contact: contact,
		Impulse: impulse,
	})
}
