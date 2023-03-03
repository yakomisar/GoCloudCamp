package main

import (
	"fmt"
	"log"
	"musicplayer/player"
	"sync"
	"time"
)

// Интерфейс Плеера с методами, которые заданы по условию задания
// Соответственно структура Playlist будет удовлетворять данному интерфейсу
// и обладать набором необходимого функционала
type Player interface {
	Play()
	Pause()
	AddSong(s *Song)
	Next()
	Prev()
}

// Структура Song, которая содержит наименование песни и ее продолжительность
type Song struct {
	Name     string
	Duration time.Duration
}

// Структура Track, которая которая представляет собой двусвязный список
// Содержит указатель на следующую песню и предыдущую
type Track struct {
	current *Song
	next    *Track
	prev    *Track
}

// Структура Playlist, которая будет обладает набором необходимых полей
// для управления плейлистом. В качестве методов управления реализованы функции
// управления, удовлетворяющие интерфейсу Player
type Playlist struct {
	head      *Track
	tail      *Track
	isPlaying bool
	isOnPause bool
	mutex     sync.Mutex
	pause     chan struct{}
	resume    chan struct{}
	done      chan struct{}
	ticker    player.Ticker
}

// Методы управления, необходимые по условию задания
// Метод AddSong - реализует добавление песни в конец плейлиста
func (p *Playlist) AddSong(s Song) {
	newSong := &Track{current: &s}
	// Проверка на наличие песен в плейлисте
	if p.head == nil {
		p.head = newSong
		p.tail = newSong
	} else {
		currentNode := p.head
		for currentNode.next != nil {
			currentNode = currentNode.next
		}
		newSong.prev = currentNode
		currentNode.next = newSong
		p.tail = newSong
	}
	log.Printf("Song %s has been added to the end.\n", s.Name)
}

// Метод Play - реализует функционал проигрывания песни (имитация)
// Для данного метода был немного усовершенствован Ticker из стандартной
// библиотеки, так как у базового тикера нет функционала временной паузы
// и продолжение с того места на котором Ticker был поставлен на паузу
func (p *Playlist) Play() {
	// Плеер не играет и стоит на паузе
	if !p.isPlaying && p.isOnPause {
		fmt.Println("Resume")
		p.resume <- struct{}{}
		return
	}
	if p.head == nil {
		log.Println("Please add song to the playlist.")
		return
	}
	p.isPlaying = true
	p.ticker = player.New(time.Second * p.head.current.Duration)
	log.Println("Playing song:", p.head.current.Name, " duration:", p.head.current.Duration)
	for {
		select {
		case <-p.ticker.Ticks():
			p.ticker.Stop()
			p.isPlaying = false
			p.isOnPause = false
			p.Next()
		case <-p.pause:
			log.Println("Pause button clicked.")
			p.isPlaying = false
			p.isOnPause = true
			<-p.resume
			log.Println("Play button clicked.")
			p.isPlaying = true
			p.isOnPause = false
		case <-p.done:
			return
		}
	}
}

// Метод Pause - реализует функционал постановки на паузу песни (имитация)
// Здесь необходим mutex так как область данного метода является критической секцией
// доступ к которой должен быть управляемым и не вызывать проблем при многократном
// вызове метода Pause со стороны пользователя
func (p *Playlist) Pause() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.isPlaying {
		p.pause <- struct{}{}
	} else {
		log.Println("Player is on pause state.")
	}
}

// Метод Next - реализует функционал переключение на следующую песню в плейлисте (имитация)
func (p *Playlist) Next() {
	log.Println("Button \"next\" has been clicked.")
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// Если пользователь нажал Next на пустом плейлисте
	if p.head == nil {
		p.done <- struct{}{}
		p.isPlaying = false
		p.isOnPause = false
		return
	}
	// Если пользователь нажал Next, а следующей песни нет
	if p.head.next == nil {
		log.Println("End of playlist. Please add song to the playlist.")
		return
	}
	// Плеер не играет и стоит на паузе
	if !p.isPlaying && p.isOnPause {
		p.ticker.Stop()
		moveNext(p)
		p.resume <- struct{}{}
		return
	}
	// Плеер играет в данный момент
	if p.isPlaying {
		moveNext(p)
		return
	}
	// Плеер не играет и не стоит на паузе
	if !p.isPlaying && !p.isOnPause {
		moveNext(p)
	}
}

// Вспомогательный метод для переставления указателя на следующую песню
// и отражения текущей информации о проигрываемом треке
func moveNext(p *Playlist) {
	p.head = p.head.next
	p.ticker = player.New(time.Second * p.head.current.Duration)
	log.Println("Playing song:", p.head.current.Name, " duration:", p.head.current.Duration)
}

// Метод Prev - реализует функционал переключение на предыдущую песню в плейлисте (имитация)
func (p *Playlist) Prev() {
	log.Println("Button \"prev\" has been clicked.")
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// Если пользователь нажал Prev на пустом плейлисте
	if p.head == nil {
		p.done <- struct{}{}
		p.isPlaying = false
		p.isOnPause = false
		return
	}
	// Если пользователь нажал Prev, а предыдущей песни нет
	if p.head.prev == nil {
		log.Println("You are at the beginning of the playlist.")
		return
	}
	// Плеер не играет и стоит на паузе
	if !p.isPlaying && p.isOnPause {
		p.ticker.Stop()
		moveBack(p)
		p.resume <- struct{}{}
		return
	}
	// Плеер играет в данный момент
	if p.isPlaying {
		moveBack(p)
		return
	}
	// Плеер не играет и не стоит на паузе
	if !p.isPlaying && !p.isOnPause {
		log.Println("PPPPREV")
		moveBack(p)
	}
}

// Вспомогательный метод для переставления указателя на предыдущую песню
// и отражения текущей информации о проигрываемом треке
func moveBack(p *Playlist) {
	p.head = p.head.prev
	p.ticker = player.New(time.Second * p.head.current.Duration)
	log.Println("Playing song:", p.head.current.Name, " duration:", p.head.current.Duration)
}

func printSongs(p Playlist) {
	if p.head == nil {
		log.Println("Playlist is empty.")
	} else {
		for p.head != nil {
			log.Printf("Song: %s with %d\n", p.head.current.Name, p.head.current.Duration)
			p.head = p.head.next
		}
	}
}

func NewPlaylist() *Playlist {
	pl := Playlist{
		head:      nil,
		tail:      nil,
		isPlaying: false,
		isOnPause: false,
		mutex:     sync.Mutex{},
		pause:     make(chan struct{}, 1),
		resume:    make(chan struct{}, 1),
		done:      make(chan struct{}, 1),
		ticker:    nil,
	}
	return &pl
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
	playlist := NewPlaylist()
	playlist.AddSong(song1)
	playlist.AddSong(song2)
	playlist.AddSong(song3)
	playlist.pause = make(chan struct{}, 1)
	playlist.resume = make(chan struct{}, 1)
	playlist.done = make(chan struct{}, 1)

	printSongs(*playlist)
	go playlist.Play()
	time.Sleep(time.Second * 18)
}
