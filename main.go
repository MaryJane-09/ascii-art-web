package main

import (
	"ascii-art-web/asciiart"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Data struct {
	Result     string
	Input      string
	Banner     string
	Theme      string
	TestResult []TestResult
	Error      string
}

type TestResult struct {
	Name string
	Art  string
}

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

var lastResult string
var lastInput string
var lastBanner string
var lastTestResult []TestResult

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/ascii-art", http.StatusSeeOther)

}

func exportHandler(w http.ResponseWriter, r *http.Request) {

	if lastResult == "" {
		http.Error(w, "No result to export", http.StatusBadRequest)
		return
	}

	body := []byte(lastResult)
	bodylen := len(body)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(bodylen))
	w.Header().Set("Content-Disposition", `attachment; filename="file.txt"`)

	w.Write(body)

}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {

	theme := "light"
	if c, err := r.Cookie("theme"); err == nil {
		theme = c.Value
	}

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodGet {

		fmt.Println("Loading:")
		fmt.Println("Input:", lastInput)
		fmt.Println("Banner:", lastBanner)
		fmt.Println("Result length:", len(lastResult))
		err := tmpl.Execute(w, Data{
			Input:  lastInput,
			Banner: lastBanner,
			Result: lastResult,
			TestResult: lastTestResult,
			Theme:  theme,
		})
		if err != nil {
			fmt.Println("Template error:", err)
			return
		}
	}

	if r.Method == http.MethodPost {
		action := r.FormValue("action")

		switch action {

		case "apply":
			theme = r.FormValue("theme")
			if theme == "" {
				theme = "light"
			}

			http.SetCookie(w, &http.Cookie{
				Name:  "theme",
				Value: theme,
				Path:  "/",
			})

			http.Redirect(w, r, "/ascii-art", http.StatusSeeOther)
			return

		case "generate":

			Inputtext := r.FormValue("input")
			bannerName := r.FormValue("banner")

			if Inputtext == "" {
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
					Input:  Inputtext,
					Banner: bannerName,
					Theme:  theme,
				})
				return
			}

			_, err := asciiart.ValidateInput(Inputtext)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				tmpl.Execute(w, Data{
					Error: err.Error(),
					Input: Inputtext,
					Theme: theme,
				})
				return
			}

			BannerFile, err := asciiart.BannerCheck("banners/" + bannerName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				tmpl.Execute(w, Data{
					Error: "Invalid Banner File. Pick a valid Banner.",
					Input: Inputtext,
					Theme: theme,
				})
				return
			}

			FinalResult := asciiart.GenerateArt(Inputtext, BannerFile)
			lastResult = FinalResult
			lastInput = Inputtext
			lastBanner = bannerName
			lastTestResult = nil

			fmt.Println("Saved:")
			fmt.Println("Input:", lastInput)
			fmt.Println("Banner:", lastBanner)
			fmt.Println("Result length:", len(lastResult))

			err = tmpl.Execute(w, Data{
				Result: FinalResult,
				Input:  Inputtext,
				Banner: bannerName,
				Theme:  theme,
			})
			if err != nil {
				fmt.Println("Template error:", err)
				return
			}

		case "testall":

			Inputtext := r.FormValue("input")

			if Inputtext == "" {
				w.WriteHeader(http.StatusBadRequest)
				tmpl.Execute(w, Data{
					Error: "Input a text",
					Theme: theme,
				})
				return
			}

			banners := []string{"Standard", "Shadow", "Thinkertoy"}

			var results []TestResult

			for _, name := range banners {
				banner, err := asciiart.BannerCheck("banners/" + name + ".txt")
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)

					tmpl.Execute(w, Data{
						Error: "Invalid Banner File. Pick a valid Banner.",
						Input: Inputtext,
						Theme: theme,
					})
					return
				}

				art := asciiart.GenerateArt(Inputtext, banner)

				results = append(results, TestResult{
					Name: name,
					Art:  art,
				})
			}
			lastInput = Inputtext
			lastTestResult = results
			lastResult = ""
			lastBanner = ""

			err := tmpl.Execute(w, Data{
				TestResult: results,
				Input:      Inputtext,
				Theme:      theme,
			})
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				tmpl.Execute(w, Data{
					Error: "Cannot Test all",
					Theme: theme,
				})
				return
			}

		case "export":

			http.Redirect(w, r, "/export", http.StatusSeeOther)
			return
		}
	}
	fmt.Println()
	fmt.Println(r.Method, "/ascii-art")
	fmt.Println(time.Now().Format(time.DateTime))

}

func main() {

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)
	http.HandleFunc("/export", exportHandler)
	fmt.Println("server running at port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
