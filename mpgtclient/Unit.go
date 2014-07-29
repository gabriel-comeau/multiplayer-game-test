package main

import (
	sf "bitbucket.org/krepa098/gosfml2"

	"github.com/gabriel-comeau/multiplayer-game-test/texturemanager"
)

// Represents a drawable unit.
type Unit struct {
	tex    *sf.Texture
	sprite *sf.Sprite
}

// Draw the unit to the render target (the window)
func (this *Unit) Draw(target sf.RenderTarget, states sf.RenderStates) {
	this.sprite.Draw(target, states)
}

// Move unit from its current position to a new one via a vector offset
func (this *Unit) Move(offset sf.Vector2f) {
	this.sprite.Move(offset)
}

// Set the unit to a new absolute position (regardless of its previous position)
func (this *Unit) SetPosition(pos sf.Vector2f) {
	this.sprite.SetPosition(pos)
}

// Get the unit's current position
func (this *Unit) GetPosition() sf.Vector2f {
	return this.sprite.GetPosition()
}

// This constructor will set up the unit with the Player texture
func NewPlayer(initialPos sf.Vector2f) *Unit {
	player := new(Unit)

	var err error
	player.tex, err = texturemanager.LoadTexture("sprites-player", "redsquare.png")
	if err != nil {
		return nil
	}

	spr, err := sf.NewSprite(player.tex)
	if err != nil {
		return nil
	}

	player.sprite = spr
	player.sprite.SetPosition(initialPos)

	return player
}

// This constructor will set up the unit with the Other texture
func NewOther(initialPos sf.Vector2f) *Unit {
	other := new(Unit)

	var err error
	other.tex, err = texturemanager.LoadTexture("sprites-other", "bluesquare.png")
	if err != nil {
		return nil
	}

	spr, err := sf.NewSprite(other.tex)
	if err != nil {
		return nil
	}

	other.sprite = spr
	other.sprite.SetPosition(initialPos)

	return other
}
