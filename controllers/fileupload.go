package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/lib/pq"
)

type Result struct {
	ID                                   int32  `db:"id"`
	Year                                 string `db:"year"`
	Exam_type                            string `db:"exam_type"`
	Group                                string `db:"groups"`
	Roll_number                          string `db:"roll_number"`
	Name                                 string `db:"name"`
	Bangla                               int32  `db:"bangla"`
	English                              int32  `db:"english"`
	Ict                                  int32  `db:"ict"`
	Physics                              int32  `db:"physics"`
	Chemistry                            int32  `db:"chemistry"`
	Biology                              int32  `db:"biology"`
	Higher_mathematics                   int32  `db:"higher_mathematics"`
	Agriculture_education                int32  `db:"agriculture_education"`
	Geography                            int32  `db:"geography"`
	Psychology                           int32  `db:"psychology"`
	Accounting                           int32  `db:"accounting"`
	Statistics                           int32  `db:"statistics"`
	Economics                            int32  `db:"economics"`
	Business_organization_and_management int32  `db:"business_organization_and_management"`
	Finance_banking_and_insurance        int32  `db:"finance_banking_and_insurance"`
	History                              int32  `db:"history"`
	Islamic_history_and_culture          int32  `db:"islamic_history_and_culture"`
	Sociology                            int32  `db:"sociology"`
	Logic                                int32  `db:"logic"`
}

func FileSave(w http.ResponseWriter, req *http.Request) {
	//var result Result
	req.ParseMultipartForm(10 << 20)

	if err := req.ParseForm(); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	file, handler, err := req.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("uploaded file %+v\n", handler.Filename)
	fmt.Printf("MIME header %+v\n", handler.Header)

	f, err := os.OpenFile("./public/upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	fi, err := excelize.OpenFile("./public/upload/" + handler.Filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	rows := fi.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	dat := make(map[string]interface{})
	examtype := strings.Join(req.Form["examtype"], " ")
	year := strings.Join(req.Form["year"], " ")
	group := strings.Join(req.Form["group"], " ")
	for i := 1; i < len(rows); i++ {
		for j := 0; j < len(rows[i]); j++ {
			dat[rows[0][j]] = rows[i][j]
		}
		query := `
		INSERT INTO result(
			year,
			exam_type,
			groups,
			roll_number,
			name,
			bangla,
			english,
			ict,
			physics,
			chemistry,
			biology,  
			higher_mathematics 

		) VALUES(
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
		)`
		_, err := DB.Exec(query, year, examtype, group, dat["Roll"], dat["Name"], dat["bangla"], dat["english"], dat["ict"], dat["physics"], dat["chemistry"], dat["biology"], dat["math"])
		if err != nil {
			fmt.Println("error occur when insert data", err)
			return
		}

	}

}
