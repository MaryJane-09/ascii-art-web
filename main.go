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
	Theme  string
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

func themeHandler(w http.ResponseWriter, r *http.Request) {
	theme := r.FormValue("theme")
	if theme == "" {
		theme = "light"
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "theme",
		Value: theme,
		Path:  "/",
	})

	http.Redirect(w, r, "/ascii-art", http.StatusSeeOther)

}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {

	theme := "light"
	if c, err := r.Cookie("theme"); err == nil{
		theme = c.Value
	}

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodGet {
		err := tmpl.Execute(w, Data{
			Theme: theme,
		})
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
				Theme: theme,
			})
			return
		}

		if bannerName == "" {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, Data{
				Error:  "Select a banner",
				Input:  input,
				Banner: bannerName,
				Theme: theme,
			})
			return
		}

		_, err := asciiart.ValidateInput(input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, Data{
				Error: err.Error(),
				Input: input,
				Theme: theme,
			})
			return
		}

		banner, err := asciiart.BannerCheck("banners/" + bannerName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, Data{
				Error: "Invalid Banner File. Pick a valid Banner.",
				Input: input,
				Theme: theme,
			})
			return
		}

		result := asciiart.GenerateArt(input, banner)

		err = tmpl.Execute(w, Data{
			Result: result,
			Input:  input,
			Banner: bannerName,
			Theme:  theme,
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

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)
	http.HandleFunc("/theme", themeHandler)
	fmt.Println("server running at port :8080")
	http.ListenAndServe(":8080", nil)
}
