package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
	"strconv"
)

type temp struct{
	downloaded int64
	all int64
	time int
}

func (tmp *temp) Write(b []byte) (int, error){
	n := len(b)
	tmp.downloaded += int64(n)
	return n, nil
}

func print(tmp *temp) {
	fmt.Println(tmp.downloaded/1024, "kilobytes downloaded from", tmp.all/1024, "kilobytes")
}

 func countingBytes(tmp *temp){
	for {
		time.Sleep(time.Second)
		tmp.time++
		print(tmp)
	}
}

func downloadFile(fileURL string) error{

	resp, err := http.Get(fileURL)
	if err != nil{
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path.Base(resp.Request.URL.String()))
	if err != nil{
		return err
	}

	contentLength := resp.Header.Get("content-length")
	length, err := strconv.Atoi(contentLength)
	if err != nil{
		return err
	}

	counter := &temp{}
	counter.all = int64(length)
	go countingBytes(counter)
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil{
		return err
	}

	fmt.Println(counter.all/1024, "kilobytes downloaded from", counter.all/1024, "kilobytes")

	out.Close()

	return nil
}

func main(){

	var fileURL string
	fmt.Println("Your url: ")
	fmt.Scanf("%s", &fileURL)

	fmt.Println()

	err := downloadFile(fileURL)
	if err != nil{
		panic(err)
	}
	
	fmt.Println()
	fmt.Println("Done!")
}