.PHONY: all windows mac linux-amd64 linux-arm64 install-tools build-linux-image

all: windows mac linux-amd64 linux-arm64

install-tools:
	go install github.com/fyne-io/fyne-cross@latest
	go install fyne.io/fyne/v2/cmd/fyne@latest

build-linux-image:
	docker build -f Dockerfile.linux-audio -t oddcipher-linux-audio .

windows:
	fyne-cross windows -arch amd64

mac:
	fyne-cross darwin -arch arm64

linux-amd64: build-linux-image
	fyne-cross linux -arch amd64 -image oddcipher-linux-audio

linux-arm64: build-linux-image
	fyne-cross linux -arch arm64 -image oddcipher-linux-audio
