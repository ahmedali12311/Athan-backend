package validator

import (
	"errors"
	"fmt"
	"math/rand"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"app/config"
	"app/pkg/interfaces"

	"github.com/disintegration/imaging"
)

func (v *Validator) AssignImage(
	key string,
	m interfaces.HasImage,
	required bool,
	allowedScopes ...string,
) error {
	v.SaveOldImgThumbDists(m)

	if !v.Data.FileExists(key) && required {
		v.Check(false, key, v.T.ValidateRequired())
	}

	if v.Data.FileExists(key) {
		v.Permit(key, allowedScopes)
		img := v.Data.GetFile(key)
		_, params, err := mime.ParseMediaType(
			img.Header.Get("Content-Disposition"),
		)
		if err != nil {
			return err
		}

		filename := params["filename"]
		imgName, err := v.imageName(filename, m.TableName())
		if err != nil {
			return err
		}
		imgBytes, err := v.Data.GetFileBytes(key)
		if err != nil {
			return err
		}

		imgVal := filepath.Join("uploads", m.TableName(), imgName)
		thumbVal := filepath.Join(
			"uploads",
			m.TableName(),
			"thumbs",
			fmt.Sprintf("thumb_%s", imgName),
		)

		m.SetImg(&imgVal)
		m.SetThumb(&thumbVal)

		// public is a hidden path on live urls are in the format:
		// https://proj.com/uploads/banners/thumbs/thumb_banners_1637_9577.jpeg
		// thats why the database value is set without it,
		// but the OS path is full
		imgDist := config.GetRootPath(filepath.Join("public", imgVal))
		thumbDist := config.GetRootPath(filepath.Join("public", thumbVal))

		distpath := filepath.Join(
			"public",
			"uploads",
			m.TableName(),
			"thumbs",
		)
		_, err = os.Stat(distpath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if err := os.MkdirAll(distpath, 0o750); err != nil {
					return err
				}
				v.Logger.Info().Msgf(
					"Directory created successfully: %s",
					distpath,
				)
			}
		}

		// Create a new file in the uploads directory
		dist, err := os.Create(filepath.Clean(imgDist))
		if err != nil {
			return err
		}
		defer dist.Close()

		if _, err := dist.WriteString(string(imgBytes)); err != nil {
			return err
		}
		if err := v.generateThumb(img, thumbDist); err != nil {
			return err
		}
		v.newImg = imgDist
		v.newThumb = thumbDist
		v.DeleteOldPicture()
	}
	return nil
}

// imageName for uploaded files.
func (v *Validator) imageName(filename, dir string) (string, error) {
	// this slice must be sorted alphabetically
	mimetypes := []string{
		// ".jfif",
		".jpe",
		".jpeg",
		".jpg",
		".png",
		// ".webp",
	}
	ext := filepath.Ext(filename)
	if ok := slices.Contains(mimetypes, ext); !ok {
		mimeError := fmt.Errorf(
			"file extension not allowed: %s, allowed files are: %s",
			ext,
			strings.Join(mimetypes, ","),
		)
		return "", mimeError
	}
	// Image data
	randomNum := rand.Int63n(1_000_000) //nolint:gosec // doesn't matter
	imgName := dir +
		"_" +
		strconv.FormatInt(config.TimeNow().UnixNano(), 10) +
		"_" +
		strconv.FormatInt(randomNum, 10) +
		ext

	return imgName, nil
}

// generateThumb for the form file provided in a post or put request.
func (v *Validator) generateThumb(
	file *multipart.FileHeader,
	dist string,
) error {
	imageFile, err := file.Open()
	if err != nil {
		return err
	}
	decodedImg, err := imaging.Decode(imageFile)
	if err != nil {
		return err
	}
	resizedThumb := imaging.Resize(decodedImg, 150, 0, imaging.Lanczos)

	f, err := os.Create(filepath.Clean(dist))
	if err != nil {
		return err
	}
	defer f.Close()

	return imaging.Encode(f, resizedThumb, imaging.PNG)
}

// deleteFile removes a single file provided dist string from system.
func (v *Validator) deleteFile(dist string) {
	if len(dist) > 0 && strings.Contains(dist, ".") {
		if err := os.Remove(dist); err != nil {
			// file is deleted
			err = fmt.Errorf("failed to delete file: %s, error: %w", dist, err)
			v.Logger.Error().
				Err(err).
				Msg("LogOnlyError")
		}
	}
}

// DeleteNewPicture removes a newly uploaded image and its thumb.
func (v *Validator) DeleteNewPicture() {
	if v.newImg != "" {
		v.deleteFile(v.newImg)
	}
	if v.newThumb != "" {
		v.deleteFile(v.newThumb)
	}
}

// SaveOldImgThumbDists sets old file path instead of url img,
// thumb values on validator.
func (v *Validator) SaveOldImgThumbDists(m interfaces.HasImage) {
	v.oldImg = m.GetImg()
	if v.oldImg != nil {
		imgNoDomain := strings.ReplaceAll(
			*v.oldImg,
			config.DOMAIN+"/",
			"",
		)
		oldImgDist := config.GetRootPath(filepath.Join("public", imgNoDomain))
		v.oldImg = &oldImgDist
		// v.newImg = oldImgDist
		m.SetImg(&imgNoDomain)
	}

	v.oldThumb = m.GetThumb()
	if v.oldThumb != nil {
		thumbNoDomain := strings.ReplaceAll(
			*v.oldThumb,
			config.DOMAIN+"/",
			"",
		)
		oldThumbDist := config.GetRootPath(
			filepath.Join("public", thumbNoDomain),
		)
		v.oldThumb = &oldThumbDist
		// v.newThumb = oldThumbDist
		m.SetThumb(&thumbNoDomain)
	}
}

// DeleteOldPicture removes an existing image and its thumb
// after successful update of new files.
func (v *Validator) DeleteOldPicture() {
	if v.oldImg != nil {
		v.deleteFile(*v.oldImg)
	}
	if v.oldThumb != nil {
		v.deleteFile(*v.oldThumb)
	}
}
