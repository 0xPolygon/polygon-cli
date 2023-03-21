#!/bin/bash

# https://www.frozentux.net/iptables-tutorial/iptables-tutorial.html#TRAVERSINGOFTABLES
# https://inai.de/projects/xtables-addons/
# https://tldp.org/HOWTO/Traffic-Control-HOWTO/

# apt install linux-headers-$(uname -r)
# apt install xtables-addons-dkms xtables-addons-common

# geth --dev --dev.period 2 --ws --ws.addr 0.0.0.0 --ws.port 8546 --http --http.addr 0.0.0.0 --http.port 8545 --http.api admin,debug,web3,eth,txpool,personal,miner,net --verbosity 5 --rpc.gascap 50000000  --rpc.txfeecap 0 --miner.gaslimit  10 --miner.gasprice 1 --gpo.blocks 1 --gpo.percentile 1 --gpo.maxprice 10 --gpo.ignoreprice 2 --dev.gaslimit 50000000
# polycli testharness --listen-ip 0.0.0.0


readonly http_port=8545
readonly ws_port=8546
readonly harness_port=11235

readonly interface=enp1s0


cleanup () {
    tc qdisc del dev $interface root

    iptables -P INPUT ACCEPT
    iptables -P OUTPUT ACCEPT
    iptables -P FORWARD ACCEPT

    iptables -t filter -t raw -F
    iptables -t raw -F
    iptables -t nat -F
    iptables -t mangle -F
}

main () {
    # Allow port 22/ssh access
    iptables -I INPUT -p tcp --dport 22 -j ACCEPT

    # Setup root qdisc
    tc qdisc add dev $interface root handle 1: htb

    # Set default chain policies
    iptables -P INPUT DROP
    iptables -P FORWARD DROP
    iptables -P OUTPUT ACCEPT

    # Accept on localhost
    iptables -A INPUT -i lo -j ACCEPT
    iptables -A OUTPUT -o lo -j ACCEPT

    # Allow established sessions to receive traffic
    iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT

    # Allow ICMP
    iptables -A INPUT -p icmp -j ACCEPT

    # Allow HTTP
    iptables -I INPUT -p tcp --dport $http_port -j ACCEPT
    iptables -I INPUT -p tcp --dport $ws_port -j ACCEPT
    iptables -I INPUT -p tcp --dport $harness_port -j ACCEPT

    # Allow 8000
    iptables -t nat -A PREROUTING -p TCP --dport 8000 -j REDIRECT --to-port $http_port

    # DROP 8001
    iptables -t filter -A INPUT -p TCP --dport 8001 -j DROP

    # REJECT 8002
    iptables -t filter -A INPUT -p TCP --dport 8002 -j REJECT

    # TARPIT 8003
    iptables -t filter -A INPUT -p TCP --dport 8003 -j TARPIT

    # DELUDE 8004
    iptables -t filter -A INPUT -p TCP --dport 8004 -j DELUDE



    # Allow 8101 and Drop 10% of the packets in the way in
    iptables -t nat -A PREROUTING -p TCP --dport 8101 -j REDIRECT --to-port $http_port
    iptables -t raw -I PREROUTING -p TCP --dport 8101 -m statistic --mode random --probability 0.1 -j DROP

    # Allow 8102 and Drop 10% of the packets on the way out
    iptables -t mangle -A PREROUTING -p TCP --dport 8102 -j CONNMARK --set-mark 8102
    iptables -t nat -A PREROUTING -p TCP --dport 8102 -j REDIRECT --to-port $http_port
    iptables -t mangle -A POSTROUTING -p TCP -m connmark --mark 8102 -m statistic --mode random --probability 0.1 -j DROP

    # Allow 8103 and add a 1 second delay
    iptables -t mangle -A PREROUTING -p TCP --dport 8103 -j CONNMARK --set-mark 8103
    iptables -t nat -A PREROUTING -p TCP --dport 8103 -j REDIRECT --to-port $http_port
    iptables -t mangle -A POSTROUTING -p tcp -m connmark --mark 8103 -j CLASSIFY --set-class 1:3
    tc class add dev $interface parent 1: classid 1:3 htb rate 1000mbps
    tc qdisc add dev $interface parent 1:3 handle 3: netem delay 1000ms
    # tc filter add dev $interface parent 1:0 protocol ip u32 match ip sport 8083 FFFF flowid 1:2

    # Allow 8104 and add a packet limit of 2
    iptables -t mangle -A PREROUTING -p TCP --dport 8104 -j CONNMARK --set-mark 8104
    iptables -t nat -A PREROUTING -p TCP --dport 8104 -j REDIRECT --to-port $http_port
    iptables -t mangle -A POSTROUTING -p tcp -m connmark --mark 8104 -j CLASSIFY --set-class 1:4
    tc class add dev $interface parent 1: classid 1:4 htb rate 1000mbps
    tc qdisc add dev $interface parent 1:4 handle 4: netem limit 2

    # Allow 8105 and lose a random 20% of packets
    iptables -t mangle -A PREROUTING -p TCP --dport 8105 -j CONNMARK --set-mark 8105
    iptables -t nat -A PREROUTING -p TCP --dport 8105 -j REDIRECT --to-port $http_port
    iptables -t mangle -A POSTROUTING -p tcp -m connmark --mark 8105 -j CLASSIFY --set-class 1:5
    tc class add dev $interface parent 1: classid 1:5 htb rate 1000mbps
    tc qdisc add dev $interface parent 1:5 handle 5: netem loss random 20%

    # Allow 8106 and corrupt a random 20% of packets
    iptables -t mangle -A PREROUTING -p TCP --dport 8106 -j CONNMARK --set-mark 8106
    iptables -t nat -A PREROUTING -p TCP --dport 8106 -j REDIRECT --to-port $http_port
    iptables -t mangle -A POSTROUTING -p tcp -m connmark --mark 8106 -j CLASSIFY --set-class 1:6
    tc class add dev $interface parent 1: classid 1:6 htb rate 1000mbps
    tc qdisc add dev $interface parent 1:6 handle 6: netem corrupt 20%

    # Allow 8107 and duplicate a random 20% of packets
    iptables -t mangle -A PREROUTING -p TCP --dport 8107 -j CONNMARK --set-mark 8107
    iptables -t nat -A PREROUTING -p TCP --dport 8107 -j REDIRECT --to-port $http_port
    iptables -t mangle -A POSTROUTING -p tcp -m connmark --mark 8107 -j CLASSIFY --set-class 1:7
    tc class add dev $interface parent 1: classid 1:7 htb rate 1000mbps
    tc qdisc add dev $interface parent 1:7 handle 7: netem duplicate 20%

    # Allow 8108 and reorder a random 50% of packets
    iptables -t mangle -A PREROUTING -p TCP --dport 8108 -j CONNMARK --set-mark 8108
    iptables -t nat -A PREROUTING -p TCP --dport 8108 -j REDIRECT --to-port $http_port
    iptables -t mangle -A POSTROUTING -p tcp -m connmark --mark 8108 -j CLASSIFY --set-class 1:8
    tc class add dev $interface parent 1: classid 1:8 htb rate 1000mbps
    tc qdisc add dev $interface parent 1:8 handle 8: netem duplicate 50%

    # Allow 8109 and use a very slow rate limit
    iptables -t mangle -A PREROUTING -p TCP --dport 8109 -j CONNMARK --set-mark 8109
    iptables -t nat -A PREROUTING -p TCP --dport 8109 -j REDIRECT --to-port $http_port
    iptables -t mangle -A POSTROUTING -p tcp -m connmark --mark 8109 -j CLASSIFY --set-class 1:9
    tc class add dev $interface parent 1: classid 1:9 htb rate 1000mbps
    tc qdisc add dev $interface parent 1:9 handle 9: netem rate 56kbit


}

cleanup;
main;


