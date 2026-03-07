package audio

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// SoundService gerencia todas as musicas //fade, play, stop
type SoundService struct {
	ctx     *audio.Context    //audio context compartilhado
	musics  map[string]*Music //mapa nome->music
	lock    sync.Mutex        //thread safe
	current *Music            //musica atual para fade
}

const fadeDuration = 1200 * time.Millisecond

// NewSoundService cria o servico de audio
func NewSoundService() *SoundService {
	return &SoundService{
		ctx:    audio.NewContext(sampleRate),
		musics: make(map[string]*Music),
	}
}

// LoadMusic carrega musica para o servico //name=identificador
func (ss *SoundService) LoadMusic(name, path string) {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	ss.musics[name] = NewMusic(ss.ctx, path)
}

// Play toca musica //loop=true repete //fade entre musicas
func (ss *SoundService) Play(name string) {
	if ss.current != nil && ss.current != ss.musics[name] {
		// fade out da música atual
		ss.current.FadeTo(0, fadeDuration)
	}

	newMusic := ss.musics[name]
	if newMusic != nil {
		// fade in da nova música
		newMusic.volume = 0
		newMusic.Play()
		newMusic.FadeTo(1, fadeDuration)
		ss.current = newMusic
	}
}

// StopCurrent faz fade out da musica atual
func (ss *SoundService) StopCurrent() {
	if ss.current != nil {
		ss.current.FadeTo(0, 500*time.Millisecond)
	}
}

// getMusic retorna ponteiro da musica pelo nome
func (ss *SoundService) GetMusic(name string) *Music {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	return ss.musics[name]
}

// fecha todas as músicas
func (ss *SoundService) CloseAll() error {
	for _, m := range ss.musics {
		err := m.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
