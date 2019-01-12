package solitaire

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"

	"github.com/fyne-io/examples/solitaire/faces"
)

const minCardWidth = 95
const minPadding = 4
const cardRatio = 142.0 / minCardWidth

var (
	cardSize = fyne.Size{Width: minCardWidth, Height: minCardWidth * cardRatio}

	smallPad = minPadding
	overlap  = smallPad * 5
	bigPad   = smallPad + overlap

	minWidth  = cardSize.Width*7 + smallPad*8
	minHeight = cardSize.Height*3 + bigPad + smallPad*2
)

type cardPosition struct {
	card  *Card
	image *canvas.Image
}

func (c *cardPosition) updatePosition(x, y int) {
	c.image.Resize(cardSize)
	c.image.Move(fyne.NewPos(x, y))
}

func (c *cardPosition) withinBounds(pos fyne.Position) bool {
	if pos.X < c.image.Position().X || pos.Y < c.image.Position().Y {
		return false
	}

	if pos.X >= c.image.Position().X+c.image.Size().Width || pos.Y >= c.image.Position().Y+c.image.Size().Height {
		return false
	}

	return true
}

func newCardPos(card *Card) *cardPosition {
	if card == nil {
		return &cardPosition{nil, &canvas.Image{}}
	}

	var face fyne.Resource
	if card.FaceUp {
		face = card.Face()
	} else {
		face = faces.ForBack()
	}
	image := &canvas.Image{Resource: face}
	image.Resize(cardSize)

	return &cardPosition{card, image}
}

func newCardSpace() *cardPosition {
	space := faces.ForSpace()
	image := &canvas.Image{Resource: space}
	image.Resize(cardSize)

	return &cardPosition{nil, image}
}

type tableRender struct {
	game *Game

	deck     *cardPosition
	selected *cardPosition

	pile1, pile2, pile3            *cardPosition
	space1, space2, space3, space4 *cardPosition

	stack1, stack2, stack3, stack4, stack5, stack6, stack7 *stackRender

	objects []fyne.CanvasObject
	table   *Table
}

func updateSizes(pad int) {
	smallPad = pad
	overlap = smallPad * 5
	bigPad = smallPad + overlap
}

func (t *tableRender) selectCard(card *cardPosition) {
	if t.selected != nil {
		t.selected.image.Translucency = 0.0
		canvas.Refresh(t.selected.image)
	}

	if card != nil {
		t.selected = card
		card.image.Translucency = 0.25
		canvas.Refresh(card.image)
	}
}

func (t *tableRender) MinSize() fyne.Size {
	return fyne.NewSize(minWidth, minHeight)
}

func (t *tableRender) Layout(size fyne.Size) {
	padding := int(float32(size.Width) * .006)
	updateSizes(padding)

	newWidth := (size.Width - smallPad*8) / 7.0
	newWidth = fyne.Max(newWidth, minCardWidth)
	cardSize = fyne.NewSize(newWidth, int(float32(newWidth)*cardRatio))

	t.deck.updatePosition(smallPad, smallPad)

	t.pile1.updatePosition(smallPad*2+cardSize.Width, smallPad)
	t.pile2.updatePosition(smallPad*2+cardSize.Width+overlap, smallPad)
	t.pile3.updatePosition(smallPad*2+cardSize.Width+overlap*2, smallPad)

	t.space1.updatePosition(size.Width-(smallPad+cardSize.Width)*4, smallPad)
	t.space2.updatePosition(size.Width-(smallPad+cardSize.Width)*3, smallPad)
	t.space3.updatePosition(size.Width-(smallPad+cardSize.Width)*2, smallPad)
	t.space4.updatePosition(size.Width-(smallPad+cardSize.Width), smallPad)

	t.stack1.Layout(fyne.NewPos(smallPad, smallPad+bigPad+cardSize.Height),
		fyne.NewSize(cardSize.Width, size.Height-(smallPad+bigPad+cardSize.Height)))
	t.stack2.Layout(fyne.NewPos(smallPad+(smallPad+cardSize.Width), smallPad+bigPad+cardSize.Height),
		fyne.NewSize(cardSize.Width, size.Height-(smallPad+bigPad+cardSize.Height)))
	t.stack3.Layout(fyne.NewPos(smallPad+(smallPad+cardSize.Width)*2, smallPad+bigPad+cardSize.Height),
		fyne.NewSize(cardSize.Width, size.Height-(smallPad+bigPad+cardSize.Height)))
	t.stack4.Layout(fyne.NewPos(smallPad+(smallPad+cardSize.Width)*3, smallPad+bigPad+cardSize.Height),
		fyne.NewSize(cardSize.Width, size.Height-(smallPad+bigPad+cardSize.Height)))
	t.stack5.Layout(fyne.NewPos(smallPad+(smallPad+cardSize.Width)*4, smallPad+bigPad+cardSize.Height),
		fyne.NewSize(cardSize.Width, size.Height-(smallPad+bigPad+cardSize.Height)))
	t.stack6.Layout(fyne.NewPos(smallPad+(smallPad+cardSize.Width)*5, smallPad+bigPad+cardSize.Height),
		fyne.NewSize(cardSize.Width, size.Height-(smallPad+bigPad+cardSize.Height)))
	t.stack7.Layout(fyne.NewPos(smallPad+(smallPad+cardSize.Width)*6, smallPad+bigPad+cardSize.Height),
		fyne.NewSize(cardSize.Width, size.Height-(smallPad+bigPad+cardSize.Height)))
}

func (t *tableRender) ApplyTheme() {
	// no-op we are a custom UI
}

func (t *tableRender) BackgroundColor() color.Color {
	return color.RGBA{0x07, 0x63, 0x24, 0xff}
}

func (t *tableRender) Refresh() {
	if len(t.game.Deck.Cards) > 0 {
		t.deck.image.Resource = faces.ForBack()
	} else {
		t.deck.image.Resource = faces.ForSpace()
	}
	canvas.Refresh(t.deck.image)

	t.pile1.image.Hidden = t.pile1.card == nil
	if t.pile1.card != nil {
		t.pile1.image.Resource = t.game.Draw1.Face()
	}
	t.pile2.image.Hidden = t.pile2.card == nil
	if t.pile2.card != nil {
		t.pile2.image.Resource = t.game.Draw2.Face()
	}
	t.pile3.image.Hidden = t.pile3.card == nil
	if t.pile3.card != nil {
		t.pile3.image.Resource = t.game.Draw3.Face()
	}
	canvas.Refresh(t.pile1.image)
	canvas.Refresh(t.pile2.image)
	canvas.Refresh(t.pile3.image)

	t.stack1.Refresh(t.game.Stack1)
	t.stack2.Refresh(t.game.Stack2)
	t.stack3.Refresh(t.game.Stack3)
	t.stack4.Refresh(t.game.Stack4)
	t.stack5.Refresh(t.game.Stack5)
	t.stack6.Refresh(t.game.Stack6)
	t.stack7.Refresh(t.game.Stack7)

	canvas.Refresh(t.table)
}

func (t *tableRender) Objects() []fyne.CanvasObject {
	return t.objects
}

func (t *tableRender) appendStack(stack *stackRender) {
	for _, card := range stack.cards {
		t.objects = append(t.objects, card.image)
	}
}

func (t *tableRender) positionForCard(card *Card) *cardPosition {

	return nil
}

func newTableRender(game *Game) *tableRender {
	render := &tableRender{}
	render.game = game
	render.deck = newCardPos(nil)

	render.pile1 = newCardPos(nil)
	render.pile2 = newCardPos(nil)
	render.pile3 = newCardPos(nil)

	render.space1 = newCardSpace()
	render.space2 = newCardSpace()
	render.space3 = newCardSpace()
	render.space4 = newCardSpace()

	render.stack1 = newStackRender()
	render.stack2 = newStackRender()
	render.stack3 = newStackRender()
	render.stack4 = newStackRender()
	render.stack5 = newStackRender()
	render.stack6 = newStackRender()
	render.stack7 = newStackRender()

	render.objects = []fyne.CanvasObject{render.deck.image, render.pile1.image, render.pile2.image, render.pile3.image,
		render.space1.image, render.space2.image, render.space3.image, render.space4.image}

	render.appendStack(render.stack1)
	render.appendStack(render.stack2)
	render.appendStack(render.stack3)
	render.appendStack(render.stack4)
	render.appendStack(render.stack5)
	render.appendStack(render.stack6)
	render.appendStack(render.stack7)

	render.Refresh()
	return render
}

type stackRender struct {
	cards [13]*cardPosition
}

func (s *stackRender) Layout(pos fyne.Position, size fyne.Size) {
	top := pos.Y
	for i := range s.cards {
		s.cards[i].updatePosition(pos.X, top)

		top += overlap
	}
}

func (s *stackRender) Refresh(stack Stack) {
	var i int
	var card *Card
	for i, card = range stack.Cards {
		s.cards[i].card = card
		if card.FaceUp {
			s.cards[i].image.Resource = card.Face()
			s.cards[i].image.Show()
		} else {
			s.cards[i].image.Resource = faces.ForBack()
			s.cards[i].image.Show()
		}
	}

	for i = i + 1; i < len(s.cards); i++ {
		s.cards[i].card = nil
		s.cards[i].image.Hide()
	}
}

func newStackRender() *stackRender {
	r := &stackRender{}
	for i := range r.cards {
		r.cards[i] = newCardPos(nil)
	}

	return r
}
