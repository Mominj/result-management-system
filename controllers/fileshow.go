package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type Quer struct {
	Year      string `db:"year"`
	Exam_type string `db:"exam_type"`
	Group     string `db:"groups"`
	Roll      string `db:"roll_number"`
}
type DD struct {
	Name string
	Gpa  float32
}
type D struct {
	ID                                   int32         `db:"id"`
	Year                                 string        `db:"year"`
	Exam_type                            string        `db:"exam_type"`
	Group                                string        `db:"groups"`
	Roll_number                          string        `db:"roll_number"`
	Name                                 string        `db:"name"`
	Bangla                               sql.NullInt32 `db:"bangla"`
	English                              sql.NullInt32 `db:"english"`
	Ict                                  sql.NullInt32 `db:"ict"`
	Physics                              sql.NullInt32 `db:"physics"`
	Chemistry                            sql.NullInt32 `db:"chemistry"`
	Biology                              sql.NullInt32 `db:"biology"`
	Higher_mathematics                   sql.NullInt32 `db:"higher_mathematics"`
	Agriculture_education                sql.NullInt32 `db:"agriculture_education"`
	Geography                            sql.NullInt32 `db:"geography"`
	Psychology                           sql.NullInt32 `db:"psychology"`
	Accounting                           sql.NullInt32 `db:"accounting"`
	Statistics                           sql.NullInt32 `db:"statistics"`
	Economics                            sql.NullInt32 `db:"economics"`
	Business_organization_and_management sql.NullInt32 `db:"business_organization_and_management"`
	Finance_banking_and_insurance        sql.NullInt32 `db:"finance_banking_and_insurance"`
	History                              sql.NullInt32 `db:"history"`
	Islamic_history_and_culture          sql.NullInt32 `db:"islamic_history_and_culture"`
	Sociology                            sql.NullInt32 `db:"sociology"`
	Logic                                sql.NullInt32 `db:"logic"`
}

func cal(x int32) string {
	var s string
	switch {
	case x <= 100 && x >= 80:
		s = "A+"
	case x <= 79 && x >= 70:
		s = "A"
	case x <= 69 && x >= 60:
		s = "A-"
	case x <= 59 && x >= 50:
		s = "B"
	case x <= 49 && x >= 40:
		s = "C"
	case x <= 39 && x >= 33:
		s = "D"
	default:
		s = "F"
	}
	return s
}

func gpaa(i int32) float32 {
	var s float32
	if i >= 80 && i <= 100 {
		s = 5
	}
	if i >= 70 && i <= 79 {
		s = 4
	}
	if i >= 60 && i <= 69 {
		s = 3.5
	}
	if i >= 50 && i <= 59 {
		s = 3
	}
	if i >= 40 && i <= 49 {
		s = 2
	}
	if i >= 33 && i <= 39 {
		s = 1
	}
	fmt.Println("ii...", i, "....", s)
	return s
}

func FileShow(w http.ResponseWriter, req *http.Request) {
	var quer Quer
	results := []D{}
	if err := json.NewDecoder(req.Body).Decode(&quer); err != nil {
		log.Println("error occur while  req body data  : ", err.Error())
		return
	}
	query := `SELECT * FROM result WHERE roll_number = $1 AND exam_type = $2 AND groups = $3 AND year=$4`

	if err := DB.Select(&results, query, quer.Roll, quer.Exam_type, quer.Group, quer.Year); err != nil {
		if err != nil {
			log.Println("error while getting RESULT from db: ", err.Error())
			return
		}
	}
	subjects := make(map[string]int32)
	resultss := make(map[string]string)

	if (results[0].Bangla.Int32) != 0 {
		subjects["Bangla"] = int32(results[0].Bangla.Int32)
	}
	if (results[0].English.Int32) != 0 {
		subjects["English"] = int32(results[0].English.Int32)
	}
	if (results[0].Ict.Int32) != 0 {
		subjects["Ict"] = int32(results[0].Ict.Int32)
	}
	if (results[0].Physics.Int32) != 0 {
		subjects["Physics"] = int32(results[0].Physics.Int32)
	}
	if (results[0].Chemistry.Int32) != 0 {
		subjects["Chemistry"] = int32(results[0].Chemistry.Int32)
	}
	if (results[0].Biology.Int32) != 0 {
		subjects["Biology"] = int32(results[0].Biology.Int32)
	}
	if (results[0].Higher_mathematics.Int32) != 0 {
		subjects["Higher_mathematics"] = int32(results[0].Higher_mathematics.Int32)
	}
	//fmt.Println(quer)
	for i, v := range subjects {
		resultss[i] = cal(v)
	}

	var sum float32
	for _, v := range subjects {
		if v < 33 {
			sum = 0
			return
		}
		sum = sum + gpaa(v)
	}
	//gpa := sum / 7

	rjson, err := json.Marshal(resultss)
	if err != nil {
		fmt.Println("err occur when convert", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write((rjson))

	/*
	   pdf
	   create
	*/
	var data [][]string
	for i, v := range resultss {
		dat := []string{}
		dat = append(dat, i)
		dat = append(dat, v)
		data = append(data, dat)
	}
	fmt.Println(data)
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 10, 20)

	buildHeading(m)
	buildFruitList(m, data)
	if err := m.OutputFileAndClose("pdfs/123.pdf"); err != nil {
		fmt.Println("could not save pdf", err.Error())
		return
	}

	fmt.Println("SAVE pdf")
}

func buildHeading(m pdf.Maroto) {
	m.RegisterHeader(func() {
		m.Row(50, func() {
			m.Col(12, func() {
				err := m.FileImage("images/images.png", props.Rect{
					Center:  true,
					Percent: 75,
				})

				if err != nil {
					fmt.Println("image file was not loaded ", err.Error())

				}
			})
		})
	})
}

func buildFruitList(m pdf.Maroto, fruits [][]string) {
	tableHeadings := []string{"Subject", "GPA"}
	contents := fruits
	lightPurpleColor := getLightPurpleColor()

	m.SetBackgroundColor(getTealColor())
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("HSC Result 2021", props.Text{
				Top:    2,
				Size:   13,
				Color:  color.NewWhite(),
				Family: consts.Courier,
				Style:  consts.Bold,
				Align:  consts.Center,
			})
		})
	})

	m.SetBackgroundColor(color.NewWhite())

	m.TableList(tableHeadings, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      13,
			GridSizes: []uint{3, 7, 2},
		},
		ContentProp: props.TableListContent{
			Size:      13,
			GridSizes: []uint{3, 7, 2},
		},
		Align:                consts.Center,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})

}

func getDarkPurpleColor() color.Color {
	return color.Color{
		Red:   88,
		Green: 80,
		Blue:  99,
	}
}
func getLightPurpleColor() color.Color {
	return color.Color{
		Red:   210,
		Green: 200,
		Blue:  230,
	}
}

func getTealColor() color.Color {
	return color.Color{
		Red:   3,
		Green: 166,
		Blue:  166,
	}
}
