import styled from 'styled-components';
import { Log } from './Log';
import { LogEntry } from '../types';
import { useMemo } from 'react';

const Wrapper = styled.div<{ $backgroundcolor: string }>`
  background-color: ${({ $backgroundcolor }) => $backgroundcolor};
  color: white;
  border-radius: 1em;
  padding: 1em;
  text-align: center;
  display: flex;
  flex-direction: column;
  width: 100%;
  font-family: monospace;
  transition: transform ease-in-out 150ms, background-color ease-in-out 150ms;
  &:active {
    transform: scale(0.95);
  }
  cursor: pointer;
`;

const SectionTitle = styled.h2`
  margin: 0;
`;

export function BigActionControlButton({
  actionTitle,
  onClick,
}: {
  actionTitle?: string;
  onClick: () => void;
}) {
  const backgroundColor = useMemo(
    () => (actionTitle == null ? '#808080' : '#c23a22'),
    [actionTitle]
  );
  const subtitle = useMemo(
    () =>
      actionTitle == null
        ? 'No action running'
        : `Action running: ${actionTitle}`,
    [actionTitle]
  );
  const mainText = useMemo(
    () => (actionTitle == null ? '_' : 'Click to STOP action'),
    [actionTitle]
  );

  return (
    <Wrapper $backgroundcolor={backgroundColor} onClick={onClick}>
      {subtitle}
      <SectionTitle>{mainText}</SectionTitle>
    </Wrapper>
  );
}
