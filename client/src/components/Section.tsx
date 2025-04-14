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
        inputPlaceholder="Authorization token"
        buttonText="Authenticate"
        enabled={true}
        secret
        onSubmit={onSubmit}
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
