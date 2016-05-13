package stns

import (
	"reflect"
	"regexp"

	"github.com/STNS/STNS/settings"
	"github.com/ant0ine/go-json-rest/rest"
)

type Handler struct {
	config *Config
}

type MetaData struct {
	ApiVersion float64 `json:"api_version"`
	Salt       bool    `json:"salt_enable"`
	Stretching int     `json:"stretching_count"`
	Result     string  `json:"result""`
}

type ResponseFormat struct {
	MetaData *MetaData   `json:"metadata"`
	Items    *Attributes `json:"items"`
}

func (h *Handler) getQuery(r *rest.Request) *Query {
	value := r.PathParam("value")
	column := r.PathParam("column")
	resource_name := r.PathParam("resource_name")
	return &Query{h.config, resource_name, column, value}
}

func (h *Handler) getListQuery(r *rest.Request) *Query {
	resource_name := r.PathParam("resource_name")
	return &Query{h.config, resource_name, "list", ""}
}

func (h *Handler) Get(w rest.ResponseWriter, r *rest.Request) {
	query := h.getQuery(r)
	h.Response(query, w, r)
}

func (h *Handler) Auth(w rest.ResponseWriter, r *rest.Request) {
	query := h.getQuery(r)
	h.AuthResponse(query, w, r)
}

func (h *Handler) GetList(w rest.ResponseWriter, r *rest.Request) {
	query := h.getListQuery(r)
	h.Response(query, w, r)
}

func (h *Handler) Response(q *Query, w rest.ResponseWriter, r *rest.Request) {
	attr := q.Get()
	if attr == nil || reflect.ValueOf(attr).IsNil() {
		rest.NotFound(w, r)
		return
	}

	v2 := regexp.MustCompile(`^/v2`)
	if v2.MatchString(r.URL.Path) {
		response := ResponseFormat{
			&MetaData{
				settings.API_VERSION,
				h.config.Salt,
				h.config.Stretching,
				settings.SUCCESS,
			},
			&attr,
		}
		w.WriteJson(response)
	} else {
		w.WriteJson(attr)
	}
}

func (h *Handler) AuthResponse(q *Query, w rest.ResponseWriter, r *rest.Request) {
	attr := q.Get()
	if attr == nil || reflect.ValueOf(attr).IsNil() {
		rest.NotFound(w, r)
		return
	}

	for _, params := range attr {
		if params.Password == r.PathParam("hash") {
			response := ResponseFormat{
				&MetaData{
					settings.API_VERSION,
					h.config.Salt,
					h.config.Stretching,
					settings.SUCCESS,
				},
				&Attributes{},
			}
			w.WriteJson(response)
			return
		}
	}
	rest.NotFound(w, r)
}
