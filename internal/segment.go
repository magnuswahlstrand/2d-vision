package internal

import "image"

type Segment struct {
	X1, Y1, X2, Y2 float64
}

func Seg(x1, y1, x2, y2 float64) Segment {
	return Segment{x1, y1, x2, y2}
}

func SegmentsFromRect(r image.Rectangle) []Segment {
	s := []Segment{}
	s = append(s, Seg(float64(r.Min.X), float64(r.Min.Y), float64(r.Min.X), float64(r.Max.Y)))
	s = append(s, Seg(float64(r.Min.X), float64(r.Max.Y), float64(r.Max.X), float64(r.Max.Y)))
	s = append(s, Seg(float64(r.Max.X), float64(r.Max.Y), float64(r.Max.X), float64(r.Min.Y)))
	s = append(s, Seg(float64(r.Max.X), float64(r.Min.Y), float64(r.Min.X), float64(r.Min.Y)))
	return s
}
