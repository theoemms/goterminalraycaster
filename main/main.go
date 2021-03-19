package main
import (
	"fmt"
	"raymarch"
	"strings"
	"math"
	"time"
	"os"
	"bufio"
)

func Clear(){
	for i := 0; i < 100; i++{
		fmt.Printf("\n")
	}
}

func InitScene() (*raymarch.Camera, *raymarch.Scene){
	var sphere1 = raymarch.Sphere{raymarch.Vector3{0, 0.5, -9}, 0.5}
	var sphere2 = raymarch.Sphere{raymarch.Vector3{-3, 0.5, -7}, 0.5}
	var sphere3 = raymarch.Sphere{raymarch.Vector3{2, 0.5, -5}, 0.5}

	var plane = raymarch.Plane{raymarch.Vector3{}, raymarch.Vector3{0, 1, 0}.Normalised()}

	var wall1 = raymarch.MakeParalelipiped(raymarch.Vector3{-0.85, 0.0, -2}, raymarch.Vector3{0.2, 0, 0},  raymarch.Vector3{0, 1, 0}, raymarch.Vector3{0, 0, -6})
	var wall2 = raymarch.MakeParalelipiped(raymarch.Vector3{0.65, 0.0, -2}, raymarch.Vector3{0.2, 0, 0},  raymarch.Vector3{0, 1, 0}, raymarch.Vector3{0, 0, -6})
	var roof = raymarch.MakeParalelipiped(raymarch.Vector3{-0.85, 1.0, -2}, raymarch.Vector3{1.7, 0, 0},  raymarch.Vector3{0, 0.2, 0}, raymarch.Vector3{0, 0, -6})

	var lights = [](raymarch.Light){}
	var sunLight = raymarch.SunLight{ raymarch.Vector3{-1, -1, -1}, 0.2}

	lights = append(lights, &sunLight)
	lights = append(lights, raymarch.PointLight{raymarch.Vector3{0, 0.5, -5}, 1.0})
	for i := 0; i < 3; i++ {
		var light1 = raymarch.PointLight{ raymarch.Vector3{-1.5, 0.8, float64(-5 - 2 * i)}, 2}
		var light2 = raymarch.PointLight{ raymarch.Vector3{1.5, 0.8, float64(-5 - 2 * i)}, 2}
		lights = append(lights, &light1)
		lights = append(lights, &light2)
	}

	var camera = raymarch.Camera{raymarch.Vector3{0, 0.5, 2}, raymarch.Vector3{0, 0, -1}, raymarch.Vector3{0, 1, 0}, math.Pi / 3}
	var scene = raymarch.Scene{
		&camera, 
		[](raymarch.Geometry){ &plane, &sphere1, &sphere2, &sphere3, &wall1, &wall2, roof}, 
		lights}
   	return &camera, &scene
}

func Draw(scene *raymarch.Scene){
	var intensityChars = []string{" ", ".", "~", ":", "o", "+", "#", "@"}
	var pixAspectRatio float64 = 2
	var resX, resY = 120, 60
    var marcher = raymarch.Raymarcher{1, 50.0, 0.01}
    var screenIntensities = make([]float64, resX * resY)
    var clipIntensity float64 = 3
    for x := 0; x < resX; x++ {
    	for y := 0; y < resY; y++ {
    		var xSS = 2 * float64(x) / float64(resX) - 1
    		var ySS = pixAspectRatio * (2 * (float64(y) + float64(resX - resY) / 2) / (float64(resX)) - 1)
    		var screenPos = raymarch.Vector3{xSS, ySS, 0}
    		var intensity = math.Min(clipIntensity, math.Log(1 + marcher.BlinnPhong(screenPos, scene, 0.1)))
    		screenIntensities[x + resX * y] = intensity
    	}
    }
	
   	var screenChars = make([]string, resX * resY)
   	for i := 0; i < len(screenChars); i++ {
		var charIx = int(screenIntensities[i] * float64(len(intensityChars) - 1))
		charIx = raymarch.Max(0, raymarch.Min(len(intensityChars) - 1, charIx))
		screenChars[i] = intensityChars[charIx]
	}

	var sb strings.Builder
	for y := 0; y < resY; y++ {
		for x := 0; x < resX; x++ {
    		sb.WriteString(screenChars[x + resX * (resY - y - 1)])
    	}
    	sb.WriteString("\n")
    }

    Clear()
	fmt.Printf(sb.String())
}

func ExecuteUserCommand(command rune, camera *raymarch.Camera){
	var turnSpeed = 1.0
	var stepSize = 1.0
	switch command {
	case 'w':
		camera.Position = raymarch.Add(camera.Position, camera.Heading.Mul(stepSize))
	case 'a':
		camera.Heading = raymarch.Add(camera.Heading, camera.Right().Mul(-turnSpeed)).Normalised()
	case 's':
		camera.Position = raymarch.Add(camera.Position, camera.Heading.Mul(-stepSize))
	case 'd':
		camera.Heading = raymarch.Add(camera.Heading, camera.Right().Mul(turnSpeed)).Normalised()
	}
}

func main() {
	var camera, scene = InitScene()
    var deltaTime = 0.1
    var totalTime float64 = 0
    
   	reader := bufio.NewReader(os.Stdin)
    for {
    	var frameStartTime = time.Now()

    	Draw(scene)
		command, _, err := reader.ReadRune()
		if err == nil{
			ExecuteUserCommand(command, camera)
		}

		deltaTime = time.Now().Sub(frameStartTime).Seconds()
		totalTime += deltaTime
    }
}
