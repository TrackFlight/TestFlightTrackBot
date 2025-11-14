package utils

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type Align int

const (
	AlignVerticallyStart Align = 1 << iota
	AlignHorizontallyEnd
	AlignHorizontallyCenter
	AlignHorizontallyStart
	AlignVerticallyCenter
	AlignVerticallyEnd
)

type Orientation int

const (
	OrientationHorizontal Orientation = iota
	OrientationVertical
)

const (
	SizeAuto = iota * -1
	SizeFitParent
)

type Rect struct {
	width, height                                    int
	x, y                                             int
	marginLeft, marginTop, marginRight, marginBottom int
	inset                                            int
	alignment                                        Align
	orientation                                      Orientation
	children                                         []*Rect
	parent                                           *Rect
	index                                            int
}

type TextRect struct {
	rect *Rect
	text string
	font font.Face
}

func (tr *TextRect) Width() int {
	return tr.rect.Width()
}

func (tr *TextRect) Height() int {
	return tr.rect.Height()
}

func (tr *TextRect) SetMargin(left, top, right, bottom int) {
	tr.rect.SetMargin(left, top, right, bottom)
}

func (tr *TextRect) X() int {
	return tr.rect.X()
}

func (tr *TextRect) Y() int {
	return tr.rect.Y()
}

func (tr *TextRect) String() string {
	return tr.text
}

func (tr *TextRect) Font() font.Face {
	return tr.font
}

func NewRect(width, height int) *Rect {
	return &Rect{
		width:       width,
		height:      height,
		orientation: OrientationVertical,
		alignment:   AlignVerticallyStart | AlignHorizontallyStart,
	}
}

func (r *Rect) SetInset(inset int) {
	r.inset = inset
}

func (r *Rect) SetAlignment(alignment Align) {
	r.alignment = alignment
}

func (r *Rect) SetOrientation(orientation Orientation) {
	r.orientation = orientation
}

func (r *Rect) SetMargin(left, top, right, bottom int) {
	r.marginLeft = left
	r.marginTop = top
	r.marginRight = right
	r.marginBottom = bottom
}

func (r *Rect) Width() int {
	return r.measure(
		r.width,
		r.orientation == OrientationHorizontal,
		func(c *Rect) int { return c.Width() },
	)
}

func (r *Rect) Height() int {
	return r.measure(
		r.height,
		r.orientation == OrientationVertical,
		func(c *Rect) int { return c.Height() },
	)
}

func (r *Rect) measure(
	fixed int,
	isStack bool,
	childSize func(*Rect) int,
) int {
	if fixed == SizeFitParent {
		return childSize(r.parent)
	} else if fixed != SizeAuto {
		return fixed
	}
	if isStack {
		total := r.inset * (len(r.children) - 1)
		for _, child := range r.children {
			total += childSize(child)
		}
		return total
	}
	m := 0
	for _, child := range r.children {
		if v := childSize(child); v > m {
			m = v
		}
	}
	return m
}

func (r *Rect) computePos(
	offset int,
	margin int,
	isStack bool,
	size func(*Rect) int,
	marginStart func(*Rect) int,
	marginEnd func(*Rect) int,
	parentSize func(*Rect) int,
	parentPos func(*Rect) int,
	alignCenter bool,
	alignEnd bool,
) int {
	if r.parent == nil {
		return offset
	}

	p := r.parent

	pos := parentPos(p) + offset + margin

	total := size(r)

	if isStack {
		for i := 0; i < r.index; i++ {
			pos += size(p.children[i]) + marginStart(p.children[i]) + p.inset
			total += size(p.children[i]) + marginStart(p.children[i]) + marginEnd(p.children[i])
		}
		for i := r.index + 1; i < len(p.children); i++ {
			total += size(p.children[i]) + marginStart(p.children[i]) + marginEnd(p.children[i])
		}
		total += p.inset * (len(p.children) - 1)
	}

	parentDim := parentSize(p)

	if alignCenter {
		pos += parentDim/2 - total/2
	} else if alignEnd {
		pos += parentDim - total
	}

	return pos
}

func (r *Rect) Children() []*Rect {
	return r.children
}

func (r *Rect) X() int {
	if r.parent == nil {
		return r.x
	}

	p := r.parent

	return r.computePos(
		r.x,
		r.marginLeft,
		p.orientation == OrientationHorizontal,
		func(c *Rect) int { return c.Width() },
		func(r *Rect) int { return r.marginLeft },
		func(r *Rect) int { return r.marginRight },
		func(p *Rect) int { return p.Width() },
		func(p *Rect) int { return p.X() },
		p.alignment&AlignHorizontallyCenter != 0,
		p.alignment&AlignHorizontallyEnd != 0,
	)
}

func (r *Rect) Y() int {
	if r.parent == nil {
		return r.y
	}

	p := r.parent

	return r.computePos(
		r.y,
		r.marginTop,
		p.orientation == OrientationVertical,
		func(c *Rect) int { return c.Height() },
		func(r *Rect) int { return r.marginTop },
		func(r *Rect) int { return r.marginBottom },
		func(p *Rect) int { return p.Height() },
		func(p *Rect) int { return p.Y() },
		p.alignment&AlignVerticallyCenter != 0,
		p.alignment&AlignVerticallyEnd != 0,
	)
}

func (r *Rect) AddRectChild(width, height int) *Rect {
	return r.addChildInternal(width, height)
}

func (r *Rect) AddTextChild(dc *gg.Context, text string, ft *truetype.Font, size float64) *TextRect {
	face := truetype.NewFace(ft, &truetype.Options{Size: size})
	dc.SetFontFace(face)
	w, _ := dc.MeasureString(text)
	m := face.Metrics()
	child := r.addChildInternal(int(w), int((m.Ascent+m.Descent)>>6))
	return &TextRect{
		child,
		text,
		face,
	}
}

func (r *Rect) AddLayoutChild(width, height int) *Rect {
	return r.addChildInternal(width, height)
}

func (r *Rect) addChildInternal(width, height int) *Rect {
	child := NewRect(width, height)
	child.parent = r
	child.index = len(r.children)
	r.children = append(r.children, child)
	return child
}
