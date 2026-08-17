package generation

import (
	"testing"

	"github.com/markbeep/syncfo/templates"
)

func makeSummaries(n int) []page.ThreadSummary {
	out := make([]page.ThreadSummary, n)
	for i := range out {
		out[i] = page.ThreadSummary{ID: int64(i + 1)}
	}
	return out
}

func ids(list []page.ThreadSummary) []int64 {
	out := make([]int64, len(list))
	for i, s := range list {
		out[i] = s.ID
	}
	return out
}

func eqInt64(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func seqDesc(hi, lo int) []int64 {
	out := make([]int64, hi-lo+1)
	for i := range out {
		out[i] = int64(hi - i)
	}
	return out
}

func TestResolvePageUnder20(t *testing.T) {
	all := makeSummaries(15)
	data := ResolvePage(all, 0)
	if data.StablePages != 0 {
		t.Fatalf("StablePages = %d, want 0", data.StablePages)
	}
	if !eqInt64(ids(data.Threads), seqDesc(15, 1)) {
		t.Errorf("page 0 threads = %v, want 15..1", ids(data.Threads))
	}
	if len(data.Preview) != 0 {
		t.Errorf("Preview = %v, want empty", ids(data.Preview))
	}
	if data.NextNumber != -1 || data.PrevNumber != -1 {
		t.Errorf("Next/Prev = %d/%d, want -1/-1", data.NextNumber, data.PrevNumber)
	}

	out := ResolvePage(all, 1)
	if len(out.Threads) != 0 || out.NextNumber != -1 || out.PrevNumber != -1 {
		t.Errorf("out-of-range page = %+v, want empty", out)
	}
}

func TestResolvePageExactly20(t *testing.T) {
	all := makeSummaries(20)
	data := ResolvePage(all, 0)
	if data.StablePages != 1 {
		t.Fatalf("StablePages = %d, want 1", data.StablePages)
	}
	if len(data.Threads) != 0 {
		t.Errorf("page 0 threads = %v, want empty", ids(data.Threads))
	}
	if len(data.Preview) != 0 {
		t.Errorf("Preview = %v, want empty", ids(data.Preview))
	}
	if data.NextNumber != 1 || data.PrevNumber != -1 {
		t.Errorf("page 0 Next/Prev = %d/%d, want 1/-1", data.NextNumber, data.PrevNumber)
	}

	p1 := ResolvePage(all, 1)
	if !eqInt64(ids(p1.Threads), seqDesc(20, 1)) {
		t.Errorf("page 1 threads = %v, want 20..1", ids(p1.Threads))
	}
	if p1.NextNumber != -1 || p1.PrevNumber != -1 {
		t.Errorf("page 1 Next/Prev = %d/%d, want -1/-1", p1.NextNumber, p1.PrevNumber)
	}
}

func TestResolvePage25(t *testing.T) {
	all := makeSummaries(25)
	data := ResolvePage(all, 0)
	if data.StablePages != 1 {
		t.Fatalf("StablePages = %d, want 1", data.StablePages)
	}
	if !eqInt64(ids(data.Threads), seqDesc(25, 21)) {
		t.Errorf("page 0 threads = %v, want 25..21", ids(data.Threads))
	}
	if !eqInt64(ids(data.Preview), seqDesc(20, 6)) {
		t.Errorf("page 0 preview = %v, want 20..6", ids(data.Preview))
	}
	if data.NextNumber != 1 || data.PrevNumber != -1 {
		t.Errorf("page 0 Next/Prev = %d/%d, want 1/-1", data.NextNumber, data.PrevNumber)
	}

	p1 := ResolvePage(all, 1)
	if !eqInt64(ids(p1.Threads), seqDesc(20, 1)) {
		t.Errorf("page 1 threads = %v, want 20..1", ids(p1.Threads))
	}
	if p1.NextNumber != -1 || p1.PrevNumber != 0 {
		t.Errorf("page 1 Next/Prev = %d/%d, want -1/0", p1.NextNumber, p1.PrevNumber)
	}
}

func TestResolvePage45(t *testing.T) {
	all := makeSummaries(45)

	p0 := ResolvePage(all, 0)
	if p0.StablePages != 2 {
		t.Fatalf("StablePages = %d, want 2", p0.StablePages)
	}
	if !eqInt64(ids(p0.Threads), seqDesc(45, 41)) {
		t.Errorf("page 0 threads = %v, want 45..41", ids(p0.Threads))
	}
	if !eqInt64(ids(p0.Preview), seqDesc(40, 26)) {
		t.Errorf("page 0 preview = %v, want 40..26", ids(p0.Preview))
	}
	if p0.NextNumber != 2 || p0.PrevNumber != -1 {
		t.Errorf("page 0 Next/Prev = %d/%d, want 2/-1", p0.NextNumber, p0.PrevNumber)
	}

	p1 := ResolvePage(all, 1)
	if !eqInt64(ids(p1.Threads), seqDesc(20, 1)) {
		t.Errorf("page 1 threads = %v, want 20..1", ids(p1.Threads))
	}
	if p1.NextNumber != -1 || p1.PrevNumber != 2 {
		t.Errorf("page 1 Next/Prev = %d/%d, want -1/2", p1.NextNumber, p1.PrevNumber)
	}

	p2 := ResolvePage(all, 2)
	if !eqInt64(ids(p2.Threads), seqDesc(40, 21)) {
		t.Errorf("page 2 threads = %v, want 40..21", ids(p2.Threads))
	}
	if p2.NextNumber != 1 || p2.PrevNumber != 0 {
		t.Errorf("page 2 Next/Prev = %d/%d, want 1/0", p2.NextNumber, p2.PrevNumber)
	}

	p3 := ResolvePage(all, 3)
	if len(p3.Threads) != 0 || p3.NextNumber != -1 || p3.PrevNumber != -1 {
		t.Errorf("out-of-range page 3 = %+v, want empty", p3)
	}
}
