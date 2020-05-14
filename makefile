# Copyright © 2020 Thilina Manamgoda
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

TOOL_VERSION=v0.1.0
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=etchosts
DEP=dep
GOLINT=golint
GOFMT=$(GOCMD) fmt

TEST_PKGS=./pkg/...
FMT_PKGS=./cmd/... ./pkg/...
LDFLAGS=-X 'github.com/ThilinaManamgoda/etchosts/cmd.Version=$(TOOL_VERSION)'

all: clean deps lint ineffassign unit-test build-linux build-darwin

build-doc:
		$(GOBUILD) -tags doc  -o gen_doc -v

build:
		$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) -v

unit-test:
		$(GOTEST) -v $(TEST_PKGS)


lint:
		$(GOGET) -u golang.org/x/lint/golint
		$(GOLINT) $(FMT_PKGS)

clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)

fmt:
		$(GOFMT) $(FMT_PKGS)

run:
		$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)

deps:
		$(DEP) ensure

build-linux:
		env GOOS="linux" GOARCH="amd64" $(GOBUILD) -ldflags "$(LDFLAGS)" -o "target/linux/$(TOOL_VERSION)/$(BINARY_NAME)" -v

build-darwin:
		env GOOS="darwin" GOARCH="amd64" $(GOBUILD) -ldflags "$(LDFLAGS)" -o "target/darwin/$(TOOL_VERSION)/$(BINARY_NAME)" -v

build-darwin-tar:
		cd "target/darwin/$(TOOL_VERSION)"; tar -cvf etchosts-$(TOOL_VERSION).tar.gz $(BINARY_NAME);

ineffassign:
		$(GOGET) -u github.com/gordonklaus/ineffassign
		ineffassign main.go cmd/* pkg/*
