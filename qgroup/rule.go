package qgroup

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strings"
)
import "github.com/Mrs4s/go-cqhttp/coolq"

func (g *GroupRule) Listern(bot *coolq.CQBot) {
	g.bot = bot
	g.bot.OnEventPush(g.OnNewUserRequest)
	g.bot.OnEventPush(g.OnGroupMsg)
	//g.bot.OnEventPush(g.onNewGroupMate)
}

func (g *GroupRule) OnNewUserRequest(e *coolq.Event) {
	m := e.RawMsg
	//"post_type":    "request",
	//	"request_type": "group",
	//	"sub_type":     "add",
	//	"group_id":     e.GroupCode,
	//	"user_id":      e.RequesterUin,
	//	"comment":      e.Message,
	//	"flag":         flag,
	//	"time":         time.Now().Unix(),
	//	"self_id":      c.Uin,
	if m["post_type"] != "request" && m["request_type"] != "group" && m["sub_type"] != "add" {
		return
	}
	uid := m["user_id"]

	for _, u := range g.BanUsers {
		if u == uid {
			g.bot.CQProcessGroupRequest(m["flag"].(string), "add", "Ban", false)
			// 拒绝
			return
		}
	}
}

func (g *GroupRule) OnGroupMsg(e *coolq.Event) {
	m := e.RawMsg
	if mType, ok := m["message_type"]; !ok || mType != "group" {
		// 过滤非群组消息
		return
	}
	if g.bot.Client.Uin == m["user_id"].(int64) {
		fmt.Println("过滤自己的数据")
		return
	}

	if g.checkBanWords(m["message"].(string)) {
		log.Infoln("敏感词撤回")
		g.bot.CQDeleteMessage(m["message_id"].(int32))
		g.bot.CQSetGroupBan(m["group_id"].(int64), m["user_id"].(int64), 60*60)
		//撤回
		//m["user_id"]
		return
	}

	if ans := g.checkQA(m["message"].(string)); ans != "" {
		log.Infoln("匹配QA")
		//匹配成功, 发消息
		g.bot.CQSendGroupMessage(m["group_id"].(int64), gjson.Result{Type: gjson.String, Str:ans}, false)
		return
	}
	//m["user_id"]
	//m["message"]
	//m["group_id"]
}

// func (g *GroupRule) onNewGroupMate(m coolq.MSG) {
// 	if m["post_type"] != "notice" {
// 		return
// 	}
// 	fmt.Println("#11111")
// 	if m["notice_type"] != "group_increase" {
// 		return
// 	}
// 	fmt.Println("#2222222")

// 	if m["sub_type"] == "approve" {
// 		fmt.Println("new guy")
// 		msg := fmt.Sprintf("欢迎@%d, 可以发送[ 攻略介绍 ]查看攻略信息 \n新手可以领取礼品码 xslb1  xslb2  xslp1  xslp2  fz001  fz002  88888888")
// 		g.bot.CQSendGroupMessage(m["group_id"].(int64), msg, false)
// 	}

// }

func (g *GroupRule) checkQA(msg string) string {
	for q, idx := range g.BotQA.KeyMap {
		if strings.Contains(msg, q) {
			return g.BotQA.AnswerList[idx]
		}
	}
	return ""
}

func (g *GroupRule) checkBanWords(msg string) bool {
	for _, w := range g.BanAdWords {
		if strings.Contains(msg, w) {
			return true
		}
	}
	return false
}
