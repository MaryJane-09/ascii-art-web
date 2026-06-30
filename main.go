package main

import (
	"ascii-art-web/asciiart"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Data struct {
	Result string
	Input  string
	Banner string
	Error  string
}

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/ascii-art", http.StatusSeeOther)

}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodGet {
		err := tmpl.Execute(w, Data{})
		if err != nil {
			fmt.Println("Template error:", err)
			return
		}
	}

	if r.Method == http.MethodPost {
		input := r.FormValue("input")
		bannerName := r.FormValue("banner")

		if input == "" {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, Data{
				Error: "Input a text",
			})
			return
		}

		if bannerName == "" {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, Data{
				Error: "Select a banner",
				Input: input,
				Banner: bannerName,
			})
			return
		}

		_, err := asciiart.ValidateInput(input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, Data{
				Error: err.Error(),
				Input: input,
			})
			return
		}

		banner, err := asciiart.BannerCheck("banners/" + bannerName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, Data{
				Error: "Invalid Banner File. Pick a valid Banner.",
				Input: input,
			})
			return
		}

		result := asciiart.GenerateArt(input, banner)

		err = tmpl.Execute(w, Data{
			Result: result,
			Input:  input,
			Banner: bannerName,
		})
		if err != nil {
			fmt.Println("Template error:", err)
			return
		}

	}
	fmt.Println(r.Method, "/ascii-art")
	fmt.Println(time.DateTime)

}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)
	fmt.Println("server running at port :8080")
	http.ListenAndServe(":8080", nil)
}
