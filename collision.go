package engoBox2dSystem

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/ByteArena/box2d"
)

// CollisionStartMessage is sent out for the box2d collision callback
// CollisionStart
type CollisionStartMessage struct {
	Contact box2d.B2ContactInterface
}

// Type implements the engo.Message interface
func (CollisionStartMessage) Type() string { return "CollisionStartMessage" }

// CollisionEndMessage is sent out for the box2d collision callback
// CollisionEnd
type CollisionEndMessage struct {
	Contact box2d.B2ContactInterface
}

// Type implements the engo.Message interface
func (CollisionEndMessage) Type() string { return "CollisionEndMessage" }

// PreSolveMessage is sent out before a step of the physics engine
type PreSolveMessage struct {
	Contact     box2d.B2ContactInterface
	OldManifold box2d.B2Manifold
}

// Type implements the engo.Message interface
func (PreSolveMessage) Type() string { return "PreSolveMessage" }

// PostSolveMessage is sent out after a step of the physics engine
type PostSolveMessage struct {
	Contact box2d.B2ContactInterface
	Impulse *box2d.B2ContactImpulse
}

// Type implements the engo.Message interface
func (PostSolveMessage) Type() string { return "PostSolveMessage" }

type collisionEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*Box2dComponent
}

// CollisionSystem is a system that handles the callbacks for box2d's
// collision system. This system does not require the physics system, but a
// they do need box2d bodies.
type CollisionSystem struct {
	entities []collisionEntity
}

// New sets the system to the contact listener for box2d, which allows the collision
// messages to be sent out.
func (c *CollisionSystem) New(w *ecs.World) {
	World.SetContactListener(c)
}

// Add adds the entity to the collision system.
// It also adds the body's user data to the BasicEntity's ID, which makes it
// easy to figure out which entities are which when comparing in the messages / callbacks
func (c *CollisionSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent, box *Box2dComponent) {
	box.Body.SetUserData(basic.ID())
	c.entities = append(c.entities, collisionEntity{basic, space, box})
}

// AddByInterface adds the entity to the collision system if it implements the Collisionable interface
func (c *CollisionSystem) AddByInterface(o Collisionable) {
	c.Add(o.GetBasicEntity(), o.GetSpaceComponent(), o.GetBox2dComponent())
}

// Remove removes the entity from the system
func (c *CollisionSystem) Remove(basic ecs.BasicEntity) {
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

// Update doesn't do anything, since the physics engine handles passing out the
// callbacks.
func (c *CollisionSystem) Update(dt float32) {}

// BeginContact implements the B2ContactListener interface.
// when a BeginContact callback is made by box2d, it sends a message containing
// the information from the callback.
func (c *CollisionSystem) BeginContact(contact box2d.B2ContactInterface) {
	engo.Mailbox.Dispatch(CollisionStartMessage{
		Contact: contact,
	})
}

// EndContact implements the B2ContactListener interface.
// when a EndContact callback is made by box2d, it sends a message containing
// the information from the callback.
func (c *CollisionSystem) EndContact(contact box2d.B2ContactInterface) {
	engo.Mailbox.Dispatch(CollisionEndMessage{
		Contact: contact,
	})
}

// PreSolve implements the B2ContactListener interface.
// this is called after a contact is updated but before it goes to the solver.
// When it is called, a message is sent containing the information from the callback
func (c *CollisionSystem) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
	engo.Mailbox.Dispatch(PreSolveMessage{
		Contact:     contact,
		OldManifold: oldManifold,
	})
}

// PostSolve implements the B2ContactListener interface.
// this is called after the solver is finished.
// When it is called, a message is sent containing the information from the callback
func (c *CollisionSystem) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
	engo.Mailbox.Dispatch(PostSolveMessage{
		Contact: contact,
		Impulse: impulse,
	})
}
