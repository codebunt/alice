package peer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/getamis/sirius/log"
	"github.com/gin-gonic/gin"
)

type Node struct {
	id          string
	port        uint32
	dkgsessions map[string]*dkgsession
}

type dkgrequest struct {
	Sessionid string   `json:"sid" binding:"required"`
	Peers     []string `json:"peers" binding:"required"`
	Threshold uint32   `json:"t" binding:"required"`
	Rank      uint32   `json:"rank"`
}

func (d *dkgrequest) getPeers() (map[string]*PeerInfo, error) {
	peermap := make(map[string]*PeerInfo)
	for i := 0; i < len(d.Peers); i++ {
		log.Warn(d.Peers[i])
		res := strings.Split(d.Peers[i], ":")
		port, err := strconv.Atoi(res[1])
		if err != nil {
			return nil, err
		}
		rank, err := strconv.Atoi(res[2])
		if err != nil {
			return nil, err
		}

		p := &PeerInfo{
			id:   res[3],
			host: res[0],
			port: uint32(port),
			rank: uint32(rank),
		}
		peermap[p.id] = p
	}
	return peermap, nil
}

func NewNode(id string, port uint32) *Node {
	node := &Node{
		id:          id,
		port:        port,
		dkgsessions: make(map[string]*dkgsession),
	}
	return node
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) StartDKGSession(peers map[string]*PeerInfo, threshold uint32, rank uint32, sessionid string) *dkgsession {
	session := newDKGSession(n, peers, threshold, rank)
	log.Warn("StartDKGSession")
	if len(strings.Trim(sessionid, " ")) > 0 {
		session.id = sessionid
	}
	n.dkgsessions[session.id] = session
	for _, peer := range peers {
		n.initiatePeerSession(peer, session)
	}
	session.GetKey()
	return session
}

func (n *Node) Start() {
	r := gin.Default()
	r.POST("/startdkg", func(c *gin.Context) {
		var req dkgrequest
		c.BindJSON(&req)
		log.Warn("..............startdkg..................")
		log.Warn("fromid - " + c.GetHeader("from-node"))
		log.Warn("toid - " + c.GetHeader("to-node"))
		if n.dkgsessions[req.Sessionid] != nil {
			c.Status(http.StatusOK)
			return
		} else {
			log.Warn("Missing seesion id creating new")
		}
		peers, err := req.getPeers()
		if err != nil {
			c.Data(http.StatusBadRequest, "application/text", []byte("error starting dkg"))
		}
		n.StartDKGSession(peers, req.Threshold, req.Rank, req.Sessionid)
		c.String(http.StatusOK, "%t", "ww")
	})

	r.POST("/dkg", func(c *gin.Context) {
		sessionid := c.GetHeader("session-id")
		log.Warn("..............dkg.................." + sessionid)
		log.Warn("fromid - " + c.GetHeader("from-node"))
		log.Warn("toid - " + c.GetHeader("to-node"))

		dkgsession := n.dkgsessions[sessionid]
		body, _ := ioutil.ReadAll(c.Request.Body)

		err := dkgsession.handle(body)
		if err != nil {
			log.Error(err.Error())
			c.Data(http.StatusBadRequest, "application/text", []byte("error handling message"))
			return
		}
		c.Data(http.StatusOK, "application/octet-stream", []byte(sessionid))
	})
	log.Warn("Server running")

	r.Run(":" + strconv.FormatUint(uint64(n.port), 10)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

func (n *Node) initiatePeerSession(target *PeerInfo, session *dkgsession) error {
	if n.id == target.id {
		return nil
	}
	log.Warn("http://" + target.host + ":" + strconv.Itoa(int(target.port)) + "/startdkg")
	log.Warn(n.ID())
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	dkgrequest := &dkgrequest{
		Sessionid: session.id,
		Peers:     session.GetPeerURLsForPeer(target),
		Threshold: session.GetThreshold(),
		Rank:      session.GetRank(),
	}
	b, err := json.Marshal(dkgrequest)
	log.Info(string(b))

	req, err := http.NewRequest("POST", "http://"+target.host+":"+strconv.Itoa(int(target.port))+"/startdkg", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("from-node", n.id)
	req.Header.Set("to-node", target.id)

	rsp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		log.Warn("http://" + target.host + ":" + strconv.Itoa(int(target.port)) + "/startdkg")
		log.Warn("Request failed with response code: " + strconv.Itoa(rsp.StatusCode))
	}
	return nil
}
