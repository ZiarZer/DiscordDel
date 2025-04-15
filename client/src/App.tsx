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
  width: fit-content;
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
const guildSectionInfoFields = [{ label: 'ID', fieldName: 'id' }];

const displayUsername = (user: User|null) =>
  user?.discriminator === '0' ? `@${user?.username}` : `${user?.username}#${user?.discriminator}`;

const getUserAvatarUrl = (user: User) =>
  `https://cdn.discordapp.com/avatars/${user.id}/${user.avatar}.png`;
const getGuildIconUrl = (guild: Guild) =>
  `https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png`;

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
  const [loadedGuild, setLoadedGuild] = useState<Guild>(null);

  const userStatusMessage = useMemo(
    () =>
      currentUser != null
        ? `Successfully logged in as ${displayUsername(currentUser)}`
        : 'Not logged in',
    [currentUser]
  );
  const guildStatusMessage = useMemo(
    () =>
      loadedGuild != null
        ? `Successfully loaded guild ${loadedGuild.name}`
        : 'No guild loaded',
    [loadedGuild]
  );

  useEffect(() => {
    if (lastJsonMessage?.type === 'LOGIN') {
      setCurrentUser(lastJsonMessage.body);
    }else if (lastJsonMessage?.type === 'GET_GUILD') {
      setLoadedGuild(lastJsonMessage.body);
    }
  }, [lastJsonMessage, lastMessage]);

  const [authorizationToken, setAuthorizationToken] = useState('');
  const [inputGuildId, setInputGuildId] = useState('');

  const sendLoginRequest = useCallback(
    () =>
      sendJsonMessage({
        type: 'LOGIN',
        body: { authorizationToken },
      }),
    [sendJsonMessage, authorizationToken],
  );

  const sendGetGuildRequest = useCallback(
    () =>
      sendJsonMessage({
        type: 'GET_GUILD',
        body: { authorizationToken, guildId: inputGuildId },
      }),
    [sendJsonMessage, authorizationToken, inputGuildId],
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
            onChange: (e: ChangeEvent) => setAuthorizationToken(e.target.value),
          }}
          statusMessage={userStatusMessage}
          currentObject={currentUser}
          infoFields={userSectionInfoFields}
          getAvatarUrl={getUserAvatarUrl}
        />
        <Section
          title="Guild"
          actionInputBar={{
            inputPlaceholder: 'Guild ID',
            buttonLabel: 'Load guild by ID',
            enabled: currentUser != null,
            onSubmit: sendGetGuildRequest,
            onChange: (e: ChangeEvent) => setInputGuildId(e.target.value),
          }}
          statusMessage={guildStatusMessage}
          currentObject={loadedGuild}
          infoFields={guildSectionInfoFields}
          getAvatarUrl={getGuildIconUrl}
        />
      </Wrapper>
      <Version>{version}</Version>
    </>
  );
}

export default App;
