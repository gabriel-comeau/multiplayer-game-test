// TextureManager takes care of texture loading so that a texture file
// only needs to be loaded once from disk.
package texturemanager

import (
	"errors"

	sf "bitbucket.org/krepa098/gosfml2"

	"github.com/gabriel-comeau/multiplayer-game-test/shared"
)

var (
	// Hold on to the loaded textures
	textures map[string]*sf.Texture
)

func init() {
	textures = make(map[string]*sf.Texture)
}

// Loads up a texture.  It will attempt a keylookup if not provided a path.
func LoadTexture(key, path string) (*sf.Texture, error) {

	tex, ok := textures[key]
	if ok {
		return tex, nil
	} else {
		if path == "" {
			return nil, errors.New("No texture loaded for key " + key + ".  Provide a path to load the texture file with!")
		}
	}

	tex, err := sf.NewTextureFromFile(shared.TEXTURE_ROOT+path, nil)
	if err != nil {
		return nil, err
	}

	textures[key] = tex
	return tex, nil
}
