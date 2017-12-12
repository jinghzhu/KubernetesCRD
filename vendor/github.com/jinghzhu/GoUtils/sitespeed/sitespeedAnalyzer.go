package sitespeed

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
)

const CountRule = 39   // num of sitespeed rules
const Threshold50 = 50 // used for analyze the percent of data above scroe 50
const WorstScore = 0   // the bottom score of performance budget
const BestScore = 90   // the top score of performance budget
const RulesNum = 39

// the following 3 strcuts are used to analyze performance budget based on sitespeed.io
// for more details, please read example test.json
type PerformanceBudget struct {
	CookieUrl string          `json:"cookieUrl,omitempty"`
	BasUrl    string          `json:"baseUrl"`
	Timeout   int             `json:"mochaTimeout,omitempty"`
	Id        string          `json:"id,omitempty"`
	Pw        string          `json:"pw,omitempty"`
	TestCases []JsonTestCases `json:"testCases"`
}

type JsonTestCases struct {
	Description string                 `json:"description,omitempty"`
	Pathname    string                 `json:"pathname"`
	UrlParams   map[string]interface{} `json:"urlParams"`
	Budget      JsonBudget             `json:"budget"`
}

type JsonBudget struct {
	Rules JsonRules `json:"rules"`
}

// all sitespeed.io rules
type JsonRules struct {
	Criticalpath                int `json:"criticalpath"`
	Spof                        int `json:"spof"`
	Cssnumreq                   int `json:"cssnumreq"`
	Cssimagesnumreq             int `json:"cssimagesnumreq"`
	Jsnumreq                    int `json:"jsnumreq"`
	Yemptysrc                   int `json:"yemptysrc"`
	Ycompress                   int `json:"ycompress"`
	Ycsstop                     int `json:"ycsstop"`
	Yjsbottom                   int `json:"yjsbottom"`
	Yexpressions                int `json:"yexpressions"`
	Ydns                        int `json:"ydns"`
	Yminify                     int `json:"yminify"`
	Redirects                   int `json:"redirects"`
	Noduplicates                int `json:"noduplicates"`
	Yetags                      int `json:"yetags"`
	Yxhr                        int `json:"yxhr"`
	Yxhrmethod                  int `json:"yxhrmethod"`
	Mindom                      int `json:"mindom"`
	Yno404                      int `json:"yno404"`
	Ymincookie                  int `json:"ymincookie"`
	Ycookiefree                 int `json:"ycookiefree"`
	Ynofilter                   int `json:"ynofilter"`
	Avoidscalingimages          int `json:"avoidscalingimages"`
	Yfavicon                    int `json:"yfavicon"`
	Thirdparyasyncjs            int `json:"thirdpartyasyncjs"`
	Csspring                    int `json:"cssprint"`
	Cssinheaddomain             int `json:"cssinheaddomain"`
	Syncjsinhead                int `json:"syncjsinhead"`
	Avoidfont                   int `json:"avoidfont"`
	Totalrequests               int `json:"totalrequests"`
	Expiresmod                  int `json:"expiresmod"`
	Longexpirehead              int `json:"longexpirehead"`
	Nodnslookupswhenfewrequests int `json:"nodnslookupswhenfewrequests"`
	Inlinecsswhenfewrequest     int `json:"inlinecsswhenfewrequest"`
	Textcontent                 int `json:"textcontent"`
	Thirdpartyversions          int `json:"thirdpartyversions"`
	Ycdn                        int `json:"ycdn"`
	Connectionclose             int `json:"connectionclose"`
	Privateheaders              int `json:"privateheaders"`
}

func main() {
	file, err := ioutil.ReadFile("./sitespeedAnalyzer.json")
	if err != nil {
		fmt.Println("File error: %v\n", err)
		return
	}

	var budgets PerformanceBudget
	err = json.Unmarshal(file, &budgets)

	if err != nil {
		fmt.Println("Unmarshal error: %v\n", err)
		return
	}

	err = GetBudgetStatics(budgets)
	// err = GetRuleScore(budgets)
	if err != nil {
		fmt.Println("Analysis error: %v\n", err)
		return
	}

}

// For a site, as it have many paths, sitespeed.io would generate many performance budgets.
// This function is trying to calcute the average value in each rule field.
func GetRuleScore(budgets PerformanceBudget) error {
	var countBudget int = len(budgets.TestCases)
	if countBudget < 1 {
		errMsg := "Error! 0 test cases"
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}

	ruleList := reflect.ValueOf(budgets.TestCases[0].Budget.Rules)

	for i := 0; i < RulesNum; i++ {
		ruleField := ruleList.Type().Field(i)
		ruleName := ruleField.Tag
		fmt.Println("\nready to deal with rule: " + ruleName.Get("json"))
		var sum int64 = 0
		for j := 0; j < countBudget; j++ {
			values := reflect.ValueOf(budgets.TestCases[j].Budget.Rules)
			value := values.Field(i)
			val := value.Int()
			sum += val
		}
		mean := sum / int64(countBudget)
		fmt.Println("mean = " + strconv.FormatInt(mean, 10))
	}

	return nil
}

// As mentioned before, if there are more than 1 paths set in test.json,
// sitespeed would generate more than one performance results (one for each path).
// This is used to analyze performance test result. It provides:
// (1) the best score value and its rule name;
// (2) how many paths meet (1);
// (3) the worst score value and its rule name;
// (4) how many paths meet (3);
// (5) the percent of rule scores in range (-1, 50) and [50, 90) respectively
func GetBudgetStatics(budgets PerformanceBudget) error {
	var countBudget int = len(budgets.TestCases)
	if countBudget < 1 {
		errMsg := "Error! 0 test cases"
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}

	var max, min, countVal049, countVal5089, count, mean, countBestSco, countWorstSco int64 = -1, 90, 0, 0, 0, 0, 0, 0

	var total int64 = 0

	for i := 0; i < countBudget; i++ {
		values := reflect.ValueOf(budgets.TestCases[i].Budget.Rules)

		for j := 0; j < values.NumField(); j++ {
			value := values.Field(j)
			valField := values.Type().Field(j)
			rule := valField.Tag

			val := value.Int()
			fmt.Println(rule.Get("json"))
			fmt.Println(val)

			if val > 90 || val < -1 {
				errMsg := "Score value should be between -1 and 90."
				fmt.Println(errMsg + " Field #" + strconv.Itoa(j+1))
			}

			if val > max {
				max = val
			} else if val < min {
				min = val
			}
			if val >= Threshold50 {
				if val == BestScore {
					countBestSco++

					fmt.Println("******")
					fmt.Println("BestScore: " + rule.Get("json") + strconv.FormatInt(val, 10))
					fmt.Println("******")
					fmt.Println("\n\n")
				} else {
					countVal5089++

					fmt.Println("******")
					fmt.Println("[50, 90): " + rule.Get("json") + strconv.FormatInt(val, 10))
					fmt.Println("******")
					fmt.Println("\n\n")
				}
			} else {
				if val == WorstScore {

					fmt.Println("******")
					fmt.Println("WorstScore: " + rule.Get("json") + strconv.FormatInt(val, 10))
					fmt.Println("******")
					fmt.Println("\n\n")
					countWorstSco++
				} else {

					fmt.Println("******")
					fmt.Println("(-1, 50): " + rule.Get("json") + strconv.FormatInt(val, 10))
					fmt.Println("******")
					fmt.Println("\n\n")
					countVal049++
				}
			}

			total += val
			count++
		}
	}

	mean = total / count
	val5089Pct := float64(countVal5089) / float64(count)
	valBestPct := float64(countBestSco) / float64(count)
	val049Pct := float64(countVal049) / float64(count)
	valWorstPct := float64(countWorstSco) / float64(count)

	fmt.Println("max = " + strconv.FormatInt(max, 10))
	fmt.Println("min = " + strconv.FormatInt(min, 10))
	fmt.Println("\n\n\n\n\n\n\n\n\n")

	fmt.Println("Num of rules data = " + strconv.FormatInt(count, 10))
	fmt.Println("\n")

	fmt.Println("Mean score of performance = " + strconv.FormatInt(mean, 10))
	fmt.Println("\n")

	fmt.Println("Num of worst performance rules(rule.score = 0 ) = " + strconv.FormatInt(countWorstSco, 10))
	fmt.Println("Percent of worst performance rules = ", valWorstPct*100)
	fmt.Println("\n")

	fmt.Println("Num of best performance rules( = 90 ) = " + strconv.FormatInt(countBestSco, 10))
	fmt.Println("Percent of best performance rules = ", valBestPct*100)
	fmt.Println("\n")

	fmt.Println("Num of rules whose score is between 50 and 90([50, 90]) = " + strconv.FormatInt(countVal5089, 10))
	fmt.Println("Percent of rules whose score is between 50 and 90([50, 90]) = ", val5089Pct*100)
	fmt.Println("\n")

	fmt.Println("Num of rules whose score is between 0 and 50([0, 50)) = " + strconv.FormatInt(countVal049, 10))
	fmt.Println("Percent of rules whose score is between 0 and 50([0, 50)) = ", val049Pct*100)
	fmt.Println("\n")

	return nil
}
