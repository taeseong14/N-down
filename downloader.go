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
)

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
	Title string  `json:"title"`
	Ep    float64 `json:"ep"`
	Date  string  `json:"date"`
	Link  string  `json:"link"`
}

type Chan struct {
	Result string
	Index  int
}

func main() {
	fmt.Println("Novelpia Downloader by taeseong14")
	fmt.Println("\n[Login]")
	var LOGINKEY, id, pw string
	dat, _ := os.ReadFile("account.txt")
	if dat != nil {
		s := strings.Split(string(dat), "\n")
		id, pw = s[0], s[1]
		fmt.Print("\nlogin with ", id, "...\n")
	} else {
		fmt.Print("\nid: ")
		fmt.Scan(&id)

		if !strings.Contains(id, "@") {
			id = id + "@gmail.com"
			fmt.Println("id:", id)
		}

		fmt.Print("pw: ")
		fmt.Scan(&pw)
		fmt.Printf("\rlogin...")
	}

	json_data, _ := json.Marshal(User{id, pw})
	resp, _ := http.Post("https://b-p.msub.kr/novelp/login", "application/json", bytes.NewBuffer(json_data))

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if res["err"] != nil {
		fmt.Println("\n\nError:", res["err"])
		dat, _ := os.ReadFile("account.txt")
		if dat != nil {
			os.Remove("account.txt")
			fmt.Println("./account.txt removed")
		}
		end()
		return
	}

	LOGINKEY = res["result"].(string)

	fmt.Printf("\rlogin success\n")

	dat, _ = os.ReadFile("account.txt")
	if dat == nil {
		fmt.Println("login data saved in ./account.txt")
		os.WriteFile("account.txt", []byte(id+"\n"+pw), 0644)
	}

	var bookid int
	fmt.Print("\nbookId: ")
	fmt.Scan(&bookid)

	fmt.Printf("\rloading book info...")

	resp, _ = http.Get("https://b-p.msub.kr/novelp/info/?id=" + strconv.Itoa(bookid))
	json.NewDecoder(resp.Body).Decode(&res)

	if res["err"] != nil {
		fmt.Printf("\rError: %s", res["err"])
		end()
		return
	}

	info := res["result"].(map[string]interface{})
	fmt.Printf("\r[%s - %s] is right? (y/n): ", info["title"], info["author"])
	var yn string
	fmt.Scan(&yn)
	if !strings.Contains(yn, "y") {
		end()
		return
	}

	fmt.Println()

	fmt.Printf("\rGet page.")

	resp, _ = http.Get("https://b-p.msub.kr/novelp/list/?p=all&id=" + strconv.Itoa(bookid))
	var resResult Results
	json.NewDecoder(resp.Body).Decode(&resResult)
	fmt.Printf("\rGet page. %.0f/%.0f\n\n", resResult.P+1, resResult.P+1)

	result := make([]string, 1000)

	ch := make(chan Chan)

	for i := range resResult.Result {
		if i%10 == 1 {
			time.Sleep(time.Millisecond * 100)
		}
		fmt.Printf("\rRequest EP. %d/%.0f [%.0f%s]", i+1, resResult.Last, float64(i)/resResult.Last*100, "%")
		go getEp(LOGINKEY, &resResult.Result[i], resResult.Last, i, ch, 1)
		time.Sleep(time.Second / 100)
	}
	fmt.Printf("\rRequest EP. %.0f/%.0f [100%s]", resResult.Last, resResult.Last, "%")

	fmt.Println()

	fmt.Printf("\rGet EP. %d/%.0f", 0, resResult.Last)

	for a := range resResult.Result {
		result_ := <-ch
		fmt.Printf("\rGet EP. %d/%.0f [%.0f%s]", a, resResult.Last, float64(a)/resResult.Last*100, "%")
		result[result_.Index] = result_.Result
	}

	os.WriteFile("result.txt", []byte(strings.TrimSpace(strings.Join(result, "\n\n\n\n\n\n\n\n\n\n"))), 0644)

	fmt.Println("\n\nDone! check ./result.txt")

	end()
}

func getEp(LOGINKEY string, page *Result, max float64, i int, ch chan Chan, tried int) {
	if tried == 3 {
		ch <- Chan{page.Title + "\n\n\n\n\nError: 소설 정보를 불러올 수 없음", i}
		return
	}
	req, _ := http.NewRequest("GET", "https://b-p.msub.kr/novelp/view/?id="+page.Link, nil)
	req.Header.Set("Cookie", "LOGINKEY="+LOGINKEY+";")

	resp, _ := http.DefaultClient.Do(req)

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if res["err"] != nil {
		fmt.Println("\n\n" + res["err"].(string))
		ch <- Chan{page.Title + "\n\n\n\n\n" + res["err"].(string), i}
		return
	}

	if res["result"] == nil {
		getEp(LOGINKEY, page, max, i, ch, tried+1)
		return
	}

	ch <- Chan{page.Title + "\n\n\n\n\n" + res["result"].(string), i}
}

func end() {
	fmt.Print("\n\npress enter to exit...")
	fmt.Scanln()
	fmt.Scanln()
}
