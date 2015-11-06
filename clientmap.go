package main

var ClientInterfaces = map[string]ClientIf{"ribd": &RibClient{}, "asicd": &AsicDClient{}}
