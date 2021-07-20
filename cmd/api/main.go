package main

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hipeople.api/internal"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Api struct {
	config       *internal.Config
	imageService *internal.ImageService
	user         string
}

func NewApi(config *internal.Config, service *internal.ImageService) *Api {
	return &Api{
		config:       config,
		imageService: service,
	}
}

func (a *Api) handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/images/", a.checkAuth(a.imagesRouter))
	mux.Handle("/images/direct/", a.checkAuth(a.direct))
	return mux
}

func (a *Api) run() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: a.handler(),
	}
	log.Fatal(server.ListenAndServe())
}

func main() {
	config := internal.NewConfig()
	// In production app I would use real DB with some cache in front of it.
	db := make(internal.FakeDB)
	imageRepository := internal.NewImageRepository(db)
	imageService := internal.NewImageService(imageRepository)
	api := NewApi(config, imageService)
	api.run()
}

func (a *Api) findUser(splitToken []string) (string, error) {
	if len(splitToken) != 2 || splitToken[0] != "Bearer" {
		return "", errors.New("findUser.notFound")
	}
	for _, token := range a.config.Tokens {
		if subtle.ConstantTimeCompare([]byte(splitToken[1]), []byte(token.Token)) == 1 {
			return token.User, nil
		}
	}
	return "", errors.New("findUser.notFound")
}

func (a *Api) checkAuth(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		forbidden := true
		if len(request.Header["Authorization"]) >= 1 {
			if user, err := a.findUser(strings.Fields(request.Header["Authorization"][0])); err == nil {
				a.user = user
				forbidden = false
			}
		}
		if forbidden {
			writer.WriteHeader(http.StatusForbidden)
			return
		}
		handler.ServeHTTP(writer, request)
	})
}

func (a *Api) direct(writer http.ResponseWriter, request *http.Request) {
	log.Printf("direct, url: %v, method: %v", request.URL, request.Method)
	if request.Method != http.MethodGet {
		a.notSupported(writer, request, "GET")
		return
	}
	imageId, err := parseImageId(request.URL.Path[strings.LastIndex(request.URL.Path, "/")+1:])
	if err != nil {
		log.Println("getImages.BadRequest,", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	image := a.imageService.SelectImage(imageId, a.user)
	if image == nil {
		log.Println("direct.NotFound")
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.Header().Set("Content-Type", image.Mime)
	dec, _ := base64.URLEncoding.DecodeString(image.Data)
	writer.Write(dec)
}

func (a *Api) imagesRouter(writer http.ResponseWriter, request *http.Request) {
	log.Printf("imagesRouter, url: %v, method: %v, user %v", request.URL, request.Method, a.user)
	lastPath := request.URL.Path[strings.LastIndex(request.URL.Path, "/")+1:]
	switch request.Method {
	case http.MethodGet:
		a.getImages(writer, request, lastPath)
	case http.MethodPost:
		a.postImages(writer, request, lastPath)
	default:
		a.notSupported(writer, request, "GET, POST")
	}
}

func (a *Api) getImages(writer http.ResponseWriter, _ *http.Request, lastPath string) {
	imageId, err := parseImageId(lastPath)
	if err != nil {
		log.Println("getImages.BadRequest,", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var result interface{}
	if imageId == 0 {
		result = a.imageService.SelectImages(a.user)
	} else {
		result = a.imageService.SelectImage(imageId, a.user)
	}
	if result == nil {
		log.Println("getImages.NotFound")
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	resultJson, _ := json.Marshal(result)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(resultJson)
}

func (a *Api) postImages(writer http.ResponseWriter, request *http.Request, lastPath string) {
	if lastPath != "" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("postImages.StatusMethodNotAllowed")
		return
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		log.Println("postImages.UnprocessableEntity,", err)
		return
	}
	image, err := internal.NewImageJson(&body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("postImages.BadRequest,", err)
		return
	}

	imageId := a.imageService.InsertImage(image, a.user)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(fmt.Sprintf("{\"image_id\": %v}", imageId)))
}

func (a *Api) notSupported(writer http.ResponseWriter, _ *http.Request, methods string) {
	writer.Header().Set("Allow", methods)
	writer.WriteHeader(http.StatusMethodNotAllowed)
	log.Println("notSupported.MethodNotAllowed")
}

func parseImageId(lastPath string) (int, error) {
	if lastPath != "" {
		if parsedInt, err := strconv.ParseInt(lastPath, 10, 64); err != nil {
			return 0, err
		} else {
			return int(parsedInt), err
		}
	}
	return 0, nil
}
