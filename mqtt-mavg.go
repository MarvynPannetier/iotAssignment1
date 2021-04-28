package main

import (
    "fmt"
    "encoding/json"
    "sync"
    "time"
    "github.com/eclipse/paho.mqtt.golang"
    "log"
    "os"
)

type Sample struct {
    Time  float64 `json:"time"`
    Value float64 `json:"value"`
}

var (
    brokers []string = []string{"tcp://127.0.0.1:1883"}
    pattern   string = "siggen/+/+"
    dispatch map[string]chan Sample = make(map[string]chan Sample)
    dispatch_mux sync.Mutex
    t0 float64 = 0.0
    client mqtt.Client
    
    // create log variables 
    receivedTimeLogger *log.Logger  
    sentTimeLogger *log.Logger  
    differenceTimeLogger *log.Logger 
    ) 

const (
    WINDOW_SIZE float64 = 3.0
)

// This function will create the log file and define a formatting of each beginning of line depending on the var log used
//***************************************************************************************
func init() {
    file, err := os.OpenFile("logs_mavg.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }
    receivedTimeLogger = log.New(file,"received : ", log.Ldate|log.Ltime|log.Lshortfile)
    sentTimeLogger = log.New(file,"sent : ", log.Ldate|log.Ltime|log.Lshortfile)
    differenceTimeLogger = log.New(file,"difference : ",log.Ldate|log.Ltime|log.Lshortfile)
}
//***************************************************************************************


func get_time () float64 {
    return float64(time.Now().UnixNano())
}


func mavg (topic string, channel chan Sample) {
    otopic := "mavg"+topic[6:]
    fmt.Println("Republishing moving average to", otopic)
    var sent bool
    var received bool
    
    // initialize window and sum
    window := make([]float64, int(WINDOW_SIZE))
    for i := range window {
        window[i] = 0
    }
    sum := 0.0
    
    // service loop
    i := 0.0
    for sample := range channel {
        value := sample.Value
        received = false
        sent = false 
        
        //Marvyn : I chose to put it here because it is the moment we "received" the sample
   	t0=get_time()
   	 // with this line we will print a line on the log file with the information between parentheses in addition to those define in the init function.
   	receivedTimeLogger.Println(otopic, "Time : ",t0/1000000000, "s value : ",value) 
        received = true
        // update window and sum
        sum += value - window[int(i)%int(WINDOW_SIZE)]
        window[int(i)%int(WINDOW_SIZE)] = value
        
        // build message
        var new_sample Sample = Sample{sample.Time, sum/WINDOW_SIZE}
        message, _ := json.Marshal(new_sample)
              
        t2 := get_time()
         // with this line we will print a line on the log file with the information between parentheses in addition to those define in the init function.
        sentTimeLogger.Println(otopic, "Time : ",t2/1000000000, "s value : ", value)
  	sent = true
        // t1 is the time between the moment we received the sample and the moment we send the average
        t1 := (t2-t0)
        if received && sent && t1 > 0.0 {
         // with this line we will print a line on the log file with the information between parentheses in addition to those define in the init function.
        differenceTimeLogger.Println(otopic, "Time : ",t1/1000,"us")
        }   
        // publish
        client.Publish(otopic, 1, false, message)
        i++
    }
}

func dispatch_sample (client mqtt.Client, message mqtt.Message) {
    var topic string = message.Topic()
    var sample Sample
    
    // unmarshal
    err := json.Unmarshal(message.Payload(), &sample)
    if err!=nil {
        fmt.Println("Unable to unmarshal incoming sample:", err)
        return
    }
    
    // make sure that channel exists
    dispatch_mux.Lock()
    channel, ok := dispatch[topic]
    if !ok {
        channel = make(chan Sample, 2)
        go mavg(topic, channel);
        dispatch[topic] = channel
    }
    dispatch_mux.Unlock()
    
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
