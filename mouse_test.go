package engoBox2dSystem

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"engo.io/engo/math"
	"engo.io/gl"

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
	}, &MouseTestScene{1})

	shifts := [][]float32{
		[]float32{5, 5, 0, 0, 0, 15},
		[]float32{5, 5, 5, 0, 0, 0},
		[]float32{5, 5, 0, 5, 0, 0},
		[]float32{5, 5, 0, 0, 2, 0},
		[]float32{5, 5, 5, 5, 0, 15},
		[]float32{5, 5, 5, 0, 2, 0},
		[]float32{5, 5, 0, 5, 2, 0},
		[]float32{5, 5, 5, 5, 2, 0},
		[]float32{5, 5, 5, 5, 2, 15},
		[]float32{5, 5, 5, 5, 2, 15},
		[]float32{5, 5, 5, 5, 2, 15},
	}

	for i := 0; i < len(shifts); i++ {
		//To test the backend
		if i == len(shifts)-3 {
			engo.Backend = "Mobile"
		} else if i == len(shifts)-2 {
			engo.Backend = "Web"
		} else {
			engo.Backend = "GLFW"
		}

		expX, expY := cameraShift(shifts[i][0], shifts[i][1], shifts[i][2], shifts[i][3], shifts[i][4], shifts[i][5])

		sys.camera.Update(updateTime)
		sys.Update(updateTime)

		if expX != sys.mouseX {
			t.Errorf("mouse X did not match expected X, want %v, have %v, test case %d", expX, sys.mouseX, i)
		}

		if expY != sys.mouseY {
			t.Errorf("mouse Y did not match expected Y, want %v, have %v, test case %d", expY, sys.mouseY, i)
		}

		resetCamera()
	}
}

func cameraShift(origX, origY, transX, transY, transZ, rotation float32) (expectedX, expectedY float32) {
	engo.Input.Mouse.X = origX
	engo.Input.Mouse.Y = origY
	expectedX = origX
	expectedY = origY

	if transZ != 0 {
		engo.Mailbox.Dispatch(common.CameraMessage{
			Axis:  common.ZAxis,
			Value: transZ,
		})
		expectedX = expectedX*sys.camera.Z() - (engo.GameWidth() * (sys.camera.Z() - 1) / 2)
		expectedY = expectedY*sys.camera.Z() - (engo.GameHeight() * (sys.camera.Z() - 1) / 2)
	}

	if transX != 0 {
		engo.Mailbox.Dispatch(common.CameraMessage{
			Axis:        common.XAxis,
			Value:       transX,
			Incremental: true,
		})
		expectedX += transX
	}

	if transY != 0 {
		engo.Mailbox.Dispatch(common.CameraMessage{
			Axis:        common.YAxis,
			Value:       transY,
			Incremental: true,
		})
		expectedY += transY
	}

	if rotation != 0 {
		engo.Mailbox.Dispatch(common.CameraMessage{
			Axis:  common.Angle,
			Value: rotation,
		})
		sin, cos := math.Sincos(rotation * math.Pi / 180)
		expectedX, expectedY = expectedX*cos+expectedY*sin, expectedY*cos-expectedX*sin
	}

	return
}

func resetCamera() {
	engo.Mailbox.Dispatch(common.CameraMessage{
		Axis:  common.XAxis,
		Value: common.CameraBounds.Max.X / 2,
	})
	engo.Mailbox.Dispatch(common.CameraMessage{
		Axis:  common.YAxis,
		Value: common.CameraBounds.Max.Y / 2,
	})
	engo.Mailbox.Dispatch(common.CameraMessage{
		Axis:  common.ZAxis,
		Value: 1,
	})
	engo.Mailbox.Dispatch(common.CameraMessage{
		Axis:  common.Angle,
		Value: 0,
	})
}

//Tracking should update the mouseEntity whenever the mouse moves
func TestMouseSystemTracking(t *testing.T) {
	updateTime := float32(1.0 / 60.0)
	engo.Run(engo.RunOptions{
		Width:        100,
		Height:       100,
		NoRun:        true,
		HeadlessMode: true,
	}, &MouseTestScene{3})

	//add tracking to two entities
	sys.entities[0].Track = true
	sys.entities[2].Track = true

	moves := []float32{
		15, 15,
		5, 5,
		25, 5,
		45, 5,
	}

	for i := 0; i < len(moves); i += 2 {
		//move the mouse
		engo.Input.Mouse.X = moves[i]
		engo.Input.Mouse.Y = moves[i+1]

		//update
		sys.Update(updateTime)

		//check X of each entity
		for _, e := range sys.entities {
			if e.Track {
				if e.MouseX != engo.Input.Mouse.X {
					t.Errorf("Tracking mouse X doesn't match input, tracking: %v, mouse: %v", e.MouseX, engo.Input.Mouse.X)
				}
				if e.MouseY != engo.Input.Mouse.Y {
					t.Errorf("Tracking mouse Y doesn't match input, tracking: %v, mouse: %v", e.MouseY, engo.Input.Mouse.Y)
				}
			} else {
				if e.MouseComponent.Hovered {
					if e.MouseX != engo.Input.Mouse.X {
						t.Errorf("While hovering, mouse should be tracked. for X, tracking: %v, mouse: %v", e.MouseX, engo.Input.Mouse.X)
					}
					if e.MouseY != engo.Input.Mouse.Y {
						t.Errorf("While hovering, mouse should be tracked. for Y, tracking: %v, mouse: %v", e.MouseY, engo.Input.Mouse.Y)
					}
				} else {
					if e.MouseX == engo.Input.Mouse.X {
						t.Errorf("Mouse should not be tracked, for X, tracking: %v, mouse: %v", e.MouseX, engo.Input.Mouse.X)
					}
					if e.MouseY == engo.Input.Mouse.Y {
						t.Errorf("Mouse should not be tracked, for Y, tracking: %v, mouse: %v", e.MouseX, engo.Input.Mouse.Y)
					}
				}
			}
		}
	}
}

// If there's no SpaceComponent, the mouse should not updated
func TestMouseSystemSpaceAndBoxComponentNil(t *testing.T) {
	updateTime := float32(1.0 / 60.0)
	engo.Run(engo.RunOptions{
		Width:        100,
		Height:       100,
		NoRun:        true,
		HeadlessMode: true,
	}, &MouseTestScene{3})

	// set components to nil
	sys.entities[0].SpaceComponent = nil
	sys.entities[2].SpaceComponent = nil
	sys.entities[1].Box2dComponent = nil
	sys.entities[2].Box2dComponent = nil

	// hover over them and check for component being updated
	for i, e := range sys.entities {
		// input setup
		engo.Input.Mouse.X = float32(i*20 + 5)
		engo.Input.Mouse.Y = float32(i*20 + 5)

		// update
		sys.Update(updateTime)

		// Check
		if e.Hovered {
			t.Errorf("updated to hovering even though there is no required component, entity: %d", i)
		}
	}
}

// testDrawable implements the common.Drawable interface without actually using
// the gl context or gpu for testing headless
type testDrawable struct {
	width, height float32
}

func (d *testDrawable) Close() {}

func (d *testDrawable) Width() float32 { return d.width }

func (d *testDrawable) Height() float32 { return d.height }

func (d *testDrawable) Texture() *gl.Texture { return nil }

func (d *testDrawable) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }

func TestMouseSystemRenderComponent(t *testing.T) {
	updateTime := float32(1.0 / 60.0)
	engo.Run(engo.RunOptions{
		Width:        100,
		Height:       100,
		NoRun:        true,
		HeadlessMode: true,
	}, &MouseTestScene{3})

	drawable := &testDrawable{width: 10, height: 10}

	sys.entities[0].RenderComponent = &common.RenderComponent{
		Drawable: drawable,
	}

	engo.Input.Mouse.X = 5
	engo.Input.Mouse.Y = 5

	sys.Update(updateTime)

	if sys.entities[0].Hovered != true {
		t.Error("render component test, entity 0 should be hovered but was not")
	}

	//Hidden
	sys.entities[1].RenderComponent = &common.RenderComponent{
		Drawable: drawable,
		Hidden:   true,
	}

	engo.Input.Mouse.X = 25
	engo.Input.Mouse.Y = 5

	sys.Update(updateTime)

	if sys.entities[1].Hovered == true {
		t.Error("render component test, entity 1 should not be hovered but was")
	}

	//hud
	sys.entities[2].RenderComponent = &common.RenderComponent{
		Drawable: drawable,
	}
	sys.entities[2].MouseComponent.IsHUDShader = true

	engo.Input.Mouse.X = 45
	engo.Input.Mouse.Y = 5

	engo.Mailbox.Dispatch(common.CameraMessage{
		Axis:        common.XAxis,
		Value:       25,
		Incremental: true,
	})

	sys.Update(updateTime)

	if sys.entities[2].Hovered != true {
		t.Error("render component test, entity 2 should be hovered but was not")
	}
}

// test mouse button presses
func TestMouseSystemButtonPresses(t *testing.T) {
	updateTime := float32(1.0 / 60.0)
	engo.Run(engo.RunOptions{
		Width:        100,
		Height:       100,
		NoRun:        true,
		HeadlessMode: true,
	}, &MouseTestScene{1})

	//Place cursor outside entity
	engo.Input.Mouse.X = 15
	engo.Input.Mouse.Y = 15

	//Left click
	engo.Input.Mouse.Button = engo.MouseButtonLeft

	sys.Update(updateTime)

	//Check that it didn't click it
	if sys.entities[0].Clicked {
		t.Error("Entity was left clicked even though the cursor was not over it")
	}

	//Click on the inside of the entity
	engo.Input.Mouse.X = 5
	engo.Input.Mouse.Y = 5

	engo.Input.Mouse.Button = engo.MouseButtonLeft
	engo.Input.Mouse.Action = engo.Press

	sys.Update(updateTime)

	if !sys.entities[0].Clicked {
		t.Error("Entity was not left clicked with cursor over it")
	}

	if !sys.entities[0].startedDragging {
		t.Error("Entity was not started dragging when left clicked with cursor over it")
	}

	if !sys.mouseDown {
		t.Error("System did not change to mousedown when mouse left clicked")
	}

	//Move mouse
	engo.Input.Mouse.X = 8
	engo.Input.Mouse.Y = 8

	engo.Input.Mouse.Action = engo.Move

	sys.Update(updateTime)

	if !sys.entities[0].Dragged {
		t.Error("Entity was not left dragged")
	}

	//Release
	engo.Input.Mouse.Action = engo.Release

	sys.Update(updateTime)

	if !sys.entities[0].Released {
		t.Error("Entity was not left released when mouse was released")
	}

	if sys.entities[0].Dragged {
		t.Error("Entity was still dragged when mouse was released")
	}

	if sys.entities[0].startedDragging {
		t.Error("Entity was still startedDragging when mouse was released")
	}

	if sys.mouseDown {
		t.Error("System was still mouseDown when mouse was released")
	}

	//Place cursor outside entity
	engo.Input.Mouse.X = 15
	engo.Input.Mouse.Y = 15

	//Right click
	engo.Input.Mouse.Button = engo.MouseButtonRight

	sys.Update(updateTime)

	//Check that it didn't click it
	if sys.entities[0].RightClicked {
		t.Error("Entity was right clicked even though the cursor was not over it")
	}

	//Click on the inside of the entity
	engo.Input.Mouse.X = 5
	engo.Input.Mouse.Y = 5

	engo.Input.Mouse.Button = engo.MouseButtonRight
	engo.Input.Mouse.Action = engo.Press

	sys.Update(updateTime)

	if !sys.entities[0].RightClicked {
		t.Error("Entity was not right clicked with cursor over it")
	}

	if !sys.entities[0].rightStartedDragging {
		t.Error("Entity was not right started dragging when right clicked with cursor over it")
	}

	if !sys.rightMouseDown {
		t.Error("System did not change to rightMousedown when mouse right clicked")
	}

	//Move mouse
	engo.Input.Mouse.X = 8
	engo.Input.Mouse.Y = 8

	engo.Input.Mouse.Action = engo.Move

	sys.Update(updateTime)

	if !sys.entities[0].RightDragged {
		t.Error("Entity was not right dragged")
	}

	//Release
	engo.Input.Mouse.Action = engo.Release

	sys.Update(updateTime)

	if !sys.entities[0].RightReleased {
		t.Error("Entity was not right released when mouse was released")
	}

	if sys.entities[0].RightDragged {
		t.Error("Entity was still right dragged when mouse was released")
	}

	if sys.entities[0].rightStartedDragging {
		t.Error("Entity was still rightStartedDragging when mouse was released")
	}

	if sys.rightMouseDown {
		t.Error("System was still rightMouseDown when mouse was released")
	}
}
