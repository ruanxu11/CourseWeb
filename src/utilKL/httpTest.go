package utilKL

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func HttpGet() string {
	resp, err := http.Get("http://localhost:2333/login?ID=3140102431&Password=3140102431&Type=student")
	if err != nil {
		// handle error
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return ""
	}
	fmt.Println(string(body), err)
	return string(body)
}

func HttpPost() {
	v := url.Values{}
	v.Set("ID", "3140102431")
	v.Set("Password", "3140102431")
	v.Set("Type", "student")
	body := ioutil.NopCloser(strings.NewReader(v.Encode())) //把form数据编下码
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:2333/login", body)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded") //这个一定要加，不加form的值post不过去，被坑了两小时
	fmt.Printf("%+v\n\n\n", req)                                        //看下发送的结构

	resp, err := client.Do(req) //发送
	defer resp.Body.Close()     //一定要关闭resp.Body
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data), err)
}

func HttpGetPost() {
	v := url.Values{}
	v.Set("ID", "3140102433")
	v.Set("Password", "3140102431")
	v.Set("Type", "student")
	body := ioutil.NopCloser(strings.NewReader(v.Encode())) //把form数据编下码
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:2333/login?ID=3140102431&ID=3140102432&Password=3140102431&Type=student", body)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded") //这个一定要加，不加form的值post不过去，被坑了两小时
	fmt.Printf("%+v\n\n\n", req)                                        //看下发送的结构

	resp, err := client.Do(req) //发送
	defer resp.Body.Close()     //一定要关闭resp.Body
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data), err)
}
