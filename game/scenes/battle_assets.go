package scenes

import (
	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/internal/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

// BattleAssets contém os recursos gráficos e animações usados na cena de batalha.
// Isso evita carregar imagens repetidamente a cada frame.
type BattleAssets struct {
	// FireAnimation é a animação de fogo usada quando um navio é atingido.
	FireAnimation *components.FireAnimation
	// HitImage é a imagem estática exibida em um acerto (X vermelho).
	HitImage *ebiten.Image
	// MissImage é a imagem estática exibida em um erro (círculo na água).
	MissImage *ebiten.Image
}

// LoadBattleAssets carrega todos os assets necessários para a batalha de uma vez.
// Retorna uma struct com ponteiros para as imagens e animações prontas.
func LoadBattleAssets() *BattleAssets {
	// Carrega os sprites da animação de fogo e seus tempos de duração
	frames, delays, _ := assets.LoadFireAnimation()
	// Carrega a imagem de acerto (X)
	hit, _ := assets.LoadHitImage()
	// Carrega a imagem de erro (água)
	miss, _ := assets.LoadMissImage()

	// Se a imagem de hit falhar, usa o primeiro frame do fogo como fallback
	if hit == nil && len(frames) > 0 {
		hit = frames[0]
	}

	// Cria o componente de animação de fogo se houver frames carregados
	var fireAnim *components.FireAnimation
	if len(frames) > 0 {
		fireAnim = components.NewFireAnimation(frames, delays)
	}

	return &BattleAssets{
		FireAnimation: fireAnim,
		HitImage:      hit,
		MissImage:     miss,
	}
}
