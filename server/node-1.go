package main

import (
	"fmt"
	"github.com/imaimen/cdn-go/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Cache struct for caching
type Cache struct {
	files map[string]string
}

// NewCache create a new Cache
func NewCache() *Cache {
	return &Cache{files: make(map[string]string)}
}

// fetch files from cache or origin server
func (c *Cache) fetchFile(filename string, originServer string) (string, error) {
	// check if file is in cache
	filePath, found := c.files[filename]
	if found {
		// verify if the cached file actually exists on the file system
		if _, err := os.Stat(filePath); err == nil {
			// if the file exists, serve it from cache
			log.Println("Serving from cache:", filename)
			return filePath, nil
		} else {
			log.Println("Cache file missing:", filename, "- Re-fetching from origin...")
		}
	}

	// if the file is not in cache or missing, fetch it from origin server
	filePath = filepath.Join("./cache", filename)
	err := c.downloadFromOrigin(filename, originServer, filePath)
	if err != nil {
		return "", err
	}

	// save to cache
	c.files[filename] = filePath
	return filePath, nil
}

func (c *Cache) downloadFromOrigin(filename, originServer, destination string) error {
	url := fmt.Sprintf("http://%s/%s", originServer, filename)
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	// create destination file
	output, err := os.Create(destination)
	if err != nil {
		return err
	}

	// close the file
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {

		}
	}(output)

	// copy content from response to the file
	_, err = io.Copy(output, response.Body)
	if err != nil {
		return err
	}

	log.Println("Downloaded", filename, "from origin server")
	return nil
}

func main() {
	// initialize redis in the node edge
	utils.Connect()
	cache := NewCache()
	originServer := "localhost:8080"

	// create cache directory
	if _, err := os.Stat("./cache"); os.IsNotExist(err) {
		err := os.Mkdir("./cache", 0755)
		if err != nil {
			return
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Path[1:] // remove the leading dash (/)
		if filename == "" {
			http.Error(w, "File not specified", http.StatusBadRequest)
			return
		}

		// fetch file from cache or origin server
		filePath, err := cache.fetchFile(filename, originServer)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// server file
		http.ServeFile(w, r, filePath)
	})

	log.Println("Edge server running on :8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		utils.Disconnect()
		log.Fatal(err)
	}
}
