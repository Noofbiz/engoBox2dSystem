package engoBox2dSystem

import (
	"testing"

	"github.com/ByteArena/box2d"
)

func TestDestroyBody(t *testing.T) {
	//Create bodies and destroy them
	for i := 0; i < 5; i++ {
		bodyDef := box2d.NewB2BodyDef()
		bodyDef.Position = box2d.B2Vec2{X: float64(i * 20), Y: float64(i * 20)}
		body := World.CreateBody(bodyDef)
		comp := Box2dComponent{Body: body}
		comp.DestroyBody()
	}

	// Check that the list is correct
	if len(listOfBodiesToRemove) != 5 {
		t.Errorf("listOfBodiesToRemove has wrong count, want: %d, got: %d", 5, len(listOfBodiesToRemove))
	}

	// Clear out list
	for _, bod := range listOfBodiesToRemove {
		World.DestroyBody(bod)
	}
	listOfBodiesToRemove = make([]*box2d.B2Body, 0)

	// Check that list was cleared
	if len(listOfBodiesToRemove) != 0 {
		t.Errorf("listOfBodiesToRemove has wrong count after clearing, want %d, got %d", 0, len(listOfBodiesToRemove))
	}
}
