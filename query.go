package grpclangserv

import (
	"bufio"
	"fmt"
	"go/build"
	"io"
	"log"
	"os"

	pb "github.com/colek42/grpclangserv/api"
	"github.com/colek42/guru"
	"github.com/spf13/viper"
)

type Cfg struct {
	GoPath    string
	RPCPort   string
	RPCListen string
}

func GetCfg() (cfg Cfg) {

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatalln("GOPATH must be set.")
	}
	viper.SetEnvPrefix("LANGSERV")

	viper.AutomaticEnv()
	viper.SetDefault("RPC_PORT", 4534)
	viper.SetDefault("RPC_LISTEN", "0.0.0.0")
	viper.SetDefault("GOPATH", gopath)

	cfg = Cfg{
		GoPath:    viper.GetString("GOPATH"),
		RPCPort:   viper.GetString("RPC_PORT"),
		RPCListen: viper.GetString("RPC_LISTEN"),
	}

	return
}

//Query takes a byte offset, query mode, and full file path and returns a definition response
//Only definition is implimented at this time
func Query(byteOffset int, mode string, file string) (*pb.DefResponse, error) {

	posn := fmt.Sprintf("%s:#%d", file, byteOffset)

	var scope []string
	ctxt := &build.Default
	query := guru.Query{
		Pos:   posn,
		Build: ctxt,
		Scope: scope,
	}

	res, err := guru.Run(mode, &query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	out := &pb.DefResponse{
		Name: res.Obj.Name(),
		Type: res.Obj.Type().String(),
		Pkg:  res.Obj.Pkg().Name(),
		Position: &pb.Position{
			FileName: res.Position.Filename,
			Offset:   int32(res.Position.Offset),
			Line:     int32(res.Position.Line),
			Column:   int32(res.Position.Column),
		},
	}
	return out, nil
}

//getPath returns the full filepath
func getPath(pkg string, fn string) string {
	gopath := viper.GetString("GOPATH")
	posn := fmt.Sprintf("%s/src/%s/%s", gopath, pkg, fn)
	return posn

}

//LineToByteOffset returns a byte offset given a line number and character number for a file.
//Lines are indexed starting at 1
//Character offsets are indexed starting at 1
//Byteoffset is indexed starting at 0
func LineToByteOffset(fn string, lineNumber int32, charNumber int32) (int32, error) {
	byteCount := int32(0)
	currLine := int32(1)

	inputFile, err := os.Open(fn)
	if err != nil {
		return 0, err
	}
	scanner := bufio.NewScanner(inputFile)
	if err != nil {
		return 0, err
	}
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		//count new lines
		if currLine < lineNumber {
			byteCount = byteCount + int32(len(scanner.Text())+1)

			currLine = currLine + 1
		} else {

			return byteCount + charNumber, nil
		}

	}
	return 0, io.EOF

}
