package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

type User struct {
	conn net.Conn
	id   string
}

type Event string

const (
	Auth           Event = "auth"
	Message        Event = "messsage"
	Help           Event = "help"
	Ok             Event = "ok"
	NotOk          Event = "notok"
	ListUserOnline Event = "Listuseronline"
)

type Data struct {
	Event   Event  `json:"event"`
	Message string `json:"message"`
}

var USERS []User

func checkError(err error) bool {
	if err != nil {
		fmt.Println("Error:", err)
		return true
	}
	return false
}

func printError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func isUserIdPresence(id string) bool {
	for _, u := range USERS {
		if u.id == id {
			return true
		}
	}
	return false
}

func clearScreen() {
	// Clear the terminal screen based on the operating system
	cmd := "clear" // Default to Unix-like systems
	if os.Getenv("OS") == "Windows_NT" {
		cmd = "cls"
	}
	if command := exec.Command(cmd); command != nil {
		command.Stdout = os.Stdout
		command.Run()
	}
}

func chatUser(connClient net.Conn) {
	clearScreen()
	fmt.Println("Write his/her name to chat")
	dataSent := Data{Event: ListUserOnline, Message: ""}

	_, err := json.Marshal(dataSent)
	printError(err)

	//waiting response from server
	var dataReceived Data
	decoder := json.NewDecoder(connClient) // Create a new JSON decoder for the connection

	err = decoder.Decode(&dataReceived)
	printError(err)

	var userSlice []string
	if dataReceived.Message == "" {
		fmt.Println("There is no user online")
	} else {
		userSlice = strings.Split(dataReceived.Message, ";")
		for index, userID := range userSlice {
			fmt.Printf("%d. %s\n", index, userID)
		}
	}

}

func displayMenu(connClient net.Conn) {
	clearScreen()
	fmt.Println("Please enter a number")
	fmt.Println("1. Chat to a user")
	fmt.Println("2. Exit")
	var number int

	for {
		_, err := fmt.Scanf("%d", &number)
		if err != nil || number < 1 || number > 2 {
			fmt.Println("Enter correct number")
		} else {
			fmt.Println("number est ", number)
			break
		}
	}

	if number == 1 {
		fmt.Println("on rentre ici number 1")
		chatUser(connClient)
	} else {
		fmt.Println("on rentre ici number 2")
		return
	}
}

func sendAllUserOnline() string {
	var users string
	if len(USERS) == 0 {
		return users
	} else {
		for _, u := range USERS {
			users = users + ";" + u.id
		}
	}
	return users
}
func clienAuth(connClient net.Conn) bool {
	fmt.Println("Enter your id")

	auth := false
	for auth != true {
		var name string
		fmt.Scanf("%s", &name)

		dataSent := Data{Event: Auth, Message: name}

		dataSentJson, err := json.Marshal(dataSent)
		printError(err)

		_, err = connClient.Write([]byte(dataSentJson))
		printError(err)

		fmt.Println("we are sending a message from client", dataSent.Event, dataSent.Message)
		//waiting response from server
		var dataReceived Data
		decoder := json.NewDecoder(connClient) // Create a new JSON decoder for the connection

		err = decoder.Decode(&dataReceived)
		printError(err)

		fmt.Printf("Data received from client: Event=%s, Message=%s\n", dataReceived.Event, dataReceived.Message)

		if dataReceived.Event == Auth {
			if dataReceived.Message == string(Ok) {
				auth = true
			} else {
				fmt.Println("id user is alread used")
				continue
			}
		}
	}
	return true
}

func main() {
	// Listen for incoming TCP connections on port 8080
	listener, err := net.Listen("tcp", ":8080")
	server := true
	var connClient net.Conn
	if err != nil {
		server = false
		fmt.Println("you're a client")

		serverAddr := "127.0.0.1:8080"
		connClient, err = net.Dial("tcp", serverAddr)

		if err != nil {
			fmt.Println("Failed to connect to server:", err)
			return
		}
		defer connClient.Close()

		for {
			if clienAuth(connClient) {
				break
			}
		}
	}

	if server {
		defer listener.Close()
		for {
			// Accept incoming connections
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Failed to accept connection:", err)
				continue
			}

			// Handle the connection concurrently
			go handleConnection(conn)
		}
	} else {

		fmt.Println("we are entering here in client side")

		displayMenu(connClient)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed when the function exits.
	fmt.Println("we are in handleconnection")

	// Loop to continuously read messages from the connection
	for {
		var dataReceived Data
		decoder := json.NewDecoder(conn) // Create a new JSON decoder for the connection

		if err := decoder.Decode(&dataReceived); err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client")
			} else {
				fmt.Println("Error decoding JSON:", err)
			}
			break // Exit the loop on any error
		}
		fmt.Printf("Data received from client: Event=%s, Message=%s\n", dataReceived.Event, dataReceived.Message)

		// fmt.Println("we have received message from user ", dataUnmarshal.Event, dataUnmarshal.Message)
		if dataReceived.Event == Auth {
			var dataSent Data
			if !isUserIdPresence(dataReceived.Message) {
				newUser := User{conn: conn, id: dataReceived.Message}
				USERS = append(USERS, newUser)
				fmt.Println("ok, id can be used")
				dataSent = Data{Event: Auth, Message: string(Ok)}
			} else {
				fmt.Println("ok, id can not be used")
				dataSent = Data{Event: Auth, Message: string(NotOk)}
			}
			dataSentJson, err := json.Marshal(dataSent)
			printError(err)
			_, err = conn.Write([]byte(dataSentJson))
			printError(err)
		} else if dataReceived.Event == ListUserOnline {
			fmt.Println("on rentre ici listuseronline")
			var datasent Data
			datasent = Data{Event: ListUserOnline, Message: sendAllUserOnline()}
			dataSentJson, err := json.Marshal(datasent)
			printError(err)
			_, err = conn.Write([]byte(dataSentJson))
			printError(err)
		}
	}
}
