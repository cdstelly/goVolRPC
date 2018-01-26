package main

import (
	"fmt"
	"net/rpc"
	"net"
	"log"
	"net/http"
	"strconv"
	"os/exec"
	"bytes"
)

type NugArg struct {
	TheData []byte
	Inode   string //this probably isn't necessary, but keeping NugArg uniform
}

type NugVol struct {
	SavedData []byte
	PathToFile string
}

func (nd *NugVol) GetDataLen(dataArg *NugArg, reply *string) error {
	*reply = strconv.Itoa(len(nd.SavedData))
	return nil
}

func (nd *NugVol) LoadData(dataArg *NugArg, reply *string) error {
	nd.SavedData = dataArg.TheData

	//we're going to hack it for now
	nd.PathToFile = "/space/m57/pat-2009-12-03.mddramimage"

	*reply = "done"
	return nil
}

func (nd *NugVol) PSList(dataArg *NugArg, reply *string) error {
	//todo: log the hash of the tool used? at least get version information
	pathToTool := "/usr/bin/volatility"
	cmd := exec.Command(pathToTool, "-f", nd.PathToFile, "pslist")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	fmt.Println(out.String())
	*reply = out.String()

	if err != nil {
		fmt.Println(err)
	}
	return err
}

func main() {
	fmt.Println("started")
	tsk := new(NugVol)
	rpc.Register(tsk)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":2002")
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	http.Serve(l,nil) //won't pass here without an error
	fmt.Println("done")
}
