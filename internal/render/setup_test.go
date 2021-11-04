package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/chirper-go/internal/config"
	"github.com/andkolbe/chirper-go/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	gob.Register(models.Chirp{})

	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true 
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false 

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

// create an object using interfaces that satisfies the requirements for a response writer
type myWriter struct {}

// Header, WriteHeader, and Write are all requirements on the http.ResponseWriter interface
func (mw *myWriter) Header() http.Header {
	var h http.Header 
	return h 
}
func (mw *myWriter) WriteHeader(i int) {}
func (mw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
