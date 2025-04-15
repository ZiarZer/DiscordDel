import { useCallback, useEffect, useMemo, useState } from 'react';
import useWebSocket from 'react-use-websocket';
import styled from 'styled-components';

import { version } from '../package.json';
import { Section } from './components/Section';
import { CHANNEL_TYPES } from './constants';

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
const channelSectionInfoFields = [
  { label: "ID", fieldName: "id" },
  { label: "Last message ID", fieldName: "last_message_id" },
  {
    label: "Type",
    fieldName: "type",
    display: (v: keyof typeof CHANNEL_TYPES) => CHANNEL_TYPES[v],
  },
  { label: "Parent ID", fieldName: "parent_id" },
  { label: "Guild ID", fieldName: "guild_id" },
  { label: "Message count", fieldName: "message_count" },
];

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
  const [loadedChannel, setLoadedChannel] = useState<Channel>(null);

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
  const channelStatusMessage = useMemo(
    () =>
      loadedChannel != null
        ? `Successfully loaded channel ${loadedChannel.name}`
        : 'No channel loaded',
    [loadedChannel]
  );

  useEffect(() => {
    if (lastJsonMessage?.type === 'LOGIN') {
      setCurrentUser(lastJsonMessage.body);
    } else if (lastJsonMessage?.type === 'GET_GUILD') {
      setLoadedGuild(lastJsonMessage.body);
    } else if (lastJsonMessage?.type === 'GET_CHANNEL') {
      setLoadedChannel(lastJsonMessage.body);
    }
  }, [lastJsonMessage, lastMessage]);

  const [authorizationToken, setAuthorizationToken] = useState('');
  const [inputGuildId, setInputGuildId] = useState('');
  const [inputChannelId, setInputChannelId] = useState('');

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
  const sendGetChannelRequest = useCallback(
    () =>
      sendJsonMessage({
        type: 'GET_CHANNEL',
        body: { authorizationToken, channelId: inputChannelId },
      }),
    [sendJsonMessage, authorizationToken, inputChannelId],
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
        <Section
          title="Channel"
          actionInputBar={{
            inputPlaceholder: 'Channel ID',
            buttonLabel: 'Load channel by ID',
            enabled: currentUser != null,
            onSubmit: sendGetChannelRequest,
            onChange: (e: ChangeEvent) => setInputChannelId(e.target.value),
          }}
          statusMessage={channelStatusMessage}
          currentObject={loadedChannel}
          infoFields={channelSectionInfoFields}
          getAvatarUrl={() => null}
        />
      </Wrapper>
      <Version>{version}</Version>
    </>
  );
}

export default App;
