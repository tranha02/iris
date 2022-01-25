package main

import (
	"context"
	"flag"
	"net"
	"os"

	logging "github.com/ipfs/go-log/v2"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"happystoic/p2pnetwork/pkg/config"
	"happystoic/p2pnetwork/pkg/node"
)

var log = logging.Logger("p2pnetwork")

// TODO: signal handling
// TODO: search for buddies from the same organisation in peer discovery
// TODO: add rate limitting/queue for each peer
// TODO: make a goroutine to purge message cache keys after some time
// TODO: finish peer discovery
// TODO: create connection manager (using reliability and some other smart algorithms)
// TODO: deal with NAT
// TODO: add intelligence msg (be careful, id of original msg has to be present when forwarding and checking if peer has already seen that)
// TODO: add redis api
// TODO: add network notiffee and send update to TL through redis when peers connecgt/disconnect
// TODO: use msgs timestamp to throw away old messages?
// TODO: can you spoof QUIC sender? If not, tell TL if stream is corrupted and cannot be deserialized so the IP can be punished

func loadConfig() (*config.Config, error) {
	var c config.Config

	configFile := flag.String("conf", "", "path to configuration file")
	flag.Parse()

	if configFile == nil || *configFile == "" {
		return nil, errors.New("missing path of configuration file")
	} else {
	}
	viper.SetConfigFile(*configFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, c.Check()
}

func checkUDPPortAvailability(port uint32) error {
	ln, err := net.ListenUDP("udp", &net.UDPAddr{Port: int(port)})
	if err != nil {
		return err
	}
	_ = ln.Close()
	return nil
}

func main() {
	// load configuration
	conf, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// check if port is available
	err = checkUDPPortAvailability(conf.Server.Port)
	if err != nil {
		log.Fatal(err)
	}

	// create p2p node
	ctx, _ := context.WithCancel(context.Background())
	localNode, err := node.NewNode(conf, ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("created node with ID: %s", localNode.ID())
	for _, addr := range localNode.Addrs() {
		log.Infof("connection string: '%s %s'", addr, localNode.ID())
	}

	// connect node to the network
	localNode.ConnectToInitPeers()

	// temporary playground
	doSomething := len(os.Getenv("DO_SOMETHING")) > 0
	localNode.Start(doSomething)

	log.Infof("finished, program terminating...")
}
