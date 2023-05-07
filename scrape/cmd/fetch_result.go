package main

import (
	resultNITH "Result-NITH"
	"Result-NITH/db"
	"context"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// ParseResultHtml depends on the current html structure of official result website
// - table 0 => last update title table
// - table 1 => row 0 => rollNo, name, fathersName
// - tables from 2 to end, not including last one have semester data
// - two table for each semester
// - first table subjects
// - second table summary
// - last table useless
func ParseResultHtml(body io.ReadCloser) (user *resultNITH.StudentHtmlParsed, lastUpdateResultName string, err error) {
	resultDoc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, "", err
	}
	user = &resultNITH.StudentHtmlParsed{}
	tableFind := resultDoc.Find("table")
	semesters := (tableFind.Length() - 3) / 2
	if semesters < 0 || semesters >= tableFind.Length() {
		return nil, "", fmt.Errorf("something went wrong")
	}
	user.SemesterResults = make([]resultNITH.SemesterResult, semesters)
	tableFind.Each(func(tableIndex int, selection *goquery.Selection) {
		if tableIndex == 0 {
			//last update title table
			lastUpdateResultName = strings.TrimSpace(selection.Find("tr td").Text())
		} else if tableIndex == 1 {
			//student roll number, name, father's name
			selection.Find("td").Each(func(cellIndex int, selection *goquery.Selection) {
				txt := strings.Replace(selection.Text(), "ROLL NUMBER", "", -1)
				txt = strings.Replace(txt, "STUDENT NAME", "", -1)
				txt = strings.Replace(txt, "FATHER NAME", "", -1)
				txt = strings.TrimSpace(txt)
				switch cellIndex {
				case 0:
					user.RollNumber = txt
				case 1:
					user.Name = txt
				case 2:
					user.FathersName = txt
				}
			})
		} else if tableIndex == tableFind.Length()-1 {
			//useless table
			return
		} else if tableIndex%2 == 0 {
			//semester result table: subjects data
			rowFind := selection.Find("tr")
			subjectsResult := make([]resultNITH.SubjectResult, rowFind.Length()-2)
			rowFind.Each(func(rowIndex int, selection *goquery.Selection) {
				if rowIndex < 2 {
					return
				}
				//each row is a subject after row index 1
				selection.Find("td").Each(func(cellIndex int, selection *goquery.Selection) {
					text := strings.TrimSpace(selection.Text())
					switch cellIndex {
					case 1:
						subjectsResult[rowIndex-2].SubjectName = text
					case 2:
						subjectsResult[rowIndex-2].SubjectCode = text
					case 3:
						{
							subPoint, _ := strconv.Atoi(text)
							subjectsResult[rowIndex-2].SubPoint = int64(subPoint)
						}
					case 4:
						subjectsResult[rowIndex-2].Grade = text
					case 5:
						{
							subGP, _ := strconv.Atoi(text)
							subjectsResult[rowIndex-2].SubGP = int64(subGP)
						}
					}
				})
			})
			user.SemesterResults[(tableIndex-2)/2].SubjectResults = subjectsResult
			user.SemesterResults[(tableIndex-2)/2].SemesterNumber = int64((tableIndex-2)/2 + 1)
		} else {
			//semester result table: semester overall data
			selection.Find("tr td").Each(func(cellIndex int, selection *goquery.Selection) {
				equalCharPosition := strings.Index(selection.Text(), "=")
				text := strings.TrimSpace(selection.Text()[equalCharPosition+1:])
				if cellIndex == 1 {
					user.SemesterResults[(tableIndex-2)/2].SGPI, _ = strconv.ParseFloat(text, 64)
				} else if cellIndex == 3 {
					user.SemesterResults[(tableIndex-2)/2].CGPI, _ = strconv.ParseFloat(text, 64)
				}
			})
		}
	})
	user.CGPI = user.SemesterResults[len(user.SemesterResults)-1].CGPI
	return
}

// StoreStudentInDb this function will contain the logic to associate roll number with branch and batch year
func StoreStudentInDb(sqlDB *sql.DB, queries *db.Queries, students []resultNITH.StudentHtmlParsed) error {
	ctx := context.Background()

	tx, err := sqlDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)
	for studentIndex, student := range students {
		//student table
		batch, branch := resultNITH.GetBatchAndBranch(student.RollNumber)
		_, err := qtx.CreateStudent(ctx, db.CreateStudentParams{
			RollNumber:     student.RollNumber,
			Name:           student.Name,
			FathersName:    student.FathersName,
			Batch:          batch,
			Branch:         branch,
			LatestSemester: int64(len(student.SemesterResults)),
			Cgpi:           student.CGPI,
		})
		if err != nil {
			println("Error when creating student entry for rollNumber ", student.RollNumber, ": ", err)
		}
		//subject table
		for _, sem := range student.SemesterResults {
			for _, subject := range sem.SubjectResults {
				_, err := qtx.CreateSubject(ctx, db.CreateSubjectParams{
					Code:    subject.SubjectCode,
					Name:    subject.SubjectName,
					Credits: subject.SubPoint,
				})
				if err != nil {
					println("Error when creating subject entry for rollNumber ", student.RollNumber, "and subject ", subject.SubjectName, ": ", err)
				}
			}
		}
		//semester result
		for semNumber, sem := range student.SemesterResults {
			_, err := qtx.CreateSemesterResultData(ctx, db.CreateSemesterResultDataParams{
				StudentRollNumber: student.RollNumber,
				Semester:          int64(semNumber + 1),
				Cgpi:              sem.CGPI,
				Sgpi:              sem.SGPI,
			})
			if err != nil {
				println("Error when creating semester entry for rollNumber ", student.RollNumber, "and semNumber ", semNumber, ": ", err)
			}
		}
		//subject result data
		for semNumber, sem := range student.SemesterResults {
			for _, subject := range sem.SubjectResults {
				_, err := qtx.CreateSubjectResultData(ctx, db.CreateSubjectResultDataParams{
					StudentRollNumber: student.RollNumber,
					SubjectCode:       subject.SubjectCode,
					Grade:             subject.Grade,
					SubGp:             subject.SubGP,
					Semester:          int64(semNumber + 1),
				})
				if err != nil {
					println("Error when creating subject result entry for rollNumber ", student.RollNumber, "and subject ", subject.SubjectName, ": ", err.Error())
				}
			}
		}
		println("Finished student Number ", studentIndex, "out of ", len(students))
	}
	return tx.Commit()
}

func getResultHtml(rollNumber string) (io.ReadCloser, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Jar: cookieJar,
	}
	path := resultNITH.GetUrlForRollNumber(rollNumber)

	//get tokens
	formPageResponse, err := httpClient.Get(path)
	if err != nil {
		return nil, err
	}
	formPageDoc, err := goquery.NewDocumentFromReader(formPageResponse.Body)
	if err != nil {
		return nil, err
	}
	csrfToken, exists := formPageDoc.Find("[name=CSRFToken]").Attr("value")
	if !exists {
		return nil, fmt.Errorf("CSRFToken not found")
	}
	verToken, exists := formPageDoc.Find("[name=RequestVerificationToken]").Attr("value")
	if !exists {
		return nil, fmt.Errorf("RequestVerificationToken not found")
	}
	//get result html
	data := url.Values{
		"RollNumber":               {rollNumber},
		"CSRFToken":                {csrfToken},
		"RequestVerificationToken": {verToken},
		"B1":                       {"Submit"},
	}
	postReq, err := http.NewRequest(http.MethodPost, path, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	postReq.Header.Set("DNT", "1")
	postReq.Header.Set("Content-Type", " application/x-www-form-urlencoded")
	postReq.AddCookie(formPageResponse.Cookies()[0])
	resp, err := httpClient.Do(postReq)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func getResultsFromWeb() []resultNITH.StudentHtmlParsed {
	//build an array of roll numbers
	rollNumbers := resultNITH.GenRollNumbers()
	println("Total roll numbers to process: ", len(rollNumbers))
	var doneRollNumbers int32 = 0
	//build an array of student objects that contain result
	var students []resultNITH.StudentHtmlParsed
	for _, rollNumber := range rollNumbers {
		time.Sleep(time.Second * time.Duration(rand.Intn(3)))
		resultHtml, err := getResultHtml(rollNumber)
		if err != nil {
			atomic.AddInt32(&doneRollNumbers, 1)
			err = fmt.Errorf("error for rollNumber %s: %w, Total done: %d\n", rollNumber, err, doneRollNumbers)
			log.Print(err)
		}
		student, _, err := ParseResultHtml(resultHtml)
		if err == nil {
			atomic.AddInt32(&doneRollNumbers, 1)
			fmt.Printf("Success for rollNumber %s, Total done: %d\n", rollNumber, doneRollNumbers)
			students = append(students, *student)
		} else {
			atomic.AddInt32(&doneRollNumbers, 1)
			err = fmt.Errorf("error for rollNumber %s: %w, Total done: %d\n", rollNumber, err, doneRollNumbers)
			log.Print(err)
		}
	}
	return students
}

func main() {
	students := getResultsFromWeb()
	println("\n\nFinished fetching students\n")

	db, queries := resultNITH.GetDbQueriesForNewDb("result.db")

	err := StoreStudentInDb(db, queries, students)
	if err != nil {
		println("Error in storeTheStudentInDb")
		println(err)
	}
}
