package server

import "fmt"

const (
	DefaultBindInterface = "0.0.0.0"
	DefaultBindPort      = 3000
	DefaultServerName    = "app"
	DefaultVersion       = "dev"
)

type StartHttpOpts struct {
	BindInterface         string
	BindPort              int
	DisableLivenessProbe  bool
	DisableReadinessProbe bool
	DisableMetrics        bool
	DisableSwagger        bool
	EventsChannel         chan Event
	LivenessProbes        []Healthcheck
	LogError              Loggerf
	LogInfo               Loggerf
	ReadinessProbes       []Healthcheck
	Router                Router
	ServerName            string
	Version               string
}

func (sho *StartHttpOpts) Init() error {
	if sho.BindInterface == "" {
		sho.BindInterface = DefaultBindInterface
	}
	if sho.BindPort <= 0 || sho.BindPort >= 65535 {
		return fmt.Errorf("failed to receive valid port: received port[%v]", sho.BindPort)
	}
	if sho.ServerName == "" {
		sho.ServerName = DefaultServerName
	}
	if sho.Version == "" {
		sho.Version = DefaultVersion
	}
	if sho.EventsChannel == nil {
		sho.EventsChannel = make(chan Event, 64)
		go defaultEventHandler(sho.EventsChannel)
	}
	if sho.LogError == nil {
		sho.LogError = defaultErrorLogger
	}
	if sho.LogInfo == nil {
		sho.LogInfo = defaultInfoLogger
	}
	return nil
}

func defaultErrorLogger(_ string, _ ...any) {}
func defaultInfoLogger(_ string, _ ...any)  {}
func defaultEventHandler(event chan Event) {
	for {
		<-event
	}
}
