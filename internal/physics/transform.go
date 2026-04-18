package physics

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

type Quat struct {
	X float64
	Y float64
	Z float64
	W float64
}

type Transform struct {
	Position Vec3
	Rotation Quat
	Scale    Vec3
}

func IdentityQuat() Quat {
	return Quat{W: 1}
}

func UnitVec3() Vec3 {
	return Vec3{X: 1, Y: 1, Z: 1}
}

func DefaultTransform(position Vec3) Transform {
	return Transform{
		Position: position,
		Rotation: IdentityQuat(),
		Scale:    UnitVec3(),
	}
}
