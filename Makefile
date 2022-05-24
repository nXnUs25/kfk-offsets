.PHONY : core cli offset lag cg 

kfk: core cli offset lag cg 
	go build -o kfkgo *.go

core:
	go build core/*.go
lag: 
	go build lag/*.go
cg:
	go build cg/*.go
offset: 
	go build offset/*.go
cli: 
	go build cmd/*.go

clean:
	rm ./kfkgo