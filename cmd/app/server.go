package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gholib/http/pkg/banners"
)

//Server .. это наш логический сервер
type Server struct {
	mux       *http.ServeMux
	bannerSvc *banners.Service
}

//NewServer .. Функция для создание нового сервера
func NewServer(m *http.ServeMux, bnrSvc *banners.Service) *Server {
	return &Server{mux: m, bannerSvc: bnrSvc}
}

//ServeHTTP ... метод для запуска сервера
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//Init .. мотод для инициализации сервера
func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveByID)
}

func (s *Server) handleGetBannerByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	idParam := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		log.Print(err)
		errorWriter(w, http.StatusBadRequest)
		return
	}

	item, err := s.bannerSvc.ByID(r.Context(), id)
	if err != nil {
		//печатаем ошибку
		log.Print(err)

		errorWriter(w, http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)

	if err != nil {
		log.Print(err)
		errorWriter(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGetAllBanners(w http.ResponseWriter, r *http.Request) {
	item, err := s.bannerSvc.All(r.Context())
	if err != nil {
		log.Print(err)

		errorWriter(w, http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)

	if err != nil {
		log.Print(err)
		errorWriter(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleSaveBanner(w http.ResponseWriter, r *http.Request) {

	//получаем данные из параметра запроса
	idP := r.PostFormValue("id")
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	button := r.PostFormValue("button")
	link := r.PostFormValue("link")

	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		log.Print(err)
		errorWriter(w, http.StatusBadRequest)
		return
	}

	if title == "" && content == "" && button == "" && link == "" {
		log.Print(err)
		errorWriter(w, http.StatusBadRequest)
		return
	}

	item := &banners.Banner{
		ID:      id,
		Title:   title,
		Content: content,
		Button:  button,
		Link:    link,
	}

	file, fileHeader, err := r.FormFile("image")

	if err == nil {
		var name = strings.Split(fileHeader.Filename, ".")
		item.Image = name[len(name)-1]
	}

	banner, err := s.bannerSvc.Save(r.Context(), item, file)

	if err != nil {
		log.Print(err)

		errorWriter(w, http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(banner)
	if err != nil {
		log.Print(err)

		errorWriter(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleRemoveByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		log.Print(err)
		errorWriter(w, http.StatusBadRequest)
		return
	}

	banner, err := s.bannerSvc.RemoveByID(r.Context(), id)

	if err != nil {
		log.Print(err)
		errorWriter(w, http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(banner)

	if err != nil {
		log.Print(err)
		errorWriter(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

//это функция для записывание ошибки в responseWriter или просто для ответа с ошиками
func errorWriter(w http.ResponseWriter, httpSts int) {
	http.Error(w, http.StatusText(httpSts), httpSts)
}
