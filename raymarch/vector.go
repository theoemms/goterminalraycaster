package raymarch
import "math"

type Vector3 struct {
    X, Y, Z float64
}

func Add(vecs ...Vector3) Vector3{
	var result Vector3
	for _, vec := range vecs{
		result.X += vec.X
		result.Y += vec.Y
		result.Z += vec.Z
	}
	return result
}

func (self Vector3) Mul(scalar float64) Vector3 {
	return Vector3{self.X * scalar, self.Y * scalar, self.Z * scalar}
}

func Div(lhs Vector3, rhs float64) Vector3{
	return lhs.Mul(1 / rhs)
}

func Sub(lhs Vector3, rhs Vector3) Vector3{
	return Add(lhs, rhs.Mul(-1))
}

func Dot(lhs Vector3, rhs Vector3) float64{
	return lhs.X * rhs.X + lhs.Y * rhs.Y + lhs.Z * rhs.Z
}

func Cross(lhs Vector3, rhs Vector3) Vector3{
	return Vector3{lhs.Y * rhs.Z - lhs.Z * rhs.Y, lhs.Z * rhs.X - lhs.X * rhs.Z, lhs.X * rhs.Y - lhs.Y * rhs.X}
}

func (self Vector3) LengthSq() float64{
	return self.X * self.X + self.Y * self.Y + self.Z * self.Z
}

func (self Vector3) Length() float64{
	return math.Sqrt(self.LengthSq())
}

func (self Vector3) Normalised() Vector3{
	return Div(self, self.Length())
}