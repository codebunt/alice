package main

/*
   #include <stdlib.h>
*/
import "C"

import (
	bufio "bufio"
	"encoding/json"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/codebunt/dart_api_dl"

	"github.com/btcsuite/btcd/btcec"
	"github.com/getamis/alice/crypto/tss/dkg"
	"github.com/getamis/alice/internal/message/types"
	kommons "github.com/getamis/alice/lib/commons"
	ksigner "github.com/getamis/alice/lib/signer"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var instance *ActiveSessions
var once sync.Once

type DkgSession struct {
	id           string
	nodeid       string
	peers        map[string]*PeerInfo
	dkg          *dkg.DKG
	done         chan struct{}
	threshold    int
	rank         int
	messageBox   map[string]*dkg.Message
	incomingMsgs map[string]*dkg.Message
	result       *kommons.DKGResult
}

type PeerInfo struct {
	Id   string `json:"id"`
	Rank uint32 `json:"rank"`
}

type ActiveSessions struct {
	sessions     map[string]*DkgSession
	pipeFile     *os.File
	callbackType string
	port         int64
	logfile      *os.File
	tcpwriter    *bufio.Writer
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
	msg, _ := message.(proto.Message)
	bs, err := proto.Marshal(msg)
	// probably unncessary
	x := &dkg.Message{}
	proto.Unmarshal(bs, x)
	msgId := p.getMessageId(peerid, x)
	println(x.Type)
	p.messageBox[msgId] = x
	if err != nil {
		println("Cannot marshal message : " + err.Error())
	} else {
		if getActiveSessions().callbackType == "PIPE" {
			getActiveSessions().pipeFile.WriteString("dkground:" + p.id + ":" + peerid + ":" + msgId + "\n")
		}
		if getActiveSessions().callbackType == "TCP" {
			getActiveSessions().tcpwriter.WriteString("dkground:" + p.id + ":" + peerid + ":" + msgId + "\n")
			getActiveSessions().tcpwriter.Flush()
		}
		if getActiveSessions().callbackType == "DART_PORT" {
			dart_api_dl.SendToPort(getActiveSessions().port, "dkground", p.id, peerid, msgId)
		}

	}
}

func (p *DkgSession) getMessageId(peerid string, msg *dkg.Message) string {
	return msg.Type.String() + "_" + peerid
}

func (p *DkgSession) OnStateChanged(oldState types.MainState, newState types.MainState) {
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
		if getActiveSessions().callbackType == "PIPE" {
			getActiveSessions().pipeFile.WriteString(status + ":" + p.id + "\n")
		}
		if getActiveSessions().callbackType == "TCP" {
			getActiveSessions().tcpwriter.WriteString(status + ":" + p.id + "\n")
			getActiveSessions().tcpwriter.Flush()
		}
		if getActiveSessions().callbackType == "DART_PORT" {
			dart_api_dl.SendToPort(getActiveSessions().port, status, p.id, "", "")
		}

		return
	}
	println("State changed", "old", oldState.String(), "new", newState.String())
}

func (p *DkgSession) fetchDKGResult(result *dkg.Result) {
	dkgResult := &kommons.DKGResult{
		Share: result.Share.String(),
		Pubkey: kommons.Pubkey{
			X: result.PublicKey.GetX().String(),
			Y: result.PublicKey.GetY().String(),
		},
		BKs: make(map[string]kommons.BK),
	}
	for peerID, bk := range result.Bks {
		dkgResult.BKs[peerID] = kommons.BK{
			X:    bk.GetX().String(),
			Rank: bk.GetRank(),
		}
	}
	p.result = dkgResult
}

// Exposed Methods
// DISTRIBUTED KEY GEN
//export NewDkgSession
func NewDkgSession(s *C.char, n *C.char, threshold int, rank int, jsoncstr *C.char) {
	sessionid := deepCopy(C.GoString(s))
	peerJsonStr := deepCopy(C.GoString(jsoncstr))
	nodeid := deepCopy(C.GoString(n))

	println("NewDkgSession........................" + sessionid)
	println(peerJsonStr)

	// C.free(unsafe.Pointer(s))
	// C.free(unsafe.Pointer(jsoncstr))
	// C.free(unsafe.Pointer(n))

	session := &DkgSession{
		id:           sessionid,
		done:         make(chan struct{}),
		threshold:    threshold,
		rank:         rank,
		peers:        make(map[string]*PeerInfo),
		messageBox:   make(map[string]*dkg.Message),
		incomingMsgs: make(map[string]*dkg.Message),
		nodeid:       nodeid,
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
	// C.free(unsafe.Pointer(s))
	// C.free(unsafe.Pointer(p))

	getActiveSessions().sessions[sessionid].peers[peerid] = &PeerInfo{
		Id:   peerid,
		Rank: uint32(peerrank),
	}
	println("addPeer")
}

//export Initialize
func Initialize(port int) int {
	println("Initialize............." + strconv.Itoa(port))
	c, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		println("client, Dial" + err.Error())
		return 0
	}
	getActiveSessions().tcpwriter = bufio.NewWriter(c)
	getActiveSessions().callbackType = "TCP"
	ksigner.InitializeWithTCPConnection(bufio.NewWriter(c))
	println("Initialize........................ Done")
	return 1
}

//export InitializeDL
func InitializeDL(api unsafe.Pointer, port int64) int {
	getActiveSessions().callbackType = "DART_PORT"
	dart_api_dl.Init(api)
	getActiveSessions().port = port
	return 1
}

//export GetMessage
func GetMessage(s *C.char, m *C.char) *C.char {
	println("GetMessage........................" + getActiveSessions().callbackType)

	sessionid := deepCopy(C.GoString(s))
	// C.free(unsafe.Pointer(s))

	msgid := deepCopy(C.GoString(m))
	// C.free(unsafe.Pointer(m))

	bs := getActiveSessions().sessions[sessionid].messageBox[msgid]
	println(bs.Type)
	mOptions := protojson.MarshalOptions{
		Multiline:       true,
		EmitUnpopulated: true,
	}
	jsonstr, _ := mOptions.Marshal(bs)

	println(string(jsonstr))

	jsoncstr := C.CString(string(jsonstr))
	// // C.free(unsafe.Pointer(jsoncstr))
	println("end GetMessage........................" + getActiveSessions().callbackType)

	return jsoncstr
}

//export HandleMessage
func HandleMessage(s *C.char, m *C.char) *C.char {
	println("HandleMessage........................" + getActiveSessions().callbackType)

	sessionid := deepCopy(C.GoString(s))
	// C.free(unsafe.Pointer(s))
	println(sessionid)

	body := deepCopy(C.GoString(m))
	// C.free(unsafe.Pointer(m))
	println(body)
	// handle data
	x := &dkg.Message{}

	// unmarshal it
	err := protojson.Unmarshal([]byte(body), x)
	if err != nil {
		println("Cannot unmarshal data", "err", err.Error())
		errstr := C.CString(err.Error())
		// C.free(unsafe.Pointer(errstr))
		return errstr
	}
	println(x.Type.String() + "_" + x.Id)
	if getActiveSessions().sessions[sessionid].incomingMsgs[x.Type.String()+"_"+x.Id] != nil {
		errstr := C.CString("Already handled \n" + x.Type.String() + "_" + x.Id)
		// C.free(unsafe.Pointer(errstr))
		return errstr
	}
	printsessions()
	err = getActiveSessions().sessions[sessionid].dkg.AddMessage(x)
	getActiveSessions().sessions[sessionid].incomingMsgs[x.Type.String()+"_"+x.Id] = x

	if err != nil {
		println("Cannot add message to DKG", "err", err)
		errstr := C.CString(err.Error())
		// C.free(unsafe.Pointer(errstr))
		return errstr
	}
	println("Added message to DKG")
	errstr := C.CString("")
	// C.free(unsafe.Pointer(errstr))
	return errstr
}

func printsessions() {
	sessions := getActiveSessions().sessions
	for k, _ := range sessions {
		println("-----" + k)
	}

}

//export GenerateKey
func GenerateKey(s *C.char) {
	sessionid := deepCopy(C.GoString(s))
	session := getActiveSessions().sessions[sessionid]
	// C.free(unsafe.Pointer(s))
	println("Generate Key........................" + getActiveSessions().callbackType)
	// 1. Start a DKG process.
	session.dkg.Start()
}

//export GetResult
func GetResult(s *C.char) *C.char {
	sessionid := deepCopy(C.GoString(s))
	session := getActiveSessions().sessions[sessionid]
	// C.free(unsafe.Pointer(s))
	if session.result == nil {
		return C.CString(string("{}"))
	}
	jsonbytes, err := json.Marshal(session.result)
	if err != nil {
		return C.CString(string("{}"))
	}
	return C.CString(string(jsonbytes))
}

// SIGNER
//export InitializeSignerWithPipe
func InitializeSignerWithPipe(s *C.char) int {
	pipeFile := deepCopy(C.GoString(s))
	println("Initialize........................" + pipeFile)

	// C.free(unsafe.Pointer(s))
	// os.Remove(pipeFile)
	// syscall.Mkfifo(pipeFile, 0666)

	f, err := os.OpenFile(pipeFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		println("error opening file: " + err.Error())
		return 0
	}
	if err != nil {
		println("Make named pipe file error:", err.Error())
		return 0
	}
	ksigner.InitializeWithPipe(f)
	return 1
}

//export InitializeSignerWithPort
func InitializeSignerWithPort(api unsafe.Pointer, port int64) int {
	dart_api_dl.Init(api)
	ksigner.InitializeSignerWithPort(port)
	return 1
}

//export NewSignerSession
func NewSignerSession(s *C.char, n *C.char, shareJson *C.char, messageToSign *C.char) {
	sessionid := deepCopy(C.GoString(s))
	messageToSignStr := deepCopy(C.GoString(messageToSign))
	shareJsonStr := deepCopy(C.GoString(shareJson))
	nodeid := deepCopy(C.GoString(n))
	println("NewSignerSession........................" + sessionid)
	println(messageToSign)
	ksigner.NewSignerSession(sessionid, nodeid, shareJsonStr, messageToSignStr)
	// C.free(unsafe.Pointer(s))
	// C.free(unsafe.Pointer(messageToSign))
	// C.free(unsafe.Pointer(shareJson))
	// C.free(unsafe.Pointer(n))
}

//export GetSignerMessage
func GetSignerMessage(s *C.char, m *C.char) *C.char {
	println("GetSignerMessage........................" + getActiveSessions().callbackType)

	sessionid := deepCopy(C.GoString(s))

	msgid := deepCopy(C.GoString(m))

	jsoncstr := C.CString(ksigner.GetMessage(sessionid, msgid))
	// TODO
	// // C.free(unsafe.Pointer(jsoncstr))
	println("end GetSignerMessage........................")
	//// C.free(unsafe.Pointer(s))
	//// C.free(unsafe.Pointer(m))
	return jsoncstr
}

//export HandleSignerMessage
func HandleSignerMessage(s *C.char, m *C.char) *C.char {
	println("HandleSignerMessage........................" + getActiveSessions().callbackType)

	sessionid := deepCopy(C.GoString(s))
	// // C.free(unsafe.Pointer(s))
	println(sessionid)

	body := deepCopy(C.GoString(m))
	// // C.free(unsafe.Pointer(m))
	println(body)

	err := ksigner.HandleMessage(sessionid, body)

	errstr := C.CString(err)
	// // C.free(unsafe.Pointer(errstr))
	return errstr
}

//export Sign
func Sign(s *C.char) {
	sessionid := deepCopy(C.GoString(s))
	ksigner.Sign(sessionid)
}

//export GetSignerResult
func GetSignerResult(s *C.char) *C.char {
	sessionid := deepCopy(C.GoString(s))
	// C.free(unsafe.Pointer(s))
	return C.CString(string(ksigner.GetResult(sessionid)))
}

func main() {

}
