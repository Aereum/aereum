package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func handleAPI(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success")
	go func() {
		defer conn.Close()
		post := NewPostStage{
			Action:    "NewStagePost",
			Stage:     "Aereum",
			Author:    "Ruben",
			TimeStamp: "17:34:12",
			Content:   "Minha primeira mensagem.",
		}
		wsutil.WriteServerText(conn, post.ToJSON())
		post = NewPostStage{
			Action:    "NewStagePost",
			Stage:     "Aereum",
			Author:    "Larissa",
			TimeStamp: "17:34:22",
			Content:   "Minha primeira resposta.",
		}
		wsutil.WriteServerText(conn, post.ToJSON())

		for {
			msg, code, err := wsutil.ReadClientData(conn)
			if err != nil {
				fmt.Println(err)
				return
			}
			if code.IsData() {
				var action Action
				if err := json.Unmarshal(msg, &action); err == nil {
					if action.Action == "createNewWallet" {
						wallet := CreateNewWallteBalance()
						wsutil.WriteServerText(conn, wallet.ToJSON())
					}
				} else {
					fmt.Println(err)
				}
			}
		}
	}()
}

type Action struct {
	Action string `json:"action"`
	Values string `json:"values"`
}

type NewWalletBalance struct {
	Action  string `json:"action"`
	Token   string `json:"token"`
	Balance string `json:"balance"`
	Hash    string `json:"hash"`
}

type NewPostStage struct {
	Action    string `json:"action"`
	Stage     string `json:"stage"`
	Author    string `json:"author"`
	TimeStamp string `json:"timestamp"`
	Content   string `json:"content"`
}

func (w NewPostStage) ToJSON() []byte {
	data, err := json.Marshal(w)
	if err != nil {
		return []byte(`{"error": true}`)
	}
	return data
}

func CreateNewWallteBalance() NewWalletBalance {
	wallet := NewWalletBalance{
		Action: "NewWalletBalance",
	}
	token := make([]byte, 32)
	rand.Read(token)
	wallet.Token = fmt.Sprintf("0x%v...", hex.EncodeToString(token)[0:20])
	wallet.Balance = "20,934"
	hash := sha256.Sum256(token)
	wallet.Hash = hex.EncodeToString(hash[:])
	return wallet
}

func (w NewWalletBalance) ToJSON() []byte {
	data, err := json.Marshal(w)
	if err != nil {
		return []byte(`{"error": true}`)
	}
	return data
}

func main() {
	/*http.Handle("/ws", http.HandlerFunc(handleAPI))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	err := http.ListenAndServe(":7000", nil)
	if err != nil {
		log.Fatal(err)
	}*/
	listener := readPlays()
	if listener == nil {
		log.Fatal("could not play sheakespeare")
	}
	for {
		instr := listener.Recieve()
		content := instr.Text
		fmt.Println(string(content))
	}
}
