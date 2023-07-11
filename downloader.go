package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v4"
)

const version string = "0.1.3"

var space string

type User struct {
	Id string `json:"id"`
	Pw string `json:"pw"`
}
type User2 struct {
	LOGINKEY string `json:"LOGINKEY"`
}

type Results struct {
	Result []Result `json:"result"`
	Last   float64  `json:"last"`
	P      float64  `json:"p"`
	Cont   bool     `json:"cont"` // continue
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

type respBody struct {
	S []struct {
		Text string `json:"text"`
	} `json:"s"`
}

var setting map[string]interface{}
var useColors bool

func main() {
	// read file: settings.txt
	dat, _ := os.ReadFile("settings.txt")
	set := string(dat)
	if dat == nil {
		set = "{\n    \"account.auto_login\": true,\n    \"account.auto_login_file\": \"account.txt\",\n    \"account.default_mail\": \"@gmail.com\",\n    \"account.login_with_cookie\": false,\n    \"cmd.check_with_yn\": true,\n    \"cmd.exit_when_finish\": true,\n    \"cmd.max_try_per_episode\": 5,\n    \"cmd.use_colors\": true,\n    \"result.image_display\": true,\n    \"result.image_format\": \"[이미지: ${link}]\",\n    \"result.directory_name\": \"result\",\n    \"result.file_name\": \"${title}.txt\",\n    \"result.space_between_episodes\": \"\\n\\n\\n\\n\\n\\n\\n\\n\\n\\n\"\n}"
		os.WriteFile("settings.txt", []byte(set), 0644)
	}
	setting = make(map[string]interface{})
	err := json.Unmarshal([]byte(set), &setting)
	if err != nil {
		fmt.Println(err)
		fmt.Println("settings.txt 파일이 손상되었습니다. settings.txt 파일을 지운 후 다시 실행해주세요.")
		fmt.Println("(참고: 따옴표나 쉼표 등 json 형식을 지켜주세요.)")
		end()
		return
	}

	loginDataFile := setting["account.auto_login_file"].(string)
	if setting["account.auto_login"].(bool) {
		dat, _ = os.ReadFile(loginDataFile)
		if dat == nil {
			os.WriteFile(loginDataFile, []byte(""), 0644)
		}
	}

	space = setting["result.space_between_episodes"].(string)
	useColors = setting["cmd.use_colors"].(bool)

	if useColors {
		fmt.Println(aurora.Cyan("Novelpia Downloader by taeseong14").Bold(), aurora.Gray(12, "v"+version), aurora.BgWhite("[Github]").Black().Hyperlink("https://github.com/taeseong14/N-down"))
		fmt.Print(aurora.BgIndex(16, "\n[Login]\n\n"))
	} else {
		fmt.Println("Novelpia Downloader by taeseong14", "v"+version, "[Github] github.com/taeseong14/N-down")
		fmt.Print("\n[Login]\n\n")
	}
	var LOGINKEY, id, pw string
	login_with_cookie := setting["account.login_with_cookie"]

	dat, _ = os.ReadFile(loginDataFile)
	if string(dat) != "" && setting["account.auto_login"].(bool) {
		s := strings.Split(string(dat), "\n")
		if login_with_cookie == true {
			LOGINKEY = s[0]
		} else {
			id, pw = s[0], s[1]
			if useColors {
				fmt.Print(aurora.BrightYellow("login with "), aurora.Cyan(id), "...\n")
			} else {
				fmt.Print("login with ", id, "...\n")
			}
		}
	} else {
		if login_with_cookie == true {
			fmt.Print("LOGINKEY: ")
			fmt.Scan(&LOGINKEY)
		} else {
			fmt.Print("\nid: ")
			fmt.Scan(&id)
			if !strings.Contains(id, "@") {
				id = id + setting["account.default_mail"].(string)
				if useColors {
					fmt.Println(aurora.Green("id: " + id))
				} else {
					fmt.Println("id: " + id)
				}
			}
			fmt.Print("pw: ")
			fmt.Scan(&pw)
			fmt.Println()

			fmt.Print("\rlogin...")
		}
	}

	var res map[string]interface{}

	if login_with_cookie == false {
		json_data, _ := json.Marshal(User{id, pw})
		resp, _ := http.Post("https://b-p.msub.kr/novelp/login?v="+version, "application/json", bytes.NewBuffer(json_data))

		json.NewDecoder(resp.Body).Decode(&res)

		if res["err"] != nil {
			if res["err"].(string) == "New Version Released" {
				if useColors {
					fmt.Println(aurora.Yellow("\rNew Version Released:"), aurora.BgWhite(res["v"]).Black().Hyperlink("https://github.com/taeseong14/N-down/releases/tag/v"+res["v"].(string)))
				} else {
					fmt.Println("\rNew Version Released: https://github.com/taeseong14/N-down/releases/tag/v" + res["v"].(string))
				}
				fmt.Println()
			} else {
				if useColors {
					fmt.Println(aurora.BrightRed("\n\nError:"), aurora.BrightRed(res["err"]))
				} else {
					fmt.Println("\n\nError:", res["err"])
				}
				// dat, _ := os.ReadFile(loginDataFile)
				// if dat != nil {
				// 	os.Remove(loginDataFile)
				// 	fmt.Println("\naccount file removed")
				// }
				end()
				return
			}
		}

		LOGINKEY = res["result"].(string)

		if useColors {
			fmt.Println(aurora.BrightGreen("\rlogin success"))
		} else {
			fmt.Println("\rlogin success")
		}
	} else {
		json_data, _ := json.Marshal(User2{LOGINKEY})
		resp, _ := http.Post("https://b-p.msub.kr/novelp/login?v="+version, "application/json", bytes.NewBuffer(json_data))

		json.NewDecoder(resp.Body).Decode(&res)

		if res["err"] != nil {
			if res["err"].(string) == "New Version Released" {
				if useColors {
					fmt.Println(aurora.Yellow("\rNew Version Released:"), aurora.BgWhite(res["v"]).Black().Hyperlink("https://github.com/taeseong14/N-down/releases/tag/v"+res["v"].(string)))
				} else {
					fmt.Println("\rNew Version Released: https://github.com/taeseong14/N-down/releases/tag/v" + res["v"].(string))
				}
				fmt.Println()
			} else {
				if useColors {
					fmt.Println(aurora.BrightRed("\n\nError:"), aurora.BrightRed(res["err"]))
				} else {
					fmt.Println("\n\nError:", res["err"])
				}
				// dat, _ := os.ReadFile(loginDataFile)
				// if dat != nil {
				// 	os.Remove(loginDataFile)
				// 	fmt.Println("\naccount file removed")
				// }
				end()
				return
			}
		}
	}

	if setting["account.auto_login"].(bool) {
		dat, _ = os.ReadFile(loginDataFile)
		if string(dat) == "" {
			if useColors {
				fmt.Println(aurora.BrightBlue("login data saved in " + loginDataFile))
			} else {
				fmt.Println("login data saved in " + loginDataFile)
			}
			if login_with_cookie == true {
				os.WriteFile(loginDataFile, []byte(LOGINKEY), 0644)
			} else {
				os.WriteFile(loginDataFile, []byte(id+"\n"+pw), 0644)
			}
		}
	}

	for {

		var bookid int
		var title, author string

		for {

			if useColors {
				fmt.Print(aurora.BrightMagenta("\nbookId: "))
			} else {
				fmt.Print("\nbookId: ")
			}
			fmt.Scan(&bookid)

			if useColors {
				fmt.Print(aurora.BrightCyan("\rloading book info..."))
			} else {
				fmt.Print("\rloading book info...")
			}

			resp, _ := http.Get("https://b-p.msub.kr/novelp/info/?id=" + strconv.Itoa(bookid))
			json.NewDecoder(resp.Body).Decode(&res)

			if res["err"] != nil {
				if useColors {
					fmt.Print(aurora.BrightRed("\rError:"), aurora.BrightRed(res["err"]))
				} else {
					fmt.Print("\rError:", res["err"])
				}
				fmt.Println()
				continue
			}

			info := res["result"].(map[string]interface{})
			title = info["title"].(string)
			title = strings.ReplaceAll(title, "/", "／")
			title = strings.ReplaceAll(title, "\\", "＼")
			title = strings.ReplaceAll(title, ":", "：")
			title = strings.ReplaceAll(title, "*", "＊")
			title = strings.ReplaceAll(title, "?", "？")
			title = strings.ReplaceAll(title, "\"", "＂")
			title = strings.ReplaceAll(title, "<", "＜")
			title = strings.ReplaceAll(title, ">", "＞")
			title = strings.ReplaceAll(title, "|", "｜")
			author = info["author"].(string)

			if setting["cmd.check_with_yn"].(bool) {
				if useColors {
					fmt.Printf("\r[%s - %s] is right? (y/n): ", aurora.Cyan(title), aurora.BrightGreen(author))
				} else {
					fmt.Printf("\r[%s - %s] is right? (y/n): ", title, author)
				}

				var yn string
				fmt.Scan(&yn)
				if !strings.Contains(yn, "y") && !strings.Contains(yn, "Y") {
					fmt.Println()
				} else {
					break
				}
			} else {
				if useColors {
					fmt.Printf("\rDownloading [%s - %s]...\n", aurora.Cyan(title), aurora.BrightGreen(author))
				} else {
					fmt.Printf("\rDownloading [%s - %s]...\n", title, author)
				}
				break
			}

		}

		fmt.Println()

		// read /result/{titleformat} file
		// fileName = setting["result.file_name"].(string) .replace("${title}", title).replace("${author}", author).replace("${id}", bookId)
		fileName := setting["result.file_name"].(string)
		fileName = strings.ReplaceAll(fileName, "${title}", title)
		fileName = strings.ReplaceAll(fileName, "${author}", author)
		fileName = strings.ReplaceAll(fileName, "${id}", strconv.Itoa(bookid))

		d, _ := os.ReadFile(setting["result.directory_name"].(string) + "/" + fileName)
		l := "0"
		if string(d) != "" {
			arr := strings.Split(string(d), "\n")
			for _, v := range arr {
				// if v startsWith("[") and contains "화]"
				if strings.HasPrefix(v, "[") && strings.Contains(v, "화]") {
					l = v[1:strings.Index(v, "화]")]
				}
			}
			if l == "BONUS" {
				l = "0"
			}
		}

		fmt.Println("last ep:", l)
		fmt.Printf("\rGet page.")
		pageLink := fmt.Sprintf("https://b-p.msub.kr/novelp/list/?p=all&last=%s&id=%d", l, bookid)
		resp, _ := http.Get(pageLink)
		var resResult Results
		json.NewDecoder(resp.Body).Decode(&resResult)
		fmt.Printf("\rGet page. %.0f/%.0f", resResult.P+1, resResult.P+1)
		if useColors {
			fmt.Print(aurora.Green(" [100%]\n\n"))
		} else {
			fmt.Print(" [100%]\n\n")
		}

		if len(resResult.Result) == 0 {
			if useColors {
				fmt.Println(aurora.BrightRed("No new episode"))
			} else {
				fmt.Println("No new episode")
			}
			if setting["cmd.exit_when_finish"].(bool) {
				end()
				return
			}
			continue
		}

		result := make([]string, 0)
		for i := 0; i < len(resResult.Result); i++ {
			// append result: i as string
			result = append(result, strconv.Itoa(i))
		}

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
				time.Sleep(time.Second / 30) // 30 req/s
			}
			fmt.Println()

			for a := 0; a < left; a++ {
				result_ := <-ch

				fmt.Printf("\r[%d of %d] %d/%d", j+1, jLimit+1, a+1, left)
				result[result_.Index] = result_.Result
			}

			if useColors {
				fmt.Print(aurora.Green(" Done!\n\n"))
			} else {
				fmt.Print(" Done!\n\n")
			}

			if len(resResult.Result) > 300 && left == 300 {
				fmt.Print("\rResting 5 sec...")
				time.Sleep(time.Second * 3)
			}
		}

		dirName := setting["result.directory_name"].(string)
		if _, err := os.Stat(dirName); os.IsNotExist(err) {
			if useColors {
				fmt.Print(aurora.BrightRed("\rresult dir does not exist"))
				os.Mkdir(dirName, 0755)
				fmt.Println(aurora.BrightGreen("\rresult directory created"))
			} else {
				fmt.Print("\rresult dir does not exist")
				os.Mkdir(dirName, 0755)
				fmt.Println("\r" + dirName + "directory created")
			}
		}

		if resResult.Cont {
			os.WriteFile(dirName+"/"+fileName, []byte(string(d)+space+strings.TrimSpace(strings.Join(result, space))), 0644)
		} else {
			os.WriteFile(dirName+"/"+fileName, []byte(strings.TrimSpace(strings.Join(result, space))), 0644)
		}

		if useColors {
			fmt.Println(aurora.Green("\n\nDone! check ./result/" + title + ".txt"))
		} else {
			fmt.Println("\n\nDone! check ./" + setting["result.directory_name"].(string) + "/" + fileName)
		}

		if setting["cmd.exit_when_finish"].(bool) {
			break
		}
	}

	end()
}

func getEp(LOGINKEY string, page *Result, max float64, i int, ch chan Chan, tried int) {
	if tried == int(setting["cmd.max_try_per_episode"].(float64)+1) {
		ch <- Chan{"[" + page.Ep + "] " + page.Title + "\n\n\n\n\nError: 소설 정보를 불러올 수 없음", i}
		return
	}
	// req, _ := http.NewRequest("GET", "https://b-p.msub.kr/novelp/view/?id="+page.Link, nil)
	// req.Header.Set("Cookie", "LOGINKEY="+LOGINKEY+";")

	// resp, _ := http.DefaultClient.Do(req)

	// var res map[string]interface{}
	// json.NewDecoder(resp.Body).Decode(&res)

	// if res["err"] != nil {
	// 	fmt.Println(aurora.BrightYellow("\n\n  at EP." + strconv.Itoa(i+1) + ":"))
	// 	fmt.Print(aurora.BrightRed(res["err"].(string) + "\n\n"))
	// 	ch <- Chan{"[" + page.Ep + "] " + page.Title + "\n\n\n\n\n" + res["err"].(string), i}
	// 	return
	// }

	// if res["result"] == nil {
	// 	getEp(LOGINKEY, page, max, i, ch, tried+1)
	// 	return
	// }

	req, _ := http.NewRequest("POST", "https://novelpia.com/proc/viewer_data/"+page.Link, nil)
	req.Header.Set("Cookie", "LOGINKEY="+LOGINKEY+";")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Novelpia Downloader)")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var body respBody
	_ = json.NewDecoder(resp.Body).Decode(&body)

	arr := make([]string, len(body.S))
	for i := 0; i < len(body.S); i++ {
		arr[i] = body.S[i].Text
	}
	str := strings.Join(arr, "")

	str = regexp.MustCompile(`src="([^"]+)"`).ReplaceAllStringFunc(str, func(s string) string {
		if !setting["result.image_display"].(bool) {
			return ""
		}
		s = strings.Replace(s, "src=\"", "", 1)
		s = strings.Replace(s, "\"", "", 1)
		imgFormat := setting["result.image_format"].(string)
		if strings.HasPrefix(s, "http") {
			return ">" + strings.Replace(imgFormat, "${link}", s, 1) + "<"
		} else {
			return ">" + strings.Replace(imgFormat, "${link}", "http:"+s, 1) + "<"
		}
	})
	str = strings.ReplaceAll(str, "커버보기", "")
	str = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(str, "")

	str = strings.ReplaceAll(str, "&nbsp;", " ")
	str = strings.ReplaceAll(str, "&lt;", "<")
	str = strings.ReplaceAll(str, "&gt;", ">")
	str = strings.ReplaceAll(str, "&amp;", "&")
	str = strings.ReplaceAll(str, "&quot;", "\"")

	if str == "" {
		getEp(LOGINKEY, page, max, i, ch, tried+1)
		return
	}

	if page.Ep != "BONUS" {
		page.Ep += "화"
	}
	ch <- Chan{"[" + page.Ep + "] " + page.Title + "\n\n\n\n\n" + str, i}
}

func end() {
	if useColors {
		fmt.Print(aurora.BgWhite("\n\npress enter to exit...").Black())
	} else {
		fmt.Print("\n\npress enter to exit...")
	}
	fmt.Scanln()
	fmt.Scanln()
}
