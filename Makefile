all: help

.PHONY:
help:
	@echo "-------------"
	@echo "ToMaster Help"
	@echo "-------------"
	@echo "build       Build an executable"
	@echo "go-install  Build and Install ToMaster using `go install`"
	@echo "install     Build and Install ToMaster as symbolic link on /usr/local/bin "
	@echo "unistall    Build and Install ToMaster as symbolic link on /usr/local/bin "

build:
	go build -o bin/tomaster -ldflags="-s -w"  cmd/*.go

go-install:
	mkdir tomaster;
	cp cmd/*.go tomaster/;
	cd tomaster && go install;
	rm -f tomaster/*.go
	rmdir tomaster

install:
	mkdir ~/.tomaster
	go build -o  ~/.tomaster/tomaster -ldflags="-s -w" cmd/*.go
	ln -s ~/.tomaster/tomaster /usr/local/bin/tomaster

uninstall:
	rm /usr/local/bin/tomaster
	rm ~/.tomaster/tomaster
	rmdir ~/.tomaster
