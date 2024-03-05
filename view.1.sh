#!/usr/bin/env bash

rm -f viewing/*
#cp -p annotated/$1*human-with.png viewing/
#cp -p annotated/$1*human-like.png viewing/
#cp -p annotated/$1*happiness.png viewing/
#cp -p annotated/$1*postal.png viewing/
#cp -p annotated/$1*fury.png viewing/
#cp -p annotated/$1*love.png viewing/
cp -p annotated/$1*human.png viewing/
cp -p annotated/$1*human2.png viewing/
open viewing
read -p "Press Enter to continue"
