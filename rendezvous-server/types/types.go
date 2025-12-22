package types

type IPAddressPair struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

type IPAddressPort struct {
	AddressPort string `json:"addressPort"`
}

type Identifier struct {
	Identifier string `json:"identifier"`
}

type PortResponse struct {
	Port int `json:"port"`
}

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}
