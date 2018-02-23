package main

import (
	"image/color"
	"log"
	"strconv"
	"sync"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"

	"github.com/Noofbiz/engoBox2dSystem"

	"github.com/ByteArena/box2d"
)

type defaultScene struct{}

type guy struct {
	ecs.BasicEntity

	common.RenderComponent
	common.SpaceComponent
	engoBox2dSystem.Box2dComponent
}

type wall struct {
	ecs.BasicEntity
	common.SpaceComponent
	engoBox2dSystem.Box2dComponent
}

var fnt *common.Font

func (*defaultScene) Preload() {
	engo.Files.Load("icon.png", "grass.png", "star_gold.png", "Oswald-SemiBold.ttf")

	engo.Input.RegisterButton("up", engo.W, engo.ArrowUp, engo.Space)
	engo.Input.RegisterButton("left", engo.A, engo.ArrowLeft)
	engo.Input.RegisterButton("right", engo.D, engo.ArrowRight)
}

func (*defaultScene) Setup(w *ecs.World) {
	bg := color.RGBA{R: 135, G: 206, B: 235}
	common.SetBackground(bg)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&controlSystem{})

	//add box2d systems
	w.AddSystem(&engoBox2dSystem.PhysicsSystem{VelocityIterations: 3, PositionIterations: 8})
	w.AddSystem(&engoBox2dSystem.CollisionSystem{})
	w.AddSystem(&starCollectionSystem{})
	w.AddSystem(&scoreSystem{})

	//Set gravity to point down
	engoBox2dSystem.World.SetGravity(box2d.B2Vec2{X: 0, Y: 10})

	// Guy Texture
	dudeTexture, err := common.LoadedSprite("icon.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	dude := guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	dude.RenderComponent = common.RenderComponent{
		Drawable: dudeTexture,
		Scale: engo.Point{
			X: 8,
			Y: 8,
		},
	}
	dude.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{
			X: 0,
			Y: 0,
		},
		Width:  dudeTexture.Width() * dude.RenderComponent.Scale.X,
		Height: dudeTexture.Height() * dude.RenderComponent.Scale.Y,
	}
	dude.SpaceComponent.SetCenter(engo.Point{
		X: 64,
		Y: 512,
	})

	//box2d component setup
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

	// Add it to appropriate systems
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

	// Grass Texture
	grassTexture, err := common.LoadedSprite("grass.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	grass := guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	grass.RenderComponent = common.RenderComponent{
		Drawable: grassTexture,
		Scale: engo.Point{
			X: 8,
			Y: 1,
		},
	}
	grass.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{
			X: 0,
			Y: 0,
		},
		Width:  grassTexture.Width() * grass.RenderComponent.Scale.X,
		Height: grassTexture.Height() * grass.RenderComponent.Scale.Y,
	}
	grass.SpaceComponent.SetCenter(engo.Point{
		X: 512,
		Y: 608,
	})

	//box2d component setup
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

	// First Floor
	grass1 := guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 3x
	grass1.RenderComponent = common.RenderComponent{
		Drawable: grassTexture,
		Scale: engo.Point{
			X: 3,
			Y: 1,
		},
	}
	grass1.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{
			X: 0,
			Y: 0,
		},
		Width:  grassTexture.Width() * grass1.RenderComponent.Scale.X,
		Height: grassTexture.Height() * grass1.RenderComponent.Scale.Y,
	}
	grass1.SpaceComponent.SetCenter(engo.Point{
		X: 192,
		Y: 350,
	})

	//box2d component setup
	grass1BodyDef := box2d.NewB2BodyDef()
	grass1BodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(grass1.SpaceComponent.Center())
	grass1BodyDef.Angle = engoBox2dSystem.Conv.DegToRad(grass1.SpaceComponent.Rotation)
	grass1.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(grass1BodyDef)
	var grass1BodyShape box2d.B2PolygonShape
	grass1BodyShape.SetAsBox(engoBox2dSystem.Conv.PxToMeters(grass1.SpaceComponent.Width/2),
		engoBox2dSystem.Conv.PxToMeters(grass1.SpaceComponent.Height/2))
	grass1FixtureDef := box2d.B2FixtureDef{Shape: &grass1BodyShape}
	grass1.Box2dComponent.Body.CreateFixtureFromDef(&grass1FixtureDef)

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&grass1.BasicEntity, &grass1.RenderComponent, &grass1.SpaceComponent)
		case *engoBox2dSystem.PhysicsSystem:
			sys.Add(&grass1.BasicEntity, &grass1.SpaceComponent, &grass1.Box2dComponent)
		}
	}

	// Second Floor
	grass2 := guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 3x
	grass2.RenderComponent = common.RenderComponent{
		Drawable: grassTexture,
		Scale: engo.Point{
			X: 3,
			Y: 1,
		},
	}
	grass2.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{
			X: 0,
			Y: 0,
		},
		Width:  grassTexture.Width() * grass2.RenderComponent.Scale.X,
		Height: grassTexture.Height() * grass2.RenderComponent.Scale.Y,
	}
	grass2.SpaceComponent.SetCenter(engo.Point{
		X: 832,
		Y: 200,
	})

	//box2d component setup
	grass2BodyDef := box2d.NewB2BodyDef()
	grass2BodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(grass2.SpaceComponent.Center())
	grass2BodyDef.Angle = engoBox2dSystem.Conv.DegToRad(grass2.SpaceComponent.Rotation)
	grass2.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(grass2BodyDef)
	var grass2BodyShape box2d.B2PolygonShape
	grass2BodyShape.SetAsBox(engoBox2dSystem.Conv.PxToMeters(grass2.SpaceComponent.Width/2),
		engoBox2dSystem.Conv.PxToMeters(grass2.SpaceComponent.Height/2))
	grass2FixtureDef := box2d.B2FixtureDef{Shape: &grass2BodyShape}
	grass2.Box2dComponent.Body.CreateFixtureFromDef(&grass2FixtureDef)

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&grass2.BasicEntity, &grass2.RenderComponent, &grass2.SpaceComponent)
		case *engoBox2dSystem.PhysicsSystem:
			sys.Add(&grass2.BasicEntity, &grass2.SpaceComponent, &grass2.Box2dComponent)
		}
	}

	// Create an entity so the guy can't fall out of bounds
	leftWall := wall{BasicEntity: ecs.NewBasic()}

	leftWall.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: -10, Y: -10},
		Width:    10,
		Height:   660,
	}

	//box2d component setup
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

	// Create an entity so the guy can't fall out of bounds
	rightWall := wall{BasicEntity: ecs.NewBasic()}

	rightWall.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 1024, Y: -10},
		Width:    10,
		Height:   660,
	}

	//box2d component setup
	rightWallBodyDef := box2d.NewB2BodyDef()
	rightWallBodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(rightWall.SpaceComponent.Center())
	rightWallBodyDef.Angle = engoBox2dSystem.Conv.DegToRad(rightWall.SpaceComponent.Rotation)
	rightWall.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(rightWallBodyDef)
	var rightWallBodyShape box2d.B2PolygonShape
	rightWallBodyShape.SetAsBox(engoBox2dSystem.Conv.PxToMeters(rightWall.SpaceComponent.Width/2),
		engoBox2dSystem.Conv.PxToMeters(rightWall.SpaceComponent.Height/2))
	rightWallFixtureDef := box2d.B2FixtureDef{Shape: &rightWallBodyShape}
	rightWall.Box2dComponent.Body.CreateFixtureFromDef(&rightWallFixtureDef)

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engoBox2dSystem.PhysicsSystem:
			sys.Add(&rightWall.BasicEntity, &rightWall.SpaceComponent, &rightWall.Box2dComponent)
		}
	}

	// Create an entity so the guy can't fall out of bounds
	ceil := wall{BasicEntity: ecs.NewBasic()}

	ceil.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: -10},
		Width:    1024,
		Height:   10,
	}

	//box2d component setup
	ceilBodyDef := box2d.NewB2BodyDef()
	ceilBodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(ceil.SpaceComponent.Center())
	ceilBodyDef.Angle = engoBox2dSystem.Conv.DegToRad(ceil.SpaceComponent.Rotation)
	ceil.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(ceilBodyDef)
	var ceilBodyShape box2d.B2PolygonShape
	ceilBodyShape.SetAsBox(engoBox2dSystem.Conv.PxToMeters(ceil.SpaceComponent.Width/2),
		engoBox2dSystem.Conv.PxToMeters(ceil.SpaceComponent.Height/2))
	ceilFixtureDef := box2d.B2FixtureDef{Shape: &ceilBodyShape}
	ceil.Box2dComponent.Body.CreateFixtureFromDef(&ceilFixtureDef)

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engoBox2dSystem.PhysicsSystem:
			sys.Add(&ceil.BasicEntity, &ceil.SpaceComponent, &ceil.Box2dComponent)
		}
	}

	starTexture, err := common.LoadedSprite("star_gold.png")
	if err != nil {
		log.Println(err)
	}

	for i := 0; i < 7; i++ {
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
	}

	fnt = &common.Font{
		URL:  "Oswald-SemiBold.ttf",
		FG:   color.Black,
		Size: 64,
	}
	fnt.CreatePreloaded()
	scoreLabel := guy{BasicEntity: ecs.NewBasic()}
	scoreLabel.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: "Score  0",
	}
	scoreLabel.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{
			X: 724,
			Y: 5,
		},
	}
	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&scoreLabel.BasicEntity, &scoreLabel.RenderComponent, &scoreLabel.SpaceComponent)
		case *scoreSystem:
			sys.Add(&scoreLabel.BasicEntity, &scoreLabel.SpaceComponent, &scoreLabel.RenderComponent)
		}
	}
}

func (*defaultScene) Type() string { return "GameWorld" }

type controlSystem struct {
	entities []controlEntity
}

type controlEntity struct {
	*ecs.BasicEntity
	*engoBox2dSystem.Box2dComponent
}

func (c *controlSystem) Add(basic *ecs.BasicEntity, box *engoBox2dSystem.Box2dComponent) {
	c.entities = append(c.entities, controlEntity{basic, box})
}

func (c *controlSystem) Remove(basic ecs.BasicEntity) {
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

type starCollectionSystem struct {
	entities []starCollectionEntity
}

type starCollectionEntity struct {
	ecs.BasicEntity
	*engoBox2dSystem.Box2dComponent

	isGuy bool
}

func (s *starCollectionSystem) New(w *ecs.World) {
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
								e2.Box2dComponent.DestroyBody()
								w.RemoveEntity(e2.BasicEntity)
								engo.Mailbox.Dispatch(scoreMessage{})
							}
						}
					}
				}
			}
		}
	})
}

func (s *starCollectionSystem) Add(basic ecs.BasicEntity, box *engoBox2dSystem.Box2dComponent, isGuy bool) {
	s.entities = append(s.entities, starCollectionEntity{basic, box, isGuy})
}

func (s *starCollectionSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *starCollectionSystem) Update(dt float32) {
}

type scoreMessage struct{}

func (scoreMessage) Type() string { return "scoreMessage" }

type scoreSystem struct {
	entities []scoreEntity

	score     int
	scoreLock sync.RWMutex
}

type scoreEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*common.RenderComponent
}

func (s *scoreSystem) New(w *ecs.World) {
	engo.Mailbox.Listen("scoreMessage", func(message engo.Message) {
		_, isScore := message.(scoreMessage)
		if isScore {
			s.scoreLock.Lock()
			s.score += 100
			s.scoreLock.Unlock()
		}
	})
}

func (s *scoreSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent, render *common.RenderComponent) {
	s.entities = append(s.entities, scoreEntity{basic, space, render})
}

func (s *scoreSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *scoreSystem) Update(dt float32) {
	for _, e := range s.entities {
		s.scoreLock.RLock()
		scoreStr := strconv.Itoa(s.score)
		s.scoreLock.RUnlock()
		e.RenderComponent.Drawable.Close()
		e.RenderComponent.Drawable = common.Text{
			Font: fnt,
			Text: "Score  " + scoreStr,
		}
		e.SpaceComponent.Position = engo.Point{
			X: 724,
			Y: 5,
		}
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Box2d-Engo Platformer Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &defaultScene{})
}
