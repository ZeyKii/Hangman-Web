package main

import (
	"hangmanweb/hangman_classic"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type hang struct { // Structure contenant toute les variables pour le Hangman
	Word       string   // le mot aléatoire
	Letter     string   // La Lettre qui est entrée
	Wordtofind string   // Le mot qui est caché
	Attempts   int      // Le nombre d'essai restant
	LetterUsed []string // Les lettres utilisées
	Mode       string   // Mode de jeu (Difficultées)
	Hang_img   string   // Etapes du pendu (images)
}

var data hang // Stockage de la structure dans une variable

func main() {
	http.HandleFunc("/", Main_Page)                           // Appel de la fonction pour la page d'entrée
	http.HandleFunc("/reset", resetHandler)                   // Appel de la fonction qui pour Reset le mot (ramène à la page d'accueil)
	http.HandleFunc("/hangman", Dataprocess)                  // Appel de fonction Dataprocess gère les variables de la structures
	http.HandleFunc("/youlose", Lose)                         // Redirection a la page Lose
	http.HandleFunc("/youwin", Win)                           // Redirection a la page Win
	http.HandleFunc("/game", Handler)                         // Page du jeu
	http.HandleFunc("/gotogame", GoToGame)                    // Fonction de redirection vers le Jeu
	fs := http.FileServer(http.Dir("./static"))               // Chargement des assets dans une variables
	http.Handle("/static/", http.StripPrefix("/static/", fs)) // Chargement des assets dans une variables
	http.ListenAndServe(":8081", nil)                         // Selection du port sur lequel on ouvre le serveur
}
func Main_Page(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("./static/main_page.html"))
	if data.Word == "" { // Génère un mot si le word est vide
		switch r.Method {
		case "POST":
			value := r.FormValue("difficulty") // Récupération de la difficulté pour le jeu
			data.Mode = value
			var lst_mot []string
			switch data.Mode { // Switch qui permet de charger les mots +/- difficile
			case "Easy":
				lst_mot = hangman_classic.Scan("./assets/word_easy.txt") // Scan d'une banque de mot
			case "Medium":
				lst_mot = hangman_classic.Scan("./assets/word_medium.txt")
			case "Hard":
				lst_mot = hangman_classic.Scan("./assets/word_hard.txt")
			}
			rand.Seed(time.Now().UnixNano())
			println(len(lst_mot))
			randomword := lst_mot[rand.Intn(len(lst_mot))] // Choisis un mot aléatoire dans la liste des mots scannez
			findword := hangman_classic.CreateWord((randomword))
			data.Word = randomword                                     // Mot aléatoire
			data.Wordtofind = findword                                 // Mot aléatoire caché
			data.Attempts = 10                                         // Essai initialiser à 10
			data.Hang_img = "/static/assets/images/hangman/0_step.png" // charge la 1er image du pendu (vide)
			http.Redirect(w, r, "/gotogame", http.StatusFound)         // Redirection vers le jeu
		}
	}
	tpl.Execute(w, nil) // Execute tout ce que l'on modifie dans le main
}
func GoToGame(w http.ResponseWriter, r *http.Request) { // Fonction de redirection vers le Jeu
	data.LetterUsed = []string{}
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}
func Handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/index.html")) // Chargement de la page du jeu
	tmpl.Execute(w, data)
}
func Dataprocess(w http.ResponseWriter, r *http.Request) {
	var err string                                                                                           // Initialisation de la varibale err
	variable := r.FormValue("input")                                                                         // Récupération du formulaire (une lettre ou un mot par exemple)
	data.Wordtofind, err = hangman_classic.IsInputOk(variable, data.Word, data.Wordtofind, &data.LetterUsed) // Vérification de correspondance des lettres entre l'Input et le Word
	if err == "fail" || err == "error" || err == "wordwrong" || err == "wordinvalid" {                       // Gestion des erreur
		data.Hang_img = "/static/assets/images/hangman/Step_" + strconv.Itoa(data.Attempts) + ".png"
		data.Attempts -= 1 // Décrémentation des Attemps en cas d'erreur
	}
	if data.Attempts <= 0 { // Redirection en cas de Attemps = 0
		http.Redirect(w, r, "/youlose", http.StatusSeeOther)
	}
	if data.Wordtofind == data.Word || variable == data.Word { // Redirection en cas de victoire
		http.Redirect(w, r, "/youwin", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/game", http.StatusSeeOther) // Redirection vers le jeu si aucune autre redirection n'est éffectué (win ou lose)
}
func resetHandler(w http.ResponseWriter, r *http.Request) { // Fonction du reset qui vide le mot
	data.Word = ""               // On vide la table afin de pouvoire regénérer le mot
	data.LetterUsed = []string{} // On vide aussi les lettres utilisées
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func Win(w http.ResponseWriter, r *http.Request) { // Fonction de redirection en cas de win
	tpl := template.Must(template.ParseFiles("./static/you_win.html"))
	tpl.Execute(w, data)
}
func Lose(w http.ResponseWriter, r *http.Request) { // Fonction de redirection en cas de lose
	tpl := template.Must(template.ParseFiles("./static/you_lose.html"))
	tpl.Execute(w, data)
}
