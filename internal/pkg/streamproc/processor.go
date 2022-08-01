package streamproc

type Bucket struct {
	Ley interface{}
	Val interface{}
}

type Stream struct {
	 Buckets []*Bucket
	 Limit int
	 Offset int
}

type Processor struct {
	stream *Stream
}