import { ChangeEvent, useCallback, useEffect, useMemo, useState } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import styled from 'styled-components';

import { version } from '../package.json';
import { Console } from './components/Console';
import { PaginatedList } from './components/PaginatedList';
import { Section } from './components/Section';
import { CHANNEL_TYPES } from './constants';
import { Channel, Guild, InfoListFieldConfig, LogEntry, User } from './types';

const Wrapper = styled.div`
  display: flex;
  gap: 1em;
`;

const LeftPanel = styled.div`
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;
  margin: auto 1em;
  gap: 1em;
  width: fit-content;
  width: 33%;
`;

const Version = styled.p`
  position: absolute;
  right: 1em;
  bottom: 0;
`;

const WEBSOCKET_URL = 'ws://127.0.0.1:8765';

const userSectionInfoFields: Array<InfoListFieldConfig<User>> = [
  { label: 'ID', fieldName: 'id' },
  { label: 'Display name', fieldName: 'global_name' },
];
const guildSectionInfoFields: Array<InfoListFieldConfig<Guild>> = [
  { label: 'ID', fieldName: 'id' },
];
const channelSectionInfoFields: Array<InfoListFieldConfig<Channel>> = [
  { label: "ID", fieldName: "id" },
  { label: "Last message ID", fieldName: "last_message_id" },
  {
    label: "Type",
    fieldName: "type",
    display: (v?: number | string) =>
      CHANNEL_TYPES[v as keyof typeof CHANNEL_TYPES],
  },
  { label: "Parent ID", fieldName: "parent_id" },
  { label: "Guild ID", fieldName: "guild_id" },
  { label: "Message count", fieldName: "message_count" },
];

const displayUsername = (user: User | null) =>
  user?.discriminator === '0'
    ? `@${user?.username}`
    : `${user?.username}#${user?.discriminator}`;

const getUserAvatarUrl = (user: User | null) =>
  user
    ? `https://cdn.discordapp.com/avatars/${user.id}/${user.avatar}.png`
    : undefined;
const getGuildIconUrl = (guild: Guild | null) =>
  guild
    ? `https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png`
    : undefined;

function App() {
  const { sendJsonMessage, lastMessage, lastJsonMessage, readyState } =
    useWebSocket(WEBSOCKET_URL, {
      share: false,
      shouldReconnect: () => true,
    }) as {
      sendJsonMessage: (message: object) => void;
      lastMessage: object;
      lastJsonMessage: null | { body: object; type: string };
      readyState: ReadyState;
    };

  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [loadedGuild, setLoadedGuild] = useState<Guild | null>(null);
  const [loadedChannel, setLoadedChannel] = useState<Channel | null>(null);

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

  const [logs, setLogs] = useState<Array<LogEntry>>([]);
  const addLog = useCallback(({ logLevel, message }: LogEntry) => {
    setLogs((currentLogs) =>
      currentLogs.length >= 50
        ? currentLogs.slice(1).concat([{ logLevel, message }])
        : [...currentLogs, { logLevel, message }]
    );
  }, []);

  useEffect(() => {
    if (lastJsonMessage?.type === 'LOGIN') {
      setCurrentUser(lastJsonMessage.body as User);
    } else if (lastJsonMessage?.type === 'GET_GUILD') {
      setLoadedGuild(lastJsonMessage.body as Guild);
    } else if (lastJsonMessage?.type === 'GET_CHANNEL') {
      setLoadedChannel(lastJsonMessage.body as Channel);
    } else if (
      lastJsonMessage?.type === 'GET_USER_GUILDS' ||
      lastJsonMessage?.type === 'GET_GUILD_CHANNELS'
    ) {
      setResultsList(lastJsonMessage.body as Array<Guild> | Array<Channel>);
      setIsChannelType(lastJsonMessage.type === 'GET_GUILD_CHANNELS');
    } else if (lastJsonMessage?.type === 'LOG') {
      addLog(lastJsonMessage.body as LogEntry);
    }
  }, [lastJsonMessage, lastMessage, addLog]);

  const [authorizationToken, setAuthorizationToken] = useState('');
  const [inputGuildId, setInputGuildId] = useState('');
  const [inputChannelId, setInputChannelId] = useState('');
  const [resultsList, setResultsList] = useState<
    Array<Guild> | Array<Channel> | null
  >(null);
  const [isChannelType, setIsChannelType] = useState(false);

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
  const sendGetUserGuildsRequest = useCallback(
    () =>
      sendJsonMessage({
        type: 'GET_USER_GUILDS',
        body: { authorizationToken },
      }),
    [sendJsonMessage, authorizationToken]
  );
  const sendGetGuildChannelsRequest = useCallback(
    () =>
      sendJsonMessage({
        type: 'GET_GUILD_CHANNELS',
        body: { authorizationToken, guildId: loadedGuild?.id },
      }),
    [sendJsonMessage, authorizationToken, loadedGuild]
  );
  const sendCrawlChannelRequest = useCallback(
    () =>
      sendJsonMessage({
        type: 'CRAWL_CHANNEL',
        body: {
          authorizationToken,
          channelId: loadedChannel?.id,
          authorIds: [currentUser?.id],
        },
      }),
    [sendJsonMessage, authorizationToken, loadedChannel, currentUser]
  );
  const sendCrawlGuildRequest = useCallback(
    () =>
      sendJsonMessage({
        type: 'CRAWL_GUILD',
        body: {
          authorizationToken,
          guildId: loadedGuild?.id,
          authorIds: [currentUser?.id],
        },
      }),
    [sendJsonMessage, authorizationToken, loadedGuild, currentUser]
  );
  const sendCrawlAllGuildsRequest = useCallback(
    () =>
      sendJsonMessage({
        type: 'CRAWL_ALL_GUILDS',
        body: {
          authorizationToken,
          authorIds: [currentUser?.id],
        },
      }),
    [sendJsonMessage, authorizationToken, currentUser]
  );

  const userSectionActions = [
    { label: 'Get user guilds', onClick: sendGetUserGuildsRequest },
    { label: 'Crawl all guilds', onClick: sendCrawlAllGuildsRequest },
  ];
  const guildSectionActions = [
    { label: 'Get guild channels', onClick: sendGetGuildChannelsRequest },
    { label: 'Crawl guild', onClick: sendCrawlGuildRequest },
  ];
  const channelSectionActions = [
    { label: 'Crawl channel', onClick: sendCrawlChannelRequest },
  ];

  return (
    <>
      <Wrapper>
        <LeftPanel>
          <Section
            title="User"
            actionInputBar={{
              inputPlaceholder: 'Authorization token',
              buttonLabel: 'Authenticate',
              enabled: readyState === ReadyState.OPEN,
              secret: true,
              onSubmit: sendLoginRequest,
              onChange: (e: ChangeEvent) =>
                setAuthorizationToken((e.target as HTMLInputElement).value),
            }}
            statusMessage={userStatusMessage}
            currentObject={currentUser}
            infoFields={userSectionInfoFields}
            getAvatarUrl={getUserAvatarUrl}
            actions={userSectionActions}
          />
          <Section
            title="Guild"
            actionInputBar={{
              inputPlaceholder: 'Guild ID',
              buttonLabel: 'Load guild by ID',
              enabled: currentUser != null,
              onSubmit: sendGetGuildRequest,
              onChange: (e: ChangeEvent) =>
                setInputGuildId((e.target as HTMLInputElement).value),
            }}
            statusMessage={guildStatusMessage}
            currentObject={loadedGuild}
            infoFields={guildSectionInfoFields}
            getAvatarUrl={getGuildIconUrl}
            actions={guildSectionActions}
          />
          <Section
            title="Channel"
            actionInputBar={{
              inputPlaceholder: 'Channel ID',
              buttonLabel: 'Load channel by ID',
              enabled: currentUser != null,
              onSubmit: sendGetChannelRequest,
              onChange: (e: ChangeEvent) =>
                setInputChannelId((e.target as HTMLInputElement).value),
            }}
            statusMessage={channelStatusMessage}
            currentObject={loadedChannel}
            infoFields={channelSectionInfoFields}
            actions={channelSectionActions}
          />
        </LeftPanel>
        <PaginatedList
          resultsList={resultsList}
          isChannelType={isChannelType}
        />
        <Console logs={logs} />
      </Wrapper>
      <Version>{version}</Version>
    </>
  );
}

export default App;
