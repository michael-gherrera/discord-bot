package datasource

type WizdaddyResponse struct {
	Symbol         string `json:"symbol"`
	ExpirationDate string `json:"expirationDate"`
	StrikePrice    string `json:"strikePrice"`
	Type           string `json:"type"`
	RobinhoodID    string `json:"robinhoodId"`
	RiskText       string `json:"riskText"`
	WizdaddyID     string `json:"wizdaddyId"`
	Timestamp      string `json:"timestamp"`
}
