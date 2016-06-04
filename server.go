package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/satori/go.uuid"
	"github.com/unrolled/render"
)

func runServer(port string) {
	r := render.New(render.Options{
		IsDevelopment: os.Getenv("ENVIRONMENT") != "production",
		Layout:        "layout",
	})

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return renderTemplate(c, r, http.StatusOK, "form", nil)
	})
	e.POST("/", processUpload)
	e.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		return renderTemplate(c, r, http.StatusOK, "result", id)
	})

	e.Static("/wavs", "wavs")
	e.Run(standard.New(":" + port))
}

func processUpload(c echo.Context) error {
	id := c.FormValue("id")
	if id == "" {
		id = uuid.NewV4().String()
	}

	outputDir := "wavs/" + id
	os.MkdirAll(outputDir, 0775)

	_, err := writeUploadedFile(c, "input", outputDir+"/input.wav")
	if err != nil {
		return err
	}

	filename, err := writeUploadedFile(c, "payload", outputDir+"/payload.bin")
	if err != nil {
		return err
	}

	// Fire goroutines to encode everything
	encodePayloads(outputDir, filename)

	return c.Redirect(http.StatusMovedPermanently, "/"+id)
}

func encodePayloads(workingDir, filename string) error {
	inputPath := workingDir + "/input.wav"
	payload, err := ioutil.ReadFile(workingDir + "/payload.bin")
	if err != nil {
		return err
	}
	wg := new(sync.WaitGroup)
	for lsbsToUse := 1; lsbsToUse <= 32; lsbsToUse++ {
		wg.Add(1)
		go func(lsbsToUse int) {
			defer wg.Done()
			outputPath := fmt.Sprintf("%s/output-%d.wav", workingDir, lsbsToUse)
			// This might error but for now we don't care
			if err := encode(inputPath, outputPath, lsbsToUse, filename, payload); err != nil {
				fmt.Println(err)
			}
		}(lsbsToUse)
	}
	wg.Wait()
	return nil
}

func writeUploadedFile(c echo.Context, fileField, outputFile string) (string, error) {
	file, err := c.FormFile(fileField)
	if err != nil {
		return "", err
	}
	source, err := file.Open()
	if err != nil {
		return "", err
	}
	defer source.Close()

	destination, err := os.Create(outputFile)
	if err != nil {
		return "", err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)

	return file.Filename, err
}

func renderTemplate(c echo.Context, r *render.Render, status int, name string, binding interface{}) error {
	if w, ok := c.Response().(*standard.Response); ok {
		return r.HTML(w.ResponseWriter, status, name, binding)
	} else {
		panic("Ooops")
		return nil
	}
}
