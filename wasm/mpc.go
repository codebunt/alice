package main

import (
	"strconv"
	"syscall/js"

	"github.com/getamis/alice/crypto/tss/dkg"
	"github.com/getamis/alice/wasm/keygen"

	"google.golang.org/protobuf/encoding/protojson"
)

// Exposed Methods
func NewDkgSession(this js.Value, inputs []js.Value) interface{} {
	sessionid := js.ValueOf(inputs[0]).String()
	nodeid := js.ValueOf(inputs[1]).String()
	threshold := js.ValueOf(inputs[2]).Int()
	rank := js.ValueOf(inputs[3]).Int()
	peerJsonStr := inputs[4]
	callback := inputs[len(inputs)-1:][0]

	println("NewDkgSession........................" + sessionid)

	peers := make(map[string]*keygen.PeerInfo)
	for i := 0; i < peerJsonStr.Length(); i++ {
		item := peerJsonStr.Index(i)
		r := strconv.Itoa(item.Get("rank").Int())
		p := &keygen.PeerInfo{
			Id:   item.Get("id").String(),
			Rank: uint32(item.Get("rank").Int()),
		}
		println(p.Id + " -- " + r)
		peers[p.Id] = p
	}
	keygen.NewDkgSession(sessionid, nodeid, threshold, rank, peers, callback)
	return nil
}

func DkgHandleMessage(this js.Value, inputs []js.Value) interface{} {
	println("HandleMessage........................")
	sessionid := js.ValueOf(inputs[0]).String()
	println(sessionid)
	body := js.ValueOf(inputs[1]).String()
	println(body)
	// handle data
	x := &dkg.Message{}
	// unmarshal it
	err := protojson.Unmarshal([]byte(body), x)
	if err != nil {
		println("Cannot unmarshal data", "err", err.Error())
		return err.Error()
	}
	err = keygen.GetActiveSessions().Sessions[sessionid].Dkg.AddMessage(x)

	if err != nil {
		println("Cannot add message to DKG", "err", err)
		return err.Error()
	}
	println("Added message to DKG")
	return err.Error()
}

func GenerateKey(this js.Value, inputs []js.Value) interface{} {
	sessionid := js.ValueOf(inputs[0]).String()
	callback := inputs[len(inputs)-1:][0]

	session := keygen.GetActiveSessions().Sessions[sessionid]
	session.ResultCallBack = callback
	println("Generate Key........................" + sessionid)
	// 1. Start a DKG process.
	session.Dkg.Start()
	return nil
}

// Expose methods
func registerMethods() {
	js.Global().Set("newDkgSession", js.FuncOf(NewDkgSession))
	js.Global().Set("addDkgMessage", js.FuncOf(DkgHandleMessage))
	js.Global().Set("generateKey", js.FuncOf(GenerateKey))
}

func main() {
	c := make(chan bool)
	registerMethods()
	<-c
}
