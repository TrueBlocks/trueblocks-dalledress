import React from "react";
import { View, ViewHeader, ViewContent, ViewStatus } from "@/components/view";

function HomeView() {
  return (
    <View>
      <ViewHeader>Home View Title</ViewHeader>
      <ViewContent>Home View Content</ViewContent>
      <ViewStatus>Status / Progress</ViewStatus>
    </View>
  );
}

export default HomeView;
