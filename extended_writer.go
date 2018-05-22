package golax

import "net/http"

// ExtendedWriter wraps http.ResponseWriter with StatusCode & Length
type ExtendedWriter struct {
	StatusCode     int
	statusCodeSent bool
	Length         int
	http.ResponseWriter
}

// NewExtendedWriter instances a new *ExtendedWriter
func NewExtendedWriter(w http.ResponseWriter) *ExtendedWriter {
	return &ExtendedWriter{
		StatusCode:     200,
		statusCodeSent: false,
		Length:         0,
		ResponseWriter: w,
	}
}

// Write replaces default behaviour of http.ResponseWriter
func (w *ExtendedWriter) Write(p []byte) (int, error) {
	n, err := w.ResponseWriter.Write(p)
	w.Length += n
	return n, err
}

// WriteHeader replaces default behaviour of http.ResponseWriter
func (w *ExtendedWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	if !w.statusCodeSent {
		w.ResponseWriter.WriteHeader(statusCode)
		w.statusCodeSent = true
	}
}
