package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
)

var rodlauncher *launcher.Launcher

func main() {
	v2ex()
	pojie_52()

}

func newBrowser(headless bool) *rod.Browser {
	//設定啓動器
	//false: 显示浏览器
	//true: 使用无头浏览器
	rodlauncher = launcher.New().
		//Proxy("192.168.2.114:7890"). // set flag "--proxy-server=127.0.0.1:8080"
		Proxy("socks5://127.0.0.1:10812"). // set flag "--proxy-server=127.0.0.1:8080"
		Delete("use-mock-keychain").       // delete flag "--use-mock-keychain"
		Headless(headless).
		UserDataDir("tmp/user").
		Set("mute-audio").
		Set("default-browser-check").
		Set("disable-gpu").
		Set("disable-web-security").
		Set("no-sandbox")

		//關閉無頭模式，顯示瀏覽器窗體
		//Delete("headless")

	//debug url
	launchers := rodlauncher.MustLaunch()
	fmt.Printf("debug url: %s\n", launchers)

	//連接到瀏覽器
	return rod.New().ControlURL(launchers).MustConnect()
}

func newBrowser_52pojie(headless bool) *rod.Browser {
	//設定啓動器
	//false: 显示浏览器
	//true: 使用无头浏览器
	rodlauncher = launcher.New().
		//Proxy("192.168.2.114:7890"). // set flag "--proxy-server=127.0.0.1:8080"
		Delete("use-mock-keychain"). // delete flag "--use-mock-keychain"
		Headless(headless).
		UserDataDir("tmp/user").
		Set("mute-audio").
		Set("default-browser-check").
		Set("disable-gpu").
		Set("disable-web-security").
		Set("no-sandbox")

	//關閉無頭模式，顯示瀏覽器窗體
	//Delete("headless")

	//debug url
	launchers := rodlauncher.MustLaunch()
	fmt.Printf("debug url: %s\n", launchers)

	//連接到瀏覽器
	return rod.New().ControlURL(launchers).MustConnect()
}

func v2ex() {
	//設定設備參數
	screen := devices.Device{
		Title:        "Laptop with MDPI screen",
		Capabilities: []string{"touch", "mobile"},
		UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
		Screen: devices.Screen{
			DevicePixelRatio: 1,
			Horizontal: devices.ScreenSize{
				Width:  1280,
				Height: 720,
			},
		},
	}

	//false: 显示浏览器
	//true: 使用无头浏览器
	browser := newBrowser(false)
	browser.DefaultDevice(screen).MustPage("")
	defer browser.Close()

	page := browser.Timeout(time.Minute).MustPage("https://www.v2ex.com/").MustWaitLoad()
	if !page.MustHasR("a", "登出|Sign Out") || !page.MustHasR("a", `领取今日的登录奖励`) {
		fmt.Println("[v2ex] 已经签过到了或者第一次运行请登陆系统！30 秒后退出!")
		time.Sleep(30 * time.Second)
		return
	}

	page.Race().ElementR("a", `领取今日的登录奖励`).MustHandle(func(el *rod.Element) {
		el.MustClick()

		page.MustElementR("input", "领取 X 铜币").MustClick()
		page.MustElementR(".message", "已成功领取每日登录奖励")
		log.Println("签到成功")
	}).Element(`.balance_area`).MustHandle(func(el *rod.Element) {
		log.Println("已经签过到了")
	}).MustDo()
}

func pojie_52() {
	//設定設備參數
	screen := devices.Device{
		Title:        "Laptop with MDPI screen",
		Capabilities: []string{"touch", "mobile"},
		UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
		Screen: devices.Screen{
			DevicePixelRatio: 1,
			Horizontal: devices.ScreenSize{
				Width:  1280,
				Height: 720,
			},
		},
	}

	//false: 显示浏览器
	//true: 使用无头浏览器
	//browser := newBrowser(false)
	browser := newBrowser_52pojie(false)
	browser.DefaultDevice(screen).MustPage("")
	defer browser.Close()

	page := browser.Timeout(time.Minute).MustPage("https://www.52pojie.cn/").MustWaitLoad()
	//println(page.MustHasR("img", `https://www.52pojie.cn/static/image/common/wbs.png`))
	//println(page.MustHasR("img", "https://www.52pojie.cn/static/image/common/qds.png"))
	// 获取一个元素的子元素
	projects := page.MustElements("img")
	for _, project := range projects {
		wbs := project.MustProperty("src").String()
		//fmt.Println(strings.Contains(wbs, "wbs.png"))
		if strings.Contains(wbs, "wbs.png") {
			fmt.Println("[52pojie] 已经签过到了或者第一次运行请登陆系统! 30 秒后退出!")
			time.Sleep(30 * time.Second)
			return
		}
	}

	fmt.Println("[52pojie] 开始签到!")
	page.MustElementX(`//*[@id="um"]/p[2]/a[1]/img`).MustClick()
	time.Sleep(10 * time.Second)
	fmt.Println("[52pojie] 执行完成")
}
