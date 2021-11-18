module golang.com/squishd/christmaslist

go 1.17

require (
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
	github.com/squishd/usersession v0.0.0-20211010052131-852c2483e60d // indirect
)

require github.com/squishd/authentication v0.0.0
replace github.com/squishd/authentication v0.0.0 => ./authentication