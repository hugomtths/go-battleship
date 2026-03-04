package components

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// FireAnimation controla a animação de fogo usada quando um navio é atingido.
// Ela gerencia uma sequência de imagens (frames) e seus tempos de duração.
type FireAnimation struct {
	// frames armazena a lista de imagens que compõem a animação.
	frames []*ebiten.Image
	// delays armazena a duração de cada frame em centésimos de segundo (1/100s).
	delays []int
}

// NewFireAnimation cria uma nova instância da animação de fogo.
// Recebe os frames e os tempos de duração correspondentes.
func NewFireAnimation(frames []*ebiten.Image, delays []int) *FireAnimation {
	return &FireAnimation{
		frames: frames,
		delays: delays,
	}
}

// CurrentFrame retorna a imagem correspondente ao momento atual da animação.
// A animação é baseada no tempo do sistema e roda em loop contínuo.
func (a *FireAnimation) CurrentFrame() *ebiten.Image {
	// Se não houver frames carregados, não há nada para exibir.
	if len(a.frames) == 0 {
		return nil
	}
	// Se a lista de delays não corresponder à lista de frames, retorna o primeiro frame estático por segurança.
	if len(a.delays) != len(a.frames) {
		return a.frames[0]
	}

	// Calcula a duração total de um ciclo completo da animação em milissegundos.
	totalDuration := 0
	for _, d := range a.delays {
		totalDuration += d * 10 // Converte de centésimos (10ms) para milissegundos.
	}
	// Evita divisão por zero se a duração for inválida.
	if totalDuration == 0 {
		totalDuration = 100
	}

	// Obtém o tempo atual em milissegundos.
	now := int(time.Now().UnixMilli())
	// Calcula a posição atual dentro do ciclo de animação usando o resto da divisão.
	cycleTime := now % totalDuration

	// Percorre os frames para encontrar qual deve ser exibido no tempo atual do ciclo.
	currentDuration := 0
	for i, d := range a.delays {
		frameDuration := d * 10 // Duração deste frame específico em ms.
		if frameDuration == 0 {
			frameDuration = 100 // Valor padrão de segurança.
		}

		// Se o tempo do ciclo for menor que o tempo acumulado até o final deste frame,
		// então este é o frame atual.
		if cycleTime < currentDuration+frameDuration {
			return a.frames[i]
		}
		// Acumula a duração para verificar o próximo frame.
		currentDuration += frameDuration
	}

	// Fallback para retornar o primeiro frame caso algo falhe no loop.
	return a.frames[0]
}
