import { useMemo } from 'react';
import useWebSocket from 'react-use-websocket';
import styled from 'styled-components';

import { Button } from './components';

const Wrapper = styled.div`
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;
  padding: 1em;
  gap: 1em;
`;

const WEBSOCKET_URL = 'ws://127.0.0.1:8765';

function App() {
  const { sendJsonMessage, lastMessage, lastJsonMessage } = useWebSocket(
    WEBSOCKET_URL,
    { share: false, shouldReconnect: () => true },
  );
  const messageReceivedAt = useMemo(
    () => (lastMessage ? new Date(Date.now()).toISOString() : '---'),
    [lastMessage],
  );

  return (
    <Wrapper>
      DiscordDel
      <Button onClick={() => sendJsonMessage('test')}>Send WS message</Button>
      <span>{`Received message: ${lastJsonMessage}`}</span>
      <span>{messageReceivedAt}</span>
    </Wrapper>
  );
}

export default App;
