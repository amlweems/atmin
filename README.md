# Automatic Test Case Minimizer

`atmin` is a generic test case minimization program which accepts a test case,
executor, and validator, and returns a minimized blob of data which passes the
validation check. Included in the `cmd/` directory are a few examples of using
`atmin`.

## `hatmin`

The first example of `atmin`, an HTTP request minimization program, allows
users to provide a single HTTP request and receive a minimized payload which
generally returns the same output. Validation is performed using string
matching, requiring the user to provide a sample string which should be returned
by the server when the payload succeeds. An example is shown below.

```bash
$ cat req.txt
GET / HTTP/1.1
Host: example.org
Connection: keep-alive
Pragma: no-cache
Cache-Control: no-cache
Accept: text/plain, */*; q=0.01
X-Requested-With: XMLHttpRequest
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36
Content-Type: application/json; charset=utf-8
Accept-Language: en-US,en;q=0.9
Cookie: _ga=GA1.2.1953734269.1506457661; _gid=GA1.2.593597176.1512947963; _gat=1

$ hatmin -request req.txt -url http://example.org -dry-run | grep '<h1>'
    <h1>Example Domain</h1>
$ hatmin -request req.txt -url http://example.org -needle 'Example Domain'
2017/12/10 20:05:46 Stage 0: Block Normalization
2017/12/10 20:05:46 Stage 1: Block Deletion
2017/12/10 20:05:49 Stage 2: Alphabet Minimization
2017/12/10 20:05:49 Stage 3: Character Minimization
2017/12/10 20:05:49 Stage 1: Block Deletion
2017/12/10 20:05:52 Stage 2: Alphabet Minimization
2017/12/10 20:05:52 Stage 3: Character Minimization
2017/12/10 20:05:52 Stage 1: Block Deletion
2017/12/10 20:05:52 Stage 2: Alphabet Minimization
2017/12/10 20:05:53 Stage 3: Character Minimization
GET / HTTP/0.0
0:0


```
