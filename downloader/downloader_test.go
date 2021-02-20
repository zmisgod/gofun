package downloader

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"testing"
)

func TestNewDownloader(t *testing.T) {
	text := "https://www.imax.com/download/file/fid/16840"
	down, err := NewDownloader(text,
		SetTimeout(60),
		SetDownloadRoutine(6),
	)
	if err != nil {
		log.Println(err)
	} else {
		err = down.SaveFile(context.Background())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(down.SaveName)
		}
	}
}

func combineFile() error {
	filename := "test.txt"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	for i := 0; i < 6; i++ {
		name := fmt.Sprintf("test%d.txt", i)
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		w := bufio.NewWriter(file)
		_, _ = io.Copy(w, bufio.NewReader(f))
		_ = w.Flush()
		_ = f.Close()
	}
	return nil
}

func TestCreateFile(t *testing.T) {
	//var wg sync.WaitGroup
	//for i := 0; i < 6; i++ {
	//	wg.Add(1)
	//	go func(i int) {
	//		defer wg.Done()
	//		if err := output(i, 1); err != nil {
	//			fmt.Println(err)
	//		}
	//	}(i)
	//}
	//wg.Wait()

	for i := 0; i < 6; i++ {
		name := fmt.Sprintf("test%d.txt", i)
		readF(name)
	}
	fmt.Println(combineFile())
	fmt.Println("-----")
	readF("test.txt")
}

func readF(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	w := bufio.NewReader(file)
	p := ReaderSource(w)
	fmt.Println(fileName)
	for v := range p {
		fmt.Println(v)
	}
}

func ReaderSource(reader io.Reader) <-chan int {
	out := make(chan int)
	go func() {
		buffer := make([]byte, 8)
		for {
			n, err := reader.Read(buffer)
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil {
				break
			}
		}
		close(out)
	}()
	return out
}

func WriteSource(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		_, _ = writer.Write(buffer)
	}
}

func output(i int, count int) error {
	file, err := os.Create(fmt.Sprintf("test%d.txt", i))
	if err != nil {
		return err
	}
	defer file.Close()
	p := RandomSource(count)
	w := bufio.NewWriter(file)
	WriteSource(w, p)
	_ = w.Flush()
	return nil
}

func RandomSource(count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}
