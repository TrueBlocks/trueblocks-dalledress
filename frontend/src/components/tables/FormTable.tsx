import React, { ReactNode } from "react";
import { Fieldset, Container, Grid, Flex, Text, Stack } from "@mantine/core";
import { Formatter, knownTypes } from "@components";
import classes from "./FormTable.module.css";

export type Progress = {
  have: number;
  want: number;
};

type FieldDefinition<T> = {
  label: string;
  accessor: keyof T;
  type?: knownTypes;
};

type CustomComponentDefinition = {
  component: ReactNode;
};

export type GroupDefinition<T> = {
  title: string;
  fields?: FieldDefinition<T>[]; // Optional to allow custom components without fields
  colSpan?: number;
  components?: CustomComponentDefinition[]; // New field for custom components
};

type FormTableProps<T> = {
  data: Partial<T>;
  definition: GroupDefinition<T>[];
  progress?: Progress;
};

export function FormTable<T>({ data, definition, progress }: FormTableProps<T>) {
  return (
    <Container style={{ maxWidth: "95%", paddingRight: "5%" }}>
      <Grid>
        {definition.map((group, index) => (
          <Grid.Col span={group.colSpan || 12} key={index}>
            <Stack>
              <Fieldset legend={group.title} bg="white" className={classes.fieldSet}>
                {group.fields?.map((field, fieldIndex) => {
                  var value = <>{data[field.accessor]}</>;
                  if (field.type !== undefined) {
                    value = <Formatter type={field.type} value={data[field.accessor]} />;
                  }
                  return (
                    <Flex key={fieldIndex} gap="md" align="center">
                      <Text style={{ minWidth: "150px" }}>{field.label}</Text>
                      <Text>{value}</Text>
                    </Flex>
                  );
                })}
                {group.components?.map((customComponent, componentIndex) => (
                  <div key={componentIndex}>{customComponent.component}</div>
                ))}
              </Fieldset>
            </Stack>
          </Grid.Col>
        ))}
      </Grid>
    </Container>
  );
}
