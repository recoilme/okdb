package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/recoilme/slowpoke"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server represents the gRPC server
type Server struct {
}

// SayOk generates response ok to a Ping request
func (s *Server) SayOk(ctx context.Context, in *Empty) (*Ok, error) {
	//log.Printf("Receive message")
	return &Ok{Message: "ok"}, nil
}

// Set store key and value in file
func (s *Server) Set(ctx context.Context, cmdSet *CmdSet) (*Empty, error) {
	err := slowpoke.Set(cmdSet.File, cmdSet.Key, cmdSet.Val)
	if err != nil {
		return &Empty{}, status.Errorf(codes.Unknown, err.Error())
	}
	return &Empty{}, nil
}

// Get get value by key
func (s *Server) Get(ctx context.Context, cmdGet *CmdGet) (*ResBytes, error) {
	bytes, err := slowpoke.Get(cmdGet.File, cmdGet.Key)
	if err != nil {
		return &ResBytes{}, status.Errorf(codes.Unknown, err.Error())
	}

	return &ResBytes{Bytes: bytes}, nil
}

// Sets - write key/value pairs -  return error if any
func (s *Server) Sets(ctx context.Context, cmdSets *CmdSets) (*Empty, error) {
	err := slowpoke.Sets(cmdSets.File, cmdSets.Keys)
	if err != nil {
		return &Empty{}, status.Errorf(codes.Unknown, err.Error())
	}
	return &Empty{}, nil
}

// Keys return keys from file
func (s *Server) Keys(ctx context.Context, cmdKeys *CmdKeys) (*ResKeys, error) {
	b, err := slowpoke.Keys(cmdKeys.File, cmdKeys.From, cmdKeys.Limit, cmdKeys.Offset, cmdKeys.Asc)
	if err != nil {
		return &ResKeys{}, status.Errorf(codes.Unknown, err.Error())
	}
	return &ResKeys{Keys: b}, nil
}

// Gets return key/value pairs
func (s *Server) Gets(ctx context.Context, cmdGets *CmdGets) (*ResPairs, error) {
	b := slowpoke.Gets(cmdGets.File, cmdGets.Keys)
	//	slowpoke.Delete()
	return &ResPairs{Pairs: b}, nil
}

// Delete key and val by key
func (s *Server) Delete(ctx context.Context, cmdDel *CmdDel) (*ResDel, error) {
	res, err := slowpoke.Delete(cmdDel.File, cmdDel.Key)
	if err != nil {
		return &ResDel{Deleted: res}, status.Errorf(codes.Unknown, err.Error())
	}
	return &ResDel{Deleted: res}, nil
}

// DeleteFile delete file by name
func (s *Server) DeleteFile(ctx context.Context, cmdDelFile *CmdDelFile) (*Empty, error) {
	err := slowpoke.DeleteFile(cmdDelFile.File)
	if err != nil {
		return &Empty{}, status.Errorf(codes.Unknown, err.Error())
	}
	return &Empty{}, nil
}

// BoltAPI contains handler for rest api to boltdb
func Parser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	method := r.Method
	urlPart := strings.Split(r.URL.Path, "/")
	var err error
	var database = ""
	var bucketstr = ""
	var keystr = ""
	if len(urlPart) == 4 {
		database = urlPart[1]
		bucketstr = urlPart[2]
		keystr = urlPart[3]
	}
	if len(urlPart) == 3 {
		database = urlPart[1]
		bucketstr = urlPart[2]
	}
	switch method {
	case "GET":

		val := get(database, bucketstr, keystr)
		if len(val) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Write(val)
		}
		return
	case "PUT":

		v, err := ioutil.ReadAll(r.Body)
		if err == nil {
			err = put(database, bucketstr, keystr, v)
		}
		if err != nil {
			http.Error(w, err.Error(), 204)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case "DELETE":

		err = delete(database, bucketstr, keystr)
		if err != nil {
			http.Error(w, err.Error(), 204)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case "POST":
		cnt := r.URL.Query().Get("cnt")
		ofs := r.URL.Query().Get("offset")
		var order = r.URL.Query().Get("order")
		var max = 100000
		var offset = 0
		var vals = r.URL.Query().Get("vals")
		m, e := strconv.Atoi(cnt)
		if e == nil {
			max = m
		}
		o, eo := strconv.Atoi(ofs)
		if eo == nil {
			offset = o
		}
		val, err := post(database, bucketstr, keystr, order, vals, max, offset)
		if err != nil {
			http.Error(w, err.Error(), 204)
			return
		}
		w.Write(val)
		return
	default:
		w.WriteHeader(404)
		return
	}
}

func post(database, bucketstr, keystr, order, vals string, max, offset int) ([]byte, error) {
	var err error
	var buffer bytes.Buffer
	switch database {
	case "slowpoke":
		var k []byte
		var asc = true
		if keystr != "Last" && keystr != "" && keystr != "First" {
			k = []byte(keystr)
		}

		if order == "" || order == "desc" {
			asc = false
		}
		keys, err := slowpoke.Keys(bucketstr, k, uint32(max), uint32(offset), asc)
		if vals == "true" {
			pairs := slowpoke.Gets(bucketstr, keys)
			buffer.WriteString("[")
			for i, val := range pairs {
				if i%2 == 0 {
					//key
					if i != 0 {
						buffer.WriteString(",")
					}
					buffer.WriteString(fmt.Sprintf("{\"%s\":", val))
				} else {
					//val
					buffer.WriteString(fmt.Sprintf("%s}", val))
				}

			}
			buffer.WriteString("]")
		} else {
			if err == nil {
				buffer.WriteString("[")
				for i, val := range keys {
					if i != 0 {
						buffer.WriteString(",")
					}
					buffer.WriteString(fmt.Sprintf("\"%s\"", val))

				}
				buffer.WriteString("]")
			}
		}
	}
	return buffer.Bytes(), err
}

func get(database, bucketstr, keystr string) []byte {
	var v []byte
	switch database {

	case "slowpoke":
		v, _ = slowpoke.Get(bucketstr, []byte(keystr))
	}
	return v
}

func put(database, bucketstr, keystr string, val []byte) (err error) {
	switch database {

	case "slowpoke":
		return slowpoke.Set(bucketstr, []byte(keystr), val)
	}
	return err
}

func delete(database, bucketstr, keystr string) (err error) {
	switch database {

	case "slowpoke":
		_, err = slowpoke.Delete(bucketstr, []byte(keystr))
	}
	return err
}
