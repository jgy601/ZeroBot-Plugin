// Package vtbwife 抽vtb老婆
package vtbwife

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/PuerkitoBio/goquery"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() { // 插件主体
	engine := control.Register("vtbwife", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "抽vtb老婆",
		Help:              "- 抽vtb",
		PrivateDataFolder: "vtb_wife",
	})
	var keys []string
	engine.OnRegex(`^抽vtb$`, fcext.DoOnceOnSuccess(
		func(ctx *zero.Ctx) bool {
			content, err := os.ReadFile(engine.DataFolder() + "wife_list.txt") // 779分界
			if err != nil {
				log.Println("[vtbwife]读取vtbwife数据文件失败: ", err)
				return false
			}
			// 将文件内容转换为单词
			keys = strings.Split(string(content), "\n")
			log.Println("[vtbwife]加载", len(keys), "位wife数据...")
			return true
		})).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		var key, u, b string
		var ok bool
		for i := 0; i < 3; i++ {
			key = keys[fcext.RandSenderPerDayN(ctx.Event.UserID, len(keys))+i]
			u, b, ok = geturl(key)
			if !ok {
				continue
			}
			break
		}
		if !ok {
			ctx.SendChain(message.Text("-获取图片链接失败"))
			return
		}
		img, err := web.GetData(u)
		if err != nil {
			ctx.SendChain(message.Text("-获取图片失败惹", err))
			return
		}
		txt := message.Text(
			"\n今天你的VTB老婆是: ", key,
		)
		if id := ctx.SendChain(message.At(ctx.Event.UserID), txt, message.ImageBytes(img), message.Text(b)); id.ID() == 0 {
			ctx.SendChain(message.At(ctx.Event.UserID), txt, message.Text("图片发送失败...\n"), message.Text(b))
		}
	})
}

func geturl(kword string) (u, brief string, ok bool) {
	resp, err := http.Get("https://zh.moegirl.org.cn/" + url.QueryEscape(kword))
	if err != nil {
		return "", "", false
	}
	defer resp.Body.Close()
	// 使用goquery解析网页内容
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", false
	}
	u, ok = doc.Find(".infobox-image").Attr("src") // class加.
	doc.Find("style").Remove()
	doc.Find("script").Remove()
	doc.Find("big").Remove()
	b := doc.Find(".moe-infobox").Find("tr").Text() // class加.
	bs := strings.Split(b, "\n")
	// 寻找"基本资料"
	var (
		k int
		f bool // 判断换行
	)
	for kk, vv := range bs {
		if vv == "基本资料" {
			brief = vv + "\n"
			k = kk + 1
		}
	}
	for ; k < len(bs); k++ {
		switch bs[k] {
		case "直播关联", "制作所属", "个人数据", "名字":
			continue
		case "活动范围":
			k += 2
			continue
		case "进入直播间":
			bs[k] = ""
		}
		if t := strings.TrimSpace(bs[k]); t == "" {
			continue
		} else {
			f = !f
			brief += t
		}
		if f {
			brief += ": "
		} else {
			brief += "\n"
		}
	}
	brief = strings.TrimSpace(brief)
	return
}
