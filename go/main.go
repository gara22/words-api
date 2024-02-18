package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func generateRandom(max int32) int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int31n(max)
}

func readLines(path string) (*[]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return &lines, scanner.Err()
}

func getRandomWord() string {
	words, err := readLines("../words.txt")
	var words_length = len(*words) - 1
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	var random = generateRandom(int32(words_length))
	var randomWord = strings.ToLower((*words)[random])
	return randomWord
}

type WordResponse struct {
	RandomWord string `json:"randomWord"`
}

func publisher(nc *nats.Conn, word string) {
	err := nc.Publish("messages", []byte(word))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Published message: %s\n", word)
}

func subscriber(nc *nats.Conn, id int) {
	_, err := nc.Subscribe("messages", func(msg *nats.Msg) {
		fmt.Printf("Received message, subscriber: %v %s\n", id, string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	// Keep the subscriber running
	select {}
}

func main() {
	// Connect to NATS server running locally
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	fmt.Println("Connected to NATS Server")

	// Start publisher goroutine

	// Start subscriber goroutine

	
	http.HandleFunc("/word", func(w http.ResponseWriter, r *http.Request) {
		word := getRandomWord()
		asd := WordResponse{
			RandomWord: word,
		}
		publisher(nc, word)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		// fmt.Println(word)
		json.NewEncoder(w).Encode(asd)
		
		// fmt.Fprintf(w, "Hello, this is asd: %s\n", r.URL.Path)
	})
	for i := 0; i < 10; i++ {
		go subscriber(nc, i)
	}

	http.ListenAndServe(":3100", nil)

	// Keep the main function running
	// select {}
}
