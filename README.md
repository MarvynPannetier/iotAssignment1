# iotAssignment1

I put the three go files i worked on 

if you want to test clone this three repo :

https://github.com/aslakjohansen/sdu-iot-mqtt-siggen.git

https://github.com/aslakjohansen/sdu-iot-mqtt-mavg.git

https://github.com/aslakjohansen/sdu-iot-mqtt-func.git


And replace the go files by the three of this repo. After that just run make produce in each of the three folders and <mosquitto -v -t "#"> in a fourth terminal

use this link for the test harness : 

https://github.com/aslakjohansen/simple-java-test-harness.git

link to the report : https://docs.google.com/document/d/1UhqOgxAPsGIccz7RKWE9ALRkPDM8nr2v3YACpNuMKZU/edit?usp=sharing

I started it and i will finish it before wednesday

The Makefile is for the "client".

For all TestHarnesses line 20 in the Logger.java file should have System.nanoTime().
The delay is lower than 1ms, and would therefore not show up if using System.currentTimeMillis().
