package price

import "github.com/jackc/pgx/v5"

type SqlFilter interface {
	Where() (args pgx.NamedArgs, where string, ok bool)
}
