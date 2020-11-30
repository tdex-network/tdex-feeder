package ports

type TickerWithPrice struct {
	Ticker string
	Price float64
}

type subscription struct {
	Name     string `json:"name"`
	Interval int    `json:"interval,omitempty"`
	Token    string `json:"token,omitempty"`
	Depth    int    `json:"depth,omitempty"`
	Snapshop bool   `json:"snapshot,omitempty"`
}

// RequestMessage is the data structure used to create
// jsons in order subscribe to market updates on Kraken
type requestMessage struct {
	Event        string        `json:"event"`
	Pair         []string      `json:"pair,omitempty"`
	Subscription *subscription `json:"subscription,omitempty"`
	Reqid        int           `json:"reqid,omitempty"`
}