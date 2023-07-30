package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sc, err := stan.Connect("test-cluster", "client2")
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, nats.DefaultURL)
	}
	defer sc.Close()

	fmt.Println("Hole id active!")
	//hole("./hole", sc)
	hole("./utils/publisher/hole", sc)
}

func hole(path string, sc stan.Conn) {
	for {
		dir, err := os.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(2 * time.Second)

		for _, fileInfo := range dir {
			if fileInfo.IsDir() {
				err = exec.Command("rm", "-rf", fileInfo.Name()).Run()
				if err != nil {
					log.Println(err)
				}
			}
			buff, err := os.ReadFile(path + "/" + fileInfo.Name())
			if err != nil {
				log.Println(err)
			}
			err = sc.Publish("updates", buff)
			if err != nil {
				log.Println(err)
			}
			log.Printf("Published [%s] : '%s'\n\n", "updates", string(buff))

			err = exec.Command("rm", path+"/"+fileInfo.Name()).Run()
			if err != nil {
				log.Println(err)
			}

		}
	}
}
