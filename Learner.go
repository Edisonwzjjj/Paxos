type Learner struct {
	lis net.listener

	id int

	acceptedMsg map[int]MsgArgs
}