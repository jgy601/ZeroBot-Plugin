// Package tupian 图片获取集合
package tupian

import (
	"github.com/FloatTech/AnimeAPI/bilibili"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	referer  = "https://weibo.com/"
	shouer   = "https://iw233.cn/api.php?sort=cat&referer"
	baisi    = "http://aikohfiosehgairl.fgimax2.fgnwctvip.com/uyfvnuvhgbuiesbrghiuudvbfkllsgdhngvbhsdfklbghdfsjksdhnvfgkhdfkslgvhhrjkdshgnverhbgkrthbklg.php?sort=ergbskjhebrgkjlhkerjsbkbregsbg"
	heisi    = "http://aikohfiosehgairl.fgimax2.fgnwctvip.com/uyfvnuvhgbuiesbrghiuudvbfkllsgdhngvbhsdfklbghdfsjksdhnvfgkhdfkslgvhhrjkdshgnverhbgkrthbklg.php?sort=rsetbgsekbjlghelkrabvfgheiv"
	siwa     = "http://aikohfiosehgairl.fgimax2.fgnwctvip.com/uyfvnuvhgbuiesbrghiuudvbfkllsgdhngvbhsdfklbghdfsjksdhnvfgkhdfkslgvhhrjkdshgnverhbgkrthbklg.php?sort=dsrgvkbaergfvyagvbkjavfwe"
	bizhi    = "https://iw233.cn/api.php?sort=iw233"
	baimao   = "https://iw233.cn/api.php?sort=yin"
	xing     = "https://iw233.cn/api.php?sort=xing"
	sese     = "http://aikohfiosehgairl.fgimax2.fgnwctvip.com/uyfvnuvhgbuiesbrghiuudvbfkllsgdhngvbhsdfklbghdfsjksdhnvfgkhdfkslgvhhrjkdshgnverhbgkrthbklg.php?sort=qwuydcuqwgbvwgqefvbwgueahvbfkbegh"
	biaoqing = "https://iw233.cn/api.php?sort=img"
	cos      = "http://aikohfiosehgairl.fgimax2.fgnwctvip.com/uyfvnuvhgbuiesbrghiuudvbfkllsgdhngvbhsdfklbghdfsjksdhnvfgkhdfkslgvhhrjkdshgnverhbgkrthbklg.php/?sort=cos"
	manghe   = "https://iw233.cn/api.php?sort=random"
)

func init() {
	engine := control.Register("tupian", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "图片",
		Help: "全部图片指令\n" +
			"- cos\n" +
			"- 兽耳\n" +
			"- 白毛\n" +
			"- 黑丝\n" +
			"- 白丝\n" +
			"- 丝袜\n" +
			"- 星空\n" +
			"- 开盲盒\n" +
			"- 随机壁纸\n" +
			"- 随机表情包\n" +
			"- 涩涩达咩/我要涩涩\n",
	})
	engine.OnFullMatchGroup([]string{"随机壁纸", "兽耳", "星空", "白毛", "我要涩涩", "涩涩达咩", "白丝", "黑丝", "丝袜", "随机表情包", "cos", "盲盒", "开盲盒"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var url string
			switch ctx.State["matched"].(string) {
			case "兽耳":
				url = shouer
			case "随机壁纸":
				url = bizhi
			case "白毛":
				url = baimao
			case "星空":
				url = xing
			case "我要涩涩", "涩涩达咩":
				url = sese
			case "白丝":
				url = baisi
			case "黑丝":
				url = heisi
			case "丝袜":
				url = siwa
			case "随机表情包":
				url = biaoqing
			case "cos":
				url = cos
			case "盲盒", "开盲盒":
				url = manghe
			}
			url2, err := bilibili.GetRealURL(url)
			if err != nil {
				ctx.SendChain(message.Text("获取图片地址失败惹", err))
				return
			}
			data, err := web.RequestDataWith(web.NewDefaultClient(), url2, "", referer, "", nil)
			if err != nil {
				ctx.SendChain(message.Text("获取图片失败惹"))
				return
			}
			ctx.SendChain(message.ImageBytes(data))
		})
}
