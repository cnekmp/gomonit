package client

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hpcloud/tail"
	"github.com/spf13/viper"
)

const SERVER = "127.0.0.1:6969"

type fileConfig struct {
	Hostname string
	Ip       []string
	Alarm    []string
	Logfiles []string
}

var config = tail.Config{
	ReOpen:    true,
	Follow:    true,
	MustExist: false,
	Poll:      true,
}

var fileConf fileConfig
var t []string

func Run() {
	viper.SetConfigFile("client_config.cfg")
	viper.SetConfigType("toml")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalln("Failed to", err)
		//panic(fmt.Errorf("client_config.cfg file does not exist: %s \n", err))
	}
	if err := viper.Unmarshal(&fileConf); err != nil {
		fmt.Printf("couldn't read config: %s", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			log.Fatalln("Failed to", err)
			//panic(fmt.Errorf("client_config.cfg file does not exist: %s \n", err))
		}
		if err := viper.Unmarshal(&fileConf); err != nil {
			fmt.Printf("couldn't read config: %s", err)
		}
		viper.WatchConfig()
		t = viper.GetStringSlice("logfiles")

	})

	fmt.Println(fileConf.Hostname, fileConf.Logfiles)
	c := make(chan string)

	// Start goroutines for tailing provided files
	//for _, v := range viper.GetStringSlice("logfiles") {
	//go tailFile(v, c, &config)
	//}
	t = viper.GetStringSlice("logfiles")
	go tailFile(&t, c, &config)

	for {
		// Create connection towards SERVER
		connection, err := net.Dial("tcp", SERVER)

		// Send received alarms to SERVER
		for {
			if err != nil {
				time.Sleep(time.Second * time.Duration(5))
				log.Println(err)
				continue
			}
			temp := <-c
			// Custom trigger

			if !alarmsCheck(&temp, &fileConf.Alarm) {
				_, err = connection.Write([]byte(time.Now().Format("20060102150405") + " " + fileConf.Hostname + " !!!!! " + temp + fileConf.Logfiles[0]))
			}
			if err != nil {
				connection = tcpReconnect()
				break
			} else {
				fmt.Fprintf(connection, "\n")

			}
		}
	}
}

// TCP Reconnect on connection issues
func tcpReconnect() net.Conn {
	connection, err := net.Dial("tcp", SERVER)
	if err != nil {
		log.Println("Failed to reconnect:", err.Error())
		time.Sleep(time.Millisecond * time.Duration(2000))
		connection = tcpReconnect()
	}
	return connection
}

// tailFile function used to tail provided files
func tailFile(s *[]string, c chan string, config *tail.Config) {
	for _, x := range *s {
		fmt.Println("Log for tail is: ", x)
		t, err := tail.TailFile(x, *config)
		if err != nil {
			panic(err)
		}
		for line := range t.Lines {
			c <- line.Text
		}
	}
}

// alarmsCheck function is used to filter out alarms
func alarmsCheck(s *string, b *[]string) bool {
	var a bool
	for _, v := range *b {
		// Ignore Alarms from config file
		if strings.TrimSpace(v) == *s {
			a = true
			fmt.Printf("%v", v)
		}
	}
	// Return filtered alarms to server
	return a
}
