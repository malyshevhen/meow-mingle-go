module github.com/malyshEvhen/meow_mingle

go 1.23.0

toolchain go1.24.0

require github.com/gorilla/mux v1.8.1

require github.com/felixge/httpsnoop v1.0.4 // indirect

require (
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

require (
	github.com/go-playground/validator/v10 v10.26.0
	github.com/goccy/go-yaml v1.18.0
	github.com/gocql/gocql v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/gorilla/handlers v1.5.2
	golang.org/x/crypto v0.39.0
)

replace github.com/gocql/gocql => github.com/scylladb/gocql v1.15.1
