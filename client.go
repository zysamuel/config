package main

type ClientIf interface {
	Initialize(name string, address string)
	ConnectToServer() bool
	IsConnectedToServer() bool
}

type Client struct {
	Name string `json:Name`
	Port int    `json:Port`
	Intf ClientIf
}
