import React, { ReactNode, useState, useEffect } from "react";
import classes from "./ViewStatus.module.css";
import { EventsOn, EventsOff } from "@runtime";
import { Text } from "@mantine/core";
import { MessageType } from "@gocode/app/App";

type Progress = {
  address: string;
  have: number;
  want: number;
};

export function ViewStatus() {
  const [statusMessage, setStatusMessage] = useState<string>("");
  const [color, setColor] = useState<string>(classes.green);

  useEffect(() => {
    const handleDone = () => {
      setStatusMessage("done");
      setColor(classes.green);
    };

    const handleProgress = (p: Progress) => {
      setStatusMessage(`Progress (${p.address}): ${p.have}/${p.want}`);
      setColor(classes.green);
    };

    const handleError = (errorStr: string) => {
      setStatusMessage(`Error: ${errorStr}`);
      setColor(classes.red);
    };

    EventsOn("Completed", handleDone);
    EventsOn("Progress", handleProgress);
    EventsOn("Error", handleError);

    return () => {
      EventsOff("Done");
      EventsOff("Progress");
      EventsOff("Error");
    };
  }, []);

  return (
    <Text size="lg">
      <div className={color}>{statusMessage}</div>
    </Text>
  );
}
