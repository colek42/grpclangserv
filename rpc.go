package grpclangserv

import (
	"log"
	"net"

	"golang.org/x/net/context"

	pb "github.com/colek42/grpclangserv/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type server struct {
	cfg Cfg
}

func StartRPCServer(cfg Cfg) {
	srv := server{
		cfg: cfg,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.RPCListen, cfg.RPCPort)
	if err != nil {
		err := errors.Wrapf(err, "Unable to listen on %s", rpcEntry)
		log.Fatalln(err)
	}

	s := grpc.NewServer()

	srv := server{
		cfg: cfg,
	}

	pb.RegisterLanguageServerServer(s, srv)
	gracefulShutdown()
	s.Serve(lis)
}

func (s server) GetDefinition(c context.Context, q *pb.Query) (*pb.DefResponse, error) {
	path := getPath(q.GetPkg(), q.GetFileName())
	byteOffset, err := LineToByteOffset(path, q.GetLineNumber(), q.GetCharNumber())

	if err != nil {
		return nil, err
	}

	res, err := Query(int(byteOffset), "definition", path)

	return res, err
}

//Catches sigterm/ sigint and performs a graceful shutdown
func gracefulShutdown() {
	var stop = make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)
	go func() {
		sig := <- stop
		fmt.Printf("Caught Sig: %v", sig)
		fmt.Println("Waiting 5 secs to finish jobs")
		time.Sleep(5*time.Second)
		os.Exit(0)
	}()
}
