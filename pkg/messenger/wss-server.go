package wss_server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	//"github.com/alex988334/messenger/pkg/messenger/db"
)

type WssServer struct {
	running bool
	hub     *Hub
	router  http.ServeMux
	server  *http.Server
}

func NewWssServer() *WssServer {
	return &WssServer{
		running: false,
		hub:     NewHub(),
		server:  nil,
	}
}

func (s *WssServer) Run() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered WSServer:", r)
		}
	}()

	go s.hub.Run()

	mux := http.NewServeMux()

	mux.HandleFunc("/"+downloadClient, func(responseWriter http.ResponseWriter, request *http.Request) {

		/*dat, err := os.ReadFile(versionFile)
		if err != nil {
			fmt.Print("ERROR WSServer! Error read file \"version_client.txt\"", err)
			responseWriter.Write([]byte("0"))
		} else {
			responseWriter.Write(dat)
		}
		*/
		request.Body.Close()
	})

	mux.HandleFunc("/"+aboutProject, func(responseWriter http.ResponseWriter, request *http.Request) {

		dat, err := os.ReadFile(aboutProjectFile)
		if err != nil {
			fmt.Print("ERROR WSServer! Error read file \"about_project.txt\"", err)
			responseWriter.Write([]byte("0"))
		} else {
			responseWriter.Write(dat)
		}

		request.Body.Close()
	})

	mux.HandleFunc("/"+versionClientPath, func(responseWriter http.ResponseWriter, request *http.Request) {

		dat, err := os.ReadFile(versionFile)
		if err != nil {
			fmt.Print("ERROR WSServer! Error read file \"version_client.txt\"", err)
			responseWriter.Write([]byte("0"))
		} else {
			responseWriter.Write(dat)
		}

		request.Body.Close()
	})

	mux.HandleFunc("/"+localChatUrl, func(responseWriter http.ResponseWriter, request *http.Request) {
		//	передаем в функцию сам хаб, карту заголовков и сам запрос
		serveWs(s.hub, responseWriter, request)
	})

	mux.HandleFunc("/orion-exit", func(responseWriter http.ResponseWriter, request *http.Request) {
		//	передаем в функцию сам хаб, карту заголовков и сам запрос
		resp, _ := json.Marshal("Orion shutdown")
		responseWriter.Write(resp)
		s.hub.stopServer <- true
		s.server.Shutdown(context.Background())
		fmt.Println("Orion Close")
	})

	s.server = &http.Server{
		Handler: mux,
		Addr:    host /* hostlocal*/ + ":" + port,
	}

	err := s.server.ListenAndServeTLS("/etc/letsencrypt/live/orion734.ru/fullchain.pem",
		"/etc/letsencrypt/live/orion734.ru/privkey.pem")
	//	*/
	if /*err := s.server.ListenAndServe();*/ err != nil { //*/ err != nil {
		fmt.Println("WSServer error:", err)
	} else {
		fmt.Println("WSServer stop Listen")
	}
}

func (s *WssServer) Stop() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal("ERROR SHUTDOWN HTTP SERVER 86273635:", err)
	}
	s.hub.stopServer <- false
	s.running = false
	// fmt.Println("stop MODULE WSS")
}

func (s *WssServer) GetStatus() bool {
	//fmt.Println("NOW WSS ", s.running)
	return s.running
}

/* echo request handler

mux.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		//	передаем в функцию сам хаб, карту заголовков и сам запрос
		//serveWs(s.hub, responseWriter, request)
		resp, _ := json.Marshal("Response Orion")
		rew := request.Body

		var p []byte = make([]byte, 256)

		rew.Read(p)
		fmt.Println("quer =>", string(p))
		p = bytes.Trim(p, "\x00")
		fmt.Println("quer =>", string(p))
		var v Rew = &Mod{Wer: Wer{}, Prop1: 3, Prop2: 4}
		by, er := json.Marshal(v)
		fmt.Println("json.Marshal(v) =>", string(by))
		er = json.Unmarshal(p, &v)
		if er != nil {
			fmt.Println("jsonerr =>", er)
		}
		by, er = json.Marshal(v)
		fmt.Println("json.Marshal(v) =>", string(by))

		responseWriter.Write(resp)
		request.Body.Close()
	}) // */

/*
	mux.HandleFunc("/"+avatarPath, func(responseWriter http.ResponseWriter, request *http.Request) {
		//	передаем в функцию сам хаб, карту заголовков и сам запрос
		serveWs(s.hub, responseWriter, request)
	})
*/

/*
func processingAvatar(responseWriter http.ResponseWriter, request *http.Request) {

	file, fileHeader, err := request.FormFile("file_name")
	request.Body.
	defer file.Close()

	// copy example
	f, err := os.OpenFile("./downloaded", os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	io.Copy(f, file)

}*/
