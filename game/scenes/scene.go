package scenes

import (
	"github.com/allanjose001/go-battleship/game/components/basic"
	"github.com/hajimehoshi/ebiten/v2"
)

// Scene define o contrato básico que toda tela do jogo deve seguir.
type Scene interface {
	// OnEnter é chamado quando esta cena passa a ser a cena ativa.
	// Serve para inicializar estados visuais e carregar dados necessários. -> como init do flutter
	// prev pode ser nil caso seja a primeira cena do jogo.
	//passa também o tamanho da tela (anterior?)
	OnEnter(prev Scene, size basic.Size)

	// OnExit é chamado antes de trocar para outra cena.
	// Serve para limpar estados temporários (ex: cancelar drag, animações). -> dispose do flutter
	// next é a próxima cena que será ativada.
	OnExit(next Scene)

	// Update é chamado a cada frame.
	// Deve tratar input, atualizar animações e chamar a lógica do jogo.
	// Não deve desenhar nada na tela.
	//sempre o update dos componentes da scene será dado um ponto vazio
	//(ver widget)
	Update() error

	// Draw é chamado a cada frame após o Update.
	// Recebe a tela onde tudo deve ser desenhado.
	// Não deve alterar regras de jogo, apenas renderizar o estado atual.
	Draw(screen *ebiten.Image)
}

//TODO: remover isso

// SwitchTo é uma função global para trocar de cena.
// Deve ser inicializada pelo pacote game.
var SwitchTo func(Scene)
