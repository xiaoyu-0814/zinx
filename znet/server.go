package znet

import (
	"fmt"
	"net"
	"zinx/iface"
	"zinx/util"
)

var _ iface.IServer = new(Server)

//IServer的接口实现，定义一个server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//给当前server添加一个router
	Router iface.IRouter
}

//初始化server模块
func NewServer(name string) iface.IServer {
	return &Server{
		Name:      util.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        util.GlobalObject.Host,
		Port:      util.GlobalObject.Port,
		Router:    nil,
	}
}

/*func CallBackToClient(conn *znet.TCPConn,data []byte,cnt int) error  {
	//回显业务
	fmt.Println("conn Handle CallBackToClient")
	if _,err := conn.Write(data[:cnt]);err != nil {
		fmt.Println("write error:",err)
		return errors.New("CallBackToClient error")
	}

	return nil
}*/

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP:%s,Port %d, is starting\n", s.IP, s.Port)
	go func() {
		//1.获取一个TCP的地址
		//将addr作为一个格式为"host"或"ipv6-host%zone"的IP地址来解析。 函数会在参数net指定的网络类型上解析，net必须是"ip"、"ip4"或"ip6"
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolvrIPAddr error : ", err)
			return
		}

		//监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("listenTCP error:%s", err.Error())
			return
		}

		fmt.Printf("start server %s succ\n", s.Name)
		var cid uint32
		cid = 0

		//阻塞的等待客户端链接，处理客户端业务
		for {
			//如果有客户端链接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}

			//将处理新链接的业务绑定
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			//启动
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	//将服务器的资源，状态或者一些已开辟的链接，进行停止或者回收

}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//做一些启动服务之后的额外服务

	//阻塞状态
	select {}
}

//添加一个路由功能
func (s *Server) AddRouter(router iface.IRouter) {
	s.Router = router
	fmt.Println("add router success")
}
