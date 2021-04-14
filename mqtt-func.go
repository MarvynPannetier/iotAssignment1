package main

import (
    "fmt"
    "encoding/json"
    "sync"
    "strings"
    "math"
    "log"
    "os"
    "time"
    "github.com/eclipse/paho.mqtt.golang"
)

type Sample struct {
    Time  float64 `json:"time"`
    Value float64 `json:"value"`
}

var (
    brokers []string = []string{"tcp://127.0.0.1:1883"}
    pattern   string = "mavg/+/+"
    dispatch map[string]chan Sample = make(map[string]chan Sample)
    dispatch_mux sync.Mutex
    client mqtt.Client
    t0 float64 = 0.0
    receivedTempTimeLogger *log.Logger  //marvyn
    receivedHumTimeLogger *log.Logger  //marvyn
    sentTimeLogger *log.Logger  //marvyn
    differenceTimeLogger *log.Logger  //marvyn
)



//***************************************************************************************
func init() {
    file, err := os.OpenFile("logs_func.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }
    receivedTempTimeLogger = log.New(file,"received : ", log.Ldate|log.Ltime|log.Lshortfile)
    receivedHumTimeLogger = log.New(file,"received : ", log.Ldate|log.Ltime|log.Lshortfile)
    sentTimeLogger = log.New(file,"sent : ", log.Ldate|log.Ltime|log.Lshortfile)
    differenceTimeLogger = log.New(file,"difference : ",log.Ldate|log.Ltime|log.Lshortfile)
}
//****************************Marvyn*****************************************************

func get_time () float64 {
    return float64(time.Now().UnixNano())
}


func publish (topic string, channel chan Sample) {
    for sample := range channel {
        message, _ := json.Marshal(sample)
        t1 := get_time()
        client.Publish(topic, 1, false, message)
        sentTimeLogger.Println(topic,"Time : ",t1,"Value : ",sample.Value)
        //fmt.Println("sent : ",topic,sample.Value) //add by marvyn
        t2:=t1-t0
        differenceTimeLogger.Println(topic,"Time : ",t2)
    }
}

// unit: g/m³
func calc_abs_hum (temp float64, rhum float64) float64 {
	
    return 6.112*math.Pow(math.E, (17.67*temp)/(temp+243.5))*rhum*2.1674/(273.15+temp)
}

func ahum (channel_temp chan Sample,
           channel_rhum chan Sample,
           channel_ahum chan Sample,
           topic string) { //add by marvyn
           t0=get_time()
    for {
        temp_sample := <- channel_temp
        receivedTempTimeLogger.Println(topic, "Time : ",t0,"value : ",temp_sample.Value)
        rhum_sample := <- channel_rhum
        receivedHumTimeLogger.Println(topic, "Time : ",t0,"value : ",rhum_sample.Value)
        
        temp := temp_sample.Value
        rhum := rhum_sample.Value
        ahum := calc_abs_hum(temp, rhum)
      //  fmt.Println("received : ",topic,"temperature : ",temp,"humidity : ",rhum) //add by marvyn
        
        channel_ahum <- Sample{temp_sample.Time, ahum}
    }
}

func dispatch_sample (client mqtt.Client, message mqtt.Message) {
    var topic string = message.Topic()
    var sample Sample
    var channel_temp chan Sample
    var channel_rhum chan Sample
    var channel_ahum chan Sample
    
    // preprocess topic
    tparts := strings.Split(topic, "/")
    modality := tparts[len(tparts)-1]
    
    // unmarshal
    err := json.Unmarshal(message.Payload(), &sample)
    if err!=nil {
        fmt.Println("Unable to unmarshal incoming sample:", err)
        return
    }
    
    // make sure that channel exists
    dispatch_mux.Lock()
    defer dispatch_mux.Unlock()
    channel, ok := dispatch[topic]
    if !ok {
        // create necessary channels
        channel_temp = make(chan Sample, 2)
        channel_rhum = make(chan Sample, 2)
        channel_ahum = make(chan Sample, 2)
        
        // choose action
        switch modality {
        case "temp":
            channel = channel_temp
        case "rhum":
            channel = channel_rhum
        default:
            return
        }
         fmt.Println("received : ",topic) //add by marvyn
        // define topic names
        topic_temp := "mavg/"+strings.Join(tparts[1:len(tparts)-1], "/")+"/temp"
        topic_rhum := "mavg/"+strings.Join(tparts[1:len(tparts)-1], "/")+"/rhum"
        topic_ahum := "func/"  +strings.Join(tparts[1:len(tparts)-1], "/")+"/ahum"
        
        // start up consumers
        
      //  room:=topic[:len(topic)-5] //add by marvyn
        go ahum(channel_temp, channel_rhum, channel_ahum, topic)//add by marvyn
        go publish(topic_ahum, channel_ahum)
        
        // register channels
        dispatch[topic_temp] = channel_temp
        dispatch[topic_rhum] = channel_rhum
        
    }
    
    // queue channel
    channel <- sample
}

func mqtt_subscribe () {
    // configure options
    options := mqtt.NewClientOptions()
    for _, broker := range brokers {
      options.AddBroker(broker)
    }
    
    // start mqtt client
    client = mqtt.NewClient(options)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
    
    // set up subscription
    if token := client.Subscribe(pattern, 2, dispatch_sample); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
}

func main () {
    mqtt_subscribe()
    
    select{} // block forever
}
