package peer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/getamis/alice/crypto/tss/dkg"
	"github.com/getamis/alice/internal/message/types"
	"github.com/getamis/sirius/log"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/protocol"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

type dkgsession struct {
	id        string
	peers     map[string]*PeerInfo
	dkg       *dkg.DKG
	done      chan struct{}
	protocol  protocol.ID
	node      *Node
	threshold uint32
	rank      uint32
}

type PeerInfo struct {
	id   string
	host string
	port uint32
	rank uint32
}

func NewPeerInfo(id string, host string, port uint32, rank uint32) *PeerInfo {
	return &PeerInfo{
		id:   id,
		host: host,
		port: port,
		rank: rank,
	}

}

type DKGResult struct {
	Share  string        `yaml:"share"`
	Pubkey Pubkey        `yaml:"pubkey"`
	BKs    map[string]BK `yaml:"bks"`
}

type Pubkey struct {
	X string `yaml:"x"`
	Y string `yaml:"y"`
}

type BK struct {
	X    string `yaml:"x"`
	Rank uint32 `yaml:"rank"`
}

func newDKGSession(node *Node, peers map[string]*PeerInfo, threshold uint32, rank uint32) *dkgsession {
	dkgPeerMgr := &dkgsession{
		id:        uuid.New().String(),
		peers:     peers,
		done:      make(chan struct{}),
		protocol:  "/dkg/1.0.0",
		node:      node,
		threshold: threshold,
		rank:      rank,
	}
	dkg, err := dkg.NewDKG(btcec.S256(), dkgPeerMgr, threshold, rank, dkgPeerMgr)
	if err != nil {
		log.Warn(err.Error())
	}
	dkgPeerMgr.dkg = dkg
	return dkgPeerMgr
}

func (p *dkgsession) GetThreshold() uint32 {
	return p.threshold
}

func (p *dkgsession) GetRank() uint32 {
	return p.rank
}

func (p *dkgsession) GetKey() {
	log.Info("Start DKG")
	// 1. Start a DKG process.
	p.dkg.Start()
	defer p.dkg.Stop()

	// 2. Wait the dkg is done or failed
	<-p.done
}

func (p *dkgsession) PeerIDs() []string {
	keys := reflect.ValueOf(p.peers).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		if keys[i].String() != p.node.id {
			strkeys[i] = keys[i].String()
			log.Warn(strkeys[i])
		}
	}
	return strkeys
}

func (p *dkgsession) GetPeerURLsForPeer(target *PeerInfo) []string {
	strkeys := make([]string, len(p.peers))
	var index = 0
	for _, peer := range p.peers {
		// Ignore the target
		if target.id != peer.id {
			strkeys[index] = peer.host + ":" + strconv.FormatUint(uint64(peer.port), 10) + ":" + strconv.FormatUint(uint64(peer.rank), 10) + ":" + peer.id
			index++
		}
	}
	// Add self to peer
	strkeys[index] = "localhost:" + strconv.FormatUint(uint64(p.node.port), 10) + ":" + strconv.FormatUint(uint64(p.rank), 10) + ":" + p.node.id
	return strkeys
}

func (p *dkgsession) NumPeers() uint32 {
	return uint32(len(p.peers))
}

func (p *dkgsession) SelfID() string {
	return p.node.id
}

func (p *dkgsession) MustSend(id string, message interface{}) {
	log.Warn(id)
	if p.node.id == id {
		return
	}
	p.send(context.Background(), p.peers[id], message, p.protocol)
}

func (p *dkgsession) send(ctx context.Context, target *PeerInfo, data interface{}, protocol protocol.ID) error {
	msg, ok := data.(proto.Message)
	if !ok {
		log.Warn("invalid proto message")
		return errors.New("invalid proto message")
	}
	bs, err := proto.Marshal(msg)
	if err != nil {
		log.Warn("Cannot marshal message", "err", err)
		return err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	tmp := fmt.Sprintf("%#v", target)
	log.Warn(tmp)
	req, err := http.NewRequest("POST", "http://"+target.host+":"+strconv.Itoa(int(target.port))+"/dkg", bytes.NewReader(bs))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("session-id", p.id)
	req.Header.Set("from-node", p.node.id)
	req.Header.Set("to-node", target.id)

	rsp, _ := client.Do(req)
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		log.Warn("http://" + target.host + ":" + strconv.Itoa(int(target.port)) + "/dkg")
		log.Warn("Request failed with response code: " + strconv.Itoa(rsp.StatusCode))
	}

	return nil
}

func (p *dkgsession) handle(body []byte) error {
	// handle data
	x := &dkg.Message{}

	// unmarshal it
	err := proto.Unmarshal(body, x)
	if err != nil {
		log.Error("Cannot unmarshal data", "err", err)
		return err
	}

	err = p.dkg.AddMessage(x)

	if err != nil {
		log.Warn("Cannot add message to DKG", "err", err)
		return err
	}
	log.Warn("Added message to DKG", "err", err)

	return nil
}

func (p *dkgsession) OnStateChanged(oldState types.MainState, newState types.MainState) {
	if newState == types.StateFailed {
		log.Error("Dkg failed", "old", oldState.String(), "new", newState.String())
		close(p.done)
		return
	} else if newState == types.StateDone {
		log.Info("Dkg done", "old", oldState.String(), "new", newState.String())
		result, err := p.dkg.GetResult()
		if err == nil {
			p.fetchDKGResult(p.node.id, result)
		} else {
			log.Warn("Failed to get result from DKG", "err", err)
		}
		close(p.done)
		return
	}
	log.Info("State changed", "old", oldState.String(), "new", newState.String())
}

func (p *dkgsession) fetchDKGResult(id string, result *dkg.Result) {
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
	WriteYamlFile(dkgResult, getFilePath(p.node.id+"_"+p.id))
}

func WriteYamlFile(yamlData interface{}, filePath string) error {
	data, err := yaml.Marshal(yamlData)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, data, 0600)
}

func getFilePath(id string) string {
	return fmt.Sprintf("%s-output.yaml", id)
}
