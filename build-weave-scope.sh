export GOPATH=~/go
mkdir $GOPATH
cd $GOPATH
echo "executing go clean"
go clean -i net
echo "executing go install"
go install -tags netgo std
echo "mkdir"
rm -rf src/github.com/weaveworks/scope
mkdir -p src/github.com/weaveworks
cd src/github.com/weaveworks
echo "cloning"
git clone http://github.com/dilgerma/scope
cd scope
make deps SUDO=
make SUDO=
