package raymarch

//Makes a trapezoidal house with dir1 and dir2 as the basis vectors
//dir2 ^------------------/
//    /                  /
//   /                  /
//  /                  /
// /-------------------> dir1
func MakeHouseShell(pos Vector3, dir1 Vector3, dir2 Vector3, wallThickness float64, houseHeight float64) Geometry{
	var wall1 = MakeParalelipiped(pos, dir1, dir2.Normalised().Mul(wallThickness), Vector3{0, 0, houseHeight})
	var wall2 = MakeParalelipiped(pos, dir2, dir1.Normalised().Mul(wallThickness), Vector3{0, 0, houseHeight})
	var wall3 = MakeParalelipiped(Add(pos, dir1), dir2, dir1.Normalised().Mul(wallThickness), Vector3{0, 0, houseHeight})
	var wall4 = MakeParalelipiped(Add(pos, dir2), dir1, dir2.Normalised().Mul(wallThickness), Vector3{0, 0, houseHeight})
	var roof = MakeParalelipiped(Add(pos, Vector3{0, 0, houseHeight}), dir1, dir2, Vector3{0, 0, wallThickness})
	return Union(wall1, wall2, wall3, wall4, roof)
}

func MakeDoorFrame(pos Vector3, dir Vector3, width float64, height float64, wallThickness float64) Geometry{
	var up = Vector3{0, 0, height}
	var right = Cross(dir, up).Normalised()
	var halfInternalWidth = width / 2 - wallThickness / 2
	var wall1 = MakeParalelipiped(Add(pos, right.Mul(halfInternalWidth)), right.Mul(wallThickness), up, dir.Mul(wallThickness))
	var wall2 = MakeParalelipiped(Add(pos, right.Mul(-halfInternalWidth)), right.Mul(-wallThickness), up, dir.Mul(wallThickness))
	var roof = MakeParalelipiped(Add(pos, up, right.Mul(width/2)), right.Mul(-width), up.Normalised().Mul(wallThickness), dir.Mul(wallThickness))
	return Union(wall1, wall2, roof)
}

func MakeHouse(){

}

