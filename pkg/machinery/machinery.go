package machinery

import (
	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	machinery_config "github.com/RichardKnop/machinery/v2/config"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/pkg/machinery/tasks/oceanengine"
)

type Machinery struct {
	worker      *machinery.Worker
	server      *machinery.Server
	oceanengine *oceanengine.OceanEngineOpenSdk
	notify      chan error
}

// NewMachinery 初始化Machinery实例
// 创建Machinery服务器，注册任务，并返回一个工作者实例
// 参数:
//   - cfg: 应用程序配置
//
// 返回:
//   - worker: 初始化后的工作者实例，可用于启动任务处理
//   - err: 初始化过程中发生的错误
func NewMachinery(cfg *config.Config, oceanengine *oceanengine.OceanEngineOpenSdk) (m *Machinery, err error) {

	m = &Machinery{
		notify: make(chan error, 1),
	}

	m.server, err = InitRedisServer(cfg)
	if err != nil {
		return nil, err
	}

	m.worker = m.server.NewWorker("worker_1", 10)

	// 注册任务处理函数
	m.RegisterTask()

	return m, nil
}

// InitRedisServer 初始化基于Redis的Machinery服务器实例
// 配置并创建任务服务器，包括消息代理、结果后端和锁机制
// 参数:
//   - cfg: 应用程序配置
//
// 返回:
//   - *machinery.Server: 初始化后的服务器实例
//   - error: 初始化过程中发生的错误
func InitRedisServer(cfg *config.Config) (*machinery.Server, error) {
	cnf := GetConfig(cfg)

	// 创建 Redis 代理
	// Redis代理负责任务消息的发送和接收
	broker := redisbroker.New(cnf, cnf.Broker, cfg.Machinery.RedisPassword, "", cfg.Machinery.RedisDB)

	// 创建 Redis 结果后端
	// 结果后端用于存储任务执行的状态和结果
	resultBackend := redis.New(cnf, cnf.ResultBackend, "", "", 0)

	// 创建锁
	// 使用本地内存锁实现，用于防止任务重复执行
	// 在分布式环境中可能需要替换为分布式锁实现
	lock := eagerlock.New()

	// 创建服务器实例
	return machinery.NewServer(cnf, broker, resultBackend, lock), nil
}

// GetConfig 返回默认的Machinery配置
// 根据应用程序配置创建Machinery所需的配置对象
func GetConfig(cfg *config.Config) *machinery_config.Config {
	return &machinery_config.Config{
		// 设置默认队列名称，用于任务发布时的默认目标队列
		DefaultQueue: cfg.Machinery.DefaultQueue,
		// 设置结果后端地址，用于存储任务执行结果
		ResultBackend: cfg.Machinery.RedisHost + ":" + cfg.Machinery.RedisPort,
		// 设置消息代理地址，用于任务队列管理
		Broker: cfg.Machinery.RedisHost + ":" + cfg.Machinery.RedisPort,
		// Redis 具体配置
		Redis: &machinery_config.RedisConfig{
			// 连接池中最大空闲连接数
			MaxIdle: 3,
			// 空闲连接超时时间（秒）
			IdleTimeout: 240,
			// 读取超时时间（秒）
			ReadTimeout: 15,
			// 写入超时时间（秒）
			WriteTimeout: 15,
			// 连接超时时间（秒）
			ConnectTimeout: 15,
			// 普通任务轮询周期（毫秒）
			NormalTasksPollPeriod: 1000,
			// 延迟任务轮询周期（毫秒）
			DelayedTasksPollPeriod: 500,
			// 延迟任务在Redis中的键名前缀
			DelayedTasksKey: cfg.Machinery.RedisKey,
		},
	}
}

// StartWorker 启动工作者
func (m *Machinery) StartWorker() {
	go func() {
		m.notify <- m.worker.Launch()
		close(m.notify)
	}()
}

// Notify 返回通知通道
func (m *Machinery) Notify() <-chan error {
	return m.notify
}
