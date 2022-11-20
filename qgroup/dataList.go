package qgroup

import (
	"encoding/json"
	"github.com/Mrs4s/go-cqhttp/coolq"
	"github.com/Mrs4s/go-cqhttp/global"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

type BotQAMap struct {
	KeyMap     map[string]int
	AnswerList map[int]string
	KeyIndex   int
}

type GroupRule struct {
	bot        *coolq.CQBot
	BanUsers   []int
	BanAdWords []string
	BotQA      BotQAMap
}

var defaultRule *GroupRule

func GetDefaultRule() *GroupRule {
	return defaultRule
}

func CreateNewRule(n string) *GroupRule {
	defaultRule = &GroupRule{}
	defaultRule.InitGroupRuleConfig()
	return defaultRule
}

// keyword1, kw2, kw3 ---> answers

func (g *GroupRule) InitGroupRuleConfig() {
	g.BanUsers = []int{}
	g.BanAdWords = []string{}
	g.BotQA = BotQAMap{AnswerList: make(map[int]string), KeyMap: make(map[string]int)}

	load("data/banUser.json", &g.BanUsers)
	load("data/banWords.json", &g.BanAdWords)
	load("data/BotQA.json", &g.BotQA)

}

func Save(p string, pData interface{}) error {
	data, err := json.MarshalIndent(pData, "", "\t")
	if err != nil {
		return err
	}
	global.WriteAllText(p, string(data))
	return nil
}
func (c *BotQAMap) addQA(q string, a string, answerIndex int) error {
	if q == "" {
		return nil
	}
	if oldKey, ok := c.KeyMap[q]; ok {
		log.Infof("覆盖 %s", oldKey)
		c.AnswerList[oldKey] = a
	} else {
		log.Infof("new key %s", q)
		//c.KeyIndex++
		c.KeyMap[q] = c.KeyIndex
		c.AnswerList[answerIndex] = a

	}
	return nil
}

func (c *BotQAMap) AddQAs(qs string, a string) error {
	if qs == "" || a == "" {
		return nil
	}
	qsarr := strings.Split(qs, "|")
	for _, v := range qsarr {
		if v != "" {
			c.addQA(v, a, c.KeyIndex)
		}
	}
	c.KeyIndex++
	Save("data/BotQA.json", c)
	return nil
}

func (g *GroupRule) AddBanUser(uid int) {
	for _, v := range g.BanUsers {
		if v == uid {
			return
		}
	}
	g.BanUsers = append(g.BanUsers, uid)
}

func (g *GroupRule) AddBanKeyWord(keyword string) {
	for _, v := range g.BanAdWords {
		if v == keyword {
			return
		}
	}
	g.BanAdWords = append(g.BanAdWords, keyword)
}

func load(p string, ptype interface{}) {
	if !global.PathExists(p) {
		log.Warnf("尝试加载配置文件 %v 失败: 文件不存在", p)
		return
	}

	err := json.Unmarshal([]byte(global.ReadAllText(p)), ptype)
	if err != nil {
		log.Warnf("尝试加载配置文件 %v 时出现错误: %v", p, err)
		log.Infoln("原文件已备份")
		os.Rename(p, p+".backup"+strconv.FormatInt(time.Now().Unix(), 10))
	}
}

func Test(g *GroupRule) {
	//LoadData()
	g.AddBanUser(1000)
	g.AddBanKeyWord("敏感")
	//g.BotQA.AddQA("怎么玩", "十万十万冲不要停下来")

	Save("data/banUser.json", g.BanUsers)
	Save("data/banWords.json", g.BanAdWords)
	Save("data/BotQA.json", g.BotQA)
}
