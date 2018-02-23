package engoBox2dSystem

import (
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"engo.io/engo/math"

	"github.com/ByteArena/box2d"
)

// MouseSystemPriority ensures the mouse system is updated before any other systems
const MouseSystemPriority = 100

// MouseComponent is the location for the MouseSystem to store its results;
// to be used / viewed by other Systems
type MouseComponent struct {
	// Clicked is true whenever the Mouse was clicked over
	// the entity space in this frame
	Clicked bool
	// Released is true whenever the left mouse button is released over the
	// entity space in this frame
	Released bool
	// Hovered is true whenever the Mouse is hovering
	// the entity space in this frame. This does not necessarily imply that
	// the mouse button was pressed down in your entity space.
	Hovered bool
	// Dragged is true whenever the entity space was left-clicked,
	// and then the mouse started moving (while holding)
	Dragged bool
	// RightClicked is true whenever the entity space was right-clicked
	// in this frame
	RightClicked bool
	// RightDragged is true whenever the entity space was right-clicked,
	// and then the mouse started moving (while holding)
	RightDragged bool
	// RightReleased is true whenever the right mouse button is released over
	// the entity space in this frame. This does not necessarily imply that
	// the mouse button was pressed down in your entity space.
	RightReleased bool
	// Enter is true whenever the Mouse entered the entity space in that frame,
	// but wasn't in that space during the previous frame
	Enter bool
	// Leave is true whenever the Mouse was in the space on the previous frame,
	// but now isn't
	Leave bool
	// Position of the mouse at any moment this is generally used
	// in conjunction with Track = true
	MouseX float32
	MouseY float32
	// Set manually this to true and your mouse component will track the mouse
	// and your entity will always be able to receive an updated mouse
	// component even if its space is not under the mouse cursor
	// WARNING: you MUST know why you want to use this because it will
	// have serious performance impacts if you have many entities with
	// a MouseComponent in tracking mode.
	// This is ideally used for a really small number of entities
	// that must really be aware of the mouse details event when the
	// mouse is not hovering them
	Track bool
	// Modifier is used to store the eventual modifiers that were pressed during
	// the same time the different click events occurred
	Modifier engo.Modifier

	// startedDragging is used internally to see if *this* is the object that is being dragged
	startedDragging bool
	// startedRightDragging is used internally to see if *this* is the object that is being right-dragged
	rightStartedDragging bool

	// IsHUDShader is used to update the mouse component properly for the common.HUDShader
	IsHUDShader bool
}

type mouseEntity struct {
	*ecs.BasicEntity
	*MouseComponent
	*common.SpaceComponent
	*common.RenderComponent
	*Box2dComponent
}

// MouseSystem listens for mouse events and changes value for MouseComponent accordingly
type MouseSystem struct {
	entities []mouseEntity
	world    *ecs.World
	camera   *common.CameraSystem

	mouseX         float32
	mouseY         float32
	mouseDown      bool
	rightMouseDown bool
}

// Priority implements prioritizer interface
func (m *MouseSystem) Priority() int { return MouseSystemPriority }

// New adds world and camera to the MouseSystem
func (m *MouseSystem) New(w *ecs.World) {
	m.world = w

	// First check to see if the CameraSystem is available
	for _, system := range m.world.Systems() {
		switch sys := system.(type) {
		case *common.CameraSystem:
			m.camera = sys
		}
	}

	if m.camera == nil {
		log.Println("ERROR: CameraSystem not found - have you added the `RenderSystem` before the `MouseSystem`?")
		return
	}
}

// Add adds a new entity to the MouseSystem
func (m *MouseSystem) Add(basic *ecs.BasicEntity, mouse *MouseComponent, space *common.SpaceComponent, render *common.RenderComponent, box *Box2dComponent) {
	m.entities = append(m.entities, mouseEntity{basic, mouse, space, render, box})
}

// AddByInterface adds the entity that implements the Mouseable interface to the MouseSystem
func (m *MouseSystem) AddByInterface(o Mouseable) {
	m.Add(o.GetBasicEntity(), o.GetMouseComponent(), o.GetSpaceComponent(), o.GetRenderComponent(), o.GetBox2dComponent())
}

// Remove removes an entity from the MouseSystem
func (m *MouseSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, entity := range m.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		m.entities = append(m.entities[:delete], m.entities[delete+1:]...)
	}
}

// Update updates the MouseComponent based on location of cursor and state of the mouse buttons
func (m *MouseSystem) Update(dt float32) {
	// Translate Mouse.X and Mouse.Y into "game coordinates"
	switch engo.Backend {
	case "GLFW":
		m.mouseX = engo.Input.Mouse.X*m.camera.Z()*(engo.GameWidth()/engo.CanvasWidth()) + m.camera.X() - (engo.GameWidth()/2)*m.camera.Z()
		m.mouseY = engo.Input.Mouse.Y*m.camera.Z()*(engo.GameHeight()/engo.CanvasHeight()) + m.camera.Y() - (engo.GameHeight()/2)*m.camera.Z()
	case "Mobile":
		m.mouseX = engo.Input.Mouse.X*m.camera.Z() + m.camera.X() - (engo.GameWidth()/2)*m.camera.Z() + (engo.ResizeXOffset / 2)
		m.mouseY = engo.Input.Mouse.Y*m.camera.Z() + m.camera.Y() - (engo.GameHeight()/2)*m.camera.Z() + (engo.ResizeYOffset / 2)
	case "Web":
		m.mouseX = engo.Input.Mouse.X*m.camera.Z() + m.camera.X() - (engo.GameWidth()/2)*m.camera.Z() + (engo.ResizeXOffset / 2)
		m.mouseY = engo.Input.Mouse.Y*m.camera.Z() + m.camera.Y() - (engo.GameHeight()/2)*m.camera.Z() + (engo.ResizeYOffset / 2)
	}

	// Rotate if needed
	if m.camera.Angle() != 0 {
		sin, cos := math.Sincos(m.camera.Angle() * math.Pi / 180)
		m.mouseX, m.mouseY = m.mouseX*cos+m.mouseY*sin, m.mouseY*cos-m.mouseX*sin
	}

	for _, e := range m.entities {
		// Reset all values except these
		*e.MouseComponent = MouseComponent{
			Track:                e.MouseComponent.Track,
			Hovered:              e.MouseComponent.Hovered,
			startedDragging:      e.MouseComponent.startedDragging,
			rightStartedDragging: e.MouseComponent.rightStartedDragging,
			IsHUDShader:          e.MouseComponent.IsHUDShader,
		}

		if e.MouseComponent.Track {
			e.MouseComponent.MouseX = m.mouseX
			e.MouseComponent.MouseY = m.mouseY
		}

		mx := m.mouseX
		my := m.mouseY

		if e.SpaceComponent == nil || e.Box2dComponent == nil {
			continue
		}

		//set box2d body to SpaceComponent's position and rotation
		e.Body.SetTransform(Conv.ToBox2d2Vec(e.Center()), Conv.DegToRad(e.Rotation))

		if e.RenderComponent != nil {
			// Hardcoded special case for the HUD | TODO: make generic instead of hardcoding
			if e.MouseComponent.IsHUDShader {
				mx = engo.Input.Mouse.X
				my = engo.Input.Mouse.Y
			}

			if e.RenderComponent.Hidden {
				continue // skip hidden components
			}
		}

		mousePoint := box2d.B2Vec2{
			X: Conv.PxToMeters(mx),
			Y: Conv.PxToMeters(my),
		}
		var containsMouse bool

		for f := e.Body.GetFixtureList(); f != nil; f = f.GetNext() {
			if f.TestPoint(mousePoint) {
				containsMouse = true
				break
			}
		}

		// If the Mouse component is a tracker we always update it
		// Check if the X-value is within range
		// and if the Y-value is within range

		if e.MouseComponent.Track || e.MouseComponent.startedDragging || containsMouse {

			e.MouseComponent.Enter = !e.MouseComponent.Hovered
			e.MouseComponent.Hovered = true
			e.MouseComponent.Released = false

			if !e.MouseComponent.Track {
				// If we're tracking, we've already set these
				e.MouseComponent.MouseX = mx
				e.MouseComponent.MouseY = my
			}

			switch engo.Input.Mouse.Action {
			case engo.Press:
				switch engo.Input.Mouse.Button {
				case engo.MouseButtonLeft:
					e.MouseComponent.Clicked = true
					e.MouseComponent.startedDragging = true
					m.mouseDown = true
				case engo.MouseButtonRight:
					e.MouseComponent.RightClicked = true
					e.MouseComponent.rightStartedDragging = true
					m.rightMouseDown = true
				}
			case engo.Release:
				switch engo.Input.Mouse.Button {
				case engo.MouseButtonLeft:
					e.MouseComponent.Released = true
				case engo.MouseButtonRight:
					e.MouseComponent.RightReleased = true
				}
			case engo.Move:
				if m.mouseDown && e.MouseComponent.startedDragging {
					e.MouseComponent.Dragged = true
				}
				if m.rightMouseDown && e.MouseComponent.rightStartedDragging {
					e.MouseComponent.RightDragged = true
				}
			}
		} else {
			if e.MouseComponent.Hovered {
				e.MouseComponent.Leave = true
			}

			e.MouseComponent.Hovered = false
		}

		if engo.Input.Mouse.Action == engo.Release {
			switch engo.Input.Mouse.Button {
			case engo.MouseButtonLeft:
				e.MouseComponent.Dragged = false
				e.MouseComponent.startedDragging = false
				m.mouseDown = false
			case engo.MouseButtonRight:
				e.MouseComponent.RightDragged = false
				e.MouseComponent.rightStartedDragging = false
				m.rightMouseDown = false
			}
		}

		// propagate the modifiers to the mouse component so that game
		// implementers can take different decisions based on those
		e.MouseComponent.Modifier = engo.Input.Mouse.Modifer
	}

	//Remove all bodies on list for removal
	removeBodies()
}
