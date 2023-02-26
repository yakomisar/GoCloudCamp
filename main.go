package main

import (
	"fmt"
	"log"
	"musicplayer/player"
	"sync"
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
	isPlaying bool
	isOnPause bool
	mutex     sync.Mutex
	pause     chan struct{}
	resume    chan struct{}
	done      chan struct{}
}

// Management methods
// Add song to the playlist
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

// Play method allows user to listen to songs one by one from the tracklist
func (p *Playlist) Play() {
	if !p.isPlaying && p.isOnPause {
		fmt.Println("Resume")
		p.resume <- struct{}{}
	} else {
		if p.trackList != nil {
			p.isPlaying = true
			//ticker := time.NewTicker(time.Second * p.trackList.current.Duration)
			ticker := player.New(time.Second * p.trackList.current.Duration)
			log.Println("Playing song:", p.trackList.current.Name, " duration:", p.trackList.current.Duration)
			for {
				select {
				case <-ticker.Ticks():
					fmt.Println("here here")
					ticker.Stop()
					go p.Next()
					return
				case <-p.pause:
					//ticker.Stop()
					fmt.Println("Pause button clicked.")
					p.isPlaying = false
					p.isOnPause = true
					<-p.resume
					fmt.Println("Play button clicked.")
					p.isPlaying = true
					p.isOnPause = false
				case <-p.done:
					ticker.Stop()
					return
				}
			}
		} else {
			log.Println("Please add song to the playlist.")
		}
	}
}

func (p *Playlist) Pause() {
	p.mutex.Lock()
	if p.isPlaying {
		p.pause <- struct{}{}
	} else {
		log.Println("Player is on pause state.")
	}
	fmt.Println("Trying to pause")
	p.mutex.Unlock()
}

// Next method allows user to play the next one song in the list if applicable
func (p *Playlist) Next() {
	log.Println("Button \"next\" has been clicked.")
	//p.mutex.Lock()
	//// Плеер не играет и стоит на паузе
	//if !p.isPlaying && p.isOnPause {
	//
	//}
	//// Плеер играет
	//
	//// Плеер не играет и не стоит на паузе
	//p.mutex.Unlock()
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
		Duration: 2,
	}
	song2 := Song{
		Name:     "ABBA",
		Duration: 3,
	}
	song3 := Song{
		Name:     "BI-2",
		Duration: 2,
	}
	playlist := Playlist{}
	playlist.AddSong(song1)
	playlist.AddSong(song2)
	playlist.AddSong(song3)
	playlist.pause = make(chan struct{}, 1)
	playlist.resume = make(chan struct{}, 1)
	playlist.done = make(chan struct{}, 1)

	printSongs(playlist)
	go playlist.Play()
	time.Sleep(time.Second * 1)
	playlist.Pause()
	//playlist.Pause()
	//playlist.Pause()
	time.Sleep(time.Second * 7)
	go playlist.Play()
	//go playlist.Pause()
	//time.Sleep(time.Second * 2)
	//go playlist.Play()
	time.Sleep(time.Second * 2)
}
