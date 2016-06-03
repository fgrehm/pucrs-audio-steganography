package main

import (
	"fmt"
	"io"
	"io/ioutil"
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
		input, err := openUploadedFile(c, "input")
		if err != nil {
			return err
		}
		defer input.Close()

		payload, err := openUploadedFile(c, "payload")
		if err != nil {
			return err
		}
		defer payload.Close()

		payloadFile, _ := c.FormFile("payload")

		id, err := processUpload(input, payload, payloadFile.Filename)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusMovedPermanently, "/"+id)
	})
	e.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		return renderTemplate(c, r, http.StatusOK, "result", id)
	})
	e.GET("/:id/:lsbs", func(c echo.Context) error {
		id := c.Param("id")
		lsbs := c.Param("lsbs")
		workingDir := "wavs/" + id
		filePath := fmt.Sprintf("%s/output-%s.wav", workingDir, lsbs)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			payload, err := ioutil.ReadFile(workingDir + "/payload.bin")
			if err != nil {
				return err
			}
			lsbsToUse, err := strconv.Atoi(lsbs)
			if err != nil {
				return err
			}
			info, err := readInfo(workingDir + "/info.json")
			if err != nil {
				return err
			}
			err = encode(workingDir + "/input.wav", filePath, lsbsToUse, info.Filename, payload)
			if err != nil {
				return err
			}
		}

		return c.File(filePath)
	})
	// TODO: REMOVE
	e.Static("/wavs", "wavs")
	e.Run(standard.New(":" + port))
}

func openUploadedFile(c echo.Context, fileField string) (io.ReadCloser, error) {
	fileHeader, err := c.FormFile(fileField)
	if err != nil {
		return nil, err
	}
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	return file, nil
}

func processUpload(input, payload io.Reader, filename string) (string, error) {
	id := uuid.NewV4().String()
	outputDir := "wavs/" + id
	os.MkdirAll(outputDir, 0775)

	if err := writeUploadedFile(input, outputDir+"/input.wav"); err != nil {
		return "", err
	}
	if err := writeUploadedFile(payload, outputDir+"/payload.bin"); err != nil {
		return "", err
	}
	if err := writeInfo(outputDir+"/input.wav", outputDir +"/payload.bin", outputDir+"/info.json", filename); err != nil {
		return "", err
	}

	return id, nil
}

func writeUploadedFile(source io.Reader, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, source)

	return err
}

func renderTemplate(c echo.Context, r *render.Render, status int, name string, binding interface{}) error {
	if w, ok := c.Response().(*standard.Response); ok {
		return r.HTML(w.ResponseWriter, status, name, binding)
	} else {
		panic("Ooops")
	}
	return nil
}
