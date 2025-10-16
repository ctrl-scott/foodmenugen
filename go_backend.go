// go mod init menuapi && go mod tidy
package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type Menu struct {
	ID         string `json:"id"`
	Restaurant struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		Phone    string `json:"phone"`
		Website  string `json:"website"`
		Currency string `json:"currency"`
		Note     string `json:"note"`
	} `json:"restaurant"`
	Settings struct {
		PriceDecimals   int  `json:"priceDecimals"`
		ShowTags        bool `json:"showTags"`
		ShowDescriptions bool `json:"showDescriptions"`
	} `json:"settings"`
	Sections []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Note  string `json:"note"`
		Items []struct {
			ID          string   `json:"id"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Price       float64  `json:"price"`
			Tags        []string `json:"tags"`
			Available   bool     `json:"available"`
		} `json:"items"`
	} `json:"sections"`
}

var (
	mem  = map[string]Menu{}
	lock sync.RWMutex
)

func main() {
	http.HandleFunc("/api/menu", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		var m Menu
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		lock.Lock()
		mem[m.ID] = m
		lock.Unlock()
		w.WriteHeader(http.StatusNoContent)
	})

	type renderReq struct{ ID string `json:"id"` }

	tpl := template.Must(template.New("menu").Parse(`
<!doctype html>
<html><head>
<meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.Restaurant.Name}}</title>
<style>{{.CSS}}</style>
</head><body>
<main class="right"><section class="menu" id="preview">
<header>
  <h1 style="font-size:28px; margin:0;">{{.Restaurant.Name}}</h1>
  <div class="meta">{{.Restaurant.Address}}{{if .Restaurant.Phone}} · {{.Restaurant.Phone}}{{end}}{{if .Restaurant.Website}} · {{.Restaurant.Website}}{{end}}</div>
  {{if .Restaurant.Note}}<div class="meta">{{.Restaurant.Note}}</div>{{end}}
</header>
{{range .Sections}}
  <h2>{{.Title}}</h2>
  {{if .Note}}<div class="meta">{{.Note}}</div>{{end}}
  {{range .Items}}
    {{if or (not .Available) (eq .Available true)}}{{/* treat zero as available */}}
    <div class="item">
      <div>
        <div class="name">{{.Name}}</div>
        {{if $.Settings.ShowDescriptions}}{{if .Description}}<div class="desc">{{.Description}}</div>{{end}}{{end}}
        {{if $.Settings.ShowTags}}{{if .Tags}}<div class="tags">{{range .Tags}}<span class="tag">{{.}}</span>{{end}}</div>{{end}}{{end}}
      </div>
      <div class="price">{{printf "%.*f" $.Settings.PriceDecimals .Price}}</div>
    </div>
    {{end}}
  {{end}}
{{end}}
</section></main>
</body></html>
`))

	// Copy/paste the <style> block from the canvas file into cssStr (keeps look consistent)
	cssStr := `:root { --bg:#fff; --fg:#111; --muted:#666; --acc:#0a7; --border:#e5e5e5;}
.right{padding:16px}.menu{max-width:900px;margin:0 auto}.menu h2{margin:24px 0 8px;border-bottom:2px solid #000;padding-bottom:6px;font-size:20px}
.meta{color:#666;font-size:12px;margin-bottom:16px}.item{display:grid;grid-template-columns:1fr auto;gap:8px;padding:8px 0;border-bottom:1px dashed #e5e5e5}
.item:last-child{border-bottom:0}.name{font-weight:600}.desc{color:#666;font-size:13px}.price{font-variant-numeric:tabular-nums;font-weight:600}
.tag{border:1px solid #e5e5e5;border-radius:999px;font-size:10px;padding:2px 6px;color:#666}`

	http.HandleFunc("/api/render", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		var rr renderReq
		if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		lock.RLock()
		m, ok := mem[rr.ID]
		lock.RUnlock()
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		// supply CSS into template data
		type tData struct {
			Menu
			CSS string
		}
		data := tData{Menu: m, CSS: cssStr}
		w.Header().Set("Content-Type", "application/json")
		out := struct{ HTML string ` + "`json:\"html\"`" + ` }{}
		var buf []byte
		tmp := &bufWriter{buf: &buf}
		if err := tpl.Execute(tmp, data); err != nil {
			http.Error(w, err.Error(), 500); return
		}
		out.HTML = string(*tmp.buf)
		json.NewEncoder(w).Encode(out)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type bufWriter struct{ buf *[]byte }
func (b *bufWriter) Write(p []byte) (int, error) { *b.buf = append(*b.buf, p...); return len(p), nil }
