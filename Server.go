package main

import (
	"bufio"
	"fmt"
	"net"
)

//定义协程池类型
type Pool struct {
	worker_num  int           // 协程池最大worker数量,限定Goroutine的个数
	JobsChannel chan net.Conn // 协程池内部的任务就绪队列
}

//创建协程池
func NewPool(cap int) *Pool {
	p := Pool{
		worker_num:  cap,
		JobsChannel: make(chan net.Conn),
	}
	return &p
}

//处理函数
func process(conn net.Conn, work_ID int) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		wokerHead := fmt.Sprintf("[Woker %d]", work_ID)
		//读取数据放到buf中，n为长度
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}

		recvStr := string(buf[:n])
		fmt.Println(wokerHead + "收到Client端发的数据" + recvStr)
		conn.Write([]byte(wokerHead + "success"))

	}

}

//协程池中每个Worker的功能
func (p *Pool) Worker(work_ID int) error {

	//worker不断的从JobsChannel内部任务队列中拿Conn
	for conn := range p.JobsChannel {
		//如果拿到Conn,则执行对应处理
		//除非连接断裂，否则会一直卡在process里，保证了连接数量
		process(conn, work_ID)
	}

	return nil
}

//协程池Pool开始工作
func (p *Pool) Run(ListenAddr string) {
	listen, err := net.Listen("tcp", ListenAddr)
	if err != nil {
		fmt.Println("server fail to listen")
		return
	}

	// 首先根据协程池的worker数量限定,开启固定数量的Worker,
	// 每一个Worker用一个Goroutine承载
	for i := 0; i < p.worker_num; i++ {
		go p.Worker(i)
	}

	// 将新申请的连接加入到就绪队列
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}

		p.JobsChannel <- conn
	}

}

//关闭协程池
func (p *Pool) Close() {
	close(p.JobsChannel)
}

func main() {

	p := NewPool(3)

	p.Run("127.0.0.1:20000")

	p.Close()

}
