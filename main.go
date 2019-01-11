package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"

	"golang.org/x/image/colornames"

	"github.com/SolarLune/resolv/resolv"
	"github.com/kyeett/2d-vision/internal"
	"github.com/kyeett/2d-vision/resources"

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
var dragging bool
var draggingType string

const size = 8

func (g *Game) update(screen *ebiten.Image) error {
	var dragX, dragY int
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		if (float64(cx)-x)*(float64(cx)-x) < size*size && (float64(cy)-y)*(float64(cy)-y) < size*size {
			dragging = true
			draggingType = "mouse"
		}
	}

	if len(inpututil.JustPressedTouchIDs()) > 0 {
		cx, cy := ebiten.TouchPosition(0)
		if (float64(cx)-x)*(float64(cx)-x) < size*size && (float64(cy)-y)*(float64(cy)-y) < size*size {
			fmt.Println("Start dragg")
			dragging = true
			draggingType = "touch"
		}
	}

	if inpututil.IsTouchJustReleased(0) {
		dragging = false
		draggingType = ""
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && draggingType == "mouse" {
		dragging = false
		draggingType = ""
	}

	var tmpX, tmpY float64
	if dragging {
		switch draggingType {
		case "mouse":
			dragX, dragY = ebiten.CursorPosition()
		case "touch":
			dragX, dragY = ebiten.TouchPosition(0)
		}

		r := g.space.Get(0).(*resolv.Rectangle)
		vX, vY := int32(dragX-int(x)), int32(dragY-int(y))
		if vX > 8 {
			vX = 8
		}

		if vX < -8 {
			vX = -8
		}

		if vY > 8 {
			vY = 8
		}

		if vY < -8 {
			vY = -8
		}

		if res := g.space.Resolve(r, vX, 0); res.Colliding() && !res.Teleporting {
			fmt.Println("Collide!")
		} else {
			r.Move(vX, 0)
			x += float64(vX)
		}

		if res := g.space.Resolve(r, 0, vY); res.Colliding() && !res.Teleporting {
			fmt.Println("Collide!")
		} else {
			r.Move(0, vY)
			y += float64(vY)
		}
	}

	// cx, cy := ebiten.CursorPosition()
	// fmt.Println(cx, cy)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("Game terminated by player")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		debug = !debug
	}

	var vx, vy float64
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		vx = 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		vx = -4
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		vy = -4
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		vy = 4
	}

	fmt.Println(tmpX, tmpY)

	// player := resolv.NewRectangle(int32(x-size), int32(y-size), int32(size*2), int32(size*2))
	r := g.space.Get(0).(*resolv.Rectangle)
	if res := g.space.Resolve(r, int32(vx), 0); res.Colliding() && !res.Teleporting {
		fmt.Println("Collide!")
	} else {
		x += vx
		r.Move(int32(vx), 0)
	}

	if res := g.space.Resolve(r, 0, int32(vy)); res.Colliding() && !res.Teleporting {
		fmt.Println("Collide!")
	} else {
		y += vy
		r.Move(0, int32(vy))
	}

	// p := g.space.Get(0)
	// p.SetXY(int32(x-size), int32(y-size))
	fmt.Println(g.space.Get(0), int32(vx), int32(vy))
	// if res := g.space.Resolve(r, 0, vy); res.Colliding() && !res.Teleporting {
	// 	overrideColor = colornames.Red
	// } else {
	// 	y += vy
	// }

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.Fill(color.RGBA{255, 0, 0, 255})
	drawFloor(screen)

	// lines := internal.BasicRayCasting(x, y, []image.Rectangle{outer, box, box2, box3})
	lines := internal.SmartRayCasting(x, y, g.objects)

	blackImage.Fill(color.Black)
	op := &ebiten.DrawImageOptions{}

	opt := &ebiten.DrawTrianglesOptions{}
	opt.Address = ebiten.AddressRepeat
	opt.CompositeMode = ebiten.CompositeModeSourceOut

	prevLine := lines[len(lines)-1]
	for _, line := range lines {

		// Area between lines
		t2 := T(x, y, prevLine.X2, prevLine.Y2, line.X2, line.Y2)
		blackImage.DrawTriangles(t2.Vertices(), []uint16{0, 1, 2}, maskedFgImage, opt)
		prevLine = line

		// Ray lines
		if debug {
			ebitenutil.DrawLine(screen, line.X1, line.Y1, line.X2, line.Y2, colorYellow)
			ebitenutil.DrawLine(screen, prevLine.X1, prevLine.Y1, prevLine.X2, prevLine.Y2, colorYellow)
			// Markers at intersection
			internal.DrawMarker(screen, line.X2, line.Y2, colorYellow, 1)
		}
	}

	op.ColorM.Scale(1, 1, 1, 0.7) // Make transparent
	screen.DrawImage(blackImage, op)

	for _, wall := range g.walls {
		ebitenutil.DrawLine(screen, float64(wall.X1), float64(wall.Y1), float64(wall.X2), float64(wall.Y2), colornames.Darkred)
	}

	for _, o := range g.objects[1:] {
		drawRect(screen, int32(o.Min.X), int32(o.Min.Y), int32(o.Dx()), int32(o.Dy()), color.RGBA{0, 0, 0, 150})
	}

	// Center marker
	internal.DrawMarker(screen, x, y, colorYellow, size)
	internal.DrawInstructions(screen, len(lines), debug)
	return nil
}

func drawRect(screen *ebiten.Image, x, y, w, h int32, c color.Color) {
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), c)
}

type Game struct {
	walls   []internal.Segment
	objects []image.Rectangle
	space   resolv.Space
}

func main() {

	g := Game{}

	padd := 10
	g.walls = []internal.Segment{}

	outer := image.Rect(0, 0, screenWidth-2*padd, screenWidth-2*padd).Add(image.Pt(padd, padd))
	g.walls = append(g.walls, internal.SegmentsFromRect(outer)...)
	g.objects = append(g.objects, outer)

	box := image.Rect(0, 0, 110, 110).Add(image.Pt(30, 30))
	g.walls = append(g.walls, internal.SegmentsFromRect(box)...)
	g.objects = append(g.objects, box)

	box2 := image.Rect(0, 0, 30, 30).Add(image.Pt(230, 200))
	g.walls = append(g.walls, internal.SegmentsFromRect(box2)...)
	g.objects = append(g.objects, box2)

	box3 := image.Rect(0, 0, 70, 70).Add(image.Pt(80, 180))
	g.walls = append(g.walls, internal.SegmentsFromRect(box3)...)
	g.objects = append(g.objects, box3)

	box4 := image.Rect(0, 0, 100, 30).Add(image.Pt(165, 30))
	g.walls = append(g.walls, internal.SegmentsFromRect(box4)...)
	g.objects = append(g.objects, box4)

	g.space.AddShape(
		resolv.NewRectangle(int32(x-size), int32(y-size), int32(size*2), int32(size*2)),
		internal.ShapeFromRect(box),
		internal.ShapeFromRect(box2),
		internal.ShapeFromRect(box3),
		internal.ShapeFromRect(box4))
	g.space.AddShape(internal.LinesFromRect(outer)...)

	if err := ebiten.Run(g.update, screenWidth, screenHeight, 1.5, "2D Raycasting Demo"); err != nil {
		log.Fatal("Game exited: ", err)

	}
}
