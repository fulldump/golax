package golax

import "net/http"

/**
 * `Extended Writer`
 * It wraps http.ResponseWriter to add: StatusCode & Length
 */
type ExtendedWriter struct {
	StatusCode     int
	statusCodeSent bool
	Length         int
	http.ResponseWriter
}

func NewExtendedWriter(w http.ResponseWriter) *ExtendedWriter {
	return &ExtendedWriter{200, false, 0, w}
}

func (this *ExtendedWriter) Write(p []byte) (int, error) {
	n, err := this.ResponseWriter.Write(p)
	this.Length += n
	return n, err
}

func (this *ExtendedWriter) WriteHeader(statusCode int) {
	this.StatusCode = statusCode
	if !this.statusCodeSent {
		this.ResponseWriter.WriteHeader(statusCode)
		this.statusCodeSent = true
	}
}
