package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// This code file is taken from the cli version of the same app.

func welcomeToTheProgram() {
	fmt.Println("Welcome!")
}

func showMenu() {
	fmt.Println()
	fmt.Println("1. Add Data")
	fmt.Println("2. Make Quiz")
	fmt.Println("3. Exit")
	fmt.Println()
}

func chooseFromMenu() string {
	fmt.Print("Choose: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println()
	userInput := input.Text()
	return userInput
}

func selectFromMenu() {
	for {
		showMenu()
		userInput := chooseFromMenu()
		if userInput == "1" {
			addData()
		} else if userInput == "2" {
			Quiz()
		} else if userInput == "3" {
			fmt.Println("Goodbye")
			os.Exit(0)
		}
		userInput = ""
	}
}

func Quiz() {
	fmt.Println("Make Quiz selected")

	allData := getAllData()
	totalQuestionCount := len(allData)
	fmt.Printf("Total question count: %d\n", totalQuestionCount)

	input := bufio.NewScanner(os.Stdin)
	fmt.Println("Question count you want to answer: ")
	input.Scan()

	userInput := input.Text()
	quizQuestionCount, err := strconv.Atoi(strings.Trim(userInput, " "))
	if err != nil {
		fmt.Println("You need to give an integer as input.")
		log.Fatal(err)
	}

	if quizQuestionCount < 4 {
		fmt.Println("Question count must be greater than or equal to 4.")
		return
	}

	fmt.Println("Prioritaze the least known questions (Y/n)")
	input.Scan()
	if input.Text() == "n" || input.Text() == "N" {
		askQuestionsRandomly(quizQuestionCount)
	} else {
		askLeastKnownQuestions(quizQuestionCount)
	}
}

func askQuestionsRandomly(chosenQuestionCount int) {
	fmt.Println("NOT IMPLEMENTED YET!")
}

func getPrioritazingValue(qaData map[string]string) int {
	wrongAnswerCount, err := strconv.Atoi(qaData["WrongAnswerCount"])
	if err != nil {
		log.Fatal(err)
	}

	trueAnswerCount, err := strconv.Atoi(qaData["TrueAnswerCount"])
	if err != nil {
		log.Fatal(err)
	}

	value := wrongAnswerCount - trueAnswerCount
	return value
}

func getPrioritazedData(qaData []map[string]string) []map[string]string {
	for i := 0; i < len(qaData); i++ {
		for j := 0; j < len(qaData); j++ {
			prioritazingValue1 := getPrioritazingValue(qaData[i])
			prioritazingValue2 := getPrioritazingValue(qaData[j])

			if prioritazingValue1 > prioritazingValue2 {
				qaData[i], qaData[j] = qaData[j], qaData[i]
			}
		}
	}
	return qaData
}

func removeByIndex(s []map[string]string, index int) []map[string]string {
	return append(s[:index], s[index+1:]...)
}

func getAnswerOptions(data []map[string]string, trueAnswer string, questionCount int, i int) []string {
	tempData := make([]map[string]string, len(data))
	copy(tempData, data)

	tempData = removeByIndex(tempData, i)

	randomAnswer1Index := rand.Intn(questionCount - 1)
	randomAnswer1 := tempData[randomAnswer1Index]["Answer"]
	tempData = removeByIndex(tempData, randomAnswer1Index)

	randomAnswer2Index := rand.Intn(questionCount - 2)
	randomAnswer2 := tempData[randomAnswer2Index]["Answer"]
	tempData = removeByIndex(tempData, randomAnswer2Index)

	randomAnswer3Index := rand.Intn(questionCount - 3)
	randomAnswer3 := tempData[randomAnswer3Index]["Answer"]

	answerOptions := []string{trueAnswer, randomAnswer1, randomAnswer2, randomAnswer3}

	// Dizideki eleman konumları rastgele değiştiriliyor
	for i := range answerOptions {
		j := rand.Intn(i + 1)
		answerOptions[i], answerOptions[j] = answerOptions[j], answerOptions[i]
	}

	return answerOptions

}

func askLeastKnownQuestions(chosenQuestionCount int) {
	fmt.Println("ASK LEAST KNOWN QUESTIONS FUNC STARTED")

	allData := getAllData()
	prioritazedData := getPrioritazedData(allData)

	for i := 0; i < chosenQuestionCount; i++ {

		trueAnswer := prioritazedData[i]["Answer"]

		answerOptions := getAnswerOptions(prioritazedData, trueAnswer, chosenQuestionCount, i)

		trueAnswerIndex := -1
		for j := 0; j < len(answerOptions); j++ {
			if answerOptions[j] == trueAnswer {
				trueAnswerIndex = j
				break
			}
		}
		if trueAnswerIndex == -1 {
			fmt.Println("Cevap şıklarının hazırlanması esnasında hata oluştu.")
			return
		}

		// Aşağıdaki kod bloğu soruyu ve cevap seçeneklerini yazdırıyor
		fmt.Println()
		fmt.Println(prioritazedData[i]["Question"] + "?")
		fmt.Println("A. " + answerOptions[0])
		fmt.Println("B. " + answerOptions[1])
		fmt.Println("C. " + answerOptions[2])
		fmt.Println("D. " + answerOptions[3])

		answerMapping := map[string]int{"A": 0, "B": 1, "C": 2, "D": 3}

		// Aşağıdaki kod bloğunda kullanıcıdan cevap alınıyor
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		userAnswer := strings.Trim(input.Text(), " ")
		userAnswer = strings.ToUpper(userAnswer)

		if userAnswer == "A" || userAnswer == "B" || userAnswer == "C" || userAnswer == "D" {
			answerIndex := answerMapping[userAnswer]
			if answerIndex == trueAnswerIndex {
				// Aşağıdaki kod bloğunda trueAnswerCount arttırılıyor
				trueAnswerCount, err := strconv.Atoi(prioritazedData[i]["TrueAnswerCount"])
				if err != nil {
					log.Fatal(err)
				}
				trueAnswerCount++
				prioritazedData[i]["TrueAnswerCount"] = strconv.Itoa(trueAnswerCount)
			} else {
				// Aşağıdaki kod bloğunda wrongAnswerCount arttırılıyor
				wrongAnswerCount, err := strconv.Atoi(prioritazedData[i]["WrongAnswerCount"])
				if err != nil {
					log.Fatal(err)
				}
				wrongAnswerCount++
				prioritazedData[i]["WrongAnswerCount"] = strconv.Itoa(wrongAnswerCount)
			}
		} else {
			// Aşağıdaki kod bloğunda wrongAnswerCount arttırılıyor
			wrongAnswerCount, err := strconv.Atoi(prioritazedData[i]["WrongAnswerCount"])
			if err != nil {
				log.Fatal(err)
			}
			wrongAnswerCount++
			prioritazedData[i]["WrongAnswerCount"] = strconv.Itoa(wrongAnswerCount)
		}

		// Aşağıdaki kod bloğunda askCount arttırılıyor
		askCount, err := strconv.Atoi(prioritazedData[i]["AskCount"])
		if err != nil {
			log.Fatal(err)
		}
		askCount++
		prioritazedData[i]["AskCount"] = strconv.Itoa(askCount)
	}

	// Aşağıdaki kod bloğu değişen veriyi allData'ya aktarıyor
	for i := 0; i < len(prioritazedData); i++ {
		for j := 0; j < len(allData); j++ {
			if prioritazedData[i]["Question"] == allData[j]["Question"] {
				allData[j] = prioritazedData[i]
			}
		}
	}

	// Aşağıdaki kod bloğu ile değişiklikler kaydediliyor
	deleteCsvFile()
	for i := 0; i < len(allData); i++ {
		saveData(allData[i])
	}
}

func addData() {
	data := getDataFromUser()
	saveData(data)
}

// gets QA from user
func getDataFromUser() map[string]string {
	qaMapping := make(map[string]string)
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Question: ")
	input.Scan()
	q := input.Text()
	fmt.Println()
	fmt.Print("Answer: ")
	input.Scan()
	a := input.Text()
	fmt.Println()
	qaMapping["Question"] = q
	qaMapping["Answer"] = a
	qaMapping["TrueAnswerCount"] = "0"
	qaMapping["WrongAnswerCount"] = "0"
	qaMapping["AskCount"] = "0"
	return qaMapping
}
