# kfk-offsets

Command line tool for LAG monitoring consumer groups listing and offset reset.

```shell
❯ ./kfkgo
Using config file: ~/kfk-offsets/kfk-offset.yaml
Usage:
  kfk-offsets [flags]
  kfk-offsets [command]

Available Commands:
  completion     Generate the autocompletion script for the specified shell
  consumer-group Create / List consumer groups.
  help           Help about any command
  lag            Checks topic lag per consumer group.
  offset         A brief description of your command

Flags:
      --config string   config file (default is $HOME/kfk-offset.yaml)
  -h, --help            help for kfk-offsets

Use "kfk-offsets [command] --help" for more information about a command.
```

Listing consumer groups:

```shell
❯ ./kfkgo cg
Using config file: ~/kfk-offsets/kfk-offset.yaml
propertest-cg1
propertest-cg1a1122
propertest-cg
propertest-cg111
propertest-cg2
propertest-cg1a11
propertest-cg1a
propertest-cg1a11a
propertest-cg1a1
monitoring Lag
```

```shell
❯ ./kfkgo lag
Using config file: ~/kfk-offsets/kfk-offset.yaml
GROUP                                    TOPIC                                    PARTITION       CURRENT-OFFSET  LOG-END-OFFSET  LAG            
propertest-cg1a1122                      propertest                               0               179             180             1              
propertest-cg1a1122                      propertest                               1               162             165             3              
propertest-cg1a1122                      propertest                               2               190             192             2              
propertest-cg1a1122                      propertest                               3               174             177             3              
propertest-cg1a1122                      propertest                               4               187             189             2              
propertest-cg1a1122                      propertest                               5               167             168             1              
propertest-cg1a1122                      propertest                               6               177             177             0              
propertest-cg1a1122                      propertest                               7               160             163             3              
propertest-cg1a1122                      propertest                               8               192             193             1              
propertest-cg1a1122                      propertest                               9               185             186             1              
propertest-cg1a1122                      propertest                               10              183             184             1              
propertest-cg1a1122                      propertest                               11              184             184             0              
TOTAL LAG:                                                                                                                        18             

GROUP                                    TOPIC                                    PARTITION       CURRENT-OFFSET  LOG-END-OFFSET  LAG            
propertest-cg1                           propertest3                              0               112             113             1              
propertest-cg1                           propertest3                              1               79              79              0              
propertest-cg1                           propertest3                              2               103             103             0              
propertest-cg1                           propertest3                              3               86              86              0              
propertest-cg1                           propertest3                              4               101             101             0              
propertest-cg1                           propertest3                              5               105             106             1              
propertest-cg1                           propertest3                              6               97              97              0              
propertest-cg1                           propertest3                              7               89              90              1              
propertest-cg1                           propertest3                              8               93              93              0              
TOTAL LAG:                                                                                                                        3              

GROUP                                    TOPIC                                    PARTITION       CURRENT-OFFSET  LOG-END-OFFSET  LAG            
propertest-cg2                           propertest2                              0               16              18              2              
propertest-cg2                           propertest2                              1               18              18              0              
propertest-cg2                           propertest2                              2               24              26              2              
propertest-cg2                           propertest2                              3               24              24              0              
propertest-cg2                           propertest2                              4               15              16              1              
propertest-cg2                           propertest2                              5               14              16              2              
propertest-cg2                           propertest2                              6               25              27              2              
propertest-cg2                           propertest2                              7               22              23              1              
propertest-cg2                           propertest2                              8               25              27              2              
propertest-cg2                           propertest2                              9               10              10              0              
propertest-cg2                           propertest2                              10              6               6               0              
propertest-cg2                           propertest2                              11              14              14              0              
TOTAL LAG:                                                                                                                        12             

GROUP                                    TOPIC                                    PARTITION       CURRENT-OFFSET  LOG-END-OFFSET  LAG            
propertest-cg111                         propertest1                              0               87              87              0              
propertest-cg111                         propertest1                              1               89              90              1              
propertest-cg111                         propertest1                              2               86              87              1              
propertest-cg111                         propertest1                              3               82              83              1              
propertest-cg111                         propertest1                              4               100             101             1              
propertest-cg111                         propertest1                              5               104             104             0              
propertest-cg111                         propertest1                              6               86              86              0              
propertest-cg111                         propertest1                              7               79              81              2              
propertest-cg111                         propertest1                              8               98              99              1              
TOTAL LAG:            
```

Moving offset between Consumer Groups as example showed with default Kafka commands and with kfkgo

Comparing offset 

```shell
Consumer group 'propertest-cg1a11' has no active members.
propertest-cg1a11 propertest      0          179             180             1               -               -               -
propertest-cg1a11 propertest      1          162             165             3               -               -               -
propertest-cg1a11 propertest      2          190             192             2               -               -               -
propertest-cg1a11 propertest      3          174             177             3               -               -               -
propertest-cg1a11 propertest      4          187             192             5               -               -               -
propertest-cg1a11 propertest      5          167             169             2               -               -               -
propertest-cg1a11 propertest      6          177             179             2               -               -               -
propertest-cg1a11 propertest      7          160             163             3               -               -               -
propertest-cg1a11 propertest      8          192             195             3               -               -               -
propertest-cg1a11 propertest      9          185             188             3               -               -               -
propertest-cg1a11 propertest      10         183             184             1               -               -               -
propertest-cg1a11 propertest      11         184             184             0               -               -               -


Consumer group 'propertest-cg' has no active members.
propertest-cg   propertest      0          179             180             1               -               -               -
propertest-cg   propertest      1          162             165             3               -               -               -
propertest-cg   propertest      2          190             192             2               -               -               -
propertest-cg   propertest      3          174             177             3               -               -               -
propertest-cg   propertest      4          187             192             5               -               -               -
propertest-cg   propertest      5          167             169             2               -               -               -
propertest-cg   propertest      6          177             179             2               -               -               -
propertest-cg   propertest      7          160             163             3               -               -               -
propertest-cg   propertest      8          192             195             3               -               -               -
propertest-cg   propertest      9          185             188             3               -               -               -
propertest-cg   propertest      10         183             184             1               -               -               -
propertest-cg   propertest      11         184             184             0               -               -               -
```

We can see that is exactly same on each partition and in total same check with above command.

```shell
❯ ./kfkgo lag
Using config file: ~/kfk-offsets/kfk-offset.yaml
GROUP                                    TOPIC                                    PARTITION       CURRENT-OFFSET  LOG-END-OFFSET  LAG
propertest-cg1a11                        propertest                               0               179             180             1
propertest-cg1a11                        propertest                               1               162             165             3
propertest-cg1a11                        propertest                               2               190             192             2
propertest-cg1a11                        propertest                               3               174             177             3
propertest-cg1a11                        propertest                               4               187             192             5
propertest-cg1a11                        propertest                               5               167             169             2
propertest-cg1a11                        propertest                               6               177             179             2
propertest-cg1a11                        propertest                               7               160             163             3
propertest-cg1a11                        propertest                               8               192             195             3
propertest-cg1a11                        propertest                               9               185             188             3
propertest-cg1a11                        propertest                               10              183             184             1
propertest-cg1a11                        propertest                               11              184             184             0
TOTAL LAG:                                                                                                                        28


❯ ./kfkgo lag -g propertest-cg -t propertest
Using config file: ~/kfk-offsets/kfk-offset.yaml
GROUP                                    TOPIC                                    PARTITION       CURRENT-OFFSET  LOG-END-OFFSET  LAG
propertest-cg                            propertest                               0               179             180             1
propertest-cg                            propertest                               1               162             165             3
propertest-cg                            propertest                               2               190             192             2
propertest-cg                            propertest                               3               174             177             3
propertest-cg                            propertest                               4               187             192             5
propertest-cg                            propertest                               5               167             169             2
propertest-cg                            propertest                               6               177             179             2
propertest-cg                            propertest                               7               160             163             3
propertest-cg                            propertest                               8               192             195             3
propertest-cg                            propertest                               9               185             188             3
propertest-cg                            propertest                               10              183             184             1
propertest-cg                            propertest                               11              184             184             0
TOTAL LAG:                                                                                                                        28
```

We have same offsets … 
Now to simulate the process we will produce messages to topic and keep consuming on one of the consumer group propertest-cg1a11 we will produce 5 messages and consume them all on that consumer group which will give us info that we consumed 
`^CProcessed a total of 33 messages 28 + 5`

```shell
❯ ./kfkgo lag
Using config file: ~/kfk-offsets/kfk-offset.yaml
GROUP                                    TOPIC                                    PARTITION       CURRENT-OFFSET  LOG-END-OFFSET  LAG
propertest-cg1a11                        propertest                               0               183             183             0
propertest-cg1a11                        propertest                               1               165             165             0
propertest-cg1a11                        propertest                               2               192             192             0
propertest-cg1a11                        propertest                               3               177             177             0
propertest-cg1a11                        propertest                               4               192             192             0
propertest-cg1a11                        propertest                               5               169             169             0
propertest-cg1a11                        propertest                               6               180             180             0
propertest-cg1a11                        propertest                               7               164             164             0
propertest-cg1a11                        propertest                               8               195             195             0
propertest-cg1a11                        propertest                               9               188             188             0
propertest-cg1a11                        propertest                               10              184             184             0
propertest-cg1a11                        propertest                               11              184             184             0
TOTAL LAG:                                                                                                                        0
```

```shell
❯ ./kfkgo lag -g propertest-cg -t propertest
Using config file: ~/kfk-offsets/kfk-offset.yaml
GROUP                                    TOPIC                                    PARTITION       CURRENT-OFFSET  LOG-END-OFFSET  LAG
propertest-cg                            propertest                               0               179             183             4
propertest-cg                            propertest                               1               162             165             3
propertest-cg                            propertest                               2               190             192             2
propertest-cg                            propertest                               3               174             177             3
propertest-cg                            propertest                               4               187             192             5
propertest-cg                            propertest                               5               167             169             2
propertest-cg                            propertest                               6               177             180             3
propertest-cg                            propertest                               7               160             164             4
propertest-cg                            propertest                               8               192             195             3
propertest-cg                            propertest                               9               185             188             3
propertest-cg                            propertest                               10              183             184             1
propertest-cg                            propertest                               11              184             184             0
TOTAL LAG:                                                                                                                        33
```

Now we will move the offset from propertest-cg back again to propertest-cg1a11 which will allow us to process same messages on that CG.

```shell
❯ ./kfkgo offset -m
Using config file: ~/kfk-offsets/kfk-offset.yaml
moving
```

and verification again :

Kafka commands: `kafka-consumer-groups.sh`
```shell 
Consumer group 'propertest-cg1a11' has no active members.
propertest-cg1a11 propertest      0          179             183             4               -               -               -
propertest-cg1a11 propertest      1          162             165             3               -               -               -
propertest-cg1a11 propertest      2          190             192             2               -               -               -
propertest-cg1a11 propertest      3          174             177             3               -               -               -
propertest-cg1a11 propertest      4          187             192             5               -               -               -
propertest-cg1a11 propertest      5          167             169             2               -               -               -
propertest-cg1a11 propertest      6          177             180             3               -               -               -
propertest-cg1a11 propertest      7          160             164             4               -               -               -
propertest-cg1a11 propertest      8          192             195             3               -               -               -
propertest-cg1a11 propertest      9          185             188             3               -               -               -
propertest-cg1a11 propertest      10         183             184             1               -               -               -
propertest-cg1a11 propertest      11         184             184             0               -               -               -


Consumer group 'propertest-cg' has no active members.
propertest-cg   propertest      0          179             183             4               -               -               -
propertest-cg   propertest      1          162             165             3               -               -               -
propertest-cg   propertest      2          190             192             2               -               -               -
propertest-cg   propertest      3          174             177             3               -               -               -
propertest-cg   propertest      4          187             192             5               -               -               -
propertest-cg   propertest      5          167             169             2               -               -               -
propertest-cg   propertest      6          177             180             3               -               -               -
propertest-cg   propertest      7          160             164             4               -               -               -
propertest-cg   propertest      8          192             195             3               -               -               -
propertest-cg   propertest      9          185             188             3               -               -               -
propertest-cg   propertest      10         183             184             1               -               -               -
propertest-cg   propertest      11         184             184             0               -               -               -
```

```shell
❯ ./kfkgo lag -g propertest-cg -t propertest
Using config file: ~/kfk-offsets/kfk-offset.yaml
GROUP                                    TOPIC                                    PARTITION       CURRENT-OFFSET  LOG-END-OFFSET  LAG
propertest-cg                            propertest                               0               179             183             4
propertest-cg                            propertest                               1               162             165             3
propertest-cg                            propertest                               2               190             192             2
propertest-cg                            propertest                               3               174             177             3
propertest-cg                            propertest                               4               187             192             5
propertest-cg                            propertest                               5               167             169             2
propertest-cg                            propertest                               6               177             180             3
propertest-cg                            propertest                               7               160             164             4
propertest-cg                            propertest                               8               192             195             3
propertest-cg                            propertest                               9               185             188             3
propertest-cg                            propertest                               10              183             184             1
propertest-cg                            propertest                               11              184             184             0
TOTAL LAG:                                                                                                                        33
```

```shell
❯ ./kfkgo lag
Using config file: ~/kfk-offsets/kfk-offset.yaml
GROUP                                    TOPIC                                    PARTITION       CURRENT-OFFSET  LOG-END-OFFSET  LAG
propertest-cg1a11                        propertest                               0               179             183             4
propertest-cg1a11                        propertest                               1               162             165             3
propertest-cg1a11                        propertest                               2               190             192             2
propertest-cg1a11                        propertest                               3               174             177             3
propertest-cg1a11                        propertest                               4               187             192             5
propertest-cg1a11                        propertest                               5               167             169             2
propertest-cg1a11                        propertest                               6               177             180             3
propertest-cg1a11                        propertest                               7               160             164             4
propertest-cg1a11                        propertest                               8               192             195             3
propertest-cg1a11                        propertest                               9               185             188             3
propertest-cg1a11                        propertest                               10              183             184             1
propertest-cg1a11                        propertest                               11              184             184             0
TOTAL LAG:                                                                                                                        33
```
Concept works i did say :grinning: 
