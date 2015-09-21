package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

var streamName string
var shardName string
var delayTime int

func main() {
	flag.StringVar(&streamName, "stream", "mystream", "what stream to tail")
	flag.StringVar(&shardName, "shard", "shardId-000000000000", "what shard to tail")
	flag.IntVar(&delayTime, "delay", 5, "time to sleep between GetRecord calls")

	flag.Parse()

	println(fmt.Sprintf("Working on kinesis stream: %s:%s", streamName, shardName))

	if *defaults.DefaultConfig.Region == "" {
		println("Could not find AWS Region in ENV. Please configure your ENV for AWS access")
		os.Exit(1)
	}

	println(fmt.Sprintf("AWS Defaults: %#v", defaults.DefaultConfig))
	println(fmt.Sprintf("AWS Region: %s", *defaults.DefaultConfig.Region))

	describeStream(streamName)

	watchStream(streamName, shardName)
}

func watchStream(streamname string, shardname string) {
	si, err := getShardIterator(streamname, shardname)
	if err != nil {
		println(fmt.Sprintf("Error w/ the GetShardIterator call: %#v", err.Error()))
		return
	}

	//var delay = 5

	println(fmt.Sprintf("Tailing the stream with a %ds loop ", delayTime))
	for true {
		println("Tick...")
		records, nsi, err := getRecords(si)
		if err != nil {
			println(fmt.Sprintf("Error w/ the GetRecords call: %#v", err.Error()))
			return
		}

		if nsi == "" {
			println("No next shard iterator. Bailing on stream")
		}

		//println(fmt.Sprintf("Next Shard Iterator: %#v", nsi))
		si = nsi

		for _, r := range records {
			println(fmt.Sprintf("%s @ Seq#: %s => Val: %s", *r.PartitionKey, *r.SequenceNumber, string(r.Data)))
		}

		// Sleep for 5 seconds.
		time.Sleep(time.Millisecond * 1000 * time.Duration(delayTime))
	}
}

func getShardIterator(streamname string, shardname string) (string, error) {
	svc := kinesis.New(nil)

	// Start a tail on the latest posts.
	gsiiParams := &kinesis.GetShardIteratorInput{
		StreamName:        aws.String(streamname),
		ShardId:           aws.String(shardname),
		ShardIteratorType: aws.String("LATEST"),
		//ShardIteratorType: aws.String("TRIM_HORIZON"),
	}

	gsiresp, err := svc.GetShardIterator(gsiiParams)

	if err != nil {
		println(fmt.Sprintf("Error w/ GetShardIterator: %#v", err.Error()))
		return "", err
	}

	return *gsiresp.ShardIterator, nil
}

func getRecords(si string) ([]*kinesis.Record, string, error) {
	svc := kinesis.New(nil)

	griParams := &kinesis.GetRecordsInput{
		ShardIterator: aws.String(si),
		Limit:         aws.Int64(10),
	}

	grresp, err := svc.GetRecords(griParams)

	if err != nil {
		println(fmt.Sprintf("Error w/ GetRecords: %#v", err.Error()))
		return []*kinesis.Record{}, "", err
	}

	//println(fmt.Sprintf("Records:  %#v", grresp))

	return grresp.Records, *grresp.NextShardIterator, nil
}

func listStreams() {
	svc := kinesis.New(nil)

	lsiParams := &kinesis.ListStreamsInput{
		ExclusiveStartStreamName: aws.String("StreamName"),
		//Limit: aws.Int64(1),
	}
	lsresp, err := svc.ListStreams(lsiParams)

	if err != nil {
		println(fmt.Sprintf("Error w/ ListStreams: %#v", err.Error()))
		return
	}

	println(fmt.Sprintf("Streams:  %#v", lsresp))
}

func describeStream(streamname string) {
	svc := kinesis.New(nil)

	dsiParams := &kinesis.DescribeStreamInput{
		StreamName: aws.String(streamname),
	}
	dsresp, err := svc.DescribeStream(dsiParams)

	if err != nil {
		println(fmt.Sprintf("Error w/ DescribeStream: %#v", err.Error()))
		return
	}

	println(fmt.Sprintf("Stream Details:  %#v", dsresp))
}
