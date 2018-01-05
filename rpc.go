package grpclangserv

import (
	"log"
	"net"

	"golang.org/x/net/context"

	pb "github.com/colek42/grpclangserv/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	rpcEntry = "localhost:4534"
)

type server struct{}

func StartRPCServer() {
	lis, err := net.Listen("tcp", rpcEntry)
	if err != nil {
		err := errors.Wrapf(err, "Unable to listen on %s", rpcEntry)
		log.Println(err)
	}

	s := grpc.NewServer()

	srv := server{}
	pb.RegisterLanguageServerServer(s, srv)

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
