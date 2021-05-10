# Remote Media

A simple yet fun little web app to allow control of your media from a distance.
This project was made because I use music streaming services on my pc, and have to move a lot in my appartment (due to the pandemic).

The working is simple: 
The service hosts a webpage on your machine, available to all other machines on your wifi. You can then visit `<IP-of-the-machine>:1337` and control your music from a distance.

## Installation
1. [install golang](https://golang.org/doc/install)
1. go to scripts `cd scripts`
1. allow execution of the `service-linux.sh` `chmod u+x service-linux.sh`
1. run the script `sudo service-linux.sh $USER`
1. make sure it is properly running by visiting [the url of the webapp](http://localhost:1337)