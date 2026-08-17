package generation

import (
	"github.com/a-h/templ"
	"github.com/markbeep/syncfo/templates"
)

// PageSize is the number of threads per list page and messages per thread page.
const PageSize = 20

// ResolvePage maps an ordered (oldest -> newest) list of thread summaries and a
// page number to the render data for that list page. Page 0 is the temporary
// newest page; stable pages are 1-based and immutable. The returned Threads and
// Preview slices are ordered newest -> oldest for display.
func ResolvePage(all []page.ThreadSummary, number int) page.ListPageData {
	total := len(all)
	stable := total / PageSize
	rem := total % PageSize

	data := page.ListPageData{
		Number:      number,
		StablePages: stable,
		NextNumber:  -1,
		PrevNumber:  -1,
	}

	if number == 0 {
		if stable == 0 {
			data.Threads = reverseSummaries(all)
			return data
		}
		data.Threads = reverseSummaries(all[total-rem : total])
		if rem > 0 {
			data.Preview = reverseSummaries(all[total-PageSize : total-rem])
		}
		data.NextNumber = stable
		return data
	}

	if number < 1 || number > stable {
		return data
	}

	data.Threads = reverseSummaries(all[(number-1)*PageSize : number*PageSize])
	if number > 1 {
		data.NextNumber = number - 1
	}
	if number < stable {
		data.PrevNumber = number + 1
	} else if rem > 0 {
		data.PrevNumber = 0
	}
	return data
}

// reverseSummaries returns a reversed (newest -> oldest) copy of in.
func reverseSummaries(in []page.ThreadSummary) []page.ThreadSummary {
	out := make([]page.ThreadSummary, len(in))
	for i, s := range in {
		out[len(in)-1-i] = s
	}
	return out
}

// ListPage renders a single list page.
func ListPage(data page.ListPageData, basePath string) templ.Component {
	return page.List(data, basePath)
}
