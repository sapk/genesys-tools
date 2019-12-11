go run main.go dump localhost:8080

go build -o genesys-tools .

gox -ldflags "-s -w" -osarch="windows/amd64" -output="build/{{.Dir}}-tools-{{.OS}}-{{.Arch}}"

dot -Tpng host-port/graph.dot > test.png  
dot -Tsvg host-port/graph.dot > test.svg
