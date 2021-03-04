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
	/*
	(function(){
	    var _id    = 309847;
	    var isHome = $('a[title="我的卡包"]').html();
	    var _temp  = isHome ? "" : $('span[title*="全部文件"]')[0].title.slice(4);
	    var _name  = '负重前行.mkv'; // 这里 '' 里面的内容改成需要下载的文件的名称
	    var _path  = encodeURIComponent(_temp + '/' + _name);
	    var _link  = 'https://pcs.baidu.com/rest/2.0/pcs/file?method=download&app_id='+_id+'&path='+_path;
	    console.log('下载地址为：');
	    console.log('%c%s','color:#00ff00;background-color:#000000;',_link);
	})();
	 */
	//负重前行
	text := "https://qdall01.baidupcs.com/file/4bc5bd9875da891152a66575c9413862?bkt=en-2e2b5030dd6ff037fdaa7a8e2a932ca0812ee4714ad9b1ad840ab735ce4cef52031a2ca0c088d534&fid=1996670861-309847-1097323969048963&time=1614694578&sign=FDTAXUGERLQlBHSKfWaqir-DCb740ccc5511e5e8fedcff06b081203-IQHo03o13JbxGMLafRNisi0o3dw%3D&to=92&size=394807668&sta_dx=394807668&sta_cs=7407&sta_ft=mkv&sta_ct=7&sta_mt=7&fm2=MH%2CYangquan%2CAnywhere%2C%2Cshanghai%2Cct&ctime=1484727719&mtime=1542705565&resv0=-1&resv1=0&resv2=rlim&resv3=5&resv4=394807668&vuk=1996670861&iv=0&htype=&randtype=em&newver=1&newfm=1&secfm=1&flow_ver=3&pkey=en-679954cd647643e89b98547c31a1ca95f7d8b0e335494b306d23e25592864be98b0c866b2692f9b9&sl=76480590&expires=8h&rt=pr&r=130657479&mlogid=1413092100370096985&vbdid=3403188840&fin=%E8%B4%9F%E9%87%8D%E5%89%8D%E8%A1%8C.mkv&fn=%E8%B4%9F%E9%87%8D%E5%89%8D%E8%A1%8C.mkv&rtype=1&dp-logid=1413092100370096985&dp-callid=0.1.1&hps=1&tsl=80&csl=80&fsl=-1&csign=dVYKgEit045y%2FYZnjUaT3WXZHfA%3D&so=0&ut=6&uter=4&serv=1&uc=3853463274&ti=9feb8afe5adad8c3fc50915f2b69a5a6cdc3f80ea34cfa97&hflag=30&from_type=0&adg=c_edc0108e9fa1ea2bf75676e893bdf053&reqlabel=309847_d_542374c7794461b783501a9d00e48223_-1_18e87ecb8b58d3c5a256ffa3f09e85da&by=themis"
	down, err := NewDownloader(text,
		SetTimeout(222601),
		SetDownloadRoutine(20),
		SetStrategyWait(true),
		SetTryTimes(200),
	)
	if err != nil {
		log.Println(err)
	} else {
		err = down.SaveFile(context.Background())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(down.option.SaveName)
		}
	}
}

func TestTurnCate(t *testing.T) {
	file, err := os.Create("./112.txt")
	if err != nil {
		t.Fatal(err)
	}
	err = file.Truncate(1024*1024)
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.WriteAt([]byte("test"), 222)
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.WriteAt([]byte("t211111111112222222222222222222222222222222222222"), 1000)
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.WriteAt([]byte("zmisgod"), 1024)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("ok")
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
