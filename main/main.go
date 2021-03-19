package main
import (
	"fmt"
	"raymarch"
	"strings"
	"math"
	"time"
)

func Clear(){
	for i := 0; i < 100; i++{
		fmt.Printf("\n")
	}
}

func main() {
	var pixAspectRatio float64 = 2
	var resX, resY = 200, 70

	var intensityChars = []string{" ", ".","-", "~",  "\"", "+", "/", "#", "@"}

	var sphere1 = raymarch.Sphere{raymarch.Vector3{}, 0.5}
	var sphere2 = raymarch.Sphere{raymarch.Vector3{-3, -0.2, -7}, 0.5}
	var sphere3 = raymarch.Sphere{raymarch.Vector3{2, 0.4, -5}, 0.5}

	var geom = [](*raymarch.Geometry){}
	append(geom, sphere1)
	append(geom, sphere2)
	append(geom, sphere3)

	var camera = raymarch.Camera{raymarch.Vector3{0, 0, 2}, raymarch.Vector3{0, 0, -1}, raymarch.Vector3{0, 1, 0}, math.Pi / 4}
    var scene = raymarch.Scene{
    	camera,
    	geom}
    var marcher = raymarch.Raymarcher{0.7, 10.0, 0.01}
    var lightDir = raymarch.Vector3{-1, -1, -1}.Normalised()
    var deltaTime = 0.1
    for {
    	var frameStartTime = time.Now()

    	sphere1.Pos.Z -= deltaTime

    	var screenIntensities = make([]float64, resX * resY)
	    for x := 0; x < resX; x++ {
	    	for y := 0; y < resY; y++ {
	    		var xSS = 2 * float64(x) / float64(resX) - 1
	    		var ySS = pixAspectRatio * (2 * (float64(y) + float64(resX - resY) / 2) / (float64(resX)) - 1)
	    		var screenPos = raymarch.Vector3{xSS, ySS, 0}
	    		screenIntensities[x + resX * y] = marcher.BlinnPhong(screenPos, scene, 0.2, lightDir)
	    	}
	    }

	   	var screenChars = make([]string, resX * resY)
	   	for i := 0; i < len(screenChars); i++ {
			var charIx = int(screenIntensities[i] * float64(len(intensityChars)))
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
		time.Sleep(100 * time.Millisecond)
		deltaTime = time.Now().Sub(frameStartTime).Seconds()
    }
}
