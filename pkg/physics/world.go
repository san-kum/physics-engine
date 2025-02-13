package physics

// World represents the simulation environment
type World struct {
	Objects  []*Object
	TimeStep float64
	Gravity  Vector
}

func NewWorld(timeStep float64, gravity Vector) *World {
	return &World{
		Objects:  make([]*Object, 0),
		TimeStep: timeStep,
		Gravity:  gravity,
	}
}

func (w *World) AddObject(object *Object) {
	w.Objects = append(w.Objects, object)
}

func (w *World) Update() {
	w.CheckCollisions()
	for _, obj := range w.Objects {
		w.applyGravity(obj)
		w.verletIntegration(obj)
	}
}

func (w *World) applyGravity(o *Object) {
	o.Acceleration = o.Acceleration.Add(w.Gravity)
}

// verletIntegration
func (w *World) verletIntegration(o *Object) {
	nextPosition := o.Position.Scale(2).Subtract(o.LastPosition).Add(o.Acceleration.Scale(w.TimeStep * w.TimeStep))
	o.LastPosition = o.Position
	o.Position = nextPosition
	o.Acceleration = Vector{0, 0}

}

func (w *World) CheckCollisions() {
	for i := 0; i < len(w.Objects); i++ {
		for j := i + 1; j < len(w.Objects); j++ {
			if w.Objects[i].collidesWith(w.Objects[j]) {
				w.resolveCollision(w.Objects[i], w.Objects[j])
			}
		}
	}
}

func (w *World) resolveCollision(o1, o2 *Object) {
	collisionNormal := o2.Position.Subtract(o1.Position)
	distance := collisionNormal.Magnitude()

	overlap := (o1.Radius + o2.Radius) - distance

	if distance == 0 {
		return
	}
	collisionNormal = collisionNormal.Scale(1.0 / distance)

	totalMass := o1.Mass + o2.Mass
	o1Movement := collisionNormal.Scale(-overlap * (o2.Mass / totalMass))
	o2Movement := collisionNormal.Scale(overlap * (o1.Mass / totalMass))

	o1.Position = o1.Position.Add(o1Movement)
	o2.Position = o2.Position.Add(o2Movement)
}
