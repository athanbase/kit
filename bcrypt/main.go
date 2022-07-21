package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	muxRouter := mux.NewRouter()
	RegisterRoutes(muxRouter)
	server := &http.Server{
		Addr:    ":8080",
		Handler: muxRouter,
	}
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		log.Println("got exist signal")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal("Shutdown server:", err)
		}
	}()
	log.Println("starting http server...")
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Print("Server closed under request")
		} else {
			log.Fatal("Server closed unexpected")
		}
	}
}

func HashPassword(passwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetHashingCost(hashedPassword []byte) int {
	cost, _ := bcrypt.Cost(hashedPassword) // 为了简单忽略错误处理
	return cost
}

func PassWordHashingHandler(w http.ResponseWriter, r *http.Request) {
	password := "secret"
	start := time.Now()
	hash, _ := HashPassword(password) // 为了简单忽略错误处理
	end := time.Now()
	fmt.Fprintln(w, "used time secondes: ", end.Sub(start).Seconds(), "ms: ", end.Sub(start).Milliseconds())
	fmt.Fprintln(w, "Password:", password)
	fmt.Fprintln(w, "Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Fprintln(w, "Match:   ", match)

	cost := GetHashingCost([]byte(hash))
	fmt.Fprintln(w, "Cost:    ", cost)

}

func RegisterRoutes(r *mux.Router) {
	indexRouter := r.PathPrefix("/index").Subrouter()
	indexRouter.HandleFunc("/password_hashing", PassWordHashingHandler)
}
