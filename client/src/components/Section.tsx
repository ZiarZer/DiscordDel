import styled from "styled-components";

import { ActionInputBar } from "./ActionInputBar";
import { StatusMessage } from "./StatusMessage";
import { InfoList } from "./InfoList";

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

export function Section({
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
  getAvatarUrl,
}) {
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
        currentObject={currentObject}
        fields={infoFields}
        getAvatarUrl={getAvatarUrl}
      />
    </Wrapper>
  );
}
