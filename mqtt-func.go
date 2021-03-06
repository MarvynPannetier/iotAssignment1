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
    
    // create log variables 
    receivedTempTimeLogger *log.Logger  
    receivedHumTimeLogger *log.Logger  
    sentTimeLogger *log.Logger  
    differenceTimeLogger *log.Logger  
    )

// This function will create the log file and define a formatting of each beginning of line depending on the var log used
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
//***************************************************************************************

func get_time () float64 {
    return float64(time.Now().UnixNano())
}


func publish (topic string, channel chan Sample) {
    for sample := range channel {
        message, _ := json.Marshal(sample)
        t1 := get_time()
        client.Publish(topic, 1, false, message)
        // with this line we will print a line on the log file with the information between parentheses in addition to those define in the init function.
        sentTimeLogger.Println(topic,"Time : ",t1/1000000000,"s Value : ",sample.Value)
        //fmt.Println("sent : ",topic,sample.Value) //add by marvyn
        t2:=t1-t0
        if t2 > 0.0 {
        // with this line we will print a line on the log file with the information between parentheses in addition to those define in the init function.
        differenceTimeLogger.Println(topic,"Time : ",t2/1000,"us")
        }
    }
}

// unit: g/m??
func calc_abs_hum (temp float64, rhum float64) float64 {
	
    return 6.112*math.Pow(math.E, (17.67*temp)/(temp+243.5))*rhum*2.1674/(273.15+temp)
}

func ahum (channel_temp chan Sample,
           channel_rhum chan Sample,
           channel_ahum chan Sample,
           topic_temp string,
           topic_rhum string) { //add by marvyn
         
    for {
      t0=get_time()
        temp_sample := <- channel_temp
        // with this line we will print a line on the log file with the information between parentheses in addition to those define in the init function.
        receivedTempTimeLogger.Println(topic_temp, "Time : ",t0/1000000000,"s value : ",temp_sample.Value)
        rhum_sample := <- channel_rhum
        // with this line we will print a line on the log file with the information between parentheses in addition to those define in the init function.
        receivedHumTimeLogger.Println(topic_rhum, "Time : ",t0/1000000000,"s value : ",rhum_sample.Value)
        
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
     fmt.Println("received : ",topic) 
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
        
        // define topic names
        topic_temp := "mavg/"+strings.Join(tparts[1:len(tparts)-1], "/")+"/temp"
        topic_rhum := "mavg/"+strings.Join(tparts[1:len(tparts)-1], "/")+"/rhum"
        topic_ahum := "func/"  +strings.Join(tparts[1:len(tparts)-1], "/")+"/ahum"
        
        // start up consumers
         
        //room:=topic[:len(topic)-5] //add by marvyn
        go ahum(channel_temp, channel_rhum, channel_ahum, topic_temp, topic_rhum)
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
