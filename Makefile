c:
		rm -rf ./down ./bin

b: c
		go build  -o "bin/down" ./main.go
		ln -s bin/down down

r: b
		exec ./bin/down


i: b
		go install


#vim:ft=bash
