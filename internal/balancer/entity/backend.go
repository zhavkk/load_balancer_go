package entity

type Backend struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	IsAlive bool
}

type BackendList struct {
	Backends []Backend `json:"backends"`
}
