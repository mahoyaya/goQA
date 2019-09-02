package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//DbrowWithLatest hogefuga
type DbrowWithLatest struct {
	ID      int    `json:"id"`
	Pid     int    `json:"parentid"`
	Title   string `json:"title"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Open    int    `json:"open"`
	Desired string `json:"desired_at"`
	Created string `json:"created_at"`
	Latest  string `json:"latest"`
}

//Dbrow hogefuga
type Dbrow struct {
	ID      int    `json:"id"`
	Pid     int    `json:"parentid"`
	Title   string `json:"title"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Open    int    `json:"open"`
	Desired string `json:"desired_at"`
	Created string `json:"created_at"`
	Latest  string `json:"latest"`
}

//Dbrows hogefuga
type Dbrows struct {
	list []Dbrow
}

func allList() *[]Dbrow {

	var sliceDbrow []Dbrow

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	q := `
        SELECT * FROM memo
	`
	rows, err := db.Query(q)

	defer rows.Close()
	for rows.Next() {
		//var id, pid, open int
		//var name, body, desired, created string

		row := Dbrow{}

		// カーソルから値を取得
		if err := rows.Scan(
			&row.ID,
			&row.Pid,
			&row.Title,
			&row.Name,
			&row.Body,
			&row.Open,
			&row.Desired,
			&row.Created); err != nil {
			log.Fatal("rows.Scan() ", err)
			return &sliceDbrow
		}
		sliceDbrow = append(sliceDbrow, row)

		//fmt.Printf("id: %d, pid: %d, name: %s, body: %s\n", row.Id, row.Pid, row.Name, row.Body)

	}

	return &sliceDbrow
}

func normalList() *[]DbrowWithLatest {

	var sliceDbrow []DbrowWithLatest

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	//  select t1.id,t2.id,t2.name,max(t2.created_at) from memo as t1 inner join memo as t2 on t1.id = t2.parentid GROUP BY t1.id;
	q := `
		SELECT a.id, a.parentid, a.title, a.name, a.body, a.open, IFNULL(MAX(b.desired_at), a.desired_at), DATE(MAX(IFNULL(b.created_at, a.created_at))), IFNULL(b.name, a.name) FROM memo as a
		LEFT OUTER JOIN memo as b ON a.id=b.parentid WHERE a.parentid=0 AND a.open=1 GROUP BY a.id ORDER BY a.created_at DESC
	`
	rows, err := db.Query(q)

	defer rows.Close()
	for rows.Next() {
		//var id, pid, open int
		//var name, body, desired, created string

		row := DbrowWithLatest{}

		// カーソルから値を取得
		if err := rows.Scan(
			&row.ID,
			&row.Pid,
			&row.Title,
			&row.Name,
			&row.Body,
			&row.Open,
			&row.Desired,
			&row.Created,
			&row.Latest); err != nil {
			log.Fatal("rows.Scan() ", err)
			return &sliceDbrow
		}
		sliceDbrow = append(sliceDbrow, row)

		//fmt.Printf("id: %d, pid: %d, name: %s, body: %s\n", row.Id, row.Pid, row.Name, row.Body)

	}

	return &sliceDbrow
}

func closedList() *[]DbrowWithLatest {

	var sliceDbrow []DbrowWithLatest

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	q := `
		SELECT a.id, a.parentid, a.title, a.name, a.body, a.open, IFNULL(MAX(b.desired_at), a.desired_at), DATE(MAX(IFNULL(b.created_at, a.created_at))), IFNULL(b.name, a.name) FROM memo as a
		LEFT OUTER JOIN memo as b ON a.id=b.parentid WHERE a.parentid=0 AND a.open=0 GROUP BY a.id
	`
	rows, err := db.Query(q)

	defer rows.Close()
	for rows.Next() {
		//var id, pid, open int
		//var name, body, desired, created string

		row := DbrowWithLatest{}

		// カーソルから値を取得
		if err := rows.Scan(
			&row.ID,
			&row.Pid,
			&row.Title,
			&row.Name,
			&row.Body,
			&row.Open,
			&row.Desired,
			&row.Created,
			&row.Latest); err != nil {
			log.Fatal("rows.Scan() ", err)
			return &sliceDbrow
		}
		sliceDbrow = append(sliceDbrow, row)

		//fmt.Printf("id: %d, pid: %d, name: %s, body: %s\n", row.Id, row.Pid, row.Name, row.Body)

	}

	return &sliceDbrow
}

func withChildList(pid int) *[]Dbrow {

	var sliceDbrow []Dbrow

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	q := `
        SELECT * FROM memo WHERE parentid=? OR id=?
	`
	rows, err := db.Query(q, pid, pid)

	defer rows.Close()
	for rows.Next() {
		//var id, pid, open int
		//var name, body, desired, created string

		row := Dbrow{}

		// カーソルから値を取得
		if err := rows.Scan(
			&row.ID,
			&row.Pid,
			&row.Title,
			&row.Name,
			&row.Body,
			&row.Open,
			&row.Desired,
			&row.Created); err != nil {
			log.Fatal("rows.Scan() ", err)
			return &sliceDbrow
		}
		sliceDbrow = append(sliceDbrow, row)

		//fmt.Printf("id: %d, pid: %d, name: %s, body: %s\n", row.Id, row.Pid, row.Name, row.Body)

	}

	return &sliceDbrow
}

func getTarget(pid int) *[]Dbrow {

	var sliceDbrow []Dbrow

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	q := `
        SELECT * FROM memo WHERE id=?
	`
	rows, err := db.Query(q, pid, pid)

	defer rows.Close()
	for rows.Next() {
		//var id, pid, open int
		//var name, body, desired, created string

		row := Dbrow{}

		// カーソルから値を取得
		if err := rows.Scan(
			&row.ID,
			&row.Pid,
			&row.Title,
			&row.Name,
			&row.Body,
			&row.Open,
			&row.Desired,
			&row.Created); err != nil {
			log.Fatal("rows.Scan() ", err)
			return &sliceDbrow
		}
		sliceDbrow = append(sliceDbrow, row)

		//fmt.Printf("id: %d, pid: %d, name: %s, body: %s\n", row.Id, row.Pid, row.Name, row.Body)

	}

	return &sliceDbrow
}

func execDB(db *sql.DB, q string) {
	if _, err := db.Exec(q); err != nil {
		log.Fatal(err)
	}
}

func iexecDB(db *sql.DB, q string) {
	// ignore error
	db.Exec(q)
}

func pSelectdb(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	p := r.Form["p"][0]

	title := "list"
	content := `<div style="width: 100%"><table border="1" style="width: 100%">{__tbody__}</table></div>`
	var tbody []string

	w.Header().Set("Content-Type", "text/html")

	base := `
		<tr style="background-color: #bde9ba;">
		<td colspan="1">番号: {__id__}</td>
		<td colspan="2">件名: {__title__}</td>
		<td colspan="1" rowspan="3"><div><form name="form{__id__}" action="/edit"><input type="hidden" name="p" value="{__id__}"><button>編集</button></form></div></td>
		<td colspan="1" rowspan="3"><div><form name="form{__id__}" action="/confirm_delete"><input type="hidden" name="p" value="{__id__}"><button>削除</button></form></div></td>
		</tr>
		<tr>
		<td>日付: {__created__}</td>
		<td>名前: {__name__}</td>
		<td>期限: {__desired__}</td>
		</tr>
		<tr>
		<td colspan="3">{__body__}</td>
		</tr>
		`
	pid, err := strconv.Atoi(p)
	if err != nil {
		fmt.Println("parse error")
	}
	var lastTitle string
	rep := regexp.MustCompile("\n")
	pDbrow := withChildList(pid)
	for i, data := range *pDbrow {
		state := "open"
		if data.Open == 0 {
			state = "closed"
		}
		fmt.Println("row: ", i)
		ht := strings.Replace(base, "{__id__}", html.EscapeString(strconv.Itoa(data.ID)), -1)
		ht = strings.Replace(ht, "{__title__}", html.EscapeString(data.Title), 1)
		ht = strings.Replace(ht, "{__name__}", html.EscapeString(data.Name), 1)
		ht = strings.Replace(ht, "{__created__}", html.EscapeString(data.Created), 1)
		ht = strings.Replace(ht, "{__desired__}", html.EscapeString(data.Desired), 1)
		ht = strings.Replace(ht, "{__open__}", html.EscapeString(state), 1)
		ht = strings.Replace(ht, "{__body__}", rep.ReplaceAllString(html.EscapeString(data.Body), "<br>"), 1)
		tbody = append(tbody, ht)
		lastTitle = data.Title
	}

	content = strings.Replace(content, "{__tbody__}", strings.Join(tbody, "\n"), 1)

	// output to browser
	fmt.Fprintf(w, outputHTML(title, content, outputFormWithChild(p, lastTitle)))
}

func dumpdb(w http.ResponseWriter, r *http.Request) {
	title := "list"
	content := `<div style="width: 100%"><h2>ダンプ情報</h2><textarea rows="16" cols="64">{__dumpdata__}</textarea></div>`

	w.Header().Set("Content-Type", "text/html")

	jdata := Dbrows{}
	pDbrow := allList()
	for i, data := range *pDbrow {
		fmt.Println("row: ", i)
		fmt.Printf("row: %v\n", jdata.list)
		jdata.list = append(jdata.list, data)
	}

	jsonBytes, err := json.Marshal(jdata.list)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", 302)
	}
	jsonStr := string(jsonBytes)
	fmt.Printf("jsonstr: %v\n", jsonStr)
	content = strings.Replace(content, "{__dumpdata__}", jsonStr, 1)

	// output to browser
	fmt.Fprintf(w, outputHTML(title, content, ""))
}

func confirmRestore(w http.ResponseWriter, r *http.Request) {
	title := "list"
	content := `<div style="width: 100%"><h2>データベースリストア</h2><br>
	You must doing initialize database before restore the database.<br>
	<form method="post" action="/restore">
	<textarea name="dbdata" rows="16" cols="64"></textarea><br>
	<button>restore</button>
	</form>
	</div>`

	w.Header().Set("Content-Type", "text/html")

	// output to browser
	fmt.Fprintf(w, outputHTML(title, content, ""))
}

func restoredb(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	jsonDbdata := r.Form["dbdata"][0]

	title := "list"
	content := `<div style="width: 100%"><h2>ダンプ情報</h2><textarea rows="16" cols="64">{__dumpdata__}</textarea></div>`

	w.Header().Set("Content-Type", "text/html")

	jdata := Dbrows{}

	json.Unmarshal([]byte(jsonDbdata), &jdata.list)
	fmt.Printf("jsonstr: %v\n", jdata.list)

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	q := `
		INSERT INTO memo (parentid, title, name, body, open, desired_at, created_at) VALUES(?, ?, ?, ?, ?,
			DATE(?,'localtime'),
			DATETIME(?,'localtime'))
	`
	for i, data := range jdata.list {
		//pid, _ := strconv.Atoi(r.Form["parentid"][0])
		fmt.Println("restore row: ", i)
		result, err := db.Exec(q, data.Pid, data.Title, data.Name, data.Body, data.Open, data.Desired, data.Created)
		fmt.Printf("restore %v\n", result)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer db.Close()

	// output to browser
	fmt.Fprintf(w, outputHTML(title, content, "finish."))
}

func selectdb(w http.ResponseWriter, r *http.Request) {
	title := "list"
	content := `<div style="width: 100%"><table border="1" style="width: 100%">
	<tr>
	<th>番号</th>
	<th>件名</th>
	<th>更新日付</th>
	<th>作成者</th>
	<th>期限</th>
	<th>状態</th>
	<th>最終更新者</th>
	<th>QAを開く</th>
	<th>QAを閉じる</th>
	<th>QAを削除する</th>
	</tr>
	{__tbody__}
	</table></div>`
	var tbody []string

	w.Header().Set("Content-Type", "text/html")

	/*
		base := `
		<tr>
		<td colspan="3">番号: {__id__}</td>
		</tr>
		<tr>
		<td>日付: {__created__}</td>
		<td>名前: {__name__}</td>
		<td>期限: {__desired__}</td>
		</tr>
		<tr>
		<td>{__body__}</td>
		</tr>
		`
	*/

	base := `
	<tr>
	<td style="background-color: #bde9ba;">番号: {__id__}</td>
	<td>{__title__}</td>
	<td>{__created__}</td>
	<td>{__name__}</td>
	<td>{__desired__}</td>
	<td>{__open__}</td>
	<td>{__latest__}</td>
	<td><button id="btn{__id__}" onclick="location.href='/parent?p={__id__}'">Open</button></td>
	<td><button id="tgls{__id__}" onclick="location.href='/state_change?p={__id__}&amp;state=0'">Close</button></td>
	<td><button id="del{__id__}" onclick="location.href='/confirm_delete?p={__id__}'">Delete</button></td>
	</tr>
	`

	rep := regexp.MustCompile("\n")
	pDbrow := normalList()
	for i, data := range *pDbrow {
		state := "open"
		if data.Open == 0 {
			state = "closed"
		}
		fmt.Println("row: ", i)
		ht := strings.Replace(base, "{__id__}", html.EscapeString(strconv.Itoa(data.ID)), -1)
		ht = strings.Replace(ht, "{__title__}", html.EscapeString(data.Title), 1)
		ht = strings.Replace(ht, "{__name__}", html.EscapeString(data.Name), 1)
		ht = strings.Replace(ht, "{__created__}", html.EscapeString(data.Created), 1)
		ht = strings.Replace(ht, "{__desired__}", html.EscapeString(data.Desired), 1)
		ht = strings.Replace(ht, "{__open__}", html.EscapeString(state), 1)
		ht = strings.Replace(ht, "{__latest__}", html.EscapeString(data.Latest), 1)
		ht = strings.Replace(ht, "{__body__}", rep.ReplaceAllString(html.EscapeString(data.Body), "<br>"), 1)
		tbody = append(tbody, ht)
	}

	content = strings.Replace(content, "{__tbody__}", strings.Join(tbody, "\n"), 1)

	// output to browser
	fmt.Fprintf(w, outputHTML(title, content, outputForm()))
}

func cSelectdb(w http.ResponseWriter, r *http.Request) {
	title := "list"
	content := `<div style="width: 100%"><table border="1" style="width: 100%">
	<tr>
	<th>番号</th>
	<th>件名</th>
	<th>更新日付</th>
	<th>作成者</th>
	<th>期限</th>
	<th>状態</th>
	<th>最終更新者</th>
	<th>QAを開く</th>
	<th>QAを再開</th>
	<th>QAを削除する</th>
	</tr>
	{__tbody__}</table></div>`
	var tbody []string

	w.Header().Set("Content-Type", "text/html")

	base := `
	<tr>
	<td style="background-color: #bde9ba;">番号: {__id__}</td>
	<td>{__title__}</td>
	<td>{__created__}</td>
	<td>{__name__}</td>
	<td>{__desired__}</td>
	<td>{__open__}</td>
	<td>{__latest__}</td>
	<td><button id="btn{__id__}" onclick="location.href='/parent?p={__id__}'">Open</button></td>
	<td><button id="tgls{__id__}" onclick="location.href='/state_change?p={__id__}&amp;state=1'">Reopen</button></td>
	<td><button id="del{__id__}" onclick="location.href='/confirm_delete?p={__id__}'">Delete</button></td>
	</tr>
	`

	rep := regexp.MustCompile("\n")
	pDbrow := closedList()
	for i, data := range *pDbrow {
		state := "open"
		if data.Open == 0 {
			state = "closed"
		}
		fmt.Println("row: ", i)
		ht := strings.Replace(base, "{__id__}", html.EscapeString(strconv.Itoa(data.ID)), -1)
		ht = strings.Replace(ht, "{__title__}", html.EscapeString(data.Title), 1)
		ht = strings.Replace(ht, "{__name__}", html.EscapeString(data.Name), 1)
		ht = strings.Replace(ht, "{__created__}", html.EscapeString(data.Created), 1)
		ht = strings.Replace(ht, "{__desired__}", html.EscapeString(data.Desired), 1)
		ht = strings.Replace(ht, "{__open__}", html.EscapeString(state), 1)
		ht = strings.Replace(ht, "{__latest__}", html.EscapeString(data.Latest), 1)
		ht = strings.Replace(ht, "{__body__}", rep.ReplaceAllString(html.EscapeString(data.Body), "<br>"), 1)
		tbody = append(tbody, ht)
	}

	content = strings.Replace(content, "{__tbody__}", strings.Join(tbody, "\n"), 1)

	// output to browser
	fmt.Fprintf(w, outputHTML(title, content, outputForm()))
}

func insertdb(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("title:", r.Form["title"])
	fmt.Println("name:", r.Form["name"])
	fmt.Println("limit:", r.Form["limit"])
	fmt.Println("body:", r.Form["body"])
	fmt.Println("pid:", r.Form["parentid"])
	/*
		fmt.Println("path", r.URL.Path)
		fmt.Println("scheme", r.URL.Scheme)
		fmt.Println(r.Form["url_long"])
	*/
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	sliceDate := strings.Split(r.Form["limit"][0], "/")
	fmt.Printf("slice_date: %v\n", sliceDate)
	limitDate := sliceDate[2] + "-" + sliceDate[0] + "-" + sliceDate[1]

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	q := `
		INSERT INTO memo (parentid, title, name, body, open, desired_at, created_at) VALUES(?, ?, ?, ?, 1,
			DATE(?,'localtime'),
			DATETIME('now','localtime'))
	`
	pid, _ := strconv.Atoi(r.Form["parentid"][0])
	result, err := db.Exec(q, pid, r.Form["title"][0], r.Form["name"][0], r.Form["body"][0], limitDate)
	if err != nil {
		log.Fatal(err)
	}

	/*
		// 親の作成日を更新する処理
		if pid != 0 {
			q := `
				UPDATE memo SET created_at=DATETIME('now','localtime')
				WHERE id=?
			`
			result, err := db.Exec(q, pid)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("update parent create_at: %v\n", result)
		}
	*/

	defer db.Close()

	fmt.Println("insert rows: ", result)

	http.Redirect(w, r, "/", 302)
}

func updatedb(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("title:", r.Form["title"])
	fmt.Println("name:", r.Form["name"])
	fmt.Println("limit:", r.Form["limit"])
	fmt.Println("body:", r.Form["body"])
	fmt.Println("pid:", r.Form["parentid"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	p := r.Form["p"][0]

	sliceDate := strings.Split(r.Form["limit"][0], "/")
	fmt.Printf("slice_date: %v\n", sliceDate)
	limitDate := sliceDate[2] + "-" + sliceDate[0] + "-" + sliceDate[1]

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	q := `
        SELECT * FROM memo WHERE id=? LIMIT 1
	`
	rows, err := db.Query(q, p)

	defer rows.Close()

	//var pid int
	for rows.Next() {
		//var id, pid, open int
		//var name, body, desired, created string

		row := Dbrow{}

		// カーソルから値を取得
		if err := rows.Scan(
			&row.ID,
			&row.Pid,
			&row.Title,
			&row.Name,
			&row.Body,
			&row.Open,
			&row.Desired,
			&row.Created); err != nil {
			log.Fatal("rows.Scan() ", err)
		}
		if err != nil {
			log.Fatal(err)
		}
		//pid = row.Pid
		//fmt.Printf("id: %d, pid: %d, name: %s, body: %s\n", row.Id, row.Pid, row.Name, row.Body)
	}

	q = `
		UPDATE memo SET title=?, name=?, body=?, desired_at=DATE(?,'localtime'), created_at=DATETIME('now','localtime')
		WHERE id=?
	`
	result, err := db.Exec(q, r.Form["title"][0], r.Form["name"][0], r.Form["body"][0], limitDate, p)
	if err != nil {
		log.Fatal(err)
	}

	/*
		// 親の更新日付を更新する処理
		if pid != 0 {
			q = `
			UPDATE memo SET created_at=DATETIME('now','localtime')
			WHERE id=?
			`
			result, err := db.Exec(q, pid)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("update parent create_at: %v\n", result)
		}
	*/

	defer db.Close()

	fmt.Println("update rows: ", result)

	http.Redirect(w, r, "/", 302)
}

func stateChange(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	p := r.Form["p"][0]
	s := r.Form["state"][0]

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	state := 1
	if len(s) > 0 {
		state, err = strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		if state > 0 {
			state = 1
		}
	}
	q := `
		UPDATE memo SET open=?, created_at=DATETIME('now','localtime')
		WHERE id=?
	`
	result, err := db.Exec(q, state, p)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Println("close rows: ", result)

	http.Redirect(w, r, "/", 302)
}

func editdata(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	p := r.Form["p"][0]

	title := "list"
	content := `<div style="width: 100%">
	<form method="post" action="/update">
	<input type="hidden" name="p" value="{__id__}">
	<table border="1" style="wodth: 100%">{__tbody__}</table>
	</form>
	</div>
	`
	var tbody []string

	w.Header().Set("Content-Type", "text/html")

	base := `
		<tr>
		<td>件名: <input type="text" name="title" require value="{__title__}"></td>
		<td rowspan="4"><button>更新</button></td>
		</tr>
		<tr>
		<td>名前: <input type="text" name="name" require value="{__name__}"></td>
		</tr>
		<tr>
		<td>期限: <input type="text" name="limit" id="datepicker" require value="{__desired__}"></td>
		</tr>
		<tr>
		<td><textarea name="body" required cols="64" rows="16">{__body__}</textarea></td>
		</tr>
		`
	pid, err := strconv.Atoi(p)
	if err != nil {
		fmt.Println("parse error")
	}

	pDbrow := getTarget(pid)
	for i, data := range *pDbrow {
		state := "open"
		if data.Open == 0 {
			state = "closed"
		}
		sliceDesired := strings.Split(data.Desired, "-")
		desired := sliceDesired[1] + "/" + sliceDesired[2] + "/" + sliceDesired[0]
		fmt.Println("row: ", i)
		ht := strings.Replace(base, "{__id__}", html.EscapeString(strconv.Itoa(data.ID)), -1)
		ht = strings.Replace(ht, "{__title__}", html.EscapeString(data.Title), 1)
		ht = strings.Replace(ht, "{__name__}", html.EscapeString(data.Name), 1)
		ht = strings.Replace(ht, "{__created__}", html.EscapeString(data.Created), 1)
		ht = strings.Replace(ht, "{__desired__}", html.EscapeString(desired), 1)
		ht = strings.Replace(ht, "{__open__}", html.EscapeString(state), 1)
		ht = strings.Replace(ht, "{__body__}", html.EscapeString(data.Body), 1)
		tbody = append(tbody, ht)
	}

	content = strings.Replace(content, "{__tbody__}", strings.Join(tbody, "\n"), 1)
	content = strings.Replace(content, "{__id__}", p, -1)

	// output to browser
	fmt.Fprintf(w, outputConfirm(title, content))
}

func confirmDelete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	p := r.Form["p"][0]

	title := "list"
	content := `<div style="width: 100%"><h2>親を削除した場合は子も削除されます。</h2><table border="1" style="wodth: 100%">{__tbody__}</table></div>
	<form action="/delete">
	<input type="hidden" name="p" value="{__id__}">
	<button>削除</button>
	</form>
	`
	var tbody []string

	w.Header().Set("Content-Type", "text/html")

	base := `
		<tr>
		<td colspan="1">番号: {__id__}</td>
		<td colspan="2">件名: {__title__}</td>
		</tr>
		<tr>
		<td>日付: {__created__}</td>
		<td>名前: {__name__}</td>
		<td>期限: {__desired__}</td>
		</tr>
		<tr>
		<td colspan="3">{__body__}</td>
		</tr>
		`
	pid, err := strconv.Atoi(p)
	if err != nil {
		fmt.Println("parse error")
	}
	rep := regexp.MustCompile("\n")
	pDbrow := getTarget(pid)
	for i, data := range *pDbrow {
		state := "open"
		if data.Open == 0 {
			state = "closed"
		}
		fmt.Println("row: ", i)
		ht := strings.Replace(base, "{__id__}", html.EscapeString(strconv.Itoa(data.ID)), 3)
		ht = strings.Replace(ht, "{__title__}", html.EscapeString(data.Title), 1)
		ht = strings.Replace(ht, "{__name__}", html.EscapeString(data.Name), 1)
		ht = strings.Replace(ht, "{__created__}", html.EscapeString(data.Created), 1)
		ht = strings.Replace(ht, "{__desired__}", html.EscapeString(data.Desired), 1)
		ht = strings.Replace(ht, "{__open__}", html.EscapeString(state), 1)
		ht = strings.Replace(ht, "{__body__}", rep.ReplaceAllString(html.EscapeString(data.Body), "<br>"), 1)
		tbody = append(tbody, ht)
	}

	content = strings.Replace(content, "{__tbody__}", strings.Join(tbody, "\n"), 1)
	content = strings.Replace(content, "{__id__}", p, -1)

	// output to browser
	fmt.Fprintf(w, outputConfirm(title, content))
}

func deletedb(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	p := r.Form["p"][0]
	fmt.Printf("delete id: %v\n", p)
	re := regexp.MustCompile(`\d`)
	if !re.MatchString(p) {
		fmt.Printf("invalid delete id: %v\n", p)
		http.Redirect(w, r, "/", 302)
	}

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	q := `
		DELETE FROM memo WHERE id=? OR parentid=?
	`
	id, _ := strconv.Atoi(p)
	result, err := db.Exec(q, id, id)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Println("delete rows: ", result)

	http.Redirect(w, r, "/", 302)
}

func outputHTML(t string, q string, f string) string {
	base := `
	<doctype html>
	<html>
	<header>
	<meta charset="utf-8">
	<title>{__title__}</title>
	<link rel="stylesheet" href="//code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css">
	<script src="https://code.jquery.com/jquery-1.12.4.js"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.js"></script>
	<script>
	$( function() {
	  $( "#datepicker" ).datepicker();
	} );
	</script>	</header>
	<body>
	<a href="/">Issue</a>&nbsp;
	<a href="/closed">Closed</a>&nbsp;
	<a href="/dump">Dump</a>&nbsp;
	<a href="/confirm_restore">Restore</a>&nbsp;
	<hr>
	{__content__}
	<div id="form1">
	{__form__}
	</div>
	</body>
	</html>
	`
	ht := strings.Replace(base, "{__title__}", t, 1)
	ht = strings.Replace(ht, "{__content__}", q, 1)
	ht = strings.Replace(ht, "{__form__}", f, 1)
	return ht
}

func outputConfirm(t string, q string) string {
	base := `
	<doctype html>
	<html>
	<header>
	<meta charset="utf-8">
	<title>{__title__}</title>
	<link rel="stylesheet" href="//code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css">
	<link rel="stylesheet" href="/resources/demos/style.css">
	<script src="https://code.jquery.com/jquery-1.12.4.js"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.js"></script>
	<script>
	$( function() {
	  $( "#datepicker" ).datepicker();
	} );
	</script>	</header>
	<body>
	<a href="/">Issue</a>&nbsp;
	<a href="/closed">Closed</a>&nbsp;
	<hr>
	{__content__}
	</body>
	</html>
	`
	ht := strings.Replace(base, "{__title__}", t, 1)
	ht = strings.Replace(ht, "{__content__}", q, 1)
	return ht
}

func outputForm() string {
	base := `
	<script>
	</script>
	<hr>
	<form method="post" action="/create">
	件名: <input type="text" id="title" name="title" required><br>
	名前: <input type="name" id="name" name="name" required><br>
	期限: <input type="text" id="datepicker" name="limit" autocomplete="off" placeholder="08/31/2019" required><br>
	本文: <textarea id="body" name="body" required rows="16" cols="64"></textarea><br>
	<input type="hidden" name="parentid" value="{__parentid__}">
	<button onsubmit="return formCheck();">作成</button>
	</form>
	`
	ht := base
	return ht
}

func outputFormWithChild(pid string, title string) string {
	base := `
	<script>
	</script>
	<hr>
	<form method="post" action="/create">
	件名: <input type="text" id="title" name="title" required value="{__title__}"><br>
	名前: <input type="name" id="name" name="name" required><br>
	期限: <input type="text" id="datepicker" name="limit" autocomplete="off" placeholder="08/31/2019" required><br>
	本文: <textarea id="body" name="body" required rows="16" cols="64"></textarea><br>
	<input type="hidden" name="parentid" value="{__parentid__}">
	<button onsubmit="return formCheck();">作成</button>
	</form>
	`
	ht := strings.Replace(base, "{__parentid__}", html.EscapeString(pid), 1)
	ht = strings.Replace(ht, "{__title__}", html.EscapeString(title), 1)

	return ht
}

func router(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		selectdb(w, r)
	} else if r.URL.Path == "/closed" {
		cSelectdb(w, r)
	} else if r.URL.Path == "/parent" {
		pSelectdb(w, r)
	} else if r.URL.Path == "/create" {
		insertdb(w, r)
	} else if r.URL.Path == "/state_change" {
		stateChange(w, r)
	} else if r.URL.Path == "/confirm_delete" {
		confirmDelete(w, r)
	} else if r.URL.Path == "/delete" {
		deletedb(w, r)
	} else if r.URL.Path == "/edit" {
		editdata(w, r)
	} else if r.URL.Path == "/update" {
		updatedb(w, r)
	} else if r.URL.Path == "/dump" {
		dumpdb(w, r)
	} else if r.URL.Path == "/confirm_restore" {
		confirmRestore(w, r)
	} else if r.URL.Path == "/restore" {
		restoredb(w, r)
	} else if r.URL.Path == "/init" {
		initdb(w, r)
	} else if r.URL.Path == "/testinit" {
		testInitdb(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	fmt.Println(os.Args)
	localport := "127.0.0.1:9000"
	for i := 0; i < len(os.Args); i++ {
		if strings.Compare(os.Args[i], "-l") == 0 && i+1 <= len(os.Args) {
			fmt.Println("len os.Args", len(os.Args))
			if m, _ := regexp.MatchString("^([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3})?:[0-9]{1,5}$", os.Args[i+1]); m {
				localport = os.Args[i+1]
				i++
			} else {
				log.Fatal("Invalid argument")
			}
		}
	}

	fmt.Println("open your browser ", localport)

	http.HandleFunc("/", router)
	err := http.ListenAndServe(localport, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initdb(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	q := `
        DROP TABLE memo
	`
	iexecDB(db, q) // エラーは無視する

	q = `
        CREATE TABLE memo (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		parentid INTEGER NOT NULL,
		title VARCHAR(2048) NOT NULL,
		name VARCHAR(1024) NOT NULL,
		body VARCHAR(9192) NOT NULL,
		open INTEGER NOT NULL DEFAULT 1,
		desired_at IMESTAMTP DEFAULT (DATE('now','localtime')),
        created_at TIMESTAMP DEFAULT (DATETIME('now','localtime'))
        )
    `
	execDB(db, q)

	db.Close()

	http.Redirect(w, r, "/", 302)

}

func testInitdb(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	q := `
        DROP TABLE memo
	`
	iexecDB(db, q) // エラーは無視する

	q = `
        CREATE TABLE memo (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		parentid INTEGER NOT NULL,
		title VARCHAR(2048) NOT NULL,
		name VARCHAR(1024) NOT NULL,
		body VARCHAR(9192) NOT NULL,
		open INTEGER NOT NULL DEFAULT 1,
		desired_at IMESTAMTP DEFAULT (DATE('now','localtime')),
        created_at TIMESTAMP DEFAULT (DATETIME('now','localtime'))
        )
    `
	execDB(db, q)

	q = `
        INSERT INTO memo (
		parentid,
		title,
		name,
		body,
		open,
		desired_at
		)
		VALUES(0, "テストのタイトルです", "hoge fuga", "テストの質問です", 1, "2019-08-31")
    `
	execDB(db, q)

	q = `
        INSERT INTO memo (
		parentid,
		title,
		name,
		body,
		open,
		desired_at
		)
		VALUES(1, "test title", "foo bar", "１の関連データです", 1, "2019-08-31")
    `
	execDB(db, q)

	q = `
        INSERT INTO memo (
		parentid,
		title,
		name,
		body,
		open,
		desired_at
		)
		VALUES(0, "2つ目のtest titleです", "foo bar", "テストの質問２です", 1, "2019-09-31")
    `
	execDB(db, q)

	db.Close()

	/*
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("location", "/")
		w.WriteHeader(http.StatusMovedPermanently) // 301 Moved Permanently
	*/
	http.Redirect(w, r, "/", 302)

}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //オプションを解析します。デフォルトでは解析しません。
	fmt.Println(r.Form) //このデータはサーバのプリント情報に出力されます。
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //ここでwに入るものがクライアントに出力されます。
}
