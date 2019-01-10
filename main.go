package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/kyeett/2d-vision/resources"
	geo "github.com/paulmach/go.geo"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var floorImage *ebiten.Image
var maskedFgImage *ebiten.Image
var blackImage *ebiten.Image

const (
	screenWidth  = 320
	screenHeight = 320
)

func init() {

	// f, err := os.Open("resources/floor.png")
	// if err != nil {
	// 	log.Fatal("failed to open file", err)
	// }
	// defer f.Close()

	img, _, err := image.Decode(bytes.NewReader(resources.Floor_png))
	if err != nil {
		log.Fatal("failed to decode image", err)
	}
	floorImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	maskedFgImage, _ = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	blackImage, _ = ebiten.NewImage(screenWidth*2, screenHeight*2, ebiten.FilterDefault)
}

func drawFloor(screen *ebiten.Image) {
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {

			op := &ebiten.DrawImageOptions{}
			b := floorImage.Bounds()
			op.GeoM.Translate(float64(x*b.Dx()), float64(y*b.Dy()))
			screen.DrawImage(floorImage, op)
		}
	}

}

var (
	x = screenWidth / 2.
	y = screenHeight / 2.
)

type Triangle struct {
	X1, Y1, X2, Y2, X3, Y3 float32
}

func T(x1, y1, x2, y2, x3, y3 float64) Triangle {
	return Triangle{
		float32(x1), float32(y1),
		float32(x2), float32(y2),
		float32(x3), float32(y3),
	}
}

func (t Triangle) Offset(dx, dy float64) Triangle {
	t.X1 += float32(dx)
	t.X2 += float32(dx)
	t.X3 += float32(dx)

	t.Y1 += float32(dy)
	t.Y2 += float32(dy)
	t.Y3 += float32(dy)
	return t
}

func (t Triangle) Vertices() []ebiten.Vertex {
	v := ebiten.Vertex{
		SrcX:   0,
		SrcY:   0,
		ColorR: 1,
		ColorG: 1,
		ColorB: 1,
		ColorA: 1,
	}

	v1 := v
	v2 := v
	v3 := v
	v1.DstX, v1.DstY = float32(t.X1), float32(t.Y1)
	v2.DstX, v2.DstY = float32(t.X2), float32(t.Y2)
	v3.DstX, v3.DstY = float32(t.X3), float32(t.Y3)
	return []ebiten.Vertex{v1, v2, v3}
}

var (
	colorRed         = color.RGBA{255, 0, 0, 255}
	colorYellow      = color.RGBA{255, 255, 0, 255}
	colorFadedYellow = color.RGBA{100, 100, 0, 100}
)

var debug bool

func update(screen *ebiten.Image) error {

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("Game terminated by player")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		debug = !debug
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		x += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		x -= 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		y += 2
	}

	screen.Fill(color.RGBA{255, 0, 0, 255})
	drawFloor(screen)

	lines := []segment{}

	sparser := 2.
	for i := 0.; i < (360 / sparser); i++ {
		length := float64(screenWidth + screenHeight)

		angle := sparser * math.Pi * i / 180
		start := geo.NewPoint(x, y)
		end := geo.NewPoint(x, y).Add(geo.NewPoint(length*math.Cos(angle), length*math.Sin(angle)))

		line := geo.NewLine(start, end)

		// Check intersection with all walls

		points := []*geo.Point{}
		for _, r := range []image.Rectangle{game, box, box2, box3} {
			path := geoPathFromRect(r)
			tmp, _ := path.Intersection(line)
			points = append(points, tmp...)
		}

		// Find closest point
		min := math.Inf(1)
		minP := &geo.Point{}
		for i, _ := range points {

			d := points[i].DistanceFrom(geo.NewPoint(x, y))

			if d < min {
				min = d
				minP = points[i]
			}

			if debug {
				drawMarker(screen, points[i].X(), points[i].Y(), colorRed)
			}
		}

		lines = append(lines, Seg(x, y, minP.X(), minP.Y()))
	}

	blackImage.Fill(color.Black)
	op := &ebiten.DrawImageOptions{}

	opt := &ebiten.DrawTrianglesOptions{}
	opt.Address = ebiten.AddressRepeat
	opt.CompositeMode = ebiten.CompositeModeSourceOut

	prevLine := lines[len(lines)-1]
	for _, line := range lines {
		// Ray lines
		if debug {
			ebitenutil.DrawLine(screen, line.X1, line.Y1, line.X2, line.Y2, colorYellow)
			// Markers at intersection
			drawMarker(screen, line.X2, line.Y2, colorYellow)

		}

		// Area between lines
		t2 := T(x, y, prevLine.X2, prevLine.Y2, line.X2, line.Y2)
		blackImage.DrawTriangles(t2.Vertices(), []uint16{0, 1, 2}, maskedFgImage, opt)
		prevLine = line
	}

	op.ColorM.Scale(1, 1, 1, 0.8) // Make transparent
	screen.DrawImage(blackImage, op)

	for _, wall := range walls {
		ebitenutil.DrawLine(screen, float64(wall.X1), float64(wall.Y1), float64(wall.X2), float64(wall.Y2), colorRed)
	}

	// Center marker
	drawMarker(screen, x, y, colorRed)

	ebitenutil.DebugPrint(screen, fmt.Sprintf(`

      WASD: Move
      Q: Toggle rays








              TPS: %0.0f
`, ebiten.CurrentTPS()))
	return nil
}

type segment struct {
	X1, Y1, X2, Y2 float64
}

func Seg(x1, y1, x2, y2 float64) segment {
	return segment{x1, y1, x2, y2}
}

var walls []segment

func segmentsFromRect(r image.Rectangle) []segment {
	s := []segment{}
	s = append(s, Seg(float64(r.Min.X), float64(r.Min.Y), float64(r.Min.X), float64(r.Max.Y)))
	s = append(s, Seg(float64(r.Min.X), float64(r.Max.Y), float64(r.Max.X), float64(r.Max.Y)))
	s = append(s, Seg(float64(r.Max.X), float64(r.Max.Y), float64(r.Max.X), float64(r.Min.Y)))
	s = append(s, Seg(float64(r.Max.X), float64(r.Min.Y), float64(r.Min.X), float64(r.Min.Y)))
	return s
}

func geoPathFromRect(r image.Rectangle) *geo.Path {
	path := geo.NewPath()
	path.Push(geo.NewPoint(float64(r.Min.X), float64(r.Min.Y)))
	path.Push(geo.NewPoint(float64(r.Max.X), float64(r.Min.Y)))
	path.Push(geo.NewPoint(float64(r.Max.X), float64(r.Max.Y)))
	path.Push(geo.NewPoint(float64(r.Min.X), float64(r.Max.Y)))
	path.Push(geo.NewPoint(float64(r.Min.X), float64(r.Min.Y)))
	return path
}

func drawMarker(screen *ebiten.Image, x, y float64, c color.Color) {
	ebitenutil.DrawLine(screen, x-1, y-1, x-1, y+1, c)
	ebitenutil.DrawLine(screen, x-1, y+1, x+1, y+1, c)
	ebitenutil.DrawLine(screen, x+1, y+1, x+1, y-1, c)
	ebitenutil.DrawLine(screen, x+1, y-1, x-1, y-1, c)
}

var game image.Rectangle
var box image.Rectangle
var box2 image.Rectangle
var box3 image.Rectangle

func main() {
	padd := 10
	walls = []segment{}

	game = image.Rect(0, 0, screenWidth-2*padd, screenWidth-2*padd).Add(image.Pt(padd, padd))
	walls = append(walls, segmentsFromRect(game)...)
	// walls = append(walls, Seg(padd, padd, padd, screenWidth-padd))
	// walls = append(walls, Seg(padd, screenWidth-padd, screenWidth-padd, screenWidth-padd))
	// walls = append(walls, Seg(screenWidth-padd, screenWidth-padd, screenWidth-padd, padd))
	// walls = append(walls, Seg(screenWidth-padd, padd, padd, padd))

	box = image.Rect(0, 0, 100, 100).Add(image.Pt(30, 30))
	walls = append(walls, segmentsFromRect(box)...)

	box2 = image.Rect(0, 0, 30, 30).Add(image.Pt(230, 200))
	walls = append(walls, segmentsFromRect(box2)...)

	box3 = image.Rect(0, 0, 70, 70).Add(image.Pt(80, 180))
	walls = append(walls, segmentsFromRect(box3)...)

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "2D Raycasting Demo"); err != nil {
		log.Fatal("Game exited: ", err)

	}
}
