package parrotserver

/* redis 集群的总槽位数量 16384 */
const CLUSTER_SLOTS = 16384
const CLUSTERMSG_TYPE_COUNT = 11
const CLUSTER_PROTO_VER = 1 /* Cluster bus protocol version. */

/* 集群在线 */
const CLUSTER_OK = 0 /* Everything looks ok */
/* 集群失效 */
const CLUSTER_FAIL = 1 /* The cluster can't work */
/* 集群节点名字长度 */
const CLUSTER_NAMELEN = 40 /* sha1 hex length */
/* 默认情况下；集群实际通信的端口 = 用户指定端口 + CLUSTER_PORT_INCR （6379 + 10000 = 16379）*/
const CLUSTER_PORT_INCR = 10000 /* Cluster port = baseport + PORT_INCR */

/* The following defines are amount of time, sometimes expressed as
 * multiplicators of the node timeout value (when ending with MULT). */
/* 下面是和时间相关的一些常量，以 _MULT 结尾的常量会作为时间值的乘法因子来使用 */
/* 节点故障报告的乘法因子 */
const CLUSTER_FAIL_REPORT_VALIDITY_MULT = 2 /* Fail report validity. */
/* 撤销主节点 FAIL 状态的乘法因子 */
const CLUSTER_FAIL_UNDO_TIME_MULT = 2 /* Undo fail if master is back. */
/* 在进行手动故障转移之前需要等待的超时时间 */
const CLUSTER_MF_TIMEOUT = 5000            /* Milliseconds to do a manual failover. */
const CLUSTER_MF_PAUSE_MULT = 2            /* Master pause manual failover mult. */
const CLUSTER_SLAVE_MIGRATION_DELAY = 5000 /* Delay for slave migration. */

/* Redirection errors returned by getNodeByQuery(). */
/* 下面的标识是节点之间做槽位转移的时候，客户端需要重定向节点，服务器返回给客户端的标识 */
/* 当前节点可以处理这个命令 */
const CLUSTER_REDIR_NONE = 0 /* Node can serve the request. */
/* 所请求的键在其他槽 */
const CLUSTER_REDIR_CROSS_SLOT = 1 /* -CROSSSLOT request. */
/* 键所处的槽正在进行 rehash */
const CLUSTER_REDIR_UNSTABLE = 2 /* -TRYAGAIN redirection required */
/* 需要进行 ASK 重定向 */
const CLUSTER_REDIR_ASK = 3 /* -ASK redirection required. */
/* 需要进行 MOVED 重定向 */
const CLUSTER_REDIR_MOVED = 4 /* -MOVED redirection required. */
/* 如果集群状态不是 OK 状态 */
const CLUSTER_REDIR_DOWN_STATE = 5 /* -CLUSTERDOWN, global state. */
/* 当前节点未分配槽位 */
const CLUSTER_REDIR_DOWN_UNBOUND = 6 /* -CLUSTERDOWN, unbound slot. */
/* 当前节点仅允许读 */
const CLUSTER_REDIR_DOWN_RO_STATE = 7 /* -CLUSTERDOWN, allow reads. */

/* Note that the PING, PONG and MEET messages are actually the same exact
 * kind of packet. PONG is the reply to ping, in the exact format as a PING,
 * while MEET is a special PING that forces the receiver to add the sender
 * as a node (if it is not already in the list). */
// 注意，PING 、 PONG 和 MEET 实际上是同一种消息。
// PONG 是对 PING 的回复，它的实际格式也为 PING 消息，
// 而 MEET 则是一种特殊的 PING 消息，用于强制消息的接收者将消息的发送者添加到集群中
// （如果节点尚未在节点列表中的话）
// PING
const CLUSTERMSG_TYPE_PING = 0 /* Ping */
// PONG （回复 PING）
const CLUSTERMSG_TYPE_PONG = 1 /* Pong (reply to Ping) */
// 请求将某个节点添加到集群中
const CLUSTERMSG_TYPE_MEET = 2 /* Meet "let's join" message */
// 将某个节点标记为 FAIL
const CLUSTERMSG_TYPE_FAIL = 3 /* Mark node xxx as failing */
// 通过发布与订阅功能广播消息
const CLUSTERMSG_TYPE_PUBLISH = 4 /* Pub/Sub Publish propagation */
// 请求进行故障转移操作，要求消息的接收者通过投票来支持消息的发送者
const CLUSTERMSG_TYPE_FAILOVER_AUTH_REQUEST = 5 /* May I failover? */
// 消息的接收者同意向消息的发送者投票
const CLUSTERMSG_TYPE_FAILOVER_AUTH_ACK = 6 /* Yes, you have my vote */
// 槽布局已经发生变化，消息发送者要求消息接收者进行相应的更新
const CLUSTERMSG_TYPE_UPDATE = 7 /* Another node slots configuration */
// 为了进行手动故障转移，暂停各个客户端
const CLUSTERMSG_TYPE_MFSTART = 8 /* Pause clients for manual failover */

/* Cluster node flags and macros. */
// 该节点为主节点
const REDIS_NODE_MASTER =1     /* The node is a master */
// 该节点为从节点
const REDIS_NODE_SLAVE= 2      /* The node is a slave */
// 该节点疑似下线，需要对它的状态进行确认
const REDIS_NODE_PFAIL= 4      /* Failure? Need acknowledge */
// 该节点已下线
const REDIS_NODE_FAIL =8       /* The node is believed to be malfunctioning */
// 该节点是当前节点自身
const REDIS_NODE_MYSELF =16    /* This node is myself */
// 该节点还未与当前节点完成第一次 PING - PONG 通讯
const REDIS_NODE_HANDSHAKE= 32 /* We have still to exchange the first ping */
// 该节点没有地址
const REDIS_NODE_NOADDR =  64  /* We don't know the address of this node */
// 当前节点还未与该节点进行过接触
// 带有这个标识会让当前节点发送 MEET 命令而不是 PING 命令
const REDIS_NODE_MEET= 128     /* Send a MEET message to this node */
// 该节点被选中为新的主节点
const REDIS_NODE_PROMOTED =256 /* Master was a slave propoted by failover */
// 空名字（在节点为主节点时，用作消息中的 slaveof 属性的值）
const REDIS_NODE_NULL_NAME=  "\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000"
