package main

import (
	"log"
	"time"
)

type Player interface {
	Play()
	Pause()
	AddSong(s Song)
	Next()
	Prev()
}

type Song struct {
	Name     string
	Duration time.Duration
}

type Track struct {
	current Song
	next    *Track
	prev    *Track
}

type Playlist struct {
	trackList *Track
}

// Management methods
func (p *Playlist) AddSong(s Song) {
	// Check our track list, if there are any songs
	if p.trackList == nil {
		p.trackList = &Track{current: s}
		log.Printf("Song %s has been added.\n", s.Name)
		return
	} else {
		tmp := p.trackList
		for tmp.next != nil {
			tmp = tmp.next
		}
		tmp.next = &Track{current: s}
		tmp.next.prev = tmp
		log.Printf("Song %s has been added to the end.\n", s.Name)
	}
}

func printSongs(p Playlist) {
	if p.trackList == nil {
		log.Println("Playlist is empty.")
	} else {
		for p.trackList != nil {
			log.Printf("Song: %s with %d\n", p.trackList.current.Name, p.trackList.current.Duration)
			p.trackList = p.trackList.next
		}
	}
}
func main() {
	song1 := Song{
		Name:     "A-ha",
		Duration: 1,
	}
	song2 := Song{
		Name:     "ABBA",
		Duration: 2,
	}
	song3 := Song{
		Name:     "BI-2",
		Duration: 3,
	}
	playlist := Playlist{}
	playlist.AddSong(song1)
	playlist.AddSong(song2)
	playlist.AddSong(song3)

	printSongs(playlist)
}
