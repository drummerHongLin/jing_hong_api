package emailservice

import (
	"jonghong/internal/pkg/errno"
	"jonghong/internal/pkg/log"
	"sync"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

// 定义邮件服务对象
type mailService struct {
	// 邮件服务器相关信息
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string

	// 协程相关信息
	taskQueue  chan *MailMessage //最大等待任务队列
	stopQueue  chan struct{}     // 停止信息
	semaphore  chan struct{}     // 最大并发数
	wg         sync.WaitGroup    // 保证协程的运行状态完毕
	running    bool              // 服务状态信息
	runningMux sync.Mutex        // 服务状态转变的互斥锁

	// 验证码相关
}

type codeInfo struct {
	code      string
	expiredAt time.Time
}

// 定义邮件消息对象
type MailMessage struct {
	From    string
	To      []string
	Subject string
	Body    string
}

// 定义邮件服务的功能接口
type MailService interface {
	Start() error
	Stop() error
	SendEmailAsync(msg *MailMessage) error
}

// 创建全局单例对象
var (
	MS   MailService
	once sync.Once
)

func InitMailService() error {
	once.Do(
		func() {
			MS = &mailService{
				smtpHost:     viper.GetString("smtp.host"),
				smtpPort:     viper.GetInt("smtp.port"),
				smtpUsername: viper.GetString("smtp.username"),
				smtpPassword: viper.GetString("smtp.password"),
			}
		})
	MS.Start()
	return nil
}

// 1. 同步发送方法
func (ms *mailService) sendMail(msg *MailMessage) error {
	// 构建gomail包的信息对象
	m := gomail.NewMessage()

	from := msg.From
	if from == "" {
		from = ms.smtpUsername
	}
	m.SetHeader("From", from)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)

	d := gomail.NewDialer(ms.smtpHost, ms.smtpPort, ms.smtpUsername, ms.smtpPassword)

	err := d.DialAndSend(m)

	return err
}

// 2. 异步发送方法
func (ms *mailService) SendEmailAsync(msg *MailMessage) error {
	// 同时间一个消息加入消息队列
	ms.runningMux.Lock()
	defer ms.runningMux.Unlock()

	if !ms.running {
		return errno.InternalServerError.SetMessage("邮件服务未启动，发送邮件失败")
	}
	// 将消息加入消息队列
	ms.taskQueue <- msg
	return nil
}

// 3. 任务调度器
func (ms *mailService) schedule() {
	// 无线循环监听
	for {
		// select 语句会监听通道上的操作，如果没有default，且一个通道操作都没，则会堵塞
		select {
		case task, ok := <-ms.taskQueue:
			// 任务通道关闭 退出监听循环
			if !ok {
				return
			}
			// 获取并发信号量, 如果并发消息数超过5个会被堵塞
			ms.semaphore <- struct{}{}
			// 等待任务对象增加一个，确保发送邮件对象执行完成
			ms.wg.Add(1)
			// 启动协程发送邮件
			go func(msg *MailMessage) {
				// 释放等待信号
				defer ms.wg.Done()
				// 释放并发信号
				defer func() { <-ms.semaphore }()
				if err := ms.sendMail(msg); err != nil {
					log.Errorw("邮件发送失败: ", msg.Subject, err)
				} else {
					log.Infow("邮件发送成功:", "主题", msg.Subject)
				}

			}(task)
		case <-ms.stopQueue:
			// 关闭通道
			// 遍历剩余任务
			for task := range ms.taskQueue {
				ms.semaphore <- struct{}{}
				ms.wg.Add(1)
				go func(msg *MailMessage) {
					defer ms.wg.Done()
					defer func() { <-ms.semaphore }()
					if err := ms.sendMail(msg); err != nil {
						log.Errorw("邮件发送失败: ", msg.Subject, err)
					} else {
						log.Infow("邮件发送成功:", "主题", msg.Subject)
					}

				}(task)
			}
			// 退出循环监听
			return
		}
	}
}

// 4. 开启邮件服务
func (ms *mailService) Start() error {
	// 开启锁
	ms.runningMux.Lock()
	defer ms.runningMux.Unlock()
	// 如果正在运行 直接返回
	if ms.running {
		return nil
	}
	// 初始化异步操作通道
	ms.taskQueue = make(chan *MailMessage, 100)
	ms.semaphore = make(chan struct{}, 5)
	ms.stopQueue = make(chan struct{})

	ms.running = true

	// 启动任务调度
	go ms.schedule()

	return nil

}

// 5. 停止邮件服务
func (ms *mailService) Stop() error {
	// 开启锁
	ms.runningMux.Lock()
	defer ms.runningMux.Unlock()

	if !ms.running {
		return nil
	}

	// 先改变运行状态吧？ 停止接受消息发送
	ms.running = false

	// 发送停止信号
	close(ms.stopQueue)

	// 等待任务完成

	ms.wg.Wait()

	// 清理验证码，删除定时器

	// 关停消息队列
	close(ms.taskQueue)

	log.Infow("邮件服务已关闭...")

	return nil

}
