# [Novelpia Downloader](https://github.com/taeseong14/N-down)

노벨피아 다운로더


### 공지

 - 설정 파일 추가..중


### 사용법

 * 1. 릴리즈([v0.1.1](https://github.com/taeseong14/N-down/releases/tag/v0.1.1))에서 "downloader.zip" 을 받는다
 * 2. 압축을 푼다
 * 3. downloader.exe 실행
 * 4. id(노벨피아 id만 가능, 구글 등 연동 ㄴㄴ) 와 password 입력
 * 5. bookId 입력(소설번호: 178143 등)

끝!
result/[소설명].txt 파일에서 확인하십숑

예제:
![예제](Example.png)
참고) "cmd.use_colors"를 사용할경우: Windows PowerShell을 이용해야 예제처럼 빤딱빤딱하게 나오빈다.

자세한건 설명서(zip파일의 README) ㄱㄱ

### 커스텀 설정

 - downloader.go를 실행하면 만들어지는 파일 settings.txt를 수정해서 다운로더/다운로드 파일을 입맛대로 쓰실수 있읍니다.
 - "key": value 형식이니 value만 바꿔가며 쓰십시오. (따옴표 등은 건들지 않는것을 추천함)

<details>
<summary style="font-size: 16px; font-weight: bold; cursor: pointer; margin-bottom: 20px;">Settings Document</summary>
<div>

> "account.auto_login": Boolean (true | false)
 - 자동 로그인 여부. 꺼져있다면(false) account 파일도 생성되지 않습니다.
 - Default: true

> "account.auto_login_file": String ("~~")
 - 로그인 정보 (id, pw)가 저장되는 파일명(혹은 루트)
 - Default: "account.txt"

> "account.default_mail": String ("@~~")
 - 로그인시 id가 @를 포함하지 않을 때 자동으로 뒤에 붙이는 문자열. (예제의 초록색 두번째줄 참고)
 - Default: "@gmail.com"

> "cmd.check_with_yn": Boolean (true | false)
 - bookId를 입력했을때 맞냐고 체크하는 부분 추가
 - Default: true

> "cmd.exit_when_finish": Boolean (true | false)
 - 한 소설 다운로드가 끝난 후 자동으로 exe파일 종료 여부
 - 한번 받을때 연속으로 많이 다운로드한다면 꺼놓는걸 추천
 - Default: true

> "cmd.max_try_per_episode": Number (123~)
 - 한 화당 최대 http 요청을 시도하는 횟수
 - 안건드는거 추천
 - Default: 5

> "cmd.use_colors": Boolean (true | false)
 - 자기가 예쁜 색은 보고싶지 않다, 혹은 powershell을 어떻게 연결하는지 모르겠다, 혹은 귀찮다 하시면 이거 끄시면 됨. <-]34 이런거 안나옵니다.
 - Default: true

> "result.image_display": Boolean (true | false)
 - 결과 텍스트파일에 이미지 링크 추가 여부
 - Default: true

> "result.image_format": String ("~~")
 - 이미지 링크 포멧 형식
 - query: ${link} 이미지 링크
 - Default: "[이미지: ${link}]" // -> [이미지: http://image.novelpia.com]

> "result.directory_name": String ("~~")
 - 결과 텍스트 파일들이 저장되는 폴더명, 혹은 루트
 - Default: "result" // -> ./result/~~.txt로 저장

> "result.file_name": String ("~~")
 - 결과 파일의 제목. txt파일로 저장되길 원하신다면 .txt를 붙이는걸 잊지 마세요.
 - query: ${title} 제목 | ${author} 작가
 - Default: "${title}.txt"

> "result.space_between_episodes": String ("~~")
 - 회차 사이에 추가하는 문자?
 - Default: "\n\n\n\n\n\n\n\n\n\n"

</div>
</details>

---

### 주의사항

 - exe 파일 다운로드 시, 실행 시에 경고창이 뜨지만 [위 downloader.go 파일](./downloader.go)을 빌드한것이니까 안심하셔도 됩니다. (정 머하면 golang 설치후 직접 go build downloader.go 쳐서 exe 만들기)
 - 동시에 여러개를 다운받으면 사용자에게 피해가 갈 수 있습니다.
 - 이 프로그램을 사용함으로써 생기는 피해는 이 프로그램의 제작자가 책임지지 않으며, 이 프로그램의 결과물의 저작권은 모두 원 저작자에게 있습니다. 무단 배포 및 전재가 금지됩니다.
 - 실행시 exe파일이 있는 폴더에 result폴더와 account.txt파일, settings.txt 파일을 생성하니 주의해주세요.


---


If there's any problem while downloading, progressing, or any additional function you want

[click here to make new issue](https://github.com/taeseong14/N-down/issues/new)

 + 기존꺼와 다른 문제/개선점이라면 새 이슈를 추가해주세요. (^)

 + 개발자 컨택: hutao@genshin.ai
