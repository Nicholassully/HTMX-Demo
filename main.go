package main

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

var (
	items = []string{"Item 1", "Item 2", "Item 3"}
	mu    sync.Mutex
	tpl   = template.Must(template.ParseFiles("templates/index.html", "templates/items.html", "templates/other.html", "templates/error.html"))
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/items", itemsHandler)
	http.HandleFunc("/add_item", addItemHandler)
	http.HandleFunc("/delete_item", deleteItemHandler)
	http.HandleFunc("/other", otherHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	tpl.ExecuteTemplate(w, "items.html", items)
}

func addItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		item := r.FormValue("item")
		if item != "" {
			mu.Lock()
			items = append(items, item)
			mu.Unlock()
			itemsHandler(w, r)
		} else {
			http.Error(w, "No item provided", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		indexStr := r.FormValue("index")
		index, err := strconv.Atoi(indexStr)
		if err != nil || index < 0 || index >= len(items) {
			http.Error(w, "Invalid index", http.StatusBadRequest)
			return
		}
		mu.Lock()
		items = append(items[:index], items[index+1:]...)
		mu.Unlock()
		itemsHandler(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func otherHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "other.html", "This is some other content that does not change.")
}
