package game

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Card represents a playing card using short code like "AS" for Ace of Spades.
type Card string

var assetDir string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	assetDir = filepath.Join(filepath.Dir(filename), "..", "..", "web", "assets")
}

// Deck returns a new ordered deck of 52 cards.
func Deck() []Card {
	suits := []string{"D", "H", "S", "C"}
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	deck := make([]Card, 0, 52)
	for _, v := range values {
		for _, s := range suits {
			deck = append(deck, Card(fmt.Sprintf("%s%s", v, s)))
		}
	}
	return deck
}

// Shuffle shuffles the given deck in place.
func Shuffle(deck []Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
}

// AssetPath returns the expected asset path for the card.
func (c Card) AssetPath() string {
	return filepath.Join(assetDir, fmt.Sprintf("%s.png", string(c)))
}

// VerifyAssets ensures all card assets exist. It fatally logs if any are missing.
func VerifyAssets() {
	for _, card := range Deck() {
		if _, err := os.Stat(card.AssetPath()); err != nil {
			log.Fatalf("missing asset for card %s: %v", card, err)
		}
	}
	// also check backs
	backs := []string{"blue", "green", "gray", "purple", "red", "yellow"}
	for _, b := range backs {
		path := filepath.Join(assetDir, fmt.Sprintf("%s_back.png", b))
		if _, err := os.Stat(path); err != nil {
			log.Fatalf("missing asset %s_back.png: %v", b, err)
		}
	}
}
