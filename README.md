kinsume
=======

kinsume is a simple AWS Kinesis stream consumer written in go. It will continually listen to a stream and give you the data.

## Usage
---
```
Usage: kinsume [args]
-delay int
   time to sleep between GetRecord calls (default 5)
-shard string
   what shard to tail (default "shardId-000000000000")
-stream string
   what stream to tail (default "mystream")
```
