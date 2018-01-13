package grpclangserv

import (
	"reflect"
	"testing"

	pb "github.com/colek42/grpclangserv/api"
	"github.com/davecgh/go-spew/spew"
)

func TestQuery(t *testing.T) {
	type args struct {
		byteOffset int
		mode       string
		file       string
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.DefResponse
		wantErr bool
	}{
		{
			name: "package",
			args: args{
				byteOffset: 1791,
				mode:       "definition",
				file:       "test_fixtures/main.go",
			},
			want: &pb.DefResponse{
				Name: "Println",
				Type: "func(a ...interface{}) (n int, err error)",
				Pkg:  "fmt",
				Position: &pb.Position{
					FileName: "/usr/local/go/src/fmt/print.go",
					Offset:   7393,
					Line:     256,
					Column:   6,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Query(tt.args.byteOffset, tt.args.mode, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query() = %v, want %s", got, spew.Sdump(tt.want))
			}
		})
	}
}

func TestLineToByteOffset(t *testing.T) {
	type args struct {
		fn         string
		lineNumber int32
		charNumber int32
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		{
			name: "First Char",
			args: args{
				fn:         "test_fixtures/main.go",
				lineNumber: 1,
				charNumber: 1,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Line 68, Println",
			args: args{
				fn:         "test_fixtures/main.go",
				lineNumber: 68,
				charNumber: 7,
			},
			want:    1791,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LineToByteOffset(tt.args.fn, tt.args.lineNumber, tt.args.charNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("LineToByteOffset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LineToByteOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}
