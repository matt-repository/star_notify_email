package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/gomail.v2"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	host := flag.String("host", "smtp.163.com", "send smtp server host")
	port := flag.Int("port", 465, "send smtp server port")
	sendMailbox := flag.String("sendMailbox", "", "send mailbox")
	receiveMailbox := flag.String("receiveMailbox", "", "receive mailbox")
	password := flag.String("password", "", "your mailbox Authorization code")
	cc := flag.String("cc", "", "Cc person")
	token := flag.String("token", "", "token")
	repository := os.Getenv("GITHUB_REPOSITORY")
	repoParts := strings.Split(repository, "/")
	user := repoParts[0]
	repo := repoParts[1]
	if *token == "" {
		fmt.Println("please input token")
		return
	}
	graphqlResponse := getGithubProjectInfo(*token, user, repo)
	lastUser := graphqlResponse.Data.Repository.Stargazers.Edges[0].Node
	subject := fmt.Sprintf("%s started", repository)
	content := fmt.Sprintf(`<div style="text-align: center;">   
			<h1>%s/%s</h1>
			<h2> 现在有 %d 个💕</h2> 
			<img style="max-width: 100%%; border-radius: 50%%" src="cid:avatar">   
			<div style="margin: 10px; font-size: x-large"> %s %s 给你点💕了</div>  
			<a href="%s" style="display: block; font-size: large">%s</a></div>
			`, user, repo, graphqlResponse.Data.Repository.StargazerCount, lastUser.Name, lastUser.Email, lastUser.URL, lastUser.URL)

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(*sendMailbox, ""))  //这个地方指定名称，会偶尔出现bug 是gomail 的bug
	m.SetHeader("To", m.FormatAddress(*receiveMailbox, "")) //主送
	if *cc != "" {
		m.SetHeader("Cc", *cc) //抄送
	}
	m.SetHeader("Subject", subject) //标题
	m.SetBody("text/html", content) // 发送html格式邮件，发送的内容
	d := gomail.NewDialer(*host, *port, *sendMailbox, *password)

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("send fail err:%s", err.Error())
	} else {
		fmt.Println("send success")
	}
}

// 获取此项目的一些信息，
func getGithubProjectInfo(token, user, repo string) (graphQLResponse GraphqlResponse) {
	requestBody, err := json.Marshal(GraphqlRequest{
		Query: fmt.Sprintf(`{
			repository(name: "%s", owner: "%s") {
				stargazerCount
				stargazers(last: 1) {
					edges {
						node {
							name
							url
							avatarUrl
							email
						}
					}
				}
			}
		}`, repo, user),
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	client := &http.Client{Timeout: time.Second * 10}
	request, err := http.NewRequest("POST", "https://api.github.com/graphql", strings.NewReader(string(requestBody)))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	request.Header.Add("User-Agent", "Go")
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&graphQLResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if graphQLResponse.Data.Repository.StargazerCount > 0 {
		lastUser := graphQLResponse.Data.Repository.Stargazers.Edges[0].Node
		fmt.Printf("Last star by %s (%s <%s>)\n", lastUser.Name, lastUser.Email, lastUser.URL)
		fmt.Printf("%d total stars\n", graphQLResponse.Data.Repository.StargazerCount)
	} else {
		fmt.Println("No stars yet")
	}
	return graphQLResponse
}
