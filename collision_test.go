package engoBox2dSystem

import (
	"testing"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"

	"github.com/ByteArena/box2d"
)

var (
	recStartMessage, recEndMessage, recPreSolveMessage, recPostSolveMessage bool
)

func TestCollisionSystem(t *testing.T) {
	updateTime := float32(1.0 / 60.0)
	//Setup engo Mailbox
	engo.Mailbox = &engo.MessageManager{}
	engo.Mailbox.Listen("CollisionStartMessage", func(message engo.Message) {
		recStartMessage = true
	})
	engo.Mailbox.Listen("CollisionEndMessage", func(message engo.Message) {
		recEndMessage = true
	})
	engo.Mailbox.Listen("PreSolveMessage", func(message engo.Message) {
		recPreSolveMessage = true
	})
	engo.Mailbox.Listen("PostSolveMessage", func(message engo.Message) {
		recPostSolveMessage = true
	})

	//Create System
	sys := &CollisionSystem{}
	sys.New(nil)

	//Need a Physics System too
	phys := &PhysicsSystem{VelocityIterations: 3, PositionIterations: 8}

	//Add some entities
	basics := make([]ecs.BasicEntity, 0)
	for i := 0; i < 5; i++ {
		basic := ecs.NewBasic()
		entity := collisionEntity{&basic, &common.SpaceComponent{}, &Box2dComponent{}}
		basics = append(basics, basic)
		entity.SpaceComponent = &common.SpaceComponent{
			Position: engo.Point{X: float32(i * 100), Y: 0},
			Width:    10,
			Height:   10,
		}
		entityBodyDef := box2d.NewB2BodyDef()
		entityBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
		entityBodyDef.Position = Conv.ToBox2d2Vec(entity.SpaceComponent.Center())
		entityBodyDef.Angle = Conv.DegToRad(entity.SpaceComponent.Rotation)
		entity.Box2dComponent.Body = World.CreateBody(entityBodyDef)
		var entityShape box2d.B2PolygonShape
		entityShape.SetAsBox(Conv.PxToMeters(entity.SpaceComponent.Width/2),
			Conv.PxToMeters(entity.SpaceComponent.Height/2))
		entityFixtureDef := box2d.B2FixtureDef{
			Shape:    &entityShape,
			Density:  1,
			Friction: 1,
		}
		entity.Box2dComponent.Body.CreateFixtureFromDef(&entityFixtureDef)
		sys.Add(entity.BasicEntity, entity.SpaceComponent, entity.Box2dComponent)
		phys.Add(entity.BasicEntity, entity.SpaceComponent, entity.Box2dComponent)
	}

	// Check that the system's entities are correct
	if len(sys.entities) != 5 {
		t.Errorf("sys.entities has wrong count, want: %d, got: %d", 5, len(sys.entities))
	}

	// Remove 2 and 3 from the system
	// queue up 2 and 3 for removal of bodies First
	for _, e := range sys.entities {
		if e.ID() == basics[2].ID() || e.ID() == basics[3].ID() {
			e.Box2dComponent.DestroyBody()
		}
	}
	sys.Remove(basics[2])
	sys.Remove(basics[3])
	phys.Remove(basics[2])
	phys.Remove(basics[3])

	// Check list of bodies to removes
	if len(listOfBodiesToRemove) != 2 {
		t.Errorf("listOfBodiesToRemove has wrong count, want: %d, got: %d", 2, len(listOfBodiesToRemove))
	}

	// Check the entity count is correct
	if len(sys.entities) != 3 {
		t.Errorf("sys.entities has wrong count after deletion, want: %d, got: %d", 3, len(sys.entities))
	}

	// Update the systems
	sys.Update(updateTime)
	phys.Update(updateTime)

	// Check the bodies were removed
	if len(listOfBodiesToRemove) != 0 {
		t.Errorf("listOfBodiesToRemove was not emptied after update, want: %d, got %d", 0, len(listOfBodiesToRemove))
	}

	// Change the space component of one of them
	toPoint := engo.Point{X: 15, Y: 0}
	for _, e := range sys.entities {
		if e.ID() == basics[1].ID() {
			e.SpaceComponent.Position = toPoint
			toPoint = e.SpaceComponent.Center()
		}
	}

	// Update the system
	sys.Update(updateTime)
	phys.Update(updateTime)

	// See if the Update changed the World Coordinates as well
	for _, e := range sys.entities {
		if e.ID() == basics[1].ID() {
			bodyPos := Conv.ToEngoPoint(e.Body.GetPosition())
			if bodyPos != toPoint {
				t.Errorf("Update did not change Body position to match Space position, want: %v, got %v", toPoint, bodyPos)
			}
		}
	}

	// Make sure no messages have been sent
	if recStartMessage || recEndMessage || recPreSolveMessage || recPostSolveMessage {
		t.Errorf("A collision was detected even though none have occured")
	}

	// Hit one so a collision occurs
	for _, e := range sys.entities {
		if e.ID() == basics[1].ID() {
			e.Body.ApplyLinearImpulseToCenter(box2d.B2Vec2{X: -10, Y: 0}, true)
		}
	}

	//Update systems
	sys.Update(updateTime)
	phys.Update(updateTime)

	//Update systems
	sys.Update(updateTime)
	phys.Update(updateTime)

	//Update systems
	sys.Update(updateTime)
	phys.Update(updateTime)

	if !recStartMessage {
		t.Errorf("did not recieve collision start message")
	}

	if !recPreSolveMessage {
		t.Errorf("did not recieve presolve message")
	}

	if !recPostSolveMessage {
		t.Errorf("did not recieve postsolve message")
	}

	// Hit one to stop collision
	for _, e := range sys.entities {
		if e.ID() == basics[1].ID() {
			e.Body.ApplyLinearImpulseToCenter(box2d.B2Vec2{X: 10, Y: 0}, true)
		}
	}

	//Update systems
	sys.Update(updateTime)
	phys.Update(updateTime)

	//Update systems
	sys.Update(updateTime)
	phys.Update(updateTime)

	if !recEndMessage {
		t.Errorf("did not recieve collision end message")
	}
}

func TestCollisionSystemAddByInterface(t *testing.T) {
	engo.Mailbox = &engo.MessageManager{}
	sys := &CollisionSystem{}
	sys.New(nil)

	phys := &PhysicsSystem{VelocityIterations: 3, PositionIterations: 8}

	basic := ecs.NewBasic()
	space := common.SpaceComponent{
		Width:  10,
		Height: 10,
	}
	entityBodyDef := box2d.NewB2BodyDef()
	entityBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	entityBodyDef.Position = Conv.ToBox2d2Vec(space.Center())
	entityBodyDef.Angle = Conv.DegToRad(space.Rotation)
	boxBody := World.CreateBody(entityBodyDef)
	var entityShape box2d.B2PolygonShape
	entityShape.SetAsBox(Conv.PxToMeters(space.Width/2),
		Conv.PxToMeters(space.Height/2))
	entityFixtureDef := box2d.B2FixtureDef{
		Shape:    &entityShape,
		Density:  1,
		Friction: 1,
	}
	boxBody.CreateFixtureFromDef(&entityFixtureDef)
	e := collisionEntity{&basic, &space, &Box2dComponent{Body: boxBody}}
	sys.AddByInterface(e)
	phys.AddByInterface(e)

	if len(sys.entities) != 1 {
		t.Errorf("AddByInterface failed for collision system; wanted %d, have %d", 1, len(sys.entities))
	}

	if len(phys.entities) != 1 {
		t.Errorf("AddByInterface failed for physics system; wanted %d, have %d", 1, len(sys.entities))
	}
}
