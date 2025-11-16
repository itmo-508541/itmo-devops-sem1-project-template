package report

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type RequestFilter struct {
	Start string `validate:"omitempty,date"`
	End   string `validate:"omitempty,date"`
	Min   int64  `validate:"omitempty,gte=0"`
	Max   int64  `validate:"omitempty,gte=0"`
}

func NewRequestFilter(r *http.Request) RequestFilter {
	q := r.URL.Query()
	f := RequestFilter{
		Start: q.Get("start"),
		End:   q.Get("end"),
	}

	if i, ok := f.getInt(q.Get("min")); ok {
		f.Min = i
	}
	if i, ok := f.getInt(q.Get("max")); ok {
		f.Max = i
	}

	return f
}

func (f RequestFilter) getInt(value string) (result int64, ok bool) {
	if len(value) > 0 {
		i, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			result = i
		} else {
			result = -1
		}
		ok = true
	}
	return result, ok
}

func (f RequestFilter) Where() (args pgx.NamedArgs, where string, ok bool) {
	args = make(pgx.NamedArgs)
	cond := make([]string, 0)

	if len(f.Start) > 0 {
		args["start"] = f.Start
		cond = append(cond, "create_date >= @start")
	}
	if len(f.End) > 0 {
		args["end"] = f.End
		cond = append(cond, "create_date <= @end")
	}
	if f.Min > 0 {
		args["min"] = f.Min
		cond = append(cond, "price >= @min")
	}
	if f.Max > 0 {
		args["max"] = f.Max
		cond = append(cond, "price <= @max")
	}
	if len(cond) > 0 {
		where = strings.Join(cond, " AND ")
		ok = true
	}

	return args, where, ok
}
