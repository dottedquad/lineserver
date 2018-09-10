# LineServer

Yout can run it like go run *.go -p <portnum> <filename> <stride>

portnum is the port to listen on
filename is the fixed file to serve
stride is how many lines to skip between index file entries

###How does your system work? (if not addressed in comments in source)
At start, it scans the file and creates an index of line beginnings. This will contain lines/stride entries. The default stride is 8.

###How will your system perform as the number of requests per second increases?
I think the main bottleneck will be disk, but I haven't tried.

###How will your system perform with a 1 GB file? a 10 GB file? a 100 GB file?
The initial index file creation will take O(n) time and O(m) size (where n is the number of bytes and m is the number of lines). During normal runtime, File size shouldn't matter. I've used 64 bit integers for file offsets in the index file so it should handle >4GB files.

###How will your system handle very long lines, e.g. > 1GB
There is a linear search from the indexed line, so long lines before that can slow that down. But after that it should be as fast as it can be.

###How will your system handle files with a very large number of lines, preventing you from keeping a line number based index in memory?
In the pessemistic case, a 1TB file can have ~500 billion lines (one character and one line ending). In that case, the index file will have ~67 billion entries and consume about 500GB of disk space.

###What documentation, websites, papers, etc did you consult in doing this assignment?
Stack exchange, Google, Go documentation, https://coderwall.com/p/wohavg/creating-a-simple-tcp-server-in-go

###How long did you spend on this exercise?
Probably about 8 hours on combined on Saturday and Sunday.
I'd spend longer on tests, but I just want to be done. I've used dependency inversion everywhere, so unit testing should be fairly easy.
I started writing a unit test for IndexedLineGetter, but I need a memory based ReadWriteSeeker, and I don't want to worry about figuring that out and you said no external dependencies. I think you can see that I've thought about how to make a testable system and I hope that suffices.
