# iotAssignment1

if you want to test clone this three repo :

https://github.com/aslakjohansen/sdu-iot-mqtt-siggen.git

https://github.com/aslakjohansen/sdu-iot-mqtt-mavg.git

https://github.com/aslakjohansen/sdu-iot-mqtt-func.git


And replace the go files by the three of this repo. After that just run make produce in each of the three folders and <mosquitto -v -t "#"> in a fourth terminal

use this link for the test harness : 

https://github.com/aslakjohansen/simple-java-test-harness.git

The Makefile is for the "client".

For all TestHarnesses line 20 in the Logger.java file should have System.nanoTime().
The delay is lower than 1ms, and would therefore not show up if using System.currentTimeMillis().
