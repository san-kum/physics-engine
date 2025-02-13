package physics

import (
	"math"
	"testing"
)

func TestVerletIntegration(t *testing.T) {
	timeStep := 0.1
	world := NewWorld(timeStep, Vector{0, 0})
	initialPosition := Vector{1, 2}
	initialVelocity := Vector{0.5, 0.3}
	obj := NewObject(initialPosition, initialVelocity, 1.0, 1.0, timeStep)
	world.AddObject(obj)

	world.Update()

	expectedPosition := initialPosition.Scale(2).Subtract(initialPosition.Subtract(initialVelocity.Scale(timeStep))).Add(Vector{0, 0}.Scale(world.TimeStep * world.TimeStep))

	if !approxEqual(obj.Position, expectedPosition, 0.0001) {
		t.Errorf("Expected position: %v, got %v", expectedPosition, obj.Position)
	}
	// last position
	if !approxEqual(obj.LastPosition, initialPosition, 0.0001) {
		t.Errorf("Expected last position: %v, got %v", initialPosition, obj.LastPosition)
	}
	world.Update()

	expectedPosition2 := expectedPosition.Scale(2).Subtract(initialPosition).Add(Vector{0, 0}.Scale(world.TimeStep * world.TimeStep))
	if !approxEqual(obj.Position, expectedPosition2, 0.0001) {
		t.Errorf("Expected positions: %v, got %v", expectedPosition2, obj.Position)

	}

	if !approxEqual(obj.LastPosition, expectedPosition, 0.0001) {
		t.Errorf("Expected last position: %v, got %v", expectedPosition, obj.LastPosition)
	}
}

func approxEqual(v1, v2 Vector, tolerance float64) bool {
	return math.Abs(v1.X-v2.X) < tolerance && math.Abs(v1.Y-v2.Y) < tolerance
}

func TestGravity(t *testing.T) {
	world := NewWorld(0.1, Vector{0, -9.81})
	obj := NewObject(Vector{0, 0}, Vector{0, 0}, 1.0, 1.0, 0.1)
	world.AddObject(obj)

	world.Update()
	expectedAcceleration := Vector{0, 0}
	if !approxEqual(obj.Acceleration, Vector{0, 0}, 0.0001) {
		t.Errorf("Expected acceleration: %v, got %v", expectedAcceleration, obj.Acceleration)
	}
}

func TestCollisionDetection(t *testing.T) {
	world := NewWorld(0.1, Vector{0, 0})
	obj1 := NewObject(Vector{0, 0}, Vector{0, 0}, 1.0, 1.0, 0.1)
	obj2 := NewObject(Vector{1.5, 0}, Vector{0, 0}, 1.0, 1.0, 0.1)
	world.AddObject(obj1)
	world.AddObject(obj2)

	if collided := obj1.collidesWith(obj2); !collided {
		t.Errorf("objects should be colliding")
	}
	world.CheckCollisions()

	if collided := obj1.collidesWith(obj2); collided {
		t.Errorf("objects should not be colliding")
	}
}

func TestNoCollision(t *testing.T) {
	world := NewWorld(0.1, Vector{0, 0})
	obj1 := NewObject(Vector{0, 0}, Vector{0, 0}, 1.0, 1.0, 0.1)
	obj2 := NewObject(Vector{3, 0}, Vector{0, 0}, 1.0, 1.0, 0.1)
	world.AddObject(obj1)
	world.AddObject(obj2)

	if collided := obj1.collidesWith(obj2); collided {
		t.Errorf("Objects should not be colliding")
	}
}
