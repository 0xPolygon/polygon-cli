package txgaschart

import (
	"gonum.org/v1/plot/vg"
	vgd "gonum.org/v1/plot/vg/draw"
)

// ThickCrossGlyph draws an 'X' with configurable stroke width.
type ThickCrossGlyph struct {
	Width vg.Length
}

// DrawGlyph implements the GlyphDrawer interface.
func (g ThickCrossGlyph) DrawGlyph(c *vgd.Canvas, sty vgd.GlyphStyle, p vg.Point) {
	if !c.Contains(p) {
		return
	}
	r := sty.Radius
	ls := vgd.LineStyle{Color: sty.Color, Width: g.Width}

	// Horizontal
	h := []vg.Point{{X: p.X - r, Y: p.Y}, {X: p.X + r, Y: p.Y}}
	// Vertical
	v := []vg.Point{{X: p.X, Y: p.Y - r}, {X: p.X, Y: p.Y + r}}
	// Diagonal 1 (top-left -> bottom-right)
	d1 := []vg.Point{{X: p.X - r, Y: p.Y + r}, {X: p.X + r, Y: p.Y - r}}
	// Diagonal 2 (bottom-left -> top-right)
	d2 := []vg.Point{{X: p.X - r, Y: p.Y - r}, {X: p.X + r, Y: p.Y + r}}

	c.StrokeLines(ls, h)
	c.StrokeLines(ls, v)
	c.StrokeLines(ls, d1)
	c.StrokeLines(ls, d2)
}
