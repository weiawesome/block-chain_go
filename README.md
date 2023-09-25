# api:
    建立連線成為礦工 接點
    上傳交易 接點
    查詢交易 接點

# 區塊資料庫
    RocksDB:
    last_block:最後區塊哈希值
    checked:最後被6區塊確認區塊
    last:[] 未確認區塊

# goroutine:

## * 接收交易->傳遞給channel A
## * 接收區塊->丟到channel D

## 1.When channel D收到資訊->進行內容確認
    if 前一區塊並不在:
        拒絕區塊
    透過該區塊值 不斷往前尋找7個
    if 在7個以內找到checked:
        更新UTXO 加入last[]
    //chains=[checked,1,2,3,4,5,6,7]
    if 在第7個找到checked:
        //刪除支線
        歷遍 last[] if 區塊(前一去塊哈希值=checked) 刪除更新UTXO
        歷遍  last[] if 區塊(前一區塊不存在) 刪除跟新UTXO
        //更新UTXO 最新線
        歷遍 chains 更新UTXO
        更新 checked 轉換成 1 last 轉換成 7
        傳遞給channel C->廣播

## 2.When channelA 接收到資訊:
        進行驗算合法性->進行深度確認->廣播
        將資訊丟進pq
        1. pq依序pop
        2. 建立flag 區塊中交易採用set型態 flag變true唯有新資訊出現
        3. flag為true 先傳遞stopChannel 則將區塊丟進channelB
        4. 將區塊內交易重新丟回pq
    when channelC 接收到資訊:
        1. 刷新區塊
        2. pq依序pop
        3. 若是value存在資訊 直接丟掉 否則新增進入交易
        4.  先傳遞stopChannel 將區塊丟進channelB
        5. 將區塊內交易重新丟回pq

## 3.When SuccessChannel 收到資訊:
        1. 傳進 StopChannel;
        2. 傳進ChannelD
        3. 將資訊廣播

## 4. When ChannelB 收到資訊:
        1. 根據nonce範圍
        When nonce符合:
            傳送SuccessChannel
       Default:
         continue

# Channel:
    channelA: 傳遞以驗證交易
    channelB: 傳遞區塊
    channelC: 傳遞驗證完成區塊
    channelD: 接收要驗證區塊
    SuccessChannel: 回傳成功
    StopChannel: 強迫礦工們停止挖礦

# 初始化:
    設定 公鑰
    接收2種資料庫資料
    根據 平行開發數 啟動4 goroutine數量
    5個goroutine啟動

# 資料庫:
    1. 區塊 鍵值對資料庫
        鍵:區塊哈希值
        值:區塊資訊
    2. UTXO 資料庫
        鍵:交易哈希
        索引:
        值:
        {
            amount:
            接收者地址:
            spent:
        }

# 交易資訊:
    公鑰:
    輸入:[UTXO 地址們]
    輸出:[地址們 與 金額]
    手續費:
    簽名:
    交易哈希:




package main

import (
"fmt"
"sync"
"time"
)

func mine(index int, stopChannel chan bool, blockChannel chan int, SuccessChannel chan int, wg *sync.WaitGroup) {
var value int
for {
select {
case val := <-blockChannel:
if value != val {
value = val
for i := 0; i < 1000; i++ {
select {
case val = <-blockChannel:
i = 0
continue

					default:
						if i == 10 {
							SuccessChannel <- index
							break
						}
						fmt.Println(index, " ", i)
						time.Sleep(1 * time.Second)
					}
				}
			} else {
				blockChannel <- val
			}
		default:
			continue
		}
	}

}

func main() {
// 获取当前时间
currentTime := time.Now()

	// 获取纳秒级的时间戳
	// 注意：UnixNano() 返回的是 int64 类型的纳秒数
	timestampNano := currentTime.UnixNano()

	// 获取秒级的时间戳
	// 注意：Unix() 返回的是 int64 类型的秒数
	timestampSec := currentTime.Unix()

	fmt.Printf("当前时间: %s\n", currentTime)
	fmt.Printf("纳秒级时间戳: %d\n", timestampNano)
	fmt.Printf("秒级时间戳: %d\n", timestampSec)
}

//func main() {
//	//stopChannel := make(chan bool)
//	//blockChannel := make(chan int, 10)
//	//successChannel := make(chan int, 10)
//	//var wg sync.WaitGroup // 創建 WaitGroup
//	//
//	//for i := 0; i < 5; i++ {
//	//	index := i
//	//	wg.Add(1)
//	//	go mine(index, stopChannel, blockChannel, successChannel, &wg)
//	//}
//	//for i := 0; i < 5; i++ {
//	//	blockChannel <- 2
//	//}
//	//for {
//	//	select {
//	//	case val := <-successChannel:
//	//		fmt.Println("GAGAG", val)
//	//		for i := 0; i < 5; i++ {
//	//			stopChannel <- true
//	//		}
//	//	default:
//	//		continue
//	//	}
//	//}
//}
