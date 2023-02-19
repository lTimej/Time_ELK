package common

type EtcdMsg struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}
