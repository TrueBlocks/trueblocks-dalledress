import React, { ReactNode } from "react";

export function View(params: { title?: string; children: ReactNode }) {
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        height: "100%",
        flex: 1,
      }}
    >
      {params.children}
    </div>
  );
}
