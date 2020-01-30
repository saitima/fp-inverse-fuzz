# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BIN=bin
LIB_NAME=montinversego
LIBARY_FULL_NAME=lib$(LIB_NAME).so
ARCHIVE_FULL_NAME=lib$(LIB_NAME).a
LIB_DIR=./lib
SRC=.
# SRC=api
CXX=gcc
CXX_FLAGS := -L$(LIB_DIR) -l$(LIB_NAME) -I$(LIB_DIR)

all: clean build test

static: clean build_static

external: clean build_external

static_external: clean build_static_external

build: 
		CGO_ENABLED=1 $(GOBUILD) -o $(LIB_DIR)/$(LIBARY_FULL_NAME) -buildmode=c-shared $(SRC)/inverse.go		
build_static: 
		CGO_ENABLED=1 $(GOBUILD) -o $(LIB_DIR)/$(ARCHIVE_FULL_NAME) -buildmode=c-archive $(SRC)/inverse.go	
build_external: 
		CGO_ENABLED=1 $(GOBUILD) -o $(LIB_DIR)/$(LIBARY_FULL_NAME) -buildmode=c-shared $(SRC)/inverse.go	
build_static_external: 
		CGO_ENABLED=1 $(GOBUILD) -o $(LIB_DIR)/$(ARCHIVE_FULL_NAME) -buildmode=c-archive $(SRC)/inverse.go	


clean: 
		$(GOCLEAN)
		rm -f $(LIB_DIR)/*
test:	
	$(CXX) $(CXX_FLAGS) -o $(LIB_DIR)/inverse $(SRC)/inverse.c
	cd lib && DEBUG=off  ./inverse