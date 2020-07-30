/*

 */

package main

import (
	"log"
	"net/http"

	"github.com/bitspawngg/tournament-bracket-manager/server"
)

func main() {
	srv := server.CreateServer()
	// listen and serve
	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {

		log.Println("server shutting down gracefully...")
	} else {

		log.Println("unexpected server shutdown...")
		log.Println("ERR: ", err)
	}
}
