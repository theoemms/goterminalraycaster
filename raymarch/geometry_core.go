package raymarch

type Geometry interface {
    SurfaceDistance(pos Vector3) float64
    Normal(pos Vector3) Vector3
}

type CompositionGeometry struct{
	A, B Geometry
	CompositionFunc func(float64, float64) bool
}

func (self CompositionGeometry) SurfaceDistance(pos Vector3) float64 {
	if self.CompositionFunc(self.A.SurfaceDistance(pos), self.B.SurfaceDistance(pos)) {
		return self.A.SurfaceDistance(pos)
	}
	return self.B.SurfaceDistance(pos)
}

func (self CompositionGeometry) Normal(pos Vector3) Vector3 {
	if self.A.SurfaceDistance(pos) < self.B.SurfaceDistance(pos){
		return self.A.Normal(pos)
	}
	return self.B.Normal(pos)
}

func Union(geoms... Geometry) Geometry {
	var unionSurfaceDistanceComposition = func(dist1 float64, dist2 float64) bool{
		return dist1 < dist2
	}

	if len(geoms) <= 1{
		return geoms[0]
	}

	return CompositionGeometry{geoms[0], Union(geoms[1:]...), unionSurfaceDistanceComposition}
}

func Intersection(geoms... Geometry) Geometry {
	var intersectionSurfaceDistanceComposition = func(dist1 float64, dist2 float64) bool{
		return dist1 > dist2
	}

	if len(geoms) <= 1{
		return geoms[0]
	}

	return CompositionGeometry{geoms[0], Intersection(geoms[1:]...), intersectionSurfaceDistanceComposition}
}

func Subtract(A Geometry, B Geometry) Geometry{
	var subtractionSurfaceDistanceComposition = func(dist1 float64, dist2 float64) bool{
		return dist1 > -dist2
	}
	return CompositionGeometry{A, B, subtractionSurfaceDistanceComposition}
}

//Primitives
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

func MakeParalelipiped(pos Vector3, v1 Vector3, v2 Vector3, v3 Vector3) Geometry{
	var planes = [](Geometry){}
	findPlaneNormal := func(tang1 Vector3, tang2 Vector3, normalDir Vector3) Vector3{
		var normal = Cross(tang1, tang2)
		if Dot(normal, normalDir) < 0{
			normal = normal.Mul(-1)
		}
		return normal.Normalised()
	}

	planes = append(planes, Plane{pos, findPlaneNormal(v2, v3, v1.Mul(-1))})
	planes = append(planes, Plane{pos, findPlaneNormal(v1, v3, v2.Mul(-1))})
	planes = append(planes, Plane{pos, findPlaneNormal(v1, v2, v3.Mul(-1))})
	planes = append(planes, Plane{Add(pos, v1), findPlaneNormal(v2, v3, v1)})
	planes = append(planes, Plane{Add(pos, v2), findPlaneNormal(v1, v3, v2)})
	planes = append(planes, Plane{Add(pos, v3), findPlaneNormal(v1, v2, v3)})

	return Intersection(planes...)
}
