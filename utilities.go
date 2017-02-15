package summer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/night-codes/govalidator"
	"golang.org/x/crypto/sha3"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/gin-gonic/gin.v1/binding"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

type (
	obj map[string]interface{}
	arr []interface{}
)

// PackagePath returns file path of Summer package location
func PackagePath() string {
	return basepath
}

// templates plugin for converts to json
func jsoner(object interface{}) string {
	j, _ := json.Marshal(object)
	return string(j)
}

// templates plugin - dummy
func getSite(name string) interface{} {
	return ""
}

// PostBind binds data from post request and validates them
func PostBind(c *gin.Context, ret interface{}) bool {
	c.BindWith(ret, binding.Form)
	if _, err := govalidator.ValidateStruct(ret); err != nil {
		ers := []string{}
		for k, v := range govalidator.ErrorsByField(err) {
			ers = append(ers, k+": "+v)
		}
		c.String(400, strings.Join(ers, "\n"))
		return false
	}
	return true
}

func indexOf(arr interface{}, v interface{}) int {
	V := reflect.ValueOf(v)
	Arr := reflect.ValueOf(arr)
	if t := reflect.TypeOf(arr).Kind(); t != reflect.Slice && t != reflect.Array {
		panic("Type Error! Second argument must be an array or a slice.")
	}
	for i := 0; i < Arr.Len(); i++ {
		if Arr.Index(i).Interface() == V.Interface() {
			return i
		}
	}
	return -1
}

func getJSON(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// extend struct data except default zero values
func extend(to interface{}, from interface{}) {
	valueTo := reflect.ValueOf(to).Elem()
	valueFrom := reflect.ValueOf(from).Elem()

	if valueTo.Kind() != reflect.Struct || valueFrom.Kind() != reflect.Struct || valueTo.Type() != valueFrom.Type() {
		panic(`Expected pointers of structs (same types)`)
	}

	for i := 0; i < valueFrom.Type().NumField(); i++ {
		fValue := valueFrom.Field(i)
		tValue := valueTo.Field(i)
		if !tValue.CanSet() || reflect.DeepEqual(fValue.Interface(), reflect.Zero(fValue.Type()).Interface()) {
			continue
		}
		valueTo.Field(i).Set(fValue)
	}
}

// H3hash create sha512 hash of string
func H3hash(s string) string {
	h3 := sha3.New512()
	io.WriteString(h3, s)
	return fmt.Sprintf("%x", h3.Sum(nil))
}

func dummy(c *gin.Context) {
	c.Next()
}
func setCookie(c *gin.Context, name string, value string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		MaxAge:  32000000,
		Expires: time.Now().AddDate(1, 0, 0),
	})
}

// Env returns environment variable value (or default value if env.variable absent)
func Env(envName string, defaultValue string) (value string) {
	value = os.Getenv(envName)
	if len(value) == 0 {
		value = defaultValue
	}
	return
}
func stripSlashes(s string) string {
	if len(s) > 0 {
		if s[len(s)-1] == '/' {
			s = s[:len(s)-1]
		}
	}
	if len(s) > 0 {
		if s[0] == '/' {
			s = s[1:]
		}
	}
	return s
}

func uniqAppend(s1, s2 []string) []string {
	m := map[string]bool{}
	for _, v := range s1 {
		m[v] = true
	}
	for _, v := range s2 {
		m[v] = true
	}
	ret := []string{}
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

func isOverlap(s1, s2 []string) bool {
	m := map[string]bool{}
	for _, v := range s1 {
		m[v] = true
	}
	for _, v := range s2 {
		if m[v] {
			return true
		}
	}
	return false
}

func checkRights(panel *Panel, modR Rights, usrR Rights) bool {
	userActions := uniqAppend(panel.Groups.Get(usrR.Groups...), usrR.Actions)
	rightsEmpty := len(modR.Groups) == 0 && len(modR.Actions) == 0
	allow := (len(modR.Groups) > 0 && isOverlap(usrR.Groups, modR.Groups)) || (len(modR.Actions) > 0 && isOverlap(userActions, modR.Actions))

	return rightsEmpty || allow
}

func gzipper(c *gin.Context) {
	filepath := stripSlashes(c.Param("filepath"))
	if len(filepath) != 0 && strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		switch filepath {
		case "build/login.css", "build/login.js", "build/main.js", "build/style.css":
			c.Header("Content-Encoding", "gzip")
			c.Header("Vary", "Accept-Encoding")

			if strings.HasSuffix(c.Param("filepath"), ".css") {
				c.Header("Content-Type", "text/css; charset=utf-8")
			} else if strings.HasSuffix(c.Param("filepath"), ".js") {
				c.Header("Content-Type", "application/x-javascript")
			}
			path0 := PackagePath() + "/files/" + filepath + ".gz"
			c.File(path0)
			c.Abort()
			return
		}
	}

}
