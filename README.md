# svchostKiller
This tools guards the ports from the black listed services. 

# Why ?
While i am building some jave webservers, it would need some specific ports to be available to go online. But this evil svchost.exe always ends up holding the port for itself. I had to manually taskill the process to get the port back and the worst part is the this svchost again spawn within few second before the server goes live and holds the port causing the server build to fail.😢

As a result, I made a small go script behind the famous `urfave/cli/v2` that will watch the port continously and kills the process if it is of svchost. I have made the code such way that you can watch multiple ports for multiple process like svc, gradle, etc.😎

# How to use
| Flags   | Purpose |
| -------- | ------- |
| --port  | The ports which need to be guarded (comma separated list) |
| --ban | the services from which it needs to be guarded (comma separated list)    |

Build and run the binary as administrator.

# Interesting part
This usecase helped me a lot in learning about the goroutinues. the program spawns new goroutinue for new port to be guarded. And once you want to end the program it will send the stop signal through channels and kills the goroutinues. clean✨.

Hope you had a good read, Thanks!!

