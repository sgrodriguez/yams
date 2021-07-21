package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// WildCardMux uses to handle all non registered resources
type WildCardMux struct {
	mux       *http.ServeMux
	WCHandler http.Handler
}

func (a *WildCardMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, pat := a.mux.Handler(r)
	if pat == "" {
		a.WCHandler.ServeHTTP(w, r)
		return
	}
	a.mux.ServeHTTP(w, r)
}

// Handle ...
func (a *WildCardMux) Handle(pattern string, handler http.Handler) {
	a.mux.Handle(pattern, handler)
}

// Mocks ...
type Mocks struct {
	mocks sync.Map
}

// Request ...
type Request struct {
	Resource string `json:"resource"`
	Method   string `json:"method"`
}

// Response ...
type Response struct {
	HTTPStatusCode int    `json:"http_status_code"`
	Body           string `json:"body"`
}

// Mock ...
type Mock struct {
	Name           string   `json:"name"`
	DelayInSeconds int      `json:"delay_in_seconds"`
	Request        Request  `json:"request"`
	Response       Response `json:"response"`
}

func main() {
	tomlPath := "example_config.toml"
	if len(os.Args) > 1 {
		tomlPath = os.Args[len(os.Args)-1]
	}
	cfg, err := LoadConfig(tomlPath)
	if err != nil {
		log.Fatal(fmt.Errorf("error loading config %w", err))
	}
	m := &Mocks{
		mocks: sync.Map{},
	}
	for i := range cfg.Mocks {
		mo := &cfg.Mocks[i]
		m.mocks.Store(mo.Request.Resource+mo.Request.Method, mo)
	}
	serverMux := &WildCardMux{http.NewServeMux(), m.MockHandler()}
	serverMux.Handle("/mock", m.AddMockHandler())
	serverMux.Handle("/mocks", m.ListMocksHandler())
	http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), serverMux)
}

// MockHandler ...
func (m *Mocks) MockHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mo, ok := m.mocks.Load(r.URL.Path + r.Method)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		mock, ok := mo.(*Mock)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if mock.DelayInSeconds > 0 {
			time.Sleep(time.Second * time.Duration(mock.DelayInSeconds))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(mock.Response.HTTPStatusCode)
		w.Write([]byte(mock.Response.Body))
	}
}

// ListMocksHandler ...
func (m *Mocks) ListMocksHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mocks := make([]*Mock, 0)
		m.mocks.Range(func(k, v interface{}) bool {
			mo, ok := v.(*Mock)
			if !ok {
				return false
			}
			mocks = append(mocks, mo)
			return true
		})
		mocksEncoded, err := json.Marshal(mocks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mocksEncoded)
	}
}

// AddMockHandler ...
func (m *Mocks) AddMockHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var mo Mock
		err := json.NewDecoder(r.Body).Decode(&mo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		m.mocks.Store(mo.Request.Resource+mo.Request.Method, &mo)
		w.WriteHeader(http.StatusCreated)
	}
}
