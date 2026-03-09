package scenes

import (
	"sync"

	"github.com/allanjose001/go-battleship/game/components"
	"github.com/allanjose001/go-battleship/internal/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	cachedBattleAssets *BattleAssets
	loadBattleOnce     sync.Once
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

	SunkShip1 *ebiten.Image
	SunkShip2 *ebiten.Image
	SunkShip3 *ebiten.Image
	SunkShip4 *ebiten.Image
}

// LoadBattleAssets carrega todos os assets necessários para a batalha de uma vez.
// Retorna uma struct com ponteiros para as imagens e animações prontas.
// Agora possui cache para que a decodificação/carregamento ocorra apenas 1 vez.
func LoadBattleAssets() *BattleAssets {
	loadBattleOnce.Do(func() {
		frames, delays, _ := assets.LoadFireAnimation()
		hit, _ := assets.LoadHitImage()
		miss, _ := assets.LoadMissImage()
		sunk1, _ := assets.LoadSunkShip1()
		sunk2, _ := assets.LoadSunkShip2()
		sunk3, _ := assets.LoadSunkShip3()
		sunk4, _ := assets.LoadSunkShip4()

		if hit == nil && len(frames) > 0 {
			hit = frames[0]
		}

		var fireAnim *components.FireAnimation
		if len(frames) > 0 {
			fireAnim = components.NewFireAnimation(frames, delays)
		}

		cachedBattleAssets = &BattleAssets{
			FireAnimation: fireAnim,
			HitImage:      hit,
			MissImage:     miss,
			SunkShip1:     sunk1,
			SunkShip2:     sunk2,
			SunkShip3:     sunk3,
			SunkShip4:     sunk4,
		}
	})

	return cachedBattleAssets
}
