#!/bin/bash
workspace=$(cd $(dirname $0) && pwd)
cd $workspace

module=example
app=$module
conf=cfg.json
pidfile=var/app.pid
logfile=var/app.log
gitversion=.gitversion

mkdir -p var &>/dev/null


## build & pack
function build() {
    #update_gitversion
    go build -o $app main.go
    sc=$?
    if [ $sc -ne 0 ];then
        echo "build error"
        exit $sc
    else
        echo -n "build ok, vsn=" 
        #version
    fi
}

function pack() {
    build
    version=`./$app -v`
    tar zcvf $app-$version.tar.gz control $app cfg.json $gitversion ./scripts/debug
}

function packbin() {
    build
    version=`./$app -v`
    tar zcvf $app-bin-$version.tar.gz $app gitversion
}

## opt
function start() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo -n "started, pid="
        cat $pidfile
        return 1
    fi

    nohup ./$app >>$logfile 2>&1 &
    echo $! > $pidfile
    echo "start ok, pid=$!"
}

function stop() {
    pid=`cat $pidfile`
    kill $pid
    echo "stoped"
}

function restart() {
    stop && sleep 1 && start
}

function reload() {
    build && stop && sleep 1 && start && sleep 1 && printf "\n"  && tailf
}

## other
function status() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo -n "running, pid="
        cat $pidfile
    else
        echo "stoped"
    fi
}

function version() {
    v=`./$app -v`
    if [ -f $gitversion ];then
        g=`cat $gitversion`
    fi
    echo "$v $g"
}

function tailf() {
    tail -f $logfile
}

## internal
function check_pid() {
    if [ -f $pidfile ];then
        pid=`cat $pidfile`
        if [ -n $pid ]; then
            running=`ps -p $pid|grep -v "PID TTY" |wc -l`
            return $running
        fi
    fi
    return 0
}

function update_gitversion() {
    git log -1 --pretty=%h > $gitversion
}

## usage
function usage() {
    echo "$0 build|pack|packbin|start|stop|restart|reload|status|tail|version"
}

## main
action=$1
case $action in
    ## build
    "build" )
        build
        ;;
    "pack" )
        pack
        ;;
    "packbin" )
        packbin
        ;;
    ## opt
    "start" )
        start
        ;;
    "stop" )
        stop
        ;;
    "restart" )
        restart
        ;;
    "reload" )
        reload
        ;;
    ## other
    "status" )
        status
        ;;
    "version" )
        version
        ;;
    "tail" )
        tailf
        ;;
    * )
        usage
        ;;
esac
