// Asset Loader: utilitários de carregamento de imagens e animações
// para a fase de batalha. Mantém o acesso aos arquivos centralizado
// e converte formatos (como GIF) em ebiten.Image utilizável.
package assets

import (
	"image/gif"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	firePath = "assets/images/Fire.gif"
	missPath = "assets/images/Ponto que já foi atingido 1.png"
)

// LoadFireAnimation:
// - Abre o GIF
// - Decodifica todos os frames e seus delays
// - Converte imagens em ebiten.Image para desenhar animado
func LoadFireAnimation() ([]*ebiten.Image, []int, error) {
	f, err := os.Open(firePath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, nil, err
	}

	frames := make([]*ebiten.Image, len(g.Image))
	delays := make([]int, len(g.Delay))

	for i, img := range g.Image {
		frames[i] = ebiten.NewImageFromImage(img)
		delays[i] = g.Delay[i]
	}

	return frames, delays, nil
}

// LoadHitImage:
// - Carrega uma imagem usada como efeito de acerto
// - Caso a animação esteja indisponível, ela pode servir como fallback
func LoadHitImage() (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(firePath)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// LoadMissImage:
// - Carrega a imagem para marcar erros (miss) nos tiros
// - Usada no renderer para desenhar marcadores de jogadas
func LoadMissImage() (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(missPath)
	if err != nil {
		return nil, err
	}
	return img, nil
}
