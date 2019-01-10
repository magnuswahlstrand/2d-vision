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

func DrawInstructions(screen *ebiten.Image, numLines int, debug bool) {

	var linesText, TPSText string
	if debug {
		linesText = fmt.Sprintf("# Lines: %d", numLines)
		TPSText = fmt.Sprintf("TPS: %0.0f", ebiten.CurrentTPS())
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf(`

      WASD/Drag: Move       %s
      Q: Toggle rays








              %s
`, linesText, TPSText))
}
