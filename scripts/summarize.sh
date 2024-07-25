#!/usr/bin/env bash

# Define the attributes
attributes=(
    "adverb"
    "adjective"
    "noun"
    "emotion"
    "occupation"
    "action"
    "artStyle1"
    "artStyle2"
    "litStyle"
    "color1"
    "color2"
    "color3"
    "orientation"
    "gaze"
    "backStyle"
    "setting"
)

# Define the paths
data_path="output/data"
selectors_path="output/selectors"
summary_path="output/summary"
histograms_path="output/histograms"

# Function to perform the summarization work
summarize() {
    local attr=$1
    echo "Processing ${attr}..."
    grep ${attr} ${data_path}/* | cut -f2 -d, >${selectors_path}/${attr}.txt
    sort -n ${selectors_path}/${attr}.txt | uniq -c | sort -n | sed 's/^ *//' | tr ' ' '\t' >${summary_path}/${attr}.txt
    cut -f1 ${summary_path}/${attr}.txt | sort -n | uniq -c | sed 's/^ *//' | tr ' ' '\t' >${histograms_path}/${attr}.txt
}

# Loop through each attribute and call the summarize function
for attr in "${attributes[@]}"; do
    summarize ${attr}
done
