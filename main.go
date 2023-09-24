package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/peanut996/CloudflareWarpSpeedTest/task"
	"github.com/peanut996/CloudflareWarpSpeedTest/utils"
)

func init() {
	var printVersion bool
	var help = `
gogloo \n` + `
测试 gogloo 所有 IP 的延迟和速度，获取最低延迟和端口

参数：
    -n 1
        延迟测速线程；越多延迟测速越快，性能弱的设备 (如路由器) 请勿太高；(默认 200 最多 1000)
    -t 1
        延迟测速次数；单个 IP 延迟测速的次数；(默认 10 次)
    -q 
        快速模式； 随机扫描5000的地址测速结果；默认打开，[-q=false] 关闭快速模式
    -tl 3000
        平均延迟上限；只输出低于指定平均延迟的 IP，各上下限条件可搭配使用；(默认 300 ms)
    -tll 0
        平均延迟下限；只输出高于指定平均延迟的 IP；(默认 0 ms)
    -tlr 1
        丢包几率上限；只输出低于/等于指定丢包率的 IP，范围 0.00~1.00，0 过滤掉任何丢包的 IP；(默认 1.00)
    -sl 5
        下载速度下限；只输出高于指定下载速度的 IP，凑够指定数量 [-dn] 才会停止测速；(默认 0.00 MB/s)

    -p 0
        显示结果数量；测速后直接显示指定数量的结果，为 0 时不显示结果直接退出；(默认 10 个)
    -f ip.txt
        IP段数据文件；如路径含有空格请加上引号；
    -ip 1.1.1.1,2.2.2.2/24,2606:4700::/32
        指定IP段数据；直接通过参数指定要测速的 IP 段数据，英文逗号分隔；(默认 空)
    -o 1.tmp
        写入结果文件；如路径含有空格请加上引号；值为空时不写入文件 [-o ""]；(默认 result.csv)
    -pri 私钥
        指定你的wireguard私钥
    -pub 私钥
        指定你的wireguard公钥, 默认为warp的公钥
    -full
        测速全部的端口；对 IP 段中的每个 IP 全部端口进行测速
    -h
        打印帮助说明
`
	var minDelay, maxDelay int
	var maxLossRate float64
	flag.IntVar(&task.Routines, "n", 1, "延迟测速线程")
	flag.IntVar(&task.PingTimes, "t", 1, "延迟测速次数")

	flag.IntVar(&maxDelay, "tl", 3000, "平均延迟上限")
	flag.IntVar(&minDelay, "tll", 0, "平均延迟下限")
	flag.Float64Var(&maxLossRate, "tlr", 1, "丢包几率上限")

	flag.BoolVar(&task.ScanAllPort, "full", false, "扫描全部端口")
	flag.BoolVar(&task.QuickMode, "q", true, "快速模式，仅随机扫描5000个IP的结果")
	flag.IntVar(&utils.PrintNum, "p", 0, "显示结果数量")
	flag.StringVar(&task.IPFile, "f", "", "IP段数据文件")
	flag.StringVar(&task.IPText, "ip", "", "指定IP段数据")
	flag.StringVar(&utils.Output, "o", "1.tmp", "输出结果文件")
	flag.StringVar(&task.PrivateKey, "pri", "", "指定private key")
	flag.StringVar(&task.PrivateKey, "pub", "", "指定public key")

	flag.BoolVar(&printVersion, "v", false, "打印程序版本")
	flag.Usage = func() { fmt.Print(help) }
	flag.Parse()

	utils.InputMaxDelay = time.Duration(maxDelay) * time.Millisecond
	utils.InputMinDelay = time.Duration(minDelay) * time.Millisecond
	utils.InputMaxLossRate = float32(maxLossRate)
}

func main() {
	task.InitRandSeed()
	task.InitHandshakePacket()

	fmt.Printf("gogloo\n\n")

	pingData := task.NewWarping().Run().FilterDelay().FilterLossRate()
	utils.ExportCsv(pingData)
	pingData.Print()
}
