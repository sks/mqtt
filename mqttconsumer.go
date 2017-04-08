package mqtt

import (
	"fmt"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/trivago/gollum/core"
)

type MqttConsumer struct {
	core.SimpleConsumer
	connectionString string
	topic            string
	client           mqtt.Client
	clientId         string
}

func init() {
	core.TypeRegistry.Register(MqttConsumer{})
}

func (cons *MqttConsumer) messageHandler(_ mqtt.Client, msg mqtt.Message) {
	cons.Enqueue(msg.Payload())
}

func (cons *MqttConsumer) connectionHandler(client mqtt.Client) {
	cons.Log.Debug.Printf("Connected to MQTT Server: %s. Listening to topic %s", cons.connectionString, cons.topic)
	token := client.Subscribe(cons.topic, 0, cons.messageHandler)
	token.Wait()
	err := token.Error()
	if err != nil {
		cons.Log.Error.Printf("Error Subsribing to topic %q: %s", cons.topic, err)
		time.AfterFunc(3*time.Second, cons.startConnection)
		return
	}
}

func (cons *MqttConsumer) startConnection() {
	cons.Log.Debug.Printf("Connecting to MQTT Server: %s", cons.connectionString)
	clientOptions := mqtt.NewClientOptions()
	clientOptions.ClientID = cons.clientId
	clientOptions.AddBroker(cons.connectionString)
	clientOptions.SetOnConnectHandler(cons.connectionHandler)
	clientOptions.SetConnectionLostHandler(cons.connectionLostHandler)

	cons.client = mqtt.NewClient(clientOptions)
}

func (cons *MqttConsumer) connectionLostHandler(client mqtt.Client, err error) {
	cons.Log.Error.Printf("Disconnected from MQTT Server: %s. Will try to Auto_reconnect", err)
}

func (cons *MqttConsumer) startListening() {
	if cons.client == nil {
		cons.startConnection()
	}
	cons.client.Connect()
}

func (cons *MqttConsumer) Configure(conf core.PluginConfigReader) error {
	cons.SimpleConsumer.Configure(conf)
	cons.clientId = fmt.Sprintf("%s_%s", "gollum_mqtt_", conf.GetID())
	cons.connectionString = conf.GetString("connectionString", "tcp://localhost:1883")
	cons.topic = conf.GetString("topic", "#")
	quietPeriod := uint(conf.GetInt("quietPeriod", 200))

	if os.Getenv("MQTT_ENABLE_LOGGING") == "true" {
		mqtt.DEBUG = cons.Log.Debug
		mqtt.WARN = cons.Log.Warning
		mqtt.CRITICAL = cons.Log.Error
		mqtt.ERROR = cons.Log.Error
	}
	cons.SetStopCallback(func() {
		if cons.client != nil {
			cons.client.Disconnect(quietPeriod)
		}
	})
	return conf.Errors.OrNil()
}

func (cons *MqttConsumer) Consume(workers *sync.WaitGroup) {
	go cons.startListening()
	cons.ControlLoop()
}
