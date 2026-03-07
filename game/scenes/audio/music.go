package audio

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

// sampleRatePadrao
const sampleRate = 44100

// Music representa uma musica OGG
// encapsula player, stream, volume e fade
type Music struct {
	player *audio.Player // player de audio
	file   *os.File
	stream *vorbis.Stream // stream decodificada
	volume float64        // volume atual (0.0 a 1.0)
	lock   sync.Mutex     // lock para thread-safe
}

// NewMusic cria uma musica a partir de arquivo OGG
func NewMusic(ctx *audio.Context, path string) *Music {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := vorbis.DecodeWithSampleRate(sampleRate, file)
	if err != nil {
		log.Fatal(err)
	}

	loop := audio.NewInfiniteLoop(stream, stream.Length())

	player, err := ctx.NewPlayer(loop)
	if err != nil {
		log.Fatal(err)
	}

	return &Music{
		file:   file,
		stream: stream,
		player: player,
		volume: 1.0,
	}
}

// Play inicia ou reinicia a musica
// loop=true repete a musica
func (m *Music) Play() {
	m.lock.Lock()
	defer m.lock.Unlock()

	err := m.player.Rewind()

	if err != nil {
		log.Fatal(err)
	}
	m.player.Play()
}

// Stop pausa e reinicia a musica
func (m *Music) Stop() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.player.Pause()
	err := m.player.Rewind()
	if err != nil {
		log.Fatal(err)
	}
}

// SetVolume define o volume da musica (0.0 a 1.0)
func (m *Music) SetVolume(vol float64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.volume = vol
	m.player.SetVolume(vol)
}

// FadeTo faz fade do volume atual para target em duration
// roda em goroutine, não precisa de Update contínuo
func (m *Music) FadeTo(target float64, duration time.Duration) {
	go func() {
		m.lock.Lock()
		startVol := m.volume
		m.lock.Unlock()

		steps := 20
		interval := duration / time.Duration(steps)
		stepVol := (target - startVol) / float64(steps)

		for i := 1; i <= steps; i++ {
			time.Sleep(interval)
			m.lock.Lock()
			m.volume = startVol + stepVol*float64(i)
			if m.volume < 0 {
				m.volume = 0
			} else if m.volume > 1 {
				m.volume = 1
			}
			m.player.SetVolume(m.volume)
			m.lock.Unlock()
		}
	}()
}

func (m *Music) Close() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	err := m.player.Close()
	if err != nil {
		return err
	}
	err2 := m.file.Close()
	if err2 != nil {
		return err2
	}
	return nil
}
