package extension

import (
	"bufio"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/wangjuelong/go-streams"
	"github.com/wangjuelong/go-streams/flow"
	"github.com/wangjuelong/go-streams/util/ospkg"
)

// FileSource streams data from the file
type FileSource struct {
	fileName string
	in       chan interface{}
}

// NewFileSource returns a new FileSource instance
func NewFileSource(fileName string) *FileSource {
	source := &FileSource{fileName, make(chan interface{})}
	source.init()
	return source
}

func (fs *FileSource) init() {
	go func() {
		file, err := os.Open(fs.fileName)
		if err != nil {
			log.Error(errors.Wrap(err, "OpenFile"))
			return
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			l, isPrefix, err := reader.ReadLine()
			if err != nil {
				close(fs.in)
				break
			}

			var msg string
			if isPrefix {
				msg = string(l)
			} else {
				msg = string(l) + ospkg.NewLine
			}

			fs.in <- msg
		}
	}()
}

// Via streams data through the given flow
func (fs *FileSource) Via(_flow streams.Flow) streams.Flow {
	flow.DoStream(fs, _flow)
	return _flow
}

// Out returns an output channel for sending data
func (fs *FileSource) Out() <-chan interface{} {
	return fs.in
}

// A FileSink writes items to a file
type FileSink struct {
	fileName string
	in       chan interface{}
}

// NewFileSink returns a new FileSink instance
func NewFileSink(fileName string) *FileSink {
	sink := &FileSink{fileName, make(chan interface{})}
	sink.init()
	return sink
}

func (fs *FileSink) init() {
	go func() {
		file, err := os.OpenFile(fs.fileName, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Error(errors.Wrap(err, "OpenFile"))
			return
		}
		defer file.Close()
		for elem := range fs.in {
			_, err = file.WriteString(elem.(string))
			if err != nil {
				log.Error(errors.Wrap(err, "WriteString"))
				return
			}
		}
	}()
}

// In returns an input channel for receiving data
func (fs *FileSink) In() chan<- interface{} {
	return fs.in
}
