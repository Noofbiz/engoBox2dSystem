# Gravity Demo

## What does it do?
It shows how to use the Physics system, how to create Box2d bodies and add them to an entity, and how to add a gravity vector to the physics system.

## What are the important aspects of the code?

##### Adding a Physics System to the World

`w.AddSystem(&engoBox2dSystem.PhysicsSystem{VelocityIterations: 3, PositionIterations: 8})`

when you add an &engoBox2dSystem.PhysicsSystem to the world, you have to also specify the VelocityIterations and the Position iterations that the Box2d World will use during simulations.

##### Adding gravity to the World

`engoBox2dSystem.World.SetGravity(box2d.B2Vec2{X: 0, Y: 10})`

This sets the gravity used by the world. Remember, `engo` renders from the top left corner of the screen, so for gravity to point down, it should be in the positive Y direction.

##### Adding a Box2dBody to an entity

First the entity declaration looks like
```go
type guy struct {
	ecs.BasicEntity

	common.RenderComponent
	common.SpaceComponent
	engoBox2dSystem.Box2dComponent
}
```

Then to create a box2d body you'd first need a body definition

`dudeBodayDef := box2d.NewB2BodyDef()`

Set the type of body to dynamic if you want it to move

`dudeBodyDef.Type = box2d.B2BodyType.B2_dynamicBody`

Set the position and angle of the body. Note that Box2d is based on the CENTER of the body. You also need to convert between the render system dimensions and the Box2d dimensions with `engoBox2dSystem.Conv`

`dudeBodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(dude.Center())`
`dudeBodyDef.Antle = engoBox2dSystem.Conv.DegToRad(dude.Rotation)`

Once you have the body complete, add it to the entitty and the World

`dude.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(dudeBodyDef)`

The body only keeps track of the position, angle, and type of the entity. It does not know anything about the shape, size, joints, fixtures, or anything else. We'll need to add those separately.

```go
var dudeBodyShape box2d.B2PolygonShape
dudeBodyShape.SetAsBox(engoBox2dSystem.Conv.PxToMeters(dude.SpaceComponent.Width/2),
	engoBox2dSystem.Conv.PxToMeters(dude.SpaceComponent.Height/2))
dudeFixtureDef := box2d.B2FixtureDef{
	Shape:    &dudeBodyShape,
	Density:  1.0,
	Friction: 0.1,
}
dude.Box2dComponent.Body.CreateFixtureFromDef(&dudeFixtureDef)
```

Last, add it to the PhysicsSystem:

```go
case *engoBox2dSystem.PhysicsSystem:
	sys.Add(&grass.BasicEntity, &grass.SpaceComponent, &grass.Box2dComponent)
```
