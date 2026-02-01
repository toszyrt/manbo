package main

import (
	"log"
	"strings"
	"time"

	"github.com/toszyrt/manbo/openclaw"
	"github.com/toszyrt/manbo/xiaoai"
)

type Manbo struct {
	xCli *xiaoai.MiClient
	oCli *openclaw.WSClient

	name     string
	hardware string
	miotDID  string
	deviceID string

	lastQuery     string
	lastTimestamp int64
}

func NewManbo() *Manbo {
	return &Manbo{}
}

func (m *Manbo) Init() {
	// 1. 初始化 xiaoai 客户端
	m.xCli = xiaoai.NewMiClient(nil)
	devices, err := m.xCli.MinaDeviceList()
	if err != nil {
		log.Fatal(err)
	}
	if len(devices.Data) == 0 {
		log.Fatal("no devices found")
	}
	for _, d := range devices.Data {
		if d.Presence == "online" {
			m.name = devices.Data[0].Name
			m.hardware = devices.Data[0].Hardware
			m.miotDID = devices.Data[0].MiotDID
			m.deviceID = devices.Data[0].DeviceID
			break
		}
	}
	if len(m.name) == 0 {
		log.Fatal("devices offline")
	}
	log.Printf("✅ xiaoai device found")
	log.Printf("name=%s, hardware=%s, miot_did=%s, device_id=%s", m.name, m.hardware, m.miotDID, m.deviceID)
	// 2. 初始化 openclaw 客户端
	m.oCli = openclaw.NewWsClient(nil)
	m.oCli.ConnectAndHandshake()
	log.Println("✅ openclaw gateway connected")
}

func (m *Manbo) Loop() {
	for {
		// 1. 1s获取1次最新对话
		time.Sleep(time.Second)
		conservation, err := m.xCli.MinaGetLatestConversation(m.hardware, m.deviceID)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(conservation.Records) != 1 {
			log.Println("records len != 1")
			continue
		}
		record := conservation.Records[0]
		query, timestamp := record.Query, record.Time
		lastQuery, lastTimestamp := m.lastQuery, m.lastTimestamp
		m.lastQuery, m.lastTimestamp = query, timestamp
		// 2. 判断是否要处理
		if len(lastQuery) == 0 && lastTimestamp == 0 {
			continue // 上线第一轮不处理
		}
		if query == lastQuery && timestamp == lastTimestamp {
			continue // 没有新的消息
		}
		evtName, err := DispatchEvent(query)
		if len(evtName) != 0 {
			time.Sleep(3 * time.Second) // 等待"好的"说完
			reply := "已下发任务, " + evtName
			if err == nil {
				reply += ", 成功"
			} else {
				reply += ", 失败, 错误信息为, " + err.Error()
			}
			if err := m.xCli.MiotAction(m.miotDID, 5, 1, []any{reply}); err != nil {
				log.Println(err)
			}
		}
		if !strings.Contains(query, "请") {
			continue // 只处理包含“请”的消息
		}
		// 3. 桥接到openclaw
		if err := m.xCli.MiotAction(m.miotDID, 5, 1, []any{"稍等"}); err != nil {
			log.Println(err)
		}
		reply, err := m.oCli.SendQueryAndWait(query)
		if err != nil {
			log.Println(err)
			continue
		}
		if err := m.xCli.MiotAction(m.miotDID, 5, 1, []any{reply}); err != nil {
			log.Println(err)
		}
	}
}
