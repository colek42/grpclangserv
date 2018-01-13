# go-grpc-language-server
A simple gRPC based language server build on top of go guru
# How to run it

1. Make sure your gopath is set
2. make serve
	- This will compile the protobufs, lint, test, and run the app
3. Run the client code, client/main.go


# How it works

The program I wrote is a gRPC service built for scalability and high availability.  The client makes a RPC call.  The rpc handler then calls the Query function (after some logic dealing with byte offset and pathing).  The query function sets up a call to my fork of GoGuru which returns the response via gRPC.

The GoGuru fork of mine https://github.com/colek42/guru is very hacky.  I was hoping to be able to use guru as is, but I faced two obsticals.

1.  Guru is a main package and not exportable.
2.  Guru's use of interfaces means there really is no common well defined response that I could marshal into .proto

So I decided it would be fun to hack it up, and it was.  With my design descsion and time available I needed to take the third part of the problem off the table (more on this later).  I decided that the "definition" mode from guru gave me everything I needed for the first two section of the challenge.  I removed the logic that gave quick responses as that logic looked like it made the problem more complicated.  I resturcted the program so that the response from a "definition" query has a well defined type.

# Tradeoffs

I did not spend much time on client code.  The client package I have is pretty much worthless, except for testing.  However, because I used a language agnostic protocall a client can be quickly written for a variety of use cases.

I have no integration testing.  I like to do integration testing when I set up CI.

This service is very slow, this let me keep it simple.  There are a few ways I could speed it up without adding additional services.
- memoize the query function.
- create a virtual file system for all of the code
- restore the "shortcut" code in my guru fork


# Alternative Approaches
1.  Use JSON Response from guru
	- This would require a more complicated build.  When my build process is simple it is easier to flush out issues.  Having a signle binary is about as siimple as it gets.
	- I did not want to marshal->unmarshal->marshal.  This takes time, and memory, and bad bugs seem to always pop up in this code.  Being able to stay type safe thoughout my service call makes me happy.
2. Impliment the code crawling logic myself.
	- The guru code is well tested and written by people with much more domain knowledge than I, and it is free.
# Issues
- I don't have a real UX, instead I have an API.  It didn't make sense to write a command line utility on top of guru, as it already exists.  Plus, I wanted to show off my ability to build scalable services.
- I am not 100% sure if the LineToByteOffset function is correct.  It might be off by one.  I would like more testing.
- I only had time for part one and two of the challenge.  For part 3 I would probably just extend the .proto to include where the token is referenced and allow the client to eneter an array of modes.
- Modes should be an Enum so we have type saftety at the endpoint.
- I am missing edge cases everywhere. For example, what happens if the file is shorter than the linenumber and charnumber entered by the user?  What happens if they are negative.  Tests cases need to be added for thing like this so the response can be well defined.
- I need to add heath/ready/version endpoints
- There is no versioning
- I did not write any orchestration
