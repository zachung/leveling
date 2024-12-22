package ui

import (
	"fmt"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	contract2 "leveling/internal/client/contract"
	"leveling/internal/client/service"
	"leveling/internal/client/ui/component"
	"leveling/internal/contract"
	"sort"
)

type World struct {
	event         contract.WorldEvent
	currentTarget string
	list          *widget.List
	entries       []any
}

type ListEntry string

func newWorld() *World {
	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage()

	entries := make([]any, 0)
	list := widget.NewList(
		// Set how wide the list should be
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(150, 0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchVertical:    true,
				Padding:            widget.NewInsetsSimple(50),
			}),
		)),
		// Set the entries in the list
		widget.ListOpts.Entries(entries),
		widget.ListOpts.ScrollContainerOpts(
			// Set the background images/color for the list
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Disabled: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Mask:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}),
		),
		widget.ListOpts.SliderOpts(
			// Set the background images/color for the background of the slider track
			widget.SliderOpts.Images(&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}, buttonImage),
			widget.SliderOpts.MinHandleSize(5),
			// Set how wide the track should be
			widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(2))),
		// Hide the horizontal slider
		widget.ListOpts.HideHorizontalSlider(),
		// Set the font for the list options
		widget.ListOpts.EntryFontFace(contract2.UiTextFace),
		// Set the colors for the list
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                   color.NRGBA{R: 0, G: 255, B: 0, A: 255},     // Foreground color for the unfocused selected entry
			Unselected:                 color.NRGBA{R: 254, G: 255, B: 255, A: 255}, // Foreground color for the unfocused unselected entry
			SelectedBackground:         color.NRGBA{R: 130, G: 130, B: 200, A: 255}, // Background color for the unfocused selected entry
			SelectingBackground:        color.NRGBA{R: 130, G: 130, B: 130, A: 255}, // Background color for the unfocused being selected entry
			SelectingFocusedBackground: color.NRGBA{R: 130, G: 140, B: 170, A: 255}, // Background color for the focused being selected entry
			SelectedFocusedBackground:  color.NRGBA{R: 130, G: 130, B: 170, A: 255}, // Background color for the focused selected entry
			FocusedBackground:          color.NRGBA{R: 170, G: 170, B: 180, A: 255}, // Background color for the focused unselected entry
			DisabledUnselected:         color.NRGBA{R: 100, G: 100, B: 100, A: 255}, // Foreground color for the disabled unselected entry
			DisabledSelected:           color.NRGBA{R: 100, G: 100, B: 100, A: 255}, // Foreground color for the disabled selected entry
			DisabledSelectedBackground: color.NRGBA{R: 100, G: 100, B: 100, A: 255}, // Background color for the disabled selected entry
		}),
		// This required function returns the string displayed in the list
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return string(e.(ListEntry))
		}),
		// Padding for each entry
		widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
		// Text position for each entry
		widget.ListOpts.EntryTextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		// This handler defines what function to run when a list item is selected.
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			entry := args.Entry.(ListEntry)
			fmt.Println("Entry Selected: ", entry)
		}),
	)

	return &World{list: list, entries: entries}
}

func selectTarget(name string) {
	event := contract.SelectTargetEvent{
		Event: contract.Event{
			Type: contract.SelectTarget,
		},
		Name: name,
	}
	service.Controller().Send(event)
}

func (w *World) UpdateWorld(event contract.WorldEvent) {
	w.event = event
}

// Focus deprecated
func (w *World) Focus() {
}

func (w *World) SelectNext() {
	if len(w.event.Heroes) == 0 {
		return
	}
	w.entries = w.entries[0:0]
	curIndex := 0
	for i, hero := range w.event.Heroes {
		w.entries = append(w.entries, hero.Name)
		if hero.Name == w.currentTarget {
			curIndex = i
		}
	}
	index := curIndex + 1
	if index >= len(w.event.Heroes) {
		index = 0
	}
	selectTarget(w.event.Heroes[index].Name)
}

func (w *World) Draw(dst *ebiten.Image) {
	event := service.EventBus().GetWorldState()
	heroes := event.Heroes
	if len(heroes) == 0 {
		return
	}
	// sort
	m := make(map[string]contract.Hero)
	keys := make([]string, 0, len(heroes))
	curName := service.Connector().GetCurName()
	for _, hero := range heroes {
		// ignore self
		if hero.Name == curName {
			continue
		}
		m[hero.Name] = hero
		keys = append(keys, hero.Name)
	}
	sort.Strings(keys)
	// make list
	listBox := component.NewListBox(400, 200)
	for _, k := range keys {
		name := m[k].Name
		mainText := fmt.Sprintf("%s(%d)", name, m[k].Health)
		if m[k].Target != nil {
			if m[k].Target.Name == curName {
				// you are the target
				mainText += "⚔️"
			}
		}
		listBox.AppendItem(mainText)
	}
	listBox.Draw(dst)
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}
