package main

import (
	"testing"
	"time"
)

// In TestAddSong, we create an empty playlist, add a song to it,
// and check if the song was added to the playlist correctly.
func TestAddSong(t *testing.T) {
	// Initialize playlist
	p := &Playlist{}

	// Add a song
	s := Song{Name: "Song1", Duration: time.Duration(3 * time.Minute)}
	p.AddSong(s)

	// Check that the song was added to the end of the playlist
	if p.head == nil || p.head.current.Name != "Song1" || p.tail.current.Name != "Song1" {
		t.Error("AddSong failed to add song to end of playlist")
	}
}

// In TestPlay, we create a playlist with one song, start playing it,
// pause it, resume it, and finally skip to the next song.
// We use channels to simulate user input.
func TestPlay(t *testing.T) {
	p := Playlist{}
	song := Song{"Test Song", 1 * time.Millisecond}
	p.AddSong(song)
	p.pause = make(chan struct{})
	p.resume = make(chan struct{})
	p.done = make(chan struct{})

	go p.Play()
	// Wait for a few milliseconds to allow Play() to start playing the song.
	time.Sleep(10 * time.Millisecond)
	p.Pause()
	// Wait for a few milliseconds to allow Play() to pause the song.
	time.Sleep(10 * time.Millisecond)
	p.Pause()
	// Wait for a few milliseconds to allow Play() to resume playing the song.
	time.Sleep(10 * time.Millisecond)
	p.Next()

	if p.isPlaying || p.isOnPause {
		t.Error("Failed to stop playing the song.")
	}

	// Make sure Play() has exited.
	p.done <- struct{}{}
}

// In TestPause, we create a playlist with one song, start playing it,
// and pause it. We check if the song is actually paused.
func TestPause(t *testing.T) {
	// Initialize playlist
	p := &Playlist{}

	// Attempt to pause when player is not playing
	p.Pause()
	if p.isOnPause {
		t.Error("Pause should not be possible when player is not playing")
	}

	// Start playing a song
	s := Song{Name: "Song1", Duration: time.Duration(3 * time.Minute)}
	p.AddSong(s)
	p.Play()

	// Pause the song
	p.Pause()
	if !p.isOnPause {
		t.Error("Pause should set isOnPause to true when called while playing")
	}
}

// In TestNext, we test the Next() method by calling it three times:
// once on an empty playlist, once on a playlist with one song, and
// once on a playlist with two songs. We check if the playlist is empty
// after each call.
func TestNext(t *testing.T) {
	p := Playlist{}
	p.Next()

	if p.head != nil || p.tail != nil {
		t.Error("Failed to stop playing the song.")
	}

	song := Song{"Test Song", 1 * time.Millisecond}
	p.AddSong(song)
	p.Next()

	if p.head != nil || p.tail != nil {
		t.Error("Failed to stop playing the song.")
	}

	p.AddSong(song)
	p.Next()

	if p.head != nil || p.tail != nil {
		t.Error("Failed to stop playing the song.")
	}
}
