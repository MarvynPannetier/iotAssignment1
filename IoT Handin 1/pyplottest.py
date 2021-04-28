from matplotlib import pyplot as plt
import pandas as pd
import seaborn as sns
import numpy as np

dfMavg = pd.read_csv('testMavg.csv', names=['timestamp', 'status', 'line'])
dfFunc = pd.read_csv('testFunc.csv', names=['timestamp', 'status', 'line', "value"])

dfSiggen = pd.read_csv('testSiggen.csv', names=['timestamp', 'status', 'line'])
dfClient = pd.read_csv('testClient.csv', names=['timestamp', 'time', 'value'])

def mavgDiffs(df):
    newDf = pd.DataFrame(columns=['line', 'diff'])
    for index, row in df.iterrows():
        if 'siggen/room1/rhum' in row['line'] and row['status'] == 'Received':
            for index2, row2 in df.loc[index:].iterrows():
                if 'mavg/room1/rhum' in row2['line'] and row2['status'] == "Sending":
                    #append row to new df with the diff
                    diff = (int(row2['timestamp']) - int(row['timestamp']))/10000
                    line = row2['line']
                    new_row = {'line' : line, 'diff' : diff}
                    newDf = newDf.append(new_row, ignore_index=True)
                    break
        elif 'siggen/room1/temp' in row['line'] and row['status'] == 'Received':
            for index2, row2 in df.loc[index:].iterrows():
                if 'mavg/room1/temp' in row2['line'] and row2['status'] == "Sending":
                    #append row to new df with the diff
                    diff = (int(row2['timestamp']) - int(row['timestamp']))/10000
                    line = row2['line']
                    new_row = {'line' : line, 'diff' : diff}
                    newDf = newDf.append(new_row, ignore_index=True)
                    break
        if 'siggen/room2/rhum' in row['line'] and row['status'] == 'Received':
            for index2, row2 in df.loc[index:].iterrows():
                if 'mavg/room2/rhum' in row2['line'] and row2['status'] == "Sending":
                    #append row to new df with the diff
                    diff = (int(row2['timestamp']) - int(row['timestamp']))/10000
                    line = row2['line']
                    new_row = {'line' : line, 'diff' : diff}
                    newDf = newDf.append(new_row, ignore_index=True)
                    break
        elif 'siggen/room2/temp' in row['line'] and row['status'] == 'Received':
            for index2, row2 in df.loc[index:].iterrows():
                if 'mavg/room2/temp' in row2['line'] and row2['status'] == "Sending":
                    #append row to new df with the diff
                    diff = (int(row2['timestamp']) - int(row['timestamp']))/10000 
                    line = row2['line']
                    new_row = {'line' : line, 'diff' : diff}
                    newDf = newDf.append(new_row, ignore_index=True)
                    break
    #print(newDf)
    newDf= newDf[newDf['diff'] > 0]
    return newDf

def funcDiffs(df):
    newDf = pd.DataFrame(columns=['line', 'diff'])
    for index, row in df.iterrows():
        nextRow = index + 1
        if 'mavg/room1/temp' in row['line'] and row['status'] == 'Received' and df.at[nextRow, 'status'] != 'Sending':
            for index2, row2 in df.loc[index:].iterrows():
                if 'func/room1/ahum' in row2['line'] and row2['status'] == "Sending" :
                    #append row to new df with the diff
                    diff = (int(row2['timestamp']) - int(row['timestamp']))/10000
                    line = row2['line']
                    new_row = {'line' : line, 'diff' : diff}
                    newDf = newDf.append(new_row, ignore_index=True)
                    break
        if 'mavg/room1/rhum' in row['line'] and row['status'] == 'Received' and df.at[nextRow, 'status'] != 'Sending':
            for index2, row2 in df.loc[index:].iterrows():
                if 'func/room1/ahum' in row2['line'] and row2['status'] == "Sending" :
                    #append row to new df with the diff
                    diff = (int(row2['timestamp']) - int(row['timestamp']))/10000
                    line = row2['line']
                    new_row = {'line' : line, 'diff' : diff}
                    newDf = newDf.append(new_row, ignore_index=True)
                    break
        if 'mavg/room2/temp' in row['line'] and row['status'] == 'Received' and df.at[nextRow, 'status'] != 'Sending':
            for index2, row2 in df.loc[index:].iterrows():
                if 'func/room2/ahum' in row2['line'] and row2['status'] == "Sending" :
                    #append row to new df with the diff
                    diff = (int(row2['timestamp']) - int(row['timestamp']))/10000
                    line = row2['line']
                    new_row = {'line' : line, 'diff' : diff}
                    newDf = newDf.append(new_row, ignore_index=True)
                    break
        if 'mavg/room2/rhum' in row['line'] and row['status'] == 'Received' and df.at[nextRow, 'status'] != 'Sending':
            for index2, row2 in df.loc[index:].iterrows():
                if 'func/room2/ahum' in row2['line'] and row2['status'] == "Sending" :
                    #append row to new df with the diff
                    diff = (int(row2['timestamp']) - int(row['timestamp']))/10000
                    line = row2['line']
                    new_row = {'line' : line, 'diff' : diff}
                    newDf = newDf.append(new_row, ignore_index=True)
                    break
    newDf= newDf[newDf['diff'] < 3000]#Removes everything above 3000
    newDf= newDf[newDf['diff'] > 0]
    #print(newDf)
    return newDf

    
def siggenDiffs(siggenDf, mavgDf): 
    newDf = pd.DataFrame(columns=['line', 'diff'])

    mavgReceive = mavgDf[~mavgDf.status.str.contains("Send")]

    mavgReceiveRoom1 = mavgReceive[~mavgReceive.line.str.contains("room2")]
    mavgReceiveRoom2 = mavgReceive[~mavgReceive.line.str.contains("room1")]

    mavgReceiveRoom1Rhum = mavgReceiveRoom1[~mavgReceiveRoom1.line.str.contains("temp")].reset_index(drop=True)
    mavgReceiveRoom2Rhum = mavgReceiveRoom2[~mavgReceiveRoom2.line.str.contains("temp")].reset_index(drop=True)

    mavgReceiveRoom1Temp = mavgReceiveRoom1[~mavgReceiveRoom1.line.str.contains("rhum")].reset_index(drop=True)
    mavgReceiveRoom2Temp = mavgReceiveRoom2[~mavgReceiveRoom2.line.str.contains("rhum")].reset_index(drop=True)


    siggenSendRoom1 = siggenDf[~siggenDf.line.str.contains("room2")]
    siggenSendRoom2 = siggenDf[~siggenDf.line.str.contains("room1")]

    siggenSendRoom1Rhum = siggenSendRoom1[~siggenSendRoom1.line.str.contains("temp")].reset_index(drop=True)
    siggenSendRoom2Rhum = siggenSendRoom2[~siggenSendRoom2.line.str.contains("temp")].reset_index(drop=True)

    siggenSendRoom1Temp = siggenSendRoom1[~siggenSendRoom1.line.str.contains("rhum")].reset_index(drop=True)
    siggenSendRoom2Temp = siggenSendRoom2[~siggenSendRoom2.line.str.contains("rhum")].reset_index(drop=True)
    

    for index, row in mavgReceiveRoom1Rhum.iterrows():
        diff = (int(row['timestamp']) - int(siggenSendRoom1Rhum.ix[index]['timestamp']))/10000
        line = row['line']
        new_row = {'line' : line, 'diff' : diff}
        newDf = newDf.append(new_row, ignore_index=True)

    for index, row in mavgReceiveRoom2Rhum.iterrows():
        diff = (int(row['timestamp']) - int(siggenSendRoom2Rhum.ix[index]['timestamp']))/10000
        line = row['line']
        new_row = {'line' : line, 'diff' : diff}
        newDf = newDf.append(new_row, ignore_index=True)

    for index, row in mavgReceiveRoom1Temp.iterrows():
        diff = (int(row['timestamp']) - int(siggenSendRoom1Temp.ix[index]['timestamp']))/10000
        line = row['line']
        new_row = {'line' : line, 'diff' : diff}
        newDf = newDf.append(new_row, ignore_index=True)

    for index, row in mavgReceiveRoom2Temp.iterrows():
        diff = (int(row['timestamp']) - int(siggenSendRoom2Temp.ix[index]['timestamp']))/10000
        line = row['line']
        new_row = {'line' : line, 'diff' : diff}
        newDf = newDf.append(new_row, ignore_index=True)

    newDf= newDf[newDf['diff'] < 300]
    newDf= newDf[newDf['diff'] > 0]
    #print(newDf)
    return newDf


def mavgToFuncDiffs(mavgDf, funcDf):
    newDf = pd.DataFrame(columns=['line', 'diff'])
    mavgSend = mavgDf[~mavgDf.status.str.contains("Receive")]

    funcReceive = funcDf[~funcDf.status.str.contains("Send")]

    mavgRoom1 = mavgSend[~mavgSend.line.str.contains("room2")]
    mavgRoom2 = mavgSend[~mavgSend.line.str.contains("room1")]

    mavgRoom1Temp = mavgRoom1[~mavgRoom1.line.str.contains("rhum")].reset_index(drop=True)
    mavgRoom2Temp = mavgRoom2[~mavgRoom2.line.str.contains("rhum")].reset_index(drop=True)

    mavgRoom1Rhum = mavgRoom1[~mavgRoom1.line.str.contains("temp")].reset_index(drop=True)
    mavgRoom2Rhum = mavgRoom2[~mavgRoom2.line.str.contains("temp")].reset_index(drop=True)

    
    funcRoom1 = funcReceive[~funcReceive.line.str.contains("room2")]
    funcRoom2 = funcReceive[~funcReceive.line.str.contains("room1")]

    funcRoom1Temp = funcRoom1[~funcRoom1.line.str.contains("rhum")].reset_index(drop=True)
    funcRoom2Temp = funcRoom2[~funcRoom2.line.str.contains("rhum")].reset_index(drop=True)

    funcRoom1Rhum = funcRoom1[~funcRoom1.line.str.contains("temp")].reset_index(drop=True)
    funcRoom2Rhum = funcRoom2[~funcRoom2.line.str.contains("temp")].reset_index(drop=True)

    for index, row in funcRoom1Temp.iterrows():
        diff = (int(row['timestamp']) - int(mavgRoom1Temp.ix[index]['timestamp']))/10000
        line = row['line']
        new_row = {'line' : line, 'diff' : diff}
        newDf = newDf.append(new_row, ignore_index=True)
    for index, row in funcRoom2Temp.iterrows():
        diff = (int(row['timestamp']) - int(mavgRoom2Temp.ix[index]['timestamp']))/10000
        line = row['line']
        new_row = {'line' : line, 'diff' : diff}
        newDf = newDf.append(new_row, ignore_index=True)
    for index, row in funcRoom1Rhum.iterrows():
        diff = (int(row['timestamp']) - int(mavgRoom1Rhum.ix[index]['timestamp']))/10000
        line = row['line']
        new_row = {'line' : line, 'diff' : diff}
        newDf = newDf.append(new_row, ignore_index=True)
    for index, row in funcRoom2Rhum.iterrows():
        diff = (int(row['timestamp']) - int(mavgRoom2Rhum.ix[index]['timestamp']))/10000
        line = row['line']
        new_row = {'line' : line, 'diff' : diff}
        newDf = newDf.append(new_row, ignore_index=True)
    newDf= newDf[newDf['diff'] < 400]
    #REMOVES NEGATIVE NUMBERS THIS IS 28 ENTRIES FOR THIS EXACT ONE
    newDf= newDf[newDf['diff'] > 0]
    return newDf


def funcToClientDiffs(funcDf, clientDf):
    newDf = pd.DataFrame(columns=['line', 'diff'])    
    funcSendDf = funcDf[~funcDf.status.str.contains("Receive")].reset_index(drop=True)

    for index, row in clientDf.iterrows():
        diff = (int(row['timestamp']) - int(funcSendDf.ix[index]['timestamp']))/10000
        line = row['value']
        
        #Two prints that prove the negative numbers are not because wrongful ordering
        #print(line, funcSendDf.ix[index]['value'])
        #print(diff)
        new_row = {'line' : line, 'diff' : diff}
        newDf = newDf.append(new_row, ignore_index=True)
    newDf= newDf[newDf['diff'] < 200]
    newDf= newDf[newDf['diff'] > 0]
    #newDf= newDf[newDf['diff'] > 0]
    return newDf

def siggenToClientDiffs(siggenDf, clientDf, funcDf):
    clientReceiveDf = pd.DataFrame(columns=['line', 'client_receive'])

    funcSendDf = funcDf[~funcDf.status.str.contains("Receive")].reset_index(drop=True)

    for index, row in funcSendDf.iterrows():
        line = row['line']
        client_receive = clientDf.ix[index]['timestamp']
        new_row = {'line' : line, 'client_receive' : client_receive}
        clientReceiveDf = clientReceiveDf.append(new_row, ignore_index=True)


    clientReceiveRoom1 = clientReceiveDf[~clientReceiveDf.line.str.contains("room2")].reset_index(drop=True)
    clientReceiveRoom2 = clientReceiveDf[~clientReceiveDf.line.str.contains("room1")].reset_index(drop=True)

    siggenSendRoom1 = siggenDf[~siggenDf.line.str.contains("room2")].iloc[::2].reset_index(drop=True)
    siggenSendRoom2 = siggenDf[~siggenDf.line.str.contains("room1")].iloc[::2].reset_index(drop=True)


    timeDiffDf = pd.DataFrame(columns=['line', 'diff'])
    for index, row in clientReceiveRoom1.iterrows():
        diff = (int(row['client_receive']) - int(siggenSendRoom1.ix[index]['timestamp']))/10000
        line = row['line']
        
        #Two prints that prove the negative numbers are not because wrongful ordering
        #print(line, funcSendDf.ix[index]['value'])
        #print(diff)
        new_row = {'line' : line, 'diff' : diff}
        timeDiffDf = timeDiffDf.append(new_row, ignore_index=True)

    for index, row in clientReceiveRoom2.iterrows():
        diff = (int(row['client_receive']) - int(siggenSendRoom2.ix[index]['timestamp']))/10000
        line = row['line']
        
        #Two prints that prove the negative numbers are not because wrongful ordering
        #print(line, funcSendDf.ix[index]['value'])
        #print(diff)
        new_row = {'line' : line, 'diff' : diff}
        timeDiffDf = timeDiffDf.append(new_row, ignore_index=True)

    timeDiffDf = timeDiffDf[timeDiffDf['diff'] > 0]
    timeDiffDf= timeDiffDf[timeDiffDf['diff'] < 700]
    return timeDiffDf


    
def drawSigToMavg(): 
    fromSigToMavg = siggenDiffs(dfSiggen, dfMavg)

    mean = fromSigToMavg["diff"].mean()
    print("Mean:" ,mean)

    counted = fromSigToMavg['diff'].value_counts().sort_index().to_frame().rename_axis('delay').reset_index('delay').rename(columns={"diff": "occurrences"})

    counted['bucket'] = pd.cut(counted.delay, 10)
    bucketed = counted[["bucket", "occurrences"]].groupby("bucket").sum()
    #print(bucketed)

    plot = bucketed.plot(kind="bar")
    plot.set_xlabel("\nTime in 1/100th of a millisecond", fontsize=14)
    plt.title("Time between siggen and mavg", fontsize=20)
    plt.xticks(rotation=60)
    plt.show()        

def drawMavgRecToSend():
    plotMavgDf = mavgDiffs(dfMavg)

    mean = plotMavgDf["diff"].mean()
    print("Mean:" ,mean)

    funcBox = plotMavgDf['diff'].value_counts().sort_index().plot.bar(rot=0)
    funcBox.set_title("Mavg computation time", fontsize=20)
    funcBox.set_xlabel("\nTime in 1/100th of a millisecond", fontsize=14)
    funcBox.set_ylabel("Occurrences")
    #print(plotMavgDf['diff'].value_counts().sort_index())
    plt.xticks(rotation=60)
    plt.show()     

def drawMavgToFunc():
    fromMavgToFunc = mavgToFuncDiffs(dfMavg, dfFunc)

    mean = fromMavgToFunc["diff"].mean()
    print("Mean:" ,mean)


    counted = fromMavgToFunc['diff'].value_counts().sort_index().to_frame().rename_axis('delay').reset_index('delay').rename(columns={"diff": "occurrences"})

    counted['bucket'] = pd.cut(counted.delay, 15)
    bucketed = counted[["bucket", "occurrences"]].groupby("bucket").sum()
    #print(bucketed)

    plot = bucketed.plot(kind="bar")
    plot.set_xlabel("\nTime in 1/100th of a millisecond", fontsize=14)
    plt.title("Time between mavg and func", fontsize=20)
    plt.xticks(rotation=60)
    plt.show()  

def drawFuncRecToSend():
    plotFuncDf = funcDiffs(dfFunc)

    mean = plotFuncDf["diff"].mean()
    print("Mean:" ,mean)

    funcBox = plotFuncDf['diff'].value_counts().sort_index().plot.bar(rot=0)
    funcBox.set_title("Func computation time", fontsize=20)
    funcBox.set_xlabel("Time in 1/100th of a millisecond", fontsize=14)
    funcBox.set_ylabel("Occurrences", fontsize=14)
    #print(plotFuncDf['diff'].value_counts().sort_index())
    plt.xticks(rotation=60)
    plt.show()

def drawFuncToClient():
    plotFuncToClient = funcToClientDiffs(dfFunc, dfClient)

    mean = plotFuncToClient["diff"].mean()
    print("Mean:" ,mean)

    counted = plotFuncToClient['diff'].value_counts().sort_index().to_frame().rename_axis('delay').reset_index('delay').rename(columns={"diff": "occurrences"})

    counted['bucket'] = pd.cut(counted.delay, 15)
    bucketed = counted[["bucket", "occurrences"]].groupby("bucket").sum()
    #print(bucketed)

    plot = bucketed.plot(kind="bar")
    plot.set_xlabel("\nTime in 1/100th of a millisecond", fontsize=14)
    plt.title("Time between func and client", fontsize=20)
    plt.xticks(rotation=60)
    plt.show()  

def drawSiggenToClient():
    plotFuncToClient = siggenToClientDiffs(dfSiggen, dfClient, dfFunc)
    counted = plotFuncToClient['diff'].value_counts().sort_index().to_frame().rename_axis('delay').reset_index('delay').rename(columns={"diff": "occurrences"})

    mean = plotFuncToClient["diff"].mean()
    print("Mean:" ,mean)
    

    counted['bucket'] = pd.cut(counted.delay, 15)
    bucketed = counted[["bucket", "occurrences"]].groupby("bucket").sum()
    #print(bucketed)

    plot = bucketed.plot(kind="bar")
    plot.set_xlabel("Time in 1/100th of a millisecond", fontsize=14)
    plt.title("Time between siggen and client", fontsize=20)
    plt.xticks(rotation=45)
    plt.show()  



#drawSigToMavg()
#drawMavgRecToSend()
drawMavgToFunc()


#drawFuncRecToSend()

#drawFuncToClient()

#drawSiggenToClient()

