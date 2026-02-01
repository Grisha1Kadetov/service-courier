package enviroment

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/pkg/log"

	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"
)

const portEnv = "PORT"
const addrEnv = "ADDR"

const pgHostEnv = "POSTGRES_HOST"
const pgPortEnv = "POSTGRES_PORT"
const pgUserEnv = "POSTGRES_USER"
const pgPasswordEnv = "POSTGRES_PASSWORD"
const pgDBEnv = "POSTGRES_DB"

const tickDeliveryEnv = "TICK_DELIVERY"
const tickOrderEnv = "TICK_ORDER"
const orderServiceHostEnv = "ORDER_SERVICE_HOST"
const retryAttemptsEnv = "RETRY_ATTEMPTS"
const retrySleepEnv = "RETRY_SLEEP"

const kafkaBrokersEnv = "KAFKA_BROKERS"
const kafkaTopicsEnv = "KAFKA_TOPICS"
const kafkaGroupIDEnv = "KAFKA_GROUP_ID"

const RpsEnv = "RPS"
const BurstEnv = "BURST"

const InDockerEnv = "IN_DOCKER"
const PprofUsernameEnv = "PPROF_USER"
const PprofPasswordEnv = "PPROF_PASSWORD"

type Environment struct {
	l                log.Logger
	Addr             string
	PprofAddr        string
	PprofUsername    string
	PprofPassword    string
	Port             int
	PostgresEnv      PostgresEnv
	KafkaEnv         KafkaEnv
	OrderServiceHost string
	TickDelivery     time.Duration
	TickOrder        time.Duration
	RetryAttempts    int
	RetrySleep       time.Duration
	Rps              float64
	Burst            float64
}

type PostgresEnv struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type KafkaEnv struct {
	Brokers []string
	GroupID string
	Topics  []string
}

func NewEnvironment(l log.Logger) *Environment {
	return &Environment{l: l}
}

func (env *Environment) Load() {
	err := godotenv.Load()
	if err != nil {
		env.l.Warn("error loading .env file, using environment variables")
	}

	env.getPostgresEnv()
	env.getKafkaEnv()

	env.PprofAddr = "127.0.0.1:6060"
	if os.Getenv(InDockerEnv) == "1" {
		env.PprofAddr = "0.0.0.0:6060"
	}

	env.PprofUsername = env.getRequeredEnv(PprofUsernameEnv)
	env.PprofPassword = env.getRequeredEnv(PprofPasswordEnv)

	var ok bool
	env.Port, ok = getFlagPort()
	if !ok {
		env.Port, err = strconv.Atoi(env.getRequeredEnv(portEnv))
		if err != nil {
			env.l.Fatal("PORT should be int", log.NewField("err", err))
		}
	}
	env.Addr = os.Getenv(addrEnv)
	env.Rps, err = strconv.ParseFloat(env.getRequeredEnv(RpsEnv), 64)
	if err != nil {
		env.l.Fatal("RPS should be float", log.NewField("err", err))
	}
	env.Burst, err = strconv.ParseFloat(env.getRequeredEnv(BurstEnv), 64)
	if err != nil {
		env.l.Fatal("BURST should be float", log.NewField("err", err))
	}
	env.TickDelivery, err = time.ParseDuration(env.getRequeredEnv(tickDeliveryEnv))
	if err != nil {
		env.l.Fatal("TickDelivery should be duration", log.NewField("err", err))
	}
	env.TickOrder, err = time.ParseDuration(env.getRequeredEnv(tickOrderEnv))
	if err != nil {
		env.l.Fatal("TickOrder should be duration", log.NewField("err", err))
	}
	env.RetrySleep, err = time.ParseDuration(env.getRequeredEnv(retrySleepEnv))
	if err != nil {
		env.l.Fatal("RetrySleep should be duration", log.NewField("err", err))
	}
	env.RetryAttempts, err = strconv.Atoi(env.getRequeredEnv(retryAttemptsEnv))
	if err != nil {
		env.l.Fatal("RetryAttempts should be int", log.NewField("err", err))
	}
	env.OrderServiceHost = env.getRequeredEnv(orderServiceHostEnv)
}

func (e *Environment) getPostgresEnv() {
	var env PostgresEnv
	var err error
	env.Host = e.getRequeredEnv(pgHostEnv)
	env.Port, err = strconv.Atoi(e.getRequeredEnv(pgPortEnv))
	if err != nil {
		e.l.Fatal("Postgres port should be int", log.NewField("err", err))
	}
	env.Password = e.getRequeredEnv(pgPasswordEnv)
	env.User = e.getRequeredEnv(pgUserEnv)
	env.DBName = e.getRequeredEnv(pgDBEnv)
	e.PostgresEnv = env
}

func (e *Environment) getKafkaEnv() {
	var env KafkaEnv

	env.GroupID = e.getRequeredEnv(kafkaGroupIDEnv)

	env.Brokers = strings.Split(e.getRequeredEnv(kafkaBrokersEnv), ",")
	if len(env.Brokers) == 0 {
		e.l.Fatal("Kafka brokers should contain at least one broker")
	}

	env.Topics = strings.Split(e.getRequeredEnv(kafkaTopicsEnv), ",")
	if len(env.Topics) == 0 {
		e.l.Fatal("Kafka topics should contain at least one topic")
	}

	e.KafkaEnv = env
}

func (e *Environment) GetDSN() string {
	p := e.PostgresEnv
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v", p.User, p.Password, p.Host, p.Port, p.DBName)
}

func getFlagPort() (int, bool) {
	portFlaged := flag.Int("port", -1, "server port")
	flag.Parse()
	ok := *portFlaged != -1
	return *portFlaged, ok
}

func (e *Environment) getRequeredEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		e.l.Fatal("this key is required", log.NewField("key", key))
	}
	return v
}
