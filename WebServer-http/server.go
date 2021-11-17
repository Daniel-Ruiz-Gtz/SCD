package main

// Import library
import (
	"./data"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Const
const (
	// Port
	PORT = ":9000"
	// Routers
	HOME = "/"
	FORM = "/form"
	DATA = "/data"
	GENERAL = "/general"
	STUDENT = "/student"
	SUBJECT = "/subject"
	// Header
	KEY = "Content-Type"
	VALUE = "text/html"
)

// Instance
var dataset data.AllData

/* loadingHTML Function
 * @param a string
 * @return string
 */
func loadingHTML(a string) string {
	html, _ := ioutil.ReadFile(a)
	return string(html)
}

/* root Function
 * @param res http.ResponseWriter
 * @param _ *http.Request
 */
func root(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set(KEY, VALUE)
	_, _ = fmt.Fprint(res, loadingHTML("index.html"))
}

/* form Function
 * @param res http.ResponseWriter
 * @param _ *http.Request
 */
func form(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set(KEY, VALUE)
	_, _ = fmt.Fprintf(res, loadingHTML("form.html"))
}

/* student Function
 * @param res http.ResponseWriter
 * @param req *http.Request
 */
func student(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)

	if err := req.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(res, "ParseForm() error %v", err)
		return
	}

	fmt.Println(req.PostForm)
	stu := req.FormValue("student")

	fmt.Println(stu)

	res.Header().Set(KEY, VALUE)
	_, _ = fmt.Fprintf(res, loadingHTML("student.html"), dataset.StudentAVG(stu))
}

/* subject Function
 * @param res http.ResponseWriter
 * @param req *http.Request
 */
func subject(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)

	if err := req.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(res, "ParseForm() error %v", err)
		return
	}

	fmt.Println(req.PostForm)
	sub := req.FormValue("subject")

	fmt.Println(sub)

	res.Header().Set(KEY, VALUE)
	_, _ = fmt.Fprintf(res, loadingHTML("subject.html"), dataset.SubjectAVG(sub))
}

/* general Function
 * @param res http.ResponseWriter
 * @param _ *http.Request
 */
func general(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set(KEY, VALUE)
	_,_ = fmt.Fprintf(res, loadingHTML("general.html"), dataset.GeneralAVG())
}

/* datasets Function
 * @param res http.ResponseWriter
 * @param req *http.Request
 */
func datasets(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)

	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			_, _ = fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}

		fmt.Println(req.PostForm)

		dts := 	data.Data{
			Student: req.FormValue("student"),
			Subject: req.FormValue("subject"),
			Grade: req.FormValue("grade")}

		dataset.Add(dts)
		fmt.Println(dataset)

		res.Header().Set(KEY, VALUE)
		_, _ = fmt.Fprintf(res, loadingHTML("register.html"), dts.Student)

	case "GET":
		res.Header().Set(KEY, VALUE)
		_, _ = fmt.Fprintf(res, loadingHTML("data.html"), dataset.String())
	}
}

// main Function
func main() {
	// Routers
	http.HandleFunc(HOME, root)
	http.HandleFunc(FORM, form)
	http.HandleFunc(DATA, datasets)
	http.HandleFunc(GENERAL, general)
	http.HandleFunc(STUDENT, student)
	http.HandleFunc(SUBJECT, subject)

	// Start Server
	fmt.Println("App listening at http://localhost" + PORT)
	_ = http.ListenAndServe(PORT, nil)
}
