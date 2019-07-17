package scentser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
	"scentrealm_bcc/utils"
)

var logger = utils.Logger
var (
	ErrNoHandlerFound = errors.New("No Handler Found ")
	ErrJsonParseError = errors.New("Can Not Parse To Json Object ")
)

const (
	CmdConnect   = "Connect"
	CmdDisconnect   = "Disconnect"
	CmdPlaySmell = "PlaySmell"
	CmdStopPlay  = "StopPlay"
	CmdWakeup    = "WakeUp"
)

type BaseRequestJson struct {
	Cmd    string                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
}

type BaseResponseJson struct {
	//BaseRequestJson
	Cmd   string `json:"cmd"`
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func (resp *BaseResponseJson) ToJson() []byte {
	str, err := json.Marshal(*resp)
	if err == nil {
		return str
	}
	return nil
}

type ScentWsServer struct {
	ServerAddr  string
	ServeUri    string
	WsServer    *websocket.Upgrader
	CmdHandlers map[string]func(pas map[string]interface{}) error
}

func (ser *ScentWsServer) Init() {
	ser.WsServer = &websocket.Upgrader{} // use default options
	ser.ServeUri = "/scent/ctl"
	ser.ServerAddr = "localhost:8080"
	ser.CmdHandlers = make(map[string]func(pas map[string]interface{}) error)
}

func (ser *ScentWsServer) AddHandler(cmd string, handle func(pas map[string]interface{}) error) {
	ser.CmdHandlers[cmd] = handle
}

func (ser *ScentWsServer) ParseCMD(str string) *BaseResponseJson {
	req := &BaseRequestJson{}
	resp := &BaseResponseJson{}

	err := json.Unmarshal([]byte(str), &req)
	if err == nil {
		resp.Cmd = req.Cmd
		if handle, ok := ser.CmdHandlers[req.Cmd]; ok {
			err = handle(req.Params)
			if err != nil {
				resp.Code = 1001
				resp.Error = err.Error()
				return resp
			} else {
				return resp
			}
		}
		err = ErrNoHandlerFound
	}
	resp.Code = 1001
	resp.Error = ErrJsonParseError.Error()
	return resp
}

func (ser *ScentWsServer) ScentWsServerHandle(w http.ResponseWriter, r *http.Request) {
	c, err := ser.WsServer.Upgrade(w, r, nil)
	if err != nil {
		logger.Print("upgrade:", err)
		return
	}
	defer func() {
		err = c.Close()
		logger.Println("defer err", err)
	}()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			logger.Println("read:", err)
			break
		}
		logger.Printf("recv: %s\n", string(message))

		resp := ser.ParseCMD(string(message))
		logger.Printf("resp: %s\n", string(resp.ToJson()))

		err = c.WriteMessage(mt, resp.ToJson())

		if err != nil {
			logger.Println("write:", err)
			break
		}
	}
	logger.Println("end server")
}

func (ser *ScentWsServer) TestPage(w http.ResponseWriter, r *http.Request) {
	err := homeTemplate.Execute(w, "ws://"+r.Host+ser.ServeUri)
	if err != nil {
		fmt.Println(err)
	}
}

func (ser *ScentWsServer) RunServer() {
	http.HandleFunc(ser.ServeUri, ser.ScentWsServerHandle)
	http.HandleFunc("/", ser.TestPage)
	logger.Fatal(http.ListenAndServe(ser.ServerAddr, nil))
}

func RunTestWS() {

	var server = &ScentWsServer{}
	server.Init()
	server.RunServer()
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
	Date.prototype.Format = function(formatStr) {
		var str = formatStr;
		var Week = [ '日', '一', '二', '三', '四', '五', '六' ];
	
		str = str.replace(/yyyy|YYYY/, this.getFullYear());
		str = str.replace(/yy|YY/, (this.getYear() % 100) > 9 ? (this
				.getYear() % 100).toString() : '0'
				+ (this.getYear() % 100));
	
		str = str.replace(/MM/, this.getMonth() > 9 ? this.getMonth()
				.toString() : '0' + this.getMonth());
		str = str.replace(/M/g, this.getMonth());
	
		str = str.replace(/w|W/g, Week[this.getDay()]);
	
		str = str.replace(/dd|DD/, this.getDate() > 9 ? this.getDate()
				.toString() : '0' + this.getDate());
		str = str.replace(/d|D/g, this.getDate());
	
		str = str.replace(/hh|HH/, this.getHours() > 9 ? this
				.getHours().toString() : '0' + this.getHours());
		str = str.replace(/h|H/g, this.getHours());
		str = str.replace(/mm/, this.getMinutes() > 9 ? this
				.getMinutes().toString() : '0' + this.getMinutes());
		str = str.replace(/m/g, this.getMinutes());
	
		str = str.replace(/ss|SS/, this.getSeconds() > 9 ? this
				.getSeconds().toString() : '0' + this.getSeconds());
		str = str.replace(/s|S/g, this.getSeconds());
		str = str.replace(/u|U/g, (this.getMilliseconds()/1000 + "00000").substring(2,5) );
		return str;
	}
	var now  = function() {
		return (new Date()).Format('HH:mm:ss.u');
	}
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var input1 = document.getElementById("input1");
    var input2 = document.getElementById("input2");
    var input3 = document.getElementById("input3");
	var input4 = document.getElementById("input4");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = '[' + now() + '] ' + message;
        // output.appendChild(d);
		output.insertBefore(d,output.children[0])
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESP: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("send1").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input1.value);
        ws.send(input1.value);
        return false;
    };
    document.getElementById("send2").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input2.value);
        ws.send(input2.value);
        return false;
    };
    document.getElementById("send3").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input3.value);
        ws.send(input3.value);
        return false;
    };
    document.getElementById("send4").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input4.value);
        ws.send(input4.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>

<p><input id="input2" type="text" value='{"cmd":"Connect"}' style="width:500px;">
<button id="send2">Send Connect CMD</button>

<p><input id="input1" type="text" value='{"cmd":"WakeUp","params":{"blocking":false}}' style="width:500px;">
<button id="send1">Send WakeUp CMD</button>

<p><input id="input" type="text" value='{"cmd":"PlaySmell","params":{"smell":2,"duration":2000,"channel":0}}' style="width:500px;">
<button id="send">Send Play Smell CMD</button>

<p><input id="input3" type="text" value='{"cmd":"StopPlay","params":{"channel":0}}' style="width:500px;">
<button id="send3">Send Stop Play CMD</button>

<p><input id="input4" type="text" value='{"cmd":"Disconnect"}' style="width:500px;">
<button id="send4">Send Disconnect CMD</button>

</form>
</td><td valign="top" width="50%">
<div id="output" style="    scroll-behavior: auto;
    max-height: 800px;
    overflow: scroll;"></div>
</td></tr></table>
</body>
</html>

`))
