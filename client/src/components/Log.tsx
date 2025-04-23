import styled from 'styled-components';
import { LogEntry } from '../types';

const TAG_STYLES = {
  DEBUG: { color: 'white' },
  INFO: { color: 'black', backgroundColor: '#6a7ec7' },
  SUCCESS: { color: 'black', backgroundColor: '#85b32b' },
  WARNING: { color: 'black', backgroundColor: 'yellow' },
  ERROR: { color: 'white', backgroundColor: '#c3265e' },
  FATAL: { color: 'white', backgroundColor: '#676767' },
};

const LogTag = styled.span`
  font-weight: 800;
`;

const Wrapper = styled.span`
  text-align: left;
  font-size: 0.8em;
`;

function padCenter(str: string, length: number, char: string) {
  return str
    .padStart(Math.floor((length + str.length) / 2), char)
    .padEnd(length, char);
}

const NBSP = ' ';

export function Log({ message, logLevel = null }: LogEntry) {
  return (
    <Wrapper>
      {logLevel == null ? (
        NBSP.repeat(9)
      ) : (
        <LogTag>
          [
          <span style={TAG_STYLES[logLevel]}>
            {padCenter(logLevel, 7, NBSP)}
          </span>
          ]
        </LogTag>
      )}{' '}
      {message}
    </Wrapper>
  );
}
