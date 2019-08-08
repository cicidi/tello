# PS4 DualShock controller + Tello
## Use PlayStation 4 DualShock controller control DJI Tello drone

### Preparation 
1. Connect dualshock controller by USB cable. (I tried to use bluetooth, the gods4 lib not able to to make a connection)
2. Power on Tello
3. Connect computer wifi to Tello-XXXXX
### Build
```
go build
```

### Run
```
go run tello
```

### Control
`L1` -> `LANDOFF`  
`R1` -> `TAKEOFF`

`L2` -> `DOWN`  
`R2` -> `UP`

`LFET STICK`

    UP -> FORWARD  
    DOWN -> BACKWARD  
    LEFT -> MOVE_LEFT   
    RIGHT -> MOVE_RIGHT

`RIGHT STICK`

    LEFT -> CLOCKWISE   
    RIGHT -> COUNTER_CLOCKWISE
### Reference
[kpeu3i/gods4](https://github.com/kpeu3i/gods4)  
[gobot](https://gobot.io/)  
[Automating DJI Tello Drone using GOBOT](https://medium.com/tarkalabs/automating-dji-tello-drone-using-gobot-2b711bf42af6)

### TO DO
I will try to use raspberrypi to control Tello later