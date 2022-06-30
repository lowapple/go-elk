# go-elk

go.elk는 notion으로 작성된 페이지를 다른 웹 서비스에 올려둘 수 있도록 static site를 생성하는 클라이언트 라이브러리 입니다.

Go 버전 1.18로 작성되었습니다. 현재 하나의 페이지만 지원합니다. 이후 페이지 내에 포함된 서브 페이지도 추적하여 Notion Webpage와 동일하게 만드는것이 목표입니다. 
기본적으로 테스트를 진행하고, 올리고 있지만 항상 테스트하고 있지는 않습니다.

## 설치
별도로 chromedriver 설치가 필요합니다. 환경에 따라 각자 설치해 주세요
* [chromedriver](https://chromedriver.chromium.org/downloads)

```sh
go get github/lowapple/go-elk@latest
```
또는 Golang이 설치가 되어있으면 소스코드를 다운로드하고, 소스코드에서 다음과 같이 실행할 수 있습니다.
```sh
go run main.go ...
```

## 시작하기

```
elk [OPTIONS] URL 
```
```
elk "https://lowapple.notion.site/e1db500567ca46bcbb59ed2f575325e4"
```

## 옵션
```
  -o string
    	notion output static page dir path (default "./dist")
```

## 라이선스

