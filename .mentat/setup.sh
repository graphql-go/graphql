curl -L https://go.dev/dl/go1.13.15.linux-amd64.tar.gz -o go.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf go.tar.gz
export PATH=$PATH:/usr/local/go/bin
go mod download