package grpclangserv

import (
	"bufio"
	"fmt"
	"go/build"
	"log"
	"os"

	pb "github.com/colek42/grpclangserv/api"
	"github.com/colek42/guru"
	"github.com/pkg/errors"
)

/*
type Cfg struct {
	GoPath string
}

func GetCfg() (cfg Cfg) {
	viper.SetEnvPrefix("SAGE")

	viper.AutomaticEnv()
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatalln("GOPATH must be set.")
	}
	return cfg
}
*/

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

func getPath(pkg string, fn string) string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatalln("GOPATH not set")
	}
	posn := fmt.Sprintf("%s/src/%s/%s", gopath, pkg, fn)
	return posn

}

func LineToByteOffset(fn string, lineNumber int32, charNumber int32) (int32, error) {
	byteCount := int32(0)
	//curent line should be 1 here, but for some reason scanner
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
	return 0, errors.New("EOF")

}
