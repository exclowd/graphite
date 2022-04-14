package main

import (
	"flag"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"strings"
)

// The full server encapsulated in a struct
type server struct {
	logger *zap.Logger

	raft *raft.Raft
	fsm  *raftFSM
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	logger.Info("Hello from zap logger")

	id := flag.Int("id", 0, "Id of the cluster")
	port := flag.String("p", "8001", "Port to listen on")
	clusterS := flag.String("cluster", "", "Comma Separated ips of the rafts of current nodes")
	flag.Parse()
	if len(*port) != 4 {
		fmt.Println("Usage server [-p] port ...")
		flag.PrintDefaults()
		os.Exit(1)
	}
	var cluster []string = strings.Split(*clusterS, ",")
	for i := range cluster {
		cluster[i] = strings.Trim(cluster[i], " ")
	}
	logger.Info(fmt.Sprintf("Running Node: %d at port: %s", *id, *port))

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *badger.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	sayHello := func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("<h1>Hello</h1>"))
		if err != nil {
			return
		}
	}

	http.HandleFunc("/", sayHello)

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
