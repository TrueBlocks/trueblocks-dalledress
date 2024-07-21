import React, { ReactNode, useState, useEffect } from "react";
import classes from "./ViewStatus.module.css";
import { EventsOn, EventsOff } from "@runtime";
import { Text } from "@mantine/core";

export function ViewStatus() {
  const [statusMessage, setStatusMessage] = useState<string>("");
  const [color, setColor] = useState<string>(classes.green);

  useEffect(() => {
    const handleDone = () => {
      setStatusMessage("done");
      setColor(classes.green);
    };

    const handleProgress = (x: number, y: number) => {
      setStatusMessage(`Progress: ${x}/${y}`);
      setColor(classes.green);
    };

    const handleError = (errorStr: string) => {
      setStatusMessage(`Error: ${errorStr}`);
      setColor(classes.red);
    };

    EventsOn("Done", handleDone);
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
