gen:
	rm -rf rpc
	mkdir rpc
	protoc ./proto/* --go_out=plugins=grpc:. --twirp_out=.

serve:
	clear
	go run main.go

format:
	go fmt ./server/*
	go fmt ./config/*

dbinit: dbinit.c
	

dbinit.c: dbinit.b
	

dbinit.b: dbinit.a
	

dbinit.a:
	