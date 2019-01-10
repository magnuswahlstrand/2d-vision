package internal

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func DrawMarker(screen *ebiten.Image, x, y float64, c color.Color, r float64) {
	ebitenutil.DrawLine(screen, x-r, y-r, x-r, y+r, c)
	ebitenutil.DrawLine(screen, x-r, y+r, x+r, y+r, c)
	ebitenutil.DrawLine(screen, x+r, y+r, x+r, y-r, c)
	ebitenutil.DrawLine(screen, x+r, y-r, x-r, y-r, c)
}

func DrawInstructions(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf(`

      WASD/Drag: Move
      Q: Toggle rays








              TPS: %0.0f
`, ebiten.CurrentTPS()))
}
