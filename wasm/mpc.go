package main

/*
   #include <stdlib.h>
*/

import (
	"reflect"
	"strconv"
	"sync"
	"syscall/js"

	"github.com/getamis/sirius/log"

	"github.com/getamis/alice/crypto/tss/dkg"
	"github.com/getamis/alice/internal/message/types"
	"google.golang.org/protobuf/proto"
)

var instance *ActiveSessions
var once sync.Once
var pipeFile = "dkg.log"

type DkgSession struct {
	id         int
	peers      map[string]*PeerInfo
	dkg        *dkg.DKG
	done       chan struct{}
	threshold  int
	rank       int
	jscallback js.Value
}

type PeerInfo struct {
	id   string
	rank uint32
}

type ActiveSessions struct {
	sessions map[int]*DkgSession
}

func getActiveSessions() *ActiveSessions {
	//singleton
	once.Do(func() {
		instance = &ActiveSessions{
			sessions: make(map[int]*DkgSession),
		}
	})
	return instance
}

// export NewDkgSession
func NewDkgSession(sessionid int, threshold int, rank int) {
	session := &DkgSession{
		id:        sessionid,
		done:      make(chan struct{}),
		threshold: threshold,
		rank:      rank,
		peers:     make(map[string]*PeerInfo),
	}
	getActiveSessions().sessions[sessionid] = session
}

func (p *DkgSession) NumPeers() uint32 {
	return uint32(len(p.peers))
}

func (p *DkgSession) PeerIDs() []string {
	keys := reflect.ValueOf(p.peers).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	return strkeys
}

func (p *DkgSession) SelfID() string {
	return strconv.Itoa(p.id)
}

func (p *DkgSession) MustSend(id string, message proto.Message) {
	msg, _ := message.(proto.Message)
	bs, _ := proto.Marshal(msg)
	arrayConstructor := js.Global().Get("Uint8Array")
	dataJS := arrayConstructor.New(len(bs))
	js.CopyBytesToJS(dataJS, bs)
	p.jscallback.Invoke(js.Null(), dataJS)
}

func (p *DkgSession) OnStateChanged(oldState types.MainState, newState types.MainState) {

}

func registerMethods() {
	js.Global().Set("newSession", js.FuncOf(newSession))
	js.Global().Set("addPeer", js.FuncOf(addPeer))
	js.Global().Set("run", js.FuncOf(run))
	js.Global().Set("registerCallback", js.FuncOf(registerCallback))
}

func newSession(this js.Value, i []js.Value) interface{} {
	sessionid := js.ValueOf(i[0]).Int()
	threshold := js.ValueOf(i[1]).Int()
	rank := js.ValueOf(i[2]).Int()
	NewDkgSession(sessionid, threshold, rank)
	println("newSession")
	return nil
}

func registerCallback(this js.Value, inputs []js.Value) interface{} {
	sessionid := js.ValueOf(inputs[0]).Int()
	callback := inputs[len(inputs)-1:][0]
	if callback.Type() == js.TypeFunction {
		getActiveSessions().sessions[sessionid].jscallback = callback
	}
	println("registerCallback")
	callback.Invoke(js.Null(), "Did you say ")
	return nil
}

func addPeer(this js.Value, i []js.Value) interface{} {
	sessionid := js.ValueOf(i[0]).Int()
	peerid := js.ValueOf(i[1]).String()
	peerrank := js.ValueOf(i[2]).Int()
	getActiveSessions().sessions[sessionid].peers[peerid] = &PeerInfo{
		id:   peerid,
		rank: uint32(peerrank),
	}
	println("addPeer")

	return nil
}

func run(this js.Value, i []js.Value) interface{} {
	sessionid := js.ValueOf(i[0]).Int()
	dkg := getActiveSessions().sessions[sessionid].dkg
	dkg.Start()
	defer dkg.Stop()
	<-getActiveSessions().sessions[sessionid].done
	println("run")
	return nil
}

func receiveMsg(this js.Value, params []js.Value) interface{} {
	sessionid := js.ValueOf(params[0]).Int()
	msg := make([]byte, params[1].Length())

	js.CopyBytesToGo(msg, params[1])

	d := getActiveSessions().sessions[sessionid].dkg
	x := &dkg.Message{}
	proto.Unmarshal(msg, x)
	err := d.AddMessage(x)
	if err != nil {
		log.Warn("Cannot add message to DKG", "err", err)
	}
	println("receiveMsg")
	return nil
}

func main() {
	c := make(chan bool)
	registerMethods()
	<-c
}
