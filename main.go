package main

import (
	"bytes"
	"context"
	"fmt"
	"go-generate-pdf-demo/htmltopdf"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/generate-pdf", handleGeneratePDF)

	server := &http.Server{
		Handler: router,
		Addr:    ":9000",
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	<-sigterm

	server.Shutdown(context.Background())
}

func handleGeneratePDF(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r := recover()
		if r != nil {
			writeError(w, fmt.Errorf("error: %v", r))
		}
	}()
	t, err := template.ParseFiles("./template/template.html")
	if err != nil {
		writeError(w, err)
	}
	asiaJakarta, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(asiaJakarta)
	title := fmt.Sprintf("try_to_generate_pdf_%d", now.Unix())
	data := map[string]interface{}{
		"title": title,
	}
	html := new(bytes.Buffer)

	if err = t.Execute(html, data); err != nil {
		writeError(w, err)
	}

	pdf, err := htmltopdf.NewHTMLToPDF().GenerateFromReader(context.Background(), html)
	if err != nil {
		writeError(w, err)
	}

	f, err := os.Create(fmt.Sprintf("%s.pdf", title))
	if err != nil {
		writeError(w, err)
	}
	defer f.Close()

	buffForRead := new(bytes.Buffer)
	_, err = buffForRead.ReadFrom(pdf)
	if err != nil {
		writeError(w, err)
	}

	_, err = f.Write(buffForRead.Bytes())
	if err != nil {
		writeError(w, err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("OK with name %s", title)))
}

func writeError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
