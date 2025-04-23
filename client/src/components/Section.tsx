import styled from "styled-components";

import { ActionInputBar } from "./ActionInputBar";
import { StatusMessage } from "./StatusMessage";
import { InfoList } from "./InfoList";
import { Button } from ".";
import { Channel, Guild, InfoListFieldConfig, User } from "../types";
import { ChangeEvent } from "react";

const Wrapper = styled.div`
  background-color: #ffffff30;
  border-radius: 1em;
  padding: 1em;
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 1em;
  width: 100%;
`;

const SectionTitle = styled.h3`
  margin: 0;
`;

const ActionsContainer = styled.div`
  display: flex;
  justify-content: space-evenly;
`;

type SectionProps<T extends User | Guild | Channel> = {
  title: string;
  actionInputBar: {
    inputPlaceholder: string;
    buttonLabel: string;
    enabled?: boolean;
    secret?: boolean;
    onSubmit: () => void;
    onChange: (e: ChangeEvent) => void;
  };
  statusMessage: string;
  currentObject: T | null;
  infoFields: Array<InfoListFieldConfig<T>>;
  getAvatarUrl?: (param: T) => string | undefined;
  actions?: Array<{ label: string; onClick: () => void }>;
};

export function Section<T extends User | Guild | Channel>({
  title,
  actionInputBar: {
    inputPlaceholder,
    buttonLabel,
    enabled = false,
    secret = false,
    onSubmit,
    onChange,
  },
  statusMessage,
  currentObject,
  infoFields,
  getAvatarUrl = () => undefined,
  actions = [],
}: SectionProps<T>) {
  return (
    <Wrapper>
      <SectionTitle>{title}</SectionTitle>
      <ActionInputBar
        inputPlaceholder={inputPlaceholder}
        buttonText={buttonLabel}
        enabled={enabled}
        secret={secret}
        onSubmit={onSubmit}
        onChange={onChange}
      />
      <StatusMessage message={statusMessage} success={currentObject != null} />
      <InfoList
        currentObject={
          currentObject as Record<keyof T, string | number | undefined>
        }
        fields={infoFields}
        getAvatarUrl={getAvatarUrl}
      />
      <ActionsContainer>
        {actions.map(({ label, onClick }) => (
          <Button disabled={currentObject == null} onClick={onClick}>{label}</Button>
        ))}
      </ActionsContainer>
    </Wrapper>
  );
}
