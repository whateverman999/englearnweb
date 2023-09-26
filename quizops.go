package main

import (
	"log"
)

func getQuizData(chosenQuestionCount int, prioritizedData []map[string]string, allData []map[string]string) []map[string]string {
	questionCount := chosenQuestionCount
	questionAndAnswerOptionsMappingArr := make([]map[string]string, questionCount)
	for j := 0; j < questionCount; j++ {
		questionAndAnswerOptionsMappingArr[j] = make(map[string]string)
	}

	for i := 0; i < questionCount; i++ {

		trueAnswer := prioritizedData[i]["Answer"]

		answerOptions := getAnswerOptions(prioritizedData, trueAnswer, questionCount, i)

		trueAnswerIndex := -1
		for j := 0; j < len(answerOptions); j++ {
			if answerOptions[j] == trueAnswer {
				trueAnswerIndex = j
				break
			}
		}
		if trueAnswerIndex == -1 {
			log.Fatal("An error occured while preparing the answer options.")
		}

		questionAndAnswerOptionsMappingArr[i]["Question"] = prioritizedData[i]["Question"]
		questionAndAnswerOptionsMappingArr[i]["TrueAnswer"] = trueAnswer
		questionAndAnswerOptionsMappingArr[i]["OptionA"] = answerOptions[0]
		questionAndAnswerOptionsMappingArr[i]["OptionB"] = answerOptions[1]
		questionAndAnswerOptionsMappingArr[i]["OptionC"] = answerOptions[2]
		questionAndAnswerOptionsMappingArr[i]["OptionD"] = answerOptions[3]
	}
	return questionAndAnswerOptionsMappingArr
}

func saveChanges(prioritizedData []map[string]string, allData []map[string]string) {
	// Aşağıdaki kod bloğu değişen veriyi allData'ya aktarıyor
	for i := 0; i < len(prioritizedData); i++ {
		for j := 0; j < len(allData); j++ {
			if prioritizedData[i]["Question"] == allData[j]["Question"] {
				allData[j] = prioritizedData[i]
			}
		}
	}

	// Aşağıdaki kod bloğu ile değişiklikler kaydediliyor
	deleteCsvFile()
	for i := 0; i < len(allData); i++ {
		saveData(allData[i])
	}
}
