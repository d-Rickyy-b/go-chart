package chart

import "github.com/wcharczuk/go-chart/drawing"

// DrawLineSeries draws a line series with a renderer.
func DrawLineSeries(r Renderer, canvasBox Box, xrange, yrange Range, s Style, vs ValueProvider) {
	if vs.Len() == 0 {
		return
	}

	cb := canvasBox.Bottom
	cl := canvasBox.Left

	v0x, v0y := vs.GetValue(0)
	x0 := cl + xrange.Translate(v0x)
	y0 := cb - yrange.Translate(v0y)

	var vx, vy float64
	var x, y int

	fill := s.GetFillColor()
	if !fill.IsZero() {
		r.SetFillColor(fill)
		r.MoveTo(x0, y0)
		for i := 1; i < vs.Len(); i++ {
			vx, vy = vs.GetValue(i)
			x = cl + xrange.Translate(vx)
			y = cb - yrange.Translate(vy)
			r.LineTo(x, y)
		}
		r.LineTo(x, cb)
		r.LineTo(x0, cb)
		r.Close()
		r.Fill()
	}

	r.SetStrokeColor(s.GetStrokeColor())
	r.SetStrokeDashArray(s.GetStrokeDashArray())
	r.SetStrokeWidth(s.GetStrokeWidth(DefaultStrokeWidth))

	r.MoveTo(x0, y0)
	for i := 1; i < vs.Len(); i++ {
		vx, vy = vs.GetValue(i)
		x = cl + xrange.Translate(vx)
		y = cb - yrange.Translate(vy)
		r.LineTo(x, y)
	}
	r.Stroke()
}

// MeasureAnnotation measures how big an annotation would be.
func MeasureAnnotation(r Renderer, canvasBox Box, s Style, lx, ly int, label string) Box {
	r.SetFillColor(s.GetFillColor(DefaultAnnotationFillColor))
	r.SetStrokeColor(s.GetStrokeColor())
	r.SetStrokeWidth(s.GetStrokeWidth())
	r.SetFont(s.GetFont())
	r.SetFontColor(s.GetFontColor(DefaultTextColor))
	r.SetFontSize(s.GetFontSize(DefaultAnnotationFontSize))

	textBox := r.MeasureText(label)
	textWidth := textBox.Width()
	textHeight := textBox.Height()
	halfTextHeight := textHeight >> 1

	pt := s.Padding.GetTop(DefaultAnnotationPadding.Top)
	pl := s.Padding.GetLeft(DefaultAnnotationPadding.Left)
	pr := s.Padding.GetRight(DefaultAnnotationPadding.Right)
	pb := s.Padding.GetBottom(DefaultAnnotationPadding.Bottom)

	strokeWidth := s.GetStrokeWidth()

	top := ly - (pt + halfTextHeight)
	right := lx + pl + pr + textWidth + DefaultAnnotationDeltaWidth + int(strokeWidth)
	bottom := ly + (pb + halfTextHeight)

	return Box{
		Top:    top,
		Left:   lx,
		Right:  right,
		Bottom: bottom,
	}
}

// DrawAnnotation draws an anotation with a renderer.
func DrawAnnotation(r Renderer, canvasBox Box, s Style, lx, ly int, label string) {
	r.SetFillColor(s.GetFillColor(DefaultAnnotationFillColor))
	r.SetStrokeColor(s.GetStrokeColor())
	r.SetStrokeWidth(s.GetStrokeWidth())
	r.SetStrokeDashArray(s.GetStrokeDashArray())

	textBox := r.MeasureText(label)
	textWidth := textBox.Width()
	halfTextHeight := textBox.Height() >> 1

	pt := s.Padding.GetTop(DefaultAnnotationPadding.Top)
	pl := s.Padding.GetLeft(DefaultAnnotationPadding.Left)
	pr := s.Padding.GetRight(DefaultAnnotationPadding.Right)
	pb := s.Padding.GetBottom(DefaultAnnotationPadding.Bottom)

	textX := lx + pl + DefaultAnnotationDeltaWidth
	textY := ly + halfTextHeight

	ltx := lx + DefaultAnnotationDeltaWidth
	lty := ly - (pt + halfTextHeight)

	rtx := lx + pl + pr + textWidth + DefaultAnnotationDeltaWidth
	rty := ly - (pt + halfTextHeight)

	rbx := lx + pl + pr + textWidth + DefaultAnnotationDeltaWidth
	rby := ly + (pb + halfTextHeight)

	lbx := lx + DefaultAnnotationDeltaWidth
	lby := ly + (pb + halfTextHeight)

	r.MoveTo(lx, ly)
	r.LineTo(ltx, lty)
	r.LineTo(rtx, rty)
	r.LineTo(rbx, rby)
	r.LineTo(lbx, lby)
	r.LineTo(lx, ly)
	r.Close()
	r.FillStroke()

	r.SetFont(s.GetFont())
	r.SetFontColor(s.GetFontColor(DefaultTextColor))
	r.SetFontSize(s.GetFontSize(DefaultAnnotationFontSize))

	r.Text(label, textX, textY)
}

// DrawBox draws a box with a given style.
func DrawBox(r Renderer, b Box, s Style) {
	r.SetFillColor(s.GetFillColor())
	r.SetStrokeColor(s.GetStrokeColor(DefaultStrokeColor))
	r.SetStrokeWidth(s.GetStrokeWidth(DefaultStrokeWidth))
	r.SetStrokeDashArray(s.GetStrokeDashArray())

	r.MoveTo(b.Left, b.Top)
	r.LineTo(b.Right, b.Top)
	r.LineTo(b.Right, b.Bottom)
	r.LineTo(b.Left, b.Bottom)
	r.LineTo(b.Left, b.Top)
	r.FillStroke()
}

// DrawText draws text with a given style.
func DrawText(r Renderer, text string, x, y int, s Style) {
	r.SetFontColor(s.GetFontColor(DefaultTextColor))
	r.SetStrokeColor(s.GetStrokeColor())
	r.SetStrokeWidth(s.GetStrokeWidth())
	r.SetFont(s.GetFont())
	r.SetFontSize(s.GetFontSize())

	r.Text(text, x, y)
}

// DrawTextCentered draws text with a given style centered.
func DrawTextCentered(r Renderer, text string, x, y int, s Style) {
	r.SetFontColor(s.GetFontColor(DefaultTextColor))
	r.SetStrokeColor(s.GetStrokeColor())
	r.SetStrokeWidth(s.GetStrokeWidth())
	r.SetFont(s.GetFont())
	r.SetFontSize(s.GetFontSize())

	tb := r.MeasureText(text)
	tx := x - (tb.Width() >> 1)
	ty := y - (tb.Height() >> 1)
	r.Text(text, tx, ty)
}

// CreateLegend returns a legend renderable function.
func CreateLegend(c *Chart, style Style) Renderable {
	return func(r Renderer, cb Box, defaults Style) {
		workingStyle := style.WithDefaultsFrom(defaults.WithDefaultsFrom(Style{
			FillColor:   drawing.ColorWhite,
			FontColor:   DefaultTextColor,
			FontSize:    8.0,
			StrokeColor: DefaultAxisColor,
			StrokeWidth: DefaultAxisLineWidth,
		}))

		// DEFAULTS
		legendPadding := 5
		lineTextGap := 5
		lineLengthMinimum := 25

		var labels []string
		var lines []Style
		for _, s := range c.Series {
			if s.GetStyle().IsZero() || s.GetStyle().Show {
				if _, isAnnotationSeries := s.(AnnotationSeries); !isAnnotationSeries {
					labels = append(labels, s.GetName())
					lines = append(lines, s.GetStyle())
				}
			}
		}

		legend := Box{
			Top:  cb.Top, //padding
			Left: cb.Left,
		}

		legendContent := Box{
			Top:  legend.Top + legendPadding,
			Left: legend.Left + legendPadding,
		}

		r.SetFontColor(workingStyle.GetFontColor())
		r.SetFontSize(workingStyle.GetFontSize())

		// measure
		for x := 0; x < len(labels); x++ {
			if len(labels[x]) > 0 {
				tb := r.MeasureText(labels[x])
				legendContent.Bottom += (tb.Height() + DefaultMinimumTickVerticalSpacing)
				rowRight := tb.Width() + legendContent.Left + lineLengthMinimum + lineTextGap
				legendContent.Right = MaxInt(legendContent.Right, rowRight)
			}
		}

		legend = legend.Grow(legendContent)
		DrawBox(r, legend, workingStyle)

		legendContent.Right = legend.Right - legendPadding
		legendContent.Bottom = legend.Bottom - legendPadding

		ycursor := legendContent.Top
		tx := legendContent.Left
		for x := 0; x < len(labels); x++ {
			if len(labels[x]) > 0 {
				tb := r.MeasureText(labels[x])
				ycursor += tb.Height()

				//r.SetFillColor(DefaultTextColor)
				r.Text(labels[x], tx, ycursor)
				th2 := tb.Height() >> 1

				lx := tx + tb.Width() + lineTextGap
				ly := ycursor - th2
				lx2 := legendContent.Right - legendPadding

				r.SetStrokeColor(lines[x].GetStrokeColor())
				r.SetStrokeWidth(lines[x].GetStrokeWidth())
				r.SetStrokeDashArray(lines[x].GetStrokeDashArray())

				r.MoveTo(lx, ly)
				r.LineTo(lx2, ly)
				r.Stroke()

				ycursor += DefaultMinimumTickVerticalSpacing
			}
		}
	}
}
