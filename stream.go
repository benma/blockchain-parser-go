package main

import (
	"io"
	"os"
)

type stream struct {
	files    [](*os.File)
	position int64
}

type Streamer interface {
	Read(p []byte) (int, error)
	Skip(offset int64) (int64, error)
	Position() int64
}

func Stream(r ...*os.File) Streamer {
	return &stream{r, 0}
}

func (mr *stream) Skip(offset int64) (int64, error) {
	if len(mr.files) == 0 {
		return 0, io.EOF
	}

	var tn int64 = 0
	for len(mr.files) > 0 {
		fi, err := mr.files[0].Stat()
		if err != nil {
			return tn, err
		}

		cur, err := mr.files[0].Seek(0, 1)
		if err != nil {
			return tn, err
		}
		if cur+offset < fi.Size() {
			n, err := mr.files[0].Seek(offset, 1)
			tn += n
			mr.position += n
			return tn, err
		} else {
			mr.files = mr.files[1:]
			offset -= (fi.Size() - cur)
			tn += fi.Size()
			mr.position += fi.Size()
		}
	}
	return tn, nil
}

func (mr *stream) Read(p []byte) (int, error) {
	if len(mr.files) == 0 {
		return 0, io.EOF
	}

	cp := p
	var tn int = 0
	for len(mr.files) > 0 {
		n, err := mr.files[0].Read(cp)
		tn += n
		mr.position += int64(n)

		if err != nil && err != io.EOF {
			return tn, err
		}

		cp = cp[n:]
		if len(cp) == 0 {
			// filled up the buffer, done
			break
		}
		mr.files = mr.files[1:]
	}
	return tn, nil
}

func (r *stream) Position() int64 {
	return r.position
}
