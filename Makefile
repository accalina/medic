# run pup script
run: main.go
	go run main.go

# do standard compilation
compile: main.go
	[ -f medic ] && rm medic
	go build -ldflags "-s -w" -o ./medic main.go
	upx -9 --lzma medic
	chmod +x medic

# do optimize compilation
compile-prod: main.go
	[ -f medic ] && rm medic
	go build -a -gcflags=all="-l -B" -ldflags "-s -w" -o ./medic main.go
	upx -9 --lzma medic
	chmod +x medic

