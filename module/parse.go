package module

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"../prods"
)

type FiveGStream struct {
	OutputChn chan []byte
	Close     chan int
	Wg        sync.WaitGroup
}

//go run main.go ./config.ini prod를 Generate 함수 등을 통해 주여진 파일 (또는 경로로 부터) 데이터를 수집해 input channel로 보냄
//producer 객채는 데이터를 수신하고 이를 kafka 서버에 보냄

////// data parse
func Generate(path string, input chan prods.Input, imsiIndex int, bearer bool) {
	if _, err := os.Stat(path); err != nil {
		// path/to/whatever exists
		panic(fmt.Sprintf("File [%v] not exists!", path))
		panic(err)
	}

	f := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			f, err := ioutil.ReadFile(path)

			if err != nil {
				panic(err)
			}

			lines := strings.Split(string(f), "\r\n")

			for _, line := range lines {
				//println(len(f))
				if len(line) == 0 {
					continue
				}

				//cols := strings.Split(line, ",")
				cols := strings.Split(line, "\\") //구분자 parsing

				//println(cols)
				Imsi := cols[imsiIndex]

				//
				i := prods.Input{

					Data: strings.Join(cols, ","),
					Imsi: Imsi,
				}
				input <- i
			}
		}
		return err
	}

	err := filepath.Walk(path, f)
	if err != nil {
		panic(err)
	}
}

//// 지정된 Path 의 xDR 차례대로 처리
func GeneratePath(fpath string, input chan prods.Input, imsiIndex int, bearer bool) {
	if _, err := os.Stat(fpath); err != nil {
		// path/to/whatever exists

		panic(fmt.Sprintf("File [%v] not exists!", fpath))
		panic(err)
	}
	files, _ := filepath.Glob(fpath + "/*")
	println(files)

	for _, path := range files {
		println(path)
		f := func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				f, err := ioutil.ReadFile(path)

				if err != nil {
					panic(err)
				}

				lines := strings.Split(string(f), "\r\n")

				for _, line := range lines {
					//println(len(f))
					if len(line) == 0 {
						continue
					}

					//cols := strings.Split(line, ",")
					cols := strings.Split(line, "\\\\") //구분자 parsing

					//control length중 이상한거 있으면 넘어간다
					if len(cols) != 49 {
						println(cols)
						continue
					}
					Imsi := cols[imsiIndex]

					//delete 모드일때는 스킵
					if string(path[len(path)-1]) != "R" {
						//println(path)

						DirIndex := 36
						Dir := cols[DirIndex]

						if Dir == "36" {
							//println(Dir)
							continue
						}

					}

					if Imsi == "" {
						Imsi = "1"
					}

					i := prods.Input{

						Data: strings.Join(cols, ","),
						Imsi: Imsi,
					}

					input <- i
				}

			}
			return err
		}

		err := filepath.Walk(path, f)
		if err != nil {
			panic(err)
		}

		println("finished:", path)

	}

}
