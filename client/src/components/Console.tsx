import styled from 'styled-components';
import { Log } from './Log';
import { LogEntry } from '../types';

const Wrapper = styled.div`
  border-radius: 1em;
  padding: 1em;
  text-align: center;
  display: flex;
  flex-direction: column;
  width: 33%;
  font-family: monospace;
`;

const SectionTitle = styled.h3`
  margin: 0;
`;

type ConsoleProps = {
  logs: Array<LogEntry>;
};

export function Console({ logs }: ConsoleProps) {
  return (
    <Wrapper>
      <SectionTitle>Logs</SectionTitle>
      {logs.map(({ logLevel, message }, index) => (
        <Log message={message} logLevel={logLevel} key={index} />
      ))}
    </Wrapper>
  );
}
