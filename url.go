package main

// ShortURL is the stored mapping between a short code and the destination URL
// it redirects to.
type ShortURL struct {
	// Code is the generated identifier that appears in the short URL path.
	Code string
	// OriginalURL is the validated destination the short URL redirects to.
	OriginalURL string
}
