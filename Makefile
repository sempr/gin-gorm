test:
	GIN_MODE=release APP_DB_TYPE=sqlite3 APP_DB_URI=":memory:" go test -v

