JATTACH_VERSION=2.2

ifneq ($(findstring Windows,$(OS)),)
  CL=cl.exe
  CFLAGS=/O2 /D_CRT_SECURE_NO_WARNINGS
  JATTACH_EXE=jattach.exe
  JATTACH_DLL=jattach.dll
else 
  JATTACH_EXE=jattach

  UNAME_S:=$(shell uname -s)
  ifeq ($(UNAME_S),Darwin)
    CFLAGS ?= -O3 -arch x86_64 -arch arm64 -mmacos-version-min=10.12
    JATTACH_DLL=libjattach.dylib
  else
    CFLAGS ?= -O3
    JATTACH_DLL=libjattach.so
  endif

  ifeq ($(UNAME_S),Linux)
    ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
    RPM_ROOT=$(ROOT_DIR)/build/rpm
    SOURCES=$(RPM_ROOT)/SOURCES
    SPEC_FILE=jattach.spec
  endif
endif


.PHONY: all dll clean rpm-dirs rpm go-build go-test go-cli go-cli-posix go-cli-windows go-cli-all go-clean

all: build build/$(JATTACH_EXE)

dll: build build/$(JATTACH_DLL)

build:
	mkdir -p build

build/jattach: src/posix/*.c src/posix/*.h
	$(CC) $(CPPFLAGS) $(CFLAGS) $(LDFLAGS) -DJATTACH_VERSION=\"$(JATTACH_VERSION)\" -o $@ src/posix/*.c

build/$(JATTACH_DLL): src/posix/*.c src/posix/*.h
	$(CC) $(CPPFLAGS) $(CFLAGS) $(LDFLAGS) -fPIC -shared -fvisibility=hidden -o $@ src/posix/*.c

build/jattach.exe: src/windows/jattach.c
	$(CL) $(CFLAGS) /DJATTACH_VERSION=\"$(JATTACH_VERSION)\" /Fobuild/jattach.obj /Fe$@ $^ advapi32.lib /link /SUBSYSTEM:CONSOLE,5.02

clean:
	rm -rf build

# Go targets
go-build:
	go build

go-test:
	go test -v

go-cli:
	cd cmd/jattach-go && go build -o jattach-go

# Cross-compile Go CLI for POSIX platforms (Linux and macOS)
# Note: CGo is used, so cross-compilation works within the same OS family
go-cli-posix: build
	@echo "Building Go CLI for POSIX platforms..."
	cd cmd/jattach-go && \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o ../../build/jattach-go-linux-amd64 && \
		GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -o ../../build/jattach-go-linux-arm64 && \
		GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o ../../build/jattach-go-darwin-amd64 && \
		GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o ../../build/jattach-go-darwin-arm64
	@echo "Go CLI POSIX binaries built in build/"

# Cross-compile Go CLI for Windows
# Note: Requires mingw-w64 cross-compiler: brew install mingw-w64 (macOS) or apt install mingw-w64 (Linux)
go-cli-windows: build
	@echo "Building Go CLI for Windows..."
	@echo "Note: This requires mingw-w64 cross-compiler toolchain"
	cd cmd/jattach-go && \
		CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o ../../build/jattach-go-windows-amd64.exe && \
		CGO_ENABLED=1 GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc go build -o ../../build/jattach-go-windows-386.exe
	@echo "Go CLI Windows binaries built in build/"

# Cross-compile Go CLI for all platforms (requires cross-compilers)
go-cli-all: go-cli-posix go-cli-windows

go-clean:
	go clean
	cd cmd/jattach-go && go clean
	rm -f cmd/jattach-go/jattach-go
	rm -f build/jattach-go-*

$(RPM_ROOT):
	mkdir -p $(RPM_ROOT)

rpm-dirs: $(RPM_ROOT)
	mkdir -p $(RPM_ROOT)/SPECS
	mkdir -p $(SOURCES)/bin
	mkdir -p $(RPM_ROOT)/BUILD
	mkdir -p $(RPM_ROOT)/SRPMS
	mkdir -p $(RPM_ROOT)/RPMS
	mkdir -p $(RPM_ROOT)/tmp

rpm: rpm-dirs build build/$(JATTACH_EXE)
	cp $(SPEC_FILE) $(RPM_ROOT)/
	cp build/jattach $(SOURCES)/bin/
	rpmbuild -bb \
                --define '_topdir $(RPM_ROOT)' \
                --define '_tmppath $(RPM_ROOT)/tmp' \
                --clean \
                --rmsource \
                --rmspec \
                --buildroot $(RPM_ROOT)/tmp/build-root \
                $(RPM_ROOT)/jattach.spec
