package calendar

import (
	"bytes"
	"fmt"
	"image/png"
	"time"

	"github.com/fogleman/gg"
)

var colors = map[string]Color{
	"green": {
		R: 0,
		G: 1,
		B: 0,
		A: 0.5,
	},
	"blue": {
		R: 0,
		G: 0,
		B: 1,
		A: 0.5,
	},
	"red": {
		R: 1,
		G: 0,
		B: 0,
		A: 0.5,
	},
	"yellow": {
		R: 1,
		G: 1,
		B: 0,
		A: 0.5,
	},
}

type Event struct {
	Day       int    `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Title     string `json:"title"`
	Color     string `json:"color"`
}

type Color struct {
	R float64
	G float64
	B float64
	A float64
}

type Margin struct {
	X int
	Y int
}

func GenerateCalendar(events []Event) ([]byte, error) {
	margin := Margin{X: 50, Y: 25}
	width, height := 1000, 600
	dc := gg.NewContext(width, height)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	daysOfWeek := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	drawCalendarGrid(dc, width, height, daysOfWeek, margin)

	err := drawEvents(dc, events, daysOfWeek, margin)
	if err != nil {
		return []byte{}, err
	}

	out := new(bytes.Buffer)
	err = png.Encode(out, dc.Image())
	if err != nil {
		return []byte{}, err
	}

	return out.Bytes(), nil
}

func drawCalendarGrid(
	dc *gg.Context,
	width,
	height int,
	daysOfWeek []string,
	margin Margin,
) {
	dc.SetRGB(0, 0, 0)

	hours := 24

	gridWidth := width - margin.X
	gridHeight := height - margin.Y

	for i := 0; i <= hours; i++ {
		y := (float64(i) * float64(gridHeight) / float64(hours)) + float64(margin.Y)
		dc.DrawLine(float64(margin.X), y, float64(width), y)
		dc.Stroke()

		dc.DrawStringAnchored(fmt.Sprintf("%01d:00", i), float64(margin.X-5), y-5, 1, 1)
	}

	for i, day := range daysOfWeek {
		x := float64(i)*float64(gridWidth)/float64(len(daysOfWeek)) + float64(margin.X)
		dc.DrawLine(x, float64(margin.Y), x, float64(gridHeight+margin.Y))
		dc.Stroke()

		dc.DrawStringAnchored(day, x+float64(gridWidth)/(2*float64(len(daysOfWeek))), float64(margin.Y)-10, 0.5, 0)
	}
}

func drawEvents(
	dc *gg.Context,
	events []Event,
	daysOfWeek []string,
	margin Margin,
) error {
	dayWidth := float64(dc.Width()-margin.X) / float64(len(daysOfWeek))

	for _, event := range events {
		if _, ok := colors[event.Color]; !ok {
			return fmt.Errorf("the color for the %s event is not in the dictionary. the choices are red, blue, green, and yellow", event.Title)
		}

		startHour, startMin, err := parseTime(event.StartTime)
		if err != nil {
			return err
		}

		endHour, endMin, err := parseTime(event.EndTime)
		if err != nil {
			return err
		}

		startY := ((float64(startHour) + float64(startMin/60)) * float64(dc.Height()-margin.Y)) / 24
		endY := (float64(endHour) + float64(endMin)/60) * float64(dc.Height()-margin.Y) / 24

		startY += float64(margin.Y)
		endY += float64(margin.Y)

		x := (float64(event.Day) * dayWidth) + float64(margin.X)

		color := colors[event.Color]

		fmt.Println()

		dc.SetRGBA(color.R, color.G, color.B, color.A)
		dc.DrawRectangle(x, startY, dayWidth, endY-startY)
		dc.Fill()

		dc.SetRGB(0, 0, 0)
		dc.DrawStringWrapped(event.Title, x+5, startY+5, 0, 0, dayWidth-10, 1.5, gg.AlignLeft)
	}

	return nil
}

func parseTime(timeStr string) (int, int, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return 0, 0, err
	}

	return t.Hour(), t.Minute(), nil
}
