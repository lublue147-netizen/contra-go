package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Ensure .wasm MIME type is registered
	_ = mime.AddExtensionType(".wasm", "application/wasm")

	dir := "./dist"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		dir = "./web"
	}

	absDir, _ := filepath.Abs(dir)
	log.Printf("Serving Contra arcade from %s on :%s\n", absDir, port)

	fs := http.FileServer(http.Dir(dir))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers for WebAssembly caching & security
		if filepath.Ext(r.URL.Path) == ".wasm" {
			w.Header().Set("Content-Type", "application/wasm")
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		fs.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
