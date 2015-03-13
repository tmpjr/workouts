package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type Exercise struct {
	Name string `json:"name"`
	Sets int    `json:"sets"`
	Reps int    `json:"reps"`
	Rest int    `json:"rest"`
}

type WorkoutType struct {
	Title     string `json:"title"`
	Exercises []Exercise
}

type Workout struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Workouts    []WorkoutType
}

var HTML string

func main() {
	var path string
	flag.StringVar(&path, "path", "./path.json", "Path to workout json file")
	flag.Parse()

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	var workout *Workout
	decoder.Decode(&workout)

	var workouts []WorkoutType = workout.Workouts
	var exercises []Exercise
	HTML = `<!DOCTYPE html>
		<html>
		  <head>
		    <title>My Workout</title>
		  </head>
		  <body>
		    <h2>%s</h2>
			<h4>%s</h4>
		    <table border="1" cellpadding="2" cellspacing="2">
			  %s
			</table>
		  </body>
		</html>`

	var tbl string
	var table string
	var row string
	for _, wo := range workouts {
		exercises = wo.Exercises
		tbl = "<tr><th style=\"text-align: left;\">" + wo.Title + "</th></tr>"
		for _, exercise := range exercises {
			row = "<tr><td>%s / %d / %d / %d</td></tr>"
			row = fmt.Sprintf(row, exercise.Name, exercise.Sets, exercise.Reps, exercise.Rest)
			tbl += row
		}
		table += tbl
	}
	HTML = fmt.Sprintf(HTML, workout.Title, workout.Description, table)

	r := mux.NewRouter()
	r.HandleFunc("/", handleIndex)
	http.ListenAndServe(":8000", r)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(HTML))
}
