package main
import (
	"fmt"
	"raymarch"
	"strings"
	"math"
	"time"
	"os"
	"bufio"
	"sync"
	//"math/rand"
)

func Clear(){
	for i := 0; i < 100; i++{
		fmt.Printf("\n")
	}
}

func InitScene() (*raymarch.Camera, *raymarch.Scene){
	//var numHills = 30;

	var sphere1 = raymarch.Sphere{raymarch.Vector3{0, 0.5, -9}, 0.5}
	var sphere2 = raymarch.Sphere{raymarch.Vector3{-3, 0.5, -3}, 0.5}
	var sphere3 = raymarch.Sphere{raymarch.Vector3{2, 0.5, -5}, 0.5}

	var plane = raymarch.Plane{raymarch.Vector3{}, raymarch.Vector3{0, 1, 0}.Normalised()}

	var sunLight = raymarch.SunLight{ raymarch.Vector3{-1, -1, -1}, 0.2}
	var skyLight = raymarch.SunLight{ raymarch.Vector3{1, -1, 1}, 0.2}
	var pointLight = raymarch.PointLight{raymarch.Vector3{0, 3, 0}, 5.0}
	var lights = [](raymarch.Light){&sunLight, &skyLight, &pointLight}

	var doorFrame = raymarch.MakeDoorFrame(raymarch.Vector3{}, raymarch.Vector3{0, 0, 1}, 0.5, 0.8, 0.05)
	var geom = []raymarch.Geometry{&plane, &sphere1, &sphere2, &sphere3, doorFrame}

	var camera = raymarch.Camera{raymarch.Vector3{0, 0.5, 2}, raymarch.Vector3{0, 0, -1}.Normalised(), raymarch.Vector3{0, 1, 0}, math.Pi / 3}
	var scene = raymarch.Scene{
		&camera, 
		geom, 
		lights}
   	return &camera, &scene
}

func Draw(scene *raymarch.Scene){
	var intensityChars = []string{" ", ".", ":", "*", "o", "?", "8"}
	var pixAspectRatio float64 = 2
	var resX, resY = 180, 80
    var marcher = raymarch.Raymarcher{1, 50.0, 0.01}
    var screenIntensities = make([]float64, resX * resY)
    var numDrawThreads = 32
    var waitGroup sync.WaitGroup
	
	drawWorker := func(n int){
		for x := n; x < resX; x += numDrawThreads {
	    	for y := 0; y < resY; y++ {
	    		var xSS = 2 * float64(x) / float64(resX) - 1
	    		var ySS = pixAspectRatio * (2 * (float64(y) + float64(resX - resY) / 2) / (float64(resX)) - 1)
	    		var screenPos = raymarch.Vector3{xSS, ySS, 0}
	    		var hit, intensity = marcher.BlinnPhong(screenPos, scene, 4.0)
	    		if hit{
	    			screenIntensities[x + resX * y] = intensity
	    		} else {
	    			screenIntensities[x + resX * y] = -1
	    		}
	    	}
		}
		waitGroup.Done()
	}

	for n := 0; n < numDrawThreads; n++{
		waitGroup.Add(1)
		go drawWorker(n)
	}

	waitGroup.Wait()

	var maxIntensity = 0.0
	var minIntensity = 10000.0
	for i := 0; i < len(screenIntensities); i++{
		if screenIntensities[i] < 0{
			continue
		}
		maxIntensity = math.Max(screenIntensities[i], maxIntensity)
		minIntensity = math.Min(screenIntensities[i], minIntensity)
	}

   	var screenChars = make([]string, resX * resY)
   	for i := 0; i < len(screenChars); i++ {
   		var normedIntensity = (screenIntensities[i] - minIntensity) / (maxIntensity - minIntensity)
		var charIx = int(normedIntensity * float64(len(intensityChars)))
		charIx = raymarch.Max(0, raymarch.Min(len(intensityChars) - 1, charIx))
   		if screenIntensities[i] < 0{
   			charIx = 0
   		}
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
	var turnAngle = math.Pi / 24
	var turnSpeed = math.Tan(turnAngle)
	var stepSize = 0.25
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
