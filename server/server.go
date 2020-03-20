package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var file, _ = os.OpenFile("messages.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

func Run() {

	viper.SetConfigFile("server_config.cfg")
	viper.SetConfigType("toml")
	viper.AddConfigPath("../")
	// Viper read config
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalln("Failed to", err)
	}
	viper.WatchConfig()
	// Re-read config on config change
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			log.Fatalln("Failed to", err)
		}
	})

	///Server
	listener, err := net.Listen("tcp", "0.0.0.0:6969")
	fmt.Println("Stated Server")
	if err != nil {
		log.Fatal("tcp server listener error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("tcp server accept error", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		conn.Close()
		log.Println(conn.RemoteAddr().String() + " has been disconnected ")
		return
	}

	message := string(bufferBytes)
	if len(strings.TrimSpace(message)) > 1 {
		clientAddr := conn.RemoteAddr().String()
		response := fmt.Sprintf(clientAddr + " " + strings.TrimSpace(message) + "\n")
		//log.Println(response)
		fmt.Printf("%v", response)
		_, err := file.WriteString(response)
		if err != nil {
			log.Println(err)
		}
	}

	handleConnection(conn)
}
