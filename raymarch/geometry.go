package raymarch
import "math"

type Geometry interface {
    SurfaceDistance(pos Vector3) float64
    Normal(pos Vector3) Vector3
}

type UnionGeometry struct {
	A, B Geometry
}

func (self UnionGeometry) SurfaceDistance(pos Vector3) float64 {
	return math.Min(self.A.SurfaceDistance(pos), self.B.SurfaceDistance(pos))
}

func (self UnionGeometry) Normal(pos Vector3) Vector3 {
	if self.A.SurfaceDistance(pos) < self.B.SurfaceDistance(pos) {
		return self.A.Normal(pos)
	}
	return self.B.Normal(pos)
}

func Union(lhs Geometry, rhs Geometry) Geometry {
	return UnionGeometry{lhs, rhs}
}

type Sphere struct {
	Pos Vector3
	Radius float64
}

func (self Sphere) SurfaceDistance(pos Vector3) float64 {
	return Sub(self.Pos, pos).Length() - self.Radius
}

func (self Sphere) Normal(pos Vector3) Vector3 {
	return Sub(pos, self.Pos).Normalised()
}

type Plane struct {
	Pos, Normal, Tangent Vector3
	Width, Length float64
}


