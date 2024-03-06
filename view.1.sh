#!/usr/bin/env bash

rm -f viewing/*
#cp -p annotated/$1*human-with.png viewing/
#cp -p annotated/$1*human-like.png viewing/
#cp -p annotated/$1*happiness.png viewing/
#cp -p annotated/$1*postal.png viewing/
cp -p annotated/$1*fury.png viewing/
#cp -p annotated/$1*love.png viewing/
#cp -p annotated/$1*human.png viewing/
#cp -p annotated/$1*human2.png viewing/
#cp -p annotated/$1*steam.png viewing/
#cp -p annotated/$1*solar.png viewing/
#cp -p annotated/$1*00000.png viewing/
#cp -p annotated/$1*00001.png viewing/
cp -p annotated/*00002.png viewing/
open viewing
read -p "Press Enter to continue"
