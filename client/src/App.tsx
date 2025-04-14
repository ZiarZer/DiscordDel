import { useCallback, useEffect, useState } from 'react';
import useWebSocket from 'react-use-websocket';
import styled from 'styled-components';

import { version } from '../package.json';
import { ActionInputBar } from './components/ActionInputBar'
import { StatusMessage } from './components/StatusMessage';
import { InfoList } from './components/InfoList';

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

const SectionWrapper = styled.div`
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

const WEBSOCKET_URL = 'ws://127.0.0.1:8765';

const userSectionInfoFields = [
  { label: 'ID', fieldName: 'id' },
  { label: 'Display name', fieldName: 'global_name' },
];

const displayUsername = (user: User|null) =>
  user?.discriminator === '0' ? `@${user?.username}` : `${user?.username}#${user?.discriminator}`;

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
        <SectionWrapper>
          <SectionTitle>User</SectionTitle>
          <ActionInputBar
            inputPlaceholder='Authorization token'
            buttonText='Authenticate'
            enabled={true}
            secret
            onSubmit={sendLoginRequest}
          />
          <StatusMessage
            message={
              currentUser != null
                ? `Successfully logged in as ${displayUsername(currentUser)}`
                : 'Not logged in'
            }
            success={currentUser != null}
          />
          <InfoList
            currentObject={currentUser}
            fields={userSectionInfoFields}
            getAvatarUrl={(user) =>
              `https://cdn.discordapp.com/avatars/${user.id}/${user.avatar}.png`
            }
          />
        </SectionWrapper>
      </Wrapper>
      <Version>{version}</Version>
    </>
  );
}

export default App;
