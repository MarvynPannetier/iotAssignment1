package main

import (
    "fmt"
    "io/ioutil"
    "encoding/json"
    "time"
    "log"
    "os"
    
    "github.com/eclipse/paho.mqtt.golang"
)

type Sample struct {
    Time  float64 `json:"time"`
    Value float64 `json:"value"`
}

type Signal struct {
    Topic   string   `json:"topic"`
    Samples []Sample `json:"samples"`
}
type Config []Signal

const (
    config_filename string = "config.json"
)

var (
    globalTimeLogger *log.Logger // create a log var
    brokers []string = []string{"tcp://127.0.0.1:1883"}
    t0 float64 = 0.0
    t1 float64 = 0.0
)

// This function will create the log file and define a formatting of each beginning of line
//***************************************************************************************
func init() {
    file, err := os.OpenFile("logs_siggen.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

    globalTimeLogger = log.New(file, "globalTime: ", log.Ldate|log.Ltime|log.Lshortfile)
    
}
//**************************************************************************************


func mqtt_connect () mqtt.Client {
    // configure options
    options := mqtt.NewClientOptions()
    for _, broker := range brokers {
      options.AddBroker(broker)
    }
    
    // start mqtt client
    client := mqtt.NewClient(options)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
    
    return client
}

func read_config (filename string) Config {
    var config Config
    
    // read config file
    data, err := ioutil.ReadFile(config_filename)
    if err!=nil {
        panic("Unable to load config file: "+err.Error())
    }
    
    // initial partial parsing of config file
    err = json.Unmarshal(data, &config)
    if err!=nil {
        panic("Unable to unmarshal config file: "+err.Error())
    }
    
    return config
}

func get_time () float64 {
    return float64(time.Now().UnixNano())
}

func produce (client mqtt.Client, signal Signal, t0 float64) {

	var t2 float64

   for i := range(signal.Samples) {
        signal.Samples[i].Time *= 1000000000
    }
    
    period := signal.Samples[len(signal.Samples)-1].Time
    var i float64 = 0.0
    fmt.Println("About to produce", signal.Topic)

    for {
    	 
        for _, sample := range(signal.Samples) {
     
            tnext := t0+i*period+sample.Time
            
            // produce payload
            var new_sample Sample = Sample{tnext, sample.Value}
            message, _ := json.Marshal(new_sample)
            
            // sleep
            t := get_time()
            tdiff := tnext-t
            time.Sleep(time.Duration(tdiff) * time.Nanosecond)
            
            // publish
            client.Publish(signal.Topic, 1, false, message)
            t2 = get_time()
            fmt.Println(signal.Topic, ":" , tnext, sample.Value )
            
            // with this line we will print a line on the log file with the information between parentheses in addition to those define in the init function.
            globalTimeLogger.Println(signal.Topic, "Time : ",(t2-t0)/1000000000,"s Value : ", sample.Value ) 
          
            
        }
        i += 1.0
    }
}

func main () {
   
    config := read_config(config_filename)
    client := mqtt_connect()
    
    
    // start up a producer for each topic
    t0 := get_time()
    for _, signal := range(config) {
        go produce(client, signal, t0)
    }
    
    select{} // block forever
}
