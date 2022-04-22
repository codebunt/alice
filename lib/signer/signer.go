package signer

/*
   #include <stdlib.h>
*/
import "C"

import (
	"bufio"
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"reflect"
	"sync"

	"github.com/btcsuite/btcd/btcec"
	"github.com/codebunt/dart_api_dl"
	"github.com/getamis/alice/crypto/birkhoffinterpolation"
	"github.com/getamis/alice/crypto/ecpointgrouplaw"
	"github.com/getamis/alice/crypto/homo/paillier"
	"github.com/getamis/alice/crypto/tss/dkg"
	"github.com/getamis/alice/crypto/tss/signer"
	kommons "github.com/getamis/alice/lib/commons"

	"github.com/getamis/alice/internal/message/types"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var instance *ActiveSessions
var once sync.Once

type SignerSession struct {
	id         string
	nodeid     string
	peers      map[string]*PeerInfo
	signer     *signer.Signer
	done       chan struct{}
	messageBox map[string]*signer.Message
	result     *SignerResult
}

type PeerInfo struct {
	Id   string `json:"id"`
	Rank uint32 `json:"rank"`
}

type ActiveSessions struct {
	sessions     map[string]*SignerSession
	pipeFile     *os.File
	callbackType string
	port         int64
	tcpwriter	 *bufio.Writer
}

type SignerResult struct {
	R string `json:"r"`
	S string `json:"s"`
}

func GetActiveSessions() *ActiveSessions {
	//singleton
	once.Do(func() {
		instance = &ActiveSessions{
			sessions: make(map[string]*SignerSession),
		}
	})
	return instance
}

func (p *SignerSession) NumPeers() uint32 {
	return uint32(len(p.peers))
}

func (p *SignerSession) PeerIDs() []string {
	keys := reflect.ValueOf(p.peers).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	return strkeys
}

func (p *SignerSession) SelfID() string {
	return p.nodeid
}

func (p *SignerSession) MustSend(peerid string, message interface{}) {
	msg, _ := message.(proto.Message)
	bs, err := proto.Marshal(msg)
	// probably unncessary
	x := &signer.Message{}
	proto.Unmarshal(bs, x)
	msgId := p.getMessageId(peerid, x)
	println(msgId)
	p.messageBox[msgId] = x
	if err != nil {
		println("Cannot marshal message : " + err.Error())
	} else {
		if GetActiveSessions().callbackType == "PIPE" {
			GetActiveSessions().pipeFile.WriteString("signerround:" + p.id + ":" + peerid + ":" + msgId + "\n")
		}
		if GetActiveSessions().callbackType == "TCP" {
			GetActiveSessions().tcpwriter.WriteString("signerround:" + p.id + ":" + peerid + ":" + msgId + "\n")
			GetActiveSessions().tcpwriter.Flush()
		}
		if GetActiveSessions().callbackType == "DART_PORT" {
			dart_api_dl.SendToPort(GetActiveSessions().port, "signerround", p.id, peerid, msgId)
		}

	}
}

func (p *SignerSession) getMessageId(peerid string, msg *signer.Message) string {
	return msg.Type.String() + "_" + peerid
}

func (p *SignerSession) OnStateChanged(oldState types.MainState, newState types.MainState) {
	if newState == types.StateFailed {
		println("Dkg failed", "old", oldState.String(), "new", newState.String())
		close(p.done)
		return
	} else if newState == types.StateDone {
		println("Signer done", "old", oldState.String(), "new", newState.String())
		result, err := p.signer.GetResult()
		var status = "signer_success"
		if err == nil {
			p.fetchSignerResult(result)
		} else {
			println("Failed to get result from DKG", "err", err)
			status = "signer_failure"
		}
		close(p.done)
		p.signer.Stop()
		if GetActiveSessions().callbackType == "PIPE" {
			GetActiveSessions().pipeFile.WriteString(status + ":" + p.id + "\n")
		}
		if GetActiveSessions().callbackType == "TCP" {
			GetActiveSessions().tcpwriter.WriteString(status + ":" + p.id + "\n")
			GetActiveSessions().tcpwriter.Flush()
		}
		if GetActiveSessions().callbackType == "DART_PORT" {
			dart_api_dl.SendToPort(GetActiveSessions().port, status, p.id, "", "")
		}

		return
	}
	println("State changed", "old", oldState.String(), "new", newState.String())
}

func (p *SignerSession) fetchSignerResult(result *signer.Result) error {
	signerResult := &SignerResult{
		R: result.R.String(),
		S: result.S.String(),
	}
	p.result = signerResult
	return nil
}

func InitializeWithPipe(f *os.File) {
	GetActiveSessions().pipeFile = f
	GetActiveSessions().callbackType = "PIPE"
}

func InitializeWithTCPConnection(f *bufio.Writer) {
	GetActiveSessions().tcpwriter = f
	GetActiveSessions().callbackType = "TCP"
}

func InitializeSignerWithPort(port int64) {
	GetActiveSessions().callbackType = "DART_PORT"
	GetActiveSessions().port = port
}

func NewSignerSession(sessionid string, nodeid string, sharejson string, messageToSign string) {
	session := &SignerSession{
		id:         sessionid,
		done:       make(chan struct{}),
		peers:      make(map[string]*PeerInfo),
		messageBox: make(map[string]*signer.Message),
		nodeid:     nodeid,
	}
	GetActiveSessions().sessions[sessionid] = session
	var dkgResult kommons.DKGResult
	json.Unmarshal([]byte(sharejson), &dkgResult)
	for peerID, bk := range dkgResult.BKs {
		if peerID != nodeid {
			session.peers[peerID] = &PeerInfo{
				Id:   peerID,
				Rank: bk.Rank,
			}
		}
	}
	// For simplicity, we use Paillier algorithm in signer.
	paillier, err := paillier.NewPaillier(2048)
	if err != nil {
		println("Cannot create a paillier function", "err", err)
		return
	}

	result, _ := convertDKGResult(dkgResult)
	ksigner, err := signer.NewSigner(session, result.PublicKey, paillier, result.Share, result.Bks, []byte(messageToSign), session)
	if err != nil {
		println(err.Error())
	}

	session.signer = ksigner

}

func convertDKGResult(dkgResult kommons.DKGResult) (*dkg.Result, error) {

	// Build public key.
	x, ok := new(big.Int).SetString(dkgResult.Pubkey.X, 10)
	if !ok {
		println("Cannot convert string to big int", "x", dkgResult.Pubkey.X)
		return nil, ErrConversion
	}
	y, ok := new(big.Int).SetString(dkgResult.Pubkey.Y, 10)
	if !ok {
		println("Cannot convert string to big int", "y", dkgResult.Pubkey.Y)
		return nil, ErrConversion
	}
	pubkey, err := ecpointgrouplaw.NewECPoint(btcec.S256(), x, y)
	if err != nil {
		println("Cannot get public key", "err", err)
		return nil, err
	}

	// Build share.
	share, ok := new(big.Int).SetString(dkgResult.Share, 10)
	if !ok {
		println("Cannot convert string to big int", "share", share)
		return nil, ErrConversion
	}

	result := &dkg.Result{
		PublicKey: pubkey,
		Share:     share,
		Bks:       make(map[string]*birkhoffinterpolation.BkParameter),
	}

	// Build bks.
	for peerID, bk := range dkgResult.BKs {
		x, ok := new(big.Int).SetString(bk.X, 10)
		if !ok {
			println("Cannot convert string to big int", "x", bk.X)
			return nil, ErrConversion
		}
		result.Bks[peerID] = birkhoffinterpolation.NewBkParameter(x, bk.Rank)
	}

	return result, nil
}

var (
	// ErrConversion for big int conversion error
	ErrConversion = errors.New("conversion error")
)

func GetMessage(sessionid string, msgid string) string {
	println("GetMessage........................" + msgid)
	bs := GetActiveSessions().sessions[sessionid].messageBox[msgid]
	println(bs)
	mOptions := protojson.MarshalOptions{
		Multiline:       true,
		EmitUnpopulated: true,
	}
	jsonstr, _ := mOptions.Marshal(bs)
	return string(jsonstr)
}

func HandleMessage(sessionid string, body string) string {
	println("Signer HandleMessage........................")
	// handle data
	x := &signer.Message{}
	// unmarshal it
	err := protojson.Unmarshal([]byte(body), x)
	if err != nil {
		println("Cannot unmarshal data", "err", err.Error())
		return err.Error()
	}
	err = GetActiveSessions().sessions[sessionid].signer.AddMessage(x)
	if err != nil {
		println("Cannot add message to DKG", "err", err)
		return err.Error()
	}
	println("Added message to DKG")
	return ""
}

func Sign(sessionid string) {
	session := GetActiveSessions().sessions[sessionid]
	// 1. Start a DKG process.
	session.signer.Start()
}

func GetResult(sessionid string) string {
	session := GetActiveSessions().sessions[sessionid]
	if session.result == nil {
		return string("{}")
	}
	jsonbytes, err := json.Marshal(session.result)
	if err != nil {
		return string("{}")
	}
	return string(jsonbytes)
}
