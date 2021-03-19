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

type IntersectionGeometry struct {
	A, B Geometry
}

func (self IntersectionGeometry) SurfaceDistance(pos Vector3) float64 {
	return math.Max(self.A.SurfaceDistance(pos), self.B.SurfaceDistance(pos))
}

func (self IntersectionGeometry) Normal(pos Vector3) Vector3 {
	if self.A.SurfaceDistance(pos) > self.B.SurfaceDistance(pos) {
		return self.A.Normal(pos)
	}
	return self.B.Normal(pos)
}

func Intersection(lhs Geometry, rhs Geometry) Geometry {
	return IntersectionGeometry{lhs, rhs}
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
	Pos, FixedNormal Vector3
}

func (self Plane) SurfaceDistance(pos Vector3) float64 {
	return Dot(Sub(pos, self.Pos), self.FixedNormal)
}

func (self Plane) Normal(pos Vector3) Vector3 {
	return self.FixedNormal
}

type Paralelipiped struct{
	Geom Geometry
}

func MakeParalelipiped(pos Vector3, v1 Vector3, v2 Vector3, v3 Vector3) Paralelipiped{
	var planes = [](Plane){}
	findPlaneNormal := func(tang1 Vector3, tang2 Vector3, normalDir Vector3) Vector3{
		var normal = Cross(tang1, tang2)
		if Dot(normal, normalDir) < 0{
			return normal.Mul(-1).Normalised()
		}
		return normal.Normalised()
	}

	planes = append(planes, Plane{pos, findPlaneNormal(v2, v3, v1.Mul(-1))})
	planes = append(planes, Plane{pos, findPlaneNormal(v1, v3, v2.Mul(-1))})
	planes = append(planes, Plane{pos, findPlaneNormal(v1, v2, v3.Mul(-1))})
	planes = append(planes, Plane{Add(pos, v1), findPlaneNormal(v2, v3, v1)})
	planes = append(planes, Plane{Add(pos, v2), findPlaneNormal(v1, v3, v2)})
	planes = append(planes, Plane{Add(pos, v3), findPlaneNormal(v1, v2, v3)})

	var Geom Geometry = planes[0]
	for i := 1; i < len(planes); i++{
		Geom = IntersectionGeometry{Geom, planes[i]}
	}
	return Paralelipiped{Geom}
}

func (self Paralelipiped) SurfaceDistance(pos Vector3) float64 {
	return self.Geom.SurfaceDistance(pos)
}

func (self Paralelipiped) Normal(pos Vector3) Vector3 {
	return self.Geom.Normal(pos)
}