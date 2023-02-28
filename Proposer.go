type Proposer struct {
	//服务器id
	id int
	//提议者已知最大轮次
	round int
	//提案编号
	number int
	acceptor []int
}