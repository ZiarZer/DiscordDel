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
  background-color: ${({ $backgroundColor }: { $backgroundColor: string }) =>
    $backgroundColor};
  color: ${({ $color }: { $color: string }) => $color};
`;

export function WebsocketReconnectBanner({
  readyState,
  retry,
}: {
  readyState: ReadyState;
  retry: () => void;
}) {
  const isDisconnected = [
    ReadyState.UNINSTANTIATED,
    ReadyState.CLOSED,
  ].includes(readyState);
  return (
    isDisconnected && (
      <Wrapper $backgroundColor='#aa0000' $color='#ffffff'>
        Websocket disconnected
        <Button onClick={retry}>Try reconnecting</Button>
      </Wrapper>
    )
  );
}
