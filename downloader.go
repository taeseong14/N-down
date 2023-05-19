package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v4"
)

const version string = "0.0.9"

type User struct {
	Id string `json:"id"`
	Pw string `json:"pw"`
}

type Results struct {
	Result []Result `json:"result"`
	Last   float64  `json:"last"`
	P      float64  `json:"p"`
}

type Result struct {
	Title string `json:"title"`
	Ep    string `json:"ep"`
	Date  string `json:"date"`
	Link  string `json:"link"`
}

type Chan struct {
	Result string
	Index  int
}

func main() {
	fmt.Println(aurora.Cyan("Novelpia Downloader by taeseong14").Bold(), aurora.Gray(12, "v"+version), aurora.BgWhite("[Github]").Black().Hyperlink("https://github.com/taeseong14/N-down"))
	fmt.Print(aurora.BgIndex(16, "\n[Login]\n\n"))
	var LOGINKEY, id, pw string
	dat, _ := os.ReadFile("account.txt")
	if dat != nil {
		s := strings.Split(string(dat), "\n")
		id, pw = s[0], s[1]
		fmt.Print(aurora.BrightYellow("login with "), aurora.Cyan(id), "...\n")
	} else {
		fmt.Print("\nid: ")
		fmt.Scan(&id)
		if !strings.Contains(id, "@") {
			id = id + "@gmail.com"
			fmt.Println(aurora.Green("id: " + id))
		}
		fmt.Print("pw: ")
		fmt.Scan(&pw)
		fmt.Println()

		fmt.Print("\rlogin...")
	}

	json_data, _ := json.Marshal(User{id, pw})
	resp, _ := http.Post("https://b-p.msub.kr/novelp/login?v="+version, "application/json", bytes.NewBuffer(json_data))

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if res["err"] != nil {
		if res["err"].(string) == "New Version Released" {
			fmt.Println(aurora.Yellow("\rNew Version Released:"), aurora.BgWhite(res["v"]).Black().Hyperlink("https://github.com/taeseong14/N-down/releases/tag/v"+res["v"].(string)))
			end()
			return
		}
		fmt.Println(aurora.BrightRed("\n\nError:"), aurora.BrightRed(res["err"]))
		dat, _ := os.ReadFile("account.txt")
		if dat != nil {
			os.Remove("account.txt")
			fmt.Println("\n./account.txt removed")
		}
		end()
		return
	}

	LOGINKEY = res["result"].(string)

	fmt.Println(aurora.BrightGreen("\rlogin success"))

	dat, _ = os.ReadFile("account.txt")
	if dat == nil {
		fmt.Println(aurora.BrightBlue("login data saved in ./account.txt"))
		os.WriteFile("account.txt", []byte(id+"\n"+pw), 0644)
	}

	var bookid int
	fmt.Print(aurora.BrightMagenta("\nbookId: "))
	fmt.Scan(&bookid)

	fmt.Print(aurora.BrightCyan("\rloading book info..."))

	resp, _ = http.Get("https://b-p.msub.kr/novelp/info/?id=" + strconv.Itoa(bookid))
	json.NewDecoder(resp.Body).Decode(&res)

	if res["err"] != nil {
		fmt.Print(aurora.BrightRed("\rError:"), aurora.BrightRed(res["err"]))
		end()
		return
	}

	info := res["result"].(map[string]interface{})
	title := info["title"].(string)
	title = strings.ReplaceAll(title, "/", "／")
	title = strings.ReplaceAll(title, "\\", "＼")
	title = strings.ReplaceAll(title, ":", "：")
	title = strings.ReplaceAll(title, "*", "＊")
	title = strings.ReplaceAll(title, "?", "？")
	title = strings.ReplaceAll(title, "\"", "＂")
	title = strings.ReplaceAll(title, "<", "＜")
	title = strings.ReplaceAll(title, ">", "＞")
	title = strings.ReplaceAll(title, "|", "｜")

	fmt.Printf("\r[%s - %s] is right? (y/n): ", aurora.Cyan(title), aurora.BrightGreen(info["author"]))

	var yn string
	fmt.Scan(&yn)
	if !strings.Contains(yn, "y") && !strings.Contains(yn, "Y") {
		end()
		return
	}

	fmt.Println()

	fmt.Printf("\rGet page.")

	resp, _ = http.Get("https://b-p.msub.kr/novelp/list/?p=all&id=" + strconv.Itoa(bookid))
	var resResult Results
	json.NewDecoder(resp.Body).Decode(&resResult)
	fmt.Printf("\rGet page. %.0f/%.0f", resResult.P+1, resResult.P+1)
	fmt.Print(aurora.Green(" [100%]\n\n"))

	result := make([]string, int(resResult.P)*20+300)

	ch := make(chan Chan)

	// 300개씩 끊어서 요청
	jLimit := len(resResult.Result) / 300
	for j := 0; j <= jLimit; j++ {
		// 회차 = j * 300 + i

		// 출력: [${j} of ${jLimit}] ${i} of rest
		left := len(resResult.Result) - j*300
		if left > 300 {
			left = 300
		}

		// for i := range resResult.Result {
		for i := 0; i < left; i++ {
			if i%50 == 0 {
				time.Sleep(time.Second / 2)
			}
			fmt.Printf("\r[%d of %d] Requesting %d/%d", j+1, jLimit+1, i+1, left)
			go getEp(LOGINKEY, &resResult.Result[i+j*300], resResult.Last, j*300+i, ch, 1)
			time.Sleep(time.Second / 50) // 50 req/s
		}
		fmt.Println()

		for a := 0; a < left; a++ {
			result_ := <-ch

			fmt.Printf("\r[%d of %d] %d/%d", j+1, jLimit+1, a+1, left)

			result[result_.Index] = result_.Result
		}

		fmt.Print(aurora.Green(" Done!\n\n"))

		if len(resResult.Result) > 300 && left == 300 {
			fmt.Print("\rResting 5 sec...")
			time.Sleep(time.Second * 3)
		}
	}

	if _, err := os.Stat("result"); os.IsNotExist(err) {
		fmt.Print(aurora.BrightRed("\rresult dir not exist"))
		os.Mkdir("result", 0755)
		fmt.Println(aurora.BrightGreen("\rresult dir created"))
	}

	os.WriteFile("result/"+title+".txt", []byte(strings.TrimSpace(strings.Join(result, "\n\n\n\n\n\n\n\n\n\n"))), 0644)

	fmt.Println(aurora.Green("\n\nDone! check ./result/" + title + ".txt"))

	end()
}

func getEp(LOGINKEY string, page *Result, max float64, i int, ch chan Chan, tried int) {
	if tried == 4 {
		ch <- Chan{page.Title + "\n\n\n\n\nError: 소설 정보를 불러올 수 없음", i}
		return
	}
	req, _ := http.NewRequest("GET", "https://b-p.msub.kr/novelp/view/?id="+page.Link, nil)
	req.Header.Set("Cookie", "LOGINKEY="+LOGINKEY+";")

	resp, _ := http.DefaultClient.Do(req)

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if res["err"] != nil {
		fmt.Println(aurora.BrightYellow("\n\n  at EP." + strconv.Itoa(i+1) + ":"))
		fmt.Print(aurora.BrightRed(res["err"].(string) + "\n\n"))
		ch <- Chan{page.Title + "\n\n\n\n\n" + res["err"].(string), i}
		return
	}

	if res["result"] == nil {
		getEp(LOGINKEY, page, max, i, ch, tried+1)
		return
	}

	if page.Ep != "BONUS" {
		page.Ep += "화"
	}
	ch <- Chan{"[" + page.Ep + "] " + page.Title + "\n\n\n\n\n" + res["result"].(string), i}
}

func end() {
	fmt.Print(aurora.BgWhite("\n\npress enter to exit...").Black())
	fmt.Scanln()
	fmt.Scanln()
}
