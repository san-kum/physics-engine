package physics

import "math"

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Add(other Vector) Vector {
	return Vector{v.X + other.X, v.Y + other.Y}
}

func (v Vector) Subtract(other Vector) Vector {
	return Vector{v.X - other.X, v.Y - other.Y}
}

func (v Vector) Scale(scalar float64) Vector {
	return Vector{v.X * scalar, v.Y * scalar}
}

func (v Vector) Dot(other Vector) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

type Object struct {
	Position     Vector
	LastPosition Vector
	Acceleration Vector
	Mass         float64
	Radius       float64
}

// NewObject creates a new object with the given initial conditions.
func NewObject(position, velocity Vector, mass, radius float64, timeStep float64) *Object {
	//  NOTE: Verlet: lastPosition is based on the initial velocity
	lastPosition := position.Subtract(velocity.Scale(timeStep))

	return &Object{
		Position:     position,
		LastPosition: lastPosition,
		Acceleration: Vector{0, 0},
		Mass:         mass,
		Radius:       radius,
	}
}

func (o *Object) collidesWith(other *Object) bool {
	distance := o.Position.Subtract(other.Position).Magnitude()
	return distance < (o.Radius + other.Radius)
}
