package tree_test

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	. "ecosia/tree"
)

// TestHandler tests the http handler. It is still a unit test.
func TestHandler(t *testing.T) {
	view := func(w io.Writer, tree string) error {
		_, _ = w.Write([]byte(tree))
		return nil
	}

	var (
		urlKey = "tree"
		tree   = "baobab"
	)

	handler := NewHandler(urlKey, view)

	t.Run("ok", func(t *testing.T) {
		resp := getRecorderResponse(urlKey, tree, handler)

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "text/html" {
			t.Fatal("document encoding is not text/html")
		}

		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != tree {
			t.Fatalf("body was expected to be '%s', got '%s'", tree, string(body))
		}
	})

	t.Run("empty tree", func(t *testing.T) {
		resp := getRecorderResponse(urlKey, "", handler)

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("missing tree", func(t *testing.T) {
		resp := getRecorderResponse(urlKey+"nope", "sasa", handler)

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("error in view", func(t *testing.T) {
		errStr := "errored"

		errorView := func(w io.Writer, tree string) error {
			return errors.New(errStr)
		}

		handler := NewHandler(urlKey, errorView)

		resp := getRecorderResponse(urlKey, "y", handler)

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		expected := fmt.Sprintf("error while writing the tree: %s\n", errStr)
		if string(body) != expected {
			t.Fatalf("body was expected to be '%s', got '%s'", expected, string(body))
		}
	})
}

func getRecorderResponse(urlKey string, tree string, handler http.Handler) *http.Response {
	req := httptest.NewRequest("GET", fmt.Sprintf("/?%s=%s", urlKey, tree), nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	return resp
}

// TestIntegration tests the actual server and its integration with the templating thing
func TestIntegration(t *testing.T) {
	var (
		urlKey = "favoriteTree"
		tree   = "baobab"
	)

	tpl, _ := template.New("tree").Parse(`
		{{- if . -}}
			It's nice to know that your favorite tree is a {{.}}
		{{- else -}}
			Please tell me your favorite tree
		{{- end -}}
	`)

	view := NewTemplateView(tpl)

	handler := NewHandler(urlKey, view)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	t.Run("ok", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/?%s=%s", ts.URL, urlKey, tree))
		if err != nil {
			log.Fatal(err)
		}

		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()

		var buf bytes.Buffer
		_ = view(&buf, tree)
		expected := buf.String()

		body := string(bodyBytes)
		if body != expected {
			t.Fatalf("expected returned body to be\n'%s'\n, got\n'%s'", expected, body)
		}
	})

	t.Run("wrong route", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/smth", ts.URL))
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("url param missing", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/?%s-invalid=%s", ts.URL, urlKey, tree))
		if err != nil {
			log.Fatal(err)
		}
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()

		var buf bytes.Buffer
		_ = view(&buf, "")
		expected := buf.String()

		body := string(bodyBytes)
		if body != expected {
			t.Fatalf("expected returned body to be\n'%s'\n, got\n'%s'", expected, body)
		}
	})

	t.Run("url param empty", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/?%s=", ts.URL, urlKey))
		if err != nil {
			log.Fatal(err)
		}
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()

		var buf bytes.Buffer
		_ = view(&buf, "")
		expected := buf.String()

		body := string(bodyBytes)
		if body != expected {
			t.Fatalf("expected returned body to be\n'%s'\n, got\n'%s'", expected, body)
		}
	})

	t.Run("default use case", func(t *testing.T) {
		// rebuild these:
		view := NewTemplateView(nil)
		handler := NewHandler(urlKey, view)
		ts := httptest.NewServer(handler)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/?%s=%s", ts.URL, urlKey, tree))
		if err != nil {
			log.Fatal(err)
		}

		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()

		body := string(bodyBytes)

		var buf bytes.Buffer
		_ = view(&buf, tree)

		expected := buf.String()
		if body != expected {
			t.Fatalf("expected returned body to be\n'%s'\n, got\n'%s'", expected, body)
		}
	})
}

func ExampleNewHandler() {
	var (
		urlKey = "favoriteTree"
		tree   = "猴麵包樹"
	)

	tpl, _ := template.New("test").Parse("<p>{{.}}</p>")

	view := NewTemplateView(tpl) // or use the default tpl: NewTemplateView(nil)
	handler := NewHandler(urlKey, view)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	tempUrl := fmt.Sprintf("%s/?%s=%s", ts.URL, urlKey, tree)
	res, err := http.Get(tempUrl)
	if err != nil {
		panic(err)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bodyBytes))
	// Output: <p>猴麵包樹</p>
}
