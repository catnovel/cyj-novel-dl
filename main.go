package main

import (
	"flag"
	"fmt"
	"github.com/catnovelapi/cyj"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"path"
	"sync"
)

type Person struct {
	bookId      string
	bookName    string
	saveDir     string
	loginToken  string
	account     string
	chaptersDir string
	client      *cyj.Client
}

func init() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("打开日志文件失败：", err.Error())
		return
	}
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func (person *Person) download(bookInfo gjson.Result) {
	person.bookId = bookInfo.Get("bookId").String()
	person.bookName = bookInfo.Get("bookName").String()
	if person.bookName == "" {
		log.Fatalf("bookId:%s,获取书籍信息失败", person.bookId)
	}
	person.chaptersDir = person.newFile(path.Join(person.saveDir, person.bookName, "chapters"))

	var wg sync.WaitGroup
	for _, chapter := range person.client.NewGetCatalogByBookIDApi(person.bookId) {
		wg.Add(1)
		go func(chapter gjson.Result) {
			defer wg.Done()
			var downloadSuccess = true
			chapterId := chapter.Get("chapterId").String()
			if chapter.Get("isFee").String() == "1" && chapter.Get("isBuy").String() != "1" {
				downloadSuccess = false
				log.Println("chapterId:", chapterId, "warning:该章节需要付费")
			}
			if person.exists(path.Join(person.chaptersDir, chapterId+".txt")) {
				downloadSuccess = false
				log.Println("chapterId:", chapterId, "warning:该章节已经下载过了")
			}
			if downloadSuccess {
				content := person.client.NewGetContentByBookIdAndChapterIdApi(chapterId, person.bookId)
				if content != "" {
					fmt.Println(chapter.Get("chapterName").String(), "下载成功")
					person.writeFile(path.Join(person.chaptersDir, chapterId+".txt"), content)
				}
			}
		}(chapter)
	}
	wg.Wait()
}
func (person *Person) outFile() {
	p := path.Join(path.Dir(person.chaptersDir), person.bookName+".txt")
	if err := os.Truncate(p, 0); err != nil {
		fmt.Println("清空文件失败:", err)
	}
	file, err := os.OpenFile(p, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()
	for _, chapter := range person.client.NewGetCatalogByBookIDApi(person.bookId) {
		var chapterId = chapter.Get("chapterId").String()
		if person.exists(path.Join(person.chaptersDir, chapterId+".txt")) {
			content, err := os.ReadFile(path.Join(person.chaptersDir, chapterId+".txt"))
			if err != nil {
				log.Println("读取文件失败:", err)
			} else {
				_, _ = file.Write([]byte(chapter.Get("chapterName").String() + "\n" + string(content) + "\n\n"))
			}
		}
	}
}

func (person *Person) newFile(name string) string {
	_, err := os.Stat(name)
	if err != nil {
		err = os.MkdirAll(name, os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败:", err)
		}
	}
	return name
}

func (person *Person) exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

func (person *Person) writeFile(name string, content string) {
	if err := os.WriteFile(name, []byte(content), 0644); err != nil {
		fmt.Println("写入文件失败:", err)
	}
}

func main() {
	var person Person
	flag.StringVar(&person.bookId, "d", "", "add bookid to download book")
	flag.StringVar(&person.saveDir, "o", "books", "save file name")
	flag.StringVar(&person.bookName, "s", "", "Search book name to download book")
	flag.StringVar(&person.account, "a", "", "set account")
	flag.StringVar(&person.loginToken, "t", "", "set login token")
	flag.Parse()

	person.client = cyj.NewCiyuanjiClient(cyj.Token(person.loginToken))
	if person.bookId != "" {
		person.download(person.client.GetBookInfoApi(person.bookId).Get("data.book"))
		person.outFile()
	} else if person.account != "" {
		res := person.client.GetPhoneCodeByPhoneNumberApi(person.account)
		if res.Get("code").String() == "200" {
			fmt.Printf("验证码已发送至%s,请注意查收\n请输入验证码：\n", person.account)
			var phoneCode string
			fmt.Scanln(&phoneCode)
			response := person.client.GetLoginByPhoneNumberAndPhoneCodeApi(person.account, phoneCode)
			if response.Get("code").String() == "200" {
				fmt.Println("次元姬账号登入成功!")
				fmt.Println("token:", response.Get("data.userInfo.token").String())
			} else {
				fmt.Println("登录失败 msg:", response.Get("msg").String())
			}
		} else {
			fmt.Println("验证码发送失败:", res.Get("msg").String())
		}
	} else if person.bookName != "" {
		resultArray := person.client.GetSearchByKeywordApi(person.bookName, "1").Get("data.esBookList").Array()
		for i, result := range resultArray {
			fmt.Println("INDEX", i, "\t\tbookName", result.Get("bookName").String())
		}
		var bookIndex int
		for {
			fmt.Printf("请输入要下载的书籍序号:")
			fmt.Scanln(&bookIndex)
			if bookIndex < len(resultArray) {
				person.download(resultArray[bookIndex])
				person.outFile()
			}
		}
	} else {
		fmt.Println("请输入参数,使用 -h 查看帮助")
	}
}
