package websocket2

import (
	"github.com/gorilla/websocket"
	"net/http"
	"promptrun-api/utils"
	"sync"
	"time"
)

var (
	// 消息通道
	news = make(map[string]chan interface{})
	// websocket 客户端链接池
	clients = make(map[string]*websocket.Conn)
	// 互斥锁，防止程序对统一资源（news、clients）同时进行读写
	mux sync.Mutex
)

// websocket Upgrader
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消 ws 跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsHandler 处理 ws 请求
func WsHandler(w http.ResponseWriter, r *http.Request, id string) {
	if clients[id] != nil {
		utils.Log().Info("", "websocket client already exist, id: %s", id)
		return
	}

	// 创建一个定时器用于服务端心跳
	pingTicker := time.NewTicker(time.Second * 30)
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.Log().Error("", "websocket upgrade fail, errMsg: %s", err.Error())
		return
	}
	// 把与客户端的链接添加到客户端链接池中
	registerClient(id, conn)

	// 获取该客户端的消息通道
	m, exist := getNewsChannel(id)
	if !exist {
		m = make(chan interface{})
		addNewsChannel(id, m)
	}

	// 设置客户端关闭ws链接回调函数
	conn.SetCloseHandler(func(code int, text string) error {
		unregisterClient(id)
		utils.Log().Info("", "websocket close by clients, id: %s", id)
		return nil
	})

	go func() {
		for {
			select {
			case content, _ := <-m:
				// 从消息通道接收消息，然后推送给前端
				utils.Log().Info("", "websocket send message, id: %s, message: %s", id, content)
				err = conn.WriteJSON(content)
				if err != nil {
					utils.Log().Error("", "websocket write message fail, errMsg: %s", err.Error())
					conn.Close()
					unregisterClient(id)
					return
				}
			case <-pingTicker.C:
				// 服务端心跳:每30秒ping一次客户端，查看其是否在线
				utils.Log().Info("", "websocket ping, id: %s", id)
				err = conn.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					utils.Log().Info("", "websocket ping fail, will close socket, errMsg: %s", err.Error())
					conn.Close()
					unregisterClient(id)
					return
				}
			}
		}
	}()
}

// 将客户端注册到客户端链接池
func registerClient(id string, conn *websocket.Conn) {
	mux.Lock()
	clients[id] = conn
	mux.Unlock()
}

// 获取指定客户端链接
func getClient(id string) (conn *websocket.Conn, exist bool) {
	mux.Lock()
	conn, exist = clients[id]
	mux.Unlock()
	return
}

// 注销客户端链接
func unregisterClient(id string) {
	mux.Lock()
	delete(clients, id)
	utils.Log().Info("", "delete clients, id: %s", id)
	mux.Unlock()
}

// 添加用户消息通道
func addNewsChannel(id string, m chan interface{}) {
	mux.Lock()
	news[id] = m
	mux.Unlock()
}

// 获取指定用户消息通道
func getNewsChannel(id string) (m chan interface{}, exist bool) {
	mux.Lock()
	m, exist = news[id]
	mux.Unlock()
	return
}

// 删除指定消息通道
func deleteNewsChannel(id string) {
	mux.Lock()
	if m, ok := news[id]; ok {
		close(m)
		delete(news, id)
	}
	mux.Unlock()
}

// SendNew 发送 websocket 消息
func SendNew(id string, new interface{}) {
	m, exist := getNewsChannel(id)
	if exist {
		m <- new
	}
}
