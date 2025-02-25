package main

import "go-zinx/znet"

func main() {
	s := znet.NewServer("[ZINXV0.1]")
	s.Serve()
}
