package utils

import (
	"github.com/dolong2110/Memoirization-Apps/account/model/apperrors"
	"github.com/google/uuid"
	"log"
	"net/url"
	"path"
)

func ObjNameFromURL(imageURL string) (string, error) {
	// if user doesn't have imageURL - create one
	// otherwise, extract last part of URL to get cloud storage object name
	if imageURL == "" {
		objID, _ := uuid.NewRandom()
		return objID.String(), nil
	}

	// split off last part of URL, which is the image's storage object ID
	urlPath, err := url.Parse(imageURL)
	if err != nil {
		log.Printf("Failed to parse objectName from imageURL: %v\n", imageURL)
		return "", apperrors.NewInternal()
	}

	// get "path" of url (everything after domain)
	// then get "base", the last part
	return path.Base(urlPath.Path), nil
}
