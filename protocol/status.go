package protocol

type PingStatus struct {
	Version            StatusVersion     `json:"version"`
	Players            StatusPlayerInfo  `json:"players"`
	Description        StatusDescription `json:"description"`
	ServerIcon         string            `json:"favicon"`
	EnforcesSecureChat bool              `json:"enforcesSecureChat"`
}

type StatusVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type StatusPlayerInfo struct {
	Max    int                  `json:"max"`
	Online int                  `json:"online"`
	Sample []StatusPlayerSample `json:"sample"`
}

type StatusPlayerSample struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type StatusDescription struct {
	Text string `json:"text"`
}