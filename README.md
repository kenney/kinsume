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


## Example output
---
```
$ go run kinsume.go -stream mystream
Working on kinesis stream: dontcrossmebro:shardId-000000000000
Could not find AWS Region in ENV. Please configure your ENV for AWS access
Streams:  {
  HasMoreStreams: false,
  StreamNames: [
    "mystream",
    "mysweetstream"
  ]
}
Stream Details:  {
  StreamDescription: {
    HasMoreShards: false,
    Shards: [{
        HashKeyRange: {
          EndingHashKey: "987654321",
          StartingHashKey: "0"
        },
        SequenceNumberRange: {
          StartingSequenceNumber: "123456789"
        },
        ShardId: "shardId-000000000000"
      }],
    StreamARN: "arn:aws:kinesis:us-east-1:000000000000:stream/mystream",
    StreamName: "mystream",
    StreamStatus: "ACTIVE"
  }
}
```
