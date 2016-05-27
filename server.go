package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

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
	e.POST("/", func(c echo.Context) error {
		sourceFileHeader, err := c.FormFile("file")
		if err != nil {
			return err
		}
		sourceFile, err := sourceFileHeader.Open()
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		lsbsToUse, err := strconv.Atoi(c.FormValue("lsbs"))
		if err != nil {
			return err
		}

		payload := c.FormValue("payload")
		id, err := embedStrPayloadOnUploadedFile(payload, sourceFile, lsbsToUse)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusMovedPermanently, "/"+id)
	})
	e.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		return renderTemplate(c, r, http.StatusOK, "result", id)
	})
	e.Static("/wavs", "wavs")
	e.Run(standard.New(":" + port))
}

func embedStrPayloadOnUploadedFile(payload string, sourceFile io.Reader, lsbsToUse int) (string, error) {
	id := uuid.NewV4().String()
	outputDir := "wavs/" + id
	os.MkdirAll(outputDir, 0775)

	inputPath := outputDir + "/input.wav"
	outputPath := fmt.Sprintf("%s/output-%d.wav", outputDir, lsbsToUse)

	// Write source file to output dir
	inputFile, err := os.Create(inputPath)
	if err != nil {
		return "", err
	}
	defer inputFile.Close()

	// Copy
	if _, err = io.Copy(inputFile, sourceFile); err != nil {
		return "", err
	}

	// Encode it
	if err := encode(inputPath, outputPath, lsbsToUse, []byte(payload)); err != nil {
		return "", err
	}

	return id, nil
}

func renderTemplate(c echo.Context, r *render.Render, status int, name string, binding interface{}) error {
	if w, ok := c.Response().(*standard.Response); ok {
		return r.HTML(w.ResponseWriter, status, name, binding)
	} else {
		panic("Ooops")
	}
	return nil
}
