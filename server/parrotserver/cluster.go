package parrotserver
//
//import (
//	"io"
//	"net"
//	"reflect"
//	"strconv"
//	"strings"
//	"time"
//	"unsafe"
//)
//
//type Node struct {
//	createTime  int64
//	name        string
//	flags       uint16
//	configEpoch uint64
//	slots       [CLUSTER_SLOTS / 8]byte
//	numslots    int
//	numslaves   int
//	slave       **Node
//	slaveof     *Node
//	ip          string
//	hostname    string
//	port        int
//	cport       int
//	/* 该节点相关的连接对象（连接状态是 established），这个 link 是 TCP 客户端发送数据的 link */
//	link *Link /* TCP/IP link established toward this node */
//	/* accept 到的连接，这个是 TCP 服务端生成的 link，用来接收数据 */
//	inbound_link *Link /* TCP/IP link accepted from this node */
//	/* 该节点保存下线通知的链表 */
//	//list *fail_reports       /* List of nodes signaling this as failing */
//}
//
///* 集群状态，每个节点都保存着一个这样的状态，记录了它们眼中集群的样子 */
//type State struct {
//	/* 指向当前节点的指针 */
//	myself *Node /* This node */
//	/* 集群当前纪元，用于故障转移 */
//	currentEpoch uint64
//	/* 集群状态 */
//	state int /* CLUSTER_OK, CLUSTER_FAIL, ... */
//	/* 集群中至少处理一个槽的节点数量 */
//	size int /* Num of master nodes with at least one slot */
//	/* 保存集群节点的字典，键是节点名称，值是 clusterNode 结构的指针 */
//	nodes map[string]*Node /* Hash table of name . clusterNode structures */
//	/* 集群节点黑名单（包括 myself），可以防止在集群中的节点二次加入集群
//	 * 黑名单可以防止被 forget 的节点重新添加到集群节点 */
//	nodes_black_list map[string]*Node /* Nodes we don't re-add for a few seconds. */
//	/* 记录要从当前节点迁移到目标节点的槽，以及迁移的目标节点 */
//	migrating_slots_to [CLUSTER_SLOTS]*Node
//	/* 记录从其他节点迁移出来的槽 */
//	importing_slots_from [CLUSTER_SLOTS]*Node
//	/* 负责处理各个槽的节点 */
//	slots [CLUSTER_SLOTS]*Node
//	//rax *slots_to_channels
//	/* The following fields are used to take the slave state on elections. */
//	/* 之前或下一次选举的时间，主要是用来限制当前节点下一次投票发起的时间 */
//	failover_auth_time int64 /* Time of previous or next election. */
//	/* 节点获得支持的票数，从节点 */
//	failover_auth_count int /* Number of votes received so far. */
//	/* 如果为 True，表示该节点已经向其他节点发送了投票请求 */
//	failover_auth_sent int /* True if we already asked for votes. */
//	/* 该从节点在当前请求中的排名，该值根据复制偏移量计算而来，最终用于确定 slave 节点发起投票的时间
//	 * 注：排名就是当前从节点对应的主节点下所有从节点复制偏移量大于当前节点复制偏移量的数量，
//	 * 也就是说复制偏移量越大，排名越前，而排名会用作 failover_auth_time 的计算，排名越后，
//	 * failover_auth_time 也就越大，发起选举的时间越晚，
//	 * 即 rank 值越小的节点通常有更大的复制偏移量，它能越早发起选举竞争主节点 */
//	failover_auth_rank int /* This slave rank for current auth request. */
//	/* 当前选举的纪元 */
//	failover_auth_epoch uint64 /* Epoch of the current election. */
//	/* 从节点不能执行故障转移的原因 */
//	cant_failover_reason int /* Why a slave is currently not able to
//	   failover. See the CANT_FAILOVER_* macros. */
//	/* Manual failover state in common. */
//	/* 为 0 表示没有正在进行手动故障转移，否则表示手动故障转移的时间限制
//	 * 代码逻辑里会使用该属性来判断是手动故障转移，还是自动故障转移 */
//	//mstime_t mf_end           /* Manual failover time limit (ms unixtime).
//	//   It is zero if there is no MF in progress. */
//	/* Manual failover state of master. */
//	/* 执行手动故障转移的从节点 */
//	mf_slave *Node /* Slave performing the manual failover. */
//	/* Manual failover state of slave. */
//	/* 从节点记录手动故障转移时的主节点偏移量 */
//	mf_master_offset int64 /* Master offset the slave needs to start MF
//	   or -1 if still not received. */
//	/* 非 0 表示可以手动故障转移 */
//	mf_can_start int /* If non-zero signal that the manual failover
//	   can start requesting masters vote. */
//	/* The following fields are used by masters to take state on elections. */
//	/* 集群最近一次投票的纪元 */
//	lastVoteEpoch uint64 /* Epoch of the last vote granted. */
//	/* 调用 clusterBeforeSleep() 所做的一些事 */
//	todo_before_sleep int /* Things to do in clusterBeforeSleep(). */
//	/* Stats */
//	/* Messages received and sent by msgtype. */
//	/* 发送的字节数 */
//	stats_bus_messages_sent [CLUSTERMSG_TYPE_COUNT]int64
//	/* 通过 cluster 接收到的消息数量 */
//	stats_bus_messages_received [CLUSTERMSG_TYPE_COUNT]int64
//	stats_pfail_nodes           int64 /* Number of nodes in PFAIL status,
//	   excluding nodes without address. */
//	stat_cluster_links_buffer_limit_exceeded uint64 /* Total number of cluster links freed due to exceeding buffer limit */
//}
//
//type Link struct {
//	createTime int64
//	Conn       net.Conn
//	sendBuf    []byte
//	rcvBuf     []byte
//	node       *Node
//	inbound    int
//}
//
//type clusterMsgDataGossip struct {
//	/* 节点名称 */
//	nodename string
//	/* 发送 ping 的时间 */
//	ping_sent uint32
//	/* 接收 pong 的时间 */
//	pong_received uint32
//	/* 节点 IP 地址 */
//	ip string
//	/* 节点的端口 */
//	port uint16
//	/* 节点监听集群通信端口 */
//	cport uint16
//	/* 节点状态 flags */
//	flags uint16 /* node.flags copy */
//	/* 如果是 TLS 协议，该属性标识实际通信端口 */
//	pport uint16 /* plaintext-port, when base port is TLS */
//	/* 预留 */
//	notused1 uint16
//}
//
///* 包体内容 */
//type ping struct {
//	/* Array of N clusterMsgDataGossip structures */
//	gossip [1]clusterMsgDataGossip
//}
//
///* 用于描述集群节点间互相通信的消息的结构，包头 */
//type clusterMsg struct {
//	/* 固定 RCmb 标识 */
//	sig [4]byte /* Signature "RCmb" (Redis Cluster message bus). */
//	/* 消息的总长度 */
//	totlen uint32 /* Total length of this message */
//	/* 协议版本，当前设置为 1 */
//	ver uint16 /* Protocol version, currently set to 1. */
//	/* 发送方监听的端口 */
//	port uint16 /* TCP base port number. */
//	/* 包类型，（接收到包后通过该属性来决定如何解析包体） */
//	messagetype uint16 /* Message msgtype */
//	/* data 中的 gossip section 个数（供 ping pong meet 包使用） */
//	count uint16 /* Only used for some kind of messages. */
//	/* 发送方节点记录的集群当前纪元 */
//	currentEpoch uint64 /* The epoch accordingly to the sending node. */
//	/* 发送方节点对应的配置纪元（如果为从节点，为该从节点对应的主节点） */
//	configEpoch uint64 /* The config epoch if it's a master, or the last
//	   epoch advertised by its master if it is a
//	   slave. */
//	/* 如果为主节点，该值标识复制偏移量，如果为从，该值表示从已处理的偏移量 */
//	offset uint64 /* Master replication offset if node is a master or
//	   processed replication offset if node is a slave. */
//	/* 发送方名称 */
//	sender string /* Name of the sender node */
//	/* 发送方提供服务的 slot 映射表，（如果为从，则为该从所对应的主提供服务的 slot 映射表） */
//	myslots [CLUSTER_SLOTS / 8]byte
//	/* 发送方如果为从，则该字段为对应的主的名称 */
//	slaveof string
//	/* 发送方 IP */
//	myip string /* Sender IP, if not all zeroed. */
//	/* 和该包一起发送的扩展数 */
//	extensions uint16 /* Number of extensions sent along with this packet. */
//	/* 预留属性 */
//	notused1 [30]byte /* 30 bytes reserved for future usage. */
//	/* 发送方实际发送数据的端口 */
//	pport uint16 /* Sender TCP plaintext port, if base port is TLS */
//	/* 发送方监听的 cluster bus 端口 */
//	cport uint16 /* Sender TCP cluster bus port */
//	/* 发送方节点所记录的 flags */
//	flags uint16 /* Sender node flags */
//	/* 发送方节点所记录的集群状态 */
//	state int /* Cluster state from the POV of the sender */
//	/* 目前只有 mflags[0] 会在手动 failover 时使用 */
//	mflags [3]byte /* Message flags: CLUSTERMSG_FLAG[012]_... */
//	/* 包体内容 */
//	data ping
//}
//
//func clusterInit() {
//	/* 获取集群通信的端口 */
//	/* 这里会在端口下创建 socket fd 并赋值给 server.cfd ，然后会监听该端口 */
//	ln, err := net.Listen("tcp", "127.0.0.1:28888")
//	if err != nil {
//		panic("failed")
//	}
//
//	for {
//		conn, e := ln.Accept()
//		if e != nil {
//
//		}
//		go serveConn(conn)
//	}
//}
//
//// 读事件处理器
//// 首先读入内容的头，以判断读入内容的长度
//// 如果内容是一个 whole packet ，那么调用函数来处理这个 packet 。
//func serveConn(conn net.Conn) error {
//	buf := make([]byte,int(unsafe.Sizeof(clusterMsg{})))
//	var link Link
//	nread, err := io.ReadAtLeast(conn, buf, 8)
//	if err != nil || nread < 8 {
//		return err
//	}
//	hdr := (*clusterMsg)(unsafe.Pointer(
//		(*reflect.SliceHeader)(unsafe.Pointer(&link.rcvBuf)).Data,
//	))
//
//	if nread == 8 {
//		nread, err :=io.ReadAtLeast(conn, buf[8:], int(hdr.totlen))
//		if err!=nil || nread != int(hdr.totlen){
//			return err
//		}
//		link.rcvBuf = buf
//		link.Conn = conn
//		if err := clusterProcessPacket(&link); err != nil {
//
//		}
//	}
//
//
//	return nil
//
//}
//
///* When this function is called, there is a packet to process starting
// * at node.rcvbuf. Releasing the buffer is up to the caller, so this
// * function should just handle the higher level stuff of processing the
// * packet, modifying the cluster state if needed.
// *
// * 当这个函数被调用时，说明 node.rcvbuf 中有一条待处理的信息。
// * 信息处理完毕之后的释放工作由调用者处理，所以这个函数只需负责处理信息就可以了。
// *
// * The function returns 1 if the link is still valid after the packet
// * was processed, otherwise 0 if the link was freed since the packet
// * processing lead to some inconsistency error (for instance a PONG
// * received from the wrong sender ID).
// *
// * 如果函数返回 1 ，那么说明处理信息时没有遇到问题，连接依然可用。
// * 如果函数返回 0 ，那么说明信息处理时遇到了不一致问题
// * （比如接收到的 PONG 是发送自不正确的发送者 ID 的），连接已经被释放。
// */
//
//func clusterProcessPacket(link *Link) error {
//	// 指向消息头
//	hdr := (*clusterMsg)(unsafe.Pointer(
//		(*reflect.SliceHeader)(unsafe.Pointer(&link.rcvBuf)).Data,
//	))
//	// 消息的长度
//	 totlen := hdr.totlen
//
//	// 消息的类型
//	 msgtype := hdr.messagetype
//
//	// 消息发送者的标识
//	// flags := hdr.flags
//
//	 senderCurrentEpoch := 0
//	 //senderConfigEpoch := 0
//
//	var  sender *Node
//
//	// 更新接受消息计数器
//	//server.cluster.stats_bus_messages_received++
//
//	//redisLog(REDIS_DEBUG,"--- Processing packet of msgtype %d, %lu bytes",
//	//msgtype, (unsigned long) totlen)
//
//	/* Perform sanity checks */
//	// 合法性检查
//	if totlen < 16{
//		return nil
//	}
//                    /* At least signature, version, totlen, count. */
//	if (hdr.ver) != 0 {return nil /* Can't handle versions other than 0.*/}
//	if int(totlen) > len(link.rcvBuf) {
//		return nil
//	}
//	if msgtype == CLUSTERMSG_TYPE_PING ||
//		msgtype == CLUSTERMSG_TYPE_PONG ||
//		msgtype == CLUSTERMSG_TYPE_MEET {
//		 count := hdr.count
//		  /* expected length of this packet */
//
//		explen := uint16(unsafe.Sizeof(clusterMsg{})-unsafe.Sizeof(ping{}))
//		curLen := uint16(unsafe.Sizeof(clusterMsgDataGossip{}))
//		explen += curLen*count
//		if totlen != uint32(explen){return nil}}
//	//} else if msgtype == CLUSTERMSG_TYPE_FAIL {
//	//	uint32_t explen = sizeof(clusterMsg)-sizeof(union clusterMsgData)
//	//
//	//	explen += sizeof(clusterMsgDataFail)
//	//	if (totlen != explen) return 1
//	//} else if (msgtype == CLUSTERMSG_TYPE_PUBLISH) {
//	//	uint32_t explen = sizeof(clusterMsg)-sizeof(union clusterMsgData)
//	//
//	//	explen += sizeof(clusterMsgDataPublish) +
//	//		ntohl(hdr.data.publish.msg.channel_len) +
//	//		ntohl(hdr.data.publish.msg.message_len)
//	//	if (totlen != explen) return 1
//	//} else if (msgtype == CLUSTERMSG_TYPE_FAILOVER_AUTH_REQUEST ||
//	//msgtype == CLUSTERMSG_TYPE_FAILOVER_AUTH_ACK ||
//	//msgtype == CLUSTERMSG_TYPE_MFSTART)
//	//{
//	//uint32_t explen = sizeof(clusterMsg)-sizeof(union clusterMsgData)
//	//
//	//if (totlen != explen) return 1
//	//} else if (msgtype == CLUSTERMSG_TYPE_UPDATE) {
//	//	uint32_t explen = sizeof(clusterMsg)-sizeof(union clusterMsgData)
//	//
//	//	explen += sizeof(clusterMsgDataUpdate)
//	//	if (totlen != explen) return 1
//	//}
//
//	/* Check if the sender is a known node. */
//	// 查找发送者节点
//	sender = clusterLookupNode(hdr.sender)
//	// 节点存在，并且不是 HANDSHAKE 节点
//	// 那么个更新节点的配置纪元信息
//	//if (sender && !nodeInHandshake(sender)) {
//	//	/* Update our curretEpoch if we see a newer epoch in the cluster. */
//	//	senderCurrentEpoch = ntohu64(hdr.currentEpoch)
//	//	senderConfigEpoch = ntohu64(hdr.configEpoch)
//	//	if (senderCurrentEpoch > server.cluster.currentEpoch)
//	//		server.cluster.currentEpoch = senderCurrentEpoch
//	//	/* Update the sender configEpoch if it is publishing a newer one. */
//	//	if (senderConfigEpoch > sender.configEpoch) {
//	//		sender.configEpoch = senderConfigEpoch
//	//		clusterDoBeforeSleep(CLUSTER_TODO_SAVE_CONFIG|
//	//			CLUSTER_TODO_FSYNC_CONFIG)
//	//	}
//	//	/* Update the replication offset info for this node. */
//	//	sender.repl_offset = ntohu64(hdr.offset)
//	//	sender.repl_offset_time = mstime()
//	//	/* If we are a slave performing a manual failover and our master
//	//	 * sent its offset while already paused, populate the MF state. */
//	//	if (server.cluster.mf_end &&
//	//		nodeIsSlave(myself) &&
//	//		myself.slaveof == sender &&
//	//		hdr.mflags[0] & CLUSTERMSG_FLAG0_PAUSED &&
//	//		server.cluster.mf_master_offset == 0)
//	//	{
//	//		server.cluster.mf_master_offset = sender.repl_offset
//	//		redisLog(REDIS_WARNING,
//	//			"Received replication offset for paused "
//	//		"master manual failover: %lld",
//	//			server.cluster.mf_master_offset)
//	//	}
//	//}
//
//	/* Process packets by msgtype. */
//	// 根据消息的类型，处理节点
//
//	// 这是一条 PING 消息或者 MEET 消息
//	if msgtype == CLUSTERMSG_TYPE_PING || msgtype == CLUSTERMSG_TYPE_MEET {
//		//redisLog(REDIS_DEBUG,"Ping packet received: %p", (void*)link.node)
//
//		/* Add this node if it is new for us and the msg msgtype is MEET.
//		 *
//		 * 如果当前节点是第一次遇见这个节点，并且对方发来的是 MEET 信息，
//		 * 那么将这个节点添加到集群的节点列表里面。
//		 *
//		 * In this stage we don't try to add the node with the right
//		 * flags, slaveof pointer, and so forth, as this details will be
//		 * resolved when we'll receive PONGs from the node.
//		 *
//		 * 节点目前的 flag 、 slaveof 等属性的值都是未设置的，
//		 * 等当前节点向对方发送 PING 命令之后，
//		 * 这些信息可以从对方回复的 PONG 信息中取得。
//		 */
//		if sender==nil && msgtype == CLUSTERMSG_TYPE_MEET {
//
//			// 创建 HANDSHAKE 状态的新节点
//			node := createClusterNode("", REDIS_NODE_HANDSHAKE)
//
//			// 设置 IP 和端口
//			conn := link.Conn
//			node.ip = conn.LocalAddr().String()
//			node.port, _ = strconv.Atoi( node.ip[strings.Index(node.ip,":"):])
//
//
//			// 将新节点添加到集群
//			GetServerInstance().Cluster.nodes[node.name]=node
//
//
//			//clusterDoBeforeSleep(CLUSTER_TODO_SAVE_CONFIG)
//		}
//
//		/* Get info from the gossip section */
//		// 分析并取出消息中的 gossip 节点信息
//		//clusterProcessGossipSection(hdr,link)
//
//		/* Anyway reply with a PONG */
//		// 向目标节点返回一个 PONG
//		clusterSendPing(link, CLUSTERMSG_TYPE_PONG)
//	}
//
//	/* PING or PONG: process config information. */
//	// 这是一条 PING 、 PONG 或者 MEET 消息
//	if msgtype == CLUSTERMSG_TYPE_PING || msgtype == CLUSTERMSG_TYPE_PONG ||
//	msgtype == CLUSTERMSG_TYPE_MEET {
//	//redisLog(REDIS_DEBUG,"%s packet received: %p",
//	//msgtype == CLUSTERMSG_TYPE_PING ? "ping" : "pong",
//	//(void*)link.node)
//
//	// 连接的 clusterNode 结构存在
//	if link.node!=nil {
//	// 节点处于 HANDSHAKE 状态
//	if (nodeInHandshake(link.node)) {
//	/* If we already have this node, try to change the
//	 * IP/port of the node with the new one. */
//	if (sender!=nil) {
//	//redisLog(REDIS_VERBOSE,
//	//"Handshake: we already know node %.40s, "
//	//"updating the address if needed.", sender.name)
//	// 如果有需要的话，更新节点的地址
//	//if nodeUpdateAddressIfNeeded(sender,link,(hdr.port))
//	//{
//	//clusterDoBeforeSleep(CLUSTER_TODO_SAVE_CONFIG|
//	//CLUSTER_TODO_UPDATE_STATE)
//	//}
//	///* Free this node as we alrady have it. This will
//	// * cause the link to be freed as well. */
//	//// 释放节点
//	//freeClusterNode(link.node)
//	//return 0
//	}
//
//	/* First thing to do is replacing the random name with the
//	 * right node name if this was a handshake stage. */
//	// 用节点的真名替换在 HANDSHAKE 时创建的随机名字
//	//clusterRenameNode(link.node, hdr.sender)
//	//redisLog(REDIS_DEBUG,"Handshake with node %.40s completed.",
//	//link.node.name)
//
//	// 关闭 HANDSHAKE 状态
//	link.node.flags &= ~REDIS_NODE_HANDSHAKE
//
//	// 设置节点的角色
//	link.node.flags |= flags&(REDIS_NODE_MASTER|REDIS_NODE_SLAVE)
//
//	clusterDoBeforeSleep(CLUSTER_TODO_SAVE_CONFIG)
//
//	// 节点已存在，但它的 id 和当前节点保存的 id 不同
//	} else if (memcmp(link.node.name,hdr.sender,
//	REDIS_CLUSTER_NAMELEN) != 0)
//	{
//	/* If the reply has a non matching node ID we
//	 * disconnect this node and set it as not having an associated
//	 * address. */
//	// 那么将这个节点设为 NOADDR
//	// 并断开连接
//	redisLog(REDIS_DEBUG,"PONG contains mismatching sender ID")
//	link.node.flags |= REDIS_NODE_NOADDR
//	link.node.ip[0] = '\0'
//	link.node.port = 0
//
//	// 断开连接
//	freeClusterLink(link)
//
//	clusterDoBeforeSleep(CLUSTER_TODO_SAVE_CONFIG)
//	return 0
//	}
//	}
//
//	/* Update the node address if it changed. */
//	// 如果发送的消息为 PING
//	// 并且发送者不在 HANDSHAKE 状态
//	// 那么更新发送者的信息
//	if (sender && msgtype == CLUSTERMSG_TYPE_PING &&
//	!nodeInHandshake(sender) &&
//	nodeUpdateAddressIfNeeded(sender,link,(hdr.port)))
//	{
//	clusterDoBeforeSleep(CLUSTER_TODO_SAVE_CONFIG|
//	CLUSTER_TODO_UPDATE_STATE)
//	}
//
//	/* Update our info about the node */
//	// 如果这是一条 PONG 消息，那么更新我们关于 node 节点的认识
//	if (link.node && msgtype == CLUSTERMSG_TYPE_PONG) {
//
//	// 最后一次接到该节点的 PONG 的时间
//	link.node.pong_received = mstime()
//
//	// 清零最近一次等待 PING 命令的时间
//	link.node.ping_sent = 0
//
//	/* The PFAIL condition can be reversed without external
//	 * help if it is momentary (that is, if it does not
//	 * turn into a FAIL state).
//	 *
//	 * 接到节点的 PONG 回复，我们可以移除节点的 PFAIL 状态。
//	 *
//	 * The FAIL condition is also reversible under specific
//	 * conditions detected by clearNodeFailureIfNeeded().
//	 *
//	 * 如果节点的状态为 FAIL ，
//	 * 那么是否撤销该状态要根据 clearNodeFailureIfNeeded() 函数来决定。
//	 */
//	if (nodeTimedOut(link.node)) {
//	// 撤销 PFAIL
//	link.node.flags &= ~REDIS_NODE_PFAIL
//
//	clusterDoBeforeSleep(CLUSTER_TODO_SAVE_CONFIG|
//	CLUSTER_TODO_UPDATE_STATE)
//	} else if (nodeFailed(link.node)) {
//	// 看是否可以撤销 FAIL
//	clearNodeFailureIfNeeded(link.node)
//	}
//	}
//
//	/* Check for role switch: slave . master or master . slave. */
//	// 检测节点的身份信息，并在需要时进行更新
//	if (sender) {
//
//	// 发送消息的节点的 slaveof 为 REDIS_NODE_NULL_NAME
//	// 那么 sender 就是一个主节点
//	if (!memcmp(hdr.slaveof,REDIS_NODE_NULL_NAME,
//	sizeof(hdr.slaveof)))
//	{
//	/* Node is a master. */
//	// 设置 sender 为主节点
//	clusterSetNodeAsMaster(sender)
//
//	// sender 的 slaveof 不为空，那么这是一个从节点
//	} else {
//
//	/* Node is a slave. */
//	// 取出 sender 的主节点
//	clusterNode *master = clusterLookupNode(hdr.slaveof)
//
//	// sender 由主节点变成了从节点，重新配置 sender
//	if (nodeIsMaster(sender)) {
//	/* Master turned into a slave! Reconfigure the node. */
//
//	// 删除所有由该节点负责的槽
//	clusterDelNodeSlots(sender)
//
//	// 更新标识
//	sender.flags &= ~REDIS_NODE_MASTER
//	sender.flags |= REDIS_NODE_SLAVE
//
//	/* Remove the list of slaves from the node. */
//	// 移除 sender 的从节点名单
//	if (sender.numslaves) clusterNodeResetSlaves(sender)
//
//	/* Update config and state. */
//	clusterDoBeforeSleep(CLUSTER_TODO_SAVE_CONFIG|
//	CLUSTER_TODO_UPDATE_STATE)
//	}
//
//	/* Master node changed for this slave? */
//
//	// 检查 sender 的主节点是否变更
//	if (master && sender.slaveof != master) {
//	// 如果 sender 之前的主节点不是现在的主节点
//	// 那么在旧主节点的从节点列表中移除 sender
//	if (sender.slaveof)
//	clusterNodeRemoveSlave(sender.slaveof,sender)
//
//	// 并在新主节点的从节点列表中添加 sender
//	clusterNodeAddSlave(master,sender)
//
//	// 更新 sender 的主节点
//	sender.slaveof = master
//
//	/* Update config. */
//	clusterDoBeforeSleep(CLUSTER_TODO_SAVE_CONFIG)
//	}
//	}
//	}
//
//	/* Update our info about served slots.
//	 *
//	 * 更新当前节点对 sender 所处理槽的认识。
//	 *
//	 * Note: this MUST happen after we update the master/slave state
//	 * so that REDIS_NODE_MASTER flag will be set.
//	 *
//	 * 这部分的更新 *必须* 在更新 sender 的主/从节点信息之后，
//	 * 因为这里需要用到 REDIS_NODE_MASTER 标识。
//	 */
//
//	/* Many checks are only needed if the set of served slots this
//	 * instance claims is different compared to the set of slots we have
//	 * for it. Check this ASAP to avoid other computational expansive
//	 * checks later. */
//	clusterNode *sender_master = NULL /* Sender or its master if slave. */
//	int dirty_slots = 0 /* Sender claimed slots don't match my view? */
//
//	if (sender) {
//	sender_master = nodeIsMaster(sender) ? sender : sender.slaveof
//	if (sender_master) {
//	dirty_slots = memcmp(sender_master.slots,
//	hdr.myslots,sizeof(hdr.myslots)) != 0
//	}
//	}
//
//	/* 1) If the sender of the message is a master, and we detected that
//	 *    the set of slots it claims changed, scan the slots to see if we
//	 *    need to update our configuration. */
//	// 如果 sender 是主节点，并且 sender 的槽布局出现了变动
//	// 那么检查当前节点对 sender 的槽布局设置，看是否需要进行更新
//	if (sender && nodeIsMaster(sender) && dirty_slots)
//	clusterUpdateSlotsConfigWith(sender,senderConfigEpoch,hdr.myslots)
//
//	/* 2) We also check for the reverse condition, that is, the sender
//	 *    claims to serve slots we know are served by a master with a
//	 *    greater configEpoch. If this happens we inform the sender.
//	 *
//	 *    检测和条件 1 的相反条件，也即是，
//	 *    sender 处理的槽的配置纪元比当前节点已知的某个节点的配置纪元要低，
//	 *    如果是这样的话，通知 sender 。
//	 *
//	 * This is useful because sometimes after a partition heals, a
//	 * reappearing master may be the last one to claim a given set of
//	 * hash slots, but with a configuration that other instances know to
//	 * be deprecated. Example:
//	 *
//	 * 这种情况可能会出现在网络分裂中，
//	 * 一个重新上线的主节点可能会带有已经过时的槽布局。
//	 *
//	 * 比如说：
//	 *
//	 * A and B are master and slave for slots 1,2,3.
//	 * A 负责槽 1 、 2 、 3 ，而 B 是 A 的从节点。
//	 *
//	 * A is partitioned away, B gets promoted.
//	 * A 从网络中分裂出去，B 被提升为主节点。
//	 *
//	 * B is partitioned away, and A returns available.
//	 * B 从网络中分裂出去， A 重新上线（但是它所使用的槽布局是旧的）。
//	 *
//	 * Usually B would PING A publishing its set of served slots and its
//	 * configEpoch, but because of the partition B can't inform A of the
//	 * new configuration, so other nodes that have an updated table must
//	 * do it. In this way A will stop to act as a master (or can try to
//	 * failover if there are the conditions to win the election).
//	 *
//	 * 在正常情况下， B 应该向 A 发送 PING 消息，告知 A ，自己（B）已经接替了
//	 * 槽 1、 2、 3 ，并且带有更更的配置纪元，但因为网络分裂的缘故，
//	 * 节点 B 没办法通知节点 A ，
//	 * 所以通知节点 A 它带有的槽布局已经更新的工作就交给其他知道 B 带有更高配置纪元的节点来做。
//	 * 当 A 接到其他节点关于节点 B 的消息时，
//	 * 节点 A 就会停止自己的主节点工作，又或者重新进行故障转移。
//	 */
//	if (sender && dirty_slots) {
//	int j
//
//	for (j = 0 j < REDIS_CLUSTER_SLOTS j++) {
//
//	// 检测 slots 中的槽 j 是否已经被指派
//	if (bitmapTestBit(hdr.myslots,j)) {
//
//	// 当前节点认为槽 j 由 sender 负责处理，
//	// 或者当前节点认为该槽未指派，那么跳过该槽
//	if (server.cluster.slots[j] == sender ||
//	server.cluster.slots[j] == NULL) continue
//
//	// 当前节点槽 j 的配置纪元比 sender 的配置纪元要大
//	if (server.cluster.slots[j].configEpoch >
//	senderConfigEpoch)
//	{
//	redisLog(REDIS_VERBOSE,
//	"Node %.40s has old slots configuration, sending "
//	"an UPDATE message about %.40s",
//	sender.name, server.cluster.slots[j].name)
//
//	// 向 sender 发送关于槽 j 的更新信息
//	clusterSendUpdate(sender.link,
//	server.cluster.slots[j])
//
//	/* TODO: instead of exiting the loop send every other
//	 * UPDATE packet for other nodes that are the new owner
//	 * of sender's slots. */
//	break
//	}
//	}
//	}
//	}
//
//	/* If our config epoch collides with the sender's try to fix
//	 * the problem. */
//	if (sender &&
//	nodeIsMaster(myself) && nodeIsMaster(sender) &&
//	senderConfigEpoch == myself.configEpoch)
//	{
//	clusterHandleConfigEpochCollision(sender)
//	}
//
//	/* Get info from the gossip section */
//	// 分析并提取出消息 gossip 协议部分的信息
//	clusterProcessGossipSection(hdr,link)
//
//	// 这是一条 FAIL 消息： sender 告知当前节点，某个节点已经进入 FAIL 状态。
//	} else if (msgtype == CLUSTERMSG_TYPE_FAIL) {
//		clusterNode *failing
//
//		if (sender) {
//
//			// 获取下线节点的消息
//			failing = clusterLookupNode(hdr.data.fail.about.nodename)
//			// 下线的节点既不是当前节点，也没有处于 FAIL 状态
//			if (failing &&
//				!(failing.flags & (REDIS_NODE_FAIL | REDIS_NODE_MYSELF)))
//			{
//				redisLog(REDIS_NOTICE,
//					"FAIL message received from %.40s about %.40s",
//					hdr.sender, hdr.data.fail.about.nodename)
//
//				// 打开 FAIL 状态
//				failing.flags |= REDIS_NODE_FAIL
//				failing.fail_time = mstime()
//				// 关闭 PFAIL 状态
//				failing.flags &= ~REDIS_NODE_PFAIL
//				clusterDoBeforeSleep(CLUSTER_TODO_SAVE_CONFIG|
//					CLUSTER_TODO_UPDATE_STATE)
//			}
//		} else {
//			redisLog(REDIS_NOTICE,
//				"Ignoring FAIL message from unknonw node %.40s about %.40s",
//				hdr.sender, hdr.data.fail.about.nodename)
//		}
//
//		// 这是一条 PUBLISH 消息
//	} else if (msgtype == CLUSTERMSG_TYPE_PUBLISH) {
//		robj *channel, *message
//		uint32_t channel_len, message_len
//
//		/* Don't bother creating useless objects if there are no
//		 * Pub/Sub subscribers. */
//		// 只在有订阅者时创建消息对象
//		if (dictSize(server.pubsub_channels) ||
//			listLength(server.pubsub_patterns))
//		{
//			// 频道长度
//			channel_len = ntohl(hdr.data.publish.msg.channel_len)
//
//			// 消息长度
//			message_len = ntohl(hdr.data.publish.msg.message_len)
//
//			// 频道
//			channel = createStringObject(
//				(char*)hdr.data.publish.msg.bulk_data,channel_len)
//
//			// 消息
//			message = createStringObject(
//				(char*)hdr.data.publish.msg.bulk_data+channel_len,
//			message_len)
//			// 发送消息
//			pubsubPublishMessage(channel,message)
//
//			decrRefCount(channel)
//			decrRefCount(message)
//		}
//
//		// 这是一条请求获得故障迁移授权的消息： sender 请求当前节点为它进行故障转移投票
//	} else if (msgtype == CLUSTERMSG_TYPE_FAILOVER_AUTH_REQUEST) {
//		if (!sender) return 1  /* We don't know that node. */
//		// 如果条件允许的话，向 sender 投票，支持它进行故障转移
//		clusterSendFailoverAuthIfNeeded(sender,hdr)
//
//		// 这是一条故障迁移投票信息： sender 支持当前节点执行故障转移操作
//	} else if (msgtype == CLUSTERMSG_TYPE_FAILOVER_AUTH_ACK) {
//		if (!sender) return 1  /* We don't know that node. */
//
//		/* We consider this vote only if the sender is a master serving
//		 * a non zero number of slots, and its currentEpoch is greater or
//		 * equal to epoch where this node started the election. */
//		// 只有正在处理至少一个槽的主节点的投票会被视为是有效投票
//		// 只有符合以下条件， sender 的投票才算有效：
//		// 1） sender 是主节点
//		// 2） sender 正在处理至少一个槽
//		// 3） sender 的配置纪元大于等于当前节点的配置纪元
//		if (nodeIsMaster(sender) && sender.numslots > 0 &&
//			senderCurrentEpoch >= server.cluster.failover_auth_epoch)
//		{
//			// 增加支持票数
//			server.cluster.failover_auth_count++
//
//			/* Maybe we reached a quorum here, set a flag to make sure
//			 * we check ASAP. */
//			clusterDoBeforeSleep(CLUSTER_TODO_HANDLE_FAILOVER)
//		}
//
//	} else if (msgtype == CLUSTERMSG_TYPE_MFSTART) {
//		/* This message is acceptable only if I'm a master and the sender
//		 * is one of my slaves. */
//		if (!sender || sender.slaveof != myself) return 1
//		/* Manual failover requested from slaves. Initialize the state
//		 * accordingly. */
//		resetManualFailover()
//		server.cluster.mf_end = mstime() + REDIS_CLUSTER_MF_TIMEOUT
//		server.cluster.mf_slave = sender
//		pauseClients(mstime()+(REDIS_CLUSTER_MF_TIMEOUT*2))
//		redisLog(REDIS_WARNING,"Manual failover requested by slave %.40s.",
//			sender.name)
//	} else if (msgtype == CLUSTERMSG_TYPE_UPDATE) {
//		clusterNode *n /* The node the update is about. */
//		uint64_t reportedConfigEpoch =
//			ntohu64(hdr.data.update.nodecfg.configEpoch)
//
//		if (!sender) return 1  /* We don't know the sender. */
//
//		// 获取需要更新的节点
//		n = clusterLookupNode(hdr.data.update.nodecfg.nodename)
//		if (!n) return 1   /* We don't know the reported node. */
//
//		// 消息的纪元并不大于节点 n 所处的配置纪元
//		// 无须更新
//		if (n.configEpoch >= reportedConfigEpoch) return 1 /* Nothing new. */
//
//		/* If in our current config the node is a slave, set it as a master. */
//		// 如果节点 n 为从节点，但它的槽配置更新了
//		// 那么说明这个节点已经变为主节点，将它设置为主节点
//		if (nodeIsSlave(n)) clusterSetNodeAsMaster(n)
//
//		/* Update the node's configEpoch. */
//		n.configEpoch = reportedConfigEpoch
//		clusterDoBeforeSleep(CLUSTER_TODO_SAVE_CONFIG|
//			CLUSTER_TODO_FSYNC_CONFIG)
//
//		/* Check the bitmap of served slots and udpate our
//		 * config accordingly. */
//		// 将消息中对 n 的槽布局与当前节点对 n 的槽布局进行对比
//		// 在有需要时更新当前节点对 n 的槽布局的认识
//		clusterUpdateSlotsConfigWith(n,reportedConfigEpoch,
//			hdr.data.update.nodecfg.slots)
//	} else {
//		redisLog(REDIS_WARNING,"Received unknown packet msgtype: %d", msgtype)
//	}
//	return 1
//	return nil
//}
//
//func clusterLookupNode(name string)*Node {
//	return GetServerInstance().Cluster.nodes[name]
//}
//
///* -----------------------------------------------------------------------------
// * CLUSTER node API
// * -------------------------------------------------------------------------- */
//
///* Create a new cluster node, with the specified flags.
// *
// * 创建一个带有指定 flag 的集群节点。
// *
// * If "nodename" is NULL this is considered a first handshake and a random
// * node name is assigned to this node (it will be fixed later when we'll
// * receive the first pong).
// *
// * 如果 nodename 参数为 NULL ，那么表示我们尚未向节点发送 PING ，
// * 集群会为节点设置一个随机的命令，
// * 这个命令在之后接收到节点的 PONG 回复之后就会被更新。
// *
// * The node is created and returned to the user, but it is not automatically
// * added to the nodes hash table.
// *
// * 函数会返回被创建的节点，但不会自动将它添加到当前节点的节点哈希表中
// * （nodes hash table）。
// */
//func createClusterNode(nodename string,  flags uint16) *Node {
//
//	return &Node{
//		createTime:   time.Now().Unix(),
//		name:         nodename,
//		flags:        flags,
//		configEpoch:  0,
//		slots:        [2048]byte{},
//		numslots:     0,
//		numslaves:    0,
//		slave:        nil,
//		slaveof:      nil,
//		port:         0,
//		cport:        0,
//		link:         nil,
//		inbound_link: nil,
//	}
//
//
//}
//
///* Send a PING or PONG packet to the specified node, making sure to add enough
// * gossip informations. */
//// 向指定节点发送一条 MEET 、 PING 或者 PONG 消息
//func clusterSendPing(link *Link, msgtype int) {
//	//buf :=make([]byte,unsafe.Sizeof(clusterMsg{}))
//	//hdr :=(*clusterMsg)(unsafe.Pointer(
//	//	(*reflect.SliceHeader)(unsafe.Pointer(&buf)).Data,
//	//))
//	//
//	//var  gossipcount int = 0
//	//var totlen int
///* freshnodes is the number of nodes we can still use to populate the
// * gossip section of the ping packet. Basically we start with the nodes
// * we have in memory minus two (ourself and the node we are sending the
// * message to). Every time we add a node we decrement the counter, so when
// * it will drop to <= zero we know there is no more gossip info we can
// * send. */
//// freshnodes 是用于发送 gossip 信息的计数器
//// 每次发送一条信息时，程序将 freshnodes 的值减一
//// 当 freshnodes 的数值小于等于 0 时，程序停止发送 gossip 信息
//// freshnodes 的数量是节点目前的 nodes 表中的节点数量减去 2
//// 这里的 2 指两个节点，一个是 myself 节点（也即是发送信息的这个节点）
//// 另一个是接受 gossip 信息的节点
//// freshnodes := len(redisServer.Cluster.nodes)-2
//
//// 如果发送的信息是 PING ，那么更新最后一次发送 PING 命令的时间戳
////if link.node!=nil && msgtype == CLUSTERMSG_TYPE_PING{
////		link.node. = mstime()
////	}
//
//
//// 将当前节点的信息（比如名字、地址、端口号、负责处理的槽）记录到消息里面
//	//clusterBuildMessageHdr(hdr,msgtype)
//	//
//	///* Populate the gossip fields */
//	//// 从当前节点已知的节点中随机选出两个节点
//	//// 并通过这条消息捎带给目标节点，从而实现 gossip 协议
//	//
//	//// 每个节点有 freshnodes 次发送 gossip 信息的机会
//	//// 每次向目标节点发送 2 个被选中节点的 gossip 信息（gossipcount 计数）
//	//for ;freshnodes > 0 && gossipcount < 3; {
//	//// 从 nodes 字典中随机选出一个节点（被选中节点）
//	//dictEntry *de = dictGetRandomKey(server.cluster->nodes)
//	//clusterNode *this = dictGetVal(de)
//	//
//	//clusterMsgDataGossip *gossip
//	//int j
//
///* In the gossip section don't include:
// * 以下节点不能作为被选中节点：
// * 1) Myself.
// *    节点本身。
// * 2) Nodes in HANDSHAKE state.
// *    处于 HANDSHAKE 状态的节点。
// * 3) Nodes with the NOADDR flag set.
// *    带有 NOADDR 标识的节点
// * 4) Disconnected nodes if they don't have configured slots.
// *    因为不处理任何槽而被断开连接的节点
// */
//	//if (this == myself ||
//	//this->flags & (REDIS_NODE_HANDSHAKE|REDIS_NODE_NOADDR) ||
//	//(this->link == NULL && this->numslots == 0))
//	//{
//	//freshnodes-- /* otherwise we may loop forever. */
//	//continue
//	//}
//	//
//	///* Check if we already added this node */
//	//// 检查被选中节点是否已经在 hdr->data.ping.gossip 数组里面
//	//// 如果是的话说明这个节点之前已经被选中了
//	//// 不要再选中它（否则就会出现重复）
//	//for j := 0 ;j < gossipcount; j++ {
//	//if (memcmp(hdr->data.ping.gossip[j].nodename,this->name,
//	//REDIS_CLUSTER_NAMELEN) == 0) {break}
//	//}
//	//if j != gossipcount continue
//	//
//	///* Add it */
//	//
//	//// 这个被选中节点有效，计数器减一
//	//freshnodes--
//	//
//	//// 指向 gossip 信息结构
//	//gossip = &(hdr->data.ping.gossip[gossipcount])
//	//
//	//// 将被选中节点的名字记录到 gossip 信息
//	//memcpy(gossip->nodename,this->name,REDIS_CLUSTER_NAMELEN)
//	//// 将被选中节点的 PING 命令发送时间戳记录到 gossip 信息
//	//gossip->ping_sent = htonl(this->ping_sent)
//	//// 将被选中节点的 PING 命令回复的时间戳记录到 gossip 信息
//	//gossip->pong_received = htonl(this->pong_received)
//	//// 将被选中节点的 IP 记录到 gossip 信息
//	//memcpy(gossip->ip,this->ip,sizeof(this->ip))
//	//// 将被选中节点的端口号记录到 gossip 信息
//	//gossip->port = htons(this->port)
//	//// 将被选中节点的标识值记录到 gossip 信息
//	//gossip->flags = htons(this->flags)
//	//
//	//// 这个被选中节点有效，计数器增一
//	//gossipcount++
//	//}
//	//
//	//// 计算信息长度
//	//totlen = sizeof(clusterMsg)-sizeof(union clusterMsgData)
//	//totlen += (sizeof(clusterMsgDataGossip)*gossipcount)
//	//// 将被选中节点的数量（gossip 信息中包含了多少个节点的信息）
//	//// 记录在 count 属性里面
//	//hdr->count = htons(gossipcount)
//	//// 将信息的长度记录到信息里面
//	//hdr->totlen = htonl(totlen)
//
//// 发送信息
////clusterSendMessage(link,buf,totlen)
//}
//
///* Put stuff into the send buffer.
// *
// * 发送信息
// *
// * It is guaranteed that this function will never have as a side effect
// * the link to be invalidated, so it is safe to call this function
// * from event handlers that will do stuff with the same link later.
// *
// * 因为发送不会对连接本身造成不良的副作用，
// * 所以可以在发送信息的处理器上做一些针对连接本身的动作。
// */
//func clusterSendMessage(link *Link,  msg []byte,  msglen int) {
//	link.sendBuf = append(link.sendBuf,msg...)
//	_,err:=link.Conn.Write(link.sendBuf)
//	if err!=nil{
//		return
//	}
//// 安装写事件处理器
//}
//
///* Send a message to all the nodes that are part of the cluster having
// * a connected link.
// *
// * 向节点连接的所有其他节点发送信息。
// *
// * It is guaranteed that this function will never have as a side effect
// * some node->link to be invalidated, so it is safe to call this function
// * from event handlers that will do stuff with node links later. */
////func clusterBroadcastMessage(void *buf, size_t len) {
////	dictIterator *di
////	dictEntry *de
////
////// 遍历所有已知节点
////di = dictGetSafeIterator(server.cluster->nodes)
////while((de = dictNext(di)) != NULL) {
////clusterNode *node = dictGetVal(de)
////
////// 不向未连接节点发送信息
////if (!node->link) continue
////
////// 不向节点自身或者 HANDSHAKE 状态的节点发送信息
////if (node->flags & (REDIS_NODE_MYSELF|REDIS_NODE_HANDSHAKE))
////continue
////
////// 发送信息
////clusterSendMessage(node->link,buf,len)
////}
////
////}
//
///* Build the message header */
//// 构建信息
//func clusterBuildMessageHdr(hdr *clusterMsg, msgtype int) {
// totlen :=uint32( 0)
////var  offset uint64
//var master *Node
//
///* If this node is a master, we send its slots bitmap and configEpoch.
// *
// * 如果这是一个主节点，那么发送该节点的槽 bitmap 和配置纪元。
// *
// * If this node is a slave we send the master's information instead (the
// * node is flagged as slave so the receiver knows that it is NOT really
// * in charge for this slots.
// * 如果这是一个从节点，
// * 那么发送这个节点的主节点的槽 bitmap 和配置纪元。
// *
// * 因为接收信息的节点通过标识可以知道这个节点是一个从节点，
// * 所以接收信息的节点不会将从节点错认作是主节点。
// */
//	if nodeIsSlave(myself) && myself.slaveof!=nil {
//		master=myself.slaveof
//	}else{
//		master=myself
//	}
//
//	hdr.sig[0] = 'R'
//	hdr.sig[1] = 'C'
//	hdr.sig[2] = 'm'
//	hdr.sig[3] = 'b'
//
//// 设置信息类型
//hdr.messagetype = uint16(msgtype)
//	// 设置信息发送者
//hdr.sender = myself.name
//
//
//// 设置当前节点负责的槽
//	hdr.myslots = master.slots
//
//// 清零 slaveof 域
//
//
//
//// 如果节点是从节点的话，那么设置 slaveof 域
//if myself.slaveof != nil {
//	hdr.slaveof = myself.slaveof.name
//}
////memcpy(hdr->slaveof,myself->slaveof->name, REDIS_CLUSTER_NAMELEN)
//
//// 设置端口号
//hdr.port = redisServer.port
//
//// 设置标识
//hdr.flags = myself.flags
//
//// 设置状态
//hdr.state = redisServer.Cluster.state
//
/////* Set the currentEpoch and configEpochs. */
////// 设置集群当前配置纪元
////hdr.currentEpoch = htonu64(server.cluster->currentEpoch)
////// 设置主节点当前配置纪元
////hdr->configEpoch = htonu64(master->configEpoch)
//
///* Set the replication offset. */
//// 设置复制偏移量
////if (nodeIsSlave(myself))
////offset = replicationGetSlaveOffset()
////else
////offset = server.master_repl_offset
////hdr->offset = htonu64(offset)
//
///* Set the message flags. */
////if (nodeIsMaster(myself) && server.cluster->mf_end)
////hdr->mflags[0] |= CLUSTERMSG_FLAG0_PAUSED
//
///* Compute the message length for certain messages. For other messages
// * this is up to the caller. */
//// 计算信息的长度
//if msgtype == CLUSTERMSG_TYPE_FAIL {
//totlen = uint32(unsafe.Sizeof(clusterMsg{})-unsafe.Sizeof(ping{}))
////totlen += unsafe.Sizeof(clusterMsgDataFail)
//} else if msgtype == CLUSTERMSG_TYPE_UPDATE {
//totlen = uint32(unsafe.Sizeof(clusterMsg{})-unsafe.Sizeof(ping{}))
////totlen += sizeof(clusterMsgDataUpdate)
//}
//
//// 设置信息的长度
//hdr.totlen = uint32(totlen)
///* For PING, PONG, and MEET, fixing the totlen field is up to the caller. */
//}
//var myself *Node
//// 用于判断节点身份和状态的一系列宏
//func nodeIsMaster(n *Node) bool {return (n.flags & REDIS_NODE_MASTER)==1}
//func nodeIsSlave(n *Node) bool {return (n.flags & REDIS_NODE_SLAVE)==1 }
//func nodeInHandshake(n *Node)  bool {return (n.flags & REDIS_NODE_HANDSHAKE)==1 }
//func nodeHasAddr(n *Node) bool {return (n.flags & REDIS_NODE_NOADDR)==0 }
//func nodeWithoutAddr(n *Node) bool { return (n.flags & REDIS_NODE_NOADDR)==1 }
//func nodeTimedOut(n *Node) bool {return (n.flags & REDIS_NODE_PFAIL)==1 }
//func nodeFailed(n *Node)bool {return (n.flags & REDIS_NODE_FAIL)==1 }