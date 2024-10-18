package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "text/template"
)

type GameState struct {
    Word           string   `json:"word"`
    Display        string   `json:"display"`
    Attempts       int      `json:"attempts"`
    LettersGuessed []string `json:"letters_guessed"`
    Won            bool     `json:"won"`
    Lost           bool     `json:"lost"`
}

var game GameState
var mu sync.Mutex

const port = ":8080"

func startGame() {
    game = GameState{
        Word:           "golang",
        Display:        "______",
        Attempts:       6,
        LettersGuessed: []string{},
        Won:            false,
        Lost:           false,
    }
}

func Home(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "Home Page")
}

func About(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "About Page")
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
    t, err := template.ParseFiles("./templates/" + tmpl + ".page.tmpl")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, nil)
}

func checkLetter(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()

    letter := r.URL.Query().Get("letter")
    if letter == "" || len(letter) != 1 {
        http.Error(w, "Invalid letter", http.StatusBadRequest)
        return
    }

    for _, l := range game.LettersGuessed {
        if l == letter {
            json.NewEncoder(w).Encode(game)
            return
        }
    }

    game.LettersGuessed = append(game.LettersGuessed, letter)

    correct := false
    newDisplay := []rune(game.Display)
    for i, c := range game.Word {
        if string(c) == letter {
            newDisplay[i] = rune(c)
            correct = true
        }
    }
    game.Display = string(newDisplay)

    if !correct {
        game.Attempts--
    }

    if game.Display == game.Word {
        game.Won = true
    } else if game.Attempts <= 0 {
        game.Lost = true
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(game)
}

func main() {
    startGame()
    http.HandleFunc("/", Home)
    http.HandleFunc("/about", About)
    http.HandleFunc("/guess", checkLetter)

    fmt.Println("Server started at http://localhost" + port)
    http.ListenAndServe(port, nil)
}
