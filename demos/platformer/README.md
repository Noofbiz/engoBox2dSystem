# Platformer Demo

## What does it do?
Stars fall from the sky. You can collect the stars and get points for each one you collect.

## What are the important aspects of the code?

##### Adding your Player to the systems

```go
	dudeBodyDef := box2d.NewB2BodyDef()
	dudeBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	dudeBodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(dude.SpaceComponent.Center())
	dudeBodyDef.Angle = engoBox2dSystem.Conv.DegToRad(dude.SpaceComponent.Rotation)
	dudeBodyDef.FixedRotation = true
	dude.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(dudeBodyDef)
	var dudeBodyShape box2d.B2PolygonShape
	dudeBodyShape.SetAsBox(engoBox2dSystem.Conv.PxToMeters(dude.SpaceComponent.Width/2),
		engoBox2dSystem.Conv.PxToMeters(dude.SpaceComponent.Height/2))
	dudeFixtureDef := box2d.B2FixtureDef{
		Shape:    &dudeBodyShape,
		Density:  1.0,
		Friction: 0.1,
	}
	dude.Box2dComponent.Body.CreateFixtureFromDef(&dudeFixtureDef)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&dude.BasicEntity, &dude.RenderComponent, &dude.SpaceComponent)
		case *engoBox2dSystem.PhysicsSystem:
			sys.Add(&dude.BasicEntity, &dude.SpaceComponent, &dude.Box2dComponent)
		case *controlSystem:
			sys.Add(&dude.BasicEntity, &dude.Box2dComponent)
		case *engoBox2dSystem.CollisionSystem:
			sys.Add(&dude.BasicEntity, &dude.SpaceComponent, &dude.Box2dComponent)
		case *starCollectionSystem:
			sys.Add(dude.BasicEntity, &dude.Box2dComponent, true)
		}
	}
```

##### Adding the floors

```go
	grassBodyDef := box2d.NewB2BodyDef()
	grassBodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(grass.SpaceComponent.Center())
	grassBodyDef.Angle = engoBox2dSystem.Conv.DegToRad(grass.SpaceComponent.Rotation)
	grass.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(grassBodyDef)
	var grassBodyShape box2d.B2PolygonShape
	grassBodyShape.SetAsBox(engoBox2dSystem.Conv.PxToMeters(grass.SpaceComponent.Width/2),
		engoBox2dSystem.Conv.PxToMeters(grass.SpaceComponent.Height/2))
	grassFixtureDef := box2d.B2FixtureDef{Shape: &grassBodyShape}
	grass.Box2dComponent.Body.CreateFixtureFromDef(&grassFixtureDef)

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&grass.BasicEntity, &grass.RenderComponent, &grass.SpaceComponent)
		case *engoBox2dSystem.PhysicsSystem:
			sys.Add(&grass.BasicEntity, &grass.SpaceComponent, &grass.Box2dComponent)
		}
	}
```

##### Adding Walls so you can't leave the play area

```go
	leftWallBodyDef := box2d.NewB2BodyDef()
	leftWallBodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(leftWall.SpaceComponent.Center())
	leftWallBodyDef.Angle = engoBox2dSystem.Conv.DegToRad(leftWall.SpaceComponent.Rotation)
	leftWall.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(leftWallBodyDef)
	var leftWallBodyShape box2d.B2PolygonShape
	leftWallBodyShape.SetAsBox(engoBox2dSystem.Conv.PxToMeters(leftWall.SpaceComponent.Width/2),
		engoBox2dSystem.Conv.PxToMeters(leftWall.SpaceComponent.Height/2))
	leftWallFixtureDef := box2d.B2FixtureDef{Shape: &leftWallBodyShape}
	leftWall.Box2dComponent.Body.CreateFixtureFromDef(&leftWallFixtureDef)

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engoBox2dSystem.PhysicsSystem:
			sys.Add(&leftWall.BasicEntity, &leftWall.SpaceComponent, &leftWall.Box2dComponent)
		}
	}
```


##### Adding stars with a pentagon shaped hitbox

```go
		star := guy{BasicEntity: ecs.NewBasic()}
		// Initialize the components, set scale to 8x
		star.RenderComponent = common.RenderComponent{
			Drawable: starTexture,
		}
		star.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{
				X: 0 + float32(150*i),
				Y: 0,
			},
			Width:  starTexture.Width(),
			Height: starTexture.Height(),
		}

		//box2d component setup
		starBodyDef := box2d.NewB2BodyDef()
		starBodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(star.SpaceComponent.Center())
		starBodyDef.Angle = engoBox2dSystem.Conv.DegToRad(star.SpaceComponent.Rotation)
		starBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
		star.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(starBodyDef)
		var starBodyShape box2d.B2PolygonShape
		var vertices []box2d.B2Vec2
		vertices = append(vertices, box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(0), Y: engoBox2dSystem.Conv.PxToMeters(-50)})
		vertices = append(vertices, box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(49), Y: engoBox2dSystem.Conv.PxToMeters(-15)})
		vertices = append(vertices, box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(30), Y: engoBox2dSystem.Conv.PxToMeters(41)})
		vertices = append(vertices, box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(-31), Y: engoBox2dSystem.Conv.PxToMeters(41)})
		vertices = append(vertices, box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(-50), Y: engoBox2dSystem.Conv.PxToMeters(-15)})
		starBodyShape.Set(vertices, 5)
		starFixtureDef := box2d.B2FixtureDef{Shape: &starBodyShape}
		star.Box2dComponent.Body.CreateFixtureFromDef(&starFixtureDef)

		// Add it to appropriate systems
		for _, system := range w.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&star.BasicEntity, &star.RenderComponent, &star.SpaceComponent)
			case *engoBox2dSystem.PhysicsSystem:
				sys.Add(&star.BasicEntity, &star.SpaceComponent, &star.Box2dComponent)
			case *engoBox2dSystem.CollisionSystem:
				sys.Add(&star.BasicEntity, &star.SpaceComponent, &star.Box2dComponent)
			case *starCollectionSystem:
				sys.Add(star.BasicEntity, &star.Box2dComponent, false)
			}
		}
```

##### Movement via Box2d
In the controlSystem, we have the Guy move by applying impulses in the Box2d world instead of updating the space component directily

```go
func (c *controlSystem) Update(dt float32) {
	for _, e := range c.entities {
		if engo.Input.Button("up").JustPressed() {
			e.Body.ApplyLinearImpulseToCenter(box2d.B2Vec2{X: 0, Y: -500}, true)
		} else if engo.Input.Button("left").JustPressed() {
			e.Body.SetLinearVelocity(box2d.B2Vec2{X: -10, Y: e.Body.GetLinearVelocity().Y})
		} else if engo.Input.Button("right").JustPressed() {
			e.Body.SetLinearVelocity(box2d.B2Vec2{X: 10, Y: e.Body.GetLinearVelocity().Y})
		}
	}
}
```

##### Collision Detection via Box2d
In the starCollectionSystem, we listen for the CollisionStartMessage and destroy the star and add to our score if they're touching

```go
	engo.Mailbox.Listen("CollisionStartMessage", func(message engo.Message) {
		c, isCollision := message.(engoBox2dSystem.CollisionStartMessage)
		if isCollision {
			if c.Contact.IsTouching() {
				a := c.Contact.GetFixtureA().GetBody().GetUserData()
				b := c.Contact.GetFixtureB().GetBody().GetUserData()
				for i1, e1 := range s.entities {
					if !e1.isGuy {
						continue
					}
					if e1.BasicEntity.ID() == a || e1.BasicEntity.ID() == b {
						for i2, e2 := range s.entities {
							if i1 == i2 {
								continue
							}
							if e2.BasicEntity.ID() == a || e2.BasicEntity.ID() == b {
								// Remove it from all the appropriate systems
								for _, system := range w.Systems() {
									switch sys := system.(type) {
									case *common.RenderSystem:
										sys.Remove(e2.BasicEntity)
									case *engoBox2dSystem.PhysicsSystem:
										sys.Remove(e2.BasicEntity)
									case *engoBox2dSystem.CollisionSystem:
										sys.Remove(e2.BasicEntity)
									case *starCollectionSystem:
										sys.Remove(e2.BasicEntity)
									}
								}
								e2.Box2dComponent.DestroyBody()
								engo.Mailbox.Dispatch(scoreMessage{})
							}
						}
					}
				}
			}
		}
	})
```
