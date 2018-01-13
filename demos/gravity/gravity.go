package main

import (
	"image/color"
	"log"

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

func (*defaultScene) Preload() {
	engo.Files.Load("icon.png", "grass.png")
}

func (*defaultScene) Setup(w *ecs.World) {
	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})

	//add box2d systems
	w.AddSystem(&engoBox2dSystem.PhysicsSystem{VelocityIterations: 3, PositionIterations: 8})

	//add downward gravity to box2d world
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
		X: 512,
		Y: 0,
	})

	//box2d component setup
	dudeBodyDef := box2d.NewB2BodyDef()
	dudeBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	dudeBodyDef.Position = engoBox2dSystem.TheConverter.ToBox2d2Vec(dude.Center())
	dudeBodyDef.Angle = engoBox2dSystem.TheConverter.DegToRad(dude.Rotation)
	dude.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(dudeBodyDef)
	var dudeBodyShape box2d.B2PolygonShape
	dudeBodyShape.SetAsBox(engoBox2dSystem.TheConverter.PxToMeters(dude.SpaceComponent.Width/2),
		engoBox2dSystem.TheConverter.PxToMeters(dude.SpaceComponent.Height/2))
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
	grassBodyDef.Position = engoBox2dSystem.TheConverter.ToBox2d2Vec(grass.Center())
	grassBodyDef.Angle = engoBox2dSystem.TheConverter.DegToRad(grass.Rotation)
	grass.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(grassBodyDef)
	var grassBodyShape box2d.B2PolygonShape
	grassBodyShape.SetAsBox(engoBox2dSystem.TheConverter.PxToMeters(grass.SpaceComponent.Width/2),
		engoBox2dSystem.TheConverter.PxToMeters(grass.SpaceComponent.Height/2))
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
}

func (*defaultScene) Type() string { return "GameWorld" }

func main() {
	opts := engo.RunOptions{
		Title:  "Box2d Engo Gravity Demo",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &defaultScene{})
}
