package mqttClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"rederinghub.io/utils/config"
	"rederinghub.io/utils/tracer"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

const (
	CmdInfo string = "info"
	CmdHeight string = "height"
	CmdSensor string = "sensor"
)


type MqttMessage struct {
	Payload map[string]string
	Type string
	Status string
	Topic   string
}

type deviceMqtt struct {
	Config config.MQTTConfig 
	Mqtt mqtt.Client
	Token mqtt.Token
	Tracer tracer.ITracer
	CMessage         chan *MqttMessage
	CErr chan error
}

type IDeviceMqtt interface {
	Subscribe(channelName string) 
	Publish(channelName string, data interface{})  error
	Disconnect() 
	MessageHandler(client mqtt.Client, msg mqtt.Message) 
	GetMessageChan()  (chan *MqttMessage, chan error)
}


func NewDeviceMqtt(conf config.MQTTConfig, tracer tracer.ITracer) (IDeviceMqtt, error) {
	c := &deviceMqtt{}
	c.Config = conf
	c.CMessage = make(chan *MqttMessage)
	c.CErr = make(chan error)
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%s", conf.Address, conf.Port)).
		SetClientID(uuid.New().String()).
		SetUsername( conf.UserName).
		SetPassword( conf.Password).
		SetAutoReconnect(true).
		SetKeepAlive(60 * 5 * time.Second).
		SetPingTimeout(time.Second * 3).
		SetMaxReconnectInterval(time.Second * 30).
		SetCleanSession(true).
		SetDefaultPublishHandler(c.MessageHandler).
		SetConnectionLostHandler(c.connectLostHandler)
		
	c.Mqtt = mqtt.NewClient(opts)
	c.Token = c.Mqtt.Connect()
	c.Tracer = tracer

	wait :=  c.Token.Wait()
	err := 	 c.Token.Error()
	if  wait && err  != nil {
		return nil, err
	}

	return c, nil
}	

func (m deviceMqtt) Subscribe(channelName string)  {
	m.Mqtt.Subscribe(channelName, 1, m.MessageHandler)
	fmt.Printf("Subscribed channel - %s", channelName)
}

func (m deviceMqtt) GetMessageChan()   (chan *MqttMessage, chan error) {
	return m.CMessage, m.CErr
}

func (m deviceMqtt)  MessageHandler(client mqtt.Client, msg mqtt.Message) {
	msgObject := &MqttMessage{}
	var err error

	defer func ()  {
		if err != nil {
			m.CErr <- err 
		}
		m.CMessage <- msgObject
	}()

	payload := make(map[string]string)
	err = json.Unmarshal(msg.Payload(), &payload)
	if  err != nil {
		return
	}

	msgType, ok := payload["type"]
	if !ok {
		err = errors.New("Cannot detect message type")
		return
	}
	
	msgStatus, ok := payload["status"]
	if !ok {
		msgStatus = "0"
	}

	msgObject.Topic =  msg.Topic()
	msgObject.Payload =  payload
	msgObject.Type =  msgType
	msgObject.Status =  msgStatus
}

func (m deviceMqtt) Publish(topicName string,  msg interface{}) error  {
	encodedMsg, ok := msg.([]byte)
	if !ok {
		encodedMsgByte, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		encodedMsg = encodedMsgByte
	}
	
	token := m.Mqtt.Publish(topicName, 0, false, encodedMsg)
	token.Wait()
	if token.Error() != nil {
		e := token.Error()
		return e
	}
	return nil
}

func (m deviceMqtt) connectLostHandler(client mqtt.Client, err error) {
	//Logger.log.Errorf("Connect mqtt server lost: %v\n", err)
}

func (m deviceMqtt) reConnectHandler(client mqtt.Client, clientOption *mqtt.ClientOptions) {
	
}

func (m deviceMqtt) Disconnect()  {
	m.Mqtt.Disconnect(100)
}

