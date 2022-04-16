package kommons

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
