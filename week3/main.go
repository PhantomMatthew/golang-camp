package main

import (
	"context"
	"encoding/json"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
)

//Start http server
func StartHttpServer(srv *http.Server) error {
	http.HandleFunc("/handler", Handler)
	log.Println("http server start")
	err := srv.ListenAndServe()
	return err
}

// Start http handler
func Handler(w http.ResponseWriter, req *http.Request) {
	log.Println(req.Body)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Handler is called"
	jsonResp, err := json.Marshal(resp)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func main() {
	ctx := context.Background()
	// To use context along with WithCancel
	ctx, cancel := context.WithCancel(ctx)

	group, errCtx := errgroup.WithContext(ctx)

	srv := &http.Server{Addr: ":18080"}

	group.Go(func() error {
		return StartHttpServer(srv)
	})

	group.Go(func() error {
		<-errCtx.Done()
		log.Println("Start http server")
		return srv.Shutdown(errCtx)
	})

	// Use os.Signal to make channel
	channel := make(chan os.Signal, 1)
	signal.Notify(channel)

	group.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
				return errCtx.Err()
			case <-channel:
				cancel()
			}
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		log.Println("Group error info: ", err)
	}
	log.Println("All Group have done!")

}
