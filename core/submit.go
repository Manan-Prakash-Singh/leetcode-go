package core

import (
	"net/http"
    "fmt"
    "encoding/json"
    "github.com/gookit/color"
    "io"
    "os"
	"github.com/Manan-Prakash-Singh/leetcode-go/utils"
)

func getSubmissionId(fileName string, submit bool) (*RunTestCaseResponse, error) {

	questionID, problemName, lang, err := utils.ParseFileName(fileName)

	if err != nil {
		return nil, err
	}

	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	fileContent, err := io.ReadAll(file)

	fileContentStr := string(fileContent)

	if err != nil {
		return nil, err
	}

	testCases, err := utils.ParseTestCases(fileContentStr)

	if err != nil {
		return nil, err
	}

	jsonReq := map[string]string{
		"data_input":  testCases,
		"lang":        lang,
		"question_id": questionID,
		"typed_code":  fileContentStr,
	}

	requestBody, err := json.Marshal(jsonReq)

	if err != nil {
		return nil, err
	}

    runUrl := "https://leetcode.com/problems/" + problemName + "/interpret_solution/"
    submitUrl := "https://leetcode.com/problems/" + problemName + "/submit/"

    var request *http.Request

	if submit {
        request, err = utils.NewAuthRequest("POST",submitUrl,requestBody)
	} else {
        request, err = utils.NewAuthRequest("POST",runUrl,requestBody)
	}
	if err != nil {
		return nil, err
	}

    body, err := utils.SendRequest(request)

    if err != nil {
        return nil, err 
    }

	var response RunTestCaseResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func execute(id string) (*SubmissionResponse,error) {
    
    submissionResult := &SubmissionResponse{}
    url := "https://leetcode.com/submissions/detail/" + id + "/check/"
    request, err := utils.NewAuthRequest("GET",url,nil)
    if err != nil {
        return nil, err
    }
    response, err := utils.SendRequest(request)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(response,&submissionResult)
    if err != nil {
        return nil, err
    }
    return submissionResult, nil 
}

func RunTestCases(fileName string, submit bool) error {

    testCaseResponse, err := getSubmissionId(fileName,submit)

    if err != nil {
        return err
    }

    var id string

    if submit {
        id = fmt.Sprint(testCaseResponse.SubmissionId)
    }else{
        id = testCaseResponse.InterpretId
    }

    Result := &SubmissionResponse{}
    
    for {
        Result, err = execute(id)
        if err != nil {
            return err
        }
        if Result.State == "SUCCESS" {
            break
        }
    }

    OutputResult(testCaseResponse,Result)
    return nil
}


func OutputResult(testCaseResponse *RunTestCaseResponse, Result *SubmissionResponse) {

	statusMsg := Result.StatusMsg

	switch statusMsg {

	case "Compile Error":

		color.Redln("Compile Error")
		color.Redln(Result.FullCompileError)

	case "Accepted":

		if testCaseResponse.InterpretId == Result.SubmissionID {
			for i, testCase := range utils.TestCaseList {

				if Result.CodeAnswer[i] != Result.ExpectedCodeAnswer[i] {
					color.Redln("Wrong Answer")
				} else {
					color.Greenln("Correct Answer")
				}

				fmt.Println("Input")
				fmt.Println(testCase)
				fmt.Println("Output")
				fmt.Println(Result.CodeAnswer[i])
				fmt.Println("Expected")
				fmt.Println(Result.ExpectedCodeAnswer[i])

				color.Yellowln("----------------------------------------------------")
			}
		}

		if fmt.Sprint(testCaseResponse.SubmissionId) == Result.SubmissionID {

			color.Greenln("Accepted")
			fmt.Printf("Test cases passed: %d/%d\n", Result.TotalCorrect, Result.TotalTestcases)
			fmt.Printf("Runtime : %v [Beats : %0.2f%%]\n", Result.StatusRuntime, Result.RuntimePercentile)
			fmt.Printf("Memory : %v [Beats : %0.2f%%]\n", Result.StatusMemory, Result.MemoryPercentile)

		}

	case "Wrong Answer":

		color.Redln("Wrong Answer")
		fmt.Printf("Test cases passed: %d/%d\n", Result.TotalCorrect, Result.TotalTestcases)
		fmt.Println("Last test case executed: ")
		fmt.Println(Result.LastTestcase)
		fmt.Println("Expected Output:")
		fmt.Println(Result.ExpectedOutput)
		fmt.Println("Your Output:")
		fmt.Println(Result.CodeOutput.(string))

	case "Time Limit Exceeded":
		if testCaseResponse.InterpretId == Result.SubmissionID {
			color.Redln("Time Limit Exceeded")
		}
		if fmt.Sprint(testCaseResponse.SubmissionId) == Result.SubmissionID {

			color.Redln("Time Limit Exceeded")
			fmt.Printf("Test cases passed: %d/%d\n", Result.TotalCorrect, Result.TotalTestcases)
			fmt.Println("Last test case executed: ")
			fmt.Println(Result.LastTestcase)
			fmt.Println("Expected Output:")
			fmt.Println(Result.ExpectedOutput)
			fmt.Println("Your Output:")
			fmt.Println(Result.CodeOutput.(string))

		}
	}

}
func SubmitCode(fileName string) error {

	if err := RunTestCases(fileName, true); err != nil {
		return err
	}

	return nil
}

