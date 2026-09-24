package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/fateevanastusha/videos-go-backend/internal/models"
	"github.com/fateevanastusha/videos-go-backend/internal/validators"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) validationError(w http.ResponseWriter, r *http.Request, errs []validators.ErrorMessage) {
	s.respond(w, r, http.StatusBadRequest, map[string][]validators.ErrorMessage{"errorsMessages": errs})
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func Start(address string) error {
	server := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
	}
	server.configureRouter()

	return http.ListenAndServe(address, server)
}

func (s *server) configureRouter() {
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/testing/all-data", s.deleteAllData).Methods("DELETE")
	s.router.HandleFunc("/videos", s.getVideos()).Methods("GET")
	s.router.HandleFunc("/videos", s.createVideo()).Methods("POST")
	s.router.HandleFunc("/videos/{id}", s.getVideoById()).Methods("GET")
	s.router.HandleFunc("/videos/{id}", s.putVideoById()).Methods("PUT")
	s.router.HandleFunc("/videos/{id}", s.deleteVideoById()).Methods("DELETE")
}

var videos []*models.Video = []*models.Video{}
var nextId = 1

func (s *server) deleteAllData(w http.ResponseWriter, r *http.Request) {
	videos = []*models.Video{}
	nextId = 1
	s.respond(w, r, 204, nil)
}

// return all videos
func (s *server) getVideos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, 200, videos)
	}
}

func findVideoById(id int) *models.Video {
	for _, v := range videos {
		if v.Id == id {
			return v
		}
	}

	return nil
}

// create new video
func (s *server) createVideo() http.HandlerFunc {
	type request struct {
		Title                string              `json:"title"`
		Author               string              `json:"author"`
		AvailableResolutions []models.Resolution `json:"availableResolutions"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, 400, err)
			return
		}

		errs := validators.ValidateTitleAuthor(req.Title, req.Author)
		errs = append(errs, validators.ValidateResolutions(req.AvailableResolutions)...)
		if len(errs) > 0 {
			s.validationError(w, r, errs)
			return
		}

		now := time.Now().UTC()
		newVideo := &models.Video{
			Id:                   nextId,
			Title:                req.Title,
			Author:               req.Author,
			AvailableResolutions: req.AvailableResolutions,
			CanBeDownloaded:      false,
			MinAgeRestriction:    nil,
			CreatedAt:            now.Format(time.RFC3339),
			PublicationDate:      now.AddDate(0, 0, 1).Format(time.RFC3339),
		}
		nextId++
		videos = append(videos, newVideo)

		s.respond(w, r, 201, newVideo)
	}
}

// return video by id
func (s *server) getVideoById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			s.error(w, r, 400, err)
			return
		}

		video := findVideoById(id)
		if video == nil {
			s.error(w, r, 404, errors.New("not found video"))
			return
		}

		s.respond(w, r, 200, video)
	}
}

// update existing video by id with InputModel
func (s *server) putVideoById() http.HandlerFunc {
	type request struct {
		Title                string              `json:"title"`
		Author               string              `json:"author"`
		AvailableResolutions []models.Resolution `json:"availableResolutions"`
		CanBeDownloaded      bool                `json:"canBeDownloaded"`
		MinAgeRestriction    *int                `json:"minAgeRestriction"`
		PublicationDate      string              `json:"publicationDate"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			s.error(w, r, 400, err)
			return
		}

		video := findVideoById(id)
		if video == nil {
			s.error(w, r, 404, errors.New("not found video"))
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, 400, err)
			return
		}

		errs := validators.ValidateTitleAuthor(req.Title, req.Author)
		errs = append(errs, validators.ValidateResolutions(req.AvailableResolutions)...)
		if req.MinAgeRestriction != nil && (*req.MinAgeRestriction < 1 || *req.MinAgeRestriction > 18) {
			errs = append(errs, validators.ErrorMessage{Message: "must be between 1 and 18", Field: "minAgeRestriction"})
		}
		if len(errs) > 0 {
			s.validationError(w, r, errs)
			return
		}

		video.Author = req.Author
		video.Title = req.Title
		video.AvailableResolutions = req.AvailableResolutions
		video.CanBeDownloaded = req.CanBeDownloaded
		video.PublicationDate = req.PublicationDate
		video.MinAgeRestriction = req.MinAgeRestriction

		s.respond(w, r, 204, nil)
	}
}

// delete video specified by id
func (s *server) deleteVideoById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			s.error(w, r, 400, err)
			return
		}

		for i, v := range videos {
			if v.Id == id {
				videos = append(videos[:i], videos[i+1:]...)
				s.respond(w, r, 204, nil)
				return
			}
		}

		s.error(w, r, 404, errors.New("not found video"))
	}
}
