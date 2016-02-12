package golax

import "net/http"

/**
 * `Extended Writer`
 * It wraps http.ResponseWriter to add: StatusCode & Length
 */
type ExtendedWriter struct {
	StatusCode int
	Length     int
	http.ResponseWriter
}

func NewExtendedWriter(w http.ResponseWriter) *ExtendedWriter {
	return &ExtendedWriter{200, 0, w}
}

func (this *ExtendedWriter) Write(p []byte) (int, error) {
	n, err := this.ResponseWriter.Write(p)
	this.Length += n
	return n, err
}

func (this *ExtendedWriter) WriteHeader(statusCode int) {
	this.StatusCode = statusCode
	this.ResponseWriter.WriteHeader(statusCode)
}
