package main

/*
   #include <stdlib.h>
*/
import "C"

import (
	"encoding/base32"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/btcsuite/btcd/btcec"
	"github.com/getamis/alice/crypto/tss/dkg"
	"github.com/getamis/alice/internal/message/types"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var instance *ActiveSessions
var once sync.Once

type DkgSession struct {
	id         string
	nodeid     string
	peers      map[string]*PeerInfo
	dkg        *dkg.DKG
	done       chan struct{}
	threshold  int
	rank       int
	messageBox map[string]*dkg.Message
	result     *DKGResult
}

type PeerInfo struct {
	Id   string `json:"id"`
	Rank uint32 `json:"rank"`
}

type ActiveSessions struct {
	sessions map[string]*DkgSession
	pipeFile *os.File
}

type DKGResult struct {
	Share  string        `json:"share"`
	Pubkey Pubkey        `json:"pubkey"`
	BKs    map[string]BK `json:"bks"`
}

type Pubkey struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type BK struct {
	X    string `json:"x"`
	Rank uint32 `json:"rank"`
}

func getActiveSessions() *ActiveSessions {
	//singleton
	once.Do(func() {
		instance = &ActiveSessions{
			sessions: make(map[string]*DkgSession),
		}
	})
	return instance
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
	return p.nodeid
}

func (p *DkgSession) MustSend(peerid string, message interface{}) {
	println("MustSend........................" + peerid)

	msg, _ := message.(proto.Message)
	bs, err := proto.Marshal(msg)
	// probably unncessary
	x := &dkg.Message{}
	proto.Unmarshal(bs, x)
	msgId := p.getMessageId(peerid, x)
	println(msgId)
	p.messageBox[msgId] = x
	if err != nil {
		println("Cannot marshal message : " + err.Error())
	} else {
		getActiveSessions().pipeFile.WriteString("dkground:" + p.id + ":" + peerid + ":" + msgId + "\n")
	}
}

func (p *DkgSession) getMessageId(peerid string, msg *dkg.Message) string {
	data := []byte(msg.Id + "_" + peerid + "_" + msg.Type.String())
	return base32.StdEncoding.EncodeToString(data)
}

func (p *DkgSession) OnStateChanged(oldState types.MainState, newState types.MainState) {
	println("OnStateChanged........................" + getActiveSessions().pipeFile.Name())

	if newState == types.StateFailed {
		println("Dkg failed", "old", oldState.String(), "new", newState.String())
		close(p.done)
		return
	} else if newState == types.StateDone {
		println("Dkg done", "old", oldState.String(), "new", newState.String())
		result, err := p.dkg.GetResult()
		var status = "result_success"
		if err == nil {
			p.fetchDKGResult(result)
		} else {
			println("Failed to get result from DKG", "err", err)
			status = "result_failure"
		}
		close(p.done)
		p.dkg.Stop()
		getActiveSessions().pipeFile.WriteString(status + ":" + p.id + "\n")
		return
	}
	println("State changed", "old", oldState.String(), "new", newState.String())
}

func (p *DkgSession) fetchDKGResult(result *dkg.Result) {
	dkgResult := &DKGResult{
		Share: result.Share.String(),
		Pubkey: Pubkey{
			X: result.PublicKey.GetX().String(),
			Y: result.PublicKey.GetY().String(),
		},
		BKs: make(map[string]BK),
	}
	for peerID, bk := range result.Bks {
		dkgResult.BKs[peerID] = BK{
			X:    bk.GetX().String(),
			Rank: bk.GetRank(),
		}
	}
	p.result = dkgResult
}

//export NewDkgSession
func NewDkgSession(s *C.char, n *C.char, threshold int, rank int, jsoncstr *C.char) {
	sessionid := deepCopy(C.GoString(s))
	peerJsonStr := deepCopy(C.GoString(jsoncstr))
	nodeid := deepCopy(C.GoString(n))

	println("NewDkgSession........................" + sessionid)
	println(peerJsonStr)

	defer C.free(unsafe.Pointer(s))
	defer C.free(unsafe.Pointer(jsoncstr))
	defer C.free(unsafe.Pointer(n))

	session := &DkgSession{
		id:         sessionid,
		done:       make(chan struct{}),
		threshold:  threshold,
		rank:       rank,
		peers:      make(map[string]*PeerInfo),
		messageBox: make(map[string]*dkg.Message),
		nodeid:     nodeid,
	}
	getActiveSessions().sessions[sessionid] = session
	// Unmarshall
	var pi []PeerInfo
	json.Unmarshal([]byte(peerJsonStr), &pi)
	for i := 0; i < len(pi); i++ {
		println(pi[i].Id + " - " + string(pi[i].Rank))
		session.peers[pi[i].Id] = &pi[i]
	}
	dkg, err := dkg.NewDKG(btcec.S256(), session, uint32(session.threshold), uint32(session.rank), session)
	if err != nil {
		println(err.Error())
	}
	session.dkg = dkg

}

func deepCopy(s string) string {
	var sb strings.Builder
	sb.WriteString(s)
	return sb.String()
}

//export AddPeer
func AddPeer(s *C.char, p *C.char, peerrank int) {

	sessionid := deepCopy(C.GoString(s))
	println("AddPeer........................" + sessionid)

	println(sessionid)

	peerid := deepCopy(C.GoString(p))
	defer C.free(unsafe.Pointer(s))
	defer C.free(unsafe.Pointer(p))

	getActiveSessions().sessions[sessionid].peers[peerid] = &PeerInfo{
		Id:   peerid,
		Rank: uint32(peerrank),
	}
	println("addPeer")
}

//export Initialize
func Initialize(s *C.char) int {
	pipeFile := deepCopy(C.GoString(s))
	println("Initialize........................" + pipeFile)

	defer C.free(unsafe.Pointer(s))
	os.Remove(pipeFile)
	syscall.Mkfifo(pipeFile, 0666)

	f, err := os.OpenFile(pipeFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		println("error opening file: " + err.Error())
		return 0
	}
	getActiveSessions().pipeFile = f
	if err != nil {
		println("Make named pipe file error:", err.Error())
		return 0
	}
	return 1
}

//export GetMessage
func GetMessage(s *C.char, m *C.char) *C.char {
	println("GetMessage........................" + getActiveSessions().pipeFile.Name())

	sessionid := deepCopy(C.GoString(s))
	defer C.free(unsafe.Pointer(s))

	msgid := deepCopy(C.GoString(m))
	defer C.free(unsafe.Pointer(m))

	bs := getActiveSessions().sessions[sessionid].messageBox[msgid]
	println(bs)
	jsonstr, _ := protojson.Marshal(bs)
	println(string(jsonstr))

	jsoncstr := C.CString(string(jsonstr))
	// defer C.free(unsafe.Pointer(jsoncstr))

	return jsoncstr
}

//export HandleMessage
func HandleMessage(s *C.char, m *C.char) *C.char {
	println("HandleMessage........................" + getActiveSessions().pipeFile.Name())

	sessionid := deepCopy(C.GoString(s))
	defer C.free(unsafe.Pointer(s))
	println(sessionid)

	body := deepCopy(C.GoString(m))
	defer C.free(unsafe.Pointer(m))
	println(body)
	// handle data
	x := &dkg.Message{}

	// unmarshal it
	err := protojson.Unmarshal([]byte(body), x)
	if err != nil {
		println("Cannot unmarshal data", "err", err.Error())
		errstr := C.CString(err.Error())
		defer C.free(unsafe.Pointer(errstr))
		return errstr
	}
	printsessions()
	err = getActiveSessions().sessions[sessionid].dkg.AddMessage(x)

	if err != nil {
		println("Cannot add message to DKG", "err", err)
		errstr := C.CString(err.Error())
		defer C.free(unsafe.Pointer(errstr))
		return errstr
	}
	println("Added message to DKG")
	errstr := C.CString("")
	defer C.free(unsafe.Pointer(errstr))
	return errstr
}

func printsessions() {
	sessions := getActiveSessions().sessions
	for k, _ := range sessions {
		println("-----" + k)
	}

}

//export InitDkg
func InitDkg(s *C.char) {
	sessionid := deepCopy(C.GoString(s))
	println("HandleMessage........................" + sessionid)
	session := getActiveSessions().sessions[sessionid]
	defer C.free(unsafe.Pointer(s))
	dkg, err := dkg.NewDKG(btcec.S256(), session, uint32(session.threshold), uint32(session.rank), session)
	if err != nil {
		println(err.Error())
	}
	session.dkg = dkg
}

//export GenerateKey
func GenerateKey(s *C.char) {
	sessionid := deepCopy(C.GoString(s))
	session := getActiveSessions().sessions[sessionid]
	defer C.free(unsafe.Pointer(s))
	println("Generate Key........................" + getActiveSessions().pipeFile.Name())
	// 1. Start a DKG process.
	session.dkg.Start()
}

//export GetResult
func GetResult(s *C.char, m *C.char) *C.char {
	sessionid := deepCopy(C.GoString(s))
	session := getActiveSessions().sessions[sessionid]
	defer C.free(unsafe.Pointer(s))
	if session.result == nil {
		return C.CString(string("{}"))
	}
	jsonbytes, err := json.Marshal(session.result)
	if err != nil {
		return C.CString(string("{}"))
	}
	return C.CString(string(jsonbytes))

}

func main() {

}
