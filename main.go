package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	router()
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func router() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/add-word", addWordPage)
	http.HandleFunc("/make-quiz", makeQuizPage)
	http.HandleFunc("/quiz", quizPage)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	}
}

func addWordPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/add-word.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		word := strings.TrimSpace(r.FormValue("word"))
		word = strings.Replace(word, string(","), string(";"), -1)
		meaning := strings.TrimSpace(r.FormValue("meaning"))
		meaning = strings.Replace(meaning, string(","), string(";"), -1)
		qaMapping := make(map[string]string)
		qaMapping["Question"] = meaning
		qaMapping["Answer"] = word
		qaMapping["TrueAnswerCount"] = "0"
		qaMapping["WrongAnswerCount"] = "0"
		qaMapping["AskCount"] = "0"
		saveData(qaMapping)
		tmpl := template.Must(template.ParseFiles("templates/add-word.html"))
		data := map[string]string{"Data": "Your word is added!"}
		tmpl.Execute(w, data)
	}
}

func makeQuizPage(w http.ResponseWriter, r *http.Request) {
	recordCount := len(getAllData())
	if r.Method == http.MethodGet {

		data := map[string]string{"RecordCount": strconv.Itoa(recordCount)}
		tmpl := template.Must(template.ParseFiles("templates/make-quiz.html"))
		tmpl.Execute(w, data)

	} else if r.Method == http.MethodPost {

		if recordCount < 4 {
			fmt.Fprintf(w, "You need to add more than or equal to 4 words to make a quiz.")
		}

		questionCountToAnswer := strings.TrimSpace(r.FormValue("questionCount"))
		urlToRedirect := "/quiz?questionCount=" + questionCountToAnswer
		http.Redirect(w, r, urlToRedirect, http.StatusFound)

	}
}

func quizPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		questionCount := r.URL.Query().Get("questionCount")

		allData := getAllData()
		prioritazedData := getPrioritazedData(allData)

		questionCountToAnswer, err := strconv.Atoi(questionCount)
		if err != nil {
			fmt.Fprintf(w, "Question count must be positive integer")
			return
		}

		questionsAndAnswerOptionsMappingArr := getQuizData(questionCountToAnswer, prioritazedData, allData)

		dataArr := make([]map[string]string, questionCountToAnswer)

		// initializing maps
		for j := 0; j < questionCountToAnswer; j++ {
			dataArr[j] = make(map[string]string)
		}

		for i := 0; i < questionCountToAnswer; i++ {
			dataArr[i] = map[string]string{
				"Index":      strconv.Itoa(i + 1),
				"Question":   questionsAndAnswerOptionsMappingArr[i]["Question"],
				"OptionA":    questionsAndAnswerOptionsMappingArr[i]["OptionA"],
				"OptionB":    questionsAndAnswerOptionsMappingArr[i]["OptionB"],
				"OptionC":    questionsAndAnswerOptionsMappingArr[i]["OptionC"],
				"OptionD":    questionsAndAnswerOptionsMappingArr[i]["OptionD"],
				"TrueAnswer": questionsAndAnswerOptionsMappingArr[i]["TrueAnswer"],
			}
		}

		// "QuestionCount": strconv.Itoa(questionCountToAnswer)
		data := map[string]interface{}{"Data": dataArr, "QuestionCount": strconv.Itoa(questionCountToAnswer)}
		tmpl := template.Must(template.ParseFiles("templates/quiz.html"))
		tmpl.Execute(w, data)
	} else if r.Method == http.MethodPost {
		questionCount, err := strconv.Atoi(r.FormValue("questionCount"))
		if err != nil || questionCount <= 0 {
			fmt.Fprintf(w, "Question count must be a positive integer")
			return
		}
		allData := getAllData()

		for i := 0; i < questionCount; i++ {
			selectedOption := r.FormValue("option" + strconv.Itoa(i+1))
			trueAnswer := r.FormValue("trueAnswer" + strconv.Itoa(i+1))

			answerData := make(map[string]bool)

			if trueAnswer == selectedOption {
				answerData[trueAnswer] = true
			} else {
				answerData[trueAnswer] = false
			}

			for j := 0; j < len(allData); j++ {
				if allData[j]["Answer"] == trueAnswer {
					if answerData[trueAnswer] {
						trueAnswerCount, err := strconv.Atoi(allData[j]["TrueAnswerCount"])
						if err != nil {
							fmt.Fprintf(w, "A programmatic error occured while increasing the TrueAnswerCount")
							return
						}
						trueAnswerCount += 1
						allData[j]["TrueAnswerCount"] = strconv.Itoa(trueAnswerCount)

					} else {
						wrongAnswerCount, err := strconv.Atoi(allData[j]["WrongAnswerCount"])
						if err != nil {
							fmt.Fprintf(w, "A programmatic error occured while increasing the WrongAnswerCount")
							return
						}
						wrongAnswerCount += 1
						allData[j]["WrongAnswerCount"] = strconv.Itoa(wrongAnswerCount)

					}

					askCount, err := strconv.Atoi(allData[j]["AskCount"])
					if err != nil {
						fmt.Fprintf(w, "A programmatic error occured while increasing the AskCount")
						return
					}
					askCount += 1
					allData[j]["AskCount"] = strconv.Itoa(askCount)
				}
			}
		}
		deleteCsvFile()
		for j := 0; j < len(allData); j++ {
			saveData(allData[j])
		}
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		data := map[string]string{"Data": "Quiz completed!"}
		tmpl.Execute(w, data)
	}
}
