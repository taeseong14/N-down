# [Novelpia Downloader](https://github.com/taeseong14/N-down)

노벨피아 다운로더

### 공지사항
과도한 요청으로 서버가 ip밴 당했다 풀림. 요청하는데 딜레이 넣어둔 상태고 이슈#18으로 서버 부담 줄일 예정임


> 개발노트

 - 설정 파일 추가예정(계정, 삽화 url 유무, 화수(1화, 2화) 유무 등등)
 - 이슈#18 다음주까지 고치기 !important


> 사용법

 * 1. 릴리즈([v0.0.8](https://github.com/taeseong14/N-down/releases/tag/v0.0.8))에서 "downloader.zip" 을 받는다
 * 2. 압축을 푼다
 * 3. id(노벨피아 only, 구글 등 연동 ㄴㄴ) 와 password 입력
 * 4. bookId(소설번호: 178143 등)
 * 5. 끝날때까지 기다린다

끝!
result/[소설명].txt 파일에서 확인하십숑

예제:
![예제](Example.png)
참고) windows powershell을 이용해야 빤딱빤딱하게 나오빈다.

자세한건 설명서(zip파일의 README) ㄱㄱ



> 주의사항

 - exe 파일 다운로드 시, 실행 시에 경고창이 뜨지만 [위 파일](./downloader.go)을 빌드한것이니까 안심하셔도 됨 (정 머하면 golang 설치후 직접 go build downloader.go 치기)
 - 되도록이면 요청을 한꺼번에 보내지 말아주세요 (임시로 막아둔상태). 지금은 제 사설섭을 거쳐가지만 이슈#18 해결하면 사용자의 ip로 요청을 보내기때문에 님 ip만 밴먹는거임 ㅇㅇ


---


If there's any problem while downloading, progressing, or any additional function you want

[click here to make new issue](https://github.com/taeseong14/N-down/issues/new)

+ 기존꺼와 다른 문제/개선점이라면 새 이슈를 추가해주세요. (^)





 + 개발자 컨택: hutao@genshin.ai
