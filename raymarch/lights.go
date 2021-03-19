package raymarch

type Light interface {
	Direction(pos Vector3) Vector3
	Intensity(pos Vector3) float64
}

type SunLight struct{
	LightDir Vector3 
	Power float64
}

func (self SunLight) Direction(pos Vector3) Vector3 {
	return self.LightDir
}

func (self SunLight) Intensity(pos Vector3) float64 {
	return self.Power
}

type PointLight struct{
	Position Vector3
	Power float64
}

func (self PointLight) Direction(pos Vector3) Vector3{
	return Sub(pos, self.Position).Normalised()
}

func (self PointLight) Intensity(pos Vector3) float64{
	return self.Power / Sub(pos, self.Position).LengthSq()
}
