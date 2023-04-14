// Package vtbwife 抽vtb老婆
package vtbwife

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
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
				logrus.Debugln("[vtbwife]读取vtbwife数据文件失败: ", err)
				return false
			}
			// 将文件内容转换为单词
			keys = strings.Split(string(content), "\n")
			logrus.Debugln("[vtbwife]加载", len(keys), "位wife数据...")
			return true
		})).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		var key, u, b string
		var ok bool
		var fix int
		for i := 0; i < 3; i++ {
			key = keys[fcext.RandSenderPerDayN(ctx.Event.UserID, len(keys))+fix]
			u, b, ok = geturl(key)
			if !ok {
				fix++
				continue
			}
			break
		}
		if !ok {
			ctx.SendChain(message.Text("获取图片链接失败"))
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
		bs := message.Text(
			b,
		)
		ctx.SendChain(message.At(ctx.Event.UserID), txt, message.ImageBytes(img), bs)
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
	u, ok = doc.Find(".infobox-image").Attr("src")  // class加.
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
		if bs[k] == "粉丝勋章" {
			break
		} else if bs[k] == "活动范围" || bs[k] == "Bilibili粉丝勋章" {
			k += 2
			continue
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

// 获取图片id,图片
/*i, err := getPicID(key)
fmt.Println(key, i)
if err != nil {
	ctx.SendChain(message.Text("ERROR: ", err))
	return
}
path, err := getPic(i[0])
if err != nil {
	ctx.SendChain(message.Text("ERROR: ", err))
	return
}*/
