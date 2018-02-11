package engoBox2dSystem

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"

	"github.com/ByteArena/box2d"
)

type MouseTestScene struct {
	entityCount int
}

func (*MouseTestScene) Preload() {}

func (s *MouseTestScene) Setup(w *ecs.World) {
	// Add systems to the world
	sys = &MouseSystem{}
	w.AddSystem(&common.CameraSystem{})
	w.AddSystem(sys)

	//Add some entities
	basics = make([]ecs.BasicEntity, 0)
	for i := 0; i < s.entityCount; i++ {
		basic := ecs.NewBasic()
		basics = append(basics, basic)
		entity := mouseEntity{&basic, &MouseComponent{}, &common.SpaceComponent{}, nil, &Box2dComponent{}}
		entity.SpaceComponent = &common.SpaceComponent{
			Position: engo.Point{X: float32(i * 20), Y: 0},
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
		for _, system := range w.Systems() {
			switch sys := system.(type) {
			case *MouseSystem:
				sys.Add(entity.BasicEntity, entity.MouseComponent, entity.SpaceComponent, nil, entity.Box2dComponent)
			}
		}
	}
}

func (*MouseTestScene) Type() string { return "MouseTestScene" }

var (
	sys    *MouseSystem
	basics []ecs.BasicEntity
)

// Should be able to add entities to the system
func TestMouseSystemAdd(t *testing.T) {
	engo.Run(engo.RunOptions{
		Width:        100,
		Height:       100,
		NoRun:        true,
		HeadlessMode: true,
	}, &MouseTestScene{5})

	if len(sys.entities) != 5 {
		t.Errorf("Entity count does not match number added, have: %d, want: %d", len(sys.entities), 5)
	}
}

// Should be able to remove entities from the systems
func TestMouseSystemRemove(t *testing.T) {
	engo.Run(engo.RunOptions{
		Width:        100,
		Height:       100,
		NoRun:        true,
		HeadlessMode: true,
	}, &MouseTestScene{5})

	sys.Remove(basics[1])
	sys.Remove(basics[2])

	if len(sys.entities) != 3 {
		t.Errorf("Entity count does not match number removed, have %d, want %d", len(sys.entities), 3)
	}

	// removing an entity not in the system should do nothing
	sys.Remove(basics[2])

	if len(sys.entities) != 3 {
		t.Errorf("Removing an entity not in the system changed the count, have %d, want %d", len(sys.entities), 3)
	}

	// make sure only the IDs that were removed were the ones removed
	ids := make(map[uint64]struct{})
	ids[basics[0].ID()] = struct{}{}
	ids[basics[3].ID()] = struct{}{}
	ids[basics[4].ID()] = struct{}{}

	for _, e := range sys.entities {
		_, ok := ids[e.ID()]
		if !ok {
			t.Errorf("Entity in system has ID other than ones that should remain")
		}
	}
}

// when the system is added before a CameraSystem is, an error is logged
func TestMouseSystemNewNoCamera(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)

	expected := "ERROR: CameraSystem not found - have you added the `RenderSystem` before the `MouseSystem`?\n"

	w := ecs.World{}
	w.AddSystem(&MouseSystem{})

	if !strings.HasSuffix(str.String(), expected) {
		t.Errorf("Log did not recieve correct message suffix, have: %v, wanted (suffix): %v", str.String(), expected)
	}
}

// Test hovering
func TestMouseSystemUpdateHoverint(t *testing.T) {
	updateTime := float32(1.0 / 60.0)
	engo.Run(engo.RunOptions{
		Width:        100,
		Height:       100,
		NoRun:        true,
		HeadlessMode: true,
	}, &MouseTestScene{2})

	for i := 0; i < 2; i++ {
		//Place mouse inside entity
		engo.Input.Mouse.X = float32(5 + (20 * i))
		engo.Input.Mouse.Y = 5

		//Update the systems
		sys.Update(updateTime)

		//The one inside should be Entered for the first frame
		for _, e := range sys.entities {
			if e.ID() == basics[i].ID() {
				if !e.MouseComponent.Enter {
					t.Errorf("mouse component not updated to enter on frame mouse entered it; entity: %v", i)
				}
				if !e.MouseComponent.Hovered {
					t.Errorf("mouse component not updated to hover on frame mouse entered it; entity: %v", i)
				}
			} else {
				if e.MouseComponent.Enter {
					t.Errorf("mouse component says entered when it should not; entity: %v", i)
				}
				if e.MouseComponent.Hovered {
					t.Errorf("mouse component says hovered when it should not; entity: %v", i)
				}
			}
		}

		//Update again, this should remove the entered but keep Hovered
		sys.Update(updateTime)

		for _, e := range sys.entities {
			if e.ID() == basics[i].ID() {
				if e.MouseComponent.Enter {
					t.Errorf("mouse component not updated to not enter on second frame since mouse entered it; entity: %v", i)
				}
				if !e.MouseComponent.Hovered {
					t.Errorf("mouse component not updated to hover on second frame mouse since entered it; entity: %v", i)
				}
			} else {
				if e.MouseComponent.Enter {
					t.Errorf("mouse component says entered when it should not; entity: %v", i)
				}
				if e.MouseComponent.Hovered {
					t.Errorf("mouse component says hovered when it should not; entity: %v", i)
				}
			}
		}

		//Move mouse out
		engo.Input.Mouse.X = 20
		engo.Input.Mouse.Y = 20

		//Update the system, this should give exit signal
		sys.Update(updateTime)

		for _, e := range sys.entities {
			if e.ID() == basics[i].ID() {
				if !e.MouseComponent.Leave {
					t.Errorf("mouse component not updated to leave on frame mouse left it; entity: %v", i)
				}
				if e.MouseComponent.Hovered {
					t.Errorf("mouse component not updated to not hover on frame mouse left it; entity: %v", i)
				}
			} else {
				if e.MouseComponent.Leave {
					t.Errorf("mouse component says leave when it should not; entity: %v", i)
				}
				if e.MouseComponent.Hovered {
					t.Errorf("mouse component says hovered when it should not; entity: %v", i)
				}
			}
		}

		//One more update, should remove leavee too
		sys.Update(updateTime)

		for _, e := range sys.entities {
			if e.ID() == basics[i].ID() {
				if e.MouseComponent.Leave {
					t.Errorf("mouse component not updated to not leave on second frame since mouse left it; entity: %v", i)
				}
				if e.MouseComponent.Hovered {
					t.Errorf("mouse component not updated to not hover by second frame since mouse left it; entity: %v", i)
				}
			} else {
				if e.MouseComponent.Leave {
					t.Errorf("mouse component says leave when it should not; entity: %v", i)
				}
				if e.MouseComponent.Hovered {
					t.Errorf("mouse component says hovered when it should not; entity: %v", i)
				}
			}
		}
	}
}

// Camera movements
func TestMouseSystemCameraMove(t *testing.T) {
	updateTime := float32(1.0 / 60.0)
	engo.Run(engo.RunOptions{
		Width:        100,
		Height:       100,
		NoRun:        true,
		HeadlessMode: true,
	}, &MouseTestScene{2})

	basic := ecs.NewBasic()
	space := &common.SpaceComponent{
		Position: engo.Point{X: 0, Y: 0},
		Rotation: 0,
	}
	sys.camera.FollowEntity(&basic, space)

	sys.camera.Update(updateTime)
	sys.Update(updateTime)

	space.Rotation = 45

	sys.camera.Update(updateTime)
	sys.Update(updateTime)

}
