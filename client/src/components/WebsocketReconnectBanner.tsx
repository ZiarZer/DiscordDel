import styled from "styled-components";
import { Button } from ".";
import { ReadyState } from "react-use-websocket";

const Wrapper = styled.div<{ $backgroundColor: string; $color: string }>`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1em;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  padding: 0.25em;
  color: white;
  background-color: ${({ $backgroundColor }: { $backgroundColor: string }) =>
    $backgroundColor};
`;

export function WebsocketReconnectBanner({
  readyState,
  retry,
}: {
  readyState: ReadyState;
  retry: () => void;
}) {
  if ([ReadyState.UNINSTANTIATED, ReadyState.CLOSED].includes(readyState)) {
    return (
      <Wrapper $backgroundColor="#aa0000">
        Websocket disconnected
        <Button onClick={retry}>Try reconnecting</Button>
      </Wrapper>
    );
  } else if (readyState === ReadyState.CONNECTING) {
    return (
      <Wrapper $backgroundColor="#ff8000">
        Websocket connecting
        <Button disabled>Try reconnecting</Button>
      </Wrapper>
    );
  }
  return null;
}
