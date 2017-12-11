package atmin

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPMinimizer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			fmt.Fprintln(w, `{"error":"invalid content type"}`)
			return
		}

		if r.Host != "api.example.org" {
			fmt.Fprintln(w, `{"error":"invalid host"}`)
			return
		}

		cook, err := r.Cookie("atk_token")
		if err != nil {
			fmt.Fprintln(w, `{"error":"missing atk_token"}`)
			return
		}
		if cook.Value != "VmKG63mFuFCtjR2ZAnlWTmu3HO2zjvpaQt1UR8KRmPI" {
			fmt.Fprintln(w, `{"error":"invalid atk_token"}`)
			return
		}

		fmt.Fprintf(w, `{"username":"admin","error":""}`)
	}))
	defer ts.Close()

	// large request body with extraneous content
	in := []byte(`GET /api/v1/examples/http HTTP/1.1
Host: api.example.org
Content-Type: application/json
Connection: keep-alive
Pragma: no-cache
Cache-Control: no-cache
Accept: text/plain, */*; q=0.01
X-Requested-With: XMLHttpRequest
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36
Content-Type: application/json; charset=utf-8
Referer: https://app.example.org/examples/http
Accept-Language: en-US,en;q=0.9
Cookie: __stripe_mid=fa7d36a5-7148-41f2-89cb-e798f76eabfe; __qca=P0-1719821274-1506457757963; signin_return_url=%252F; atk_token=VmKG63mFuFCtjR2ZAnlWTmu3HO2zjvpaQt1UR8KRmPI; _ga=GA1.2.1953734269.1506457661; _gid=GA1.2.593597176.1512947963; _gat=1

`)

	addr := strings.Split(ts.URL, "//")[1]
	m := NewMinimizer(in).ExecuteHTTP(addr, false).ValidateString(`"username":"admin"`)
	min := m.Minimize()

	if len(min) > len(in)/4 {
		t.Error("minimized request not small enough")
	}
	t.Logf("minimized request: %s", min)

	ex := HTTPExecutor{Addr: addr, TLS: false}
	out := ex.Execute(min)
	if !bytes.Contains(out, []byte(`"username":"admin"`)) {
		t.Errorf("output did not match: %s", out)
	}
}
