package paxos

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)


type Acceptor struct {
	lis net.listener
	id int
	minProposal int // 0 表示没收到信息， 接受者承诺的提案编号
	acceptedNumber int // 已接受的提案编号
	acceptedValue interface{} //没有是nil

	learners []int
}

func newAcceptor(id int , learners []int) *Acceptor {
	acceptor := &Acceptor{
		id : id 
		learners : learners
	}
	acceptor.server()
	return acceptor
}

func (a *Acceptor) server() {
	rpcs := rpc.NewServer()
	rpcs.Register(a)
	addr := fmt.Sprintf(": %d". a.id)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("Listen error :", e)
	}
	a.lis = l;

	go func () {
		for {
			conn, err := a.lis.Accept()
			if err != nil {
				continue
			}
			go rpcs.ServeConn(conn)
		}
	}()
}
func (a *Acceptor) Accept(args *MsgArgs, reply *MsgReply) error {
	if args.Number >= a.minProposal {
		a.minProposal = args.Number
		a.acceptedNumber = args.Number
		a.acceptedValue = args.Value
		reply.Ok = true
		// 后台转发接受的提案给学习者
		for _, lid := range a.learners {
			go func(learner int) {
				addr := fmt.Sprintf("127.0.0.1:%d", learner)
				args.From = a.id
				args.To = learner
				resp := new(MsgReply)
				ok := call(addr, "Learner.Learn", args, resp)
				if !ok {
					return
				}
			}(lid)
		}
	} else {
		reply.Ok = false
	}
	return nil
}

func (a *Acceptor) close () {
	a.lis.Close()
}
