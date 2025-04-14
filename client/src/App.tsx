import { useCallback, useEffect, useMemo, useState } from 'react';
import useWebSocket from 'react-use-websocket';
import styled from 'styled-components';

import { version } from '../package.json';
import { Section } from './components/Section';

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

const userSectionInfoFields = [
  { label: 'ID', fieldName: 'id' },
  { label: 'Display name', fieldName: 'global_name' },
];

const displayUsername = (user: User|null) =>
  user?.discriminator === '0' ? `@${user?.username}` : `${user?.username}#${user?.discriminator}`;

const getUserAvatarUrl = (user: User) =>
  `https://cdn.discordapp.com/avatars/${user.id}/${user.avatar}.png`;

function App() {
  const { sendJsonMessage, lastMessage, lastJsonMessage } = useWebSocket(WEBSOCKET_URL, {
    share: false,
    shouldReconnect: () => true,
  }) as {
    sendJsonMessage: any
    lastMessage: any
    lastJsonMessage: null | { body: object; type: string }
  };

  const [currentUser, setCurrentUser] = useState<User>(null);

  const userStatusMessage = useMemo(
    () =>
      currentUser != null
        ? `Successfully logged in as ${displayUsername(currentUser)}`
        : 'Not logged in',
    [currentUser]
  );

  useEffect(() => {
    if (lastJsonMessage?.type === 'LOGIN') {
      setCurrentUser(lastJsonMessage.body);
    }
  }, [lastJsonMessage, lastMessage]);

  const sendLoginRequest = useCallback(
    (authorizationToken: string) =>
      sendJsonMessage({
        type: 'LOGIN',
        body: { authorizationToken },
      }),
    [sendJsonMessage],
  );

  return (
    <>
      <Wrapper>
        <Section
          title="User"
          actionInputBar={{
            inputPlaceholder: 'Authorization token',
            buttonLabel: 'Authenticate',
            enabled: true,
            secret: true,
            onSubmit: sendLoginRequest,
          }}
          statusMessage={userStatusMessage}
          currentObject={currentUser}
          infoFields={userSectionInfoFields}
          getAvatarUrl={getUserAvatarUrl}
        />
      </Wrapper>
      <Version>{version}</Version>
    </>
  );
}

export default App;
