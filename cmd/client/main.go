package main

import (
	"bufio"
	"fmt"
	"github.com/wan6sta/go-url/internal/config"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	cfg := config.NewConfig()
	data := url.Values{}
	fmt.Println("Введите длинный URL")

	reader := bufio.NewReader(os.Stdin)
	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	long = strings.TrimSuffix(long, "\n")

	data.Set("url", long)
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s", cfg.Address), strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	fmt.Println("Статус-код ", response.Status)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
