package keygen

/*
   #include <stdlib.h>
*/

import (
	"encoding/base32"
	"encoding/json"
	"os"
	"reflect"
	"sync"
	"syscall/js"

	"github.com/btcsuite/btcd/btcec"
	"github.com/getamis/alice/crypto/tss/dkg"
	"github.com/getamis/alice/internal/message/types"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var instance *ActiveSessions
var once sync.Once

type DkgSession struct {
	id             string
	nodeid         string
	peers          map[string]*PeerInfo
	Dkg            *dkg.DKG
	done           chan struct{}
	threshold      int
	rank           int
	messageBox     map[string]*dkg.Message
	result         *DKGResult
	jscallback     js.Value
	ResultCallBack js.Value
}

type PeerInfo struct {
	Id   string `json:"id"`
	Rank uint32 `json:"rank"`
}

type ActiveSessions struct {
	Sessions map[string]*DkgSession
	PipeFile *os.File
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

func GetActiveSessions() *ActiveSessions {
	//singleton
	once.Do(func() {
		instance = &ActiveSessions{
			Sessions: make(map[string]*DkgSession),
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
	if err != nil {
		println("Cannot marshal message : " + err.Error())
	}
	// probably unncessary
	x := &dkg.Message{}
	proto.Unmarshal(bs, x)
	jsonstr, _ := protojson.Marshal(x)
	p.jscallback.Invoke(js.Null(), string(jsonstr))
}

func (p *DkgSession) getMessageId(peerid string, msg *dkg.Message) string {
	data := []byte(msg.Id + "_" + peerid + "_" + msg.Type.String())
	return base32.StdEncoding.EncodeToString(data)
}

func (p *DkgSession) OnStateChanged(oldState types.MainState, newState types.MainState) {
	if newState == types.StateFailed {
		println("Dkg failed", "old", oldState.String(), "new", newState.String())
		close(p.done)
		p.ResultCallBack.Invoke("Dkg failed", js.Null())
		return
	} else if newState == types.StateDone {
		println("Dkg done", "old", oldState.String(), "new", newState.String())
		result, err := p.Dkg.GetResult()
		if err == nil {
			p.fetchDKGResult(result)
		} else {
			println("Failed to get result from DKG", "err", err)
			p.ResultCallBack.Invoke(err.Error(), js.Null())
		}
		close(p.done)
		p.Dkg.Stop()
		jsonbytes, err := json.Marshal(p.result)
		if err != nil {
			p.ResultCallBack.Invoke(err.Error(), js.Null())
			return
		}
		p.ResultCallBack.Invoke(js.Null(), string(jsonbytes))
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

func NewDkgSession(sessionid string, nodeid string, threshold int, rank int, peers map[string]*PeerInfo, callback js.Value) *DkgSession {
	session := &DkgSession{
		id:         sessionid,
		done:       make(chan struct{}),
		threshold:  threshold,
		rank:       rank,
		peers:      peers,
		messageBox: make(map[string]*dkg.Message),
		nodeid:     nodeid,
		jscallback: callback,
	}
	dkg, err := dkg.NewDKG(btcec.S256(), session, uint32(session.threshold), uint32(session.rank), session)
	if err != nil {
		println(err.Error())
	}
	session.Dkg = dkg
	GetActiveSessions().Sessions[sessionid] = session
	return session
}
