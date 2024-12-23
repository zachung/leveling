package ui

import (
	"fmt"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"image/color"
	"leveling/internal/client/contract"
	"leveling/internal/client/service"
)

func layoutRoot() *widget.Container {
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 0, 0})),
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			//Define number of columns in the grid
			widget.GridLayoutOpts.Columns(1),
			//Define how much padding to inset the child content
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			//Define how far apart the rows and columns should be
			widget.GridLayoutOpts.Spacing(20, 10),
			//Define how to stretch the rows and columns.
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false, false}),
		)),
	)

	headContainer := newHeadContainer()

	bodyContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(50, 50),
		),
	)

	chatContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 0, 0})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxHeight: 150,
			}),
			widget.WidgetOpts.MinSize(50, 200),
		),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			//Define number of columns in the grid
			widget.GridLayoutOpts.Columns(2),
			//Define how to stretch the rows and columns.
			widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{true}),
		)),
	)
	chatTextarea := textContainer()
	bus := service.EventBus()
	bus.AddObserver(contract.OnReportAppend, func() {
		chatTextarea.SetText(bus.GetReport())
	})

	chatContainer.AddChild(chatTextarea)

	footerContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxHeight: 20,
			}),
			widget.WidgetOpts.MinSize(20, 20),
		),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			//Define number of columns in the grid
			widget.GridLayoutOpts.Columns(2),
			//Define how to stretch the rows and columns.
			widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{true}),
		)),
	)

	rootContainer.AddChild(
		headContainer,
		bodyContainer,
		chatContainer,
		footerContainer,
	)

	return rootContainer
}

func newHeadContainer() *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxWidth: 200,
			}),
			widget.WidgetOpts.MinSize(50, 50),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(5)),
			widget.RowLayoutOpts.Spacing(5),
		)),
	)

	labelHealth := widget.NewText(
		widget.TextOpts.Text("Name: Health will be here", contract.UiTextFace, color.White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
	)
	container.AddChild(labelHealth)

	labelTarget := widget.NewText(
		widget.TextOpts.Text("Target: Health will be here", contract.UiTextFace, color.White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
	)
	container.AddChild(labelTarget)

	bus := service.EventBus()
	bus.AddObserver(contract.OnStateChanged, func() {
		event := bus.GetState()
		labelHealth.Label = fmt.Sprintf("%v: %d", event.Name, event.Health)
		if event.Target.Name != "" {
			labelTarget.Label = fmt.Sprintf("%v: %d", event.Target.Name, event.Target.Health)
		} else {
			labelTarget.Label = ""
		}
	})

	return container
}

type TextView struct {
	textarea *widget.TextArea
}

func textContainer() *widget.TextArea {
	textarea := widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				//Set the layout data for the textarea
				//including a max height to ensure the scroll bar is visible
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
				//Set the minimum size for the widget
				widget.WidgetOpts.MinSize(300, 100),
			),
		),
		//Set gap between scrollbar and text
		widget.TextAreaOpts.ControlWidgetSpacing(2),
		//Tell the textarea to display bbcodes
		widget.TextAreaOpts.ProcessBBCode(true),
		//Set the font color
		widget.TextAreaOpts.FontColor(color.Black),
		//Set the font face (size) to use
		widget.TextAreaOpts.FontFace(contract.UiTextFace),
		//Set the initial text for the textarea
		//It will automatically line wrap and process newlines characters
		//If ProcessBBCode is true it will parse out bbcode
		//widget.TextAreaOpts.Text(service.EventBus().GetReport()),
		//Tell the TextArea to show the vertical scrollbar
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(widget.NewInsetsSimple(10)),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Mask: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}),
		),
		//This sets the images to use for the sliders
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(
				// Set the track images
				&widget.SliderTrackImage{
					Idle:  image.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
					Hover: image.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
				},
				// Set the handle images
				&widget.ButtonImage{
					Idle:    image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
					Hover:   image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
					Pressed: image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				},
			),
		),
	)

	return textarea
}
