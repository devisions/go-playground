
go build -ldflags "-pluginpath=plugin/hot-$(uuidgen)" \
         -buildmode=plugin \
         -o plugin1.so main.go
