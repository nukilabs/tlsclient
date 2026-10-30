package tests

import (
	"context"
	"errors"
	"io"
	stdhttp "net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nukilabs/http"
	"github.com/nukilabs/tlsclient"
	"github.com/nukilabs/tlsclient/profiles"
)

const hookWait = 5 * time.Second

// barrier releases its waiters once n of them have arrived. A waiter gives up
// after hookWait, so a hook that never overlaps fails the test instead of
// hanging it.
type barrier struct {
	n       int32
	arrived atomic.Int32
	full    chan struct{}
}

func newBarrier(n int) *barrier {
	return &barrier{n: int32(n), full: make(chan struct{})}
}

func (b *barrier) wait() error {
	if b.arrived.Add(1) == b.n {
		close(b.full)
	}
	select {
	case <-b.full:
		return nil
	case <-time.After(hookWait):
		return errors.New("timed out waiting for concurrent hooks")
	}
}

func hookServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.URL.Path == "/redirect" {
			stdhttp.Redirect(w, r, "/target", stdhttp.StatusFound)
			return
		}
		io.WriteString(w, "ok")
	}))
	t.Cleanup(srv.Close)
	return srv
}

func hookClient(t *testing.T) *tlsclient.Client {
	t.Helper()
	c := tlsclient.New(profiles.Chrome154)
	t.Cleanup(c.CloseIdleConnections)
	return c
}

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

func get(c doer, url string) error {
	// Its own context, so nothing links a nested request to the hook's request.
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = io.Copy(io.Discard, res.Body)
	return err
}

func TestHooksRunForEveryConcurrentRequest(t *testing.T) {
	const requests = 32

	srv := hookServer(t)
	c := hookClient(t)

	var counts [4]atomic.Int32
	barriers := [4]*barrier{}
	for i := range barriers {
		barriers[i] = newBarrier(requests)
	}
	pre := func(i int) tlsclient.PreHook {
		return func(_ *tlsclient.Hook, req *http.Request) (*http.Request, error) {
			counts[i].Add(1)
			return req, barriers[i].wait()
		}
	}
	post := func(i int) tlsclient.PostHook {
		return func(_ *tlsclient.Hook, _ *http.Request, res *http.Response) (*http.Response, error) {
			counts[i].Add(1)
			return res, barriers[i].wait()
		}
	}
	c.SetPreHooks(pre(0), pre(1))
	c.SetPostHooks(post(2), post(3))

	var wg sync.WaitGroup
	errs := make(chan error, requests)
	for range requests {
		wg.Go(func() {
			if err := get(c, srv.URL); err != nil {
				errs <- err
			}
		})
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
	for i := range counts {
		if got := counts[i].Load(); got != requests {
			t.Errorf("hook %d ran %d times, want %d", i, got, requests)
		}
	}
}

func TestNestedRequestSkipsHooks(t *testing.T) {
	srv := hookServer(t)
	c := hookClient(t)

	var pre, post atomic.Int32
	c.SetPreHooks(func(h *tlsclient.Hook, req *http.Request) (*http.Request, error) {
		pre.Add(1)
		return req, get(h, srv.URL)
	})
	c.SetPostHooks(func(_ *tlsclient.Hook, _ *http.Request, res *http.Response) (*http.Response, error) {
		post.Add(1)
		return res, nil
	})

	if err := get(c, srv.URL); err != nil {
		t.Fatal(err)
	}
	if got := pre.Load(); got != 1 {
		t.Errorf("pre-hook ran %d times, want 1", got)
	}
	if got := post.Load(); got != 1 {
		t.Errorf("post-hook ran %d times, want 1", got)
	}
}

func TestClientRequestFromHookRunsHooks(t *testing.T) {
	srv := hookServer(t)
	c := hookClient(t)

	var pre, post atomic.Int32
	c.SetPreHooks(func(_ *tlsclient.Hook, req *http.Request) (*http.Request, error) {
		pre.Add(1)
		return req, nil
	})
	c.SetPostHooks(func(h *tlsclient.Hook, _ *http.Request, res *http.Response) (*http.Response, error) {
		if post.Add(1) > 1 {
			return res, nil
		}
		nested, err := h.Client().Post(srv.URL, "text/plain", nil)
		if err != nil {
			return nil, err
		}
		nested.Body.Close()
		return res, nil
	})

	if err := get(c, srv.URL); err != nil {
		t.Fatal(err)
	}
	if got := pre.Load(); got != 2 {
		t.Errorf("pre-hook ran %d times, want 2", got)
	}
	if got := post.Load(); got != 2 {
		t.Errorf("post-hook ran %d times, want 2", got)
	}
}

func TestSlowHookDoesNotBlockOtherRequests(t *testing.T) {
	srv := hookServer(t)
	c := hookClient(t)

	entered := make(chan struct{})
	release := make(chan struct{})
	c.SetPreHooks(func(_ *tlsclient.Hook, req *http.Request) (*http.Request, error) {
		if req.URL.Path == "/slow" {
			close(entered)
			<-release
		}
		return req, nil
	})

	slow := make(chan error, 1)
	go func() { slow <- get(c, srv.URL+"/slow") }()
	<-entered

	fast := make(chan error, 1)
	go func() { fast <- get(c, srv.URL+"/fast") }()

	select {
	case err := <-fast:
		if err != nil {
			t.Error(err)
		}
	case <-time.After(hookWait):
		t.Error("request blocked behind another request's slow hook")
	}
	close(release)
	if err := <-slow; err != nil {
		t.Error(err)
	}
}

func TestRedirectHooks(t *testing.T) {
	srv := hookServer(t)

	newClient := func(t *testing.T) (*tlsclient.Client, *atomic.Int32) {
		c := hookClient(t)
		c.SetRedirectFunc(func(*http.Request, []*http.Request) error { return nil })
		return c, &atomic.Int32{}
	}

	t.Run("top level request runs pre-hooks on each hop", func(t *testing.T) {
		c, count := newClient(t)
		c.SetPreHooks(func(_ *tlsclient.Hook, req *http.Request) (*http.Request, error) {
			count.Add(1)
			return req, nil
		})
		if err := get(c, srv.URL+"/redirect"); err != nil {
			t.Fatal(err)
		}
		if got := count.Load(); got != 2 {
			t.Errorf("pre-hook ran %d times, want 2", got)
		}
	})

	t.Run("nested request skips pre-hooks on redirect", func(t *testing.T) {
		c, count := newClient(t)
		c.SetPreHooks(func(h *tlsclient.Hook, req *http.Request) (*http.Request, error) {
			if count.Add(1) > 1 {
				return req, nil
			}
			return req, get(h, srv.URL+"/redirect")
		})
		if err := get(c, srv.URL); err != nil {
			t.Fatal(err)
		}
		if got := count.Load(); got != 1 {
			t.Errorf("pre-hook ran %d times, want 1", got)
		}
	})

	t.Run("pre-hook on a redirect is protected from recursion", func(t *testing.T) {
		c, count := newClient(t)
		c.SetPreHooks(func(h *tlsclient.Hook, req *http.Request) (*http.Request, error) {
			if count.Add(1) > 3 {
				return nil, errors.New("hook recursed")
			}
			if req.URL.Path != "/target" {
				return req, nil
			}
			return req, get(h, srv.URL+"/redirect")
		})
		if err := get(c, srv.URL+"/redirect"); err != nil {
			t.Fatal(err)
		}
		if got := count.Load(); got != 2 {
			t.Errorf("pre-hook ran %d times, want 2", got)
		}
	})
}

func TestHookDoEndsWithTheRequest(t *testing.T) {
	srv := hookServer(t)
	c := hookClient(t)
	c.SetRedirectFunc(func(*http.Request, []*http.Request) error { return nil })

	var pre, post atomic.Int32
	c.SetPreHooks(func(_ *tlsclient.Hook, req *http.Request) (*http.Request, error) {
		pre.Add(1)
		return req, nil
	})
	c.SetPostHooks(func(h *tlsclient.Hook, _ *http.Request, res *http.Response) (*http.Response, error) {
		if post.Add(1) > 1 {
			return res, nil
		}
		nested, err := http.NewRequest(http.MethodGet, srv.URL+"/redirect", nil)
		if err != nil {
			return nil, err
		}
		skipped, err := h.Do(nested)
		if err != nil {
			return nil, err
		}
		skipped.Body.Close()
		if got := pre.Load(); got != 1 {
			t.Errorf("Hook.Do ran pre-hooks: %d runs, want 1", got)
		}
		hooked, err := h.Client().Do(nested)
		if err != nil {
			return nil, err
		}
		hooked.Body.Close()
		return res, nil
	})

	if err := get(c, srv.URL); err != nil {
		t.Fatal(err)
	}
	// The top level request, then both hops of the request sent through the client.
	if got := pre.Load(); got != 3 {
		t.Errorf("pre-hook ran %d times, want 3", got)
	}
}

// The skip must not outlive the request: a context taken from the response runs
// redirect hooks like any other.
func TestHookDoLeavesNoMarkOnResponse(t *testing.T) {
	srv := hookServer(t)
	c := hookClient(t)
	// A timeout cancels the response's context when its body closes.
	c.Timeout = 0
	c.SetRedirectFunc(func(*http.Request, []*http.Request) error { return nil })

	var pre, post atomic.Int32
	c.SetPreHooks(func(_ *tlsclient.Hook, req *http.Request) (*http.Request, error) {
		pre.Add(1)
		return req, nil
	})
	c.SetPostHooks(func(h *tlsclient.Hook, _ *http.Request, res *http.Response) (*http.Response, error) {
		if post.Add(1) > 1 {
			return res, nil
		}
		nested, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		if err != nil {
			return nil, err
		}
		skipped, err := h.Do(nested)
		if err != nil {
			return nil, err
		}
		skipped.Body.Close()

		reused, err := http.NewRequestWithContext(skipped.Request.Context(), http.MethodGet, srv.URL+"/redirect", nil)
		if err != nil {
			return nil, err
		}
		hooked, err := h.Client().Do(reused)
		if err != nil {
			return nil, err
		}
		hooked.Body.Close()
		return res, nil
	})

	if err := get(c, srv.URL); err != nil {
		t.Fatal(err)
	}
	// The top level request, then both hops of the reused context.
	if got := pre.Load(); got != 3 {
		t.Errorf("pre-hook ran %d times, want 3", got)
	}
}

func TestCloneRedirects(t *testing.T) {
	srv := hookServer(t)

	setup := func(t *testing.T) (*tlsclient.Client, *atomic.Int32) {
		c := hookClient(t)
		c.SetRedirectFunc(func(*http.Request, []*http.Request) error { return nil })
		var original atomic.Int32
		c.SetPreHooks(func(_ *tlsclient.Hook, req *http.Request) (*http.Request, error) {
			original.Add(1)
			return req, nil
		})
		return c, &original
	}

	t.Run("run the clone's hooks on the clone", func(t *testing.T) {
		c, original := setup(t)
		clone := c.Clone()
		t.Cleanup(clone.CloseIdleConnections)

		var hops, wrongClient atomic.Int32
		clone.SetPreHooks(func(h *tlsclient.Hook, req *http.Request) (*http.Request, error) {
			hops.Add(1)
			if h.Client() != clone {
				wrongClient.Add(1)
			}
			return req, nil
		})

		if err := get(clone, srv.URL+"/redirect"); err != nil {
			t.Fatal(err)
		}
		if got := original.Load(); got != 0 {
			t.Errorf("original's pre-hook ran %d times, want 0", got)
		}
		if got := hops.Load(); got != 2 {
			t.Errorf("clone's pre-hook ran %d times, want 2", got)
		}
		if got := wrongClient.Load(); got != 0 {
			t.Errorf("hook saw the wrong client %d times", got)
		}
	})

	t.Run("skip hooks for a nested request on the clone", func(t *testing.T) {
		c, _ := setup(t)
		clone := c.Clone()
		t.Cleanup(clone.CloseIdleConnections)

		var count atomic.Int32
		clone.SetPreHooks(func(h *tlsclient.Hook, req *http.Request) (*http.Request, error) {
			if count.Add(1) > 1 {
				return req, nil
			}
			return req, get(h, srv.URL+"/redirect")
		})

		if err := get(clone, srv.URL); err != nil {
			t.Fatal(err)
		}
		if got := count.Load(); got != 1 {
			t.Errorf("pre-hook ran %d times, want 1", got)
		}
	})

	t.Run("keep SetFollowRedirects(false)", func(t *testing.T) {
		c, _ := setup(t)
		c.SetFollowRedirects(false)
		clone := c.Clone()
		t.Cleanup(clone.CloseIdleConnections)

		res, err := clone.Get(srv.URL + "/redirect")
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusFound {
			t.Errorf("status %d, want %d", res.StatusCode, http.StatusFound)
		}

		clone.SetFollowRedirects(true)
		res, err = clone.Get(srv.URL + "/redirect")
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Errorf("status %d, want %d", res.StatusCode, http.StatusOK)
		}
	})
}

func TestHookReachesItsClient(t *testing.T) {
	srv := hookServer(t)
	c := hookClient(t)

	var same atomic.Bool
	c.SetPreHooks(func(h *tlsclient.Hook, req *http.Request) (*http.Request, error) {
		same.Store(h.Client() == c)
		// Post-hooks added by a pre-hook apply to the same request.
		h.Client().AddPostHooks(func(_ *tlsclient.Hook, _ *http.Request, res *http.Response) (*http.Response, error) {
			res.Header.Set("X-Added", "1")
			return res, nil
		})
		return req, nil
	})

	res, err := c.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if !same.Load() {
		t.Error("Hook.Client is not the client the hook runs on")
	}
	if res.Header.Get("X-Added") != "1" {
		t.Error("AddPostHooks from a pre-hook did not apply to the request")
	}
}

// Hooks run in parallel, so several may add hooks at once.
func TestConcurrentHooksCanAddHooks(t *testing.T) {
	const requests = 16

	srv := hookServer(t)
	c := hookClient(t)

	bar := newBarrier(requests)
	var added atomic.Int32
	c.SetPreHooks(func(h *tlsclient.Hook, req *http.Request) (*http.Request, error) {
		if err := bar.wait(); err != nil {
			return nil, err
		}
		h.Client().AddPostHooks(func(_ *tlsclient.Hook, _ *http.Request, res *http.Response) (*http.Response, error) {
			added.Add(1)
			return res, nil
		})
		return req, nil
	})

	var wg sync.WaitGroup
	for range requests {
		wg.Go(func() {
			if err := get(c, srv.URL); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()

	// One more request runs every hook that was added, so a lost append shows.
	added.Store(0)
	c.DeletePreHooks()
	if err := get(c, srv.URL); err != nil {
		t.Fatal(err)
	}
	if got := added.Load(); got != requests {
		t.Errorf("%d hooks added, want %d", got, requests)
	}
}
