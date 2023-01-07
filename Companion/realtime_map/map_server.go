package realtime_map

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type MapServer struct {
	httpServer *http.Server
}

func NewMapServer() (*MapServer, error) {
	curExePath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	curExeDir := filepath.Dir(curExePath)

	if err != nil {
		return nil, err
	}

	fs := http.FileServer(http.Dir(path.Join(curExeDir, "map")))
	server := &http.Server{
		Handler: fs,
		Addr:    ":8000",
	}

	return &MapServer{
		httpServer: server,
	}, nil
}

func (ms *MapServer) Start() {
	go func() {
		ms.httpServer.ListenAndServe()
		log.Println("stopping map server")
	}()
}

func (ms *MapServer) Stop() error {
	return ms.httpServer.Close()
}
