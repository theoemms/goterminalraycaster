package raymarch

type Vector3 struct {
    X, Y, Z float64
}

type Camera struct {
    Position, Heading, Up Vector3
    FOV float64
}

type Geometry interface {
    SurfaceDistance(pos Vector3)
    Normal(pos Vector3)
}

type Scene struct{
    Cam Camera
    Geom []*Geometry
}

type Ray struct{
    Pos, Dir Vector3
}

type Raymarcher struct {
    StepSize, FarDist, ConvergeDist float64

}

type RayHit{
    Pos Vector3
    Geom Geometry
}

func (self Raymarcher) CastRay (ray Ray, scene Scene) RayHit {
    var geometry := scene.Geom
    var pos := ray.Pos
    for {
        var minDist = self.FarDist
        for i, g in geometry {
            var surfaceDist := g.SurfaceDistance(pos)
            if surfaceDist < self.ConvergeDist{
                return RayHit{pos, g}
            }
            if surfaceDist < minDist{
                minDist := surfaceDist
            }
        }
        if minDist == self.FarDist{
            return null
        }
        pos = pos + self.StepSize * ray.Dir * minDist
    }
}

