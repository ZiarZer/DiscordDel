import { useCallback, useEffect, useState } from 'react';
import useWebSocket from 'react-use-websocket';
import styled from 'styled-components';

import { version } from '../package.json';
import { ActionInputBar } from './components/ActionInputBar'

const Wrapper = styled.div`
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;
  padding: 1em;
  gap: 1em;
`;

const Version = styled.p`
  position: absolute;
  right: 1em;
  bottom: 0;
`;

const WEBSOCKET_URL = 'ws://127.0.0.1:8765';

function App() {
  const { sendJsonMessage, lastMessage, lastJsonMessage } = useWebSocket(WEBSOCKET_URL, {
    share: false,
    shouldReconnect: () => true,
  }) as {
    sendJsonMessage: any
    lastMessage: any
    lastJsonMessage: null | { body: object; type: string }
  };

  const [authorizationToken, setAuthorizationToken] = useState<string>('');
  useEffect(() => {
    if (lastJsonMessage?.type === 'LOGIN') {
      console.log(lastJsonMessage.body);
    }
  }, [lastJsonMessage, lastMessage]);

  const sendLoginRequest = useCallback(
    () =>
      sendJsonMessage({
        type: 'LOGIN',
        body: { authorizationToken },
      }),
    [authorizationToken, sendJsonMessage],
  );

  return (
    <>
      <Wrapper>
        DiscordDel
        <ActionInputBar
          inputPlaceholder='Authorization token'
          buttonText='Authenticate'
          enabled={true}
          secret
          onEdit={(e) => setAuthorizationToken(e?.target.value)}
          onSubmit={sendLoginRequest}
        />
      </Wrapper>
      <Version>{version}</Version>
    </>
  );
}

export default App;
