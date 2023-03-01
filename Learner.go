package paxos

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)


type Learner struct {
	lis net.listener

	id int

	acceptedMsg map[int]MsgArgs
}


func (l * Learner) Learn (args *MsgArgs, reply *MsgReply) error {
	a := l.acceptedMsg[args.From]
	if a.Number < args.Number {
		l.acceptedMsg[args.From] = *args
		reply.Ok = true
	} else {
		reply.Ok = false
	}
	return nil

}

func (l *Learner) chosen() interface{} {
	acceptCount := make(map[int]int)
	acceptMsg := make(map[int]MsgArgs)

	for _, accepted := range l.acceptedMsg {
		if accepted.Number != 0 {
			acceptCount[accepted.Number]++
			acceptMsg[accepted.Number] = accepted

		}

	}

	for n, count := range acceptCount {
		if count >= l.majority() {
			return acceptMsg[n].Value
		}
	}
	return nil
}


func (l* Learner) majority() int {
	return len(l.acceptedMsg) / 2 + 1
}


func newLearner(id int , acceptorIds []int) *Learner {
	Learner := &Learner{
		id : id,
		acceptedMsg: make(map[int]MsgArgs), 
	}
	for _, aid := range 
}