## Introduction
*As of July 25th, 2018, this project is considered complete. Core features have been implemented. Any outstanding issues are considered enhancements or low-priority.*

This personal project implements a distributed file system. This project draws inspiration from an [architecture and interface described by Ivan Beschastnikh](http://www.cs.ubc.ca/~bestchai/teaching/cs416_2017w2/assign2/index.html). 

The file system exposes 2 interfaces to users: (1) the dfs API and (2) dfs file API. Detailed descriptions of each API may be found in the section, "dfslib API". Users are able to mount and dismount an instance of the distributed file system. Upon mounting, users are able to open ".dfs" files in 3 modes: (1) read (2) write and (3) disconnected read. 

Files consist of "chunks", which are length-32 byte arrays. A file contains 256 chunks. Users may read and write with per-chunk granularity. For each file, there may be one writer and many concurrent readers. 

In read and write mode, dfs guarantees strong consistency. All writes to a file will only occur successfully if the system can guarantee that future reads to this chunk return the updated value. However, if users wish to avoid the latency incurred by this guarantee, they may optionally open the file in disconnected read mode. Disconnected read mode offers users the ability to improve read latency at the expense of potentially stale data. 

A single server serializes all file operations from clients participating in the dfs application. File data is cached on each client. No file data is stored on the server. The server maintains a minimal set of metadata regarding each client to facilitate dfs services. 

## Assumptions
- The server does not fail
- Client nodes may fail-stop, but do not experience byzantine failures or partial failures
- Opening a file in read or write mode will create always new file, potentially overwriting existing files
- To open a file in disconnected read mode, the file must have previously been created

## How To Run
*please ensure you have installed Go 1.9.2 or later on your system; these instructions additionally assume you have added the go command to your system path*
1. Open 2 or 3 command terminals.
2. In one command terminal, navigate to the directory containing the server file.
3. Input the following command in the terminal to run the server: go run server.go 127.0.0.1:3000
4. In a separate command terminal, navigate to the directory containing the application files.
5. Input the following command to run a sample application: go run app.go

## Sample Applications

1. app and app2 are intended to be run in tandem. Please run app first and app2 immediately afterwards. These two applications demonstrate dfs helper method functionality (e.g. LocalFileExists) and demonstrates basic read and write file operations with two concurrent users. 
2. app3 is intended to be run standalone. This application demonstrates the disconnected read operation.
3. app4 and app5 are intended to be run in tandem. Please run app4 first and app5 immediately afterwards. These two applications exercise all file operation functionality: write, read, and disconnected read with two concurrent clients.

## Directory Structure
- app
  - app.go: Contains sample applications that demonstrate different functionality of the dfs application
  - ...
  - app5.go
- dfslib
  - dfslib.go: Implements the dfs file system API
- server
  - server.go: Implements the single, centralized server to which clients connect to
- tmp: Contains dfs files for a client
- tmp2: Contains dfs files for a second client
- test: Contains miscellaneous test files

## System Topology
The dfs application consists of 2 nodes. 
- client node
  - application layer: The user may import the dfs library by mounting an instance of dfs. 
  - dfs library layer: The dfslib API exposes several primitives to the user to create, open, and modify ".dfs" files.
- server node

The client nodes are connected to the server node in a star topological layout. Each client has no knowledge of other clients utilizing the services of the distributed file system. The server transparently mediates communication between each client using bi-directional RPC.

## dfslib API

- MountDFS(serverAddr string, localIP string, localPath string) : (dfs DFS, err error)

- DFS
  - Open(fname string, mode FileMode) : (f DFSFile, err error) - Mounts an instance of 
  - LocalFileExists(fname string)     : (exists bool, err error)
  - GlobalFileExists(fname string)    : (exists bool, err error)
  - UMountDFS()                       : (err error)
  
- DFSFile
  - Read(chunkNum uint8, chunk \*Chunk)  : (err error)
  - Write(chunkNum uint8, chunk \*Chunk) : (err error)
  - Dread(chunkNum uint8, chunk \*Chunk) : (err error)
  - Close()                              : (err error)
