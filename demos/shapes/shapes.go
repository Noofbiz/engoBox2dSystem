package main

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/Noofbiz/engoBox2dSystem"

	"github.com/ByteArena/box2d"
)

type defaultScene struct{}

type guy struct {
	ecs.BasicEntity

	common.RenderComponent
	common.SpaceComponent
	engoBox2dSystem.Box2dComponent
	engoBox2dSystem.MouseComponent
}

type wall struct {
	ecs.BasicEntity
	common.SpaceComponent
	engoBox2dSystem.Box2dComponent
}

func (*defaultScene) Preload() {
	engo.Files.Load("food.png")
}

func (*defaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&engoBox2dSystem.MouseSystem{})
	w.AddSystem(&controlSystem{})

	//setup sprite sheet
	var spriteRegions []common.SpriteRegion
	//apple
	spriteRegions = append(spriteRegions, common.SpriteRegion{
		Position: engo.Point{X: 0, Y: 0},
		Width:    28,
		Height:   32,
	})
	//cheese
	spriteRegions = append(spriteRegions, common.SpriteRegion{
		Position: engo.Point{X: 36, Y: 2},
		Width:    38,
		Height:   33,
	})
	//lemon
	spriteRegions = append(spriteRegions, common.SpriteRegion{
		Position: engo.Point{X: 80, Y: 4},
		Width:    36,
		Height:   31,
	})
	//carrot
	spriteRegions = append(spriteRegions, common.SpriteRegion{
		Position: engo.Point{X: 119, Y: 1},
		Width:    40,
		Height:   40,
	})
	//steak
	spriteRegions = append(spriteRegions, common.SpriteRegion{
		Position: engo.Point{X: 160, Y: 0},
		Width:    37,
		Height:   35,
	})
	//grapes
	spriteRegions = append(spriteRegions, common.SpriteRegion{
		Position: engo.Point{X: 205, Y: 3},
		Width:    36,
		Height:   39,
	})
	sprites := common.NewAsymmetricSpritesheetFromFile("food.png", spriteRegions)

	//add Apple
	appleTexture := sprites.Drawable(0)
	apple := guy{BasicEntity: ecs.NewBasic()}
	apple.RenderComponent = common.RenderComponent{
		Drawable: appleTexture,
		Scale:    engo.Point{X: 2, Y: 2},
	}
	apple.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 20, Y: 20},
		Width:    appleTexture.Width() * apple.RenderComponent.Scale.X,
		Height:   appleTexture.Height() * apple.RenderComponent.Scale.Y,
	}

	// apple's box2d Body
	appleBodyDef := box2d.NewB2BodyDef()
	appleBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	appleBodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(apple.SpaceComponent.Center())
	appleBodyDef.Angle = engoBox2dSystem.Conv.DegToRad(apple.SpaceComponent.Rotation)
	apple.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(appleBodyDef)
	var appleShape1 box2d.B2CircleShape
	appleShape1.SetRadius(engoBox2dSystem.Conv.PxToMeters(25.5))
	appleShape1.M_p.Set(engoBox2dSystem.Conv.PxToMeters(-1.5), engoBox2dSystem.Conv.PxToMeters(4.5))
	appleFixture1Def := box2d.B2FixtureDef{
		Shape:    appleShape1,
		Density:  1.0,
		Friction: 0.5,
	}
	apple.Body.CreateFixtureFromDef(&appleFixture1Def)
	var appleShape2 box2d.B2PolygonShape
	var appleShape2Verts []box2d.B2Vec2
	appleShape2Verts = append(appleShape2Verts,
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(-15), Y: engoBox2dSystem.Conv.PxToMeters(-19)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(-20), Y: engoBox2dSystem.Conv.PxToMeters(-24)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(-20), Y: engoBox2dSystem.Conv.PxToMeters(-29)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(-17), Y: engoBox2dSystem.Conv.PxToMeters(-32)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(-11), Y: engoBox2dSystem.Conv.PxToMeters(-32)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(-1), Y: engoBox2dSystem.Conv.PxToMeters(-27)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(-1), Y: engoBox2dSystem.Conv.PxToMeters(-19)})
	appleShape2.Set(appleShape2Verts, 7)
	appleShape2.M_centroid.Set(engoBox2dSystem.Conv.PxToMeters(-18.5), engoBox2dSystem.Conv.PxToMeters(-25.5))
	appleFixture2Def := box2d.B2FixtureDef{
		Shape:    &appleShape2,
		Density:  1.0,
		Friction: 0.5,
	}
	apple.Body.CreateFixtureFromDef(&appleFixture2Def)
	var appleShape3 box2d.B2PolygonShape
	var appleShape3Verts []box2d.B2Vec2
	appleShape3Verts = append(appleShape3Verts,
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(2), Y: engoBox2dSystem.Conv.PxToMeters(-21)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(4), Y: engoBox2dSystem.Conv.PxToMeters(-28)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(6), Y: engoBox2dSystem.Conv.PxToMeters(-28)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(8), Y: engoBox2dSystem.Conv.PxToMeters(-24)})
	appleShape3.Set(appleShape3Verts, 4)
	appleShape3.M_centroid.Set(engoBox2dSystem.Conv.PxToMeters(5), engoBox2dSystem.Conv.PxToMeters(-24.5))
	appleFixture3Def := box2d.B2FixtureDef{
		Shape:    &appleShape3,
		Density:  1.0,
		Friction: 0.5,
	}
	apple.Body.CreateFixtureFromDef(&appleFixture3Def)
	var appleShape4 box2d.B2PolygonShape
	var appleShape4Verts []box2d.B2Vec2
	appleShape4Verts = append(appleShape4Verts,
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(9), Y: engoBox2dSystem.Conv.PxToMeters(-24)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(11), Y: engoBox2dSystem.Conv.PxToMeters(-26)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(15), Y: engoBox2dSystem.Conv.PxToMeters(-26)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(15), Y: engoBox2dSystem.Conv.PxToMeters(-21)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(12), Y: engoBox2dSystem.Conv.PxToMeters(-18)})
	appleShape4.Set(appleShape4Verts, 5)
	appleShape4.M_centroid.Set(engoBox2dSystem.Conv.PxToMeters(12), engoBox2dSystem.Conv.PxToMeters(-22))
	appleFixture4Def := box2d.B2FixtureDef{
		Shape:    &appleShape4,
		Density:  1.0,
		Friction: 0.5,
	}
	apple.Body.CreateFixtureFromDef(&appleFixture4Def)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&apple.BasicEntity, &apple.RenderComponent, &apple.SpaceComponent)
		case *engoBox2dSystem.MouseSystem:
			sys.Add(&apple.BasicEntity, &apple.MouseComponent, &apple.SpaceComponent, &apple.RenderComponent, &apple.Box2dComponent)
		case *controlSystem:
			sys.Add(&apple.BasicEntity, &apple.SpaceComponent, &apple.MouseComponent)
		}
	}

	//add Cheese
	cheeseTexture := sprites.Drawable(1)
	cheese := guy{BasicEntity: ecs.NewBasic()}
	cheese.RenderComponent = common.RenderComponent{
		Drawable: cheeseTexture,
		Scale:    engo.Point{X: 2, Y: 2},
	}
	cheese.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 420, Y: 20},
		Width:    cheeseTexture.Width() * cheese.RenderComponent.Scale.X,
		Height:   cheeseTexture.Height() * cheese.RenderComponent.Scale.Y,
	}

	// cheese's box2d Body
	cheeseBodyDef := box2d.NewB2BodyDef()
	cheeseBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	cheeseBodyDef.Position = engoBox2dSystem.Conv.ToBox2d2Vec(cheese.SpaceComponent.Center())
	cheeseBodyDef.Angle = engoBox2dSystem.Conv.DegToRad(cheese.SpaceComponent.Rotation)
	cheese.Box2dComponent.Body = engoBox2dSystem.World.CreateBody(cheeseBodyDef)
	var cheeseShape box2d.B2PolygonShape
	var cheeseShapeVerts []box2d.B2Vec2
	cheeseShapeVerts = append(cheeseShapeVerts,
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(-33), Y: engoBox2dSystem.Conv.PxToMeters(-29)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(-9), Y: engoBox2dSystem.Conv.PxToMeters(-29)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(35), Y: engoBox2dSystem.Conv.PxToMeters(-13)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(35), Y: engoBox2dSystem.Conv.PxToMeters(18)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(6), Y: engoBox2dSystem.Conv.PxToMeters(30)},
		box2d.B2Vec2{X: engoBox2dSystem.Conv.PxToMeters(-33), Y: engoBox2dSystem.Conv.PxToMeters(11)})
	cheeseShape.Set(cheeseShapeVerts, 6)
	cheeseShape.M_centroid.Set(engoBox2dSystem.Conv.PxToMeters(1), engoBox2dSystem.Conv.PxToMeters(-0.5))
	cheeseFixtureDef := box2d.B2FixtureDef{
		Shape:    &cheeseShape,
		Density:  1.0,
		Friction: 0.5,
	}
	cheese.Body.CreateFixtureFromDef(&cheeseFixtureDef)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&cheese.BasicEntity, &cheese.RenderComponent, &cheese.SpaceComponent)
		case *engoBox2dSystem.MouseSystem:
			sys.Add(&cheese.BasicEntity, &cheese.MouseComponent, &cheese.SpaceComponent, &cheese.RenderComponent, &cheese.Box2dComponent)
		case *controlSystem:
			sys.Add(&cheese.BasicEntity, &cheese.SpaceComponent, &cheese.MouseComponent)
		}
	}
}

type controlSystem struct {
	entities []controlEntity
}

type controlEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*engoBox2dSystem.MouseComponent
	*controlComponent
}

type controlComponent struct {
	following  bool
	xoff, yoff float32
}

func (c *controlSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent, mouse *engoBox2dSystem.MouseComponent) {
	c.entities = append(c.entities, controlEntity{basic, space, mouse, &controlComponent{}})
}

func (c *controlSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, entity := range c.entities {
		if entity.ID() == basic.ID() {
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
		if e.MouseComponent.Clicked {
			e.following = true
			e.xoff = engo.Input.Mouse.X - e.SpaceComponent.Position.X
			e.yoff = engo.Input.Mouse.Y - e.SpaceComponent.Position.Y
		}
		if e.MouseComponent.Released {
			e.following = false
			e.xoff = 0
			e.yoff = 0
		}
		if e.following {
			e.SpaceComponent.Position.Set(engo.Input.Mouse.X-e.xoff, engo.Input.Mouse.Y-e.yoff)
		}
	}
}

func (*defaultScene) Type() string { return "GameWorld" }

func main() {
	opts := engo.RunOptions{
		Title:       "Box2d-Engo Mouse and Shape Demo",
		Width:       1024,
		Height:      640,
		GlobalScale: engo.Point{2, 2},
	}
	engo.Run(opts, &defaultScene{})
}
