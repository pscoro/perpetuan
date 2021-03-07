package main

import (
	"encoding/json"
	"github.com/google/uuid"

	//"fmt"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
	//"strings"
)

//var connectionMap = make(map[string] player)
var clientMap = make(map[string]*client)
var server *socketio.Server
var serverPlayer = player {
	Username: "Server",
	Pos: position{
		Pos:     [2]int{0,0},
		Setting: "No Where",
	},
	Room: "",
}

var connect = initDB()

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	initDB()

	//var testGrid = [][]gridLocation{
	//	{},
	//}

	//addSetting("test", nil, 100, 100, testGrid)

	var settingList []setting

	settingList, ok := updateSettingList(settingList)
	if ok {
		log.Println("setting list updated: ", settingList)
	}

	var err error
	server, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("HI")
		s.SetContext("")

		clientMap[s.ID()] = &client{messages: make(chan message), commands: make(chan message), priority: make(chan message),
			session: session{
				ID:     uuid.UUID{},
				Player: player{},
				Room:   "",
			},
			conn: s,
		}

		server.JoinRoom("/", s.ID(), s)

		clientMap[s.ID()].session.Room = s.ID()
		//go clientMap[conn.ID()].run()

		err := handleClient(clientMap[s.ID()])
		server.JoinRoom("/", clientMap[s.ID()].session.Player.Room, s)

		if err != nil {
			log.Println("Error handling client")
			log.Fatal(err)
		}


		log.Println("connected:", s.ID())


		//conn.Emit("initSession", sess)
		//connectionMap[conn.ID()] = player{Name: "", Pos: position{[2]int{0,0}, "test"}}
		//message, err := json.Marshal(message{Sender: player{Name: "Server", Pos: position{[2]int{0,0}, "test"}}, Text: "Perpetuan<br>Choose starting region:<br>" + strings.Join(settings, " ")})
		if err != nil {
			log.Println(err)
			return err
		}
		message, _ := json.Marshal(message{
			Sender: player{},
			Text: "TESTING STICKY",
		})
		server.BroadcastToRoom("/", s.ID(), "pri", string(message))
		return nil
	})

	server.OnEvent("/", "msg", func(s socketio.Conn, receive message){
		receive.Sender = clientMap[s.ID()].session.Player
		server.BroadcastToRoom("/", clientMap[s.ID()].session.Room, "msg", receive) // does this emit globally or just to this connection
	})

	server.OnEvent("/", "pri", func(s socketio.Conn, receive message){
		receive.Sender = clientMap[s.ID()].session.Player
		msg, _ := json.Marshal(receive)
		server.BroadcastToRoom("", s.ID(), "msg", string(msg))
		server.BroadcastToRoom("", s.ID(), "unsticky")
		clientMap[s.ID()].priority <- receive
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
			log.Println("closed", reason)
	})
	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./api/assets")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

