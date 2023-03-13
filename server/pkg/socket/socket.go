package socket

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"server/pkg/code"
	"sync"
	"time"
)

// http升级websocket协议的配置，101
var Upgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 客户端读写消息
type Message struct {
	Type int    `json:"type"` // 类型
	Data []byte `json:"data"` // 数据
}

// Connection 客户端连接
type Connection struct {
	wsSocket *websocket.Conn // 底层websocket
	inChan   chan *Message   // 读队列
	outChan  chan *Message   // 写队列

	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	closeChan chan byte // 关闭通知
}

// ReadLoop
func (wsConn *Connection) ReadLoop() {
	for {
		// 读一个message
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			goto error
		}
		req := &Message{
			msgType,
			data,
		}
		// 放入请求队列
		select {
		case wsConn.inChan <- req:
		case <-wsConn.closeChan:
			goto closed
		}
	}
error:
	wsConn.Close()
closed:
}

// WriteLoop
func (wsConn *Connection) WriteLoop() {
	for {
		select {
		// 取一个应答
		case msg := <-wsConn.outChan:
			// 写给websocket
			if err := wsConn.wsSocket.WriteMessage(msg.Type, msg.Data); err != nil {
				goto error
			}
		case <-wsConn.closeChan:
			goto closed
		}
	}
error:
	wsConn.Close()
closed:
}

// procLoop
func (wsConn *Connection) procLoop() {
	// 这是一个同步处理模型（只是一个例子），如果希望并行处理可以每个请求一个gorutine，注意控制并发goroutine的数量!!!
	for {
		msg, err := wsConn.Read()
		if err != nil {
			fmt.Println("read fail")
			break
		}
		fmt.Println("接收消息", string(msg.Data))

		// 字符串装json
		type Data struct {
			Type string `json:"type"`
		}
		var a Data
		jsonData := []byte(string(msg.Data))
		err = json.Unmarshal(jsonData, &a)
		if err != nil {
			return
		}

		// 消息类型
		switch {
		case a.Type == "hello":
			data := map[string]interface{}{
				"type": "hello",
				"data": "hello",
				"code": code.SUCCESS,
				"msg":  code.GetErrMsg(code.SUCCESS),
			}

			dataJson, _ := json.Marshal(data)
			err = wsConn.Write(msg.Type, []byte(string(dataJson))) // 读到数据后，同步的去写数据。应该写成异步
			if err != nil {
				fmt.Println("write fail")
				break
			}
			break
		default:
			err := wsConn.Write(msg.Type, msg.Data)
			if err != nil {
				return
			} // 读到数据后，同步的去写数据。应该写成异步

			fmt.Println("错误")
			break
		}

	}
}

// HeartBeat
func (wsConn *Connection) HeartBeat() {
	// 服务端连接保活
	go func() {
		for { //不断向客户端写数据，其实没有它也是一样的，客户端可以检测到断开
			time.Sleep(5 * time.Second)
			data := map[string]interface{}{
				"type": "heartBeat",
				"data": "heartBeat",
				"code": code.SUCCESS,
				"msg":  code.GetErrMsg(code.SUCCESS),
			}

			dataJson, _ := json.Marshal(data)

			if err := wsConn.Write(websocket.TextMessage, []byte(string(dataJson))); err != nil {
				wsConn.Close()
				break
			}
		}
	}()

}

// 写数据
func (wsConn *Connection) Write(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &Message{messageType, data}:
	case <-wsConn.closeChan:
		return errors.New("websocket closed")
	}
	return nil
}

// Read 读数据
func (wsConn *Connection) Read() (*Message, error) {
	select {
	case msg := <-wsConn.inChan:
		return msg, nil
	case <-wsConn.closeChan:
	}
	return nil, errors.New("websocket closed")
}

// Close 连接关闭
func (wsConn *Connection) Close() {
	err := wsConn.wsSocket.Close()
	if err != nil {
		return
	}

	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}

// Handler 接口逻辑
func Handler(c *gin.Context) {
	// 应答客户端告知升级连接为websocket
	wsSocket, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	wsConn := &Connection{
		wsSocket:  wsSocket,
		inChan:    make(chan *Message, 1000),
		outChan:   make(chan *Message, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
	}
	// 处理器
	go wsConn.procLoop()
	// 读协程
	go wsConn.ReadLoop()
	// 心跳保活
	go wsConn.HeartBeat()
	// 写协程
	go wsConn.WriteLoop()
}
