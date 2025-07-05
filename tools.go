package toolkit

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const RandomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789+_"

// Tools is the type used to instantiate this module, any variable of type Tools will have access
// to the reciever methods defined on the *Tools type

type Tools struct {
	MaxFileSize      int
	AllowedFileTypes []string
}

// RandomString takes an integer n and returns a string of n random characters defined in the
// RandomStringSource
func (t *Tools) RandomString(n int) string {
	s, r := make([]rune, n), []rune(RandomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}

type UploadedFile struct {
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

func (t *Tools) UploadFile(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	var UploadedFiles []*UploadedFile

	if t.MaxFileSize == 0 {
		t.MaxFileSize = 1024 * 1024 * 1024
	}

	err := r.ParseMultipartForm(int64(t.MaxFileSize))
	if err != nil {
		return nil, errors.New("file uploaded is to large")
	}

	for _, fileHeader := range r.MultipartForm.File {
		for _, hdr := range fileHeader {
			UploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile
				infile, err := hdr.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()
				buf := make([]byte, 512)
				_, err = infile.Read(buf)
				if err != nil {
					return nil, err
				}

				allowed := false
				fileType := http.DetectContentType(buf)

				if len(t.AllowedFileTypes) > 0 {
					for _, x := range t.AllowedFileTypes {
						if strings.EqualFold(fileType, x) {
							allowed = true
						}
					}
				} else {
					allowed = true
				}

				if !allowed {
					return nil, errors.New("the uploaded file type is not permitted")
				}

				_, err = infile.Seek(0, 0)
				if err != nil {
					return nil, err
				}

				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", t.RandomString(25), filepath.Ext(hdr.Filename))
				} else {
					uploadedFile.NewFileName = hdr.Filename
				}

				var outputFile *os.File
				defer outputFile.Close()

				if outputFile, err := os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName)); err != nil {
					return nil, err
				} else {
					fileSize, err := io.Copy(outputFile, infile)
					if err != nil {
						return nil, err
					}
					uploadedFile.FileSize = fileSize
				}
				uploadedFiles = append(uploadedFiles, &uploadedFile)
				return uploadedFiles, nil
			}(UploadedFiles)
			if err != nil {
				return nil, err
			}
		}
	}
	return UploadedFiles, nil
}
