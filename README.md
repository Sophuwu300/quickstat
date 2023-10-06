A simple program that will show the current RAM and CPU usage as well as frequency and temperature of the CPU. Written in Go, using /proc and /sys to get the information. It takes 1 second to show the information as it needs to average the CPU usage over 1 second. Without this it would show the CPU usage as 100%. It should work on modern Linux systems using both Intel and AMD CPUs.

Linux shows the temperature information in different ways depending on the CPU. For Intel i series CPUs it shows the temperature under /sys/class/tharmal/thermal_zoneX/temp, the X is a number that can be different depending on the system. So I have to check the type of thermal zone before getting the temp. For AMD Ryzen CPUs it shows the temperature it shows the temperature under /sys/class/hwmon/hwmonX/temp1_input, the X is a number that can be different depending on the system. So I have to check the temp1_label to see if it is a CPU temp.

Running the program results in an output similar to this:

```   Mem: 23 % USED   39 % BUFF   38 % FREE   CPU: 10 %   34 C   2665 MHz```

Compiled binaries with [gpg signitures](https://sophuwu.site/pgp.txt) are available at [sophuwu.site](https://sophuwu.site/quickstat/).
