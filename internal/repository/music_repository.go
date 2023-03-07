package repository

import (
	"bufio"
	"io"
	"log"
	"netradio/internal/model"
	"netradio/pkg/errors"
	"netradio/pkg/generics/slices"
	"os"
)

const (
	chunkSize = 1 << 20 // 512kb
)

type MusicDB interface {
	GetPodcasts() []model.MusicInfo
	LoadMusicBatch(info model.MusicInfo) (<-chan []byte, error)
}

func NewMusicDB() MusicDB {
	return &MusicDBImpl{}
}

type MusicDBImpl struct{}

func (s MusicDBImpl) GetPodcasts() []model.MusicInfo {
	return nil
}

func (s MusicDBImpl) LoadMusicBatch(info model.MusicInfo) (<-chan []byte, error) {
	file, err := os.Open("./music.mp3")
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	r := bufio.NewReader(file)
	chunks := make(chan []byte)
	go func() {
		defer close(chunks)

		buf := make([]byte, 0, chunkSize)
		for {
			n, err := r.Read(buf[:cap(buf)])
			buf = buf[:n]
			if n == 0 {
				if err == nil {
					continue
				}
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}

			chunks <- slices.Copy(buf)
		}
	}()
	return chunks, nil
}
