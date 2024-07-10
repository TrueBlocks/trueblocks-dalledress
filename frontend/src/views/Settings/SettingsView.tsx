import React from "react";
import { Checkbox, InputLabel } from "@mantine/core";
import { View, ViewHeader, ViewContent, ViewStatus } from "@/components/view";

function SettingsView() {
  return (
    <View>
      <ViewHeader>Settings View Header</ViewHeader>
      <ViewContent>
        <InputLabel>
          <Checkbox label={"A checkbox"} />
        </InputLabel>
      </ViewContent>
      <ViewStatus>Status / Progress</ViewStatus>
    </View>
  );
}

export default SettingsView;
