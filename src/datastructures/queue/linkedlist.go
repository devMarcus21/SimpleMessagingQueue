package queue

type Node struct {
	next    *Node
	prev    *Node
	message QueueMessage
}

type LinkedList struct {
	size int
	tail *Node
	head *Node
}

func NewLinkedList() LinkedList {
	head := Node{}
	tail := Node{}

	tail.next = &head
	tail.prev = nil

	head.next = nil
	head.prev = &tail

	return LinkedList{
		size: 0,
		tail: &tail,
		head: &head,
	}
}

func (list *LinkedList) addBackOfList(message QueueMessage) {
	tailNode := list.tail
	nodeAfter := tailNode.next

	newNode := Node{}
	newNode.message = message

	tailNode.next = &newNode
	newNode.prev = tailNode

	nodeAfter.prev = &newNode
	newNode.next = nodeAfter

	list.size++
}

func (list *LinkedList) removeFrontOfList() (QueueMessage, bool) {
	if list.IsEmpty() {
		return QueueMessage{}, false
	}

	frontNode := list.head
	removeNode := frontNode.prev
	backNode := removeNode.prev

	message := removeNode.message

	backNode.next = frontNode
	frontNode.prev = backNode

	list.size--
	return message, true
}

func (list *LinkedList) Offer(message QueueMessage) {
	list.addBackOfList(message)
}

func (list *LinkedList) Poll() (QueueMessage, bool) {
	return list.removeFrontOfList()
}

func (list *LinkedList) Size() int {
	return list.size
}

func (list *LinkedList) IsEmpty() bool {
	return list.size == 0
}
