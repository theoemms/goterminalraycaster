package raymarch
import (
	"math"
	//"fmt"
)

//Need a maths utils file...
func Min(ints ...int) int{
	var output = ints[0]
	for _, n := range ints{
		if n < output{
			output = n
		}
	}
	return output
}

func Max(ints ...int) int{
	var output = ints[0]
	for _, n := range ints{
		if n > output{
			output = n
		}
	}
	return output
}

type Camera struct {
    Position, Heading, Up Vector3
    FOV float64
}

func (self Camera) Right() Vector3{
	return Cross(self.Heading, self.Up).Normalised()
}

func (self Camera) GetRay(screenPos Vector3) Ray{
	var width = math.Tan(self.FOV / 2)
	var rayDir = Add(
		self.Heading, 
		self.Right().Mul(screenPos.X * width),
		self.Up.Mul(screenPos.Y * width))
	return Ray{self.Position, rayDir.Normalised()}
}

type Scene struct{
    Cam *Camera
    Geom []Geometry
    Lights []Light
}

type Ray struct{
    Pos, Dir Vector3
}

type Raymarcher struct {
    StepSize, FarDist, ConvergeDist float64
}

type RayHit struct {
    Pos Vector3
    Geom Geometry
}

func (self Raymarcher) CastRay (ray Ray, scene *Scene) *RayHit {
    var geometry = scene.Geom
    var pos = ray.Pos
    for {
        var minDist = self.FarDist
        if Sub(ray.Pos, pos).Length() >= self.FarDist{
            return nil
        }

        for _, g := range geometry {
            var surfaceDist = g.SurfaceDistance(pos)
            if math.Abs(surfaceDist) < self.ConvergeDist{
                return &RayHit{pos, g}
            }

            if math.Abs(surfaceDist) < minDist{
                minDist = math.Abs(surfaceDist)
            }
        }

        pos = Add(pos, ray.Dir.Mul(self.StepSize * minDist))
    }
}

//Screenpos X and Y are screen-space centered at 0, 0 with min/max x and y going from -1 to 1
func (self Raymarcher) BlinnPhong(screenPos Vector3, scene *Scene, fogDist float64) (bool, float64) {
	var ray = scene.Cam.GetRay(screenPos)
	var rayHit = self.CastRay(ray, scene)
	if rayHit == nil{
		return false, 0
	}
	var geometry = rayHit.Geom
	var normal = geometry.Normal(rayHit.Pos)

	var intensity = 0.0
	for _, light := range scene.Lights{
		var lightDir = light.Direction(rayHit.Pos)
		if self.CastRay(Ray{Add(rayHit.Pos, lightDir.Mul(-2 * self.ConvergeDist)), lightDir.Mul(-1)}, scene) == nil{
			intensity += math.Max(0, Dot(lightDir, normal.Mul(-1)) * light.Intensity(rayHit.Pos))
		} 
	}

	var dist = Sub(rayHit.Pos, scene.Cam.Position).Length()
	return true, intensity * math.Exp(- dist / fogDist)
}
