package core

import (
	"io"
)

func copyWithProgress(dst io.Writer, src io.Reader, total int64, cb func(downloaded, total int64)) (int64, error) {
	const bufSize = 32 * 1024
	buf := make([]byte, bufSize)
	var downloaded int64
	for {
		n, readErr := src.Read(buf)
		if n > 0 {
			wn, writeErr := dst.Write(buf[:n])
			downloaded += int64(wn)
			if cb != nil {
				cb(downloaded, total)
			}
			if writeErr != nil {
				return downloaded, writeErr
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				return downloaded, nil
			}
			return downloaded, readErr
		}
	}
}